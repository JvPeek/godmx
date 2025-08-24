// web/js/editor.js

document.addEventListener('DOMContentLoaded', () => {
    const chainsContainer = document.getElementById('chainsContainer');
    const addChainButton = document.getElementById('addChain');
    const saveConfigButton = document.getElementById('saveConfig');

    let currentConfig = {}; // In-memory representation of the config
    let effectSchemas = {}; // Stores effect schemas fetched from backend

    // --- API Calls ---

    async function fetchEffectSchemas() {
        try {
            const response = await fetch('/api/effects/schema');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            effectSchemas = await response.json();
            console.log('Fetched Effect Schemas:', effectSchemas);
        } catch (error) {
            console.error('Error fetching effect schemas:', error);
            alert('Failed to load effect schemas. Check console for details.');
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
        chainsContainer.innerHTML = ''; // Clear existing chains
        currentConfig.Chains.forEach(chain => renderChain(chain));
    }

    function renderChain(chain) {
        const chainDiv = document.createElement('div');
        chainDiv.className = 'chain-item';
        chainDiv.dataset.chainId = chain.ID;
        chainDiv.innerHTML = `
            <h3>Chain: ${chain.ID}</h3>
            <p>Priority: <input type="number" class="chain-priority" value="${chain.Priority}" data-chain-id="${chain.ID}"></p>
            <p>TickRate: <input type="number" class="chain-tickrate" value="${chain.TickRate}" data-chain-id="${chain.ID}"></p>
            <p>NumLamps: <input type="number" class="chain-numlamps" value="${chain.NumLamps}" data-chain-id="${chain.ID}"></p>
            <h4>Effects</h4>
            <button class="add-effect" data-chain-id="${chain.ID}">Add Effect</button>
            <div class="effects-container" data-chain-id="${chain.ID}">
                <!-- Effects will be rendered here -->
            </div>
            <button class="remove-chain" data-chain-id="${chain.ID}">Remove Chain</button>
        `;
        chainsContainer.appendChild(chainDiv);

        const effectsContainer = chainDiv.querySelector('.effects-container');
        chain.Effects.forEach(effect => renderEffect(effect, effectsContainer, chain.ID));

        // Add event listeners for chain properties
        chainDiv.querySelectorAll('input').forEach(input => {
            input.addEventListener('change', (e) => {
                const chainId = e.target.dataset.chainId;
                const prop = e.target.className.replace('chain-', '');
                const value = e.target.type === 'number' ? parseInt(e.target.value) : e.target.value;
                updateChainProperty(chainId, prop, value);
            });
        });

        // Add event listener for Add Effect button
        chainDiv.querySelector('.add-effect').addEventListener('click', (e) => {
            const chainId = e.target.dataset.chainId;
            addEffectToChain(chainId);
        });

        // Add event listener for Remove Chain button
        chainDiv.querySelector('.remove-chain').addEventListener('click', (e) => {
            const chainId = e.target.dataset.chainId;
            removeChain(chainId);
        });
    }

    function renderEffect(effect, effectsContainer, chainId) {
        const effectDiv = document.createElement('div');
        effectDiv.className = 'effect-item';
        effectDiv.dataset.effectId = effect.ID;
        effectDiv.dataset.chainId = chainId;

        const effectTypeSelect = document.createElement('select');
        effectTypeSelect.className = 'effect-type';
        effectTypeSelect.dataset.effectId = effect.ID;
        effectTypeSelect.dataset.chainId = chainId;

        // Populate effect types from schemas
        for (const type in effectSchemas) {
            const option = document.createElement('option');
            option.value = type;
            option.textContent = type;
            if (type === effect.Type) {
                option.selected = true;
            }
            effectTypeSelect.appendChild(option);
        }

        effectDiv.innerHTML = `
            <h5>Effect: ${effect.ID || 'New Effect'} (Type: <span class="effect-type-display">${effect.Type}</span>)</h5>
            <p>ID: <input type="text" class="effect-id" value="${effect.ID || ''}" data-effect-id="${effect.ID}" data-chain-id="${chainId}"></p>
            <p>Enabled: <input type="checkbox" class="effect-enabled" ${effect.Enabled ? 'checked' : ''} data-effect-id="${effect.ID}" data-chain-id="${chainId}"></p>
            <p>Group: <input type="text" class="effect-group" value="${effect.Group || ''}" data-effect-id="${effect.ID}" data-chain-id="${chainId}"></p>
            <h6>Arguments</h6>
            <div class="effect-args-container" data-effect-id="${effect.ID}" data-chain-id="${chainId}">
                <!-- Args will be rendered here -->
            </div>
            <button class="remove-effect" data-effect-id="${effect.ID}" data-chain-id="${chainId}">Remove Effect</button>
        `;
        effectsContainer.appendChild(effectDiv);

        // Replace the placeholder for effect type display with the actual select element
        effectDiv.querySelector('.effect-type-display').replaceWith(effectTypeSelect);

        const argsContainer = effectDiv.querySelector('.effect-args-container');
        renderEffectArgs(effect, argsContainer, chainId);

        // Add event listeners for effect properties
        effectDiv.querySelectorAll('input, select').forEach(input => {
            input.addEventListener('change', (e) => {
                const chainId = e.target.dataset.chainId;
                const effectId = e.target.dataset.effectId;
                const prop = e.target.className.replace('effect-', '');
                let value;
                if (e.target.type === 'checkbox') {
                    value = e.target.checked;
                } else if (e.target.type === 'number') {
                    value = parseFloat(e.target.value);
                } else {
                    value = e.target.value;
                }
                updateEffectProperty(chainId, effectId, prop, value);
            });
        });

        // Add event listener for Remove Effect button
        effectDiv.querySelector('.remove-effect').addEventListener('click', (e) => {
            const chainId = e.target.dataset.chainId;
            const effectId = e.target.dataset.effectId;
            removeEffect(chainId, effectId);
        });

        // Handle effect type change to re-render args
        effectTypeSelect.addEventListener('change', (e) => {
            const chainId = e.target.dataset.chainId;
            const effectId = e.target.dataset.effectId;
            const newType = e.target.value;
            updateEffectType(chainId, effectId, newType);
        });
    }

    function renderEffectArgs(effect, argsContainer, chainId) {
        argsContainer.innerHTML = ''; // Clear existing args
        const schema = effectSchemas[effect.Type];
        if (!schema || !schema.args) {
            argsContainer.innerHTML = '<p>No arguments for this effect type.</p>';
            return;
        }

        for (const argName in schema.args) {
            const argSchema = schema.args[argName];
            const argValue = effect.Args ? effect.Args[argName] : '';
            const argDiv = document.createElement('div');
            argDiv.className = 'effect-arg-item';
            argDiv.innerHTML = `<label>${argName}:</label>`;

            let inputElement;
            switch (argSchema.type) {
                case 'string':
                    inputElement = document.createElement('input');
                    inputElement.type = 'text';
                    inputElement.value = argValue;
                    break;
                case 'number':
                case 'integer':
                    inputElement = document.createElement('input');
                    inputElement.type = 'number';
                    inputElement.value = argValue;
                    break;
                case 'boolean':
                    inputElement = document.createElement('input');
                    inputElement.type = 'checkbox';
                    inputElement.checked = argValue;
                    break;
                default:
                    inputElement = document.createElement('input');
                    inputElement.type = 'text';
                    inputElement.value = JSON.stringify(argValue); // Fallback for complex types
                    break;
            }
            inputElement.dataset.chainId = chainId;
            inputElement.dataset.effectId = effect.ID;
            inputElement.dataset.argName = argName;
            inputElement.className = 'effect-arg-input';

            inputElement.addEventListener('change', (e) => {
                const chainId = e.target.dataset.chainId;
                const effectId = e.target.dataset.effectId;
                const argName = e.target.dataset.argName;
                let value;
                if (e.target.type === 'checkbox') {
                    value = e.target.checked;
                } else if (e.target.type === 'number') {
                    value = parseFloat(e.target.value);
                } else {
                    value = e.target.value;
                }
                updateEffectArg(chainId, effectId, argName, value);
            });
            argDiv.appendChild(inputElement);
            argsContainer.appendChild(argDiv);
        }
    }

    // --- Data Manipulation Functions ---

    function updateChainProperty(chainId, prop, value) {
        const chain = currentConfig.Chains.find(c => c.ID === chainId);
        if (chain) {
            chain[prop] = value;
            console.log(`Updated Chain ${chainId} property ${prop} to ${value}`);
        }
    }

    function updateEffectProperty(chainId, effectId, prop, value) {
        const chain = currentConfig.Chains.find(c => c.ID === chainId);
        if (chain) {
            const effect = chain.Effects.find(e => e.ID === effectId);
            if (effect) {
                effect[prop] = value;
                console.log(`Updated Effect ${effectId} property ${prop} to ${value} in Chain ${chainId}`);
            }
        }
    }

    function updateEffectArg(chainId, effectId, argName, value) {
        const chain = currentConfig.Chains.find(c => c.ID === chainId);
        if (chain) {
            const effect = chain.Effects.find(e => e.ID === effectId);
            if (effect) {
                if (!effect.Args) {
                    effect.Args = {};
                }
                effect.Args[argName] = value;
                console.log(`Updated Effect ${effectId} arg ${argName} to ${value} in Chain ${chainId}`);
            }
        }
    }

    function updateEffectType(chainId, effectId, newType) {
        const chain = currentConfig.Chains.find(c => c.ID === chainId);
        if (chain) {
            const effect = chain.Effects.find(e => e.ID === effectId);
            if (effect) {
                effect.Type = newType;
                effect.Args = {}; // Clear args when type changes
                console.log(`Updated Effect ${effectId} type to ${newType} in Chain ${chainId}`);
                // Re-render the specific effect to update its args section
                const effectsContainer = document.querySelector(`.effects-container[data-chain-id="${chainId}"]`);
                const oldEffectDiv = effectsContainer.querySelector(`.effect-item[data-effect-id="${effectId}"]`);
                if (oldEffectDiv) {
                    effectsContainer.removeChild(oldEffectDiv);
                }
                renderEffect(effect, effectsContainer, chainId);
            }
        }
    }

    function addChain() {
        const newChainId = `newChain${currentConfig.Chains.length + 1}`;
        const newChain = {
            ID: newChainId,
            Priority: 0,
            TickRate: 100,
            NumLamps: 1,
            Output: { Type: "artnet", Args: { "ip": "127.0.0.1" }, ChannelMapping: "RGB", NumChannelsPerLamp: 3 }, // Default output
            Effects: []
        };
        currentConfig.Chains.push(newChain);
        renderChain(newChain);
        console.log(`Added new chain: ${newChainId}`);
    }

    function removeChain(chainId) {
        currentConfig.Chains = currentConfig.Chains.filter(c => c.ID !== chainId);
        document.querySelector(`.chain-item[data-chain-id="${chainId}"]`).remove();
        console.log(`Removed chain: ${chainId}`);
    }

    function addEffectToChain(chainId) {
        const chain = currentConfig.Chains.find(c => c.ID === chainId);
        if (chain) {
            const newEffectId = `newEffect${chain.Effects.length + 1}`;
            const newEffect = {
                ID: newEffectId,
                Type: Object.keys(effectSchemas)[0] || 'solid_color', // Default to first available or solid_color
                Args: {},
                Enabled: true,
                Group: ''
            };
            chain.Effects.push(newEffect);
            const effectsContainer = document.querySelector(`.effects-container[data-chain-id="${chainId}"]`);
            renderEffect(newEffect, effectsContainer, chainId);
            console.log(`Added new effect ${newEffectId} to chain ${chainId}`);
        }
    }

    function removeEffect(chainId, effectId) {
        const chain = currentConfig.Chains.find(c => c.ID === chainId);
        if (chain) {
            chain.Effects = chain.Effects.filter(e => e.ID !== effectId);
            document.querySelector(`.effect-item[data-effect-id="${effectId}"][data-chain-id="${chainId}"]`).remove();
            console.log(`Removed effect ${effectId} from chain ${chainId}`);
        }
    }

    // --- Event Listeners ---
    addChainButton.addEventListener('click', addChain);
    saveConfigButton.addEventListener('click', saveConfig);

    // --- Initialization ---
    async function init() {
        await fetchEffectSchemas();
        await fetchConfig();
    }

    init();
});
