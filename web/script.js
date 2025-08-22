document.addEventListener('DOMContentLoaded', () => {
    const bpmValueSpan = document.getElementById('bpm-value');
    const bpmDownButton = document.getElementById('bpm-down');
    const bpmUpButton = document.getElementById('bpm-up');

    let currentBPM = 0;

    const fetchBPM = () => {
        fetch('/api/bpm')
            .then(response => response.json())
            .then(data => {
                currentBPM = data.bpm;
                bpmValueSpan.textContent = currentBPM.toFixed(2);
            })
            .catch(error => console.error('Error fetching BPM:', error));
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

    // Initial fetch of BPM
    fetchBPM();

    fetch('/api/chains')
        .then(response => response.json())
        .then(chains => {
            const container = document.getElementById('chains-container');
            chains.forEach(chain => {
                const chainBox = document.createElement('div');
                chainBox.className = 'chain-box';
                chainBox.innerHTML = `
                    <h2>Chain ID: ${chain.ID}</h2>
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
                container.appendChild(chainBox);
            });
        })
        .catch(error => console.error('Error fetching chains:', error));
});