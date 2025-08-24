// web/js/event_action_editor.js

document.addEventListener('DOMContentLoaded', () => {
    const eventsContainer = document.getElementById('eventsContainer');
    const addEventButton = document.getElementById('addEvent');
    const saveConfigButton = document.getElementById('saveConfig');

    let currentConfig = {}; // In-memory representation of the config
    let actionSchemas = {}; // Stores action schemas fetched from backend

    // --- API Calls ---

    async function fetchActionSchemas() {
        try {
            const response = await fetch('/api/actions/schema');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            actionSchemas = await response.json();
            console.log('Fetched Action Schemas:', actionSchemas);
        } catch (error) {
            console.error('Error fetching action schemas:', error);
            alert('Failed to load action schemas. Check console for details.');
        }
    }

    async function fetchConfig() {
        try {
            const response = await fetch('/api/config');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            currentConfig = await response.json();
            console.log('Fetched Config:', currentConfig);
            renderConfig();
        } catch (error) {
            console.error('Error fetching config:', error);
            alert('Failed to load configuration. Check console for details.');
        }
    }

    async function saveConfig() {
        try {
            const response = await fetch('/api/config', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(currentConfig)
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const result = await response.json();
            console.log('Save Config Result:', result);
            alert(result.message || 'Configuration saved successfully!');
        } catch (error) {
            console.error('Error saving config:', error);
            alert('Failed to save configuration. Check console for details.');
        }
    }

    // --- Rendering Functions ---

    function renderConfig() {
        eventsContainer.innerHTML = ''; // Clear existing events
        for (const eventName in currentConfig.events) {
            renderEvent(eventName, currentConfig.events[eventName]);
        }
    }

    function renderEvent(eventName, actions) {
        const eventDiv = document.createElement('div');
        eventDiv.className = 'event-item';
        eventDiv.dataset.eventName = eventName;
        eventDiv.innerHTML = `
            <h3>Event: ${eventName}</h3>
            <button class="add-action" data-event-name="${eventName}">Add Action</button>
            <div class="actions-container" data-event-name="${eventName}">
                <!-- Actions will be rendered here -->
            </div>
            <button class="remove-event" data-event-name="${eventName}">Remove Event</button>
        `;
        eventsContainer.appendChild(eventDiv);

        const actionsContainer = eventDiv.querySelector('.actions-container');
        actions.forEach(action => renderAction(action, actionsContainer, eventName));

        // Add event listeners
        eventDiv.querySelector('.add-action').addEventListener('click', (e) => {
            const eventName = e.target.dataset.eventName;
            addActionToEvent(eventName);
        });
        eventDiv.querySelector('.remove-event').addEventListener('click', (e) => {
            const eventName = e.target.dataset.eventName;
            removeEvent(eventName);
        });
    }

    function renderAction(action, actionsContainer, eventName) {
        const actionDiv = document.createElement('div');
        actionDiv.className = 'action-item';
        actionDiv.dataset.eventName = eventName;
        const actionIndex = actionsContainer.children.length; 
        actionDiv.dataset.actionIndex = actionIndex;

        const h5Element = document.createElement('h5');
        h5Element.textContent = 'Action: ';

        const actionTypeSelect = document.createElement('select');
        actionTypeSelect.className = 'action-type';
        actionTypeSelect.dataset.eventName = eventName;
        actionTypeSelect.dataset.actionIndex = actionIndex;

        for (const type in actionSchemas) {
            const option = document.createElement('option');
            option.value = type;
            option.textContent = actionSchemas[type].human_readable_name || type;
            if (type === action.type) {
                option.selected = true;
            }
            actionTypeSelect.appendChild(option);
        }
        h5Element.appendChild(actionTypeSelect);
        actionDiv.appendChild(h5Element);

        const paramsContainer = document.createElement('div');
        paramsContainer.className = 'action-params-container';
        paramsContainer.dataset.eventName = eventName;
        paramsContainer.dataset.actionIndex = actionIndex;
        actionDiv.appendChild(paramsContainer);

        const removeButton = document.createElement('button');
        removeButton.className = 'remove-action';
        removeButton.dataset.eventName = eventName;
        removeButton.dataset.actionIndex = actionIndex;
        removeButton.textContent = 'Remove Action';
        actionDiv.appendChild(removeButton);

        actionsContainer.appendChild(actionDiv);

        renderActionParams(action, paramsContainer, eventName, actionIndex);

        // Add event listeners for action type change
        actionTypeSelect.addEventListener('change', (e) => {
            const eventName = e.target.dataset.eventName;
            const actionIndex = parseInt(e.target.dataset.actionIndex);
            const newType = e.target.value;
            updateActionType(eventName, actionIndex, newType);
        });

        // Add event listener for Remove Action button
        removeButton.addEventListener('click', (e) => {
            const eventName = e.target.dataset.eventName;
            const actionIndex = parseInt(e.target.dataset.actionIndex);
            removeAction(eventName, actionIndex);
        });
    }

    function renderActionParams(action, paramsContainer, eventName, actionIndex) {
        paramsContainer.innerHTML = ''; // Clear existing params
        const schema = actionSchemas[action.Type];
        if (!schema || !schema.Parameters) {
            paramsContainer.innerHTML = '<p>No parameters for this action type.</p>';
            return;
        }

        schema.Parameters.forEach(paramSchema => {
            const paramName = paramSchema.internal_name;
            const paramValue = action.params ? action.params[paramName] : paramSchema.default_value;
            const paramDiv = document.createElement('div');
            paramDiv.className = 'action-param-item';
            paramDiv.innerHTML = `<label>${paramSchema.display_name || paramName}:</label>`;

            let inputElement;
            switch (paramSchema.data_type) {
                case 'string':
                    inputElement = document.createElement('input');
                    inputElement.type = 'text';
                    inputElement.value = paramValue;
                    break;
                case 'float64':
                case 'int':
                    inputElement = document.createElement('input');
                    inputElement.type = 'number';
                    inputElement.value = paramValue;
                    // Add min/max if available
                    if (paramSchema.min_value !== undefined && paramSchema.min_value !== null) {
                        inputElement.min = paramSchema.min_value;
                    }
                    if (paramSchema.max_value !== undefined && paramSchema.max_value !== null) {
                        inputElement.max = paramSchema.max_value;
                    }
                    break;
                case 'bool':
                    inputElement = document.createElement('input');
                    inputElement.type = 'checkbox';
                    inputElement.checked = paramValue;
                    break;
                case 'object': // For complex types like EffectConfig in add_effect
                    // This will require more sophisticated rendering, possibly a nested form
                    inputElement = document.createElement('textarea');
                    inputElement.value = JSON.stringify(paramValue, null, 2);
                    inputElement.placeholder = "JSON object";
                    break;
                default:
                    inputElement = document.createElement('input');
                    inputElement.type = 'text';
                    inputElement.value = JSON.stringify(paramValue); // Fallback
                    break;
            }
            inputElement.dataset.eventName = eventName;
            inputElement.dataset.actionIndex = actionIndex;
            inputElement.dataset.paramName = paramName;
            inputElement.className = 'action-param-input';

            inputElement.addEventListener('change', (e) => {
                const eventName = e.target.dataset.eventName;
                const actionIndex = parseInt(e.target.dataset.actionIndex);
                const paramName = e.target.dataset.paramName;
                let value;
                if (e.target.type === 'checkbox') {
                    value = e.target.checked;
                } else if (e.target.type === 'number') {
                    value = parseFloat(e.target.value);
                } else if (paramSchema.data_type === 'object') {
                    try {
                        value = JSON.parse(e.target.value);
                    } catch (err) {
                        console.error('Invalid JSON for object parameter:', err);
                        alert('Invalid JSON input. Please correct it.');
                        return;
                    }
                } else {
                    value = e.target.value;
                }
                updateActionParam(eventName, actionIndex, paramName, value);
            });
            paramDiv.appendChild(inputElement);
            paramsContainer.appendChild(paramDiv);
        });
    }

    // --- Data Manipulation Functions ---

    function updateActionParam(eventName, actionIndex, paramName, value) {
        const actions = currentConfig.events[eventName];
        if (actions && actions[actionIndex]) {
            if (!actions[actionIndex].params) {
                actions[actionIndex].params = {};
            }
            actions[actionIndex].params[paramName] = value;
            console.log(`Updated Event ${eventName} Action ${actionIndex} param ${paramName} to ${value}`);
        }
    }

    function updateActionType(eventName, actionIndex, newType) {
        const actions = currentConfig.events[eventName];
        if (actions && actions[actionIndex]) {
            actions[actionIndex].type = newType;
            actions[actionIndex].params = {}; // Clear params when type changes
            console.log(`Updated Event ${eventName} Action ${actionIndex} type to ${newType}`);
            // Re-render the specific action to update its params section
            const actionsContainer = document.querySelector(`.actions-container[data-event-name="${eventName}"]`);
            const oldActionDiv = actionsContainer.querySelector(`.action-item[data-action-index="${actionIndex}"]`);
            if (oldActionDiv) {
                actionsContainer.removeChild(oldActionDiv);
            }
            // Re-render with a new index to avoid conflicts if elements shift
            // This is a simplified re-render, a more robust solution might re-index all actions
            renderAction(actions[actionIndex], actionsContainer, eventName);
        }
    }

    function addEvent() {
        const newEventName = `newEvent${Object.keys(currentConfig.events).length + 1}`;
        currentConfig.events[newEventName] = [];
        renderEvent(newEventName, []);
        console.log(`Added new event: ${newEventName}`);
    }

    function removeEvent(eventName) {
        delete currentConfig.events[eventName];
        document.querySelector(`.event-item[data-event-name="${eventName}"]`).remove();
        console.log(`Removed event: ${eventName}`);
    }

    function addActionToEvent(eventName) {
        const actions = currentConfig.events[eventName];
        if (actions) {
            const newAction = {
                type: Object.keys(actionSchemas)[0] || 'set_global', // Default to first available or set_global
                params: {}
            };
            actions.push(newAction);
            const actionsContainer = document.querySelector(`.actions-container[data-event-name="${eventName}"]`);
            renderAction(newAction, actionsContainer, eventName);
            console.log(`Added new action to event ${eventName}`);
        }
    }

    function removeAction(eventName, actionIndex) {
        const actions = currentConfig.events[eventName];
        if (actions && actions[actionIndex]) {
            actions.splice(actionIndex, 1);
            // Re-render all actions for this event to update indices
            const actionsContainer = document.querySelector(`.actions-container[data-event-name="${eventName}"]`);
            actionsContainer.innerHTML = '';
            actions.forEach((action, idx) => renderAction(action, actionsContainer, eventName, idx));
            console.log(`Removed action ${actionIndex} from event ${eventName}`);
        }
    }

    // --- Event Listeners ---
    addEventButton.addEventListener('click', addEvent);
    saveConfigButton.addEventListener('click', saveConfig);

    // --- Initialization ---
    async function init() {
        await fetchActionSchemas();
        await fetchConfig();
    }

    init();
});
