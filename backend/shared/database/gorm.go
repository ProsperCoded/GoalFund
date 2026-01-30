package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gofund/shared/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	LogLevel logger.LogLevel
}

// NewGormDB creates a new GORM database connection
func NewGormDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// createEnumTypes creates PostgreSQL enum types required by the models before migration
func createEnumTypes(db *gorm.DB) error {
	log.Println("Creating enum types...")
	
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create user_role enum if it doesn't exist
	// Using DO block to check existence and create atomically
	query := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_type 
				WHERE typname = 'user_role'
			) THEN
				CREATE TYPE user_role AS ENUM ('user', 'admin');
			END IF;
		END $$;
	`
	
	if _, err := sqlDB.Exec(query); err != nil {
		// Check if enum already exists (handles race conditions or if it was created manually)
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'user_role')`
		if checkErr := sqlDB.QueryRow(checkQuery).Scan(&exists); checkErr == nil && exists {
			log.Println("Enum user_role already exists, continuing...")
			return nil
		}
		// If enum doesn't exist and creation failed, return error
		return fmt.Errorf("failed to create user_role enum: %w", err)
	}

	log.Println("Enum types created successfully")
	return nil
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running GORM auto-migration...")

	// User service models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.PasswordResetToken{},
	); err != nil {
		return fmt.Errorf("failed to migrate user models: %w", err)
	}

	// Goal service models
	if err := db.AutoMigrate(
		&models.Goal{},
		&models.Contribution{},
		&models.Proof{},
		&models.Vote{},
	); err != nil {
		return fmt.Errorf("failed to migrate goal models: %w", err)
	}

	// Ledger service models
	if err := db.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.LedgerEntry{},
		&models.BalanceSnapshot{},
	); err != nil {
		return fmt.Errorf("failed to migrate ledger models: %w", err)
	}

	log.Println("GORM auto-migration completed successfully")
	return nil
}



// SetupDatabase initializes the database with migrations and default data
func SetupDatabase(cfg Config) (*gorm.DB, error) {
	db, err := NewGormDB(cfg)
	if err != nil {
		return nil, err
	}

	// Create enum types before running migrations (GORM doesn't create enums automatically)
	if err := createEnumTypes(db); err != nil {
		return nil, fmt.Errorf("failed to create enum types: %w", err)
	}

	// Run auto-migration
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}



	return db, nil
}