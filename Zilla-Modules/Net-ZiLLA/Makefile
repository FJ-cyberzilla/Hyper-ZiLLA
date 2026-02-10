# Net-ZiLLA Professional Enterprise Makefile
# Build system with features

.DEFAULT_GOAL := help
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -c

# Project Configuration
BINARY_NAME := netzilla
MAIN_PATH := cmd/netzilla/main.go
INSTALL_PATH := /usr/local/bin/$(BINARY_NAME)
CONFIG_DIR := /etc/netzilla
REPORTS_DIR := $(HOME)/.netzilla/reports
DB_DIR := $(HOME)/.netzilla/data
LOG_DIR := $(HOME)/.netzilla/logs

# Build Configuration
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION ?= $(shell go version | awk '{print $$3}')
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Race detector support (not available on all platforms, e.g., android/arm64)
HAS_RACE := $(shell go help build | grep -q "\-race" && [ "$(GOOS)/$(GOARCH)" != "android/arm64" ] && echo "yes" || echo "no")
ifeq ($(HAS_RACE),yes)
    TEST_RACE := -race
else
    TEST_RACE :=
endif

# Build Flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"
GOFLAGS := -v
CGO_ENABLED := 1

# Code Quality Tools
GOLANGCI_LINT_VERSION := v1.59.0
GOFMT := gofmt
GOVET := go vet

