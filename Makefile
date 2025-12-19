# Simple Makefile for a Go project
include .env

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main services/catalog-service/cmd/main.go

# Run the application
run:
	@go run services/catalog-service/cmd/main.go

# Run the identity service
run-identity:
	@GOOGLE_APPLICATION_CREDENTIALS=./service-account.json go run services/identity/cmd/main.go

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

DB_URL=postgres://$(BLUEPRINT_DB_USERNAME):$(BLUEPRINT_DB_PASSWORD)@$(BLUEPRINT_DB_HOST):$(BLUEPRINT_DB_PORT)/$(BLUEPRINT_DB_DATABASE)?sslmode=disable&search_path=$(BLUEPRINT_DB_SCHEMA)

dbml:
	@echo "Generating schema.dbml..."
	@docker compose exec -T -e PGPASSWORD=$(BLUEPRINT_DB_PASSWORD) psql_bp pg_dump -U $(BLUEPRINT_DB_USERNAME) -d $(BLUEPRINT_DB_DATABASE) --schema-only > schema.sql
	@sql2dbml schema.sql --postgres -o schema.dbml
	@rm schema.sql
	@echo "Done."

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./pkg/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Tidy all modules
tidy:
	@echo "Tidying modules..."
	@cd services/catalog-service && go mod tidy
	@cd services/migration && go mod tidy
	@cd pkg/database && go mod tidy
	@cd pkg/models && go mod tidy
	@cd pkg/middlewares && go mod tidy
	@cd pkg/storage && go mod tidy
	@echo "Done."

# Generate Protobuf code
proto:
	@echo "Generating Protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/auth/auth.proto
	@echo "Done."

# Generate Swagger docs
swagger:
	@echo "Generating Swagger docs..."
	@cd services/catalog-service && go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go --output docs --parseDependency --parseInternal
	@echo "Done."

.PHONY: all build run test clean watch docker-run docker-down itest tidy proto migrate migrate-down migrate-status migrate-new swagger

# Migration targets
migrate:
	@echo "Running migrations..."
	@go run services/migration/cmd/main.go -cmd=up

migrate-down:
	@echo "Rolling back last migration..."
	@go run services/migration/cmd/main.go -cmd=down

migrate-status:
	@echo "Checking migration status..."
	@go run services/migration/cmd/main.go -cmd=status

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=migration_name"; \
		exit 1; \
	fi
	@version=$$(date +%Y%m%d%H%M%S); \
	file="services/migration/internal/migrations/$${version}_$(name).go"; \
	echo "Creating migration: $$file"; \
	echo 'package migrations' > $$file; \
	echo '' >> $$file; \
	echo 'import (' >> $$file; \
	echo '	"gorm.io/gorm"' >> $$file; \
	echo ')' >> $$file; \
	echo '' >> $$file; \
	echo '// TODO: Add your migration struct name' >> $$file; \
	echo 'type Migration'$$version' struct{}' >> $$file; \
	echo '' >> $$file; \
	echo 'func (m *Migration'$$version') Version() string {' >> $$file; \
	echo '	return "'$$version'"' >> $$file; \
	echo '}' >> $$file; \
	echo '' >> $$file; \
	echo 'func (m *Migration'$$version') Name() string {' >> $$file; \
	echo '	return "$(name)"' >> $$file; \
	echo '}' >> $$file; \
	echo '' >> $$file; \
	echo 'func (m *Migration'$$version') Up(db *gorm.DB) error {' >> $$file; \
	echo '	// TODO: Implement migration' >> $$file; \
	echo '	return nil' >> $$file; \
	echo '}' >> $$file; \
	echo '' >> $$file; \
	echo 'func (m *Migration'$$version') Down(db *gorm.DB) error {' >> $$file; \
	echo '	// TODO: Implement rollback' >> $$file; \
	echo '	return nil' >> $$file; \
	echo '}' >> $$file; \
	echo "Created migration: $$file"; \
	echo "Don't forget to register it in registry.go!"
