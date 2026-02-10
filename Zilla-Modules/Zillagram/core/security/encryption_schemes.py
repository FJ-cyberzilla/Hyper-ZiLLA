# core/security/encryption_schemes.py
class TransitEncryptionScheme:
    """
    Encryption for data in transit (network communications)
    Implements AES-GCM for authenticated encryption
    """
    
    def __init__(self, key_manager: KeyManagementSystem):
        self.key_manager = key_manager
        self.algorithm = "AES-256-GCM"
    
    async def encrypt(self, plaintext: bytes, sensitivity: SensitivityLevel) -> EncryptionResult:
        """Encrypt data for transit"""
        # Generate unique key for this operation
        key_record = await self.key_manager.generate_data_key("AES-256", {
            'purpose': 'transit_encryption',
            'sensitivity': sensitivity.value
        })
        
        key_material = await self.key_manager.get_key(key_record['key_id'])
        iv = secrets.token_bytes(16)  # 128-bit IV
        
        # AES-GCM encryption
        cipher = Cipher(algorithms.AES(key_material), modes.GCM(iv), backend=default_backend())
        encryptor = cipher.encryptor()
        
        # Add additional authenticated data
        encryptor.authenticate_additional_data(self._generate_aad(sensitivity))
        
        ciphertext = encryptor.update(plaintext) + encryptor.finalize()
        
        return EncryptionResult(
            encrypted_data=ciphertext,
            iv=iv,
            tag=encryptor.tag,
            key_id=key_record['key_id'],
            encryption_metadata={
                'algorithm': self.algorithm,
                'key_size': 256,
                'mode': 'GCM'
            }
        )
    
    async def decrypt(self, ciphertext: bytes, iv: bytes, tag: bytes, 
                     key_material: Dict = None) -> bytes:
        """Decrypt data from transit"""
        if not key_material:
            raise ValueError("Key material required for decryption")
        
        key = await self.key_manager.get_key(key_material['key_id'])
        
        cipher = Cipher(algorithms.AES(key), modes.GCM(iv, tag), backend=default_backend())
        decryptor = cipher.decryptor()
        
        # Verify additional authenticated data
        decryptor.authenticate_additional_data(self._generate_aad(key_material.get('sensitivity')))
        
        plaintext = decryptor.update(ciphertext) + decryptor.finalize()
        return plaintext
    
    def _generate_aad(self, sensitivity: SensitivityLevel) -> bytes:
        """Generate Additional Authenticated Data"""
        aad_data = {
            'sensitivity': sensitivity.value,
            'timestamp': int(time.time()),
            'purpose': 'telegram_osint_transit'
        }
        return json.dumps(aad_data, sort_keys=True).encode()

class AtRestEncryptionScheme:
    """
    Encryption for data at rest (database storage)
    Implements AES-CBC with HMAC for integrity
    """
    
    def __init__(self, key_manager: KeyManagementSystem):
        self.key_manager = key_manager
        self.algorithm = "AES-256-CBC-HMAC-SHA256"
    
    async def encrypt(self, plaintext: bytes, sensitivity: SensitivityLevel) -> EncryptionResult:
        """Encrypt data for storage"""
        key_record = await self.key_manager.generate_data_key("AES-256", {
            'purpose': 'at_rest_encryption', 
            'sensitivity': sensitivity.value
        })
        
        key_material = await self.key_manager.get_key(key_record['key_id'])
        iv = secrets.token_bytes(16)
        
        # Pad data to block size
        padder = padding.PKCS7(128).padder()
        padded_data = padder.update(plaintext) + padder.finalize()
        
        # AES-CBC encryption
        cipher = Cipher(algorithms.AES(key_material), modes.CBC(iv), backend=default_backend())
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(padded_data) + encryptor.finalize()
        
        # Generate HMAC for integrity
        hmac_key = hashlib.sha256(key_material + iv).digest()
        hmac_value = hmac.new(hmac_key, ciphertext, hashlib.sha256).digest()
        
        return EncryptionResult(
            encrypted_data=ciphertext + hmac_value,
            iv=iv,
            key_id=key_record['key_id'],
            encryption_metadata={
                'algorithm': self.algorithm,
                'key_size': 256,
                'mode': 'CBC',
                'integrity_check': 'HMAC-SHA256'
            }
        )

class MemoryEncryptionScheme:
    """
    Encryption for sensitive data in memory
    Implements protection against memory scraping attacks
    """
    
    def __init__(self, key_manager: KeyManagementSystem):
        self.key_manager = key_manager
        self.algorithm = "AES-256-CTR"
        self.ephemeral_keys = {}
    
    async def encrypt(self, plaintext: bytes, sensitivity: SensitivityLevel) -> EncryptionResult:
        """Encrypt data in memory with ephemeral keys"""
        # Use ephemeral keys for memory protection
        key_id = f"memory_{int(time.time() * 1000)}"
        key_material = secrets.token_bytes(32)
        
        # Store ephemeral key (short-lived)
        self.ephemeral_keys[key_id] = {
            'key': key_material,
            'created_at': time.time(),
            'sensitivity': sensitivity
        }
        
        iv = secrets.token_bytes(16)
        
        # AES-CTR for efficient memory encryption
        cipher = Cipher(algorithms.AES(key_material), modes.CTR(iv), backend=default_backend())
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(plaintext) + encryptor.finalize()
        
        # Schedule key cleanup
        asyncio.create_task(self._cleanup_ephemeral_key(key_id))
        
        return EncryptionResult(
            encrypted_data=ciphertext,
            iv=iv,
            key_id=key_id,
            encryption_metadata={
                'algorithm': self.algorithm,
                'ephemeral': True,
                'key_size': 256,
                'mode': 'CTR'
            }
        )
    
    async def _cleanup_ephemeral_key(self, key_id: str, ttl: int = 300):
        """Clean up ephemeral keys after TTL"""
        await asyncio.sleep(ttl)
        if key_id in self.ephemeral_keys:
            # Securely wipe from memory
            self.ephemeral_keys[key_id]['key'] = b'\x00' * 32
            del self.ephemeral_keys[key_id]
