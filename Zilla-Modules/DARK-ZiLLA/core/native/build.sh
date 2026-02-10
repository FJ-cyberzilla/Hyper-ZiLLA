#!/bin/bash
#
# CodeZilla Build Script
# Complete setup, build, and installation script for CodeZilla
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
RESET='\033[0m'

# Project information
PROJECT_NAME="CodeZilla"
VERSION="2.0.0"
BUILD_DIR="build"
INSTALL_PREFIX="/usr/local"

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Functions
print_header() {
    echo -e "${CYAN}${BOLD}"
    echo "╔════════════════════════════════════════════════════════════════════╗"
    echo "║                     CodeZilla Build System                         ║"
    echo "║                    Version $VERSION - Enterprise Edition                    ║"
    echo "╚════════════════════════════════════════════════════════════════════╝"
    echo -e "${RESET}"
}

print_step() {
    echo -e "${BLUE}${BOLD}[*]${RESET} $1"
}

print_success() {
    echo -e "${GREEN}${BOLD}[✓]${RESET} $1"
}

print_error() {
    echo -e "${RED}${BOLD}[✗]${RESET} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}[!]${RESET} $1"
}

print_info() {
    echo -e "${CYAN}[i]${RESET} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check system requirements
check_requirements() {
    print_step "Checking system requirements..."
    
    local missing_deps=()
    
    # Check for required commands
    if ! command_exists cmake; then
        missing_deps+=("cmake")
    else
        print_success "CMake found: $(cmake --version | head -n1)"
    fi
    
    if ! command_exists g++; then
        missing_deps+=("g++")
    else
        print_success "G++ found: $(g++ --version | head -n1)"
    fi
    
    if ! command_exists python3; then
        missing_deps+=("python3")
    else
        print_success "Python3 found: $(python3 --version)"
    fi
    
    if ! command_exists pkg-config; then
        missing_deps+=("pkg-config")
    else
        print_success "pkg-config found"
    fi
    
    # Check for required libraries
    if ! pkg-config --exists sqlite3; then
        missing_deps+=("libsqlite3-dev")
    else
        print_success "SQLite3 found: $(pkg-config --modversion sqlite3)"
    fi
    
    if ! pkg-config --exists openssl; then
        missing_deps+=("libssl-dev")
    else
        print_success "OpenSSL found: $(pkg-config --modversion openssl)"
    fi
    
    # Check for pthread
    if ! ldconfig -p | grep -q libpthread; then
        missing_deps+=("libpthread (usually included in libc6-dev)")
    else
        print_success "pthread library found"
    fi
    
    # Report missing dependencies
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_error "Missing dependencies detected!"
        echo -e "${YELLOW}Missing packages:${RESET}"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo ""
        echo -e "${CYAN}Install them with:${RESET}"
        if command_exists apt-get; then
            echo "  sudo apt-get update"
            echo "  sudo apt-get install -y build-essential cmake libsqlite3-dev libssl-dev python3 pkg-config"
        elif command_exists yum; then
            echo "  sudo yum install -y gcc-c++ cmake sqlite-devel openssl-devel python3 pkgconfig"
        elif command_exists pacman; then
            echo "  sudo pacman -S base-devel cmake sqlite openssl python"
        elif command_exists brew; then
            echo "  brew install cmake sqlite openssl python3"
        else
            echo "  Please install the missing packages using your system's package manager"
        fi
        return 1
    fi
    
    print_success "All system requirements satisfied"
    echo ""
    return 0
}

# Verify source files exist
verify_source_files() {
    print_step "Verifying source files..."
    
    local required_files=(
        "main.cpp"
        "CMakeLists.txt"
        "src/analysis/ai/ai_engine.cpp"
        "src/analysis/ai/ai_engine.h"
        "src/analysis/ai/ai_service.py"
        "src/ui/menu_system.cpp"
        "src/core/configuration_manager.cpp"
        "src/core/database_manager.cpp"
        "src/utils/logger.cpp"
    )
    
    local missing_files=()
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ]; then
            missing_files+=("$file")
        fi
    done
    
    if [ ${#missing_files[@]} -gt 0 ]; then
        print_error "Missing source files detected!"
        echo -e "${YELLOW}Missing files:${RESET}"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
        return 1
    fi
    
    print_success "All required source files present"
    
    # Make Python script executable
    if [ -f "src/analysis/ai/ai_service.py" ]; then
        chmod +x src/analysis/ai/ai_service.py
        print_success "AI service script made executable"
    fi
    
    echo ""
    return 0
}

# Create configuration file if it doesn't exist
create_config() {
    print_step "Checking configuration file..."
    
    if [ ! -f "config.json" ]; then
        print_warning "config.json not found, creating default configuration..."
        
        cat > config.json << 'EOF'
{
  "application": {
    "name": "CodeZilla",
    "version": "2.0.0",
    "environment": "production"
  },
  "database": {
    "path": "codezilla.db",
    "connection_timeout": 30,
    "max_connections": 10,
    "auto_vacuum": true,
    "journal_mode": "WAL"
  },
  "logging": {
    "level": "INFO",
    "file": "codezilla.log",
    "max_file_size_mb": 50,
    "max_backup_files": 5,
    "console_output": true,
    "timestamp_format": "%Y-%m-%d %H:%M:%S"
  },
  "ai_engine": {
    "python_executable": "python3",
    "service_path": "src/analysis/ai/ai_service.py",
    "model_type": "advanced",
    "timeout_seconds": 30,
    "max_retries": 3,
    "retry_delay_ms": 100,
    "enable_caching": true,
    "enable_learning": true,
    "cache_max_size": 1000,
    "cache_ttl_seconds": 3600,
    "warmup_on_startup": false,
    "health_check_interval_seconds": 300
  },
  "analysis": {
    "supported_languages": ["cpp", "python", "javascript", "java", "go"],
    "max_file_size_mb": 10,
    "recursive_depth_limit": 10,
    "exclude_patterns": ["*.git*", "*.svn*", "*node_modules*", "*__pycache__*", "*.min.js", "*.min.css", "*build*", "*dist*"],
    "include_patterns": ["*.cpp", "*.cc", "*.cxx", "*.h", "*.hpp", "*.py", "*.js", "*.java", "*.go"],
    "parallel_analysis": true,
    "max_parallel_jobs": 4
  },
  "security": {
    "enable_vulnerability_detection": true,
    "severity_threshold": "LOW",
    "enable_cwe_mapping": true,
    "enable_cvss_scoring": false
  },
  "reporting": {
    "output_directory": "reports",
    "default_format": "json",
    "include_metrics": true,
    "include_recommendations": true,
    "timestamp_reports": true
  },
  "performance": {
    "enable_profiling": false,
    "memory_limit_mb": 2048,
    "cpu_limit_percent": 80
  },
  "ui": {
    "enable_colors": true,
    "clear_screen_on_menu": true,
    "auto_scroll": true,
    "items_per_page": 20
  },
  "paths": {
    "temp_directory": "/tmp/codezilla",
    "cache_directory": ".cache/codezilla",
    "report_directory": "reports"
  }
}
EOF
        print_success "Created default config.json"
    else
        print_success "config.json already exists"
    fi
    echo ""
}

# Clean build directory
clean_build() {
    print_step "Cleaning build directory..."
    
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        print_success "Removed old build directory"
    fi
    
    mkdir -p "$BUILD_DIR"
    print_success "Created fresh build directory"
    echo ""
}

# Configure with CMake
configure_build() {
    print_step "Configuring build with CMake..."
    
    cd "$BUILD_DIR"
    
    local cmake_args=()
    
    # Add build type
    if [ "$BUILD_TYPE" = "Debug" ]; then
        cmake_args+=("-DCMAKE_BUILD_TYPE=Debug")
        print_info "Build type: Debug"
    else
        cmake_args+=("-DCMAKE_BUILD_TYPE=Release")
        print_info "Build type: Release"
    fi
    
    # Add install prefix if provided
    if [ -n "$INSTALL_PREFIX" ]; then
        cmake_args+=("-DCMAKE_INSTALL_PREFIX=$INSTALL_PREFIX")
        print_info "Install prefix: $INSTALL_PREFIX"
    fi
    
    # Run CMake
    if cmake "${cmake_args[@]}" ..; then
        print_success "CMake configuration successful"
    else
        print_error "CMake configuration failed"
        cd ..
        return 1
    fi
    
    cd ..
    echo ""
    return 0
}

# Build the project
build_project() {
    print_step "Building $PROJECT_NAME..."
    
    cd "$BUILD_DIR"
    
    # Determine number of parallel jobs
    if command_exists nproc; then
        JOBS=$(nproc)
    elif command_exists sysctl; then
        JOBS=$(sysctl -n hw.ncpu)
    else
        JOBS=2
    fi
    
    print_info "Using $JOBS parallel jobs"
    
    # Build
    if make -j"$JOBS"; then
        print_success "Build completed successfully"
    else
        print_error "Build failed"
        cd ..
        return 1
    fi
    
    cd ..
    echo ""
    return 0
}

# Run tests
run_tests() {
    print_step "Running tests..."
    
    cd "$BUILD_DIR"
    
    if [ -f "system_test" ]; then
        print_info "Running system tests..."
        if ./system_test; then
            print_success "All tests passed"
        else
            print_warning "Some tests failed"
        fi
    else
        print_warning "Test executable not found, skipping tests"
    fi
    
    cd ..
    echo ""
}

# Install the project
install_project() {
    print_step "Installing $PROJECT_NAME..."
    
    cd "$BUILD_DIR"
    
    if [ "$EUID" -ne 0 ] && [ "$INSTALL_PREFIX" = "/usr/local" ]; then
        print_warning "Root privileges required for system-wide installation"
        print_info "Run: sudo make install"
        print_info "Or use: ./build.sh --install-prefix=\$HOME/.local"
    else
        if make install; then
            print_success "Installation completed"
            print_info "Binary installed to: $INSTALL_PREFIX/bin/codezilla"
        else
            print_error "Installation failed"
            cd ..
            return 1
        fi
    fi
    
    cd ..
    echo ""
}

# Print build summary
print_summary() {
    echo -e "${CYAN}${BOLD}"
    echo "╔════════════════════════════════════════════════════════════════════╗"
    echo "║                        Build Summary                               ║"
    echo "╚════════════════════════════════════════════════════════════════════╝"
    echo -e "${RESET}"
    
    if [ -f "$BUILD_DIR/codezilla" ]; then
        echo -e "${GREEN}✓ Executable:${RESET} $BUILD_DIR/codezilla"
        
        # Get file size
        if command_exists du; then
            SIZE=$(du -h "$BUILD_DIR/codezilla" | cut -f1)
            echo -e "${GREEN}✓ Size:${RESET} $SIZE"
        fi
        
        echo ""
        echo -e "${CYAN}Run the application:${RESET}"
        echo "  cd $BUILD_DIR && ./codezilla"
        echo ""
        echo -e "${CYAN}Or analyze a file directly:${RESET}"
        echo "  cd $BUILD_DIR && ./codezilla --analyze ../test_files/vulnerable.cpp"
        echo ""
        echo -e "${CYAN}Show help:${RESET}"
        echo "  cd $BUILD_DIR && ./codezilla --help"
    else
        print_error "Build failed - executable not found"
    fi
    
    echo ""
}

# Usage information
print_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -c, --clean             Clean build directory before building"
    echo "  -d, --debug             Build in debug mode"
    echo "  -r, --release           Build in release mode (default)"
    echo "  -t, --test              Run tests after building"
    echo "  -i, --install           Install after building"
    echo "  --install-prefix=PATH   Set installation prefix (default: /usr/local)"
    echo "  --skip-checks           Skip dependency checks"
    echo "  -j, --jobs=N            Number of parallel jobs (default: auto)"
    echo ""
    echo "Examples:"
    echo "  $0                      # Standard release build"
    echo "  $0 -c -d -t             # Clean debug build with tests"
    echo "  $0 --install-prefix=\$HOME/.local -i  # Install to home directory"
    echo ""
}

# Main script
main() {
    # Default options
    BUILD_TYPE="Release"
    CLEAN_BUILD=false
    RUN_TESTS=false
    DO_INSTALL=false
    SKIP_CHECKS=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                print_usage
                exit 0
                ;;
            -c|--clean)
                CLEAN_BUILD=true
                shift
                ;;
            -d|--debug)
                BUILD_TYPE="Debug"
                shift
                ;;
            -r|--release)
                BUILD_TYPE="Release"
                shift
                ;;
            -t|--test)
                RUN_TESTS=true
                shift
                ;;
            -i|--install)
                DO_INSTALL=true
                shift
                ;;
            --install-prefix=*)
                INSTALL_PREFIX="${1#*=}"
                shift
                ;;
            --skip-checks)
                SKIP_CHECKS=true
                shift
                ;;
            -j|--jobs)
                JOBS="$2"
                shift 2
                ;;
            --jobs=*)
                JOBS="${1#*=}"
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                print_usage
                exit 1
                ;;
        esac
    done
    
    # Print header
    print_header
    
    # Check requirements
    if [ "$SKIP_CHECKS" = false ]; then
        if ! check_requirements; then
            exit 1
        fi
    else
        print_warning "Skipping dependency checks"
        echo ""
    fi
    
    # Verify source files
    if ! verify_source_files; then
        exit 1
    fi
    
    # Create configuration
    create_config
    
    # Clean build if requested
    if [ "$CLEAN_BUILD" = true ]; then
        clean_build
    fi
    
    # Ensure build directory exists
    if [ ! -d "$BUILD_DIR" ]; then
        mkdir -p "$BUILD_DIR"
    fi
    
    # Configure
    if ! configure_build; then
        exit 1
    fi
    
    # Build
    if ! build_project; then
        exit 1
    fi
    
    # Run tests if requested
    if [ "$RUN_TESTS" = true ]; then
        run_tests
    fi
    
    # Install if requested
    if [ "$DO_INSTALL" = true ]; then
        install_project
    fi
    
    # Print summary
    print_summary
    
    print_success "Build process completed successfully!"
}

# Run main function
main "$@"
