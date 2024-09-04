package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const maxTries = 5

type Config struct {
	Driver       string
	Port         int
	Host         string
	Username     string
	Password     string
	DatabaseName string
	MaxOpen      int
	MaxIdle      int
	MaxLifetime  time.Duration
}

type Database struct {
	*sql.DB
	config *Config
}

func NewDatabase(dbConfig Config) (*Database, error) {
	connString := connectionString(dbConfig)
	var db *sql.DB
	var err error
	for range maxTries {
		db, err = connect(dbConfig.Driver, connString)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil || db == nil {
		return &Database{}, fmt.Errorf("error connecting to DB: %w", err)
	}

	db.SetMaxOpenConns(dbConfig.MaxOpen)
	db.SetMaxIdleConns(dbConfig.MaxIdle)
	db.SetConnMaxLifetime(dbConfig.MaxLifetime)

	return &Database{
		DB:     db,
		config: &dbConfig,
	}, nil
}

func connect(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening db connection: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("error closing db connection: %w, closeErr: %w", err, closeErr)
		}
		return nil, fmt.Errorf("error pinging db: %w", err)
	}
	return db, nil
}

func connectionString(config Config) string {
	switch config.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.Username, config.Password, config.DatabaseName,
		)
	default:
		return ""
	}
}
