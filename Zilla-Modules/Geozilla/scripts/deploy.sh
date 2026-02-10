#!/bin/bash
# scripts/deploy.sh

echo "🦖 DEPLOYING COGNIZILLA SOVEREIGN SYSTEM"
echo "🔐 EXCLUSIVE TO: FJ-CYBERZILLA"
echo "========================================"

# Check prerequisites
check_prerequisites() {
    echo "🔍 Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker not found. Please install Docker."
        exit 1
    fi
    
    if ! command -v git &> /dev/null; then
        echo "❌ Git not found. Please install Git."
        exit 1
    fi
    
    echo "✅ Prerequisites satisfied"
}

# Initialize quantum security
init_quantum_security() {
    echo "🔐 Initializing quantum security..."
    
    # Generate FJ-Cyberzilla exclusive keys
    mkdir -p configuration/quantum-keys
    openssl genrsa -out configuration/quantum-keys/fj-cyberzilla.key 4096
    openssl rsa -in configuration/quantum-keys/fj-cyberzilla.key -pubout -out configuration/quantum-keys/fj-cyberzilla.pub
    
    # Set secure permissions
    chmod 600 configuration/quantum-keys/*
    
    echo "✅ Quantum security initialized"
}

# Build and deploy
deploy_cognizilla() {
    echo "🏗️ Building Cognizilla Monster..."
    
    # Build Docker image
    docker build -t fjcyberzilla/cognizilla:latest .
    
    echo "🚀 Starting Cognizilla Services..."
    
    # Deploy with Docker Compose
    docker-compose up -d
    
    # Wait for services to be ready
    sleep 10
    
    # Health check
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Cognizilla is alive and quantum-entangled!"
    else
        echo "❌ Cognizilla failed to start"
        exit 1
    fi
}

# Display access information
show_access_info() {
    echo ""
    echo "🎉 COGNIZILLA DEPLOYMENT COMPLETE!"
    echo "=================================="
    echo "🦖 System: COGNIZILLA QUANTUM"
    echo "👑 Owner: FJ-CYBERZILLA"
    echo "🌐 Dashboard: https://localhost:8443"
    echo "🔧 API: http://localhost:8080"
    echo ""
    echo "🔐 Quantum Identity:"
    docker logs cognizilla-monster 2>&1 | grep "Quantum Identity" | tail -1
    echo ""
    echo "🚀 Next steps:"
    echo "  1. Access the dashboard at https://localhost:8443"
    echo "  2. Activate the Pixel Monster"
    echo "  3. Monitor conscious agent activity"
    echo ""
}

main() {
    check_prerequisites
    init_quantum_security
    deploy_cognizilla
    show_access_info
}

main "$@"
