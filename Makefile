tidy:
	cd src/adapter && go mod tidy
	cd src/core && go mod tidy
	cd src/migration && go mod tidy
	cd src/rest && go mod tidy

test:
	cd src/adapter && go mod tidy && go test ./...
	cd src/core && go mod tidy && go test ./...
	cd src/migration && go mod tidy && go test ./...
	cd src/rest && go mod tidy && go test ./...

migrate:
	cd src/migration && go mod tidy && go run main.go

rest:
	cd src/rest && go mod tidy && go run main.go

swagger:
	cd src/rest && swag init --parseDependency --parseDepth=3

MIGRATION_DIR=src/migration/migrations

TIMESTAMP=$(shell date +%Y%m%d%H%M)

migration:
	@read -p "Enter migration description (e.g., create_users_table): " desc; \
	touch $(MIGRATION_DIR)/$(TIMESTAMP)_$${desc}.up.sql; \
	echo "-- Migration: $${desc}" >> $(MIGRATION_DIR)/$(TIMESTAMP)_$${desc}.up.sql; \
	echo "Created migration: $(MIGRATION_DIR)/$(TIMESTAMP)_$${desc}.up.sql"