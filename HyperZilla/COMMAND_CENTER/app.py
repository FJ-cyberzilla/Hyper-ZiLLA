"""
Hyper-ZiLLA Command Center Web Interface
Proprietary AI Monitoring and Control
"""

from flask import Flask, render_template, jsonify, request, g
import sys
import os
from pathlib import Path

# Add project root to Python path
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))


def get_ai_systems():
    """
    Initializes and returns the AI systems.
    Uses Flask's application context 'g' to store the systems.
    """
    if 'director_ai' not in g:
        try:
            from HyperZilla.ZILLA_CORE.ai_hierarchy.director_ai import DirectorAI
            from HyperZilla.OPERATIONS_ARM.activation_test import SystemActivationTest
            g.director_ai = DirectorAI()
            g.system_test = SystemActivationTest()
            g.ai_available = True
        except ImportError as e:
            print(f"Warning: AI modules not available: {e}")
            g.director_ai = None
            g.system_test = None
            g.ai_available = False
    return g.director_ai, g.system_test, g.ai_available


def create_app():
    """
    Application factory for the Flask app.
    """
    app = Flask(__name__)

    from config.settings import config_by_name
    env = os.environ.get('FLASK_ENV', 'development')
    app.config.from_object(config_by_name[env])

    from HyperZilla.database import db
    from HyperZilla.models import Event
    db.init_app(app)

    with app.app_context():
        db.create_all()

    with app.app_context():
        director_ai, _, ai_available = get_ai_systems()
        if director_ai and ai_available:
            director_ai.initialize_system()

    @app.route('/')
    def dashboard():
        """Main dashboard page"""
        _, _, ai_available = get_ai_systems()
        return render_template(
            'dashboard.html',
            ai_available=ai_available,
            system_name="Hyper-ZiLLA Proprietary AI"
        )

    @app.route('/facial-intel')
    def facial_intel_panel():
        """Facial intelligence panel"""
        _, _, ai_available = get_ai_systems()
        return render_template(
            'facial_intel_panel.html',
            ai_available=ai_available
        )

    @app.route('/api/system-status')
    def system_status():
        """API endpoint for system status"""
        director_ai, _, ai_available = get_ai_systems()
        if not ai_available:
            return jsonify({
                "status": "ai_unavailable",
                "message": "AI systems not initialized",
                "modules": []
            })

        try:
            status = director_ai.get_system_status() if director_ai else {}
            return jsonify(
                            {
                                "status": "operational",
                                "system_status": status.get('system_status', 'unknown'),
                                "modules": status.get('active_modules', []),
                                "uptime": status.get('uptime', '0:00:00'),
                                "ai_technology": "Hyper-ZiLLA Proprietary AI"
                            }
                        )
        except Exception as e:
            return jsonify({
                "status": "error",
                "error": str(e)
            })

    @app.route('/api/run-test', methods=['POST'])
    def run_system_test():
        """API endpoint to run system tests"""
        _, system_test, ai_available = get_ai_systems()
        if not ai_available:
            return jsonify({
                "status": "error",
                "message": "AI systems not available"
            })

        try:
            system_test.run_comprehensive_test()
            return jsonify({
                "status": "success",
                "message": "System tests completed"
            })
        except Exception as e:
            return jsonify({
                "status": "error",
                "error": str(e)
            })

    @app.route('/api/initialize-ai', methods=['POST'])
    def initialize_ai():
        """API endpoint to initialize AI systems"""
        director_ai, _, ai_available = get_ai_systems()

        try:
            if not ai_available:
                # Re-initialize within the context
                from HyperZilla.ZILLA_CORE.ai_hierarchy.director_ai import DirectorAI
                g.director_ai = DirectorAI()
                g.ai_available = True

            g.director_ai.initialize_system()

            return jsonify({
                "status": "success",
                "message": "Hyper-ZiLLA AI initialized successfully"
            })
        except Exception as e:
            return jsonify({
                "status": "error",
                "error": str(e)
            })

    @app.route('/api/ai-capabilities')
    def ai_capabilities():
        """API endpoint listing AI capabilities"""
        capabilities = {
            "proprietary_ai": True,
            "custom_neural_networks": True,
            "facial_recognition": True,
            "threat_intelligence": True,
            "pattern_analysis": True,
            "strategic_decision_making": True,
            "external_dependencies": False,
            "technology_owner": "Hyper-ZiLLA Team"
        }

        return jsonify(capabilities)

    @app.route('/api/events')
    def get_events():
        """API endpoint for system events"""
        try:
            events = Event.query.order_by(Event.timestamp.desc()).all()
            return jsonify([{
                'id': event.id,
                'type': event.type,
                'description': event.description,
                'timestamp': event.timestamp.isoformat()
            } for event in events])
        except Exception as e:
            return jsonify({
                "status": "error",
                "error": str(e)
            })

    return app