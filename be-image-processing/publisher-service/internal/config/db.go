package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
}

type InitDatabaseParams struct {
	Conf *DatabaseConfig
}

func InitDB(params *InitDatabaseParams) (*sql.DB, error) {
	// Get database connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		slog.Error("DB_HOST config is required")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		slog.Error("DB_PORT config is required")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		slog.Error("DB_USER config is required")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		slog.Error("DB_PASSWORD config is required")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		slog.Error("DB_NAME config is required")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params.Conf.Host, params.Conf.Port, params.Conf.Username, params.Conf.Password, params.Conf.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("Successfully connected to database")
	return db, nil
}
