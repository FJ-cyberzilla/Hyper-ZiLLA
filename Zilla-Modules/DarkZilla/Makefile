# Entynet Hacker Tools Enterprise - Makefile
# Enterprise Level Penetration Testing Framework

# Configuration
BINARY_NAME = entynet-enterprise
BINARY_DIR = target
RELEASE_DIR = $(BINARY_DIR)/release
DEBUG_DIR = $(BINARY_DIR)/debug
INSTALL_DIR = /usr/local/bin
CONFIG_DIR = /etc/entynet
DATA_DIR = /var/lib/entynet
LOG_DIR = /var/log/entynet

# Build flags
RUSTFLAGS = -C target-cpu=native -C opt-level=3
CARGO_FLAGS = 

# Colors for output
GREEN = \033[0;32m
BLUE = \033[0;34m
RED = \033[0;31m
YELLOW = \033[0;33m
PINK = \033[38;5;206m
NC = \033[0m

# Test configuration
test-config:
	@echo "$(YELLOW)Testing configuration system...$(NC)"
	@cargo test --test config_tests

# Generate default config
gen-config:
	@echo "$(YELLOW)Generating default configuration...$(NC)"
	@cargo run --bin gen-config

# Validate config
validate-config:
	@echo "$(YELLOW)Validating configuration...$(NC)"
	@cargo run --bin validate-config
# Platform detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	PLATFORM = linux
	INSTALL_USER = root
	SERVICE_DIR = /etc/systemd/system
else ifeq ($(UNAME_S),Darwin)
	PLATFORM = macos
	INSTALL_USER = $(USER)
else
	PLATFORM = windows
endif

.PHONY: all build release debug install uninstall clean test lint fmt audit docs run dev setup deps update dist docker help

# Default target
all: banner build

# Display banner
banner:
	@echo "$(PINK)"
	@echo "░█▀▄░█▀█░█▀▄░█░█░▀▀█░▀█▀░█░░░█░░░█▀█░░"
	@echo "░█░█░█▀█░█▀▄░█▀▄░▄▀░░░█░░█░░░█░░░█▀█░░"
	@echo "░▀▀░░▀░▀░▀░▀░▀░▀░▀▀▀░▀▀▀░▀▀▀░▀▀▀░▀░▀░░"
	@echo "$(NC)"
	@echo "$(BLUE)=== Entynet Hacker Tools Enterprise Build System ===$(NC)"

# Build in debug mode
build: banner
	@echo "$(YELLOW)Building in debug mode...$(NC)"
	@cargo build $(CARGO_FLAGS)
	@echo "$(GREEN)Build completed! Binary: $(BINARY_DIR)/debug/$(BINARY_NAME)$(NC)"

# Build in release mode
release: banner
	@echo "$(YELLOW)Building optimized release...$(NC)"
	@RUSTFLAGS="$(RUSTFLAGS)" cargo build --release $(CARGO_FLAGS)
	@strip $(RELEASE_DIR)/$(BINARY_NAME) 2>/dev/null || true
	@echo "$(GREEN)Release build completed! Binary: $(RELEASE_DIR)/$(BINARY_NAME)$(NC)"

# Debug build with full symbols
debug: banner
	@echo "$(YELLOW)Building debug version with symbols...$(NC)"
	@cargo build $(CARGO_FLAGS)
	@echo "$(GREEN)Debug build completed!$(NC)"

# Install system-wide
install: release
	@echo "$(YELLOW)Installing $(BINARY_NAME) system-wide...$(NC)"
	@if [ "$(PLATFORM)" = "linux" ]; then \
		sudo install -D -m 755 $(RELEASE_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME); \
		sudo mkdir -p $(CONFIG_DIR) $(DATA_DIR) $(LOG_DIR); \
		sudo chown $(INSTALL_USER):$(INSTALL_USER) $(CONFIG_DIR) $(DATA_DIR) $(LOG_DIR); \
		echo "$(GREEN)Installation completed!$(NC)"; \
		echo "$(BLUE)Binary: $(INSTALL_DIR)/$(BINARY_NAME)$(NC)"; \
		echo "$(BLUE)Config: $(CONFIG_DIR)$(NC)"; \
		echo "$(BLUE)Data: $(DATA_DIR)$(NC)"; \
		echo "$(BLUE)Logs: $(LOG_DIR)$(NC)"; \
	else \
		install -D -m 755 $(RELEASE_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME); \
		echo "$(GREEN)Installation completed!$(NC)"; \
	fi

# Uninstall
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME)...$(NC)"
	@if [ "$(PLATFORM)" = "linux" ]; then \
		sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
		sudo rm -rf $(CONFIG_DIR) $(DATA_DIR) $(LOG_DIR); \
	else \
		rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)Uninstallation completed!$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@cargo clean
	@rm -rf dist/ *.deb *.rpm *.tar.gz
	@echo "$(GREEN)Clean completed!$(NC)"

# Run tests
test: banner
	@echo "$(YELLOW)Running tests...$(NC)"
	@cargo test -- --test-threads=1
	@echo "$(GREEN)Tests completed!$(NC)"

# Run integration tests
test-integration:
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@cargo test --test '*' -- --test-threads=1 --nocapture
	@echo "$(GREEN)Integration tests completed!$(NC)"

# Lint code
lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@cargo clippy -- -D warnings
	@cargo fmt -- --check
	@echo "$(GREEN)Lint completed!$(NC)"

# Format code
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	@cargo fmt
	@echo "$(GREEN)Formatting completed!$(NC)"

# Security audit
audit:
	@echo "$(YELLOW)Running security audit...$(NC)"
	@cargo audit
	@echo "$(GREEN)Security audit completed!$(NC)"

# Generate documentation
docs:
	@echo "$(YELLOW)Generating documentation...$(NC)"
	@cargo doc --no-deps --open
	@echo "$(GREEN)Documentation generated!$(NC)"

# Run the application
run: build
	@echo "$(YELLOW)Running application...$(NC)"
	@cargo run

# Development run with watch
dev:
	@echo "$(YELLOW)Starting development mode with watch...$(NC)"
	@cargo watch -x run

# Setup development environment
setup:
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@rustup update
	@cargo install cargo-watch cargo-audit cargo-outdated
	@cargo install --locked cargo-deny
	@echo "$(GREEN)Development environment setup completed!$(NC)"

# Install dependencies
deps:
	@echo "$(YELLOW)Installing system dependencies...$(NC)"
	@if [ "$(PLATFORM)" = "linux" ]; then \
		if command -v apt-get >/dev/null 2>&1; then \
			sudo apt-get update && sudo apt-get install -y \
				libssl-dev \
				pkg-config \
				build-essential \
				cmake \
				libpcap-dev; \
		elif command -v yum >/dev/null 2>&1; then \
			sudo yum install -y \
				openssl-devel \
				pkgconfig \
				gcc-c++ \
				cmake \
				libpcap-devel; \
		elif command -v dnf >/dev/null 2>&1; then \
			sudo dnf install -y \
				openssl-devel \
				pkgconfig \
				gcc-c++ \
				cmake \
				libpcap-devel; \
		fi \
	fi
	@echo "$(GREEN)Dependencies installed!$(NC)"

# Update dependencies
update:
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	@cargo update
	@cargo outdated -R
	@echo "$(GREEN)Dependencies updated!$(NC)"

# Create distribution packages
dist: release
	@echo "$(YELLOW)Creating distribution packages...$(NC)"
	@mkdir -p dist
	@if [ "$(PLATFORM)" = "linux" ]; then \
		./scripts/create-deb.sh; \
		./scripts/create-rpm.sh; \
		./scripts/create-tarball.sh; \
		echo "$(GREEN)Distribution packages created in dist/$(NC)"; \
	else \
		./scripts/create-tarball.sh; \
		echo "$(GREEN)Distribution package created in dist/$(NC)"; \
	fi

# Build Docker image
docker:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t entynet/enterprise:latest .
	@echo "$(GREEN)Docker image built: entynet/enterprise:latest$(NC)"

# Run in Docker
docker-run:
	@echo "$(YELLOW)Running in Docker...$(NC)"
	@docker run -it --rm --privileged --network host entynet/enterprise:latest

# Performance benchmark
bench:
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	@cargo bench
	@echo "$(GREEN)Benchmarks completed!$(NC)"

# Check for outdated dependencies
outdated:
	@echo "$(YELLOW)Checking for outdated dependencies...$(NC)"
	@cargo outdated -R

# Size analysis
size:
	@echo "$(YELLOW)Analyzing binary size...$(NC)"
	@if command -v size >/dev/null 2>&1; then \
		size $(RELEASE_DIR)/$(BINARY_NAME); \
	else \
		echo "size command not available"; \
	fi
	@if command -v strip >/dev/null 2>&1; then \
		cp $(RELEASE_DIR)/$(BINARY_NAME) $(RELEASE_DIR)/$(BINARY_NAME).stripped; \
		strip $(RELEASE_DIR)/$(BINARY_NAME).stripped; \
		echo "Stripped size:"; \
		ls -lh $(RELEASE_DIR)/$(BINARY_NAME).stripped | awk '{print $$5}'; \
	fi

# Code statistics
stats:
	@echo "$(YELLOW)Code statistics:$(NC)"
	@echo "Lines of Rust code:"
	@find . -name "*.rs" -not -path "./target/*" -exec cat {} \; | wc -l
	@echo "Number of Rust files:"
	@find . -name "*.rs" -not -path "./target/*" | wc -l

# Create service file (Linux only)
service:
	@if [ "$(PLATFORM)" = "linux" ]; then \
		echo "$(YELLOW)Creating systemd service...$(NC)"; \
		sudo tee $(SERVICE_DIR)/entynet.service > /dev/null << EOF \
		[Unit] \
		Description=Entynet Hacker Tools Enterprise \
		After=network.target \
		\
		[Service] \
		Type=simple \
		User=root \
		ExecStart=$(INSTALL_DIR)/$(BINARY_NAME) \
		Restart=on-failure \
		RestartSec=5s \
		\
		[Install] \
		WantedBy=multi-user.target \
		EOF \
		echo "$(GREEN)Service file created at $(SERVICE_DIR)/entynet.service$(NC)"; \
		echo "$(YELLOW)To enable: sudo systemctl enable entynet.service$(NC)"; \
	else \
		echo "$(RED)Service creation only available on Linux$(NC)"; \
	fi

# Help target
help: banner
	@echo "$(BLUE)Available targets:$(NC)"
	@echo "$(GREEN)  all$(NC)        - Default build (debug mode)"
	@echo "$(GREEN)  build$(NC)     - Build in debug mode"
	@echo "$(GREEN)  release$(NC)   - Build optimized release"
	@echo "$(GREEN)  install$(NC)   - Install system-wide"
	@echo "$(GREEN)  uninstall$(NC) - Uninstall"
	@echo "$(GREEN)  clean$(NC)     - Clean build artifacts"
	@echo "$(GREEN)  test$(NC)      - Run tests"
	@echo "$(GREEN)  lint$(NC)      - Run clippy and format check"
	@echo "$(GREEN)  fmt$(NC)       - Format code"
	@echo "$(GREEN)  audit$(NC)     - Security audit"
	@echo "$(GREEN)  run$(NC)       - Run application"
	@echo "$(GREEN)  dev$(NC)       - Development mode with watch"
	@echo "$(GREEN)  setup$(NC)     - Setup development environment"
	@echo "$(GREEN)  deps$(NC)      - Install system dependencies"
	@echo "$(GREEN)  docker$(NC)    - Build Docker image"
	@echo "$(GREEN)  dist$(NC)      - Create distribution packages"
	@echo "$(GREEN)  bench$(NC)     - Run benchmarks"
	@echo "$(GREEN)  help$(NC)      - Show this help"

# Print build info
info:
	@echo "$(BLUE)Build Information:$(NC)"
	@echo "Binary: $(BINARY_NAME)"
	@echo "Platform: $(PLATFORM)"
	@echo "Install Dir: $(INSTALL_DIR)"
	@echo "Config Dir: $(CONFIG_DIR)"
	@echo "Rust Version: $(shell rustc --version)"
	@echo "Cargo Version: $(shell cargo --version)"

# Create distribution script directory
$(shell mkdir -p scripts)

# Include additional makefiles if they exist
-include makefile.local
