#!/bin/bash
# Create RPM package (basic version)

set -e

VERSION=$(grep version Cargo.toml | head -1 | cut -d '"' -f2)
ARCH=$(uname -m)
RPM_DIR="rpm_build"

echo "Creating RPM structure..."

mkdir -p $RPM_DIR/BUILD $RPM_DIR/RPMS $RPM_DIR/SOURCES $RPM_DIR/SPECS $RPM_DIR/SRPMS

# Create spec file
cat > $RPM_DIR/SPECS/entynet-enterprise.spec << EOF
Name: entynet-enterprise
Version: $VERSION
Release: 1%{?dist}
Summary: Enterprise Level Penetration Testing Framework
License: MIT
URL: https://entrynetproject.simplesite.com
Source0: entynet-enterprise

%description
Entynet Hacker Tools Enterprise is a comprehensive penetration
testing framework for security professionals.

%prep

%build

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 755 %{SOURCE0} %{buildroot}/usr/local/bin/entynet-enterprise

%files
/usr/local/bin/entynet-enterprise

%changelog
* $(date +"%a %b %d %Y") Entynetproject <entynetproject@example.com> - $VERSION-1
- Initial package
EOF

# Copy binary to sources
cp target/release/entynet-enterprise $RPM_DIR/SOURCES/

# Build RPM (simplified - in real scenario use rpmbuild)
echo "RPM creation would require rpmbuild tool"
echo "RPM spec file created at: $RPM_DIR/SPECS/entynet-enterprise.spec"
