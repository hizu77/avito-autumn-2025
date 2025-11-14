package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	Config struct {
		Postgres `envPrefix:"POSTGRES_"`
		HTTP     `envPrefix:"HTTP_"`
		Admin    `envPrefix:"ADMIN_"`
	}

	Postgres struct {
		URL string `env:"URL"`
	}

	HTTP struct {
		Host string `env:"HOST"`
		Port string `env:"PORT"`
	}

	Admin struct {
		Secret          string `env:"SECRET"`
		DefaultID       string `env:"DEFAULT_ID"`
		DefaultPassword string `env:"DEFAULT_PASSWORD"`
	}
)

func New(logger *zap.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn(
			"warning: no .env file, skipping loading",
		)
	}

	cfg := Config{}

	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		return nil, errors.Wrap(err, "parse env")
	}

	return &cfg, nil
}
