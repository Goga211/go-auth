package main

import (
	"errors"
	"fmt"

	"github.com/Goga211/go-auth/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoadMigrator()

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.StoragePath, cfg.MigrationsTable),
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
