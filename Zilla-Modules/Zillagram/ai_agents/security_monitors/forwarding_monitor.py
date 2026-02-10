# ai_agents/security_monitors/forwarding_monitor.py
import asyncio
from datetime import datetime, timedelta

class RealTimeForwardingMonitor:
    def __init__(self, correlation_engine):
        self.correlation_engine = correlation_engine
        self.suspected_leaks = {}
        self.alert_threshold = 0.7
        
    async def monitor_conversation_leaks(self, chat_id: int, user_id: int):
        """Monitor specific conversation for leakage"""
        while True:
            recent_messages = await self._get_recent_messages(chat_id, hours=24)
            
            for message in recent_messages:
                leak_analysis = await self.analyze_message_leakage(message)
                
                if leak_analysis['leak_confidence'] > self.alert_threshold:
                    await self._trigger_leak_alert(leak_analysis)
                    
                    # Take countermeasures
                    await self._implement_countermeasures(leak_analysis)
            
            await asyncio.sleep(300)  # Check every 5 minutes
    
    async def analyze_message_leakage(self, message: Dict) -> Dict:
        """Analyze if a message has been leaked/forwarded"""
        propagation_analysis = await self.correlation_engine.track_message_propagation(message)
        
        leak_indicators = {
            "cross_chat_appearance": len(propagation_analysis.get('similar_messages', [])) > 1,
            "temporal_proximity": self._check_temporal_proximity(propagation_analysis),
            "content_modification": self._detect_content_modification(propagation_analysis),
            "trust_network_violation": self._check_trust_network(message, propagation_analysis)
        }
        
        return {
            "message_id": message['id'],
            "leak_confidence": self._calculate_leak_confidence(leak_indicators),
            "leak_indicators": leak_indicators,
            "propagation_path": propagation_analysis.get('propagation_path', []),
            "suspected_leaker": self._identify_primary_leaker(propagation_analysis),
            "recommended_actions": self._suggest_countermeasures(leak_indicators)
        }
    
    def _check_trust_network(self, original_message: Dict, propagation: Dict) -> bool:
        """Check if forwarding violates expected trust relationships"""
        original_chat = original_message['chat_id']
        original_user = original_message['user_id']
        
        for similar_msg in propagation.get('similar_messages', []):
            if similar_msg['chat_id'] != original_chat:
                # Check if this cross-chat sharing is suspicious
                relationship_score = self._calculate_relationship_score(
                    original_user, similar_msg['user_id']
                )
                
                if relationship_score < 0.3:  # Low relationship score
                    return True
        
        return False
