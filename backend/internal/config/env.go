package config

import (
	"fmt"
	"os"
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

	// BUG PRÉ-EXISTENTE: o compose.yaml injeta DB_HOST/DB_PORT/DB_NAME/
	// DB_USER/DB_PASSWORD, mas aqui só se lia DATABASE_URL. Dentro do
	// container o .env do host não existe (.dockerignore exclui .env),
	// então DATABASE_URL ficava vazia e o GORM caía no default
	// "user=root database=" a tentar ligar a um socket unix em /tmp.
	//
	// Composto a partir das variáveis que o compose de facto passa,
	// com o host "postgres" (o hostname interno da rede do compose) e
	// sslmode=disable, que é o modo do pg_hba.conf do container.
	if cfg.DATABASE_URL == "" {
		host := envOr("DB_HOST", "postgres")
		port := envOr("DB_PORT", "5432")
		name := os.Getenv("DB_NAME")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")

		if name != "" && user != "" {
			cfg.DATABASE_URL = fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s?sslmode=disable",
				user, pass, host, port, name,
			)
		}
	}

	// Sem URL nenhuma não há o que ligar. Falhar aqui é melhor do que
	// deixar o GORM tentar um socket unix e devolver um erro que não
	// diz nada sobre a causa.
	if cfg.DATABASE_URL == "" {
		return nil, fmt.Errorf(
			"sem DATABASE_URL: define-a, ou define DB_NAME/DB_USER (e DB_PASSWORD) para o compose montar a URL",
		)
	}

	return &cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
