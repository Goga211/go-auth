package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.StringVar(&migrationsTable, "migrations-table", "", "Path to a table containing the migration files")
	flag.Parse()

	validatePaths(storagePath, migrationsPath)

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations found")
			return
		}

		panic(err)
	}

	fmt.Println("Migrations applied successfully")
}

func validatePaths(storagePath string, migrationsPath string) {
	if storagePath == "" {
		panic("storage-path flag is required")
	}
	if migrationsPath == "" {
		panic("migrations-path flag is required")
	}
}
