# core/security/session_manager.py
class SecureSessionManager:
    """Enterprise session management with advanced security controls"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.active_sessions = {}
        self.session_store = EncryptedSessionStore()
        self.threat_detector = SessionThreatDetector()
        
    async def create_secure_session(self, auth_result: Dict) -> Session:
        """Create a secure session with comprehensive controls"""
        session_id = self._generate_secure_session_id()
        
        session = Session(
            session_id=session_id,
            user_id=auth_result['user_id'],
            created_at=datetime.now(),
            expires_at=datetime.now() + timedelta(hours=8),
            roles=auth_result['assigned_roles'],
            context={
                'ip_address': auth_result.get('ip_address'),
                'user_agent': auth_result.get('user_agent'),
                'device_fingerprint': auth_result.get('device_fingerprint')
            },
            security_attributes={
                'mfa_strength': auth_result.get('authentication_strength'),
                'risk_level': auth_result.get('risk_level'),
                'required_reauthentication': self._calculate_reauth_frequency(auth_result)
            }
        )
        
        # Store encrypted session
        await self.session_store.store_session(session)
        
        # Start session monitoring
        asyncio.create_task(self._monitor_session_security(session))
        
        return session
    
    async def validate_session(self, session_token: str) -> Optional[Session]:
        """Validate session with comprehensive security checks"""
        session = await self.session_store.retrieve_session(session_token)
        
        if not session:
            return None
        
        # Check session expiration
        if datetime.now() > session.expires_at:
            await self._terminate_session(session_token)
            return None
        
        # Check for suspicious activity
        threat_detected = await self.threat_detector.analyze_session_activity(session)
        if threat_detected:
            await self._terminate_session(session_token)
            await self._alert_security_team(session, threat_detected)
            return None
        
        # Check for required reauthentication
        if await self._requires_reauthentication(session):
            return None
        
        # Update last activity
        session.last_activity = datetime.now()
        await self.session_store.update_session(session)
        
        return session
    
    async def _monitor_session_security(self, session: Session):
        """Continuous session security monitoring"""
        while datetime.now() < session.expires_at:
            try:
                # Check for geographic anomalies
                if await self._detect_geographic_anomaly(session):
                    await self._force_reauthentication(session.session_id)
                    break
                
                # Check for behavior anomalies
                if await self._detect_behavioral_anomaly(session):
                    await self._force_reauthentication(session.session_id)
                    break
                
                # Check for device changes
                if await self._detect_device_change(session):
                    await self._terminate_session(session.session_id)
                    break
                
                await asyncio.sleep(60)  # Check every minute
                
            except Exception as e:
                logger.error(f"Session monitoring error: {e}")
                await asyncio.sleep(30)
