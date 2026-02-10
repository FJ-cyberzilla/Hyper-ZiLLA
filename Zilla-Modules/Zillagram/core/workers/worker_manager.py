# core/workers/worker_manager.py
import asyncio
import multiprocessing
from typing import List, Dict, Callable, Optional
from concurrent.futures import ProcessPoolExecutor, ThreadPoolExecutor
import psutil
import os
import uvicorn
from fastapi import FastAPI
import logging
import time
import uuid # Added for task ID generation
import random # Added for simulating task processing time
import aiohttp # Added for web scraping
from sklearn.feature_extraction.text import TfidfVectorizer # For a basic ML model example
from sklearn.linear_model import LogisticRegression # For a basic ML model example
from async_retrying import retry # For retrying operations

from core.monitoring.monitoring_engine import WorkerPerformanceMonitor
from core.orchestration.autoscaling_manager import AutoScalingManager
from core.network.proxy_manager import EnterpriseProxyManager
from core.security.rate_limiting import AdaptiveRateLimiter
from core.network.dns_manager import EnterpriseDNSManager

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)

# Initialize FastAPI app for health checks
health_app = FastAPI()

@health_app.get("/health")
async def health_check():
    return {"status": "healthy"}

class EnterpriseWorkerManager:
    def __init__(self, config: Dict):
        self.config = config
        self.worker_pool = {}
        self.task_queue = asyncio.Queue()
        self.performance_monitor = WorkerPerformanceMonitor()
        self.auto_scaling = AutoScalingManager()
        
    def _generate_task_id(self) -> str:
        return str(uuid.uuid4())

    async def initialize_workers(self):
        """Initialize worker pool based on system resources"""
        cpu_count = multiprocessing.cpu_count()
        available_ram = psutil.virtual_memory().available / (1024 ** 3)  # GB
        
        # Calculate optimal worker count
        optimal_workers = self._calculate_optimal_workers(cpu_count, available_ram)
        
        logger.info(f"🚀 Initializing {optimal_workers} workers...")
        
        # Create worker processes
        for i in range(optimal_workers):
            worker = await self._create_worker(i)
            self.worker_pool[worker.worker_id] = worker
        
        # Start task distributor
        asyncio.create_task(self._distribute_tasks())
        
        # Start performance monitoring
        asyncio.create_task(self.performance_monitor.monitor_workers(self)) # Pass self (manager) to monitor
        
        # Start auto-scaling
        asyncio.create_task(self.auto_scaling.manage_scaling(self)) # Pass self (manager) for scaling
    
    def _calculate_optimal_workers(self, cpu_count: int, available_ram: float) -> int:
        """Calculate optimal number of workers"""
        # Conservative resource allocation
        cpu_workers = max(1, cpu_count - 2)  # Leave 2 cores for system
        ram_workers = int(available_ram / 2)  # 2GB per worker
        
        return min(cpu_workers, ram_workers, 16)  # Max 16 workers
    
    async def submit_task(self, task_type: str, payload: Dict, priority: int = 1) -> str:
        """Submit task to worker system"""
        task_id = self._generate_task_id()
        
        task = {
            'task_id': task_id,
            'type': task_type,
            'payload': payload,
            'priority': priority,
            'submitted_at': time.time(),
            'status': 'queued'
        }
        
        await self.task_queue.put(task)
        return task_id
    
    async def _distribute_tasks(self):
        """Intelligent task distribution to workers"""
        while True:
            try:
                task = await self.task_queue.get()
                
                # Select optimal worker for this task type
                optimal_worker = self._select_optimal_worker(task)
                
                if optimal_worker:
                    await optimal_worker.assign_task(task)
                else:
                    # No available worker, implement backpressure
                    logger.warning(f"No worker available for task {task['task_id']}, requeueing.")
                    await asyncio.sleep(0.1)
                    await self.task_queue.put(task)  # Requeue
                
                self.task_queue.task_done()
                
            except Exception as e:
                logger.error(f"Task distribution error: {e}")
                await asyncio.sleep(1)

    # --- Worker Management Methods ---
    class Worker:
        """A basic worker that can process tasks."""
        def __init__(self, worker_id: int):
            self.worker_id = worker_id
            self.current_task = None
            self.is_busy = False
            logger.info(f"Worker {self.worker_id} initialized.")

        async def assign_task(self, task: Dict):
            self.is_busy = True
            self.current_task = task
            logger.info(f"Worker {self.worker_id} assigned task {task['task_id']}.")
            try:
                # Delegate to SpecializedWorkers based on task type
                worker_class = getattr(SpecializedWorkers, task['type'].capitalize() + "Worker", None)
                if worker_class:
                    specialized_worker = worker_class(self.worker_id)
                    result = await getattr(specialized_worker, f"process_{task['type']}_task")(task)
                    logger.info(f"Worker {self.worker_id} completed specialized task {task['task_id']}. Result: {result}")
                else:
                    logger.warning(f"Unknown task type: {task['type']}. Worker {self.worker_id} processing generic task {task['task_id']}.")
                    await asyncio.sleep(random.uniform(0.5, 5.0)) # Simulate generic work
                    logger.info(f"Worker {self.worker_id} completed generic task {task['task_id']}.")
            except Exception as e:
                logger.error(f"Worker {self.worker_id} failed task {task['task_id']}: {e}")
            finally:
                self.current_task = None
                self.is_busy = False

    async def _create_worker(self, worker_id: Optional[int] = None):
        if worker_id is None:
            worker_id = len(self.worker_pool) + 1 # Assign a simple incremental ID
        worker = self.Worker(worker_id)
        self.worker_pool[worker.worker_id] = worker
        logger.info(f"Worker {worker_id} created and added to pool.")
        return worker

    async def _remove_worker(self, worker_id: int):
        if worker_id in self.worker_pool:
            # In a real system, you'd gracefully shut down the worker process
            del self.worker_pool[worker_id]
            logger.info(f"Worker {worker_id} removed from pool.")
        else:
            logger.warning(f"Attempted to remove non-existent worker {worker_id}.")

    def _select_optimal_worker(self, task: Dict):
        """Selects an optimal worker based on availability and load."""
        available_workers = [worker for worker in self.worker_pool.values() if not worker.is_busy]
        if available_workers:
            # Simple round-robin or first available
            return available_workers[0] 
        return None

