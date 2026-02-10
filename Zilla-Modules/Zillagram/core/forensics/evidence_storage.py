# core/forensics/evidence_storage.py
class SecureEvidenceStore:
    """
    Secure storage for forensic evidence with legal-grade protection
    Implements evidence preservation and access controls
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.encryption_engine = EnterpriseEncryptionEngine(config)
        self.access_log = EvidenceAccessLog(config)
        
    async def store_evidence(self, evidence: ForensicEvidence) -> bool:
        """
        Store evidence with comprehensive security measures
        """
        try:
            # Encrypt evidence data
            encrypted_evidence = await self._encrypt_evidence(evidence)
            
            # Generate storage metadata
            storage_metadata = {
                'storage_timestamp': datetime.utcnow().isoformat(),
                'storage_location': self._get_secure_storage_location(),
                'encryption_metadata': encrypted_evidence.encryption_metadata,
                'access_controls': await self._generate_access_controls(evidence),
                'retention_policy': await self._get_retention_policy(evidence.evidence_type)
            }
            
            # Store in secure repository
            storage_result = await self._secure_storage_write(encrypted_evidence, storage_metadata)
            
            # Log storage operation
            await self.access_log.record_storage_operation(evidence.evidence_id, storage_metadata)
            
            return storage_result
            
        except Exception as e:
            logger.error(f"Evidence storage failed: {e}")
            raise EvidenceStorageError(f"Failed to store evidence {evidence.evidence_id}") from e
    
    async def retrieve_evidence(self, evidence_id: str, requester: str, purpose: str) -> Optional[ForensicEvidence]:
        """
        Retrieve evidence with access controls and logging
        """
        # Verify access rights
        if not await self._verify_retrieval_rights(evidence_id, requester, purpose):
            raise UnauthorizedAccessError(f"Unauthorized evidence retrieval attempt by {requester}")
        
        # Retrieve encrypted evidence
        encrypted_evidence = await self._secure_storage_read(evidence_id)
        if not encrypted_evidence:
            return None
        
        # Decrypt evidence
        evidence = await self._decrypt_evidence(encrypted_evidence)
        
        # Log retrieval
        await self.access_log.record_retrieval_operation(evidence_id, requester, purpose)
        
        return evidence
    
    async def _encrypt_evidence(self, evidence: ForensicEvidence) -> Dict:
        """Encrypt evidence for secure storage"""
        # Serialize evidence
        evidence_data = {
            'evidence_object': asdict(evidence),
            'raw_data': base64.b64encode(evidence.raw_data).decode()
        }
        
        serialized_data = json.dumps(evidence_data).encode()
        
        # Encrypt with enterprise encryption
        encryption_result = await self.encryption_engine.encrypt_data(
            serialized_data,
            SensitivityLevel.CRITICAL,
            [EncryptionLayer.AT_REST, EncryptionLayer.BACKUP]
        )
        
        return {
            'evidence_id': evidence.evidence_id,
            'encrypted_data': encryption_result.encrypted_data,
            'encryption_metadata': encryption_result.encryption_metadata,
            'storage_timestamp': datetime.utcnow()
        }
