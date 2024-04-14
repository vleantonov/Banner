package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	ServerCfg  ServerConfig
	StorageCfg StorageConfig
	QueueCfg   QueueConfig
}

type ServerConfig struct {
	Host      string        `env:"APP_HOST" env-default:"localhost"`
	Port      int           `env:"APP_PORT" env-default:"8080"`
	AppSecret string        `env:"APP_SECRET" env-required:"true"`
	TTLToken  time.Duration `env:"TTL_TOKEN" env-default:"1h"`
}

type StorageConfig struct {
	MigrationsPath  string        `env:"MIGRATIONS_PATH" env-default:"./migrations"`
	MigrationsTable string        `env:"MIGRATIONS_TABLE" env-default:"versions"`
	PGUrl           string        `env:"PG_URL" env-required:"true"`
	TTLCache        time.Duration `env:"TTL_CACHE" env-default:"5m"`
}

type QueueConfig struct {
	RMQUrl string `env:"RMQ_URL" env-required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env error: %w", err)
	}

	return &cfg, nil
}
