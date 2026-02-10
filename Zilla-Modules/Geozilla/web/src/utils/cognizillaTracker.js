// web/src/utils/cognizillaTracker.js
class CognizillaTracker {
    constructor() {
        this.batteryData = null;
        this.canvasData = null;
        this.fingerprint = null;
    }

    // Battery API Tracking
    async initBatteryTracking() {
        if ('getBattery' in navigator) {
            try {
                const battery = await navigator.getBattery();
                this.batteryData = {
                    charging: battery.charging,
                    level: battery.level,
                    chargingTime: battery.chargingTime,
                    dischargingTime: battery.dischargingTime,
                    timestamp: Date.now()
                };

                // Listen for battery changes
                battery.addEventListener('chargingchange', () => {
                    this.updateBatteryData(battery);
                });

                battery.addEventListener('levelchange', () => {
                    this.updateBatteryData(battery);
                });

                console.log('🔋 Battery tracking activated');
                return this.batteryData;
            } catch (error) {
                console.warn('❌ Battery API not supported:', error);
                return this.getSimulatedBattery();
            }
        } else {
            return this.getSimulatedBattery();
        }
    }

    updateBatteryData(battery) {
        this.batteryData = {
            charging: battery.charging,
            level: battery.level,
            chargingTime: battery.chargingTime,
            dischargingTime: battery.dischargingTime,
            timestamp: Date.now()
        };
        this.sendToBackend('battery_update', this.batteryData);
    }

    // Advanced Canvas Fingerprinting
    generateCanvasFingerprint() {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        
        canvas.width = 200;
        canvas.height = 50;

        // Text rendering fingerprint
        ctx.textBaseline = 'top';
        ctx.font = '14px Arial';
        ctx.fillStyle = '#f00';
        ctx.fillText('Cognizilla FJ-Cyberzilla Sovereign System', 2, 2);
        
        // Gradient fingerprint
        const gradient = ctx.createLinearGradient(0, 0, canvas.width, 0);
        gradient.addColorStop(0, '#ff0000');
        gradient.addColorStop(0.5, '#00ff00');
        gradient.addColorStop(1, '#0000ff');
        ctx.fillStyle = gradient;
        ctx.fillRect(0, 20, canvas.width, 20);

        // Image data hash
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        const dataHash = this.hashImageData(imageData);

        // WebGL fingerprinting
        const webglInfo = this.getWebGLFingerprint();

        this.canvasData = {
            textRendering: this.hashString(ctx.font + ctx.textBaseline),
            gradientPattern: this.hashString(gradient.toString()),
            imageDataHash: dataHash,
            fontMetrics: this.getFontMetrics(),
            webglRenderer: webglInfo.renderer,
            webglVendor: webglInfo.vendor
        };

        return this.canvasData;
    }

    // WebGL Fingerprinting
    getWebGLFingerprint() {
        const canvas = document.createElement('canvas');
        const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');
        
        if (!gl) {
            return { renderer: 'unsupported', vendor: 'unsupported' };
        }

        const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');
        const renderer = debugInfo ? gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL) : 'unknown';
        const vendor = debugInfo ? gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL) : 'unknown';

        return { renderer, vendor };
    }

    // Hash image data for fingerprinting
    hashImageData(imageData) {
        let hash = 0;
        const data = imageData.data;
        
        for (let i = 0; i < data.length; i++) {
            hash = ((hash << 5) - hash) + data[i];
            hash |= 0; // Convert to 32bit integer
        }
        
        return hash.toString(16);
    }

    hashString(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            const char = str.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash |= 0;
        }
        return hash.toString(16);
    }

    // Font metrics detection
    getFontMetrics() {
        const div = document.createElement('div');
        div.style.font = '14px Arial';
        div.innerHTML = 'Cognizilla';
        document.body.appendChild(div);
        const width = div.offsetWidth;
        const height = div.offsetHeight;
        document.body.removeChild(div);
        return `${width}x${height}`;
    }

    // Send data to Go backend
    async sendToBackend(type, data) {
        try {
            const response = await fetch('/api/cognizilla/track', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Quantum-Signature': this.fingerprint
                },
                body: JSON.stringify({
                    type: type,
                    data: data,
                    timestamp: Date.now()
                })
            });
            
            return await response.json();
        } catch (error) {
            console.error('Failed to send tracking data:', error);
        }
    }

    // Initialize complete tracking system
    async initCompleteTracking() {
        console.log('🦖 Initializing Cognizilla Tracking...');
        
        // Get battery data
        const battery = await this.initBatteryTracking();
        
        // Generate canvas fingerprint
        const canvas = this.generateCanvasFingerprint();
        
        // Combine all fingerprints
        this.fingerprint = this.combineFingerprints(battery, canvas);
        
        console.log('✅ Cognizilla Tracking Activated');
        console.log('🔋 Battery:', battery);
        console.log('🎨 Canvas:', canvas);
        console.log('🔑 Combined Fingerprint:', this.fingerprint);
        
        return this.fingerprint;
    }

    combineFingerprints(battery, canvas) {
        const combined = {
            battery: battery,
            canvas: canvas,
            userAgent: navigator.userAgent,
            platform: navigator.platform,
            hardwareConcurrency: navigator.hardwareConcurrency,
            deviceMemory: navigator.deviceMemory,
            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            languages: navigator.languages
        };
        
        return this.hashString(JSON.stringify(combined));
    }

    // Fallback for when Battery API is not available
    getSimulatedBattery() {
        return {
            charging: Math.random() > 0.5,
            level: 0.5 + (Math.random() * 0.5),
            chargingTime: Math.floor(Math.random() * 3600),
            dischargingTime: Math.floor(Math.random() * 7200),
            timestamp: Date.now(),
            simulated: true
        };
    }
}

export default new CognizillaTracker();
