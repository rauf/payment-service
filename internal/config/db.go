package config

import (
	"cmp"
	"os"

	"github.com/rauf/payment-service/internal/database"
)

type Config struct {
	Database database.Config
}

func NewConfig() *Config {
	return &Config{
		Database: database.Config{
			Driver:       "postgres",
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         5432,
			Username:     "postgres",
			Password:     "postgres",
			DatabaseName: "payment",
		},
	}
}

func getEnv(key, fallback string) string {
	return cmp.Or(os.Getenv(key), fallback)
}
