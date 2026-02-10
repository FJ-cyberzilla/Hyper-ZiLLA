# ai_agents/analysis_engines/shared_accounts_detector.py
import hashlib
from collections import Counter
import numpy as np
from sklearn.cluster import DBSCAN

class SharedAccountsDetector:
    def __init__(self):
        self.behavioral_fingerprints = {}
        self.writing_style_profiles = {}
        
    def detect_shared_accounts(self, user_messages: List[Dict]) -> Dict:
        """Detect if multiple people use the same account"""
        analysis = {
            "writing_style_analysis": self._analyze_writing_styles(user_messages),
            "temporal_fingerprinting": self._analyze_usage_times(user_messages),
            "behavioral_clustering": self._cluster_behavior_patterns(user_messages),
            "device_fingerprint_analysis": self._analyze_device_patterns(user_messages),
            "shared_account_probability": 0.0
        }
        
        # Calculate probability of shared account
        analysis["shared_account_probability"] = self._calculate_shared_probability(analysis)
        analysis["suspicious_sessions"] = self._identify_suspicious_sessions(user_messages)
        
        return analysis
    
    def _analyze_writing_styles(self, messages: List[Dict]) -> Dict:
        """Analyze writing style consistency"""
        style_metrics = []
        
        for msg in messages:
            style = {
                "avg_sentence_length": self._avg_sentence_length(msg['text']),
                "punctuation_ratio": self._punctuation_ratio(msg['text']),
                "capitalization_pattern": self._capitalization_analysis(msg['text']),
                "vocabulary_complexity": self._vocabulary_complexity(msg['text']),
                "emoji_usage_pattern": self._emoji_analysis(msg['text']),
                "typing_rhythm": self._analyze_typing_rhythm(msg)
            }
            style_metrics.append(style)
        
        return {
            "style_consistency": self._calculate_style_consistency(style_metrics),
            "detected_writing_styles": self._cluster_writing_styles(style_metrics),
            "style_switches": self._detect_style_switches(style_metrics)
        }
    
    def _analyze_usage_times(self, messages: List[Dict]) -> Dict:
        """Analyze usage time patterns for multiple operators"""
        timestamps = [msg['date'] for msg in messages]
        hours = [ts.hour for ts in timestamps]
        
        # Detect multiple usage patterns (different timezones/habits)
        hour_clusters = self._cluster_hours(hours)
        
        return {
            "active_hours_clusters": hour_clusters,
            "timezone_consistency": self._analyze_timezone_consistency(hours),
            "usage_burst_patterns": self._detect_usage_bursts(timestamps),
            "operator_count_estimate": len(hour_clusters)
        }
