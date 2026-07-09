package inspector

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simulation Inspector</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #1a1a2e;
            color: #eee;
            padding: 20px;
            line-height: 1.6;
        }
        h1 {
            text-align: center;
            margin-bottom: 30px;
            color: #00d4ff;
            font-weight: 300;
            letter-spacing: 2px;
        }
        .dashboard {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            max-width: 1400px;
            margin: 0 auto;
        }
        .card {
            background: #16213e;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
        }
        .card h2 {
            color: #00d4ff;
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 15px;
            padding-bottom: 10px;
            border-bottom: 1px solid #0f3460;
        }
        .metric {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #0f3460;
        }
        .metric:last-child {
            border-bottom: none;
        }
        .metric .label {
            color: #888;
            font-size: 13px;
        }
        .metric .value {
            font-family: 'Courier New', monospace;
            font-weight: bold;
            color: #fff;
        }
        .metric .value.good {
            color: #00ff88;
        }
        .metric .value.warning {
            color: #ffaa00;
        }
        .metric .value.bad {
            color: #ff4444;
        }
        .status-indicator {
            display: inline-block;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            margin-right: 8px;
        }
        .status-indicator.active {
            background: #00ff88;
            box-shadow: 0 0 10px #00ff88;
        }
        .status-indicator.inactive {
            background: #ff4444;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #666;
            font-size: 12px;
        }
        .connection-status {
            text-align: center;
            margin-bottom: 20px;
        }
        .connection-status .status {
            display: inline-block;
            padding: 5px 15px;
            background: #0f3460;
            border-radius: 20px;
            font-size: 12px;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        .updating {
            animation: pulse 1s infinite;
        }
    </style>
</head>
<body>
    <h1>Simulation Inspector</h1>
    
    <div class="connection-status">
        <span id="connectionStatus" class="status">Connecting...</span>
    </div>

    <div class="dashboard">
        <!-- Clock Card -->
        <div class="card">
            <h2>Simulation Clock</h2>
            <div class="metric">
                <span class="label">Elapsed Time</span>
                <span class="value" id="clockElapsed">--</span>
            </div>
            <div class="metric">
                <span class="label">Tick Count</span>
                <span class="value" id="clockTicks">--</span>
            </div>
            <div class="metric">
                <span class="label">Mode</span>
                <span class="value" id="clockMode">--</span>
            </div>
            <div class="metric">
                <span class="label">Status</span>
                <span class="value" id="clockPaused">
                    <span class="status-indicator active"></span>Running
                </span>
            </div>
        </div>

        <!-- Sun Card -->
        <div class="card">
            <h2>Sun Model</h2>
            <div class="metric">
                <span class="label">Elevation</span>
                <span class="value" id="sunElevation">--</span>
            </div>
            <div class="metric">
                <span class="label">Azimuth</span>
                <span class="value" id="sunAzimuth">--</span>
            </div>
            <div class="metric">
                <span class="label">Irradiance (GHI)</span>
                <span class="value" id="sunIrradiance">--</span>
            </div>
            <div class="metric">
                <span class="label">Direct Normal (DNI)</span>
                <span class="value" id="sunDNI">--</span>
            </div>
            <div class="metric">
                <span class="label">Diffuse</span>
                <span class="value" id="sunDiffuse">--</span>
            </div>
            <div class="metric">
                <span class="label">Daytime</span>
                <span class="value" id="sunDaytime">
                    <span class="status-indicator inactive"></span>Night
                </span>
            </div>
        </div>

        <!-- Weather Card -->
        <div class="card">
            <h2>Weather Model</h2>
            <div class="metric">
                <span class="label">Temperature</span>
                <span class="value" id="weatherTemp">--</span>
            </div>
            <div class="metric">
                <span class="label">Humidity</span>
                <span class="value" id="weatherHumidity">--</span>
            </div>
            <div class="metric">
                <span class="label">Pressure</span>
                <span class="value" id="weatherPressure">--</span>
            </div>
            <div class="metric">
                <span class="label">Cloud Cover</span>
                <span class="value" id="weatherClouds">--</span>
            </div>
            <div class="metric">
                <span class="label">Wind Speed</span>
                <span class="value" id="weatherWind">--</span>
            </div>
            <div class="metric">
                <span class="label">Wind Direction</span>
                <span class="value" id="weatherWindDir">--</span>
            </div>
            <div class="metric">
                <span class="label">Raining</span>
                <span class="value" id="weatherRain">
                    <span class="status-indicator inactive"></span>No
                </span>
            </div>
        </div>

        <!-- Grid Card -->
        <div class="card">
            <h2>Grid Model</h2>
            <div class="metric">
                <span class="label">Voltage</span>
                <span class="value" id="gridVoltage">--</span>
            </div>
            <div class="metric">
                <span class="label">Frequency</span>
                <span class="value" id="gridFreq">--</span>
            </div>
            <div class="metric">
                <span class="label">Voltage (PU)</span>
                <span class="value" id="gridVoltagePU">--</span>
            </div>
            <div class="metric">
                <span class="label">Frequency (PU)</span>
                <span class="value" id="gridFreqPU">--</span>
            </div>
            <div class="metric">
                <span class="label">P Balance</span>
                <span class="value" id="gridPBalance">--</span>
            </div>
            <div class="metric">
                <span class="label">Q Balance</span>
                <span class="value" id="gridQBalance">--</span>
            </div>
            <div class="metric">
                <span class="label">Status</span>
                <span class="value" id="gridStable">
                    <span class="status-indicator active"></span>Stable
                </span>
            </div>
        </div>

        <!-- Devices Card -->
        <div class="card">
            <h2>Virtual Devices</h2>
            <div class="metric">
                <span class="label">Device Count</span>
                <span class="value" id="deviceCount">0</span>
            </div>
            <div id="deviceList">
                <div class="metric" style="color: #666; font-style: italic;">
                    No devices registered
                </div>
            </div>
        </div>
    </div>

    <div class="footer">
        Simulation Inspector - Development Tool | Read-only view of simulation state
    </div>

    <script>
        // WebSocket connection
        let ws;
        let reconnectAttempts = 0;
        const maxReconnectAttempts = 10;
        const reconnectDelay = 1000;

        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                document.getElementById('connectionStatus').textContent = 'Connected';
                document.getElementById('connectionStatus').style.background = '#00ff88';
                document.getElementById('connectionStatus').style.color = '#000';
                reconnectAttempts = 0;
            };
            
            ws.onclose = function() {
                document.getElementById('connectionStatus').textContent = 'Disconnected';
                document.getElementById('connectionStatus').style.background = '#ff4444';
                document.getElementById('connectionStatus').style.color = '#fff';
                
                // Attempt reconnection
                if (reconnectAttempts < maxReconnectAttempts) {
                    reconnectAttempts++;
                    setTimeout(connect, reconnectDelay);
                }
            };
            
            ws.onerror = function(err) {
                console.error('WebSocket error:', err);
            };
            
            ws.onmessage = function(event) {
                try {
                    const state = JSON.parse(event.data);
                    updateDisplay(state);
                } catch (e) {
                    console.error('Failed to parse state:', e);
                }
            };
        }

        function updateDisplay(state) {
            // Clock
            const elapsed = state.clock.elapsed / 1000000000; // nanoseconds to seconds
            document.getElementById('clockElapsed').textContent = formatDuration(state.clock.elapsed);
            document.getElementById('clockTicks').textContent = state.clock.tick_count.toLocaleString();
            document.getElementById('clockMode').textContent = state.clock.mode;
            
            const pausedEl = document.getElementById('clockPaused');
            if (state.clock.is_paused) {
                pausedEl.innerHTML = '<span class="status-indicator inactive"></span>Paused';
            } else {
                pausedEl.innerHTML = '<span class="status-indicator active"></span>Running';
            }

            // Sun
            document.getElementById('sunElevation').textContent = state.sun.elevation.toFixed(2) + '°';
            document.getElementById('sunAzimuth').textContent = state.sun.azimuth.toFixed(2) + '°';
            document.getElementById('sunIrradiance').textContent = state.sun.irradiance.toFixed(1) + ' W/m²';
            document.getElementById('sunDNI').textContent = state.sun.direct_normal.toFixed(1) + ' W/m²';
            document.getElementById('sunDiffuse').textContent = state.sun.diffuse.toFixed(1) + ' W/m²';
            
            const daytimeEl = document.getElementById('sunDaytime');
            if (state.sun.is_daytime) {
                daytimeEl.innerHTML = '<span class="status-indicator active"></span>Day';
                daytimeEl.className = 'value good';
            } else {
                daytimeEl.innerHTML = '<span class="status-indicator inactive"></span>Night';
                daytimeEl.className = 'value';
            }

            // Weather
            document.getElementById('weatherTemp').textContent = state.weather.temperature.toFixed(1) + ' °C';
            document.getElementById('weatherHumidity').textContent = state.weather.humidity.toFixed(1) + ' %';
            document.getElementById('weatherPressure').textContent = state.weather.pressure.toFixed(1) + ' hPa';
            document.getElementById('weatherClouds').textContent = (state.weather.cloud_cover * 100).toFixed(0) + ' %';
            document.getElementById('weatherWind').textContent = state.weather.wind_speed.toFixed(1) + ' m/s';
            document.getElementById('weatherWindDir').textContent = state.weather.wind_direction.toFixed(0) + '°';
            
            const rainEl = document.getElementById('weatherRain');
            if (state.weather.is_raining) {
                rainEl.innerHTML = '<span class="status-indicator active"></span>Yes';
                rainEl.className = 'value good';
            } else {
                rainEl.innerHTML = '<span class="status-indicator inactive"></span>No';
                rainEl.className = 'value';
            }

            // Grid
            const voltageEl = document.getElementById('gridVoltage');
            voltageEl.textContent = state.grid.voltage.toFixed(1) + ' V';
            voltageEl.className = 'value ' + getVoltageClass(state.grid.voltage_pu);
            
            const freqEl = document.getElementById('gridFreq');
            freqEl.textContent = state.grid.frequency.toFixed(3) + ' Hz';
            freqEl.className = 'value ' + getFrequencyClass(state.grid.frequency_pu);
            
            document.getElementById('gridVoltagePU').textContent = state.grid.voltage_pu.toFixed(4) + ' pu';
            document.getElementById('gridFreqPU').textContent = state.grid.frequency_pu.toFixed(4) + ' pu';
            
            const pBalanceEl = document.getElementById('gridPBalance');
            pBalanceEl.textContent = state.grid.active_balance.toFixed(2) + ' MW';
            pBalanceEl.className = 'value ' + (state.grid.active_balance === 0 ? 'good' : 'warning');
            
            const qBalanceEl = document.getElementById('gridQBalance');
            qBalanceEl.textContent = state.grid.reactive_balance.toFixed(2) + ' MVAr';
            qBalanceEl.className = 'value ' + (state.grid.reactive_balance === 0 ? 'good' : 'warning');
            
            const stableEl = document.getElementById('gridStable');
            if (state.grid.is_stable) {
                stableEl.innerHTML = '<span class="status-indicator active"></span>Stable';
                stableEl.className = 'value good';
            } else {
                stableEl.innerHTML = '<span class="status-indicator inactive"></span>Unstable';
                stableEl.className = 'value bad';
            }

            // Devices
            document.getElementById('deviceCount').textContent = state.devices.count;
            
            const deviceList = document.getElementById('deviceList');
            if (state.devices.devices && state.devices.devices.length > 0) {
                let html = '';
                for (const device of state.devices.devices) {
                    const statusClass = device.state === 'Running' ? 'active' : 'inactive';
                    html += '<div class="metric">';
                    html += '<span class="label">' + escapeHtml(device.name) + '</span>';
                    html += '<span class="value"><span class="status-indicator ' + statusClass + '"></span>' + escapeHtml(device.type) + '</span>';
                    html += '</div>';
                }
                deviceList.innerHTML = html;
            } else {
                deviceList.innerHTML = '<div class="metric" style="color: #666; font-style: italic;">No devices registered</div>';
            }
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        function formatDuration(durationNs) {
            const totalSeconds = Math.floor(durationNs / 1000000000);
            const days = Math.floor(totalSeconds / 86400);
            const hours = Math.floor((totalSeconds % 86400) / 3600);
            const minutes = Math.floor((totalSeconds % 3600) / 60);
            const seconds = totalSeconds % 60;
            
            if (days > 0) {
                return days + 'd ' + hours + 'h ' + minutes + 'm';
            } else if (hours > 0) {
                return hours + 'h ' + minutes + 'm ' + seconds + 's';
            } else {
                return minutes + 'm ' + seconds + 's';
            }
        }

        function getVoltageClass(pu) {
            if (pu >= 0.95 && pu <= 1.05) return 'good';
            if (pu >= 0.9 && pu <= 1.1) return 'warning';
            return 'bad';
        }

        function getFrequencyClass(pu) {
            if (pu >= 0.995 && pu <= 1.005) return 'good';
            if (pu >= 0.99 && pu <= 1.01) return 'warning';
            return 'bad';
        }

        // Initialize
        connect();
    </script>
</body>
</html>
`
