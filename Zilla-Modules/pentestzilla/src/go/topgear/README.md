```markdown

# Zillabreah ⚡🦖
## Mobile Security Laboratory 

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go).
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)

> Professional security
testing platform
for authorized
research and
educational purposes only.

## 🚀 Features

- **Real-time Network Scanning** - Active port scanning and service detection
- **Vulnerability Assessment** - Risk-based security analysis
- **Professional Reporting** - Color-coded terminal interface
- **Cross-Platform** - Runs on Linux, Windows, macOS, and Android/Termux
- **Educational Focus** - Designed for security research and learning

## 📦 Installation

### Pre-built Binaries
Download the latest release from [GitHub Releases](https://github.com/FJ-cyberzilla/topgear/releases)

### From Source
```bash
git clone https://github.com/FJ-cyberzilla/topgear.git
cd topgear
go build -o topgear
```

Termux (Android)

```bash
pkg install golang git
git clone https://github.com/FJ-cyberzilla/topgear.git
cd topgear
go build -o topgear
```

🎯 Quick Start

```bash
# Basic scan
./topgear --target localhost

# Scan specific host
./topgear --target 192.168.1.1

# Scan domain
./topgear --target example.com
```

📋 Usage

```
Usage:
  topgear [flags]

Flags:
  -t, --target string   Target to scan (IP, hostname, or domain)
  -o, --output string   Output format (json, html, text)
  -h, --help           help for topgear
```

🛡️ Legal & Ethics

⚠️ IMPORTANT: TOPGEAR is designed for:

· Security research and education
· Testing your own systems and devices
· Authorized penetration testing

🚫 DO NOT USE for:

· Unauthorized scanning of systems you don't own
· Malicious activities
· Any illegal purposes

🏗️ Architecture

TOPGEAR uses a modular architecture:

· Core Engine - Orchestrates scanning processes
· Scanner Modules - Network, port, and service detection
· Analyzer Engine - Vulnerability assessment and risk classification
· UI Layer - Professional terminal interface
