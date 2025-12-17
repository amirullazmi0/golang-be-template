package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/amirullazmi0/kratify-backend/config"
	"github.com/amirullazmi0/kratify-backend/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection using database/sql
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Connection pool settings
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connected successfully", zap.String("host", cfg.Host))

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}
