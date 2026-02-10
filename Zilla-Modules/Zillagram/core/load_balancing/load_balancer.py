# core/load_balancing/load_balancer.py
import asyncio
from typing import Dict, List
from datetime import datetime
import statistics

class IntelligentLoadBalancer:
    """
    Intelligent load balancing with dynamic routing and auto-scaling
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.traffic_analyzer = TrafficAnalyzer(config)
        self.auto_scaler = AutoScaler(config)
        self.health_checker = LoadBalancerHealthChecker(config)
        
    async def configure_load_balancing(self) -> Dict:
        """Configure intelligent load balancing strategies"""
        lb_configuration = {
            "round_robin": await self._configure_round_robin(),
            "least_connections": await self._configure_least_connections(),
            "weighted_distribution": await self._configure_weighted_distribution(),
            "geographic_routing": await self._configure_geo_routing(),
            "adaptive_load_balancing": await self._configure_adaptive_lb()
        }
        
        # Start traffic analysis
        asyncio.create_task(self.traffic_analyzer.analyze_traffic_patterns())
        
        return lb_configuration
    
    async def _configure_adaptive_lb(self) -> Dict:
        """Configure adaptive load balancing based on real-time metrics"""
        adaptive_config = {
            "algorithm": "adaptive",
            "metrics": ["cpu_usage", "memory_usage", "response_time", "error_rate"],
            "adjustment_interval": "10s",
            "weight_calculation": await self._create_weight_calculation_function()
        }
        
        return await self._apply_adaptive_configuration(adaptive_config)
    
    async def handle_traffic_spike(self, traffic_data: Dict) -> Dict:
        """Handle traffic spikes with intelligent scaling"""
        spike_response = {
            "spike_detected": True,
            "timestamp": datetime.utcnow(),
            "traffic_increase": traffic_data['increase_percentage'],
            "immediate_actions": await self._take_immediate_actions(traffic_data),
            "scaling_decisions": await self.auto_scaler.assess_scaling_needs(traffic_data),
            "routing_adjustments": await self._adjust_routing_strategy(traffic_data)
        }
        
        return spike_response

class AutoScaler:
    """Intelligent auto-scaling based on multiple metrics"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.metrics_analyzer = MetricsAnalyzer(config)
        self.scaling_predictor = ScalingPredictor(config)
        
    async def assess_scaling_needs(self, metrics: Dict) -> Dict:
        """Assess scaling needs based on comprehensive metrics"""
        scaling_analysis = {
            "cpu_based": await self._analyze_cpu_scaling(metrics),
            "memory_based": await self._analyze_memory_scaling(metrics),
            "custom_metrics": await self._analyze_custom_metrics(metrics),
            "predictive_scaling": await self.scaling_predictor.predict_scaling_needs(metrics)
        }
        
        scaling_decisions = self._make_scaling_decisions(scaling_analysis)
        
        # Execute scaling if needed
        if scaling_decisions['scale_out'] or scaling_decisions['scale_in']:
            await self._execute_scaling(scaling_decisions)
        
        return scaling_decisions
    
    async def _analyze_cpu_scaling(self, metrics: Dict) -> Dict:
        """Analyze CPU-based scaling needs"""
        cpu_metrics = metrics.get('cpu', {})
        avg_cpu = statistics.mean(cpu_metrics.get('usage_percentages', [0]))
        
        scaling_recommendation = {
            "current_usage": avg_cpu,
            "threshold_breached": avg_cpu > self.config['cpu_threshold'],
            "recommended_action": "scale_out" if avg_cpu > 70 else "scale_in" if avg_cpu < 30 else "maintain",
            "confidence": min(avg_cpu / 100, 1.0)
        }
        
        return scaling_recommendation
    
    async def setup_predictive_scaling(self) -> Dict:
        """Setup predictive scaling based on historical patterns"""
        predictive_config = {
            "algorithm": "arima",  # or prophet, lstm
            "training_data_days": 30,
            "prediction_horizon": "24h",
            "confidence_threshold": 0.8,
            "seasonality_analysis": True
        }
        
        # Train predictive model
        model = await self.scaling_predictor.train_predictive_model()
        
        return {
            "predictive_model_trained": model is not None,
            "config": predictive_config,
            "model_accuracy": await self.scaling_predictor.get_model_accuracy()
        }
