# core/security/rate_limiting.py
import asyncio
import time
import logging
from typing import Dict, Any, Optional

logger = logging.getLogger(__name__)

# --- Stubs for missing classes/context ---
class APISecurityContext:
    def __init__(self, rate_limit_tier: str = "default", threat_level: float = 0.1):
        self.rate_limit_tier = rate_limit_tier
        self.threat_level = threat_level

class ClientBehaviorAnalyzer:
    async def get_behavior_profile(self, client_id: str) -> Dict:
        logger.info(f"Analyzing behavior for client {client_id} (stub).")
        return {"trust_score": 0.5, "consistent_usage": True} # Placeholder

class RateLimitAnomalyDetector:
    def detect_anomaly(self, client_id: str, current_count: int, limit: int) -> bool:
        logger.info(f"Detecting anomaly for client {client_id} (stub).")
        return False # Placeholder
# --- End Stubs ---

class AdaptiveRateLimiter:
    """
    AI-driven adaptive rate limiting based on client behavior and threat level
    """
    
    def __init__(self, config: Dict = None):
        self.config = config if config is not None else {}
        self.rate_windows = self._initialize_rate_windows()
        self.behavior_analyzer = ClientBehaviorAnalyzer()
        self.anomaly_detector = RateLimitAnomalyDetector()
        self.request_counts: Dict[str, Dict[str, Dict[str, int]]] = {} # client_id -> endpoint -> window -> count
        
    def _initialize_rate_windows(self) -> Dict:
        return {
            "second": {"window": 1, "limits": {"default": 10}},
            "minute": {"window": 60, "limits": {"default": 100}},
            "hour": {"window": 3600, "limits": {"default": 1000}},
            "day": {"window": 86400, "limits": {"default": 10000}}
        }
    
    def _get_base_limits(self, tier: str, endpoint: str) -> Dict:
        limits = {}
        for window_name, window_config in self.rate_windows.items():
            limits[window_name] = window_config["limits"].get(tier, window_config["limits"]["default"])
        return limits

    async def _get_request_count(self, client_id: str, endpoint: str, window_name: str) -> int:
        # This would typically query a time-series database or Redis
        return self.request_counts.get(client_id, {}).get(endpoint, {}).get(window_name, 0)

    async def _increment_request_counters(self, client_id: str, endpoint: str):
        if client_id not in self.request_counts:
            self.request_counts[client_id] = {}
        if endpoint not in self.request_counts[client_id]:
            self.request_counts[client_id][endpoint] = {}
        
        for window_name, window_config in self.rate_windows.items():
            current_count = self.request_counts[client_id][endpoint].get(window_name, 0)
            self.request_counts[client_id][endpoint][window_name] = current_count + 1
            
            # In a real system, you'd also manage expiration of these counts

    async def _analyze_client_behavior(self, client_id: str, endpoint: str, context: APISecurityContext):
        await self.behavior_analyzer.get_behavior_profile(client_id) # Simulate analysis

    async def _get_current_usage(self, client_id: str, endpoint: str) -> Dict:
        usage = {}
        for window_name in self.rate_windows:
            usage[window_name] = await self._get_request_count(client_id, endpoint, window_name)
        return usage

    async def wait_if_needed(self, target: str):
        # This method is specifically for worker_manager's use case with external targets
        # It's a simplified version of check_rate_limit, assuming the worker itself is the "client"
        # and the target is what needs rate limiting.
        client_id = "worker-default" # Or some worker-specific ID
        endpoint = target # Treat the target as the endpoint
        
        # This should ideally integrate with the full check_rate_limit logic
        # For now, a simple delay based on some configuration
        await asyncio.sleep(self.config.get("worker_rate_limit_delay", 0.1))
        logger.info(f"Rate limiting applied for target {target} (simulated delay).")
    
    async def check_rate_limit(self, client_id: str, endpoint: str, context: APISecurityContext) -> Dict:
        """Check adaptive rate limits for client and endpoint"""
        # Get base limits for client tier
        base_limits = self._get_base_limits(context.rate_limit_tier, endpoint)
        
        # Adjust limits based on behavior and threat
        adaptive_limits = await self._calculate_adaptive_limits(client_id, base_limits, context)
        
        # Check all rate windows
        for window_name, window_config in self.rate_windows.items():
            current_count = await self._get_request_count(client_id, endpoint, window_name)
            window_limit = adaptive_limits[window_name]
            
            if current_count >= window_limit:
                return {
                    "allowed": False,
                    "window": window_name,
                    "current": current_count,
                    "limit": window_limit,
                    "retry_after": window_config['window']
                }
        
        # Increment counters
        await self._increment_request_counters(client_id, endpoint)
        
        # Analyze behavior for future adjustments
        asyncio.create_task(self._analyze_client_behavior(client_id, endpoint, context))
        
        return {
            "allowed": True,
            "limits": adaptive_limits,
            "current_usage": await self._get_current_usage(client_id, endpoint)
        }
    
    async def _calculate_adaptive_limits(self, client_id: str, base_limits: Dict, context: APISecurityContext) -> Dict:
        """Calculate adaptive limits based on client behavior"""
        behavior_profile = await self.behavior_analyzer.get_behavior_profile(client_id)
        threat_level = context.threat_level
        
        adaptive_limits = base_limits.copy()
        
        # Adjust based on threat level
        if threat_level > 0.7:
            # High threat - reduce limits
            for window in adaptive_limits:
                adaptive_limits[window] = int(adaptive_limits[window] * 0.5)
        elif threat_level < 0.3 and behavior_profile.get('trust_score', 0) > 0.8:
            # Low threat + trusted behavior - increase limits
            for window in adaptive_limits:
                adaptive_limits[window] = int(adaptive_limits[window] * 1.5)
        
        # Adjust based on historical behavior
        if behavior_profile.get('consistent_usage', False):
            for window in adaptive_limits:
                adaptive_limits[window] = int(adaptive_limits[window] * 1.2)
        
        return adaptive_limits
