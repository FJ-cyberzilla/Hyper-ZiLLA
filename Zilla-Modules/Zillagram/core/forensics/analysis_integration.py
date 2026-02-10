# core/forensics/analysis_integration.py
class ForensicAnalysisIntegration:
    """
    Integrates forensic evidence with analysis systems
    Maintains legal integrity during analysis
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.forensic_engine = EnterpriseForensicEngine(config)
        self.analysis_tools = ForensicAnalysisTools(config)
        
    async def analyze_evidence_forensically(self, evidence_id: str, analysis_type: str, analyst: str) -> Dict:
        """
        Perform forensic analysis while maintaining legal integrity
        """
        # Record analysis start in chain of custody
        access_record = await self.chain_of_custody.record_evidence_access(
            evidence_id, analyst, f"forensic_analysis_{analysis_type}", timedelta(hours=2)
        )
        
        try:
            # Retrieve evidence for analysis
            evidence = await self.evidence_store.retrieve_evidence(evidence_id, analyst, "forensic_analysis")
            
            # Perform analysis based on type
            analysis_results = await self._perform_forensic_analysis(evidence, analysis_type, analyst)
            
            # Verify post-analysis integrity
            post_analysis_integrity = await self.integrity_verifier.verify_evidence_integrity(evidence_id)
            
            # Update access record with completion
            access_record['post_access_integrity'] = post_analysis_integrity
            access_record['analysis_results_reference'] = analysis_results['analysis_id']
            await self.custody_log.update_access_record(access_record)
            
            return {
                'analysis_id': analysis_results['analysis_id'],
                'evidence_id': evidence_id,
                'analysis_type': analysis_type,
                'analyst': analyst,
                'results': analysis_results,
                'integrity_maintained': post_analysis_integrity['overall_integrity'],
                'legal_admissibility': analysis_results.get('legal_admissibility', True)
            }
            
        except Exception as e:
            # Record analysis failure
            access_record['analysis_failed'] = True
            access_record['failure_reason'] = str(e)
            await self.custody_log.update_access_record(access_record)
            
            raise ForensicAnalysisError(f"Forensic analysis failed: {e}") from e
    
    async def _perform_forensic_analysis(self, evidence: ForensicEvidence, analysis_type: str, analyst: str) -> Dict:
        """Perform specific forensic analysis types"""
        analysis_methods = {
            'message_authentication': self._analyze_message_authenticity,
            'behavioral_patterns': self._analyze_behavioral_patterns,
            'metadata_analysis': self._analyze_metadata_forensically,
            'timeline_reconstruction': self._reconstruct_timeline,
            'network_forensics': self._analyze_network_forensics,
            'digital_fingerprinting': self._perform_digital_fingerprinting
        }
        
        analysis_method = analysis_methods.get(analysis_type)
        if not analysis_method:
            raise ValueError(f"Unknown analysis type: {analysis_type}")
        
        return await analysis_method(evidence, analyst)
    
    async def _analyze_message_authenticity(self, evidence: ForensicEvidence, analyst: str) -> Dict:
        """Forensic analysis of message authenticity"""
        analysis_results = {
            'analysis_id': f"auth_analysis_{evidence.evidence_id}",
            'analysis_type': 'message_authentication',
            'analyst': analyst,
            'analysis_timestamp': datetime.utcnow().isoformat(),
            'authenticity_indicators': {
                'digital_signature_valid': await self._verify_message_signature(evidence),
                'timestamp_consistency': await self._analyze_timestamp_consistency(evidence),
                'content_integrity': await self._verify_content_integrity(evidence),
                'metadata_authenticity': await self._analyze_metadata_authenticity(evidence),
                'behavioral_consistency': await self._analyze_behavioral_consistency(evidence)
            },
            'authenticity_score': 0.0,  # Calculated based on indicators
            'expert_opinion': await self._generate_expert_opinion(evidence),
            'legal_interpretation': await self._provide_legal_interpretation(evidence)
        }
        
        # Calculate overall authenticity score
        indicators = analysis_results['authenticity_indicators']
        analysis_results['authenticity_score'] = sum(indicators.values()) / len(indicators)
        
        return analysis_results
