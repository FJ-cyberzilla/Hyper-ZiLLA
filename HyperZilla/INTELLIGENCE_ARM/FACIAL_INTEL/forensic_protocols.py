# ~/HyperZilla/INTELLIGENCE_ARM/FACIAL_INTEL/forensic_protocols.py
import cv2
import numpy as np
from PIL import Image, ImageEnhance, ImageFilter
import json
from datetime import datetime
from typing import Dict, List, Tuple
import asyncio

class ForensicImageAnalyzer:
    """ANALYZE PROBE IMAGES - Professional image assessment"""
    
    def __init__(self):
        self.quality_metrics = {}
        self.enhancement_log = []
    
    async def analyze_probe_image(self, image_path: str) -> Dict:
        """Comprehensive probe image analysis per professional protocols"""
        print("🔍 ANALYZING PROBE IMAGE FOR FORENSIC SUITABILITY...")
        
        analysis = {
            'basic_metrics': await self._calculate_basic_metrics(image_path),
            'facial_visibility': await self._assess_facial_visibility(image_path),
            'image_quality': await self._assess_image_quality(image_path),
            'enhancement_recommendations': [],
            'suitability_score': 0.0,
            'professional_judgment': ''
        }
        
        # Calculate overall suitability score
        analysis['suitability_score'] = self._calculate_suitability_score(analysis)
        analysis['professional_judgment'] = self._provide_professional_judgment(analysis)
        analysis['enhancement_recommendations'] = self._generate_enhancement_recommendations(analysis)
        
        return analysis
    
    async def _calculate_basic_metrics(self, image_path: str) -> Dict:
        """Calculate basic image metrics"""
        image = cv2.imread(image_path)
        if image is None:
            return {}
            
        height, width, channels = image.shape
        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        
        return {
            'resolution': f"{width}x{height}",
            'file_size_kb': os.path.getsize(image_path) / 1024,
            'brightness': np.mean(gray),
            'contrast': np.std(gray),
            'sharpness': self._calculate_sharpness(gray),
            'color_channels': channels
        }
    
    async def _assess_facial_visibility(self, image_path: str) -> Dict:
        """Assess facial visibility and orientation"""
        try:
            from .facial_ai_engine import HyperZillaFacialAI
            facial_ai = HyperZillaFacialAI()
            face_data = await facial_ai._extract_facial_features(image_path)
            
            if not face_data:
                return {'faces_detected': 0, 'primary_face_quality': 'NO_FACE'}
            
            primary_face = self._analyze_primary_face(face_data)
            
            return {
                'faces_detected': len(face_data['encodings']),
                'primary_face_quality': primary_face['quality'],
                'face_angle': primary_face['angle'],
                'visibility_score': primary_face['visibility_score'],
                'landmarks_detected': primary_face['landmarks_count']
            }
            
        except Exception as e:
            return {'faces_detected': 0, 'primary_face_quality': 'ANALYSIS_ERROR'}
    
    def _analyze_primary_face(self, face_data: Dict) -> Dict:
        """Analyze the primary face in detail"""
        if not face_data.get('dlib_rects'):
            return {'quality': 'NO_FACE', 'angle': 0, 'visibility_score': 0}
        
        rect = face_data['dlib_rects'][0]
        landmarks = self._get_facial_landmarks(face_data)
        
        # Calculate face angle (simplified)
        face_angle = self._estimate_face_angle(landmarks)
        
        # Calculate visibility score
        visibility_score = self._calculate_visibility_score(rect, landmarks, face_angle)
        
        # Determine quality category
        quality = self._determine_face_quality(visibility_score, face_angle)
        
        return {
            'quality': quality,
            'angle': face_angle,
            'visibility_score': visibility_score,
            'landmarks_count': landmarks.num_parts if landmarks else 0
        }
    
    def _estimate_face_angle(self, landmarks) -> float:
        """Estimate face rotation angle"""
        if not landmarks or landmarks.num_parts < 68:
            return 0.0
        
        # Simple angle estimation based on eye positions
        try:
            left_eye_center = np.mean([(landmarks.part(i).x, landmarks.part(i).y) for i in range(36, 42)], axis=0)
            right_eye_center = np.mean([(landmarks.part(i).x, landmarks.part(i).y) for i in range(42, 48)], axis=0)
            
            dx = right_eye_center[0] - left_eye_center[0]
            angle = np.degrees(np.arctan2(dx, 50))  # Simplified calculation
            
            return abs(angle)
        except:
            return 0.0
    
    def _calculate_visibility_score(self, rect, landmarks, angle: float) -> float:
        """Calculate face visibility score (0-100)"""
        score = 100.0
        
        # Penalize off-axis faces
        score -= min(angle * 2, 40)  # Up to 40 point penalty for angle
        
        # Penalize small faces
        face_size = rect.width() * rect.height()
        if face_size < 10000:  # Small face
            score -= 20
        elif face_size < 5000:  # Very small face
            score -= 40
        
        # Penalize incomplete landmarks
        if landmarks and landmarks.num_parts < 68:
            score -= (68 - landmarks.num_parts) * 0.5
        
        return max(score, 0)
    
    def _determine_face_quality(self, visibility_score: float, angle: float) -> str:
        """Determine professional face quality category"""
        if visibility_score >= 80 and angle <= 15:
            return "EXCELLENT"
        elif visibility_score >= 60 and angle <= 30:
            return "GOOD" 
        elif visibility_score >= 40 and angle <= 45:
            return "FAIR"
        elif visibility_score >= 20:
            return "POOR"
        else:
            return "UNSUITABLE"
    
    async def _assess_image_quality(self, image_path: str) -> Dict:
        """Assess overall image quality"""
        image = cv2.imread(image_path)
        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        
        # Calculate blurriness
        blur_value = cv2.Laplacian(gray, cv2.CV_64F).var()
        
        # Calculate noise
        noise_value = np.std(gray)
        
        # Calculate dynamic range
        hist = cv2.calcHist([gray], [0], None, [256], [0, 256])
        dynamic_range = np.sum(hist > 0) / 256.0
        
        return {
            'blurriness': blur_value,
            'blur_quality': 'SHARP' if blur_value > 100 else 'BLURRED' if blur_value > 50 else 'VERY_BLURRED',
            'noise_level': noise_value,
            'noise_quality': 'LOW' if noise_value < 20 else 'MEDIUM' if noise_value < 40 else 'HIGH',
            'dynamic_range': dynamic_range,
            'exposure_quality': 'GOOD' if 0.1 < dynamic_range < 0.9 else 'POOR'
        }
    
    def _calculate_suitability_score(self, analysis: Dict) -> float:
        """Calculate overall suitability score for facial recognition"""
        score = 0.0
        
        # Face visibility contributes 60%
        face_vis = analysis['facial_visibility']
        if face_vis.get('faces_detected', 0) > 0:
            score += (face_vis.get('visibility_score', 0) / 100) * 60
        
        # Image quality contributes 40%
        img_quality = analysis['image_quality']
        quality_factors = [
            1.0 if img_quality.get('blur_quality') == 'SHARP' else 0.5,
            1.0 if img_quality.get('noise_quality') == 'LOW' else 0.5,
            1.0 if img_quality.get('exposure_quality') == 'GOOD' else 0.5
        ]
        score += (sum(quality_factors) / len(quality_factors)) * 40
        
        return score
    
    def _provide_professional_judgment(self, analysis: Dict) -> str:
        """Provide professional judgment per forensic protocols"""
        score = analysis['suitability_score']
        face_quality = analysis['facial_visibility'].get('primary_face_quality', 'UNSUITABLE')
        
        if score >= 80 and face_quality in ['EXCELLENT', 'GOOD']:
            return "HIGHLY SUITABLE - Proceed with facial recognition search"
        elif score >= 60 and face_quality in ['GOOD', 'FAIR']:
            return "MODERATELY SUITABLE - Consider image enhancement before search"
        elif score >= 40:
            return "MARGINALLY SUITABLE - Enhancement required, results may be limited"
        else:
            return "UNSUITABLE - Do not proceed without significant enhancement or alternative imagery"
    
    def _generate_enhancement_recommendations(self, analysis: Dict) -> List[str]:
        """Generate professional enhancement recommendations"""
        recommendations = []
        face_vis = analysis['facial_visibility']
        img_quality = analysis['image_quality']
        
        # Face-related recommendations
        if face_vis.get('faces_detected', 0) == 0:
            recommendations.append("NO FACE DETECTED - Cannot proceed with facial recognition")
            return recommendations
        
        if face_vis.get('primary_face_quality') in ['POOR', 'UNSUITABLE']:
            recommendations.append("FACE QUALITY ISSUE - Consider alternative source imagery")
        
        if face_vis.get('face_angle', 0) > 30:
            recommendations.append(f"FACE ANGLE {face_vis['face_angle']:.1f}° - Consider rotation correction")
        
        # Image quality recommendations
        if img_quality.get('blur_quality') in ['BLURRED', 'VERY_BLURRED']:
            recommendations.append("IMAGE BLUR DETECTED - Apply deblurring enhancement")
        
        if img_quality.get('noise_quality') == 'HIGH':
            recommendations.append("HIGH NOISE LEVEL - Apply noise reduction")
        
        if img_quality.get('exposure_quality') == 'POOR':
            recommendations.append("POOR EXPOSURE - Adjust brightness/contrast")
        
        return recommendations

