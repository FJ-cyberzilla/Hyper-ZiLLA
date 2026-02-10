# core/security/encryption_engine.py
import os
import base64
import json
import hashlib
import hmac
from typing import Dict, Any, Optional, Tuple
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import padding
from cryptography.hazmat.backends import default_backend
import asyncio
from dataclasses import dataclass
from enum import Enum

class EncryptionLayer(Enum):
    IN_TRANSIT = "in_transit"
    AT_REST = "at_rest" 
    IN_MEMORY = "in_memory"
    BACKUP = "backup"

class SensitivityLevel(Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"

@dataclass
class EncryptionResult:
    encrypted_data: bytes
    iv: bytes
    tag: Optional[bytes] = None
    key_id: str = None
    encryption_metadata: Dict = None

class EnterpriseEncryptionEngine:
    """
    Enterprise-grade multi-layer encryption system
    Implements NIST-approved cryptographic standards
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.backend = default_backend()
        self.key_manager = KeyManagementSystem(config)
        self.performance_monitor = EncryptionPerformanceMonitor()
        
        # Initialize encryption schemes
        self.encryption_schemes = {
            EncryptionLayer.IN_TRANSIT: self._init_transit_encryption(),
            EncryptionLayer.AT_REST: self._init_at_rest_encryption(),
            EncryptionLayer.IN_MEMORY: self._init_memory_encryption(),
            EncryptionLayer.BACKUP: self._init_backup_encryption()
        }
        
        logger.info("🔐 Enterprise Encryption Engine initialized")
    
    async def encrypt_data(self, 
                         plaintext: bytes, 
                         sensitivity: SensitivityLevel,
                         layers: List[EncryptionLayer] = None) -> EncryptionResult:
        """
        Encrypt data with appropriate layers based on sensitivity
        """
        if layers is None:
            layers = self._get_layers_for_sensitivity(sensitivity)
        
        encryption_chain = []
        current_data = plaintext
        
        # Apply encryption layers sequentially
        for layer in layers:
            try:
                encryptor = self.encryption_schemes[layer]
                result = await encryptor.encrypt(current_data, sensitivity)
                encryption_chain.append({
                    'layer': layer.value,
                    'algorithm': result.encryption_metadata['algorithm'],
                    'key_id': result.key_id
                })
                current_data = result.encrypted_data
                
            except Exception as e:
                logger.error(f"Encryption failed at layer {layer}: {e}")
                raise EncryptionError(f"Layer {layer} encryption failed") from e
        
        final_result = EncryptionResult(
            encrypted_data=current_data,
            iv=result.iv,
            tag=result.tag,
            key_id=result.key_id,
            encryption_metadata={
                'sensitivity': sensitivity.value,
                'layers_applied': encryption_chain,
                'timestamp': self._get_current_timestamp()
            }
        )
        
        await self.performance_monitor.record_encryption_operation(
            len(plaintext), sensitivity, layers
        )
        
        return final_result
    
    async def decrypt_data(self, 
                         encrypted_result: EncryptionResult,
                         key_material: Dict = None) -> bytes:
        """
        Decrypt data through the encryption layers
        """
        layers_applied = encrypted_result.encryption_metadata['layers_applied']
        current_data = encrypted_result.encrypted_data
        
        # Apply decryption in reverse order
        for layer_info in reversed(layers_applied):
            try:
                layer = EncryptionLayer(layer_info['layer'])
                decryptor = self.encryption_schemes[layer]
                current_data = await decryptor.decrypt(
                    current_data, 
                    encrypted_result.iv,
                    encrypted_result.tag,
                    key_material
                )
                
            except Exception as e:
                logger.error(f"Decryption failed at layer {layer}: {e}")
                raise DecryptionError(f"Layer {layer} decryption failed") from e
        
        return current_data
    
    def _get_layers_for_sensitivity(self, sensitivity: SensitivityLevel) -> List[EncryptionLayer]:
        """
        Determine encryption layers based on data sensitivity
        """
        sensitivity_profiles = {
            SensitivityLevel.LOW: [EncryptionLayer.IN_TRANSIT],
            SensitivityLevel.MEDIUM: [EncryptionLayer.IN_TRANSIT, EncryptionLayer.AT_REST],
            SensitivityLevel.HIGH: [EncryptionLayer.IN_TRANSIT, EncryptionLayer.AT_REST, 
                                  EncryptionLayer.IN_MEMORY],
            SensitivityLevel.CRITICAL: [EncryptionLayer.IN_TRANSIT, EncryptionLayer.AT_REST,
                                      EncryptionLayer.IN_MEMORY, EncryptionLayer.BACKUP]
        }
        
        return sensitivity_profiles[sensitivity]
