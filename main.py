#!/usr/bin/env python3
"""
Hyper-ZiLLA Main Entry Point
Proprietary AI Security Intelligence Platform
"""

import argparse
import sys
import os
import logging
from pathlib import Path
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Add project root to Python path
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

from HyperZilla.COMMAND_CENTER.app import create_app
from HyperZilla.OPERATIONS_ARM.activation_test import SystemActivationTest
from HyperZilla.ZILLA_CORE.ai_hierarchy.director_ai import DirectorAI

def show_banner():
    """Display Hyper-ZiLLA proprietary AI banner"""
    banner = """
    ╔══════════════════════════════════════════════════════════════╗
    ║                   HYPER-ZiLLA AI PLATFORM                    ║
    ║           Proprietary Artificial Intelligence System         ║
    ║          🔒 100% Custom AI • Zero External Dependencies 🔒   ║
    ╚══════════════════════════════════════════════════════════════╝
    """
    print(banner)

def setup_logging():
    """Configure application logging"""
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        handlers=[
            logging.FileHandler('hyperzilla.log'),
            logging.StreamHandler(sys.stdout)
        ]
    )

def main():
    """Main application entry point"""
    parser = argparse.ArgumentParser(
        description='Hyper-ZiLLA - Proprietary AI Security Intelligence Platform',
        epilog='All AI technology is 100% proprietary Hyper-ZiLLA implementation'
    )
    parser.add_argument('--mode', choices=['web', 'cli', 'test', 'api'], 
                       default='web', help='Operation mode')
    parser.add_argument('--port', type=int, default=5000, 
                       help='Web server port (web mode only)')
    parser.add_argument('--host', type=str, default='127.0.0.1',
                       help='Web server host (web mode only)')
    parser.add_argument('--debug', action='store_true',
                       help='Enable debug mode')
    
    args = parser.parse_args()
    
    # Show proprietary banner
    show_banner()
    
    # Setup logging
    setup_logging()
    logger = logging.getLogger(__name__)
    
    logger.info(f"Starting Hyper-ZiLLA Proprietary AI Platform in {args.mode} mode")
    logger.info("All AI systems: 100% Custom Hyper-ZiLLA Technology")
    
    # If HZ_DATABASE_URL is set to the default PostgreSQL example, remove it
    # to force SQLite for development unless a custom HZ_DATABASE_URL is explicitly provided
    if os.environ.get('HZ_DATABASE_URL') == 'postgresql://hyperzilla:hyperzilla-password@localhost/hyperzilla':
        del os.environ['HZ_DATABASE_URL']

    # Ensure HZ_DATABASE_URL is set for development if not already present
    if not os.environ.get('HZ_DATABASE_URL'):
        from config.settings import BaseConfig
        os.environ['HZ_DATABASE_URL'] = f'sqlite:///{BaseConfig.DATA_DIR}/hyperzilla_dev.db'
    
    try:
        if args.mode == 'web':
            logger.info(f"Starting proprietary AI web interface on {args.host}:{args.port}")
            app = create_app()
            app.run(
                host=args.host,
                port=args.port,
                debug=args.debug
            )
            
        elif args.mode == 'test':
            logger.info("Running proprietary AI system activation test")
            test = SystemActivationTest()
            test.run_comprehensive_test()
            
        elif args.mode == 'cli':
            logger.info("Starting proprietary AI CLI mode")
            director = DirectorAI()
            director.initialize_system()
            director.start_monitoring()
            
        else:
            logger.info("Proprietary AI API mode selected")
            app = create_app()
            app.run(host=args.host, port=args.port, debug=args.debug)
            
    except KeyboardInterrupt:
        logger.info("Proprietary AI system shutdown requested by user")
    except Exception as e:
        logger.error(f"Proprietary AI system error: {e}")
        if args.debug:
            raise
        sys.exit(1)

if __name__ == '__main__':
    main()
