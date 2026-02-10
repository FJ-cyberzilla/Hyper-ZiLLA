# core/monitoring/monitoring_engine.py
import asyncio
import time
import psutil
import platform
from typing import Dict, List, Optional
from dataclasses import dataclass
from datetime import datetime, timedelta
import statistics
from prometheus_client import Counter, Histogram, Gauge, generate_latest
import logging
import os # Added for psutil

logger = logging.getLogger(__name__)

class WorkerPerformanceMonitor:
    def __init__(self, config: Optional[Dict] = None):
        self.config = config if config is not None else {}
        self.metrics = {
            'worker_cpu_usage': Gauge('worker_cpu_usage_percent', 'Worker CPU usage percentage by ID', ['worker_id']),
            'worker_memory_usage': Gauge('worker_memory_usage_percent', 'Worker memory usage percentage by ID', ['worker_id']),
            'worker_task_queue_length': Gauge('worker_task_queue_length', 'Worker task queue length by ID', ['worker_id']),
            'worker_tasks_completed_total': Counter('worker_tasks_completed_total', 'Total tasks completed by worker ID', ['worker_id']),
            'worker_task_duration_seconds': Histogram('worker_task_duration_seconds', 'Worker task duration in seconds by ID', ['worker_id']),
        }
        self.process_metrics = {} # To store psutil.Process objects

    async def monitor_workers(self, worker_manager):
        logger.info("WorkerPerformanceMonitor monitoring started.")
        while True:
            # Iterate over a copy of worker_pool items to safely handle modifications
            for worker_id, worker_info in list(worker_manager.worker_pool.items()):
                if worker_id not in self.process_metrics:
                    try:
                        # Assuming worker_info has a PID attribute or similar to get the process
                        # For now, let's assume worker_manager.worker_pool stores worker objects
                        # and we can get pid from them or directly from psutil if workers are actual processes
                        # For the stub, we'll just simulate with current process pid
                        # In a real scenario, you'd get the actual PID of the worker process
                        self.process_metrics[worker_id] = psutil.Process(os.getpid()) # Placeholder
                        logger.info(f"Monitoring process for worker {worker_id} (PID: {os.getpid()}).")
                    except psutil.NoSuchProcess:
                        logger.warning(f"Process for worker {worker_id} not found, skipping monitoring.")
                        continue
                
                try:
                    process = self.process_metrics[worker_id]
                    cpu_percent = process.cpu_percent(interval=None) # Non-blocking
                    memory_percent = process.memory_percent()
                    
                    self.metrics['worker_cpu_usage'].labels(worker_id=worker_id).set(cpu_percent)
                    self.metrics['worker_memory_usage'].labels(worker_id=worker_id).set(memory_percent)
                    # For task queue length, task completed, task duration, it needs integration with worker's internal state
                    # These metrics would typically be pushed by the worker itself, or pulled via a more sophisticated IPC.

                except psutil.NoSuchProcess:
                    logger.warning(f"Process for worker {worker_id} disappeared, removing from monitoring.")
                    if worker_id in self.process_metrics:
                        del self.process_metrics[worker_id]
                except Exception as e:
                    logger.error(f"Error monitoring worker {worker_id}: {e}")
            
            await asyncio.sleep(self.config.get("monitor_interval", 10)) # Monitor every 10 seconds

@dataclass
class SystemMetrics:
    timestamp: datetime
    cpu_percent: float
    memory_usage: float
    disk_usage: float
    network_io: Dict
    process_count: int
    system_load: List[float]

