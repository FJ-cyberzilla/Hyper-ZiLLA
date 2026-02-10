# core/enterprise_security.py
class EnterpriseSecurityManager:
    """
    Unified enterprise security management
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.compliance_engine = EnterpriseComplianceEngine(config)
        self.access_control = EnterpriseAccessControl(config)
        self.api_gateway = EnterpriseAPISecurityGateway(config)
        self.encryption_engine = EnterpriseEncryptionEngine(config)
        
        logger.info("🛡️ Enterprise Security System Initialized")
    
    async def secure_data_operation(self, operation: str, data: Dict, user_context: Dict) -> Dict:
        """End-to-end secure data operation"""
        try:
            # 1. Compliance check
            compliance_result = await self.compliance_engine.process_data_operation(
                operation, data, user_context
            )
            
            # 2. Access control check
            access_granted = await self.access_control.authorize_access(
                user_context['session_token'], 
                data['resource'],
                operation
            )
            
            if not access_granted:
                raise AccessDeniedError("Insufficient permissions")
            
            # 3. Encrypt sensitive data
            encrypted_data = await self.encryption_engine.encrypt_data(
                data, 
                self._determine_sensitivity(data)
            )
            
            return {
                "success": True,
                "compliance": compliance_result,
                "encryption": encrypted_data.encryption_metadata,
                "audit_trail": await self._generate_audit_trail(operation, data, user_context)
            }
            
        except Exception as e:
            await self._handle_security_incident(operation, data, user_context, str(e))
            raise

# FastAPI Integration Example
from fastapi import FastAPI, Depends, HTTPException
from fastapi.security import HTTPBearer

app = FastAPI()
security = HTTPBearer()
enterprise_security = EnterpriseSecurityManager(load_config())

@app.middleware("http")
async def security_middleware(request: Request, call_next):
    """Enterprise security middleware for all requests"""
    security_result = await enterprise_security.api_gateway.process_api_request(request)
    
    if not security_result.get('allowed', True):
        raise HTTPException(
            status_code=429 if security_result.get('reason') == 'rate_limit_exceeded' else 403,
            detail=security_result.get('reason', 'access_denied')
        )
    
    response = await call_next(request)
    return response

@app.post("/api/v1/analyze")
async def analyze_message(
    message_data: Dict,
    credentials: HTTPAuthorizationCredentials = Depends(security)
):
    """Secure API endpoint example"""
    user_context = await enterprise_security.access_control.authenticate_user({
        'token': credentials.credentials
    })
    
    result = await enterprise_security.secure_data_operation(
        "analyze", 
        message_data, 
        user_context
    )
    
    return result
