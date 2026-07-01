# -----------------------------
# Application
# -----------------------------

run:
	cd backend && go run ./cmd/server

build:
	cd backend && go build -o bin/server ./cmd/server

# -----------------------------
# Frontend
# -----------------------------

frontend:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

# -----------------------------
# Docker
# -----------------------------

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up -d

logs:
	docker compose logs -f

# -----------------------------
# Database
# -----------------------------

psql:
	docker compose exec postgres psql -U postgres reservation

# -----------------------------
# Migration
# -----------------------------

migrate-up:
	migrate -path backend/migrations \
	-database "postgres://postgres:password@localhost:5432/reservation?sslmode=disable" up

migrate-down:
	migrate -path backend/migrations \
	-database "postgres://postgres:password@localhost:5432/reservation?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir backend/migrations -seq $(name)

# -----------------------------
# Test
# -----------------------------

test:
	cd backend && go test ./...

test-cover:
	cd backend && go test ./... -cover

# -----------------------------
# Lint
# -----------------------------

lint:
	cd backend && golangci-lint run

fmt:
	cd backend && go fmt ./...

# -----------------------------
# Swagger
# -----------------------------

swagger:
	cd backend && swag init -g cmd/server/main.go

# -----------------------------
# Batch
# -----------------------------

batch:
	cd backend && go run ./cmd/batch/daily_report

# -----------------------------
# Clean
# -----------------------------

clean:
	rm -rf backend/bin

# -----------------------------
# Setup
# -----------------------------

setup:
	cd backend && go mod tidy
	cd frontend && npm install
	docker compose up -d
	make migrate-up
	make swagger