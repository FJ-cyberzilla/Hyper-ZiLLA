#!/bin/bash
# Create Debian package

set -e

cd "$(dirname "$0")/.."  # CD to project root

VERSION=$(grep -m1 version Cargo.toml | head -1 | cut -d '"' -f2)
ARCH=$(uname -m)
PACKAGE_NAME="entynet-pentest_${VERSION}_${ARCH}.deb"
BUILD_DIR="deb_build"

echo "Creating Debian package: $PACKAGE_NAME"

# Build release first
echo "Building release binary..."
cargo build --release

# Create build directory
mkdir -p $BUILD_DIR/DEBIAN
mkdir -p $BUILD_DIR/usr/local/bin
mkdir -p $BUILD_DIR/etc/entynet
mkdir -p $BUILD_DIR/var/lib/entynet
mkdir -p $BUILD_DIR/var/log/entynet

# Copy binary
cp target/release/entynet-pentest $BUILD_DIR/usr/local/bin/

# Create control file
cat > $BUILD_DIR/DEBIAN/control << EOF
Package: entynet-pentest
Version: $VERSION
Section: security
Priority: optional
Architecture: $ARCH
Depends: libssl3, libpcap0.8
Maintainer: Entynetproject <entynetproject@example.com>
Description: Enterprise Level Penetration Testing Framework
 Entynet Hacker Tools Enterprise is a comprehensive penetration
 testing framework for security professionals.
EOF

# Create postinst script
cat > $BUILD_DIR/DEBIAN/postinst << EOF
#!/bin/bash
chmod 755 /usr/local/bin/entynet-pentest
EOF
chmod 755 $BUILD_DIR/DEBIAN/postinst

# Build package
dpkg-deb --build $BUILD_DIR dist/$PACKAGE_NAME

# Cleanup
rm -rf $BUILD_DIR

echo "Debian package created: dist/$PACKAGE_NAME"
