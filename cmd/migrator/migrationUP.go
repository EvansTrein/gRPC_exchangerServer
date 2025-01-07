package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var createForDbPath string
	var fileMigrationPath string

	flag.StringVar(&createForDbPath, "storage-path", "", "table creation path")
	flag.StringVar(&fileMigrationPath, "migrations-path", "", "path to migration file")
	flag.Parse()

	if createForDbPath == "" || fileMigrationPath == "" {
		panic("the path of the file with migrations or the path for database creation is not specified")
	}

	migrateDb, err := migrate.New("file://"+fileMigrationPath, "sqlite3://"+createForDbPath)
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
