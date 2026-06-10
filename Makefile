.PHONY: up down seed test lint

up:
	docker-compose up --build -d

down:
	docker-compose down

seed:
	docker-compose exec -T postgres psql -U postgres -d meeting_booking < scripts/seed.sql
	@echo "Тестовые данные загружены"

test:
	go test -v ./...

lint:
	golangci-lint run ./...