package config

import "github.com/rauf/payment-service/internal/database"

type Config struct {
	Database database.Config
}

func NewConfig() *Config {
	return &Config{
		Database: database.Config{
			Driver:       "postgres",
			Host:         "localhost",
			Port:         5432,
			Username:     "postgres",
			Password:     "postgres",
			DatabaseName: "payment",
		},
	}
}
