# Run api server
.PHONY: run
run:
	@go run ./cmd/app/main.go

###### ---- SWAGGER ---- ######
# Generate Swagger Docs
.PHONY: swagger
swagger:
	@swag init -g cmd/app/main.go --output docs/httpdoc

###### ---- DATABASE MIGRATION ---- ######
# Variables
MIGRATION_CMD=migrate
MIGRATION_RUNNER=./cmd/migration/main.go
MIGRATIONS_DIR=./database/migrations
VERSION =?
EXTENSION ?= sql

# Generate a Migration
.PHONY: migration-create
migration-create:
	@if [ "$(name)" = "" ]; then \
		echo "Error: name= parameter is required (e.g., make migration-create name=create_users_table)"; \
		exit 1; \
	fi
	$(MIGRATION_CMD) create -ext $(EXTENSION) -dir $(MIGRATIONS_DIR) -seq $(name)

# Apply all up migrations
.PHONY: migration-up
migration-up:
	@go run ${MIGRATION_RUNNER} up ${VERSION}

# Rollback the last migration
.PHONY: migration-down
migration-down:
	@go run ${MIGRATION_RUNNER} down

# Force a specific version
.PHONY: migration-force
migration-force:
	@go run ${MIGRATION_RUNNER} force ${VERSION}

# Drop the database
.PHONY: migration-drop
migration-drop:
	@go run ${MIGRATION_RUNNER} drop

###### ---- DATABASE SEEDER ---- ######
SEEDER_RUNNER=./cmd/seeder/main.go

# run list of seeders
.PHONY: seeder-list
seeder-list:
	@go run ${SEEDER_RUNNER} list

# run all seeder
.PHONY: seeder-run-all
seeder-run-all:
	@go run ${SEEDER_RUNNER} run:all

# run specific seeder
.PHONY: seeder-run-one
seeder-run-one:
	@if [ "$(name)" = "" ]; then \
		echo "Error: name= parameter is required (e.g., make seeder-run name=<seeder_name>)"; \
		exit 1; \
	fi
	@go run ${SEEDER_RUNNER} run:one ${name}
