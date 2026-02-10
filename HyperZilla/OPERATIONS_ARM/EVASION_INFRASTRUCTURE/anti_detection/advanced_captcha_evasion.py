# ~/HyperZilla/OPERATIONS_ARM/EVASION_INFRASTRUCTURE/anti_detection/advanced_captcha_evasion.py
import time
import random


from selenium.webdriver.common.action_chains import ActionChains

import numpy as np



import requests

class AdvancedCaptchaIntelligence:
    def __init__(self):
        self.captcha_cookies = {}  # Track CAPTCHA cookies across sessions
        self.behavioral_baselines = self._load_human_baselines()
        self.vpn_rotation = VPNManager()
        self.captcha_solving_service = CaptchaServiceAPI()
        self.history_emulator = BrowserHistoryEmulator()
        
    def analyze_recaptcha_v3_risk(self, driver):
        """Analyze reCAPTCHA v3 risk factors and return score improvement strategy"""
        risk_factors = {
            'mouse_movement_entropy': self._calculate_mouse_entropy(driver),
            'browser_fingerprint_consistency': self._check_fingerprint_consistency(driver),
            'cookie_history_depth': self._analyze_cookie_history(driver),
            'ip_reputation': self.vpn_rotation.get_ip_reputation(),
            'behavioral_anomalies': self._detect_behavioral_anomalies(driver)
        }
        
        risk_score = sum(risk_factors.values()) / len(risk_factors)
        return {
            'risk_score': risk_score,
            'improvement_actions': self._generate_improvement_plan(risk_factors),
            'expected_v3_score': max(0.9, 1.0 - risk_score)  # Target 0.9+ score
        }
    
    def _calculate_mouse_entropy(self, driver):
        """Calculate mouse movement randomness (human-like vs bot-like)"""
        # Record mouse movements during page interaction
        movements = self._record_mouse_movements(driver)
        
        # Calculate entropy of movement patterns
        if len(movements) < 10:
            return 0.8  # High risk - insufficient movement data
            
        movement_angles = []
        for i in range(1, len(movements)):
            dx = movements[i][0] - movements[i-1][0]
            dy = movements[i][1] - movements[i-1][1]
            angle = np.arctan2(dy, dx)
            movement_angles.append(angle)
        
        # Humans have higher entropy in movement patterns
        entropy = np.std(movement_angles)
        normalized_entropy = min(1.0, entropy / np.pi)  # Normalize to 0-1
        
        return 1.0 - normalized_entropy  # Lower entropy = higher risk
    
    def _check_fingerprint_consistency(self, driver):
        """Verify browser fingerprint consistency with historical data"""
        current_fingerprint = self._generate_browser_fingerprint(driver)
        
        if not hasattr(self, 'historical_fingerprints'):
            self.historical_fingerprints = []
            
        # Check consistency with previous fingerprints
        consistency_scores = []
        for historical_fp in self.historical_fingerprints[-5:]:  # Last 5 sessions
            similarity = self._fingerprint_similarity(current_fingerprint, historical_fp)
            consistency_scores.append(similarity)
            
        avg_consistency = np.mean(consistency_scores) if consistency_scores else 0.5
        self.historical_fingerprints.append(current_fingerprint)
        
        return 1.0 - avg_consistency  # Lower consistency = higher risk
    
    def _analyze_cookie_history(self, driver):
        """Analyze cookie history depth and consistency"""
        try:
            cookies = driver.get_cookies()
            captcha_cookies = [c for c in cookies if 'captcha' in c['name'].lower() or 'cf_' in c['name']]
            
            if not captcha_cookies:
                return 0.9  # High risk - no CAPTCHA history
                
            # Analyze cookie age and diversity
            cookie_ages = [c.get('expiry', 0) - time.time() for c in captcha_cookies if c.get('expiry')]
            avg_cookie_age = np.mean(cookie_ages) if cookie_ages else 0
            
            # Older, diverse cookies = lower risk
            if avg_cookie_age > 86400:  # More than 1 day old
                return 0.1
            elif avg_cookie_age > 3600:  # More than 1 hour old
                return 0.3
            else:
                return 0.7  # Fresh cookies = higher risk
                
        except Exception as e:
            print(f"Error analyzing cookie history: {e}")
            return 0.8  # Error case = assume higher risk

