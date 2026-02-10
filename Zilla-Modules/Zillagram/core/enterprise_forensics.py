# core/enterprise_forensics.py
class EnterpriseForensicManager:
    """
    Main forensic evidence management interface
    Integrates with OSINT collection and analysis
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.forensic_engine = EnterpriseForensicEngine(config)
        self.compliance_engine = EnterpriseComplianceEngine(config)
        self.access_control = EnterpriseAccessControl(config)
        
    async def collect_osint_evidence(self, 
                                   intelligence_data: Dict,
                                   case_id: str,
                                   collector: str) -> ForensicEvidence:
        """
        Collect OSINT intelligence as forensic evidence
        """
        # Verify legal basis for collection
        legal_basis = await self.compliance_engine.verify_legal_basis(
            "intelligence_collection", intelligence_data, {"collector": collector}
        )
        
        if not legal_basis:
            raise LegalComplianceError("No legal basis for intelligence collection")
        
        # Convert intelligence to forensic evidence
        evidence_data = await self._prepare_intelligence_evidence(intelligence_data)
        
        # Collect as forensic evidence
        evidence = await self.forensic_engine.collect_evidence(
            evidence_type=EvidenceType.DIGITAL_MESSAGE,
            raw_data=evidence_data['raw_data'],
            metadata=evidence_data['metadata'],
            case_id=case_id,
            collector=collector
        )
        
        logger.info(f"🔍 Collected OSINT evidence {evidence.evidence_id} for case {case_id}")
        return evidence
    
    async def prepare_legal_case_package(self, case_id: str, jurisdiction: str) -> Dict:
        """
        Prepare complete legal case package for court proceedings
        """
        # Get all evidence for case
        case_evidence = await self.forensic_engine.get_case_evidence(case_id)
        
        case_package = {
            'case_id': case_id,
            'preparation_date': datetime.utcnow().isoformat(),
            'jurisdiction': jurisdiction,
            'case_summary': await self._generate_case_summary(case_id),
            'evidence_inventory': [],
            'chain_of_custody_reports': [],
            'legal_affidavits': [],
            'expert_analysis_reports': [],
            'discovery_documents': [],
            'court_submission_guide': await self._generate_court_submission_guide(jurisdiction)
        }
        
        # Process each piece of evidence
        for evidence in case_evidence:
            # Prepare court submission for each evidence
            evidence_submission = await self.forensic_engine.legal_compliance.prepare_court_submission(
                evidence.evidence_id, jurisdiction
            )
            
            case_package['evidence_inventory'].append({
                'evidence_id': evidence.evidence_id,
                'evidence_type': evidence.evidence_type.value,
                'collection_date': evidence.collection_timestamp.isoformat(),
                'submission_package': evidence_submission
            })
        
        # Generate comprehensive chain of custody report
        case_package['master_chain_of_custody'] = await self._generate_master_custody_report(case_id)
        
        return case_package

# Usage in main OSINT system
async def demonstrate_forensic_workflow():
    """Demonstrate complete forensic evidence workflow"""
    forensic_mgr = EnterpriseForensicManager(load_config())
    
    # Collect intelligence as evidence
    intelligence_data = {
        'message_content': 'Suspicious communication detected',
        'sender_info': {'user_id': 'user123', 'username': 'suspicious_user'},
        'timestamp': '2024-01-15T10:30:00Z',
        'metadata': {'source': 'telegram_channel', 'collection_method': 'automated_monitoring'}
    }
    
    evidence = await forensic_mgr.collect_osint_evidence(
        intelligence_data, 
        case_id="case_2024_001",
        collector="ai_analyst_001"
    )
    
    # Perform forensic analysis
    analysis_results = await forensic_mgr.forensic_engine.analysis_integration.analyze_evidence_forensically(
        evidence.evidence_id, "message_authentication", "forensic_analyst_001"
    )
    
    # Prepare for court
    court_package = await forensic_mgr.prepare_legal_case_package("case_2024_001", "US_Federal")
    
    logger.info(f"⚖️ Legal case package prepared with {len(court_package['evidence_inventory'])} evidence items")
    return court_package
