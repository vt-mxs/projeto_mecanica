package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
	BACKEND_PORT string `mapstructure:"BACKEND_PORT"`
}

func LoadENV() (*Config, error) {
	var cfg Config

	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")

	// ALTERADO: _ = gotenv.Load em vez de erro fatal -> permite env vars diretas em produção (Docker)
	_ = gotenv.Load(filepath.Join(root, ".env"))

	viper.SetDefault("BACKEND_PORT", "8080")

	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("BACKEND_PORT")

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("\nError to get variables: %w", err)
	}
	if cfg.BACKEND_PORT == "" {
		cfg.BACKEND_PORT = "8080"
	}

	return &cfg, nil
}
