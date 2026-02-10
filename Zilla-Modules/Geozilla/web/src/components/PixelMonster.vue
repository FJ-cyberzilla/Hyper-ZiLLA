<!-- web/src/components/PixelMonster.vue -->
<template>
  <div class="galaxy-container">
    <!-- Animated Galaxy Background -->
    <div class="distant-galaxy">
      <div class="stars"></div>
      <div class="nebula"></div>
      <div class="black-hole"></div>
      <div class="pulsars"></div>
      <div class="asteroid-field"></div>
    </div>
    
    <!-- Cognizilla Pixel Monster -->
    <div class="pixel-monster" :class="{ 'active': isActive, 'scanning': isScanning }">
      <!-- Monster Body -->
      <div class="monster-body">
        <div class="pixel-eye left-eye" :style="eyeStyle"></div>
        <div class="pixel-eye right-eye" :style="eyeStyle"></div>
        <div class="energy-core" :style="coreStyle"></div>
        <div class="tentacles">
          <div v-for="i in 8" :key="i" class="tentacle" :style="getTentacleStyle(i)"></div>
        </div>
      </div>
      
      <!-- Status Effects -->
      <div class="status-effects">
        <div class="quantum-pulse" v-if="isScanning"></div>
        <div class="data-stream" v-if="isTransmitting"></div>
        <div class="shield-bubble" v-if="isProtected"></div>
      </div>
    </div>
    
    <!-- Interactive Controls -->
    <div class="monster-controls">
      <button @click="toggleActivation" class="control-btn activate">
        {{ isActive ? 'DEACTIVATE' : 'ACTIVATE COGNIZILLA' }}
      </button>
      <button @click="startScan" class="control-btn scan" :disabled="!isActive">
        QUANTUM SCAN
      </button>
      <button @click="showDashboard" class="control-btn dashboard">
        SOVEREIGN DASHBOARD
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PixelMonster',
  data() {
    return {
      isActive: false,
      isScanning: false,
      isTransmitting: false,
      isProtected: true,
      eyeColor: '#00ff88',
      coreEnergy: 0,
      tentaclePhase: 0
    }
  },
  computed: {
    eyeStyle() {
      return {
        'background-color': this.eyeColor,
        'box-shadow': `0 0 20px ${this.eyeColor}, 0 0 40px ${this.eyeColor}`
      }
    },
    coreStyle() {
      return {
        transform: `scale(${1 + Math.sin(this.coreEnergy) * 0.3})`,
        'background-color': `hsl(${this.coreEnergy * 360}, 100%, 50%)`
      }
    }
  },
  mounted() {
    this.animateMonster();
    this.startEnergyPulse();
  },
  methods: {
    animateMonster() {
      // Continuous animation loop
      this.animationFrame = requestAnimationFrame(() => {
        this.tentaclePhase += 0.02;
        this.coreEnergy += 0.01;
        this.animateMonster();
      });
    },
    
    startEnergyPulse() {
      setInterval(() => {
        if (this.isActive) {
          this.eyeColor = `hsl(${Math.random() * 60 + 120}, 100%, 50%)`;
        }
      }, 2000);
    },
    
    getTentacleStyle(index) {
      const phase = this.tentaclePhase + (index * 0.5);
      return {
        transform: `rotate(${Math.sin(phase) * 15}deg) scaleY(${0.8 + Math.sin(phase) * 0.2})`
      };
    },
    
    toggleActivation() {
      this.isActive = !this.isActive;
      if (this.isActive) {
        this.$emit('monster-activated');
        this.startQuantumBootSequence();
      } else {
        this.$emit('monster-deactivated');
      }
    },
    
    startScan() {
      this.isScanning = true;
      this.$emit('quantum-scan-started');
      
      // Simulate scan completion
      setTimeout(() => {
        this.isScanning = false;
        this.$emit('quantum-scan-completed');
      }, 3000);
    },
    
    startQuantumBootSequence() {
      const bootSteps = [
        { delay: 500, action: 'Quantum Core Initializing...' },
        { delay: 1000, action: 'Neural Networks Booting...' },
        { delay: 1500, action: 'Security Shields Activated...' },
        { delay: 2000, action: 'COGNIZILLA READY' }
      ];
      
      bootSteps.forEach(step => {
        setTimeout(() => {
          this.$emit('boot-progress', step.action);
        }, step.delay);
      });
    },
    
    showDashboard() {
      this.$emit('show-dashboard');
    }
  },
  beforeDestroy() {
    cancelAnimationFrame(this.animationFrame);
  }
}
</script>

