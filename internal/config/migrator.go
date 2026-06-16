package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type MigratorConfig struct {
	StoragePath     string `env:"STORAGE_PATH" env-required:"true"`
	MigrationsPath  string `env:"MIGRATIONS_PATH" env-required:"true"`
	MigrationsTable string `env:"MIGRATIONS_TABLE" env-default:"migrations"`
}

func MustLoadMigrator() *MigratorConfig {
	_ = godotenv.Load()

	var cfg MigratorConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read migrator config: " + err.Error())
	}

	return &cfg
}
