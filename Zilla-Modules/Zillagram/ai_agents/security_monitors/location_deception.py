# ai_agents/security_monitors/location_deception.py
import random
import json
from typing import Dict, Tuple, Optional
from dataclasses import dataclass
import hashlib
from geopy.geocoders import Nominatim
from geopy.distance import geodesic
import time

@dataclass
class FakeLocation:
    latitude: float
    longitude: float
    city: str
    country: str
    timezone: str
    confidence_score: float
    deception_type: str

class LocationDeceptionEngine:
    def __init__(self):
        self.geolocator = Nominatim(user_agent="telegram_osint_platform")
        self.known_spy_patterns = self._load_spy_indicators()
        self.active_deceptions = {}
        self.counter_intel_log = {}
        
    def generate_plausible_fake_location(self, real_location: Tuple[float, float], 
                                       target_region: str = None) -> FakeLocation:
        """Generate believable fake location based on real location"""
        real_lat, real_lon = real_location
        
        # Generate locations in different deception categories
        deception_strategies = {
            "urban_confusion": self._generate_urban_confusion(real_location),
            "rural_misleading": self._generate_rural_misleading(real_location),
            "border_confusion": self._generate_border_confusion(real_location),
            "tourist_trap": self._generate_tourist_trap(real_location),
            "military_base_proximity": self._generate_military_proximity(real_location)
        }
        
        # Select strategy based on context
        strategy = random.choice(list(deception_strategies.keys()))
        fake_loc = deception_strategies[strategy]
        
        # Add realistic metadata
        fake_loc = self._enhance_location_plausibility(fake_loc)
        
        return fake_loc
    
    def _generate_urban_confusion(self, real_location: Tuple[float, float]) -> FakeLocation:
        """Generate fake location in dense urban area for confusion"""
        # Major cities with similar timezones
        urban_decoys = [
            (40.7128, -74.0060, "New York", "USA", "America/New_York"),  # NYC
            (34.0522, -118.2437, "Los Angeles", "USA", "America/Los_Angeles"),  # LA
            (51.5074, -0.1278, "London", "UK", "Europe/London"),  # London
            (48.8566, 2.3522, "Paris", "France", "Europe/Paris"),  # Paris
            (35.6762, 139.6503, "Tokyo", "Japan", "Asia/Tokyo")  # Tokyo
        ]
        
        city = random.choice(urban_decoys)
        return FakeLocation(
            latitude=city[0],
            longitude=city[1],
            city=city[2],
            country=city[3],
            timezone=city[4],
            confidence_score=0.85,
            deception_type="urban_confusion"
        )
    
    def _generate_military_proximity(self, real_location: Tuple[float, float]) -> FakeLocation:
        """Generate location near military bases for intimidation"""
        military_proximity_locations = [
            (38.9517, -76.8610, "Fort Meade", "USA", "America/New_York"),  # NSA HQ
            (38.8895, -77.0353, "Pentagon", "USA", "America/New_York"),  # Pentagon
            (51.4992, -0.1246, "MI6 Headquarters", "UK", "Europe/London"),  # MI6
            (55.7558, 37.6173, "Moscow", "Russia", "Europe/Moscow"),  # FSB areas
            (31.2304, 121.4737, "Shanghai", "China", "Asia/Shanghai")  # PLA bases
        ]
        
        location = random.choice(military_proximity_locations)
        # Add slight random offset
        lat_offset = random.uniform(-0.1, 0.1)
        lon_offset = random.uniform(-0.1, 0.1)
        
        return FakeLocation(
            latitude=location[0] + lat_offset,
            longitude=location[1] + lon_offset,
            city=location[2],
            country=location[3],
            timezone=location[4],
            confidence_score=0.92,
            deception_type="military_proximity"
        )
