# core/security/secure_data_manager.py
class SecureDataManager:
    """
    High-level interface for secure data management
    Integrates encryption with business logic
    """
    
    def __init__(self, encryption_engine: EnterpriseEncryptionEngine):
        self.encryption_engine = encryption_engine
        self.sensitivity_classifier = DataSensitivityClassifier()
    
    async def store_sensitive_data(self, data: Dict, context: str) -> Dict:
        """
        Store sensitive data with appropriate encryption
        """
        # Classify data sensitivity
        sensitivity = await self.sensitivity_classifier.classify_data(data, context)
        
        # Serialize to bytes
        plaintext = json.dumps(data, sort_keys=True).encode('utf-8')
        
        # Encrypt with appropriate layers
        encryption_result = await self.encryption_engine.encrypt_data(
            plaintext, sensitivity
        )
        
        # Prepare storage record
        storage_record = {
            'encrypted_data': base64.b64encode(encryption_result.encrypted_data).decode(),
            'iv': base64.b64encode(encryption_result.iv).decode(),
            'tag': base64.b64encode(encryption_result.tag).decode() if encryption_result.tag else None,
            'key_id': encryption_result.key_id,
            'encryption_metadata': encryption_result.encryption_metadata,
            'sensitivity': sensitivity.value,
            'context': context,
            'timestamp': time.time()
        }
        
        return storage_record
    
    async def retrieve_sensitive_data(self, storage_record: Dict) -> Dict:
        """
        Retrieve and decrypt sensitive data
        """
        try:
            # Reconstruct encryption result
            encrypted_data = base64.b64decode(storage_record['encrypted_data'])
            iv = base64.b64decode(storage_record['iv'])
            tag = base64.b64decode(storage_record['tag']) if storage_record.get('tag') else None
            
            encryption_result = EncryptionResult(
                encrypted_data=encrypted_data,
                iv=iv,
                tag=tag,
                key_id=storage_record['key_id'],
                encryption_metadata=storage_record['encryption_metadata']
            )
            
            # Decrypt data
            plaintext = await self.encryption_engine.decrypt_data(encryption_result)
            
            # Deserialize from bytes
            data = json.loads(plaintext.decode('utf-8'))
            
            return data
            
        except Exception as e:
            logger.error(f"Data retrieval failed: {e}")
            raise DataAccessError("Failed to retrieve sensitive data") from e

class DataSensitivityClassifier:
    """
    Classify data sensitivity for appropriate encryption levels
    """
    
    def __init__(self):
        self.classification_rules = self._load_classification_rules()
    
    async def classify_data(self, data: Dict, context: str) -> SensitivityLevel:
        """Classify data sensitivity based on content and context"""
        
        sensitivity_indicators = {
            'personal_identifiers': self._contains_pii(data),
            'communication_content': self._analyze_communication_sensitivity(data),
            'intelligence_data': self._assess_intelligence_value(data, context),
            'operational_secrets': self._contains_operational_data(data)
        }
        
        # Score based on indicators
        sensitivity_score = sum(
            1 for indicator, present in sensitivity_indicators.items() if present
        )
        
        if sensitivity_score >= 3:
            return SensitivityLevel.CRITICAL
        elif sensitivity_score == 2:
            return SensitivityLevel.HIGH
        elif sensitivity_score == 1:
            return SensitivityLevel.MEDIUM
        else:
            return SensitivityLevel.LOW
    
    def _contains_pii(self, data: Dict) -> bool:
        """Check for Personally Identifiable Information"""
        pii_indicators = ['phone', 'email', 'address', 'location', 'user_id', 'username']
        data_str = json.dumps(data).lower()
        return any(indicator in data_str for indicator in pii_indicators)
