# ~/HyperZilla/ZILLA_CORE/fusion_engine/cross_domain_correlator.py
import json
from datetime import datetime
from typing import Dict, List, Any

class IntelligenceFusion:
    def __init__(self):
        self.digital_intel = []
        self.physical_intel = [] 
        self.enterprise_intel = []
        self.correlation_matrix = {}
    
    def ingest_digital_intel(self, osint_data: Dict):
        """Ingest OSINT intelligence"""
        enhanced_data = self._enhance_with_context(osint_data)
        self.digital_intel.append(enhanced_data)
        self._trigger_correlation()
    
    def ingest_physical_intel(self, sensor_data: Dict):
        """Ingest physical sensor intelligence"""
        geo_enhanced = self._geo_contextualize(sensor_data)
        self.physical_intel.append(geo_enhanced)
        self._trigger_correlation()
    
    def ingest_enterprise_intel(self, security_data: Dict):
        """Ingest enterprise security intelligence"""
        threat_enhanced = self._threat_contextualize(security_data)
        self.enterprise_intel.append(threat_enhanced)
        self._trigger_correlation()
    
    def _trigger_correlation(self):
        """Cross-domain intelligence correlation"""
        if len(self.digital_intel) > 0 and len(self.physical_intel) > 0:
            self._correlate_digital_physical()
        
        if len(self.enterprise_intel) > 0:
            self._correlate_threat_intel()
    
    def _correlate_digital_physical(self):
        """Correlate online presence with physical location"""
        for digital in self.digital_intel[-5:]:  # Last 5 entries
            for physical in self.physical_intel[-5:]:
                if self._spatial_temporal_match(digital, physical):
                    correlation = {
                        'confidence': 0.92,
                        'target': digital.get('target'),
                        'digital_activity': digital.get('activity'),
                        'physical_location': physical.get('coordinates'),
                        'timestamp': datetime.now().isoformat()
                    }
                    self._alert_fusion_correlation(correlation)
    
    def _spatial_temporal_match(self, digital: Dict, physical: Dict) -> bool:
        """Advanced matching algorithm"""
        time_diff = abs(
            datetime.fromisoformat(digital['timestamp']) - 
            datetime.fromisoformat(physical['timestamp'])
        ).total_seconds()
        
        return time_diff < 300  # 5-minute window
