export class HapticFeedbackController {
    constructor() {
        this.feedbackProfiles = new Map();
        this.setupFeedbackProfiles();
    }

    setupFeedbackProfiles() {
        // Define haptic feedback patterns for different events
        this.feedbackProfiles.set('target_acquired', {
            pattern: [100, 50, 100],
            intensity: 0.8,
            description: 'Double pulse - target locked'
        });

        this.feedbackProfiles.set('threat_detected', {
            pattern: [200, 100, 200, 100, 200],
            intensity: 1.0,
            description: 'Rapid pulses - threat alert'
        });

        this.feedbackProfiles.set('operation_complete', {
            pattern: [300],
            intensity: 0.6,
            description: 'Single pulse - mission success'
        });

        this.feedbackProfiles.set('system_error', {
            pattern: [100, 100, 100, 100],
            intensity: 0.9,
            description: 'Staccato pulses - system error'
        });

        this.feedbackProfiles.set('data_received', {
            pattern: [50],
            intensity: 0.3,
            description: 'Quick pulse - data incoming'
        });
    }

    async triggerHapticFeedback(event, options = {}) {
        const profile = this.feedbackProfiles.get(event);
        if (!profile) {
            console.log(`❌ Haptic profile not found: ${event}`);
            return;
        }

        if (this.supportsHapticFeedback()) {
            await this.executeHapticPattern(profile, options);
        } else {
            await this.fallbackFeedback(profile, event);
        }
    }

    async executeHapticPattern(profile, options) {
        const { pattern, intensity } = profile;
        const actualIntensity = options.intensity || intensity;

        try {
            if (navigator.vibrate) {
                // Web Vibration API
                navigator.vibrate(pattern);
            } else if (this.isMobileDevice()) {
                // Mobile-specific haptic feedback
                await this.mobileHapticFeedback(pattern, actualIntensity);
            } else {
                // Desktop fallback
                await this.desktopHapticFeedback(profile, options);
            }

            await this.logHapticEvent(event, 'EXECUTED');

        } catch (error) {
            console.log(`❌ Haptic feedback error: ${error.message}`);
            await this.fallbackFeedback(profile, event);
        }
    }

    async mobileHapticFeedback(pattern, intensity) {
        // iOS-specific haptic feedback
        if (window.WebKitHapticFeedback) {
            pattern.forEach(duration => {
                window.WebKitHapticFeedback.tapSoft();
            });
        }
        // Android-specific vibration
        else if (navigator.vibrate) {
            navigator.vibrate(pattern);
        }
    }

    async desktopHapticFeedback(profile, options) {
        // Simulate haptic feedback on desktop
        console.log(`🎮 HAPTIC: ${profile.description}`);
        
        // Could integrate with gaming controllers or specialized hardware
        if (options.controller) {
            await this.controllerHapticFeedback(profile, options.controller);
        }
    }

    fallbackFeedback(profile, event) {
        // Visual feedback when haptic isn't available
        const visualMap = {
            'target_acquired': '🎯 [TARGET ACQUIRED - HAPTIC]',
            'threat_detected': '🚨 [THREAT DETECTED - HAPTIC]',
            'operation_complete': '✅ [OPERATION COMPLETE - HAPTIC]',
            'system_error': '⚠️ [SYSTEM ERROR - HAPTIC]',
            'data_received': '📥 [DATA RECEIVED - HAPTIC]'
        };

        console.log(visualMap[event] || `🎮 [HAPTIC: ${event}]`);
    }

    supportsHapticFeedback() {
        return typeof navigator !== 'undefined' && 
               (navigator.vibrate || 
                window.WebKitHapticFeedback ||
                this.isMobileDevice());
    }

    isMobileDevice() {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    }

    async logHapticEvent(event, status) {
        const logEntry = {
            timestamp: new Date().toISOString(),
            event: event,
            status: status,
            platform: process.platform,
            haptic_support: this.supportsHapticFeedback()
        };
        
        console.log(`🎮 Haptic Event: ${event} - ${status}`);
    }
}
