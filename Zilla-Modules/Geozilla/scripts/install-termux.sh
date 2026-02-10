#!/bin/bash
# scripts/install-termux.sh

echo "📱 COGNIZILLA TERMUX INSTALLATION"
echo "🔐 FJ-Cyberzilla Mobile Sovereign System"
echo "========================================"

# Update Termux
pkg update && pkg upgrade -y

# Install dependencies
pkg install -y \
    git \
    golang \
    nodejs \
    openssl-tool \
    termux-api

# Setup Go environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Clone Cognizilla
git clone https://github.com/FJ-cyberzilla/cognizilla.git
cd cognizilla

# Build for Termux
make install-termux

# Setup storage access
termux-setup-storage

echo ""
echo "🎉 COGNIZILLA MOBILE READY!"
echo "=========================="
echo "📱 Platform: Termux (Android)"
echo "👑 Owner: FJ-Cyberzilla"
echo "🦖 Quantum: Mobile-Optimized"
echo ""
echo "🚀 Start with: cognizilla --mobile --quantum"
echo ""
echo "💡 Enable 'Stay awake' in Termux settings for background operation"
