"""
Core Anti-Tracking System
Enhanced with proper Python syntax and error handling
"""

import re
import random
import time
import hashlib
from typing import List, Dict
from urllib.parse import urlparse, urlunparse
import json

class AntiTrackingCore:
    def __init__(self):
        self.tracker_patterns = [
            r'google-analytics\.com',
            r'googlesyndication\.com',
            r'doubleclick\.net',
            r'facebook\.com/tr',
            r'connect\.facebook\.net',
            r'twitter\.com',
            r'analytics\.twitter\.com',
            r'linkedin\.com/analytics',
            r'bscore\.net',
            r'hotjar\.com',
            r'clicktale\.net',
            r'addthis\.com',
            r'sharethis\.com',
            r'googletagmanager\.com',
            r'googletagservices\.com',
            r'googleadservices\.com',
            r'googlesyndication\.com'
        ]
        
        self.compiled_patterns = [re.compile(pattern) for pattern in self.tracker_patterns]
        self.blocked_requests = set()
        self.user_agents = self._load_user_agents()
        
    def _load_user_agents(self) -> List[str]:
        """Load realistic user agents for rotation"""
        return [
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0"
        ]
    
    def is_tracker_url(self, url: str) -> bool:
        """
        Check if URL matches known tracker patterns
        
        Args:
            url: URL to check
            
        Returns:
            Boolean indicating if URL is a tracker
        """
        try:
            for pattern in self.compiled_patterns:
                if pattern.search(url):
                    return True
            return False
        except Exception as e:
            print(f"Error checking tracker URL {url}: {e}")
            return False
    
    def sanitize_url(self, url: str) -> str:
        """
        Remove tracking parameters from URL
        
        Args:
            url: Original URL with tracking parameters
            
        Returns:
            Sanitized URL without tracking parameters
        """
        try:
            parsed = urlparse(url)
            
            # Remove common tracking parameters
            tracking_params = {
                'utm_source', 'utm_medium', 'utm_campaign', 'utm_term', 'utm_content',
                'fbclid', 'gclid', 'msclkid', 'trk', 'mc_cid', 'mc_eid',
                '_ga', '_gl', 'yclid', 'igshid', 'si'
            }
            
            # Filter query parameters
            query_params = []
            if parsed.query:
                for param in parsed.query.split('&'):
                    key = param.split('=')[0] if '=' in param else param
                    if key not in tracking_params:
                        query_params.append(param)
            
            new_query = '&'.join(query_params)
            sanitized = urlunparse((
                parsed.scheme,
                parsed.netloc,
                parsed.path,
                parsed.params,
                new_query,
                parsed.fragment
            ))
            
            return sanitized
        except Exception as e:
            print(f"Error sanitizing URL {url}: {e}")
            return url
    
    def rotate_user_agent(self) -> str:
        """Return a random user agent for request rotation"""
        return random.choice(self.user_agents)
    
    def generate_fingerprint_hash(self) -> str:
        """
        Generate a randomized browser fingerprint hash
        
        Returns:
            Randomized fingerprint hash
        """
        fingerprint_components = [
            str(random.randint(1000, 9999)),  # Screen dimensions
            str(random.choice([60, 120, 144])),  # Refresh rate
            random.choice(['en-US', 'en-GB', 'fr-FR', 'de-DE']),  # Language
            str(random.choice([8, 16, 32])),  # Color depth
            str(time.time())  # Timestamp for uniqueness
        ]
        
        fingerprint_string = '|'.join(fingerprint_components)
        return hashlib.sha256(fingerprint_string.encode()).hexdigest()
    
    def block_tracking_requests(self, request_urls: List[str]) -> List[str]:
        """
        Filter out tracking requests from a list of URLs
        
        Args:
            request_urls: List of URLs to filter
            
        Returns:
            List of non-tracking URLs
        """
        clean_urls = []
        
        for url in request_urls:
            if not self.is_tracker_url(url):
                clean_url = self.sanitize_url(url)
                clean_urls.append(clean_url)
            else:
                self.blocked_requests.add(url)
                print(f"Blocked tracker: {url}")
        
        return clean_urls
    
    def get_stealth_headers(self) -> Dict[str, str]:
        """
        Generate stealth headers to avoid detection
        
        Returns:
            Dictionary of stealth headers
        """
        return {
            'User-Agent': self.rotate_user_agent(),
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
            'Accept-Language': 'en-US,en;q=0.5',
            'Accept-Encoding': 'gzip, deflate, br',
            'DNT': '1',
            'Connection': 'keep-alive',
            'Upgrade-Insecure-Requests': '1',
            'Sec-Fetch-Dest': 'document',
            'Sec-Fetch-Mode': 'navigate',
            'Sec-Fetch-Site': 'none',
            'Cache-Control': 'max-age=0'
        }
    
    def clear_tracking_data(self):
        """Clear all tracking-related data"""
        self.blocked_requests.clear()
    
    def get_blocking_stats(self) -> Dict:
        """Get statistics about blocked trackers"""
        return {
            'total_blocked': len(self.blocked_requests),
            'blocked_domains': list(self.blocked_requests),
            'patterns_loaded': len(self.compiled_patterns),
            'user_agents_available': len(self.user_agents)
        }

# Advanced anti-tracking features
class AdvancedAntiTracking:
    def __init__(self):
        self.core = AntiTrackingCore()
        self.session_fingerprint = self.core.generate_fingerprint_hash()
        
    def create_stealth_session(self) -> Dict:
        """Create a complete stealth browsing session"""
        return {
            'headers': self.core.get_stealth_headers(),
            'fingerprint': self.session_fingerprint,
            'session_id': hashlib.md5(str(time.time()).encode()).hexdigest(),
            'timestamp': time.time()
        }
    
    def advanced_url_sanitization(self, url: str, level: str = 'aggressive') -> str:
        """
        Advanced URL sanitization with multiple levels
        
        Args:
            url: URL to sanitize
            level: Sanitization level ('basic', 'aggressive', 'paranoid')
        """
        sanitized = self.core.sanitize_url(url)
        
        if level == 'aggressive':
            # Remove all query parameters
            parsed = urlparse(sanitized)
            sanitized = urlunparse((
                parsed.scheme,
                parsed.netloc,
                parsed.path,
                parsed.params,
                '',  # Empty query
                parsed.fragment
            ))
        elif level == 'paranoid':
            # Only keep scheme and netloc
            parsed = urlparse(sanitized)
            sanitized = urlunparse((
                parsed.scheme,
                parsed.netloc,
                '', '', '', ''
            ))
        
        return sanitized

# Example usage and testing
if __name__ == "__main__":
    # Test basic anti-tracking
    tracker = AntiTrackingCore()
    
    test_urls = [
        "https://example.com/page?utm_source=google",
        "https://google-analytics.com/collect",
        "https://example.com/clean-page"
    ]
    
    print("Testing URL sanitization:")
    for url in test_urls:
        clean = tracker.sanitize_url(url)
        is_tracker = tracker.is_tracker_url(url)
        print(f"Original: {url}")
        print(f"Clean: {clean}")
        print(f"Is Tracker: {is_tracker}")
        print("---")
    
    # Test advanced features
    advanced = AdvancedAntiTracking()
    session = advanced.create_stealth_session()
    print("Stealth Session:", json.dumps(session, indent=2))
    
    stats = tracker.get_blocking_stats()
    print("Blocking Stats:", json.dumps(stats, indent=2))
