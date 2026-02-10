#!/bin/bash

echo "🚀 Hyper-ZiLLA Quick Start"
echo "=========================="

# Check if Python is available
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is required but not installed."
    exit 1
fi

# Create virtual environment
echo "📦 Setting up virtual environment..."
python3 -m venv venv
source venv/bin/activate

# Install dependencies
echo "📋 Installing dependencies..."
pip install --upgrade pip
pip install -r requirements.txt

# Create environment file
if [ ! -f .env ]; then
    echo "⚙️ Creating environment configuration..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your configuration"
fi

# Create necessary directories
mkdir -p logs data cache

# Run basic tests
echo "🧪 Running basic tests..."
python -m pytest tests/test_basic.py -v

echo ""
echo "✅ Setup complete!"
echo ""
echo "To start Hyper-ZiLLA:"
echo "  source venv/bin/activate"
echo "  python main.py --mode web"
echo ""
echo "Or use the quick commands:"
echo "  ./main.py --mode web"
echo ""
