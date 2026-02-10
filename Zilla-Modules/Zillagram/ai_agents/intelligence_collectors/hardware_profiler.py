# ai_agents/intelligence_collectors/hardware_profiler.py
import psutil
import platform
import subprocess
import re
from typing import Dict, Optional

class HardwareProfiler:
    def __init__(self):
        self.system_profile = {}
        
    def collect_hardware_intel(self) -> Dict:
        """Collect hardware intelligence for device fingerprinting"""
        return {
            "battery_profile": self._analyze_battery_behavior(),
            "wifi_fingerprint": self._analyze_wifi_patterns(),
            "system_metrics": self._collect_system_metrics(),
            "hardware_signature": self._generate_hardware_signature()
        }
    
    def _analyze_battery_behavior(self) -> Dict:
        """Analyze battery usage patterns for bot detection"""
        try:
            if hasattr(psutil, "sensors_battery"):
                battery = psutil.sensors_battery()
                if battery:
                    return {
                        "plugged": battery.power_plugged,
                        "percent": battery.percent,
                        "time_remaining": battery.secsleft if battery.secsleft != psutil.POWER_TIME_UNLIMITED else "unlimited",
                        "power_consumption_rate": self._calculate_power_consumption(),
                        "charging_pattern": self._detect_charging_pattern()
                    }
        except:
            pass
        return {"error": "Battery info unavailable"}
    
    def _analyze_wifi_patterns(self) -> Dict:
        """Analyze WiFi connectivity patterns"""
        wifi_intel = {
            "ssid_history": self._get_wifi_history(),
            "connection_stability": self._calculate_connection_stability(),
            "signal_strength_patterns": self._analyze_signal_patterns(),
            "network_switching_frequency": self._detect_network_hopping()
        }
        
        # Cross-reference with geographic patterns
        wifi_intel["suspicious_roaming"] = self._detect_suspicious_roaming(wifi_intel)
        return wifi_intel
    
    def _get_wifi_history(self) -> List[str]:
        """Extract WiFi connection history"""
        try:
            if platform.system() == "Windows":
                return self._get_windows_wifi_history()
            elif platform.system() == "Linux":
                return self._get_linux_wifi_history()
            elif platform.system() == "Darwin":
                return self._get_macos_wifi_history()
        except Exception as e:
            return [f"Error: {str(e)}"]
    
    def _get_windows_wifi_history(self) -> List[str]:
        """Extract Windows WiFi profiles"""
        try:
            result = subprocess.run(
                ["netsh", "wlan", "show", "profiles"], 
                capture_output=True, 
                text=True, 
                timeout=10
            )
            profiles = []
            for line in result.stdout.split('\n'):
                if "All User Profile" in line:
                    profiles.append(line.split(":")[1].strip())
            return profiles
        except:
            return []
