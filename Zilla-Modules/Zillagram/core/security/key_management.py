# core/security/key_management.py
import secrets
import time
from typing import Dict, List, Optional
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa, ec
import asyncio
import aiosqlite

class KeyManagementSystem:
    """
    Enterprise Key Management System
    Implements key rotation, secure storage, and access controls
    """
    
    def __init__(self, config: Dict):
        self.config = config
        self.db_path = config['encryption']['key_database']
        self.master_key = self._load_master_key()
        self.key_cache = {}
        self.rotation_scheduler = KeyRotationScheduler()
        
        # Initialize key database
        asyncio.create_task(self._initialize_key_database())
    
    async def generate_data_key(self, 
                              key_type: str = "AES-256",
                              metadata: Dict = None) -> Dict:
        """
        Generate a new data encryption key
        """
        key_id = self._generate_key_id()
        
        if key_type == "AES-256":
            key_material = secrets.token_bytes(32)  # 256-bit key
        elif key_type == "AES-192":
            key_material = secrets.token_bytes(24)  # 192-bit key
        elif key_type == "AES-128":
            key_material = secrets.token_bytes(16)  # 128-bit key
        else:
            raise ValueError(f"Unsupported key type: {key_type}")
        
        # Encrypt the data key with master key
        encrypted_key = await self._encrypt_with_master_key(key_material)
        
        key_record = {
            'key_id': key_id,
            'key_type': key_type,
            'encrypted_key': encrypted_key,
            'created_at': time.time(),
            'metadata': metadata or {},
            'status': 'active',
            'version': 1
        }
        
        # Store in database
        await self._store_key_record(key_record)
        
        # Cache for performance
        self.key_cache[key_id] = {
            'plaintext': key_material,
            'record': key_record
        }
        
        logger.info(f"🔑 Generated new data key: {key_id}")
        return key_record
    
    async def get_key(self, key_id: str) -> bytes:
        """
        Retrieve decrypted key material
        """
        # Check cache first
        if key_id in self.key_cache:
            return self.key_cache[key_id]['plaintext']
        
        # Retrieve from database
        key_record = await self._retrieve_key_record(key_id)
        if not key_record:
            raise KeyError(f"Key not found: {key_id}")
        
        # Decrypt with master key
        plaintext_key = await self._decrypt_with_master_key(key_record['encrypted_key'])
        
        # Update cache
        self.key_cache[key_id] = {
            'plaintext': plaintext_key,
            'record': key_record
        }
        
        return plaintext_key
    
    async def rotate_key(self, key_id: str) -> str:
        """
        Rotate encryption key (cryptographic best practice)
        """
        old_key_record = await self._retrieve_key_record(key_id)
        if not old_key_record:
            raise KeyError(f"Key not found for rotation: {key_id}")
        
        # Generate new key
        new_key_record = await self.generate_data_key(
            key_type=old_key_record['key_type'],
            metadata=old_key_record['metadata']
        )
        
        # Mark old key as deprecated but keep for decryption
        old_key_record['status'] = 'deprecated'
        old_key_record['superseded_by'] = new_key_record['key_id']
        await self._update_key_record(old_key_record)
        
        logger.info(f"🔄 Rotated key {key_id} -> {new_key_record['key_id']}")
        return new_key_record['key_id']
    
    async def _encrypt_with_master_key(self, data: bytes) -> bytes:
        """
        Encrypt data using the master key
        """
        # Implementation depends on master key storage strategy
        # This could use HSM, cloud KMS, or secure local storage
        iv = secrets.token_bytes(16)
        cipher = Cipher(algorithms.AES(self.master_key), modes.GCM(iv), backend=self.backend)
        encryptor = cipher.encryptor()
        
        encrypted = encryptor.update(data) + encryptor.finalize()
        return iv + encryptor.tag + encrypted
    
    def _generate_key_id(self) -> str:
        """Generate unique key identifier"""
        timestamp = int(time.time() * 1000)
        random = secrets.token_bytes(8)
        return f"key_{timestamp}_{random.hex()}"
    
    async def _initialize_key_database(self):
        """Initialize secure key database"""
        async with aiosqlite.connect(self.db_path) as db:
            await db.execute('''
                CREATE TABLE IF NOT EXISTS encryption_keys (
                    key_id TEXT PRIMARY KEY,
                    key_type TEXT NOT NULL,
                    encrypted_key BLOB NOT NULL,
                    created_at REAL NOT NULL,
                    metadata TEXT,
                    status TEXT DEFAULT 'active',
                    superseded_by TEXT,
                    version INTEGER DEFAULT 1,
                    last_used REAL
                )
            ''')
            await db.execute('''
                CREATE INDEX IF NOT EXISTS idx_key_status 
                ON encryption_keys(status)
            ''')
            await db.commit()
