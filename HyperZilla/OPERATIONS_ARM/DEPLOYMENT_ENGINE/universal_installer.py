# ~/HyperZilla/OPERATIONS_ARM/DEPLOYMENT_ENGINE/universal_installer.py
import platform
import subprocess
import sys
import os
from pathlib import Path

class HyperZillaInstaller:
    def __init__(self):
        self.system = platform.system().lower()
        self.architecture = platform.machine()
        self.install_path = self._get_install_path()
        
    def _get_install_path(self):
        """Get platform-specific install path"""
        paths = {
            'linux': '/opt/HyperZilla',
            'darwin': '/Applications/HyperZilla', 
            'windows': 'C:\\Program Files\\HyperZilla'
        }
        return Path(paths.get(self.system, './HyperZilla'))
    
    def detect_platform(self):
        """Comprehensive platform detection"""
        platform_info = {
            'os': self.system,
            'arch': self.architecture,
            'python_version': sys.version,
            'env': 'production'
        }
        return platform_info
    
    def install_dependencies(self):
        """Platform-specific dependency installation"""
        install_commands = {
            'linux': 'apt-get update && apt-get install -y python3-go python3-rust',
            'darwin': 'brew install go rust',
            'windows': 'choco install golang rust'
        }
        
        try:
            subprocess.run(install_commands[self.system], shell=True, check=True)
            return True
        except subprocess.CalledProcessError:
            return self._fallback_install()
    
    def deploy_components(self):
        """Deploy all HyperZilla components"""
        components = [
            'ZILLA_CORE',
            'INTELLIGENCE_ARM/DIGITAL_OSINT', 
            'INTELLIGENCE_ARM/PHYSICAL_INTEL',
            'INTELLIGENCE_ARM/ENTERPRISE_SECURITY',
            'OPERATIONS_ARM'
        ]
        
        for component in components:
            self._deploy_component(component)
            
        self._create_startup_scripts()
    
    def _deploy_component(self, component):
        """Deploy individual component with integrity check"""
        source = Path(f'./{component}')
        destination = self.install_path / component
        
        if source.exists():
            # Copy with permissions preservation
            subprocess.run(f'cp -r {source} {destination}', shell=True)
            
            # Verify integrity
            if self._verify_component_integrity(destination):
                print(f"✓ {component} deployed successfully")
            else:
                print(f"✗ {component} integrity check failed")
                self._trigger_rollback(component)

# Run installer
if __name__ == "__main__":
    installer = HyperZillaInstaller()
    installer.deploy_components()
