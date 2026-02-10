# Makefile

.PHONY: help build test deploy clean

# Default target
help:
	@echo "Enterprise OSINT Platform Management"
	@echo ""
	@echo "Targets:"
	@echo "  build     - Build Docker images"
	@echo "  test      - Run tests"
	@echo "  deploy    - Deploy to environment"
	@echo "  clean     - Clean up resources"
	@echo "  logs      - View service logs"
	@echo "  monitor   - Open monitoring dashboard"

# Build images
build:
	docker-compose build

# Development deployment
dev: build
	docker-compose -f docker-compose.yml up -d

# Production deployment
prod: build
	docker-compose -f docker-compose.prod.yml up -d

# Run tests
test:
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# Clean up
clean:
	docker-compose down -v
	docker system prune -f

# View logs
logs:
	docker-compose logs -f

# Monitoring
monitor:
	open http://localhost:3000

# Database operations
db-migrate:
	docker-compose exec api alembic upgrade head

db-rollback:
	docker-compose exec api alembic downgrade -1

# Backup
backup:
	./scripts/backup.sh

# Security scan
security-scan:
	docker scan $(docker images -q)

# Performance test
perf-test:
	docker-compose -f docker-compose.perf.yml up
