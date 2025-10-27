// ~/HyperZilla/COMMAND_CENTER/static/js/facial_intel/facial_ui.js
class FacialIntelligenceUI {
    constructor() {
        this.currentImage = null;
        this.analysisInProgress = false;
        this.setupEventListeners();
    }

    setupEventListeners() {
        // File upload handling
        const uploadArea = document.getElementById('uploadArea');
        const fileInput = document.getElementById('faceImageInput');
        
        uploadArea.addEventListener('click', () => fileInput.click());
        uploadArea.addEventListener('dragover', (e) => this.handleDragOver(e));
        uploadArea.addEventListener('drop', (e) => this.handleFileDrop(e));
        
        fileInput.addEventListener('change', (e) => this.handleFileSelect(e));
    }

    handleDragOver(e) {
        e.preventDefault();
        e.currentTarget.classList.add('drag-over');
    }

    handleFileDrop(e) {
        e.preventDefault();
        e.currentTarget.classList.remove('drag-over');
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            this.processImageFile(files[0]);
        }
    }

    handleFileSelect(e) {
        const files = e.target.files;
        if (files.length > 0) {
            this.processImageFile(files[0]);
        }
    }

    processImageFile(file) {
        if (!file.type.startsWith('image/')) {
            this.showAlert('Please select a valid image file', 'error');
            return;
        }

        const reader = new FileReader();
        reader.onload = (e) => {
            this.displayImagePreview(e.target.result);
            this.currentImage = file;
        };
        reader.readAsDataURL(file);
    }

    displayImagePreview(imageData) {
        const uploadArea = document.getElementById('uploadArea');
        const preview = document.getElementById('imagePreview');
        const previewImage = document.getElementById('previewImage');
        
        uploadArea.style.display = 'none';
        preview.style.display = 'block';
        previewImage.src = imageData;
    }

    clearImage() {
        const uploadArea = document.getElementById('uploadArea');
        const preview = document.getElementById('imagePreview');
        const fileInput = document.getElementById('faceImageInput');
        
        uploadArea.style.display = 'block';
        preview.style.display = 'none';
        fileInput.value = '';
        this.currentImage = null;
    }

    async startFacialAnalysis() {
        if (!this.currentImage) {
            this.showAlert('Please select an image first', 'error');
            return;
        }

        if (this.analysisInProgress) {
            this.showAlert('Analysis already in progress', 'warning');
            return;
        }

        this.analysisInProgress = true;
        this.showAlert('Starting facial intelligence analysis...', 'info');

        // Update UI for analysis start
        this.updateAnalysisStatus('in_progress');

        try {
            const formData = new FormData();
            formData.append('image', this.currentImage);
            formData.append('enhancement_data', JSON.stringify(this.getEnhancementData()));
            formData.append('search_depth', document.getElementById('searchDepth').value);

            const response = await fetch('/api/facial-analysis', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();
            
            if (result.success) {
                this.displayFacialResults(result.data);
                this.showAlert('Facial analysis completed successfully!', 'success');
            } else {
                this.showAlert(`Analysis failed: ${result.error}`, 'error');
                this.updateAnalysisStatus('failed');
            }

        } catch (error) {
            this.showAlert('Network error during analysis', 'error');
            this.updateAnalysisStatus('failed');
            console.error('Facial analysis error:', error);
        } finally {
            this.analysisInProgress = false;
        }
    }

    getEnhancementData() {
        return {
            email: document.getElementById('knownEmail').value,
            username: document.getElementById('knownUsername').value,
            location: document.getElementById('knownLocation').value
        };
    }

    updateAnalysisStatus(status) {
        const statusElement = document.getElementById('faceDetectionStatus');
        const statusConfig = {
            'pending': { text: 'Waiting for analysis', icon: 'question-circle', class: 'status-pending' },
            'in_progress': { text: 'Analysis in progress...', icon: 'sync fa-spin', class: 'status-progress' },
            'completed': { text: 'Analysis complete', icon: 'check-circle', class: 'status-complete' },
            'failed': { text: 'Analysis failed', icon: 'times-circle', class: 'status-failed' }
        };

        const config = statusConfig[status] || statusConfig.pending;
        statusElement.className = config.class;
        statusElement.innerHTML = `<i class="fas fa-${config.icon}"></i> ${config.text}`;
    }

    displayFacialResults(results) {
        this.updateAnalysisStatus('completed');
        
        // Update summary tab
        this.updateSummaryTab(results);
        
        // Update matches tab
        this.updateMatchesTab(results.database_matches);
        
        // Update analysis tab
        this.updateAnalysisTab(results.primary_face_analysis);
        
        // Update OSINT tab
        this.updateOSINTTab(results);
    }

    updateSummaryTab(results) {
        // Update confidence score
        const confidence = document.getElementById('confidenceScore');
        confidence.textContent = `${Math.round(results.confidence_score * 100)}%`;
        confidence.className = `score-${this.getConfidenceLevel(results.confidence_score)}`;
        
        // Update best match
        const bestMatch = document.getElementById('bestMatch');
        if (results.database_matches && results.database_matches.length > 0) {
            const match = results.database_matches[0];
            bestMatch.innerHTML = `
                <strong>${match.identity}</strong><br>
                <small>Confidence: ${Math.round(match.confidence * 100)}%</small>
            `;
        } else {
            bestMatch.textContent = 'No matches found';
        }
        
        // Update actionable insights
        const insights = document.getElementById('actionableInsights');
        if (results.actionable_intelligence) {
            insights.innerHTML = results.actionable_intelligence.recommended_next_steps
                .map(step => `<div class="insight-item">${step}</div>`)
                .join('');
        }
    }

    updateMatchesTab(matches) {
        const matchesList = document.getElementById('matchesList');
        
        if (!matches || matches.length === 0) {
            matchesList.innerHTML = '<div class="no-matches">No database matches found</div>';
            return;
        }
        
        matchesList.innerHTML = matches.map(match => `
            <div class="match-item">
                <div class="match-identity">${match.identity}</div>
                <div class="match-confidence">
                    <div class="confidence-bar">
                        <div class="confidence-fill" style="width: ${match.confidence * 100}%"></div>
                    </div>
                    <span class="confidence-value">${Math.round(match.confidence * 100)}%</span>
                </div>
                <div class="match-source">Source: ${match.source}</div>
                ${match.metadata ? `<div class="match-metadata">${JSON.stringify(match.metadata)}</div>` : ''}
            </div>
        `).join('');
    }

    updateAnalysisTab(analysis) {
        // Update demographics
        if (analysis.demographics) {
            const demo = analysis.demographics;
            document.getElementById('demographicsAnalysis').innerHTML = `
                <div>Age: ${demo.age || 'Unknown'}</div>
                <div>Gender: ${demo.gender || 'Unknown'}</div>
                <div>Ethnicity: ${demo.ethnicity || 'Unknown'}</div>
            `;
        }
        
        // Update emotions
        if (analysis.emotions) {
            const emotions = analysis.emotions;
            document.getElementById('emotionalAnalysis').innerHTML = `
                <div>Dominant: ${emotions.dominant_emotion}</div>
                <div>Intensity: ${Math.round(emotions.emotional_intensity * 100)}%</div>
            `;
        }
    }

    updateOSINTTab(results) {
        const osintResults = document.getElementById('osintResults');
        
        let osintHTML = '';
        
        if (results.enhanced_osint && results.enhanced_osint.social_media_profiles) {
            const profiles = results.enhanced_osint.social_media_profiles;
            osintHTML += `
                <h5>Social Media Profiles</h5>
                ${profiles.map(profile => `
                    <div class="osint-profile">
                        <strong>${profile.platform}</strong>: ${profile.username}
                        ${profile.url ? `<br><a href="${profile.url}" target="_blank">View Profile</a>` : ''}
                    </div>
                `).join('')}
            `;
        }
        
        if (results.enhanced_osint && results.enhanced_osint.verified_identities) {
            const verified = results.enhanced_osint.verified_identities;
            osintHTML += `
                <h5>Verified Identities</h5>
                ${verified.map(identity => `
                    <div class="verified-identity">${identity}</div>
                `).join('')}
            `;
        }
        
        osintResults.innerHTML = osintHTML || '<div class="no-osint">No OSINT data available</div>';
    }

    getConfidenceLevel(score) {
        if (score >= 0.8) return 'high';
        if (score >= 0.6) return 'medium';
        return 'low';
    }

    showAlert(message, type) {
        // Use the existing alert system from main UI
        if (window.hyperZillaUI) {
            window.hyperZillaUI.showAlert(message, type);
        }
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.facialIntelUI = new FacialIntelligenceUI();
});

// Global functions for HTML onclick handlers
function startFacialAnalysis() {
    if (window.facialIntelUI) {
        window.facialIntelUI.startFacialAnalysis();
    }
}

function clearImage() {
    if (window.facialIntelUI) {
        window.facialIntelUI.clearImage();
    }
}

function showFacialDatabase() {
    // Implementation for showing facial database management
    console.log('Show facial database');
}
