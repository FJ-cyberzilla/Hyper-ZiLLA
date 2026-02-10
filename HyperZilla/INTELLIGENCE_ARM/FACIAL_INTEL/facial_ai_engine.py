# ~/HyperZilla/INTELLIGENCE_ARM/FACIAL_INTEL/facial_ai_engine.py
import cv2
import numpy as np
import face_recognition
import dlib
from deepface import DeepFace
import torch

from PIL import Image
import os
import json
from typing import Dict, List
import asyncio
from HyperZilla.INTELLIGENCE_ARM.DIGITAL_OSINT.integration_bridge import DigitalOSINTBridge as OSINTIntegration

class HyperZillaFacialAI:
    def __init__(self):
        self.face_detector = dlib.get_frontal_face_detector()
        self.shape_predictor = dlib.shape_predictor(self._get_model_path('shape_predictor_68_face_landmarks.dat'))
        self.face_recognizer = dlib.face_recognition_model_v1(self._get_model_path('dlib_face_recognition_resnet_model_v1.dat'))
        
        self.face_database = FaceDatabase()
        self.search_engines = FacialSearchEngines()
        self.analysis_pipeline = FacialAnalysisPipeline()
        
    def _get_model_path(self, model_name):
        """Get path to pre-trained models"""
        return f"INTELLIGENCE_ARM/FACIAL_INTEL/ai_models/{model_name}"
    
    async def identify_person_from_image(self, image_path: str, enhance_with_osint: bool = True) -> Dict:
        """Main facial identification pipeline"""
        print("🎯 STARTING FACIAL IDENTIFICATION PROTOCOL...")
        
        # Step 1: Face detection and encoding
        face_data = await self._extract_facial_features(image_path)
        
        if not face_data:
            return {"error": "No faces detected in image"}
        
        # Step 2: Database matching
        matches = await self._search_face_databases(face_data['encodings'][0])
        
        # Step 3: AI-powered analysis
        analysis = await self.analysis_pipeline.analyze_face(image_path, face_data)
        
        # Step 4: OSINT enhancement if requested
        if enhance_with_osint:
            osint_data = await self._enhance_with_osint(analysis, matches)
            analysis.update(osint_data)
        
        return {
            'success': True,
            'face_detected': True,
            'face_count': len(face_data['encodings']),
            'primary_face_analysis': analysis,
            'database_matches': matches,
            'confidence_score': self._calculate_confidence(analysis, matches),
            'next_steps': self._suggest_next_steps(analysis, matches)
        }
    
    async def _extract_facial_features(self, image_path: str) -> Dict:
        """Extract facial features and encodings"""
        try:
            # Load image
            image = face_recognition.load_image_file(image_path)
            
            # Detect faces
            face_locations = face_recognition.face_locations(image, model='hog')
            face_encodings = face_recognition.face_encodings(image, face_locations)
            
            if not face_encodings:
                return None
            
            # Convert to dlib format for advanced analysis
            dlib_rects = [dlib.rectangle(loc[3], loc[0], loc[1], loc[2]) for loc in face_locations]
            
            facial_data = {
                'locations': face_locations,
                'encodings': face_encodings,
                'dlib_rects': dlib_rects,
                'image_shape': image.shape
            }
            
            return facial_data
            
        except Exception as e:
            print(f"Face extraction error: {e}")
            return None
    
    async def _search_face_databases(self, face_encoding: np.ndarray) -> List[Dict]:
        """Search across multiple face databases"""
        matches = []
        
        # 1. Internal database search
        internal_matches = await self.face_database.find_similar_faces(face_encoding)
        matches.extend(internal_matches)
        
        # 2. External OSINT database search
        external_matches = await self.search_engines.search_public_databases(face_encoding)
        matches.extend(external_matches)
        
        # 3. Social media reverse image search
        social_matches = await self.search_engines.reverse_image_search(face_encoding)
        matches.extend(social_matches)
        
        return sorted(matches, key=lambda x: x['confidence'], reverse=True)[:10]  # Top 10 matches
    
    async def _enhance_with_osint(self, analysis: Dict, matches: List[Dict]) -> Dict:
        """Enhance facial analysis with OSINT data"""
        enhanced_data = {}
        
        # Use demographic data to narrow OSINT searches
        if 'demographics' in analysis:
            demographics = analysis['demographics']
            
            # Search social media with demographic filters
            social_searches = []
            if 'age_range' in demographics:
                social_searches.append(f"age:{demographics['age_range'][0]}-{demographics['age_range'][1]}")
            if 'gender' in demographics:
                social_searches.append(f"gender:{demographics['gender']}")
            if 'ethnicity' in demographics:
                social_searches.append(f"ethnicity:{demographics['ethnicity']}")
            
            # Enhanced social media search
            enhanced_data['social_media_profiles'] = await self._search_social_media_with_filters(social_searches)
        
        # Cross-reference with database matches
        if matches:
            best_match = matches[0]
            if best_match['confidence'] > 0.8:  # High confidence match
                enhanced_data['verified_identities'] = await self._verify_identity_osint(best_match['identity'])
        
        return enhanced_data
    
    async def _search_social_media_with_filters(self, filters: List[str]) -> List[Dict]:
        """Search social media with demographic filters"""
        # This would integrate with your existing OSINT systems
        from INTELLIGENCE_ARM.DIGITAL_OSINT.integration_bridge import DigitalOSINTBridge
        
        osint_bridge = DigitalOSINTBridge()
        search_queries = self._generate_search_queries(filters)
        
        profiles = []
        for query in search_queries:
            try:
                results = await osint_bridge.search_social_media(query)
                profiles.extend(results.get('profiles', []))
            except Exception as e:
                print(f"Social media search error: {e}")
        
        return profiles
    
    def _generate_search_queries(self, filters: List[str]) -> List[str]:
        """Generate search queries from demographic filters"""
        # Implement query generation logic based on filters
        base_queries = [
            "site:linkedin.com",
            "site:facebook.com", 
            "site:instagram.com",
            "site:twitter.com"
        ]
        
        enhanced_queries = []
        for base in base_queries:
            for filter_str in filters:
                enhanced_queries.append(f"{base} {filter_str}")
        
        return enhanced_queries

