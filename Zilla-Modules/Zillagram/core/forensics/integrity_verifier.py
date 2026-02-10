# core/forensics/integrity_verifier.py
class EvidenceIntegrityVerifier:
    """
    Continuous evidence integrity verification for legal admissibility
    Implements cryptographic proof of evidence authenticity
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.verification_log = IntegrityVerificationLog(config)
        self.tamper_detector = TamperDetectionEngine(config)
        
    async def verify_evidence_integrity(self, evidence_id: str) -> Dict:
        """
        Comprehensive evidence integrity verification
        """
        evidence = await self.evidence_store.retrieve_evidence(evidence_id)
        if not evidence:
            raise EvidenceNotFoundError(f"Evidence {evidence_id} not found")
        
        verification_checks = {
            'content_hash_verification': await self._verify_content_hash(evidence),
            'digital_signature_verification': await self._verify_digital_signature(evidence),
            'tamper_detection': await self.tamper_detector.analyze_for_tampering(evidence),
            'metadata_consistency': await self._verify_metadata_consistency(evidence),
            'timestamp_validation': await self._verify_timestamps(evidence),
            'chain_of_custody_integrity': await self._verify_chain_integrity(evidence)
        }
        
        verification_result = {
            'evidence_id': evidence_id,
            'verification_timestamp': datetime.utcnow(),
            'verification_checks': verification_checks,
            'overall_integrity': all(verification_checks.values()),
            'integrity_score': self._calculate_integrity_score(verification_checks),
            'legal_admissibility': await self._assess_legal_admissibility(verification_checks)
        }
        
        # Log verification for audit trail
        await self.verification_log.record_verification(verification_result)
        
        return verification_result
    
    async def _verify_content_hash(self, evidence: ForensicEvidence) -> bool:
        """Verify evidence content hasn't changed using multiple hash algorithms"""
        current_hashes = json.loads(evidence.content_hash)
        
        # Calculate current hashes
        recalculated_hashes = {
            'md5': hashlib.md5(evidence.raw_data).hexdigest(),
            'sha1': hashlib.sha1(evidence.raw_data).hexdigest(),
            'sha256': hashlib.sha256(evidence.raw_data).hexdigest(),
            'sha512': hashlib.sha512(evidence.raw_data).hexdigest()
        }
        
        # Compare with original hashes
        for algorithm in current_hashes.keys():
            if current_hashes[algorithm] != recalculated_hashes[algorithm]:
                logger.error(f"Hash mismatch for {algorithm} on evidence {evidence.evidence_id}")
                return False
        
        return True
    
    async def _verify_digital_signature(self, evidence: ForensicEvidence) -> bool:
        """Verify digital signature for evidence authenticity"""
        signature_data = evidence.metadata.get('digital_signature', {})
        if not signature_data:
            logger.warning(f"No digital signature found for evidence {evidence.evidence_id}")
            return False
        
        try:
            # Recreate signed data
            verification_data = {
                'evidence_id': evidence.evidence_id,
                'content_hash': evidence.content_hash,
                'timestamp': evidence.collection_timestamp.isoformat(),
                'collector': evidence.collected_by
            }
            
            # Verify signature
            signature = base64.b64decode(signature_data['signature'])
            public_key = self._load_verification_certificate()
            
            public_key.verify(
                signature,
                json.dumps(verification_data, sort_keys=True).encode(),
                padding.PSS(
                    mgf=padding.MGF1(hashes.SHA256()),
                    salt_length=padding.PSS.MAX_LENGTH
                ),
                hashes.SHA256()
            )
            
            return True
            
        except Exception as e:
            logger.error(f"Digital signature verification failed: {e}")
            return False
    
    async def continuous_integrity_monitoring(self):
        """Continuous monitoring of all evidence integrity"""
        while True:
            try:
                # Get all active evidence
                active_evidence = await self.evidence_store.get_active_evidence()
                
                for evidence in active_evidence:
                    integrity_result = await self.verify_evidence_integrity(evidence.evidence_id)
                    
                    if not integrity_result['overall_integrity']:
                        await self._handle_integrity_breach(evidence, integrity_result)
                
                # Check every hour
                await asyncio.sleep(3600)
                
            except Exception as e:
                logger.error(f"Integrity monitoring error: {e}")
                await asyncio.sleep(300)  # Retry after 5 minutes
