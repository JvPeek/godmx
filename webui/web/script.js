document.addEventListener('DOMContentLoaded', () => {
    const bpmValueSpan = document.getElementById('bpm-value');
    const bpmDownButton = document.getElementById('bpm-down');
    const bpmUpButton = document.getElementById('bpm-up');
    const chainsContainer = document.getElementById('chains-container');
    const eventsContainer = document.getElementById('events-container'); // New: Get events container

    let currentBPM = 0;
    let currentChains = [];
    let currentEvents = []; // New: To track current events

    const renderArgs = (args) => {
        if (!args || Object.keys(args).length === 0) {
            return '';
        }
        return `
            <ul class="args-list">
                ${Object.entries(args).map(([key, value]) => `
                    <li><strong>${key}:</strong> ${JSON.stringify(value)}</li>
                `).join('')}
            </ul>
        `;
    };

    const renderChains = (chains) => {
        chainsContainer.innerHTML = ''; // Clear existing chains
        chains.forEach(chain => {
            const chainBox = document.createElement('div');
            chainBox.className = 'chain-box';
            chainBox.innerHTML = `
                <h2>${chain.ID}</h2>
                <p>Priority: ${chain.Priority}</p>
                
                <h3>Chain Flow:</h3>
                <div class="chain-flow">
                    <div class="chain-element">
                        <h4>Tick</h4>
                        <p><strong>Rate:</strong> ${chain.TickRate} FPS</p>
                        <p><strong>Lamps:</strong> ${chain.NumLamps}</p>
                    </div>
                    
                    ${chain.Effects.map(effect => `
                        <div class="chain-element ${effect.Enabled ? '' : 'disabled'}">
                            <h4>Effect: ${effect.Type}</h4>
                            ${renderArgs(effect.Args)}
                        </div>
                    `).join('')}

                    <div class="chain-element">
                        <h4>Output: ${chain.Output.Type}</h4>
                        <p><strong>Channel Mapping:</strong> ${chain.Output.ChannelMapping}</p>
                        <p><strong>Channels per Lamp:</strong> ${chain.Output.NumChannelsPerLamp}</p>
                        ${renderArgs(chain.Output.Args)}
                    </div>
                </div>
            `;
            chainsContainer.appendChild(chainBox);
        });
    };

    const fetchBPM = async () => {
        try {
            const response = await fetch('/api/bpm');
            const data = await response.json();
            if (currentBPM !== data.bpm) {
                currentBPM = data.bpm;
                bpmValueSpan.textContent = currentBPM.toFixed(2);
            }
        } catch (error) {
            console.error('Error fetching BPM:', error);
        }
    };

    const fetchChains = async () => {
        try {
            const response = await fetch('/api/chains');
            const chains = await response.json();
            if (JSON.stringify(currentChains) !== JSON.stringify(chains)) {
                currentChains = chains;
                renderChains(currentChains);
            }
        } catch (error) {
            console.error('Error fetching chains:', error);
        }
    };

    // New: Function to fetch events and render buttons
    const fetchEventsAndRenderButtons = async () => {
        try {
            const response = await fetch('/api/events');
            const events = await response.json();
            // Only re-render if events have changed
            if (JSON.stringify(currentEvents) !== JSON.stringify(events)) {
                currentEvents = events;
                eventsContainer.innerHTML = ''; // Clear existing buttons
                events.forEach(eventName => {
                    const button = document.createElement('button');
                    button.textContent = eventName.replace(/_/g, ' '); // Make it more readable
                    button.className = 'event-button';
                    button.addEventListener('click', async () => {
                        try {
                            const triggerResponse = await fetch('/api/trigger', {
                                method: 'POST',
                                headers: {
                                    'Content-Type': 'application/json',
                                },
                                body: JSON.stringify({ eventName: eventName }),
                            });
                            const result = await triggerResponse.json();
                            console.log(`Event '${eventName}' triggered:`, result);
                            // After triggering, refresh all data to see changes
                            refreshAll();
                        } catch (error) {
                            console.error(`Error triggering event '${eventName}':`, error);
                        }
                    });
                    eventsContainer.appendChild(button);
                });
            }
        } catch (error) {
            console.error('Error fetching events:', error);
        }
    };

    const updateBPM = async (newBPM) => {
        try {
            const response = await fetch('/api/bpm', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ bpm: newBPM }),
            });
            const data = await response.json();
            currentBPM = data.bpm;
            bpmValueSpan.textContent = currentBPM.toFixed(2);
        } catch (error) {
            console.error('Error updating BPM:', error);
        }
    };

    bpmDownButton.addEventListener('click', () => {
        updateBPM(currentBPM - 5);
    });

    bpmUpButton.addEventListener('click', () => {
        updateBPM(currentBPM + 5);
    });

    // New: Function to refresh all data
    const refreshAll = () => {
        fetchBPM();
        fetchChains();
        fetchEventsAndRenderButtons();
    };

    // Initial fetch and poll every second
    refreshAll();
    setInterval(refreshAll, 1000);
});