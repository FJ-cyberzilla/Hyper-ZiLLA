#!/bin/bash
# scripts/install-wsl.sh

echo "🪟 COGNIZILLA WSL/WSL2 INSTALLATION"
echo "🔐 FJ-Cyberzilla Sovereign System"
echo "=================================="

# Check if we're in WSL
if ! grep -q Microsoft /proc/version; then
    echo "❌ This script must be run in WSL/WSL2"
    exit 1
fi

# Update system
echo "🔄 Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install dependencies
echo "📦 Installing dependencies..."
sudo apt install -y \
    git \
    curl \
    wget \
    build-essential \
    libssl-dev \
    pkg-config \
    docker.io \
    docker-compose

# Install Go
if ! command -v go &> /dev/null; then
    echo "🐹 Installing Go..."
    wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
    source ~/.bashrc
fi

# Install Node.js
if ! command -v node &> /dev/null; then
    echo "📟 Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
fi

# Start Docker
echo "🐳 Starting Docker..."
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Clone Cognizilla
echo "📥 Cloning Cognizilla..."
git clone https://github.com/FJ-cyberzilla/cognizilla.git
cd cognizilla

# Build and install
echo "🏗️ Building Cognizilla..."
make install-wsl

echo ""
echo "🎉 COGNIZILLA INSTALLATION COMPLETE!"
echo "===================================="
echo "🦖 System: Cognizilla Quantum"
echo "👑 Owner: FJ-Cyberzilla"
echo "🐧 Platform: WSL/WSL2"
echo ""
echo "🚀 Quick Start:"
echo "   cognizilla --quantum --sovereign"
echo ""
echo "🌐 Dashboard: https://localhost:8443"
echo "🔧 API: http://localhost:8080"
