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
		fmt.Sprintf("%s&x-migrations-table=%s",
			cfg.StorageCfg.PGUrl,
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
