# config/encryption_config.py
ENCRYPTION_CONFIG = {
    'encryption': {
        'key_database': 'data/keys/encryption_keys.db',
        'master_key_path': 'secure/master.key',
        'key_rotation_days': 90,
        'performance_monitoring': True,
        'audit_logging': True,
        
        'algorithms': {
            'transit': 'AES-256-GCM',
            'at_rest': 'AES-256-CBC-HMAC-SHA256', 
            'in_memory': 'AES-256-CTR',
            'backup': 'AES-256-GCM'
        },
        
        'key_sizes': {
            'AES': 256,
            'HMAC': 256
        },
        
        'security_parameters': {
            'min_iv_size': 16,
            'min_key_size': 32,
            'max_key_age_days': 365
        }
    }
}

# Initialize in main system
async def initialize_enterprise_encryption():
    """Initialize enterprise encryption system"""
    encryption_engine = EnterpriseEncryptionEngine(ENCRYPTION_CONFIG)
    secure_data_manager = SecureDataManager(encryption_engine)
    
    logger.info("🔐 Enterprise Encryption System Ready")
    return secure_data_manager
