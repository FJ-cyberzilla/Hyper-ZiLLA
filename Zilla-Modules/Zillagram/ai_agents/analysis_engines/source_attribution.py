# ai_agents/analysis_engines/source_attribution.py
class SourceAttributionEngine:
    def __init__(self, db_path: str):
        self.db_path = db_path
        self.fingerprinter = MessageFingerprintingEngine()
        
    async def attribute_message_source(self, suspicious_message: Dict) -> Dict:
        """Attribute a message to its original source"""
        attribution_analysis = {
            "original_author": await self._find_original_author(suspicious_message),
            "forwarding_chain": await self._reconstruct_forwarding_chain(suspicious_message),
            "content_origin": self._analyze_content_origin(suspicious_message),
            "confidence_metrics": {},
            "supporting_evidence": []
        }
        
        # Calculate confidence scores
        attribution_analysis["confidence_metrics"] = self._calculate_attribution_confidence(
            attribution_analysis
        )
        
        return attribution_analysis
    
    async def _find_original_author(self, message: Dict) -> Dict:
        """Find the original author of the message content"""
        message_fingerprint = self.fingerprinter.create_message_fingerprint(message)
        
        # Search database for earliest occurrence
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        cursor.execute('''
            SELECT * FROM messages 
            WHERE (content_hash = ? OR semantic_hash = ?)
            AND date < ?
            ORDER BY date ASC
            LIMIT 1
        ''', (
            message_fingerprint.content_hash,
            message_fingerprint.semantic_hash, 
            message['date']
        ))
        
        original_message = cursor.fetchone()
        conn.close()
        
        if original_message:
            return {
                "user_id": original_message['user_id'],
                "chat_id": original_message['chat_id'],
                "timestamp": original_message['date'],
                "confidence": 0.95
            }
        
        # If no exact match, use semantic search
        return await self._semantic_author_search(message_fingerprint, message['date'])
    
    async def _reconstruct_forwarding_chain(self, message: Dict) -> List[Dict]:
        """Reconstruct the complete forwarding chain"""
        chain = []
        current_message = message
        
        max_iterations = 10  # Prevent infinite loops
        iteration = 0
        
        while current_message and iteration < max_iterations:
            previous_version = await self._find_previous_version(current_message)
            
            if previous_version:
                chain.append({
                    "message": previous_version,
                    "forwarder": current_message['user_id'],
                    "timestamp": current_message['date'],
                    "modifications": self._detect_modifications(previous_version, current_message)
                })
                current_message = previous_version
            else:
                # Found original
                chain.append({
                    "message": current_message,
                    "status": "original",
                    "timestamp": current_message['date']
                })
                break
            
            iteration += 1
        
        return list(reversed(chain))  # Return in chronological order
    
    def _analyze_content_origin(self, message: Dict) -> Dict:
        """Analyze where the content originally came from"""
        text = message.get('text', '')
        
        origin_analysis = {
            "external_sources": self._detect_external_sources(text),
            "common_phrases": self._check_common_phrases(text),
            "cultural_references": self._analyze_cultural_references(text),
            "technical_origin": self._analyze_technical_origin(text),
            "likely_platform": self._identify_likely_platform(text)
        }
        
        return origin_analysis
    
    def _detect_external_sources(self, text: str) -> List[Dict]:
        """Detect if content came from external sources"""
        external_indicators = []
        
        # Check for URLs
        url_pattern = r'http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\\(\\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+'
        urls = re.findall(url_pattern, text)
        
        for url in urls:
            external_indicators.append({
                "type": "url",
                "content": url,
                "platform": self._classify_url_platform(url)
            })
        
        # Check for social media patterns
        social_media_indicators = [
            ("RT @", "twitter_retweet"),
            ("via @", "twitter_via"),
            ("Shared from", "facebook_share"),
            ("Forwarded", "telegram_forward"),
            ("@instagram", "instagram_mention")
        ]
        
        for pattern, platform in social_media_indicators:
            if pattern in text:
                external_indicators.append({
                    "type": "social_media",
                    "platform": platform,
                    "confidence": 0.8
                })
        
        return external_indicators
