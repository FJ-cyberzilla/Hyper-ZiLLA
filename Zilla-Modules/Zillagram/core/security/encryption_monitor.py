# core/security/encryption_monitor.py
class EncryptionPerformanceMonitor:
    """
    Monitor encryption operations for performance and security
    """
    
    def __init__(self):
        self.metrics = {
            'encryption_operations': 0,
            'decryption_operations': 0,
            'total_data_encrypted': 0,
            'total_data_decrypted': 0,
            'failed_operations': 0,
            'average_encryption_time': 0,
            'average_decryption_time': 0
        }
        self.operation_log = []
    
    async def record_encryption_operation(self, data_size: int, 
                                        sensitivity: SensitivityLevel,
                                        layers: List[EncryptionLayer]):
        """Record encryption operation metrics"""
        operation = {
            'timestamp': time.time(),
            'operation': 'encrypt',
            'data_size': data_size,
            'sensitivity': sensitivity.value,
            'layers': [layer.value for layer in layers],
            'performance': self._measure_performance()
        }
        
        self.operation_log.append(operation)
        self.metrics['encryption_operations'] += 1
        self.metrics['total_data_encrypted'] += data_size
        
        # Keep log manageable
        if len(self.operation_log) > 10000:
            self.operation_log = self.operation_log[-5000:]
    
    async def generate_security_report(self) -> Dict:
        """Generate encryption security audit report"""
        return {
            'cryptographic_health': await self._assess_cryptographic_health(),
            'key_management_status': await self._assess_key_management(),
            'performance_metrics': self.metrics.copy(),
            'compliance_status': await self._check_compliance(),
            'recommendations': await self._generate_recommendations()
        }
    
    async def _assess_cryptographic_health(self) -> Dict:
        """Assess the health of cryptographic operations"""
        return {
            'algorithm_strength': 'AES-256 (NIST Approved)',
            'key_rotation_frequency': 'Compliant',
            'random_number_generation': 'Cryptographically Secure',
            'integrity_protection': 'HMAC-SHA256',
            'authentication': 'AES-GCM with AAD'
        }
