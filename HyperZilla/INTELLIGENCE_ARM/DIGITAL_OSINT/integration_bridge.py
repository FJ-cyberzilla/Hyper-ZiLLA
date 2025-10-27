# ~/HyperZilla/INTELLIGENCE_ARM/DIGITAL_OSINT/integration_bridge.py
import sys
import os
sys.path.append('/Millitary')  # Your original stack

from core.comprehensive_orchestrator import AdvancedOrchestrator
from correlation_engine import CorrelationEngine
from evasion_engine import EvasionEngine

class DigitalOSINTBridge:
    def __init__(self):
        self.orchestrator = AdvancedOrchestrator()
        self.correlation = CorrelationEngine() 
        self.evasion = EvasionEngine()
        self.activation_status = "STANDBY"
    
    def activate_full_suite(self):
        """Activate all OSINT capabilities"""
        # Initialize evasion protocols
        self.evasion.enable_military_evasion()
        
        # Start correlation engine
        self.correlation.activate_behavioral_analysis()
        
        # Launch orchestrator
        mission_params = {
            'depth': 'deep_dive',
            'evasion_level': 'military',
            'correlation': 'cross_platform'
        }
        self.orchestrator.initialize_mission(mission_params)
        
        self.activation_status = "ACTIVE"
        return {"status": "OSINT_ACTIVE", "capabilities": "FULL_SPECTRUM"}
    
    def execute_intelligence_gathering(self, target: dict):
        """Unified intelligence collection"""
        # Evasion first
        evasion_success = self.evasion.rotate_identities()
        
        if evasion_success:
            # Collect intelligence
            intel_data = self.orchestrator.execute_deep_collection(target)
            
            # Correlate and enhance
            enhanced_intel = self.correlation.cross_reference_intel(intel_data)
            
            return {
                "status": "SUCCESS",
                "raw_intel": intel_data,
                "enhanced_intel": enhanced_intel,
                "evasion_metrics": self.evasion.get_operational_metrics()
            }
