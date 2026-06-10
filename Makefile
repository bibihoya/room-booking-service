.PHONY: up down seed test test-cover lint clean integration e2e load

# Docker commands
up:
	docker-compose up --build -d

down:
	docker-compose down

seed:
	docker-compose exec -T postgres psql -U postgres -d meeting_booking < scripts/seed.sql
	@echo "Тестовые данные загружены"

# Testing commands
test:
	go test -v ./...

test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html and coverage.out"

# Linting
lint:
	@echo "Running linter..."
	golangci-lint run ./... --timeout=5m

# Integration tests
integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

# E2E tests
e2e:
	@echo "Running E2E tests..."
	go test -v -tags=e2e ./tests/e2e/...

# Load tests
load:
	@echo "Running load tests..."
	k6 run --duration 10s --vus 10 tests/load/slot-test.js

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html
	docker-compose down -v
	go clean -testcache

# Setup for CI (run everything)
ci: lint test-cover integration e2e
	@echo "CI pipeline completed successfully"

# Help
help:
	@echo "Available commands:"
	@echo "  make up          - Start services with Docker Compose"
	@echo "  make down        - Stop services"
	@echo "  make seed        - Seed test data into database"
	@echo "  make test        - Run all tests"
	@echo "  make test-cover  - Run tests with coverage report"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make integration - Run integration tests"
	@echo "  make e2e         - Run E2E tests"
	@echo "  make load        - Run load tests"
	@echo "  make clean       - Clean up generated files and containers"
	@echo "  make ci          - Run all CI checks"