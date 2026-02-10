# core/configuration/drift_detector.py
class ConfigDriftDetector:
    """Advanced configuration drift detection and analysis"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.core_v1 = client.CoreV1Api()
        
    async def detect_drift(self) -> Dict:
        """Detect configuration drift across all systems"""
        drift_analysis = {
            "kubernetes_drift": await self._detect_k8s_drift(),
            "application_drift": await self._detect_application_drift(),
            "infrastructure_drift": await self._detect_infrastructure_drift(),
            "security_drift": await self._detect_security_drift()
        }
        
        return {
            "drift_detected": any(drift_analysis.values()),
            "timestamp": datetime.utcnow(),
            "detailed_analysis": drift_analysis,
            "severity": self._calculate_drift_severity(drift_analysis)
        }
    
    async def _detect_k8s_drift(self) -> Dict:
        """Detect Kubernetes configuration drift"""
        k8s_drift = {}
        
        # Check ConfigMap drift
        configmaps = await self._get_managed_configmaps()
        for cm in configmaps:
            current_config = await self._get_current_configmap(cm)
            expected_config = await self._get_expected_configmap(cm)
            
            if current_config != expected_config:
                k8s_drift[cm] = {
                    "type": "configmap_drift",
                    "differences": self._compare_configurations(current_config, expected_config),
                    "severity": "medium"
                }
        
        # Check Secret drift
        secrets = await self._get_managed_secrets()
        for secret in secrets:
            # Compare secret checksums (not actual values for security)
            current_checksum = await self._get_secret_checksum(secret)
            expected_checksum = await self._get_expected_secret_checksum(secret)
            
            if current_checksum != expected_checksum:
                k8s_drift[secret] = {
                    "type": "secret_drift", 
                    "severity": "high"
                }
        
        return k8s_drift
    
    async def _detect_application_drift(self) -> Dict:
        """Detect application configuration drift"""
        app_drift = {}
        
        # Check environment variables
        deployments = await self._get_managed_deployments()
        for deployment in deployments:
            env_drift = await self._check_environment_drift(deployment)
            if env_drift:
                app_drift[deployment] = env_drift
        
        # Check resource limits
        resource_drift = await self._check_resource_limits_drift()
        if resource_drift:
            app_drift['resource_limits'] = resource_drift
        
        return app_drift
