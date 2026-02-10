#!/bin/bash
# Create distribution tarball

set -e

VERSION=$(grep version Cargo.toml | head -1 | cut -d '"' -f2)
ARCH=$(uname -m)
TARBALL_NAME="entynet-enterprise-${VERSION}-${ARCH}.tar.gz"
DIST_DIR="entynet-enterprise-${VERSION}"

echo "Creating tarball: $TARBALL_NAME"

# Create distribution directory
mkdir -p $DIST_DIR
mkdir -p $DIST_DIR/config
mkdir -p $DIST_DIR/logs

# Copy binary and files
cp target/release/entynet-enterprise $DIST_DIR/
cp README.md $DIST_DIR/
cp LICENSE $DIST_DIR/
cp -r examples $DIST_DIR/ 2>/dev/null || true

# Create install script
cat > $DIST_DIR/install.sh << 'EOF'
#!/bin/bash
# Installation script for Entynet Enterprise

set -e

echo "Installing Entynet Enterprise..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root for system-wide installation"
    exit 1
fi

# Copy binary
cp entynet-enterprise /usr/local/bin/
chmod 755 /usr/local/bin/entynet-enterprise

# Create directories
mkdir -p /etc/entynet
mkdir -p /var/lib/entynet
mkdir -p /var/log/entynet

echo "Installation completed!"
echo "Run: entynet-enterprise"
EOF
chmod +x $DIST_DIR/install.sh

# Create tarball
tar czf dist/$TARBALL_NAME $DIST_DIR

# Cleanup
rm -rf $DIST_DIR

echo "Tarball created: dist/$TARBALL_NAME"