class FaceDatabase:
    """Manage internal face database"""
    
    def __init__(self):
        self.db_path = "INTELLIGENCE_ARM/FACIAL_INTEL/face_database/encodings.json"
        self.encodings = self._load_encodings()
    
    async def find_similar_faces(self, query_encoding: np.ndarray, threshold: float = 0.6) -> List[Dict]:
        """Find similar faces in database"""
        matches = []
        
        for identity, data in self.encodings.items():
            for encoding in data['encodings']:
                distance = np.linalg.norm(query_encoding - np.array(encoding))
                similarity = 1 - distance
                
                if similarity > threshold:
                    matches.append({
                        'identity': identity,
                        'confidence': similarity,
                        'source': 'internal_database',
                        'metadata': data.get('metadata', {})
                    })
        
        return matches
    
    def _load_encodings(self) -> Dict:
        """Load face encodings from database"""
        try:
            with open(self.db_path, 'r') as f:
                return json.load(f)
        except Exception as e:
            print(f"Error loading face encodings from {self.db_path}: {e}")
            return {}

class FacialSearchEngines:
    """Interface with external facial search engines"""
    
    async def search_public_databases(self, face_encoding: np.ndarray) -> List[Dict]:
        """Search public face databases (PimEyes, FaceCheck, etc.)"""
        matches = []
        
        # PimEyes-like search (would require API integration)
        try:
            pimeyes_results = await self._search_pimeyes(face_encoding)
            matches.extend(pimeyes_results)
        except Exception as e:
            print(f"PimEyes search error: {e}")
        
        # FaceCheck.id search
        try:
            facecheck_results = await self._search_facecheck(face_encoding)
            matches.extend(facecheck_results)
        except Exception as e:
            print(f"FaceCheck search error: {e}")
        
        return matches
    
    async def reverse_image_search(self, face_encoding: np.ndarray) -> List[Dict]:
        """Perform reverse image search on social media"""
        # Convert encoding back to image for reverse search
        temp_image_path = self._encoding_to_image(face_encoding)
        
        matches = []
        
        # Google Reverse Image Search
        try:
            google_results = await self._google_reverse_search(temp_image_path)
            matches.extend(google_results)
        except Exception as e:
            print(f"Google reverse search error: {e}")
        
        # Social media specific searches
        platforms = ['facebook', 'instagram', 'linkedin', 'twitter']
        for platform in platforms:
            try:
                platform_results = await self._platform_specific_search(temp_image_path, platform)
                matches.extend(platform_results)
            except Exception as e:
                print(f"{platform} search error: {e}")
        
        # Cleanup temp file
        os.unlink(temp_image_path)
        
        return matches

