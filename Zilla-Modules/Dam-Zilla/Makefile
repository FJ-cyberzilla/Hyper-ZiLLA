# Multi-Language Project Makefile
.PHONY: help install test lint build clean

# Variables
NODE_MODULES = node_modules
PYTHON_VENV = venv
JULIA_DEPS = Manifest.toml

# Default target
help:
	@echo "Available targets:"
	@echo "  install - Install all dependencies"
	@echo "  test    - Run all test suites"
	@echo "  lint    - Run all linters"
	@echo "  build   - Build the project"
	@echo "  clean   - Clean generated files"
	@echo "  dev     - Start development servers"

# Installation targets
install: install-js install-python install-julia

install-js:
	@echo "📦 Installing JavaScript dependencies..."
	npm ci

install-python:
	@echo "🐍 Setting up Python environment..."
	python -m venv $(PYTHON_VENV) || true
	. $(PYTHON_VENV)/bin/activate && pip install --upgrade pip
	if [ -f requirements.txt ]; then \
		. $(PYTHON_VENV)/bin/activate && pip install -r requirements.txt; \
	fi

install-julia:
	@echo "🔬 Installing Julia dependencies..."
	if [ -f Project.toml ]; then \
		julia -e 'using Pkg; Pkg.activate("."); Pkg.instantiate()'; \
	fi

# Testing targets
test: test-js test-ts test-python test-julia test-shell

test-js:
	@echo "🧪 Running JavaScript tests..."
	npm test

test-ts:
	@echo "📘 Running TypeScript checks..."
	if [ -f tsconfig.json ]; then \
		npx tsc --noEmit; \
		npm run test:typescript || true; \
	fi

test-python:
	@echo "🐍 Running Python tests..."
	if [ -f pytest.ini ] || [ -f requirements.txt ]; then \
		. $(PYTHON_VENV)/bin/activate && python -m pytest tests/ -v || true; \
	fi

test-julia:
	@echo "🔬 Running Julia tests..."
	if [ -f Project.toml ]; then \
		julia -e 'using Pkg; Pkg.test()'; \
	fi

test-shell:
	@echo "🐚 Checking shell scripts..."
	find . -name "*.sh" -type f -exec shellcheck {} \; || true

# Linting targets
lint: lint-js lint-python lint-julia

lint-js:
	@echo "📏 Linting JavaScript/TypeScript..."
	npx eslint . --ext .js,.jsx,.ts,.tsx --fix --max-warnings=0
	npx prettier --write .

lint-python:
	@echo "🐍 Formatting Python..."
	if [ -d $(PYTHON_VENV) ]; then \
		. $(PYTHON_VENV)/bin/activate && \
		find . -name "*.py" -type f -exec black {} \; && \
		find . -name "*.py" -type f -exec pylint {} \; || true; \
	fi

lint-julia:
	@echo "🔬 Formatting Julia..."
	julia -e 'using Pkg; Pkg.add("JuliaFormatter"); using JuliaFormatter; format(".")'

# Build targets
build: build-js
	@echo "🏗️ Building project..."
	mkdir -p dist
	npm run build || true
	[ -d build ] && cp -r build/* dist/ || true

build-js:
	@echo "📦 Building JavaScript/TypeScript..."
	if [ -f package.json ]; then \
		npm run build || true; \
	fi

# Development
dev:
	@echo "🚀 Starting development servers..."
	npm run dev || true

# Cleanup
clean:
	@echo "🧹 Cleaning up..."
	rm -rf $(NODE_MODULES)
	rm -rf $(PYTHON_VENV)
	rm -rf dist build .pytest_cache __pycache__
	rm -f *.log
	julia -e 'using Pkg; Pkg.activate("."); Pkg.gc()'

# Utility targets
deps-tree:
	@echo "🌳 Dependency trees:"
	npm list --depth=1
	if [ -d $(PYTHON_VENV) ]; then \
		. $(PYTHON_VENV)/bin/activate && pip list; \
	fi

security-check:
	@echo "🛡️ Security checks..."
	npm audit
	if [ -d $(PYTHON_VENV) ]; then \
		. $(PYTHON_VENV)/bin/activate && safety check || true; \
	fi

# Quick setup for new developers
setup: install test lint
	@echo "✅ Project setup complete!"
