# core/forensics/chain_of_custody.py
class ChainOfCustodyManager:
    """
    Maintains legal chain of custody for forensic evidence
    Tracks every access, transfer, and analysis of evidence
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.custody_log = CustodyLogDatabase(config)
        self.access_controls = EvidenceAccessControls(config)
        
    async def record_custody_transfer(self, 
                                    evidence_id: str,
                                    from_actor: str,
                                    to_actor: str,
                                    transfer_reason: str,
                                    transfer_method: str) -> Dict:
        """
        Record evidence transfer with full legal documentation
        """
        transfer_record = {
            'transfer_id': self._generate_transfer_id(),
            'evidence_id': evidence_id,
            'from_actor': from_actor,
            'to_actor': to_actor,
            'transfer_timestamp': datetime.utcnow(),
            'transfer_reason': transfer_reason,
            'transfer_method': transfer_method,
            'pre_transfer_verification': await self._verify_evidence_integrity(evidence_id),
            'post_transfer_verification': None,  # Will be set by receiver
            'witnesses': await self._get_transfer_witnesses(),
            'location': self._get_transfer_location(),
            'legal_authorization': await self._get_legal_authorization(transfer_reason)
        }
        
        # Store transfer record
        await self.custody_log.record_transfer(transfer_record)
        
        # Update evidence chain of custody
        await self._update_evidence_chain(evidence_id, transfer_record)
        
        logger.info(f"📦 Evidence {evidence_id} transferred from {from_actor} to {to_actor}")
        return transfer_record
    
    async def record_evidence_access(self, 
                                   evidence_id: str,
                                   accessing_actor: str,
                                   access_purpose: str,
                                   access_duration: timedelta) -> Dict:
        """
        Record every access to evidence for complete audit trail
        """
        access_record = {
            'access_id': self._generate_access_id(),
            'evidence_id': evidence_id,
            'accessing_actor': accessing_actor,
            'access_timestamp': datetime.utcnow(),
            'access_purpose': access_purpose,
            'access_duration': access_duration,
            'pre_access_integrity': await self._verify_evidence_integrity(evidence_id),
            'post_access_integrity': None,
            'tools_used': await self._record_analysis_tools(access_purpose),
            'legal_justification': await self._get_legal_justification(access_purpose)
        }
        
        # Verify access authorization
        if not await self.access_controls.verify_access_rights(accessing_actor, evidence_id):
            raise UnauthorizedAccessError(f"Actor {accessing_actor} not authorized for evidence {evidence_id}")
        
        await self.custody_log.record_access(access_record)
        
        return access_record
    
    async def generate_chain_of_custody_report(self, evidence_id: str) -> Dict:
        """
        Generate comprehensive chain of custody report for legal proceedings
        """
        custody_chain = await self.custody_log.get_evidence_custody_chain(evidence_id)
        
        report = {
            'evidence_id': evidence_id,
            'report_generated': datetime.utcnow().isoformat(),
            'generated_by': self.config['forensics']['reporting_authority'],
            'custody_chain': custody_chain,
            'integrity_verification': await self._verify_complete_chain_integrity(custody_chain),
            'gaps_analysis': await self._analyze_custody_gaps(custody_chain),
            'legal_compliance': await self._assess_legal_compliance(custody_chain),
            'witness_verification': await self._verify_witness_authenticity(custody_chain)
        }
        
        # Digital signature for report authenticity
        report['digital_signature'] = await self._sign_legal_report(report)
        
        return report
    
    async def _verify_complete_chain_integrity(self, custody_chain: List[Dict]) -> Dict:
        """Verify integrity throughout the entire chain of custody"""
        integrity_checks = {}
        
        for i, custody_event in enumerate(custody_chain):
            if i > 0:
                # Verify continuity between custody events
                prev_event = custody_chain[i-1]
                current_event = custody_event
                
                continuity_check = await self._verify_custody_continuity(prev_event, current_event)
                integrity_checks[f"continuity_{i-1}_to_{i}"] = continuity_check
            
            # Verify individual event integrity
            event_integrity = await self._verify_custody_event_integrity(custody_event)
            integrity_checks[f"event_{i}_integrity"] = event_integrity
        
        return {
            'all_checks_passed': all(integrity_checks.values()),
            'detailed_checks': integrity_checks,
            'overall_confidence': sum(integrity_checks.values()) / len(integrity_checks) if integrity_checks else 1.0
        }