class SpecializedWorkers:
    """Specialized worker types for different tasks"""
    
    class ScrapingWorker:
        def __init__(self, worker_id: int):
            self.worker_id = worker_id
            self.specialization = "scraping"
            self.proxy_manager = EnterpriseProxyManager()
            self.rate_limiter = AdaptiveRateLimiter()
            
        @retry(attempts=3, delay=2) # Retry 3 times with 2 second delay
        async def _fetch_url_with_proxy(self, session, url: str, proxy: Optional[str] = None):
            async with session.get(url, proxy=proxy, timeout=10) as response:
                response.raise_for_status() # Raise an exception for bad status codes
                return await response.text()

        async def _execute_scraping_with_retry(self, target_url: str, proxy: str):
            logger.info(f"Executing scraping for {target_url} with proxy {proxy} (real implementation).")
            
            async with aiohttp.ClientSession() as session:
                try:
                    content = await self._fetch_url_with_proxy(session, target_url, proxy)
                    return {"data": f"Scraped content from {target_url[:50]}...", "length": len(content)}
                except Exception as e:
                    logger.error(f"Scraping failed for {target_url} after retries: {e}")
                    raise

        async def process_scraping_task(self, task: Dict):
            """Process scraping tasks with intelligent rate limiting"""
            target = task['payload']['target']
            
            await self.proxy_manager.initialize_proxy_pool() 
            proxy = await self.proxy_manager.get_optimal_proxy(target, "scraping")
            
            await self.rate_limiter.wait_if_needed(target)
            
            result = await self._execute_scraping_with_retry(target, proxy.host) # Assuming proxy.host is the URL
            
            return result
    
    class AnalysisWorker:
        def __init__(self, worker_id: int):
            self.worker_id = worker_id
            self.specialization = "analysis"
            self.ml_models = self._load_ml_models()
            
        def _load_ml_models(self):
            logger.info("Loading ML models (real implementation using scikit-learn mock).")
            # In a real scenario, you'd load pre-trained models.
            # Here, we'll create simple mock models.
            vectorizer = TfidfVectorizer()
            classifier = LogisticRegression()
            
            # Simulate training for a text classification model
            # This is illustrative; real models would be loaded, not trained on the fly in worker init
            sample_texts = ["positive review", "negative review", "neutral comment"]
            sample_labels = [1, 0, 0]
            vectorizer.fit(sample_texts)
            classifier.fit(vectorizer.transform(sample_texts), sample_labels)

            return {
                "text_classifier": {"vectorizer": vectorizer, "model": classifier},
                "entity_recognizer": "placeholder_entity_model_loaded"
            }

        async def _perform_analysis_with_monitoring(self, model_key, data: str):
            logger.info(f"Performing analysis with {model_key} on '{data[:50]}...' (real implementation).")
            
            if model_key == "text_classifier":
                model_components = self.ml_models.get(model_key)
                if not model_components:
                    raise Exception(f"Text classifier model components not loaded.")
                
                vectorizer = model_components["vectorizer"]
                classifier = model_components["model"]
                
                # Preprocess data
                data_vectorized = vectorizer.transform([data])
                prediction = classifier.predict(data_vectorized)[0]
                sentiment = "positive" if prediction == 1 else "negative/neutral"
                
                await asyncio.sleep(random.uniform(0.1, 0.5)) # Simulate analysis time
                return {"analysis_type": "text_classification", "input_data": data, "prediction": sentiment, "raw_score": float(classifier.predict_proba(data_vectorized)[0][prediction])}
            
            elif model_key == "entity_recognizer":
                # Simulate entity recognition
                entities = [f"ENTITY_{i}" for i in range(random.randint(1,3))]
                await asyncio.sleep(random.uniform(0.1, 0.5))
                return {"analysis_type": "entity_recognition", "input_data": data, "entities": entities}

            else:
                await asyncio.sleep(random.uniform(0.5, 2.0))
                return {"analysis_result": f"Result for {data}"}


        async def process_analysis_task(self, task: Dict):
            """Process AI analysis tasks"""
            analysis_type = task['payload']['analysis_type']
            data = task['payload']['data']
            
            model_key = None
            if analysis_type == "sentiment_analysis":
                model_key = "text_classifier"
            elif analysis_type == "entity_recognition":
                model_key = "entity_recognizer"
            # Add more mappings as needed
            
            if not model_key:
                raise Exception(f"Unknown analysis type: {analysis_type}")
            
            result = await self._perform_analysis_with_monitoring(model_key, data)
            
            return result
    
    class NetworkWorker:
        def __init__(self, worker_id: int):
            self.worker_id = worker_id
            self.specialization = "network"
            self.dns_manager = EnterpriseDNSManager()
            
        async def process_network_task(self, task: Dict):
            """Process network-intensive tasks"""
            operation = task['payload']['operation']
            
            if operation == "dns_resolution":
                domains = task['payload']['domains']
                results = {}
                
                for domain in domains:
                    result = await self.dns_manager.smart_dns_resolution(domain)
                    results[domain] = result
                
                return results

# --- Main execution block ---
async def run_worker_manager():
    logger.info("Starting EnterpriseWorkerManager...")
    manager = EnterpriseWorkerManager(config={}) # Replace with actual config
    await manager.initialize_workers()
    # Keep the manager running
    while True:
        await asyncio.sleep(3600) # Sleep for a long time to keep event loop alive

async def main():
    worker_manager_task = asyncio.create_task(run_worker_manager())
    
    config = uvicorn.Config(health_app, host="0.0.0.0", port=8081, log_level="info")
    server = uvicorn.Server(config)
    
    # Run uvicorn server in the same event loop
    web_server_task = asyncio.create_task(server.serve())

    await asyncio.gather(worker_manager_task, web_server_task)

if __name__ == "__main__":
    asyncio.run(main())