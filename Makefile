default: run

migrate:
	go run cmd/migrator/migrationup.go

run:
	go run cmd/main.go