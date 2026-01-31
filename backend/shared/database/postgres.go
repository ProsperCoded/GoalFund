package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresConfig holds database configuration for standard sql.DB connections
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establishes a PostgreSQL connection
func Connect(cfg PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}

// NewPostgresConnection establishes a PostgreSQL connection using sqlx
func NewPostgresConnection(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
