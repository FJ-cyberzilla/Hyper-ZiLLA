# ai_agents/analysis_engines/timezone_detector.py
import pytz
from datetime import datetime, timedelta
import asyncio
from typing import Dict, List, Optional
import sqlite3
import statistics

class EnterpriseTimezoneDetector:
    def __init__(self, db_path: str):
        self.db_path = db_path
        self.timezone_patterns = self._load_timezone_patterns()
        self.clock_anomaly_threshold = 0.85
        
    async def detect_timezone_manipulation(self, user_id: int, current_message: Dict) -> Dict:
        """Enterprise-grade timezone manipulation detection"""
        user_messages = await self._get_user_message_history(user_id)
        
        analysis = {
            "timezone_consistency": await self._analyze_timezone_consistency(user_messages),
            "response_time_analysis": self._analyze_response_timing(user_messages),
            "activity_pattern_correlation": self._correlate_activity_patterns(user_messages),
            "clock_skew_detection": await self._detect_clock_skew(user_id, current_message),
            "geolocation_timezone_mismatch": await self._detect_geo_timezone_mismatch(user_id)
        }
        
        analysis["manipulation_confidence"] = self._calculate_manipulation_confidence(analysis)
        analysis["recommended_action"] = self._suggest_action(analysis)
        
        return analysis
    
    async def _analyze_timezone_consistency(self, messages: List[Dict]) -> Dict:
        """Analyze timezone patterns for consistency"""
        if len(messages) < 10:
            return {"confidence": 0.0, "pattern": "insufficient_data"}
        
        timezones = []
        for msg in messages:
            tz_info = await self._extract_timezone_from_message(msg)
            if tz_info:
                timezones.append(tz_info)
        
        if not timezones:
            return {"confidence": 0.0, "pattern": "no_timezone_data"}
        
        # Calculate timezone consistency
        primary_tz = statistics.mode([tz['offset'] for tz in timezones])
        consistency_score = sum(1 for tz in timezones if tz['offset'] == primary_tz) / len(timezones)
        
        return {
            "primary_timezone": primary_tz,
            "consistency_score": consistency_score,
            "timezone_changes": len(set(tz['offset'] for tz in timezones)),
            "pattern": "consistent" if consistency_score > 0.9 else "suspicious"
        }
    
    def _analyze_response_timing(self, messages: List[Dict]) -> Dict:
        """Analyze response timing for human vs bot patterns"""
        response_times = []
        
        for i in range(1, len(messages)):
            time_diff = messages[i]['date'] - messages[i-1]['date']
            response_times.append(time_diff.total_seconds())
        
        if not response_times:
            return {"pattern": "insufficient_data", "confidence": 0.0}
        
        # Statistical analysis
        avg_response = statistics.mean(response_times)
        std_dev = statistics.stdev(response_times) if len(response_times) > 1 else 0
        
        # Human patterns: varied response times, some quick, some slow
        # Bot patterns: consistent timing, immediate responses
        variability_score = std_dev / avg_response if avg_response > 0 else 0
        
        return {
            "average_response_seconds": avg_response,
            "response_std_dev": std_dev,
            "variability_score": variability_score,
            "pattern": "human" if variability_score > 0.3 else "bot_like",
            "confidence": min(variability_score * 3, 1.0)  # Scale to 0-1
        }
    
    async def _detect_clock_skew(self, user_id: int, current_message: Dict) -> Dict:
        """Detect system clock manipulation"""
        # Compare user's claimed time with network time
        user_time = current_message['date']
        actual_time = datetime.now(pytz.UTC)
        
        time_diff = (user_time - actual_time).total_seconds()
        
        # Check if difference is suspicious
        skew_analysis = {
            "claimed_time": user_time.isoformat(),
            "actual_time": actual_time.isoformat(),
            "time_difference_seconds": time_diff,
            "suspicious": abs(time_diff) > 300,  # More than 5 minutes difference
            "skew_direction": "ahead" if time_diff > 0 else "behind"
        }
        
        # Check for patterns of gradual time changes (simulating travel)
        historical_skew = await self._analyze_historical_skew(user_id)
        skew_analysis.update(historical_skew)
        
        return skew_analysis
    
    async def _detect_geo_timezone_mismatch(self, user_id: int) -> Dict:
        """Detect mismatch between claimed location and timezone"""
        # This would integrate with IP geolocation and other location signals
        user_location = await self._get_user_location_data(user_id)
        user_timezone = await self._get_user_timezone(user_id)
        
        if not user_location or not user_timezone:
            return {"detectable": False, "confidence": 0.0}
        
        # Check if timezone matches geographic location
        expected_tz = self._get_timezone_for_location(user_location)
        mismatch = expected_tz != user_timezone
        
        return {
            "detectable": True,
            "claimed_timezone": user_timezone,
            "expected_timezone": expected_tz,
            "mismatch": mismatch,
            "confidence": 0.8 if mismatch else 0.2
        }
