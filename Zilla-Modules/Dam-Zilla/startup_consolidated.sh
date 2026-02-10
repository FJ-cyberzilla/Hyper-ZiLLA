#!/bin/bash

set -e  # Exit on any error

echo "🐉 ZILLA-DAM CONSOLIDATED STARTUP SEQUENCE"
echo "══════════════════════════════════════════"

# Function to check dependencies
check_dependencies() {
    echo "🔍 Checking dependencies..."
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js not found. Please install Node.js 18 or higher."
        exit 1
    fi
    
    # Check Python
    if ! command -v python3 &> /dev/null; then
        echo "❌ Python3 not found. Please install Python 3.8 or higher."
        exit 1
    fi
    
    # Check Docker (optional)
    if command -v docker &> /dev/null; then
        echo "✅ Docker available"
    else
        echo "⚠️  Docker not available - some features may be limited"
    fi
    
    echo "✅ All critical dependencies found"
}

# Function to setup environment
setup_environment() {
    echo "🔧 Setting up environment..."
    
    # Create necessary directories
    mkdir -p logs data temp operations encrypted_vault/encrypted_data
    
    # Set permissions
    chmod 750 logs data temp operations
    chmod 700 encrypted_vault/encrypted_data
    
    # Copy environment file if not exists
    if [ ! -f .env ]; then
        cp .env.example .env
        echo "⚠️  Please configure .env file before proceeding"
        exit 1
    fi
    
    # Load environment variables
    set -a
    source .env
    set +a
}

# Function to install dependencies
install_dependencies() {
    echo "📦 Installing dependencies..."
    
    # Node.js dependencies
    if [ ! -d "node_modules" ]; then
        npm install
    else
        echo "✅ Node.js dependencies already installed"
    fi
    
    # Python dependencies
    if [ ! -d "venv" ]; then
        python3 -m venv venv
        source venv/bin/activate
        pip install -r requirements.txt
    else
        source venv/bin/activate
        echo "✅ Python dependencies already installed"
    fi
}

# Function to initialize database
initialize_database() {
    echo "🗄️ Initializing database..."
    
    if [ "$USE_POSTGRES" = "true" ] && [ -n "$POSTGRES_URL" ]; then
        echo "🔗 Using PostgreSQL"
        node core/database/migrations/postgres.js
    else
        echo "💾 Using SQLite (in-memory)"
        node core/database/migrations/sqlite.js
    fi
}

# Function to start services
start_services() {
    echo "🚀 Starting services..."
    
    # Start main application with optimized settings
    exec node --max-old-space-size=4096 \
              --max-semi-space-size=128 \
              --optimize-for-size \
              core/orchestration/zilla_orchestrator.js "$@"
}

# Main execution flow
main() {
    check_dependencies
    setup_environment
    install_dependencies
    initialize_database
    start_services "$@"
}

# Run main function with all arguments
main "$@"
