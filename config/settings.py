"""
Hyper-ZiLLA Proprietary AI Configuration
All settings for custom AI implementations
"""

import os
from pathlib import Path

class BaseConfig:
    """Base configuration class"""
    BASE_DIR = Path(__file__).parent.parent
    LOGS_DIR = BASE_DIR / "logs"
    DATA_DIR = BASE_DIR / "data"
    CACHE_DIR = BASE_DIR / "cache"
    AI_MODELS_DIR = BASE_DIR / "models"

    # Ensure directories exist
    for directory in [LOGS_DIR, DATA_DIR, CACHE_DIR, AI_MODELS_DIR]:
        directory.mkdir(exist_ok=True)

    SECRET_KEY = os.environ.get('HZ_SECRET_KEY', 'dev-secret-key-change-in-production')
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    
class DevelopmentConfig(BaseConfig):
    """Development configuration"""
    DEBUG = True
    SQLALCHEMY_DATABASE_URI = os.environ.get('HZ_DATABASE_URL', f'sqlite:///{BaseConfig.DATA_DIR}/hyperzilla_dev.db')
class TestingConfig(BaseConfig):
    """Testing configuration"""
    TESTING = True
    SQLALCHEMY_DATABASE_URI = os.environ.get('HZ_DATABASE_URL', 'sqlite:///:memory:')

class ProductionConfig(BaseConfig):
    """Production configuration"""
    DEBUG = False
    SQLALCHEMY_DATABASE_URI = os.environ.get('HZ_DATABASE_URL')

config_by_name = {
    'development': DevelopmentConfig,
    'testing': TestingConfig,
    'production': ProductionConfig,
    'default': DevelopmentConfig
}