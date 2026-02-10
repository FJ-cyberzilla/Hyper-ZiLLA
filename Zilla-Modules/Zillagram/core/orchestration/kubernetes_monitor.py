# core/orchestration/kubernetes_monitor.py
class KubernetesMonitor:
    """Comprehensive Kubernetes cluster monitoring and optimization"""
    
    def __init__(self, orchestrator: KubernetesOrchestrator):
        self.orchestrator = orchestrator
        self.metrics_client = client.CustomObjectsApi()
        
    async def monitor_cluster_health(self):
        """Continuous cluster health monitoring"""
        while True:
            try:
                health_metrics = await self._collect_cluster_metrics()
                
                # Check node health
                node_health = await self._check_node_health()
                if not node_health['all_healthy']:
                    await self._handle_node_issues(node_health)
                
                # Check resource utilization
                resource_utilization = await self._analyze_resource_usage()
                if resource_utilization['high_usage']:
                    await self._optimize_resource_allocation(resource_utilization)
                
                # Check pod status
                pod_health = await self._check_pod_health()
                if pod_health['unhealthy_pods']:
                    await self._restart_unhealthy_pods(pod_health)
                
                await asyncio.sleep(60)  # Check every minute
                
            except Exception as e:
                logging.error(f"Cluster monitoring error: {e}")
                await asyncio.sleep(30)
    
    async def auto_scale_based_on_metrics(self):
        """Intelligent auto-scaling based on custom metrics"""
        while True:
            try:
                # Get custom metrics from Prometheus
                metrics = await self._get_custom_metrics()
                
                # Scale based on message queue length
                if metrics['queue_length'] > 1000:
                    await self._scale_workers(metrics['queue_length'] // 100)
                
                # Scale based on processing latency
                if metrics['avg_latency'] > 5.0:  # 5 seconds threshold
                    await self._scale_workers(min(20, metrics['active_workers'] * 2))
                
                # Scale down during low activity
                if metrics['queue_length'] < 100 and metrics['active_workers'] > 3:
                    await self._scale_workers(max(3, metrics['active_workers'] // 2))
                
                await asyncio.sleep(30)  # Check every 30 seconds
                
            except Exception as e:
                logging.error(f"Auto-scaling error: {e}")
                await asyncio.sleep(60)
