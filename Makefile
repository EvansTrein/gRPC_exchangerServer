default: run
PATH_DB = ./internal/storages/exchanger.db
FILE_MIGRATIONS = ./migrations

.PHONY: run run-default

migrate:
	go run cmd/migrator/migrationup.go -storage-path $(PATH_DB) -migrations-path $(FILE_MIGRATIONS)

run:
	go run cmd/main.go -config ./config.yaml

run-default:
	go run cmd/main.go -config default