<style scoped>
.galaxy-container {
  position: relative;
  width: 100%;
  height: 100vh;
  background: radial-gradient(ellipse at center, #0c0c2e 0%, #000000 100%);
  overflow: hidden;
}

/* Distant Galaxy Background */
.distant-galaxy {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.stars {
  position: absolute;
  width: 100%;
  height: 100%;
  background-image: 
    radial-gradient(2px 2px at 20% 30%, #fff, transparent),
    radial-gradient(2px 2px at 40% 70%, #fff, transparent),
    radial-gradient(1px 1px at 60% 20%, #fff, transparent),
    radial-gradient(1px 1px at 80% 90%, #fff, transparent);
  animation: twinkle 4s infinite alternate;
}

.nebula {
  position: absolute;
  top: 20%;
  right: 10%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(74, 0, 124, 0.3) 0%, transparent 70%);
  filter: blur(40px);
  animation: float 20s infinite ease-in-out;
}

.black-hole {
  position: absolute;
  bottom: 15%;
  left: 10%;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, #000 0%, #330066 70%);
  border-radius: 50%;
  filter: blur(20px);
  animation: pulsate 8s infinite;
}

/* Pixel Monster Styles */
.pixel-monster {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  transition: all 0.5s ease;
}

.monster-body {
  position: relative;
  width: 120px;
  height: 120px;
  background: #1a1a3a;
  border: 4px solid #00ff88;
  border-radius: 35% 65% 45% 55% / 65% 35% 65% 35%;
  box-shadow: 
    0 0 30px rgba(0, 255, 136, 0.5),
    inset 0 0 20px rgba(0, 255, 136, 0.2);
  animation: float 6s infinite ease-in-out;
}

.pixel-eye {
  position: absolute;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  top: 35%;
}

.left-eye { left: 30%; }
.right-eye { right: 30%; }

.energy-core {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 30px;
  height: 30px;
  background: #ff0088;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  transition: all 0.3s ease;
}

.tentacles {
  position: absolute;
  bottom: -20px;
  left: 50%;
  transform: translateX(-50%);
}

.tentacle {
  position: absolute;
  width: 8px;
  height: 40px;
  background: #00ff88;
  border-radius: 4px;
  transform-origin: top center;
  transition: all 0.5s ease;
}

/* Monster States */
.pixel-monster.active .monster-body {
  border-color: #ff0088;
  box-shadow: 
    0 0 50px rgba(255, 0, 136, 0.7),
    inset 0 0 30px rgba(255, 0, 136, 0.3);
}

.pixel-monster.scanning {
  animation: scan-pulse 1s infinite;
}

/* Controls */
.monster-controls {
  position: absolute;
  bottom: 50px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 15px;
}

.control-btn {
  padding: 12px 24px;
  border: 2px solid #00ff88;
  background: rgba(0, 255, 136, 0.1);
  color: #00ff88;
  font-family: 'Courier New', monospace;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
}

.control-btn:hover {
  background: rgba(0, 255, 136, 0.3);
  box-shadow: 0 0 20px rgba(0, 255, 136, 0.5);
}

.control-btn.activate {
  border-color: #ff0088;
  color: #ff0088;
  background: rgba(255, 0, 136, 0.1);
}

.control-btn.activate:hover {
  background: rgba(255, 0, 136, 0.3);
}

/* Animations */
@keyframes twinkle {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

@keyframes float {
  0%, 100% { transform: translate(-50%, -50%) translateY(0px); }
  50% { transform: translate(-50%, -50%) translateY(-20px); }
}

@keyframes pulsate {
  0%, 100% { transform: scale(1); opacity: 0.7; }
  50% { transform: scale(1.1); opacity: 1; }
}

@keyframes scan-pulse {
  0%, 100% { filter: drop-shadow(0 0 10px #00ff88); }
  50% { filter: drop-shadow(0 0 30px #ff0088); }
}
</style>
