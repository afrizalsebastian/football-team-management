include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_PATH=./migration
SQLC_CONFIG=./sqlc.yaml
MAIN=./main.go

.PHONY: help migrate-create migrate-up migrate-down migrate-drop sqlc-generate sqlc-verify swag-init swag-fmt dep run
help:
	@echo "Available commands:"
	@echo "  make migrate-create name=<name>  Create migration"
	@echo "  make migrate-up                  Run migrations"
	@echo "  make migrate-down                Rollback 1 migration"
	@echo "  make migrate-drop                Drop database migrations"
	@echo "  make sqlc-generate               Generate sqlc code"
	@echo "  make sqlc-verify                 Verify sqlc"
	@echo "  make dep                					Install Go dependencies"
	@echo "  make run                         Run application"

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose down 1

migrate-drop:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose drop -f

sqlc-generate:
	sqlc generate -f $(SQLC_CONFIG)

sqlc-verify:
	sqlc vet -f $(SQLC_CONFIG)

swag-init:
	swag init -g $(MAIN) -o docs

swag-fmt:
	swag fmt -g $(MAIN)

dep:
	go mod tidy
	go mod download

run: dep swag-fmt swag-init
	go run main.go