class EnterpriseMonitoring:
    """
    Comprehensive monitoring for system, network, application, and business metrics
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.metrics_collector = MetricsCollector(config)
        self.alert_manager = AlertManager(config)
        self.performance_analyzer = PerformanceAnalyzer(config)
        self.dashboard_reporter = DashboardReporter(config)
        
        # Prometheus metrics
        self.metrics = self._initialize_prometheus_metrics()
        
    def _initialize_prometheus_metrics(self) -> Dict:
        """Initialize Prometheus metrics for monitoring"""
        return {
            'requests_total': Counter('http_requests_total', 'Total HTTP requests', ['method', 'endpoint', 'status']),
            'request_duration': Histogram('http_request_duration_seconds', 'HTTP request duration'),
            'system_cpu': Gauge('system_cpu_usage', 'System CPU usage percentage'),
            'system_memory': Gauge('system_memory_usage', 'System memory usage percentage'),
            'active_connections': Gauge('active_connections', 'Active network connections'),
            'queue_size': Gauge('task_queue_size', 'Current task queue size'),
            'error_rate': Gauge('error_rate', 'Application error rate')
        }
    
    async def start_comprehensive_monitoring(self):
        """Start all monitoring systems"""
        monitoring_tasks = [
            self._monitor_system_resources(),
            self._monitor_network_activity(),
            self._monitor_application_performance(),
            self._monitor_business_metrics(),
            self._monitor_security_events(),
            self._monitor_compliance_metrics()
        ]
        
        for task in monitoring_tasks:
            asyncio.create_task(task)
        
        logger.info("📊 Comprehensive monitoring system started")
    
    async def _monitor_system_resources(self):
        """Monitor system resource usage with predictive analytics"""
        while True:
            try:
                metrics = await self.metrics_collector.collect_system_metrics()
                
                # Update Prometheus metrics
                self.metrics['system_cpu'].set(metrics.cpu_percent)
                self.metrics['system_memory'].set(metrics.memory_usage)
                
                # Check thresholds and trigger alerts
                await self._check_system_thresholds(metrics)
                
                # Predictive capacity planning
                await self._analyze_capacity_trends(metrics)
                
                await asyncio.sleep(30)  # Collect every 30 seconds
                
            except Exception as e:
                logger.error(f"System monitoring error: {e}")
                await asyncio.sleep(60)
    
    async def _monitor_network_activity(self):
        """Monitor network performance and security"""
        while True:
            try:
                network_metrics = await self.metrics_collector.collect_network_metrics()
                
                # Monitor connection patterns
                connection_analysis = await self._analyze_network_patterns(network_metrics)
                
                # Detect network anomalies
                anomalies = await self._detect_network_anomalies(network_metrics)
                if anomalies:
                    await self.alert_manager.trigger_alert("NETWORK_ANOMALY", anomalies)
                
                # Monitor bandwidth usage
                await self._analyze_bandwidth_usage(network_metrics)
                
                await asyncio.sleep(60)  # Check every minute
                
            except Exception as e:
                logger.error(f"Network monitoring error: {e}")
                await asyncio.sleep(60)
    
    async def _monitor_application_performance(self):
        """Monitor application performance and health"""
        while True:
            try:
                app_metrics = await self.metrics_collector.collect_application_metrics()
                
                # Track key performance indicators
                kpis = {
                    'response_time': app_metrics.get('avg_response_time', 0),
                    'throughput': app_metrics.get('requests_per_second', 0),
                    'error_rate': app_metrics.get('error_percentage', 0),
                    'concurrent_users': app_metrics.get('active_users', 0)
                }
                
                # Update Prometheus metrics
                self.metrics['error_rate'].set(kpis['error_rate'])
                
                # Check performance thresholds
                if kpis['response_time'] > 2.0:  # 2 seconds threshold
                    await self.alert_manager.trigger_alert("HIGH_RESPONSE_TIME", kpis)
                
                if kpis['error_rate'] > 1.0:  # 1% error rate threshold
                    await self.alert_manager.trigger_alert("HIGH_ERROR_RATE", kpis)
                
                # Performance trend analysis
                await self._analyze_performance_trends(app_metrics)
                
                await asyncio.sleep(30)  # Check every 30 seconds
                
            except Exception as e:
                logger.error(f"Application monitoring error: {e}")
                await asyncio.sleep(60)
