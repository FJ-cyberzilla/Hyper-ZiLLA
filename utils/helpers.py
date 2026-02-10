"""
General Helper Functions
Common utilities used across Hyper-ZiLLA
"""

import json
import time
from typing import Any, Dict, List, Optional
from pathlib import Path
import functools


def timer(func):
    """Decorator to measure function execution time"""

    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.perf_counter()
        result = func(*args, **kwargs)
        end_time = time.perf_counter()
        print(f"Function {func.__name__} took {end_time - start_time:.4f} seconds")
        return result

    return wrapper


def save_json(data: Any, filepath: Path) -> bool:
    """Save data to JSON file with error handling"""
    try:
        with open(filepath, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        return True
    except Exception as e:
        print(f"Error saving JSON to {filepath}: {e}")
        return False


def load_json(filepath: Path) -> Optional[Any]:
    """Load data from JSON file with error handling"""
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading JSON from {filepath}: {e}")
        return None


def chunk_list(lst: List, chunk_size: int) -> List[List]:
    """Split list into chunks of specified size"""
    return [lst[i : i + chunk_size] for i in range(0, len(lst), chunk_size)]


def get_file_size(filepath: Path) -> int:
    """Get file size in bytes"""
    return filepath.stat().st_size if filepath.exists() else 0


class PerformanceMonitor:
    """Simple performance monitoring"""

    def __init__(self):
        self.metrics = {}

    def start_timer(self, name: str):
        """Start timer for named operation"""
        self.metrics[name] = {"start": time.perf_counter()}

    def stop_timer(self, name: str) -> float:
        """Stop timer and return duration"""
        if name in self.metrics and "start" in self.metrics[name]:
            duration = time.perf_counter() - self.metrics[name]["start"]
            self.metrics[name]["duration"] = duration
            self.metrics[name]["end"] = time.perf_counter()
            return duration
        return 0.0

    def get_metrics(self) -> Dict[str, Any]:
        """Get all performance metrics"""
        return self.metrics.copy()


# Example usage
if __name__ == "__main__":
    # Test helpers
    monitor = PerformanceMonitor()

    monitor.start_timer("test_operation")
    time.sleep(0.1)
    duration = monitor.stop_timer("test_operation")

    print(f"Operation took: {duration:.4f} seconds")
    print(f"All metrics: {monitor.get_metrics()}")
