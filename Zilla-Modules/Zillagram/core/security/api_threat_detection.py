# core/security/api_threat_detection.py
class APIThreatDetector:
    """Advanced API threat detection with machine learning"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.threat_models = self._load_threat_models()
        self.attack_patterns = self._load_attack_patterns()
        self.behavioral_analyzer = BehavioralThreatAnalyzer()
        
    async def analyze_request(self, request: Request, context: APISecurityContext) -> Dict:
        """Comprehensive API threat analysis"""
        threat_indicators = {
            "injection_attempts": await self._detect_injection_attempts(request),
            "malicious_payloads": await self._detect_malicious_payloads(request),
            "behavioral_anomalies": await self._detect_behavioral_anomalies(request, context),
            "reconnaissance_activity": await self._detect_reconnaissance(request),
            "data_exfiltration_signs": await self._detect_data_exfiltration(request)
        }
        
        # Calculate overall threat score
        threat_score = self._calculate_threat_score(threat_indicators)
        context.threat_level = threat_score
        
        return {
            "threat_detected": threat_score > self.config['thresholds']['high_risk'],
            "threat_score": threat_score,
            "indicators": threat_indicators,
            "recommended_action": self._get_recommended_action(threat_score, threat_indicators)
        }
    
    async def _detect_injection_attempts(self, request: Request) -> List[Dict]:
        """Detect SQL injection, command injection, etc."""
        injection_patterns = [
            (r'(\%27)|(\')|(\-\-)|(\%23)|(#)', 'sql_injection'),
            (r'(\|\||\/\*|\*\/|;|\-\-)', 'sql_injection_advanced'),
            (r'(union.*select)', 'union_sql_injection'),
            (r'(\;|\|\||\&\&|\`|\$\(|\$\{)', 'command_injection'),
            (r'(\.\.\/|\.\.\\|\\\.\.|\/\.\.)', 'path_traversal')
        ]
        
        detected_attempts = []
        request_data = await self._extract_request_data(request)
        
        for pattern, attack_type in injection_patterns:
            if re.search(pattern, request_data, re.IGNORECASE):
                detected_attempts.append({
                    "type": attack_type,
                    "pattern": pattern,
                    "confidence": 0.85
                })
        
        return detected_attempts
    
    async def _detect_behavioral_anomalies(self, request: Request, context: APISecurityContext) -> List[Dict]:
        """Detect behavioral anomalies using machine learning"""
        behavioral_features = {
            "request_frequency": await self._calculate_request_frequency(context.client_id),
            "endpoint_access_pattern": await self._analyze_endpoint_access(context.client_id),
            "parameter_usage": await self._analyze_parameter_usage(request),
            "temporal_pattern": await self._analyze_temporal_pattern(context.client_id),
            "geographic_consistency": await self._check_geographic_consistency(request, context)
        }
        
        # Use ML model to detect anomalies
        anomaly_score = await self.behavioral_analyzer.detect_anomalies(behavioral_features)
        
        if anomaly_score > 0.7:
            return [{
                "type": "behavioral_anomaly",
                "score": anomaly_score,
                "features": behavioral_features,
                "confidence": 0.9
            }]
        
        return []
