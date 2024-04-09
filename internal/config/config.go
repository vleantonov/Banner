package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configFlag            = "config"
	emptyConfigPath       = ""
	configFlagDescription = "path to config file"
)

type Config struct {
	Env        string        `yaml:"env" env:"ENV" env-default:"development"`
	ServerCfg  ServerConfig  `yaml:"server"`
	StorageCfg StorageConfig `yaml:"storage"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"APP_HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"APP_PORT" env-default:"8080"`
}

type StorageConfig struct {
	MigrationsPath  string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-default:"./migrations"`
	MigrationsTable string `yaml:"migrations_table" env:"MIGRATIONS_TABLE" env-default:"versions"`
	PGUrl           string `yaml:"pg_url" env:"PG_URL" env-required:"true"`
}

func New(path string) (*Config, error) {
	var cfg Config
	if path != emptyConfigPath {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("read config error: %w", err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("read env error: %w", err)
		}
	}

	return &cfg, nil
}

func FetchConfigPath() string {
	var res string

	flag.StringVar(&res, configFlag, emptyConfigPath, configFlagDescription)
	flag.Parse()

	return res
}
