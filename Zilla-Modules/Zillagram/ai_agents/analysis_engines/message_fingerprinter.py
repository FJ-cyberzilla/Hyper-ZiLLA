# ai_agents/analysis_engines/message_fingerprinter.py
import hashlib
import json
from typing import Dict, List, Optional
from dataclasses import dataclass
import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
import re

@dataclass
class MessageFingerprint:
    content_hash: str
    semantic_hash: str
    style_signature: str
    temporal_pattern: str
    metadata_hash: str
    composite_fingerprint: str

class MessageFingerprintingEngine:
    def __init__(self):
        self.vectorizer = TfidfVectorizer(analyzer='char', ngram_range=(2, 4))
        self.known_forwarded_patterns = self._load_forwarding_indicators()
        
    def create_message_fingerprint(self, message: Dict) -> MessageFingerprint:
        """Create comprehensive fingerprint for message tracking"""
        text = message.get('text', '')
        
        return MessageFingerprint(
            content_hash=self._generate_content_hash(text),
            semantic_hash=self._generate_semantic_hash(text),
            style_signature=self._analyze_writing_style(text),
            temporal_pattern=self._analyze_temporal_pattern(message),
            metadata_hash=self._generate_metadata_hash(message),
            composite_fingerprint=self._generate_composite_fingerprint(message)
        )
    
    def _generate_content_hash(self, text: str) -> str:
        """Generate content-based hash resistant to minor modifications"""
        # Normalize text
        normalized = self._normalize_text(text)
        
        # Create multiple hash versions for robustness
        hashes = [
            hashlib.sha256(normalized.encode()).hexdigest(),
            hashlib.sha256(text.lower().encode()).hexdigest(),
            self._generate_ngram_hash(text, n=3)
        ]
        
        return hashlib.sha256(''.join(hashes).encode()).hexdigest()[:32]
    
    def _generate_semantic_hash(self, text: str) -> str:
        """Generate semantic hash that captures meaning"""
        # Extract key semantic elements
        semantic_elements = [
            self._extract_key_entities(text),
            self._extract_sentiment_profile(text),
            self._extract_topic_signature(text),
            self._extract_syntactic_pattern(text)
        ]
        
        semantic_string = json.dumps(semantic_elements, sort_keys=True)
        return hashlib.sha256(semantic_string.encode()).hexdigest()[:32]
    
    def _analyze_writing_style(self, text: str) -> str:
        """Analyze writing style for authorship attribution"""
        style_metrics = {
            "avg_sentence_length": self._calculate_avg_sentence_length(text),
            "punctuation_ratio": self._calculate_punctuation_ratio(text),
            "capitalization_pattern": self._analyze_capitalization(text),
            "vocabulary_richness": self._calculate_vocabulary_richness(text),
            "emoji_usage_pattern": self._analyze_emoji_usage(text),
            "typing_rhythm": self._estimate_typing_rhythm(text)
        }
        
        return hashlib.sha256(json.dumps(style_metrics, sort_keys=True).encode()).hexdigest()[:32]
    
    def detect_forwarded_content(self, message: Dict, chat_context: List[Dict]) -> Dict:
        """Detect if message contains forwarded content"""
        fingerprint = self.create_message_fingerprint(message)
        
        analysis = {
            "direct_forward_detection": self._check_direct_forward_indicators(message),
            "content_similarity_analysis": self._analyze_content_similarity(message, chat_context),
            "temporal_anomalies": self._detect_temporal_anomalies(message, chat_context),
            "style_inconsistencies": self._detect_style_inconsistencies(message, chat_context),
            "cross_chat_correlation": self._correlate_across_chats(fingerprint)
        }
        
        analysis["forward_probability"] = self._calculate_forward_probability(analysis)
        analysis["potential_sources"] = self._identify_potential_sources(fingerprint)
        
        return analysis
