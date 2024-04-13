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

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config is required: %v", err)
	}

	log.Printf("try to migrate from %s to %s\n", cfg.StorageCfg.MigrationsPath, cfg.StorageCfg.PGUrl)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.StorageCfg.MigrationsPath),
		fmt.Sprintf("%s?x-migrations-table=%s&sslmode=disable",
			cfg.StorageCfg.PGUrl,
			cfg.StorageCfg.MigrationsTable,
		),
	)

	if err != nil {
		log.Fatalf("can't create migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatal(err)
	}

	log.Println("migrations successfully applied")
}
