# core/threat_intelligence/threat_engine.py
import asyncio
from typing import Dict, List, Optional
from datetime import datetime, timedelta
import aiohttp
import json

class ThreatIntelligenceEngine:
    """
    Advanced threat intelligence with correlation to known campaigns
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.intel_sources = self._initialize_intel_sources()
        self.correlation_engine = ThreatCorrelationEngine(config)
        self.ioc_manager = IOCManager(config)
        
    def _initialize_intel_sources(self) -> Dict:
        """Initialize multiple threat intelligence sources"""
        return {
            "internal": InternalThreatFeed(),
            "commercial": CommercialIntelFeeds([
                "Recorded Future",
                "CrowdStrike", 
                "FireEye",
                "Mandiant"
            ]),
            "open_source": OSINTIntelFeeds([
                "AlienVault OTX",
                "VirusTotal",
                "ThreatConnect",
                "MISP"
            ]),
            "government": GovernmentIntelShares([
                "CISA AIS",
                "NCSC",
                "ENISA"
            ])
        }
    
    async def analyze_emerging_threats(self, user_behavior: Dict) -> Dict:
        """Correlate user behavior with known threat patterns"""
        threat_analysis = {
            "known_threat_actors": await self._check_known_actors(user_behavior),
            "behavioral_patterns": await self._analyze_threat_patterns(user_behavior),
            "campaign_correlation": await self._correlate_with_campaigns(user_behavior),
            "ioc_matches": await self._check_ioc_matches(user_behavior),
            "risk_scoring": await self._calculate_threat_risk(user_behavior)
        }
        
        return threat_analysis
    
    async def _correlate_with_campaigns(self, user_behavior: Dict) -> List[Dict]:
        """Correlate with known threat campaigns"""
        campaigns = await self._get_active_campaigns()
        
        correlated_campaigns = []
        for campaign in campaigns:
            similarity_score = self._calculate_campaign_similarity(user_behavior, campaign)
            if similarity_score > 0.7:
                correlated_campaigns.append({
                    "campaign": campaign['name'],
                    "similarity_score": similarity_score,
                    "attribution": campaign.get('attribution', 'Unknown'),
                    "tactics": campaign.get('tactics', []),
                    "indicators": campaign.get('indicators', [])
                })
        
        return correlated_campaigns
    
    async def _get_active_campaigns(self) -> List[Dict]:
        """Get active threat campaigns from all sources"""
        all_campaigns = []
        
        for source_name, source in self.intel_sources.items():
            try:
                campaigns = await source.get_active_campaigns()
                all_campaigns.extend(campaigns)
            except Exception as e:
                logger.error(f"Failed to get campaigns from {source_name}: {e}")
        
        # Deduplicate and prioritize campaigns
        return self._deduplicate_campaigns(all_campaigns)
    
    async def check_ioc_matches(self, indicators: Dict) -> Dict:
        """Check Indicators of Compromise against threat intelligence"""
        ioc_matches = {
            "ip_addresses": await self._check_ip_indicators(indicators.get('ip_addresses', [])),
            "domains": await self._check_domain_indicators(indicators.get('domains', [])),
            "hashes": await self._check_hash_indicators(indicators.get('hashes', [])),
            "behavioral_indicators": await self._check_behavioral_indicators(indicators.get('behaviors', []))
        }
        
        return {
            "matches_found": any(ioc_matches.values()),
            "detailed_matches": ioc_matches,
            "confidence_score": self._calculate_ioc_confidence(ioc_matches)
        }

class ThreatCorrelationEngine:
    """Correlate threats across multiple dimensions"""
    
    async def correlate_threat_events(self, events: List[Dict]) -> Dict:
        """Correlate multiple threat events for pattern recognition"""
        correlation_analysis = {
            "temporal_correlation": await self._analyze_temporal_patterns(events),
            "spatial_correlation": await self._analyze_geographic_patterns(events),
            "behavioral_correlation": await self._analyze_behavioral_patterns(events),
            "infrastructure_correlation": await self._analyze_infrastructure_patterns(events),
            "campaign_attribution": await self._attribute_to_campaigns(events)
        }
        
        return {
            "correlation_score": self._calculate_correlation_score(correlation_analysis),
            "analysis": correlation_analysis,
            "recommended_actions": await self._generate_threat_response(correlation_analysis)
        }
