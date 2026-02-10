#!/bin/bash
# scripts/install-macos.sh

echo "🍎 COGNIZILLA macOS INSTALLATION"
echo "🔐 FJ-Cyberzilla Sovereign System"
echo "================================"

# Check for Homebrew
if ! command -v brew &> /dev/null; then
    echo "🍺 Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Install dependencies
echo "📦 Installing dependencies..."
brew install \
    go \
    node \
    docker \
    docker-compose

# Start Docker
echo "🐳 Starting Docker..."
open -a Docker

# Wait for Docker to start
sleep 30

# Clone Cognizilla
echo "📥 Cloning Cognizilla..."
git clone https://github.com/FJ-cyberzilla/cognizilla.git
cd cognizilla

# Build and install
echo "🏗️ Building Cognizilla..."
make install-macos

echo ""
echo "🎉 COGNIZILLA INSTALLATION COMPLETE!"
echo "===================================="
echo "🦖 System: Cognizilla Quantum"
echo "👑 Owner: FJ-Cyberzilla"
echo "🍎 Platform: macOS"
echo ""
echo "🚀 Quick Start:"
echo "   cognizilla --quantum --sovereign"
echo ""
echo "🌐 Dashboard: https://localhost:8443"
echo "🔧 API: http://localhost:8080"
