Name: entynet-enterprise
Version: 
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
* Sat Nov 15 2025 Entynetproject <entynetproject@example.com> - -1
- Initial package
