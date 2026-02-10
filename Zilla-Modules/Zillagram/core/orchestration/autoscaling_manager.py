# core/orchestration/autoscaling_manager.py
import asyncio
import logging
from typing import Dict, Any

logger = logging.getLogger(__name__)

class AutoScalingManager:
    """
    Manages auto-scaling for worker processes based on various metrics.
    This implementation focuses on internal worker management within a single host
    or across a set of hosts where the manager has control over process creation/destruction.
    For Kubernetes, a Horizontal Pod Autoscaler (HPA) would typically be used.
    """
    def __init__(self, config: Dict = None):
        self.config = config if config is not None else {}
        self.min_workers = self.config.get("min_workers", 1)
        self.max_workers = self.config.get("max_workers", 5)
        self.cpu_threshold_up = self.config.get("cpu_threshold_up", 80) # %
        self.cpu_threshold_down = self.config.get("cpu_threshold_down", 20) # %
        self.scale_cooldown_seconds = self.config.get("scale_cooldown_seconds", 300) # 5 minutes
        self._last_scale_action_time = 0

    async def manage_scaling(self, worker_manager):
        """
        Monitors worker metrics and scales workers up or down.
        'worker_manager' is expected to have methods like get_worker_metrics(),
        add_worker(), remove_worker().
        """
        logger.info("AutoScalingManager managing scaling started.")
        while True:
            await asyncio.sleep(self.config.get("scaling_interval", 60)) # Check every minute

            if (asyncio.get_event_loop().time() - self._last_scale_action_time) < self.scale_cooldown_seconds:
                logger.debug("Scaling is in cooldown period.")
                continue

            current_workers = len(worker_manager.worker_pool)
            if not current_workers: # No workers, something is wrong, try to add min_workers
                logger.warning("No workers found in pool, attempting to re-initialize min_workers.")
                for _ in range(self.min_workers):
                    await worker_manager._create_worker(None) # Pass None for ID, manager assigns
                self._last_scale_action_time = asyncio.get_event_loop().time()
                continue

            # In a real scenario, you'd get actual load metrics from WorkerPerformanceMonitor
            # For this basic implementation, we'll use a placeholder or average of known metrics.
            # Assuming worker_manager.performance_monitor has access to actual metrics
            avg_cpu_usage = await self._get_average_worker_cpu_usage(worker_manager)
            
            logger.info(f"Current workers: {current_workers}, Average CPU usage: {avg_cpu_usage:.2f}%")

            if avg_cpu_usage > self.cpu_threshold_up and current_workers < self.max_workers:
                # Scale up
                new_workers_to_add = max(1, min(self.max_workers - current_workers, (current_workers // 2)))
                for _ in range(new_workers_to_add):
                    await worker_manager._create_worker(None) # Manager assigns ID
                logger.info(f"Scaled up: added {new_workers_to_add} workers. Total: {len(worker_manager.worker_pool)}")
                self._last_scale_action_time = asyncio.get_event_loop().time()
            elif avg_cpu_usage < self.cpu_threshold_down and current_workers > self.min_workers:
                # Scale down
                new_workers_to_remove = max(1, min(current_workers - self.min_workers, (current_workers // 4)))
                for _ in range(new_workers_to_remove):
                    worker_id_to_remove = self._select_worker_to_remove(worker_manager)
                    if worker_id_to_remove:
                        await worker_manager._remove_worker(worker_id_to_remove)
                logger.info(f"Scaled down: removed {new_workers_to_remove} workers. Total: {len(worker_manager.worker_pool)}")
                self._last_scale_action_time = asyncio.get_event_loop().time()

    async def _get_average_worker_cpu_usage(self, worker_manager) -> float:
        """
        Placeholder for getting actual average CPU usage from workers.
        In a real scenario, this would query the WorkerPerformanceMonitor
        or a metrics store like Prometheus.
        """
        if not worker_manager.worker_pool:
            return 0.0
        
        # For a more realistic stub, generate a random CPU usage
        import random
        return random.uniform(10, 90) # Simulate CPU usage

    def _select_worker_to_remove(self, worker_manager) -> Optional[int]:
        """
        Placeholder for selecting a worker to remove during scale-down.
        Could be based on least load, oldest, etc.
        """
        if worker_manager.worker_pool:
            # For now, just remove the first one
            return list(worker_manager.worker_pool.keys())[0]
        return None
