# Cognizilla Quantum Makefile
# FJ-Cyberzilla Sovereign Build System

# Variables
APP_NAME = cognizilla
VERSION = 1.0.0-sovereign
QUANTUM_FLAGS = -tags=quantum,sovereign
BUILD_FLAGS = -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +%Y%m%d.%H%M%S)"

# Colors
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
RESET := \033[0m

.PHONY: all build clean test deploy install uninstall help

all: deps build

## Build targets
build: 
	@echo "$(CYAN)🦖 Building Cognizilla Quantum...$(RESET)"
	@cd backend && go build $(BUILD_FLAGS) $(QUANTUM_FLAGS) -o ../bin/$(APP_NAME) ./cmd/cognizilla/
	@echo "$(GREEN)✅ Cognizilla built successfully!$(RESET)"

build-frontend:
	@echo "$(CYAN)🌐 Building Quantum Dashboard...$(RESET)"
	@cd frontend && npm run build
	@echo "$(GREEN)✅ Dashboard built successfully!$(RESET)"

build-all: build build-frontend

## Platform specific builds
build-linux:
	@echo "$(CYAN)🐧 Building for Linux...$(RESET)"
	@cd backend && GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(QUANTUM_FLAGS) -o ../bin/$(APP_NAME)-linux ./cmd/cognizilla/

build-macos:
	@echo "$(CYAN)🍎 Building for macOS...$(RESET)"
	@cd backend && GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(QUANTUM_FLAGS) -o ../bin/$(APP_NAME)-darwin ./cmd/cognizilla/

build-windows:
	@echo "$(CYAN)🪟 Building for Windows...$(RESET)"
	@cd backend && GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(QUANTUM_FLAGS) -o ../bin/$(APP_NAME)-windows.exe ./cmd/cognizilla/

build-termux:
	@echo "$(CYAN)📱 Building for Termux (Android)...$(RESET)"
	@cd backend && GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) $(QUANTUM_FLAGS) -o ../bin/$(APP_NAME)-termux ./cmd/cognizilla/

build-all-platforms: build-linux build-macos build-windows build-termux

## Installation targets
install-linux: build-linux
	@echo "$(CYAN)🐧 Installing on Linux...$(RESET)"
	@sudo cp bin/$(APP_NAME)-linux /usr/local/bin/$(APP_NAME)
	@sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)✅ Cognizilla installed to /usr/local/bin/$(APP_NAME)$(RESET)"

install-macos: build-macos
	@echo "$(CYAN)🍎 Installing on macOS...$(RESET)"
	@sudo cp bin/$(APP_NAME)-darwin /usr/local/bin/$(APP_NAME)
	@sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)✅ Cognizilla installed to /usr/local/bin/$(APP_NAME)$(RESET)"

install-termux: build-termux
	@echo "$(CYAN)📱 Installing on Termux...$(RESET)"
	@cp bin/$(APP_NAME)-termux $PREFIX/bin/$(APP_NAME)
	@chmod +x $PREFIX/bin/$(APP_NAME)
	@echo "$(GREEN)✅ Cognizilla installed to $PREFIX/bin/$(APP_NAME)$(RESET)"

install-wsl: build-linux
	@echo "$(CYAN)🪟 Installing on WSL...$(RESET)"
	@sudo cp bin/$(APP_NAME)-linux /usr/local/bin/$(APP_NAME)
	@sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)✅ Cognizilla installed to /usr/local/bin/$(APP_NAME)$(RESET)"

## Docker targets
docker-build:
	@echo "$(CYAN)🐳 Building Docker image...$(RESET)"
	@docker build -t fjcyberzilla/$(APP_NAME):$(VERSION) .
	@echo "$(GREEN)✅ Docker image built: fjcyberzilla/$(APP_NAME):$(VERSION)$(RESET)"

docker-run:
	@echo "$(CYAN)🚀 Starting Cognizilla in Docker...$(RESET)"
	@docker run -d -p 8080:8080 -p 8443:8443 --name cognizilla-monster fjcyberzilla/$(APP_NAME):$(VERSION)

docker-stop:
	@docker stop cognizilla-monster
	@docker rm cognizilla-monster

## Development targets
deps:
	@echo "$(CYAN)📦 Installing dependencies...$(RESET)"
	@cd backend && go mod tidy
	@cd frontend && npm install
	@echo "$(GREEN)✅ Dependencies installed$(RESET)"

test:
	@echo "$(YELLOW)🧪 Running tests...$(RESET)"
	@cd backend && go test ./... -v
	@echo "$(GREEN)✅ Tests completed$(RESET)"

clean:
	@echo "$(RED)🧹 Cleaning build artifacts...$(RESET)"
	@rm -rf bin/ dist/ coverage/
	@docker system prune -f
	@echo "$(GREEN)✅ Clean completed$(RESET)"

## Deployment targets
deploy: build-all docker-build docker-run
	@echo "$(GREEN)🎉 Cognizilla deployed successfully!$(RESET)"
	@echo "$(CYAN)🌐 Dashboard: https://localhost:8443$(RESET)"
	@echo "$(CYAN)🔧 API: http://localhost:8080$(RESET)"

## Help
help:
	@echo "$(MAGENTA)🦖 Cognizilla Quantum Build System$(RESET)"
	@echo ""
	@echo "$(CYAN)Usage:$(RESET)"
	@echo "  make [target]"
	@echo ""
	@echo "$(CYAN)Targets:$(RESET)"
	@echo "  build               - Build Cognizilla"
	@echo "  build-all-platforms - Build for all platforms"
	@echo "  install-linux       - Install on Linux"
	@echo "  install-macos       - Install on macOS"
	@echo "  install-termux      - Install on Termux (Android)"
	@echo "  install-wsl         - Install on WSL"
	@echo "  docker-build        - Build Docker image"
	@echo "  docker-run          - Run in Docker"
	@echo "  deploy              - Full deployment"
	@echo "  deps                - Install dependencies"
	@echo "  test                - Run tests"
	@echo "  clean               - Clean build artifacts"
	@echo "  help                - Show this help"
