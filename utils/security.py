"""
Security Utilities for Hyper-ZiLLA
Cryptography, hashing, and security helpers
"""

import hashlib
import hmac
import secrets
import string
from typing import Optional
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import base64


class SecurityUtils:
    """Security utility functions"""

    @staticmethod
    def generate_secure_token(length: int = 32) -> str:
        """Generate cryptographically secure random token"""
        alphabet = string.ascii_letters + string.digits
        return "".join(secrets.choice(alphabet) for _ in range(length))

    @staticmethod
    def hash_password(password: str, salt: Optional[bytes] = None) -> tuple[bytes, bytes]:
        """Hash password with salt using PBKDF2"""
        if salt is None:
            salt = secrets.token_bytes(16)

        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=32,
            salt=salt,
            iterations=100000,
        )
        key = base64.urlsafe_b64encode(kdf.derive(password.encode()))
        return key, salt

    @staticmethod
    def verify_password(password: str, key: bytes, salt: bytes) -> bool:
        """Verify password against stored hash"""
        try:
            new_key, _ = SecurityUtils.hash_password(password, salt)
            return hmac.compare_digest(key, new_key)
        except Exception:
            return False

    @staticmethod
    def create_fernet_key(password: str, salt: Optional[bytes] = None) -> Fernet:
        """Create Fernet encryption instance from password"""
        key, salt = SecurityUtils.hash_password(password, salt)
        return Fernet(key)

    @staticmethod
    def secure_compare(a: str, b: str) -> bool:
        """Constant-time string comparison to prevent timing attacks"""
        return hmac.compare_digest(a.encode(), b.encode())


class DataSanitizer:
    """Data sanitization and validation utilities"""

    @staticmethod
    def sanitize_filename(filename: str) -> str:
        """Sanitize filename to prevent path traversal"""
        import re

        # Remove path components and special characters
        sanitized = re.sub(r"[^\w\-. ]", "", filename)
        return sanitized.replace("..", "").strip()

    @staticmethod
    def validate_url(url: str) -> bool:
        """Validate URL format"""
        import re

        pattern = re.compile(
            r"^(https?://)"  # http:// or https://
            r"(([A-Z0-9]([A-Z0-9-]{0,61}[A-Z0-9])?\.)+[A-Z]{2,6}\.?|"  # domain
            r"localhost|"  # localhost
            r"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})"  # IP
            r"(?::\d+)?"  # port
            r"(?:/?|[/?]\S+)$",
            re.IGNORECASE,
        )
        return pattern.match(url) is not None


# Example usage
if __name__ == "__main__":
    # Test security utilities
    utils = SecurityUtils()
    token = utils.generate_secure_token()
    print(f"Generated token: {token}")

    password = "my_secure_password"
    key, salt = utils.hash_password(password)
    print(f"Password verification: {utils.verify_password(password, key, salt)}")
