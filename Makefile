.PHONY: run build migrate-up migrate-down migrate-create test clean

# Run application
run:
	go run cmd/api/main.go

# Build application
build:
	go build -o bin/pos-api cmd/api/main.go

# Run migrations up
migrate-up:
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" up

# Run migrations down
migrate-down:
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" down

# Create new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/database/migrations -seq $$name

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install migrate tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Docker build
docker-build:
	docker build -t pos-backend:latest .

# Docker run
docker-run:
	docker run -p 8080:8080 --env-file .env pos-backend:latest
