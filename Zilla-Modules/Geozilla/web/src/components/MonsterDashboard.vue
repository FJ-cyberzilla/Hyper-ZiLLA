// web/src/components/MonsterDashboard.vue
<template>
  <div class="ai-command-center">
    <RealTimeMonitor :neuralData="cognitiveData" />
    <PredictiveDisplay :forecasts="aiPredictions" />
    <JackpotController @execute="triggerJackpot" />
  </div>
</template>

<script>
export default {
  methods: {
    async triggerJackpot() {
      const result = await this.$aiOrchestrator.executeAutonomousStrategy();
      this.$aiVisualizer.displayResults(result);
    }
  }
}
</script>