# Color Definitions
BOLD := \033[1m
CYAN := \033[0;36m
GREEN := \033[0;32m
RED := \033[0;31m
YELLOW := \033[1;33m
MAGENTA := \033[0;35m
NC := \033[0m

# Utility Functions
define log_header
	printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)\n"; \
	printf "$(BOLD)$(CYAN)▸ $(1)$(NC)\n"; \
	printf "$(BOLD)$(CYAN)━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━$(NC)\n"
endef

define log_success
	printf "$(GREEN)✓ $(1)$(NC)\n"
endef

define log_error
	printf "$(RED)✗ $(1)$(NC)\n"
endef

define log_info
	printf "$(CYAN)ℹ $(1)$(NC)\n"
endef

define log_warning
	printf "$(YELLOW)⚠ $(1)$(NC)\n"
endef

# Phony Targets
.PHONY: help all build install uninstall clean test lint fmt vet coverage \
		security benchmark docker check-deps version release debug

# ============================================================================
# PRIMARY TARGETS
# ============================================================================

help:
	@printf "\n$(BOLD)Net-ZiLLA Build System$(NC)\n"
	@printf "$(BOLD)Version:$(NC) $(VERSION) | $(BOLD)Go:$(NC) $(GO_VERSION)\n\n"
	@printf "$(BOLD)Core Targets:$(NC)\n"
	@printf "  $(CYAN)build$(NC)           Build the application binary\n"
	@printf "  $(CYAN)install$(NC)         Build and install to system (requires sudo)\n"
	@printf "  $(CYAN)uninstall$(NC)       Remove application and system config\n"
	@printf "  $(CYAN)clean$(NC)           Remove build artifacts and temporary files\n"
	@printf "\n$(BOLD)Quality & Testing:$(NC)\n"
	@printf "  $(CYAN)test$(NC)            Run comprehensive test suite\n"
	@printf "  $(CYAN)coverage$(NC)        Generate test coverage report\n"
	@printf "  $(CYAN)lint$(NC)            Run golangci-lint static analysis\n"
	@printf "  $(CYAN)fmt$(NC)             Format code with gofmt\n"
	@printf "  $(CYAN)vet$(NC)             Run go vet analysis\n"
	@printf "  $(CYAN)security$(NC)        Run security vulnerability scan\n"
	@printf "\n$(BOLD)Advanced:$(NC)\n"
	@printf "  $(CYAN)benchmark$(NC)       Run performance benchmarks\n"
	@printf "  $(CYAN)docker$(NC)          Build Docker image\n"
	@printf "  $(CYAN)check-deps$(NC)      Verify dependencies and updates\n"
	@printf "  $(CYAN)debug$(NC)           Build with debug symbols\n"
	@printf "  $(CYAN)release$(NC)         Create optimized release build\n"
	@printf "  $(CYAN)version$(NC)         Display version information\n"
	@printf "\n"

version:
	@$(call log_header,Version Information)
	@printf "$(BOLD)Binary:$(NC)      $(BINARY_NAME)\n"
	@printf "$(BOLD)Version:$(NC)     $(VERSION)\n"
	@printf "$(BOLD)Commit:$(NC)      $(COMMIT)\n"
	@printf "$(BOLD)Build Time:$(NC)  $(BUILD_TIME)\n"
	@printf "$(BOLD)Go Version:$(NC)  $(GO_VERSION)\n"
	@printf "\n"

# ============================================================================
# BUILD TARGETS
# ============================================================================

all: clean build

build: check-deps
	@$(call log_header,Building $(BINARY_NAME))
	@printf "$(CYAN)Target:$(NC) $(MAIN_PATH)\n"
	@printf "$(CYAN)Output:$(NC) ./$(BINARY_NAME)\n"
	@printf "$(CYAN)Flags:$(NC) $(LDFLAGS)\n\n"
	@if go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH); then \
		$(call log_success,Build completed); \
		ls -lh $(BINARY_NAME); \
	else \
		$(call log_error,Build failed); \
		exit 1; \
	fi
	@printf "\n"

debug: check-deps
	@$(call log_header,Debug Build)
	@printf "$(YELLOW)Note:$(NC) Building with debug symbols and no optimization\n\n"
	@GCFLAGS="-N -l" go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME)-debug $(MAIN_PATH)
	@$(call log_success,Debug build complete: ./$(BINARY_NAME)-debug)
	@printf "\n"

release: check-deps
	@$(call log_header,Release Build)
	@printf "$(CYAN)Creating optimized production binary$(NC)\n\n"
	@go build $(GOFLAGS) $(LDFLAGS) -ldflags "$(LDFLAGS) -s -w" -trimpath -o $(BINARY_NAME)-$(VERSION) $(MAIN_PATH)
	@$(call log_success,Release build: ./$(BINARY_NAME)-$(VERSION))
	@ls -lh $(BINARY_NAME)-$(VERSION)
	@printf "\n"

# ============================================================================
# INSTALLATION TARGETS
# ============================================================================

install: build
	@$(call log_header,Installing $(BINARY_NAME))
	@if [ ! -w /usr/local/bin ]; then \
		$(call log_warning,Requires sudo for system installation); \
	fi
	@printf "$(CYAN)Installing binary...$(NC)\n"
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)
	@sudo chmod 755 $(INSTALL_PATH)
	@$(call log_success,Binary installed: $(INSTALL_PATH))
	@printf "\n$(CYAN)Setting up directories...$(NC)\n"
	@sudo mkdir -p $(CONFIG_DIR)
	@mkdir -p $(REPORTS_DIR) $(DB_DIR) $(LOG_DIR)
	@chmod 700 $(REPORTS_DIR) $(DB_DIR) $(LOG_DIR)
	@$(call log_success,Directories created)
	@printf "\n$(CYAN)Configuring system...$(NC)\n"
	@if [ ! -f $(CONFIG_DIR)/config.yaml ]; then \
		if [ -f config.yaml ]; then \
			sudo cp config.yaml $(CONFIG_DIR)/config.yaml; \
			$(call log_success,Configuration deployed); \
		else \
			$(call log_warning,config.yaml not found in project root); \
		fi \
	else \
		$(call log_info,Using existing configuration); \
	fi
	@printf "\n"
	@$(call log_success,Installation complete!)
	@printf "$(CYAN)Quick Start:$(NC) $(BINARY_NAME) --help\n"
	@printf "$(CYAN)Config:$(NC) sudo cat $(CONFIG_DIR)/config.yaml\n"
	@printf "$(CYAN)Data:$(NC) ls -la $(REPORTS_DIR)\n\n"

uninstall:
	@$(call log_header,Uninstalling $(BINARY_NAME))
	@if [ -f $(INSTALL_PATH) ]; then \
		$(call log_info,Removing binary); \
		sudo rm -f $(INSTALL_PATH); \
		$(call log_success,Binary removed); \
	else \
		$(call log_warning,Binary not found at $(INSTALL_PATH)); \
	fi
	@if [ -d $(CONFIG_DIR) ]; then \
		$(call log_info,Removing system configuration); \
		sudo rm -rf $(CONFIG_DIR); \
		$(call log_success,System config removed); \
	fi
	@printf "\n$(YELLOW)User data preserved:$(NC)\n"
	@printf "  Reports: $(REPORTS_DIR)\n"
	@printf "  Database: $(DB_DIR)\n"
	@printf "  Logs: $(LOG_DIR)\n"
	@printf "$(CYAN)Remove manually if desired$(NC)\n\n"
	@$(call log_success,Uninstallation complete)
	@printf "\n"

# ============================================================================
# CODE QUALITY & TESTING
# ============================================================================

test: check-deps
	@$(call log_header,Running Test Suite)
	@printf "$(CYAN)Coverage Mode: $(NC)enabled\n"
	@printf "$(CYAN)Timeout: $(NC)5m\n"
	@if [ -z "$(TEST_RACE)" ]; then printf "$(YELLOW)Race Detector: $(NC)disabled (not supported on $(GOOS)/$(GOARCH))\n"; else printf "$(CYAN)Race Detector: $(NC)enabled\n"; fi
	@printf "\n"
	@go test -v $(TEST_RACE) -timeout 5m -coverprofile=coverage.out ./... 2>&1 | tee test-results.log | awk '{ \
		if ($$1 == "PASS" || $$0 ~ /ok.*coverage/) print "\033[0;32m" $$0 "\033[0m"; \
		else if ($$1 == "FAIL" || $$0 ~ /FAIL/) print "\033[0;31m" $$0 "\033[0m"; \
		else if ($$0 ~ /^===.*RUN/) print "\033[0;33m" $$0 "\033[0m"; \
		else print $$0 \
	}'
	@if [ -f coverage.out ]; then \
		$(call log_success,Coverage report generated); \
		go tool cover -func=coverage.out | tail -1; \
	fi
	@printf "\n"

coverage: test
	@$(call log_header,Coverage Analysis)
	@if [ -f coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		$(call log_success,HTML coverage report generated); \
		printf "$(CYAN)Open: $(NC)coverage.html\n\n"; \
	else \
		$(call log_error,Coverage file not found. Run 'make test' first); \
		exit 1; \
	fi

lint: check-deps
	@$(call log_header,Static Analysis)
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run ./... --timeout=5m; \
		$(call log_success,Linting passed); \
	else \
		$(call log_warning,golangci-lint not installed); \
		printf "$(CYAN)Install:$(NC) curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin\n\n"; \
	fi

fmt:
	@$(call log_header,Code Formatting)
	@if [ -z "$$($(GOFMT) -l .)" ]; then \
		$(call log_success,All files properly formatted); \
	else \
		$(call log_info,Formatting files...); \
		$(GOFMT) -w .; \
		$(call log_success,Formatting complete); \
	fi
	@printf "\n"

vet: check-deps
	@$(call log_header,Go Vet Analysis)
	@$(GOVET) ./... && $(call log_success,No issues detected) || { $(call log_error,Issues found); exit 1; }
	@printf "\n"

security:
	@$(call log_header,Security Vulnerability Scan)
	@if command -v gosec &> /dev/null; then \
		gosec -fmt sarif -out gosec-report.json ./... > /dev/null 2>&1; \
		$(call log_success,Security scan complete); \
		printf "$(CYAN)Report:$(NC) gosec-report.json\n\n"; \
	else \
		$(call log_warning,gosec not installed); \
		printf "$(CYAN)Install:$(NC) go install github.com/securego/gosec/v2/cmd/gosec@latest\n\n"; \
	fi

benchmark:
	@$(call log_header,Performance Benchmarks)
	@go test -bench=. -benchmem -benchtime=3s ./... 2>&1 | tee benchmark-results.txt
	@$(call log_success,Benchmark complete - see benchmark-results.txt)
	@printf "\n"

# ============================================================================
# MAINTENANCE TARGETS
# ============================================================================

clean:
	@$(call log_header,Cleanup)
	@printf "$(CYAN)Removing build artifacts...$(NC)\n"
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-* *.out *.log *.json *.html
	@go clean -testcache -modcache
	@$(call log_success,Cleanup complete)
	@printf "\n"

check-deps:
	@if ! command -v go &> /dev/null; then \
		$(call log_error,Go not installed); \
		exit 1; \
	fi
	@if [ ! -f go.mod ]; then \
		$(call log_error,go.mod not found); \
		exit 1; \
	fi
	@go mod verify > /dev/null 2>&1 || { $(call log_warning,Syncing dependencies); go mod tidy; }

docker:
	@$(call log_header,Docker Build)
	@if [ -f Dockerfile ]; then \
		docker build -t $(BINARY_NAME):$(VERSION) \
			--build-arg VERSION=$(VERSION) \
			--build-arg COMMIT=$(COMMIT) \
			--build-arg BUILD_TIME=$(BUILD_TIME) .; \
		$(call log_success,Docker image built: $(BINARY_NAME):$(VERSION)); \
	else \
		$(call log_error,Dockerfile not found); \
		exit 1; \
	fi
	@printf "\n"

.SILENT: help version
