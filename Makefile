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

docker-start:
	docker run -d --name gRPC_server -p 44000:44000 grpc-server

docker-build:
	docker build . -t grpc-server --no-cache