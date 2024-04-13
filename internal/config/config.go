package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env        string `env:"ENV" env-default:"development"`
	ServerCfg  ServerConfig
	StorageCfg StorageConfig
}

type ServerConfig struct {
	Host       string `env:"APP_HOST" env-default:"localhost"`
	Port       int    `env:"APP_PORT" env-default:"8080"`
	UserToken  string `env:"USER_TOKEN" env-required:"true"`
	AdminToken string `env:"ADMIN_TOKEN" env-required:"true"`
}

type StorageConfig struct {
	MigrationsPath  string        `env:"MIGRATIONS_PATH" env-default:"./migrations"`
	MigrationsTable string        `env:"MIGRATIONS_TABLE" env-default:"versions"`
	PGUrl           string        `env:"PG_URL" env-required:"true"`
	TTLCache        time.Duration `env:"TTL_CACHE" env-default:"5m"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
