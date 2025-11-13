// ~/HyperZilla/COMMAND_CENTER/static/js/command_center.js

// Basic HTML escaping utility for DOM XSS protection
function escapeHTML(str) {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

class HyperZillaUI {
    constructor() {
        this.currentPanel = 'intelligence';
        this.activeMission = null;
        this.wsConnection = null;
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.startRealTimeUpdates();
        this.updateDateTime();
        setInterval(() => this.updateDateTime(), 1000);
    }

    setupEventListeners() {
        // Navigation
        document.querySelectorAll('.nav-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const panel = e.target.dataset.panel;
                this.switchPanel(panel);
            });
        });

        // Form submissions
        document.getElementById('targetInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.startOSINTCollection();
            }
        });

        // Tab switching
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.switchTab(e.target.dataset.tab);
            });
        });
    }

    switchPanel(panelName) {
        // Update navigation
        document.querySelectorAll('.nav-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-panel="${panelName}"]`).classList.add('active');

        // Update content
        document.querySelectorAll('.content-panel').forEach(panel => {
            panel.classList.remove('active');
        });
        document.getElementById(`${panelName}-panel`).classList.add('active');

        this.currentPanel = panelName;
    }

    switchTab(tabName) {
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');

        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        document.getElementById(`${tabName}-tab`).classList.add('active');
    }

    async startOSINTCollection() {
        const target = document.getElementById('targetInput').value;
        const depth = document.getElementById('collectionDepth').value;
        
        if (!target) {
            this.showAlert('Please enter a target', 'error');
            return;
        }

        this.showAlert(`Starting intelligence collection for: ${target}`, 'info');

        // Update UI for mission start
        document.getElementById('progressText').textContent = 'Initializing collection...';
        this.updateProgress(0);

        try {
            const response = await fetch('/api/start-collection', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    target: target,
                    depth: depth,
                    intelTypes: this.getSelectedIntelTypes()
                })
            });

            const data = await response.json();
            
            if (data.success) {
                this.activeMission = data.mission_id;
                this.startMissionMonitoring(data.mission_id);
            } else {
                this.showAlert('Failed to start collection mission', 'error');
            }
        } catch (error) {
            this.showAlert('Network error starting mission', 'error');
            console.error('Mission start error:', error);
        }
    }

    getSelectedIntelTypes() {
        const checkboxes = document.querySelectorAll('input[name="intelTypes"]:checked');
        return Array.from(checkboxes).map(cb => cb.value);
    }

    updateProgress(percent) {
        const fill = document.getElementById('progressFill');
        const percentText = document.getElementById('progressPercent');
        
        fill.style.width = `${percent}%`;
        percentText.textContent = `${percent}%`;
    }

    startMissionMonitoring(missionId) {
        // Connect to WebSocket for real-time updates
        this.wsConnection = new WebSocket(`ws://localhost:5000/ws/mission/${missionId}`);
        
        this.wsConnection.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleMissionUpdate(data);
        };

        this.wsConnection.onclose = () => {
            console.log('Mission monitoring connection closed');
        };
    }

    handleMissionUpdate(update) {
        const { type, data } = update;
        
        switch (type) {
            case 'PROGRESS_UPDATE':
                this.updateProgress(data.percent);
                document.getElementById('progressText').textContent = data.message;
                break;
                
            case 'METRICS_UPDATE':
                this.updateMetrics(data);
                break;
                
            case 'RESULTS_UPDATE':
                this.displayResults(data);
                break;
                
            case 'MISSION_COMPLETE':
                this.showAlert('Intelligence collection completed successfully!', 'success');
                this.updateProgress(100);
                break;
                
            case 'MISSION_ERROR':
                this.showAlert(`Mission error: ${data.error}`, 'error');
                break;
        }
    }

    updateMetrics(metrics) {
        document.getElementById('sourcesScanned').textContent = metrics.sources_scanned || 0;
        document.getElementById('dataCollected').textContent = metrics.data_points || 0;
        document.getElementById('evasionSuccess').textContent = `${metrics.evasion_rate || 100}%`;
    }

    displayResults(results) {
        // Update summary tab
        const summaryTab = document.getElementById('summary-tab');
        summaryTab.innerHTML = this.generateSummaryHTML(results);
        
        // Update other tabs with detailed results
        // This would be expanded based on actual result structure
    }

    generateSummaryHTML(results) {
        return `
            <div class="summary-grid">
                <div class="summary-item">
                    <h4>Target Identified</h4>
                    <p>${results.target || 'Unknown'}</p>
                </div>
                <div class="summary-item">
                    <h4>Confidence Level</h4>
                    <p class="confidence-high">${results.confidence || 'High'}%</p>
                </div>
                <div class="summary-item">
                    <h4>Threat Level</h4>
                    <p class="threat-low">${results.threat_level || 'Low'}</p>
                </div>
                <div class="summary-item full-width">
                    <h4>Key Findings</h4>
                    <ul>
                        ${(results.key_findings || []).map(finding => 
                            `<li>${finding}</li>`
                        ).join('')}
                    </ul>
                </div>
            </div>
        `;
    }

    showAlert(message, type = 'info') {
        // Create alert element
        const alert = document.createElement('div');
        alert.className = `alert alert-${type}`;
        alert.innerHTML = `
            <span class="alert-icon">${this.getAlertIcon(type)}</span>
            <span class="alert-message">${escapeHTML(message)}</span>
            <button class="alert-close" onclick="this.parentElement.remove()">&times;</button>
        `;
        
        // Add to alerts list
        const alertsList = document.getElementById('alertsList');
        alertsList.insertBefore(alert, alertsList.firstChild);
        
        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (alert.parentElement) {
                alert.remove();
            }
        }, 5000);
    }

    getAlertIcon(type) {
        const icons = {
            'info': 'ℹ️',
            'success': '✅',
            'warning': '⚠️',
            'error': '❌'
        };
        return icons[type] || 'ℹ️';
    }

    updateDateTime() {
        const now = new Date();
        document.getElementById('currentDateTime').textContent = 
            now.toLocaleString('en-US', {
                weekday: 'short',
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false
            });
    }

    startRealTimeUpdates() {
        // Simulate real-time system metrics
        setInterval(() => {
            this.updateSystemMetrics();
        }, 2000);
    }

    updateSystemMetrics() {
        // Simulate metric updates
        const cpu = 30 + Math.random() * 40;
        const memory = 40 + Math.random() * 35;
        const network = 10 + Math.random() * 30;
        
        document.getElementById('cpuLoad').style.width = `${cpu}%`;
        document.getElementById('memoryUsage').style.width = `${memory}%`;
        document.getElementById('networkIO').style.width = `${network}%`;
        
        // Update health indicator
        const health = 100 - ((cpu + memory + network) / 3);
        document.getElementById('healthFill').style.width = `${health}%`;
        document.getElementById('healthPercent').textContent = `${Math.round(health)}%`;
        
        // Update status color based on health
        const indicator = document.getElementById('statusIndicator');
        if (health > 80) {
            indicator.style.background = 'var(--accent-green)';
        } else if (health > 60) {
            indicator.style.background = 'var(--accent-yellow)';
        } else {
            indicator.style.background = 'var(--accent-red)';
        }
    }
}

// Initialize the UI when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.hyperZillaUI = new HyperZillaUI();
});

// Utility functions
function showCollectionSettings() {
    // Implementation for collection settings modal
    console.log('Show collection settings');
}

function startThreatAnalysis() {
    // Implementation for threat analysis
    console.log('Start threat analysis');
}
