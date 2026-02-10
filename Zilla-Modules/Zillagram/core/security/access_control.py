# core/security/access_control.py
from enum import Enum
from typing import Dict, List, Set
import asyncio
from dataclasses import dataclass
from datetime import datetime, timedelta

class PermissionLevel(Enum):
    VIEW = "view"
    EDIT = "edit" 
    DELETE = "delete"
    ADMINISTER = "administer"
    SUPER_USER = "super_user"

@dataclass
class Role:
    role_id: str
    name: str
    permissions: Set[PermissionLevel]
    data_scopes: List[str]
    temporal_constraints: Dict
    inheritance: List[str] = None

class EnterpriseAccessControl:
    """
    Combined RBAC (Role-Based) and ABAC (Attribute-Based) Access Control
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.role_registry = self._initialize_roles()
        self.policy_engine = PolicyDecisionPoint()
        self.session_manager = SecureSessionManager()
        self.mfa_engine = MultiFactorAuthEngine()
        
    def _initialize_roles(self) -> Dict[str, Role]:
        """Initialize enterprise roles with least privilege"""
        return {
            "analyst": Role(
                role_id="analyst",
                name="Intelligence Analyst",
                permissions={PermissionLevel.VIEW, PermissionLevel.EDIT},
                data_scopes=["telegram_data", "analysis_results"],
                temporal_constraints={"max_session_hours": 8}
            ),
            "supervisor": Role(
                role_id="supervisor",
                name="Team Supervisor", 
                permissions={PermissionLevel.VIEW, PermissionLevel.EDIT, PermissionLevel.DELETE},
                data_scopes=["telegram_data", "analysis_results", "user_management"],
                temporal_constraints={"max_session_hours": 12},
                inheritance=["analyst"]
            ),
            "admin": Role(
                role_id="admin",
                name="System Administrator",
                permissions={PermissionLevel.VIEW, PermissionLevel.EDIT, PermissionLevel.DELETE, PermissionLevel.ADMINISTER},
                data_scopes=["all_data"],
                temporal_constraints={"max_session_hours": 24},
                inheritance=["supervisor"]
            ),
            "auditor": Role(
                role_id="auditor",
                name="Compliance Auditor",
                permissions={PermissionLevel.VIEW},
                data_scopes=["audit_logs", "compliance_data"],
                temporal_constraints={"max_session_hours": 4}
            )
        }
    
    async def authenticate_user(self, credentials: Dict) -> Dict:
        """Multi-factor authentication with risk assessment"""
        auth_steps = [
            self._verify_primary_credentials(credentials),
            self._check_mfa_requirement(credentials),
            self._assess_authentication_risk(credentials),
            self._verify_device_compliance(credentials)
        ]
        
        results = await asyncio.gather(*auth_steps)
        
        if all(results):
            session = await self.session_manager.create_secure_session(credentials)
            return {
                "authenticated": True,
                "session_token": session.token,
                "expires_at": session.expires_at,
                "assigned_roles": await self._determine_user_roles(credentials)
            }
        else:
            await self._handle_failed_authentication(credentials, results)
            return {"authenticated": False, "reason": "authentication_failed"}
    
    async def authorize_access(self, session_token: str, resource: str, action: PermissionLevel) -> bool:
        """Authorize access using RBAC + ABAC"""
        # Validate session
        session = await self.session_manager.validate_session(session_token)
        if not session:
            return False
        
        # Get user roles
        user_roles = await self._get_user_roles(session.user_id)
        
        # Check RBAC permissions
        rbac_allowed = await self._check_rbac_permissions(user_roles, resource, action)
        
        # Check ABAC policies
        abac_allowed = await self.policy_engine.evaluate_access_policies(
            session.user_id, resource, action, session.context
        )
        
        # Check temporal constraints
        temporal_allowed = await self._check_temporal_constraints(session, user_roles)
        
        return rbac_allowed and abac_allowed and temporal_allowed
    
    async def _check_rbac_permissions(self, roles: List[Role], resource: str, action: PermissionLevel) -> bool:
        """Check Role-Based Access Control permissions"""
        for role in roles:
            # Check direct permissions
            if action in role.permissions:
                # Check data scope
                if self._is_resource_in_scope(resource, role.data_scopes):
                    return True
            
            # Check inherited roles
            if role.inheritance:
                inherited_roles = [self.role_registry[inherited] for inherited in role.inheritance]
                if await self._check_rbac_permissions(inherited_roles, resource, action):
                    return True
        
        return False
