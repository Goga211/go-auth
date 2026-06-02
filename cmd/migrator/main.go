package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	godotenv.Load()
	viper.AutomaticEnv()

	storagePath := viper.GetString("STORAGE_PATH")
	migrationsPath := viper.GetString("MIGRATIONS_PATH")
	migrationsTable := viper.GetString("MIGRATIONS_TABLE")

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
		panic("STORAGE_PATH is required")
	}
	if migrationsPath == "" {
		panic("MIGRATIONS_PATH is required")
	}
}
