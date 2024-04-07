package main

import (
	"banner/internal/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {

	cfg, err := config.New(config.FetchConfigPath())
	if err != nil {
		log.Fatal("config is required")
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.StorageCfg.MigrationsPath),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&x-migrations-table=%s",
			cfg.StorageCfg.PostgresCfg.Username,
			cfg.StorageCfg.PostgresCfg.Password,
			cfg.StorageCfg.PostgresCfg.Host,
			cfg.StorageCfg.PostgresCfg.Port,
			cfg.StorageCfg.PostgresCfg.DBName,
			cfg.StorageCfg.PostgresCfg.SSLMode,
			cfg.StorageCfg.PostgresCfg.SSLMode,
			cfg.StorageCfg.MigrationsTable,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatal(err)
	}
}