class ForensicImageEnhancer:
    """IN-LINE IMAGE ENHANCEMENT with audit trail"""
    
    def __init__(self):
        self.enhancement_log = []
        self.original_image = None
    
    async def enhance_probe_image(self, image_path: str, enhancements: Dict) -> Dict:
        """Apply professional image enhancements with audit trail"""
        print("🛠️ APPLYING FORENSIC IMAGE ENHANCEMENTS...")
        
        # Load original image
        self.original_image = Image.open(image_path)
        enhanced_image = self.original_image.copy()
        
        # Apply requested enhancements
        enhancement_results = {}
        
        for enhancement, params in enhancements.items():
            if hasattr(self, f'_enhance_{enhancement}'):
                enhanced_image, result = await getattr(self, f'_enhance_{enhancement}')(enhanced_image, params)
                enhancement_results[enhancement] = result
        
        # Save enhanced image
        enhanced_path = f"/tmp/enhanced_{os.path.basename(image_path)}"
        enhanced_image.save(enhanced_path)
        
        # Log all enhancements
        self._log_enhancements(enhancements, enhancement_results)
        
        return {
            'enhanced_image_path': enhanced_path,
            'enhancements_applied': list(enhancements.keys()),
            'enhancement_results': enhancement_results,
            'audit_trail': self.enhancement_log
        }
    
    async def _enhance_rotate(self, image: Image.Image, params: Dict) -> Tuple[Image.Image, Dict]:
        """Rotate image"""
        angle = params.get('angle', 0)
        rotated = image.rotate(angle, expand=True)
        
        result = {
            'original_orientation': 'as_captured',
            'rotation_applied': angle,
            'timestamp': datetime.now().isoformat()
        }
        
        self.enhancement_log.append({
            'operation': 'rotate',
            'parameters': params,
            'result': result,
            'timestamp': datetime.now().isoformat()
        })
        
        return rotated, result
    
    async def _enhance_brightness(self, image: Image.Image, params: Dict) -> Tuple[Image.Image, Dict]:
        """Adjust brightness"""
        factor = params.get('factor', 1.0)
        enhancer = ImageEnhance.Brightness(image)
        enhanced = enhancer.enhance(factor)
        
        result = {
            'original_brightness': 'as_captured',
            'brightness_factor': factor,
            'timestamp': datetime.now().isoformat()
        }
        
        self.enhancement_log.append({
            'operation': 'brightness',
            'parameters': params,
            'result': result,
            'timestamp': datetime.now().isoformat()
        })
        
        return enhanced, result
    
    async def _enhance_contrast(self, image: Image.Image, params: Dict) -> Tuple[Image.Image, Dict]:
        """Adjust contrast"""
        factor = params.get('factor', 1.0)
        enhancer = ImageEnhance.Contrast(image)
        enhanced = enhancer.enhance(factor)
        
        result = {
            'original_contrast': 'as_captured',
            'contrast_factor': factor,
            'timestamp': datetime.now().isoformat()
        }
        
        self.enhancement_log.append({
            'operation': 'contrast',
            'parameters': params,
            'result': result,
            'timestamp': datetime.now().isoformat()
        })
        
        return enhanced, result
    
    async def _enhance_sharpen(self, image: Image.Image, params: Dict) -> Tuple[Image.Image, Dict]:
        """Sharpen image"""
        factor = params.get('factor', 1.0)
        enhancer = ImageEnhance.Sharpness(image)
        enhanced = enhancer.enhance(factor)
        
        result = {
            'original_sharpness': 'as_captured',
            'sharpness_factor': factor,
            'timestamp': datetime.now().isoformat()
        }
        
        self.enhancement_log.append({
            'operation': 'sharpen',
            'parameters': params,
            'result': result,
            'timestamp': datetime.now().isoformat()
        })
        
        return enhanced, result
    
    async def _enhance_deblur(self, image: Image.Image, params: Dict) -> Tuple[Image.Image, Dict]:
        """Apply deblurring"""
        # Simple deblurring using sharpening filter
        enhanced = image.filter(ImageFilter.SHARPEN)
        
        result = {
            'deblurring_method': 'sharpening_filter',
            'iterations': 1,
            'timestamp': datetime.now().isoformat()
        }
        
        self.enhancement_log.append({
            'operation': 'deblur',
            'parameters': params,
            'result': result,
            'timestamp': datetime.now().isoformat()
        })
        
        return enhanced, result
    
    def _log_enhancements(self, enhancements: Dict, results: Dict):
        """Log all enhancements for audit trail"""
        audit_entry = {
            'timestamp': datetime.now().isoformat(),
            'original_image_properties': self._get_image_properties(self.original_image),
            'enhancements_requested': enhancements,
            'enhancements_applied': results,
            'analyst_notes': 'Enhanced for facial recognition suitability'
        }
        
        self.enhancement_log.append(audit_entry)