class VPNManager:
    def __init__(self):
        self.dedicated_ips = self._load_dedicated_ips()
        self.current_ip_index = 0
        self.ip_reputation_cache = {}
        
    def _load_dedicated_ips(self):
        """Load dedicated IPs to avoid shared IP CAPTCHA triggers"""
        # Your dedicated IP list - critical for reducing CAPTCHAs
        return [
            '192.168.1.100:8080',  # Residential IP 1
            '192.168.1.101:8080',  # Residential IP 2
            '192.168.1.102:8080',  # Residential IP 3
        ]
    
    def get_next_dedicated_ip(self):
        """Rotate through dedicated IPs"""
        ip = self.dedicated_ips[self.current_ip_index]
        self.current_ip_index = (self.current_ip_index + 1) % len(self.dedicated_ips)
        return ip
    
    def get_ip_reputation(self):
        """Get current IP's CAPTCHA reputation score"""
        current_ip = self.dedicated_ips[self.current_ip_index]
        if current_ip in self.ip_reputation_cache:
            return self.ip_reputation_cache[current_ip]
        
        # Simulate IP reputation check (0.0 = bad, 1.0 = good)
        reputation = random.uniform(0.7, 0.95)  # Dedicated IPs have better reputation
        self.ip_reputation_cache[current_ip] = reputation
        return reputation

class BrowserHistoryEmulator:
    def __init__(self):
        self.history_profiles = self._load_legitimate_history_profiles()
        
    def _load_legitimate_history_profiles(self):
        """Load patterns of legitimate browsing history"""
        return {
            'business_user': [
                'https://linkedin.com', 'https://gmail.com', 'https://calendar.google.com',
                'https://docs.google.com', 'https://slack.com', 'https://github.com'
            ],
            'casual_user': [
                'https://youtube.com', 'https://facebook.com', 'https://amazon.com',
                'https://reddit.com', 'https://twitter.com', 'https://instagram.com'
            ],
            'research_user': [
                'https://wikipedia.org', 'https://arxiv.org', 'https://stackoverflow.com',
                'https://medium.com', 'https://towardsdatascience.com'
            ]
        }
    
    def emulate_browsing_history(self, driver, profile_type='business_user'):
        """Emulate legitimate browsing history through cookie creation"""
        sites = self.history_profiles[profile_type]
        
        for site in random.sample(sites, min(3, len(sites))):
            try:
                driver.get(site)
                time.sleep(random.uniform(2, 5))  # Realistic browsing time
                
                # Simulate some interaction
                actions = ActionChains(driver)
                actions.move_by_offset(random.randint(10, 100), random.randint(10, 100))
                actions.click()
                actions.perform()
                
            except Exception:
                continue  # Continue with next site if one fails

