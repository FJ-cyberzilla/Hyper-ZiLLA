# core/forensics/evidence_engine.py
import hashlib
import json
import time
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, asdict
from enum import Enum
from datetime import datetime, timedelta
import asyncio
import base64
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives.serialization import load_pem_private_key
import sqlite3
import os

class EvidenceType(Enum):
    DIGITAL_MESSAGE = "digital_message"
    USER_BEHAVIOR = "user_behavior"
    NETWORK_ACTIVITY = "network_activity"
    SYSTEM_LOGS = "system_logs"
    MEDIA_CONTENT = "media_content"
    METADATA = "metadata"
    COMMUNICATION_PATTERN = "communication_pattern"

class EvidenceStatus(Enum):
    COLLECTED = "collected"
    VERIFIED = "verified"
    PRESERVED = "preserved"
    ANALYZED = "analyzed"
    CHAIN_VERIFIED = "chain_verified"
    COURT_READY = "court_ready"
    ARCHIVED = "archived"

@dataclass
class ForensicEvidence:
    evidence_id: str
    case_id: str
    evidence_type: EvidenceType
    source: str
    collected_by: str
    collection_timestamp: datetime
    content_hash: str
    metadata: Dict
    raw_data: bytes
    status: EvidenceStatus
    chain_of_custody: List[Dict]
    legal_notes: str = ""
    sensitivity_level: str = "high"

class EnterpriseForensicEngine:
    """
    Legal-grade forensic evidence management system
    Maintains chain of custody and evidence integrity for court admissibility
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.evidence_store = SecureEvidenceStore(config)
        self.chain_of_custody = ChainOfCustodyManager(config)
        self.integrity_verifier = EvidenceIntegrityVerifier(config)
        self.legal_compliance = ForensicLegalCompliance(config)
        
        # Digital signatures for legal authenticity
        self.signing_key = self._load_signing_certificate()
        
    async def collect_evidence(self, 
                             evidence_type: EvidenceType,
                             raw_data: bytes,
                             metadata: Dict,
                             case_id: str,
                             collector: str) -> ForensicEvidence:
        """
        Collect forensic evidence with legal-grade preservation
        """
        # Generate unique evidence ID
        evidence_id = self._generate_evidence_id(evidence_type, case_id)
        
        # Create content hash for integrity
        content_hash = self._calculate_forensic_hash(raw_data)
        
        # Create evidence object
        evidence = ForensicEvidence(
            evidence_id=evidence_id,
            case_id=case_id,
            evidence_type=evidence_type,
            source=metadata.get('source', 'unknown'),
            collected_by=collector,
            collection_timestamp=datetime.utcnow(),
            content_hash=content_hash,
            metadata=metadata,
            raw_data=raw_data,
            status=EvidenceStatus.COLLECTED,
            chain_of_custody=[{
                'action': 'collection',
                'actor': collector,
                'timestamp': datetime.utcnow(),
                'location': self._get_collection_location(),
                'verification_hash': content_hash
            }]
        )
        
        # Apply legal preservation measures
        preserved_evidence = await self._apply_legal_preservation(evidence)
        
        # Store in secure evidence repository
        await self.evidence_store.store_evidence(preserved_evidence)
        
        logger.info(f"🔍 Collected forensic evidence: {evidence_id} for case: {case_id}")
        return preserved_evidence
    
    async def _apply_legal_preservation(self, evidence: ForensicEvidence) -> ForensicEvidence:
        """
        Apply legal-grade evidence preservation techniques
        """
        preservation_steps = [
            self._create_digital_signature(evidence),
            self._generate_integrity_checksums(evidence),
            self._apply_tamper_evidence(evidence),
            self._create_legal_metadata(evidence),
            self._encrypt_sensitive_components(evidence)
        ]
        
        preserved_evidence = evidence
        for step in preservation_steps:
            preserved_evidence = await step(preserved_evidence)
        
        preserved_evidence.status = EvidenceStatus.PRESERVED
        return preserved_evidence
    
    async def _create_digital_signature(self, evidence: ForensicEvidence) -> ForensicEvidence:
        """Create digital signature for legal authenticity"""
        signature_data = {
            'evidence_id': evidence.evidence_id,
            'content_hash': evidence.content_hash,
            'timestamp': evidence.collection_timestamp.isoformat(),
            'collector': evidence.collected_by
        }
        
        signature = self._sign_data(json.dumps(signature_data, sort_keys=True))
        
        # Add signature to metadata
        evidence.metadata['digital_signature'] = {
            'signature': base64.b64encode(signature).decode(),
            'algorithm': 'RSA-SHA256',
            'signing_timestamp': datetime.utcnow().isoformat(),
            'signing_authority': self.config['forensics']['signing_authority']
        }
        
        return evidence
    
    def _calculate_forensic_hash(self, data: bytes) -> str:
        """Calculate multiple forensic hashes for evidence integrity"""
        hash_algorithms = {
            'md5': hashlib.md5(data).hexdigest(),
            'sha1': hashlib.sha1(data).hexdigest(),
            'sha256': hashlib.sha256(data).hexdigest(),
            'sha512': hashlib.sha512(data).hexdigest()
        }
        
        return json.dumps(hash_algorithms, sort_keys=True)
