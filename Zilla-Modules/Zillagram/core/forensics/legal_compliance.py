# core/forensics/legal_compliance.py
class ForensicLegalCompliance:
    """
    Ensures forensic evidence meets legal standards for court admissibility
    Implements rules of evidence and legal requirements
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.legal_requirements = self._load_legal_requirements()
        self.jurisdiction_rules = self._load_jurisdiction_rules()
        
    async def prepare_court_submission(self, evidence_id: str, jurisdiction: str) -> Dict:
        """
        Prepare evidence package for court submission
        """
        evidence = await self.evidence_store.retrieve_evidence(evidence_id)
        if not evidence:
            raise EvidenceNotFoundError(f"Evidence {evidence_id} not found")
        
        # Verify legal compliance for specific jurisdiction
        compliance_check = await self._verify_jurisdiction_compliance(evidence, jurisdiction)
        
        if not compliance_check['compliant']:
            raise LegalComplianceError(f"Evidence not compliant with {jurisdiction} rules")
        
        # Prepare court submission package
        submission_package = {
            'submission_id': self._generate_submission_id(),
            'evidence_id': evidence_id,
            'jurisdiction': jurisdiction,
            'prepared_by': self.config['forensics']['legal_authority'],
            'preparation_date': datetime.utcnow().isoformat(),
            'evidence_summary': await self._generate_evidence_summary(evidence),
            'chain_of_custody_report': await self.chain_of_custody.generate_chain_of_custody_report(evidence_id),
            'integrity_verification_report': await self.integrity_verifier.verify_evidence_integrity(evidence_id),
            'legal_affidavits': await self._prepare_legal_affidavits(evidence),
            'expert_testimony_preparation': await self._prepare_expert_testimony(evidence),
            'discovery_documents': await self._prepare_discovery_documents(evidence),
            'court_instructions': await self._generate_court_instructions(evidence, jurisdiction)
        }
        
        # Apply legal formatting and requirements
        formatted_package = await self._apply_legal_formatting(submission_package, jurisdiction)
        
        logger.info(f"⚖️ Court submission package prepared for evidence {evidence_id}")
        return formatted_package
    
    async def _verify_jurisdiction_compliance(self, evidence: ForensicEvidence, jurisdiction: str) -> Dict:
        """Verify evidence meets specific jurisdiction's legal requirements"""
        jurisdiction_rules = self.jurisdiction_rules.get(jurisdiction, {})
        
        compliance_checks = {
            'authentication_requirement': await self._verify_authentication_requirements(evidence, jurisdiction_rules),
            'hearsay_exception': await self._verify_hearsay_exception(evidence, jurisdiction_rules),
            'best_evidence_rule': await self._verify_best_evidence_rule(evidence),
            'relevance_requirement': await self._verify_relevance(evidence),
            'prejudice_assessment': await self._assess_prejudice_impact(evidence),
            'expert_testimony_requirement': await self._check_expert_testimony_need(evidence)
        }
        
        return {
            'jurisdiction': jurisdiction,
            'compliant': all(compliance_checks.values()),
            'compliance_checks': compliance_checks,
            'legal_citations': await self._provide_legal_citations(compliance_checks, jurisdiction)
        }
    
    async def _prepare_legal_affidavits(self, evidence: ForensicEvidence) -> List[Dict]:
        """Prepare legal affidavits for evidence authentication"""
        affidavits = []
        
        # Collector affidavit
        collector_affidavit = {
            'affidavit_type': 'collection_authentication',
            'affiant': evidence.collected_by,
            'affidavit_date': datetime.utcnow().isoformat(),
            'sworn_statement': f"I, {evidence.collected_by}, hereby swear that I collected the evidence with ID {evidence.evidence_id} on {evidence.collection_timestamp.isoformat()} using forensically sound methods.",
            'exhibits_attached': [evidence.evidence_id],
            'notarization_required': True,
            'jurisdiction': 'To be determined'
        }
        affidavits.append(collector_affidavit)
        
        # Custody chain affidavits
        custody_chain = await self.chain_of_custody.get_evidence_custody_chain(evidence.evidence_id)
        for custody_event in custody_chain:
            if custody_event['action'] in ['transfer', 'receipt']:
                custody_affidavit = {
                    'affidavit_type': 'custody_verification',
                    'affiant': custody_event['actor'],
                    'affidavit_date': datetime.utcnow().isoformat(),
                    'sworn_statement': self._generate_custody_affidavit_statement(custody_event),
                    'exhibits_attached': [evidence.evidence_id],
                    'notarization_required': True
                }
                affidavits.append(custody_affidavit)
        
        return affidavits
