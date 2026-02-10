#!/bin/bash
echo "Testing Entynet Enterprise build..."
cargo check
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "You can now run: cargo run"
else
    echo "❌ Build failed. Check errors above."
fi
