# installation/install_tracker.py
import json
import os
import sqlite3
from datetime import datetime
import shutil
from pathlib import Path
import platform

class InstallationTracker:
    def __init__(self, base_path: str):
        self.base_path = Path(base_path)
        self.metrics_file = self.base_path / "config" / "installation_metrics.json"
        self.ensure_directories()
        
    def ensure_directories(self):
        """Ensure required directories exist"""
        (self.base_path / "config").mkdir(parents=True, exist_ok=True)
        (self.base_path / "logs").mkdir(parents=True, exist_ok=True)
    
    def record_installation(self):
        """Record installation details"""
        install_data = {
            "installation_date": datetime.now().isoformat(),
            "system_info": self._get_system_info(),
            "installation_path": str(self.base_path),
            "python_version": platform.python_version(),
            "platform": platform.platform(),
            "user_home": str(Path.home()),
            "total_disk_space": self._get_disk_space(),
            "first_run": True,
            "sessions": [],
            "user_metrics": {
                "total_messages_processed": 0,
                "total_users_analyzed": 0,
                "total_analysis_run": 0,
                "total_runtime_seconds": 0
            }
        }
        
        self._save_metrics(install_data)
        self._create_uninstall_script()
        
    def record_session_start(self):
        """Record session start"""
        metrics = self._load_metrics()
        session_id = f"session_{len(metrics.get('sessions', [])) + 1}"
        
        session_data = {
            "session_id": session_id,
            "start_time": datetime.now().isoformat(),
            "end_time": None,
            "duration_seconds": 0,
            "messages_processed": 0,
            "users_analyzed": 0
        }
        
        metrics.setdefault("sessions", []).append(session_data)
        metrics["current_session"] = session_id
        self._save_metrics(metrics)
        
        return session_id
    
    def record_session_end(self, session_id: str, stats: Dict):
        """Record session end with statistics"""
        metrics = self._load_metrics()
        
        for session in metrics.get("sessions", []):
            if session["session_id"] == session_id:
                session["end_time"] = datetime.now().isoformat()
                start_dt = datetime.fromisoformat(session["start_time"])
                end_dt = datetime.fromisoformat(session["end_time"])
                session["duration_seconds"] = (end_dt - start_dt).total_seconds()
                
                # Update session stats
                session.update(stats)
                
                # Update global metrics
                metrics["user_metrics"]["total_messages_processed"] += stats.get("messages_processed", 0)
                metrics["user_metrics"]["total_users_analyzed"] += stats.get("users_analyzed", 0)
                metrics["user_metrics"]["total_analysis_run"] += stats.get("analysis_run", 0)
                metrics["user_metrics"]["total_runtime_seconds"] += session["duration_seconds"]
                
                break
        
        metrics.pop("current_session", None)
        self._save_metrics(metrics)
    
    def _get_system_info(self) -> Dict:
        """Get comprehensive system information"""
        return {
            "os": platform.system(),
            "os_version": platform.version(),
            "architecture": platform.architecture(),
            "processor": platform.processor(),
            "hostname": platform.node(),
            "username": os.getenv('USER') or os.getenv('USERNAME'),
            "wsl": "microsoft" in platform.release().lower(),
            "termux": "com.termux" in os.environ.get('PREFIX', '')
        }
    
    def _get_disk_space(self) -> Dict:
        """Get disk space information"""
        try:
            total, used, free = shutil.disk_usage(self.base_path)
            return {
                "total_gb": total // (2**30),
                "used_gb": used // (2**30),
                "free_gb": free // (2**30)
            }
        except:
            return {"error": "Unable to get disk space"}
    
    def _load_metrics(self) -> Dict:
        """Load installation metrics"""
        try:
            with open(self.metrics_file, 'r') as f:
                return json.load(f)
        except:
            return {}
    
    def _save_metrics(self, metrics: Dict):
        """Save installation metrics"""
        with open(self.metrics_file, 'w') as f:
            json.dump(metrics, f, indent=2)