class FacialAnalysisPipeline:
    """Advanced facial analysis using multiple AI models"""
    
    def __init__(self):
        self.deepface_models = ['VGG-Face', 'Facenet', 'OpenFace', 'DeepFace']
        self.attribute_analyzer = FaceAttributeAnalyzer()
    
    async def analyze_face(self, image_path: str, face_data: Dict) -> Dict:
        """Comprehensive facial analysis"""
        analysis = {}
        
        # Basic face analysis
        analysis['basic'] = await self._basic_face_analysis(face_data)
        
        # Demographic analysis
        analysis['demographics'] = await self.attribute_analyzer.predict_demographics(image_path)
        
        # Emotional analysis
        analysis['emotions'] = await self.attribute_analyzer.analyze_emotions(image_path)
        
        # Unique facial features
        analysis['unique_features'] = await self._extract_unique_features(face_data)
        
        # Quality assessment
        analysis['quality'] = await self._assess_image_quality(image_path)
        
        return analysis
    
    async def _basic_face_analysis(self, face_data: Dict) -> Dict:
        """Basic face measurements and landmarks"""
        if not face_data.get('dlib_rects'):
            return {}
        
        rect = face_data['dlib_rects'][0]  # Primary face
        shape = self._get_facial_landmarks(face_data)
        
        return {
            'face_width': rect.width(),
            'face_height': rect.height(),
            'landmark_count': shape.num_parts if shape else 0,
            'facial_symmetry': self._calculate_symmetry(shape),
            'prominent_features': self._identify_prominent_features(shape)
        }

class FaceAttributeAnalyzer:
    """Analyze face attributes using DeepFace and custom models"""
    
    async def predict_demographics(self, image_path: str) -> Dict:
        """Predict age, gender, ethnicity"""
        try:
            analysis = DeepFace.analyze(
                img_path=image_path,
                actions=['age', 'gender', 'race'],
                enforce_detection=False
            )
            
            if analysis:
                primary_face = analysis[0]
                return {
                    'age': primary_face.get('age'),
                    'gender': primary_face.get('dominant_gender'),
                    'ethnicity': primary_face.get('dominant_race'),
                    'confidence': primary_face.get('face_confidence', 0)
                }
        except Exception as e:
            print(f"Demographic analysis error: {e}")
        
        return {}
    
    async def analyze_emotions(self, image_path: str) -> Dict:
        """Analyze facial emotions"""
        try:
            analysis = DeepFace.analyze(
                img_path=image_path,
                actions=['emotion'],
                enforce_detection=False
            )
            
            if analysis:
                emotions = analysis[0].get('emotion', {})
                return {
                    'dominant_emotion': max(emotions, key=emotions.get) if emotions else 'neutral',
                    'emotion_scores': emotions,
                    'emotional_intensity': max(emotions.values()) if emotions else 0
                }
        except Exception as e:
            print(f"Emotion analysis error: {e}")
        
        return {}

# MAIN FACIAL INTEL ORCHESTRATOR
class FacialIntelligenceOrchestrator:
    """Orchestrate all facial intelligence operations"""
    
    def __init__(self):
        self.facial_ai = HyperZillaFacialAI()
        self.osint_integration = OSINTIntegration()
    
    async def full_facial_intelligence_report(self, image_path: str, enhancement_data: Dict = None) -> Dict:
        """Generate complete facial intelligence report"""
        print("🧠 GENERATING COMPREHENSIVE FACIAL INTELLIGENCE REPORT...")
        
        # Step 1: Facial identification
        facial_id = await self.facial_ai.identify_person_from_image(image_path, enhance_with_osint=True)
        
        # Step 2: Enhanced OSINT with additional data
        if enhancement_data:
            enhanced_osint = await self.osint_integration.enhance_with_user_data(
                facial_id, 
                enhancement_data
            )
            facial_id['enhanced_osint'] = enhanced_osint
        
        # Step 3: Generate actionable intelligence
        facial_id['actionable_intelligence'] = await self._generate_actionable_insights(facial_id)
        
        return facial_id
    
    async def _generate_actionable_insights(self, facial_data: Dict) -> Dict:
        """Generate actionable intelligence from facial data"""
        insights = {
            'identification_confidence': facial_data.get('confidence_score', 0),
            'recommended_next_steps': [],
            'potential_risks': [],
            'verification_suggestions': []
        }
        
        # Analyze matches and suggest actions
        matches = facial_data.get('database_matches', [])
        if matches:
            best_match = matches[0]
            if best_match['confidence'] > 0.8:
                insights['recommended_next_steps'].append(f"High confidence match: Verify identity {best_match['identity']}")
            else:
                insights['recommended_next_steps'].append("Perform enhanced OSINT with demographic filters")
        
        # Suggest verification methods
        analysis = facial_data.get('primary_face_analysis', {})
        if analysis.get('demographics'):
            demographics = analysis['demographics']
            insights['verification_suggestions'].extend(
                self._generate_verification_suggestions(demographics)
            )
        
        return insights

# GLOBAL FACIAL INTEL INSTANCE
FACIAL_INTEL = FacialIntelligenceOrchestrator()