class UniqueMarksAnalyzer:
    """UNIQUE MARKS ANALYSIS - Scars, moles, tattoos, etc."""
    
    def __init__(self):
        self.feature_detector = cv2.SIFT_create()
    
    async def analyze_unique_marks(self, image_path: str, face_data: Dict) -> Dict:
        """Analyze unique facial marks and features"""
        print("🔎 ANALYZING UNIQUE FACIAL MARKS...")
        
        image = cv2.imread(image_path)
        if image is None:
            return {}
        
        analysis = {
            'facial_marks': await self._detect_facial_marks(image, face_data),
            'ear_characteristics': await self._analyze_ear_characteristics(image, face_data),
            'hair_characteristics': await self._analyze_hair_characteristics(image, face_data),
            'skin_texture': await self._analyze_skin_texture(image, face_data),
            'unique_identifiers': []
        }
        
        # Compile unique identifiers
        analysis['unique_identifiers'] = self._compile_unique_identifiers(analysis)
        
        return analysis
    
    async def _detect_facial_marks(self, image: np.ndarray, face_data: Dict) -> Dict:
        """Detect scars, moles, tattoos"""
        marks = {
            'scars': [],
            'moles': [],
            'tattoos': [],
            'other_marks': []
        }
        
        if not face_data.get('dlib_rects'):
            return marks
        
        # Extract face region
        face_rect = face_data['dlib_rects'][0]
        face_roi = image[face_rect.top():face_rect.bottom(), face_rect.left():face_rect.right()]
        
        if face_roi.size == 0:
            return marks
        
        # Convert to grayscale for analysis
        gray_face = cv2.cvtColor(face_roi, cv2.COLOR_BGR2GRAY)
        
        # Detect dark spots (potential moles)
        moles = self._detect_dark_spots(gray_face)
        marks['moles'] = moles
        
        # Detect linear features (potential scars)
        scars = self._detect_linear_features(gray_face)
        marks['scars'] = scars
        
        return marks
    
    def _detect_dark_spots(self, gray_image: np.ndarray) -> List[Dict]:
        """Detect dark spots that could be moles"""
        # Simple threshold-based detection
        _, thresh = cv2.threshold(gray_image, 50, 255, cv2.THRESH_BINARY_INV)
        contours, _ = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
        
        moles = []
        for contour in contours:
            area = cv2.contourArea(contour)
            if 10 < area < 1000:  # Reasonable mole size range
                x, y, w, h = cv2.boundingRect(contour)
                moles.append({
                    'type': 'mole',
                    'position': (x, y),
                    'size': area,
                    'bounding_box': (x, y, w, h)
                })
        
        return moles
    
    def _detect_linear_features(self, gray_image: np.ndarray) -> List[Dict]:
        """Detect linear features that could be scars"""
        # Edge detection for linear features
        edges = cv2.Canny(gray_image, 50, 150)
        lines = cv2.HoughLinesP(edges, 1, np.pi/180, threshold=30, minLineLength=20, maxLineGap=10)
        
        scars = []
        if lines is not None:
            for line in lines:
                x1, y1, x2, y2 = line[0]
                length = np.sqrt((x2-x1)**2 + (y2-y1)**2)
                if length > 15:  # Minimum scar length
                    scars.append({
                        'type': 'linear_feature',
                        'start_point': (x1, y1),
                        'end_point': (x2, y2),
                        'length': length
                    })
        
        return scars
    
    async def _analyze_ear_characteristics(self, image: np.ndarray, face_data: Dict) -> Dict:
        """Analyze ear shape and characteristics"""
        # This would require ear detection and analysis
        # Simplified implementation
        return {
            'ear_visibility': 'PARTIAL',  # FULL, PARTIAL, HIDDEN
            'ear_shape': 'unknown',
            'unique_characteristics': []
        }
    
    async def _analyze_hair_characteristics(self, image: np.ndarray, face_data: Dict) -> Dict:
        """Analyze hairline and hair characteristics"""
        return {
            'hair_visibility': 'VISIBLE',
            'hairline_shape': 'unknown',  # straight, widow's peak, etc.
            'hair_color': 'unknown',
            'hair_texture': 'unknown'
        }
    
    async def _analyze_skin_texture(self, image: np.ndarray, face_data: Dict) -> Dict:
        """Analyze skin texture