class HumanMouseEmulator:
    def __init__(self):
        self.trajectory_generator = MouseTrajectoryGenerator()
        
    def perform_human_click(self, driver, element):
        """Perform human-like mouse movement to element and click"""
        # Get element location
        location = element.location
        size = element.size
        
        # Calculate target point (slightly random within element)
        target_x = location['x'] + random.randint(size['width']//4, 3*size['width']//4)
        target_y = location['y'] + random.randint(size['height']//4, 3*size['height']//4)
        
        # Generate human-like trajectory
        start_x, start_y = 0, 0  # Current mouse position
        trajectory = self.trajectory_generator.generate_trajectory(start_x, start_y, target_x, target_y)
        
        # Execute trajectory
        actions = ActionChains(driver)
        for point in trajectory:
            actions.move_by_offset(point[0], point[1])
            # Add micro-pauses
            if random.random() < 0.1:
                actions.pause(random.uniform(0.01, 0.05))
        
        actions.click()
        actions.perform()

class MouseTrajectoryGenerator:
    def generate_trajectory(self, start_x, start_y, end_x, end_y, points=50):
        """Generate human-like mouse trajectory using Bezier curves with noise"""
        # Control points for Bezier curve
        control1_x = start_x + (end_x - start_x) * random.uniform(0.3, 0.7)
        control1_y = start_y + (end_y - start_y) * random.uniform(0.2, 0.8)
        
        control2_x = start_x + (end_x - start_x) * random.uniform(0.3, 0.7)
        control2_y = start_y + (end_y - start_y) * random.uniform(0.2, 0.8)
        
        trajectory = []
        for i in range(points):
            t = i / (points - 1)
            
            # Bezier curve calculation
            x = (1-t)**3 * start_x + 3*(1-t)**2*t * control1_x + 3*(1-t)*t**2 * control2_x + t**3 * end_x
            y = (1-t)**3 * start_y + 3*(1-t)**2*t * control1_y + 3*(1-t)*t**2 * control2_y + t**3 * end_y
            
            # Add human-like noise
            if 0.2 < t < 0.8:  # Add more noise in the middle of movement
                x += random.gauss(0, 2)
                y += random.gauss(0, 2)
                
            trajectory.append((int(x), int(y)))
            
        return trajectory

class CaptchaServiceAPI:
    def __init__(self):
        self.api_key = "YOUR_2CAPTCHA_API_KEY"  # Replace with actual key
        self.service_url = "http://2captcha.com/in.php"
        self.result_url = "http://2captcha.com/res.php"
        
    def solve_captcha_via_service(self, image_data, captcha_type='image'):
        """Use 2Captcha service for difficult CAPTCHAs"""
        # Implementation for 2Captcha API
        payload = {
            'key': self.api_key,
            'method': 'base64',
            'body': image_data,
            'json': 1
        }
        
        try:
            # Send CAPTCHA to solving service
            response = requests.post(self.service_url, data=payload)
            request_id = response.json().get('request')
            
            # Poll for solution
            for _ in range(30):  # 30 attempts, 5 seconds apart
                time.sleep(5)
                result_payload = {
                    'key': self.api_key,
                    'action': 'get',
                    'id': request_id,
                    'json': 1
                }
                result_response = requests.get(self.result_url, params=result_payload)
                result_data = result_response.json()
                
                if result_data.get('status') == 1:
                    return result_data.get('request')  # Solved CAPTCHA text
                    
        except Exception as e:
            print(f"CAPTCHA service error: {e}")
            
        return None

# ENHANCED EVASION STRATEGY INTEGRATION
class ComprehensiveCaptchaEvasion:
    def __init__(self):
        self.captcha_intel = AdvancedCaptchaIntelligence()
        self.mouse_emulator = HumanMouseEmulator()
        self.vpn_manager = VPNManager()
        
    def execute_evasion_protocol(self, driver, target_url):
        """Execute comprehensive CAPTCHA evasion protocol"""
        # Phase 1: Pre-emptive measures
        self._apply_preemptive_evasion(driver)
        
        # Phase 2: Behavioral preparation
        self._establish_behavioral_credibility(driver)
        
        # Phase 3: Execute access attempt
        return self._access_with_evasion(driver, target_url)
    
    def _apply_preemptive_evasion(self, driver):
        """Apply measures before CAPTCHA encounter"""
        # Use dedicated IP
        dedicated_ip = self.vpn_manager.get_next_dedicated_ip()
        self._configure_proxy(driver, dedicated_ip)
        
        # Emulate browsing history
        self.captcha_intel.history_emulator.emulate_browsing_history(driver)
        
    def _establish_behavioral_credibility(self, driver):
        """Establish human-like behavioral patterns"""
        # Generate initial mouse movements
        self._perform_credibility_movements(driver)
        
        # Simulate reading behavior (scrolling, pauses)
        self._simulate_reading_behavior(driver)

    def _simulate_reading_behavior(self, driver):
        """Simulate human-like reading behavior (scrolling, pauses)"""
        print("Simulating reading behavior (scrolling and pauses)...")
        try:
            # Scroll down the page multiple times
            for _ in range(random.randint(2, 5)):
                scroll_amount = random.randint(200, 800)
                driver.execute_script(f"window.scrollBy(0, {scroll_amount});")
                time.sleep(random.uniform(0.5, 2.0)) # Pause to simulate reading

            # Scroll back up a bit sometimes
            if random.random() > 0.5:
                scroll_amount = random.randint(100, 300)
                driver.execute_script(f"window.scrollBy(0, -{scroll_amount});")
                time.sleep(random.uniform(0.3, 1.0))

        except Exception as e:
            print(f"Error simulating reading behavior: {e}")
        
    def _access_with_evasion(self, driver, url):
        """Access target with comprehensive evasion"""
        driver.get(url)
        
        # Analyze reCAPTCHA v3 risk
        risk_analysis = self.captcha_intel.analyze_recaptcha_v3_risk(driver)
        
        if risk_analysis['risk_score'] > 0.3:
            # High risk - implement aggressive evasion
            self._execute_aggressive_evasion(driver, risk_analysis)
        else:
            # Low risk - standard human emulation
            self._execute_standard_evasion(driver)
            
        return risk_analysis
