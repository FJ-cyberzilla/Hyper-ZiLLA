# core/security/api_gateway.py
from typing import Dict, List, Optional
import asyncio
import time
from dataclasses import dataclass
from fastapi import Request, Response
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

@dataclass
class APISecurityContext:
    client_id: str
    user_id: str
    permissions: List[str]
    rate_limit_tier: str
    threat_level: float
    authentication_strength: float

class EnterpriseAPISecurityGateway:
    """
    Comprehensive API security gateway with threat detection and rate limiting
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.rate_limiter = AdaptiveRateLimiter(config)
        self.threat_detector = APIThreatDetector(config)
        self.request_validator = APIRequestValidator(config)
        self.auth_provider = APIAuthProvider(config)
        
        # Security middleware chain
        self.middleware_chain = [
            self._validate_authentication,
            self._check_rate_limits,
            self._detect_threats,
            self._validate_request,
            self._enforce_policies
        ]
    
    async def process_api_request(self, request: Request) -> Dict:
        """Process API request through security gateway"""
        security_context = await self._create_security_context(request)
        
        # Execute security middleware chain
        for middleware in self.middleware_chain:
            result = await middleware(request, security_context)
            if not result.get('allowed', True):
                await self._handle_blocked_request(request, security_context, result)
                return result
        
        # Request passed all security checks
        await self._log_approved_request(request, security_context)
        
        return {
            "allowed": True,
            "security_context": security_context,
            "processing_tier": await self._determine_processing_tier(security_context)
        }
    
    async def _validate_authentication(self, request: Request, context: APISecurityContext) -> Dict:
        """Validate API authentication and authorization"""
        try:
            # Extract credentials
            credentials = await self._extract_credentials(request)
            
            # Authenticate
            auth_result = await self.auth_provider.authenticate(credentials)
            if not auth_result.authenticated:
                return {"allowed": False, "reason": "authentication_failed"}
            
            # Update security context
            context.user_id = auth_result.user_id
            context.permissions = auth_result.permissions
            context.authentication_strength = auth_result.strength
            
            # Check authorization for requested endpoint
            endpoint = request.url.path
            method = request.method
            
            if not await self.auth_provider.authorize(auth_result.user_id, endpoint, method):
                return {"allowed": False, "reason": "authorization_failed"}
            
            return {"allowed": True}
            
        except Exception as e:
            logger.error(f"Authentication validation failed: {e}")
            return {"allowed": False, "reason": "authentication_error"}
    
    async def _check_rate_limits(self, request: Request, context: APISecurityContext) -> Dict:
        """Apply adaptive rate limiting"""
        client_identifier = await self._get_client_identifier(request)
        endpoint = request.url.path
        
        rate_limit_result = await self.rate_limiter.check_rate_limit(
            client_identifier, endpoint, context
        )
        
        if not rate_limit_result.allowed:
            return {
                "allowed": False,
                "reason": "rate_limit_exceeded",
                "retry_after": rate_limit_result.retry_after,
                "limit_details": rate_limit_result.limits
            }
        
        return {"allowed": True}
