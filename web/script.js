document.addEventListener('DOMContentLoaded', () => {
    const bpmValueSpan = document.getElementById('bpm-value');
    const bpmDownButton = document.getElementById('bpm-down');
    const bpmUpButton = document.getElementById('bpm-up');
    const chainsContainer = document.getElementById('chains-container');

    let currentBPM = 0;
    let currentChains = [];

    const renderChains = (chains) => {
        chainsContainer.innerHTML = ''; // Clear existing chains
        chains.forEach(chain => {
            const chainBox = document.createElement('div');
            chainBox.className = 'chain-box';
            chainBox.innerHTML = `
                <h2>${chain.ID}</h2>
                <p>Priority: ${chain.Priority}</p>
                <p>Tick Rate: ${chain.TickRate} FPS</p>
                <p>Num Lamps: ${chain.NumLamps}</p>
                <h3>Output:</h3>
                <p>Type: ${chain.Output.Type}</p>
                <p>Channel Mapping: ${chain.Output.ChannelMapping}</p>
                <p>Channels per Lamp: ${chain.Output.NumChannelsPerLamp}</p>
                <h3>Effects:</h3>
                <ul>
                    ${chain.Effects.map(effect => `<li>${effect.Type}</li>`).join('')}
                </ul>
            `;
            chainsContainer.appendChild(chainBox);
        });
    };

    const fetchBPM = () => {
        fetch('/api/bpm')
            .then(response => response.json())
            .then(data => {
                if (currentBPM !== data.bpm) {
                    currentBPM = data.bpm;
                    bpmValueSpan.textContent = currentBPM.toFixed(2);
                }
            })
            .catch(error => console.error('Error fetching BPM:', error));
    };

    const fetchChains = () => {
        fetch('/api/chains')
            .then(response => response.json())
            .then(chains => {
                // Simple deep comparison for now. For complex objects, a library might be better.
                if (JSON.stringify(currentChains) !== JSON.stringify(chains)) {
                    currentChains = chains;
                    renderChains(currentChains);
                }
            })
            .catch(error => console.error('Error fetching chains:', error));
    };

    const updateBPM = (newBPM) => {
        fetch('/api/bpm', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ bpm: newBPM }),
        })
        .then(response => response.json())
        .then(data => {
            // Update immediately after successful POST, then polling will keep it consistent
            currentBPM = data.bpm;
            bpmValueSpan.textContent = currentBPM.toFixed(2);
        })
        .catch(error => console.error('Error updating BPM:', error));
    };

    bpmDownButton.addEventListener('click', () => {
        updateBPM(currentBPM - 5);
    });

    bpmUpButton.addEventListener('click', () => {
        updateBPM(currentBPM + 5);
    });

    // Initial fetches
    fetchBPM();
    fetchChains();

    // Poll every second
    setInterval(fetchBPM, 1000);
    setInterval(fetchChains, 1000);
});