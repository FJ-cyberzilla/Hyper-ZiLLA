# installation/uninstall_manager.py
import os
import json
import shutil
import subprocess
from datetime import datetime
from pathlib import Path

class UninstallManager:
    def __init__(self, base_path: str):
        self.base_path = Path(base_path)
        self.metrics_file = self.base_path / "config" / "installation_metrics.json"
        
    def generate_uninstall_report(self) -> Dict:
        """Generate comprehensive uninstall report"""
        metrics = self._load_metrics()
        
        report = {
            "uninstallation_date": datetime.now().isoformat(),
            "installation_duration": self._calculate_installation_duration(metrics),
            "total_usage_metrics": metrics.get("user_metrics", {}),
            "session_history": metrics.get("sessions", []),
            "system_info": metrics.get("system_info", {}),
            "files_removed": self._scan_installation_files(),
            "backup_created": False
        }
        
        return report
    
    def _calculate_installation_duration(self, metrics: Dict) -> str:
        """Calculate total installation duration"""
        install_date = metrics.get("installation_date")
        if not install_date:
            return "Unknown"
        
        try:
            install_dt = datetime.fromisoformat(install_date)
            uninstall_dt = datetime.now()
            duration = uninstall_dt - install_dt
            
            days = duration.days
            hours = duration.seconds // 3600
            minutes = (duration.seconds % 3600) // 60
            
            return f"{days} days, {hours} hours, {minutes} minutes"
        except:
            return "Unknown"
    
    def _scan_installation_files(self) -> List[str]:
        """Scan and list all installation files"""
        file_list = []
        
        for root, dirs, files in os.walk(self.base_path):
            for file in files:
                if file.endswith(('.py', '.json', '.db', '.log', '.txt', '.toml')):
                    rel_path = os.path.relpath(os.path.join(root, file), self.base_path)
                    file_list.append(rel_path)
        
        return file_list
    
    def execute_uninstall(self, backup: bool = True) -> bool:
        """Execute the uninstallation process"""
        try:
            # Create backup if requested
            if backup:
                self._create_backup()
            
            # Remove installation files
            self._remove_installation_files()
            
            # Record final metrics
            self._record_uninstallation()
            
            return True
            
        except Exception as e:
            print(f"Uninstallation failed: {e}")
            return False
    
    def _create_backup(self):
        """Create backup of important data"""
        backup_dir = self.base_path.parent / f"telegram_osint_backup_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
        backup_dir.mkdir(exist_ok=True)
        
        # Backup database
        db_file = self.base_path / "telegram_osint.db"
        if db_file.exists():
            shutil.copy2(db_file, backup_dir / "telegram_osint.db")
        
        # Backup config
        config_dir = self.base_path / "config"
        if config_dir.exists():
            shutil.copytree(config_dir, backup_dir / "config", dirs_exist_ok=True)
        
        # Backup logs
        logs_dir = self.base_path / "logs"
        if logs_dir.exists():
            shutil.copytree(logs_dir, backup_dir / "logs", dirs_exist_ok=True)
        
        return backup_dir
    
    def _remove_installation_files(self):
        """Remove all installation files"""
        # Remove main directories
        directories = [
            "ai_agents", "core", "ml_models", "cli", 
            "installation", "security", "utils", "data"
        ]
        
        for dir_name in directories:
            dir_path = self.base_path / dir_name
            if dir_path.exists():
                shutil.rmtree(dir_path)
        
        # Remove individual files
        files = [
            "telegram_osint.db", "telegram_osint.log", 
            "requirements.txt", "pyproject.toml", "Makefile"
        ]
        
        for file_name in files:
            file_path = self.base_path / file_name
            if file_path.exists():
                file_path.unlink()
        
        # Remove Python cache files
        for pycache in self.base_path.rglob("__pycache__"):
            shutil.rmtree(pycache)
        for pyc_file in self.base_path.rglob("*.pyc"):
            pyc_file.unlink()
        for pyo_file in self.base_path.rglob("*.pyo"):
            pyo_file.unlink()
