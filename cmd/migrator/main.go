package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
)

var pathCreateForDb = "./internal/storages/exchanger.db"

func main() {

	migrateDb, err := migrate.New("file://./migrations", "sqlite3://"+pathCreateForDb)
	if err != nil {
		panic(err)
	}

	if err := migrateDb.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	log.Println("migrations have been successfully applied")
}
