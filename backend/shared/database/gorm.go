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

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running GORM auto-migration...")

	// User service models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
		&models.Session{},
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

// CreateDefaultRoles creates default system roles
func CreateDefaultRoles(db *gorm.DB) error {
	defaultRoles := []models.Role{
		{
			Name:        "user",
			Description: "Standard user role",
			Permissions: map[string]interface{}{
				"goals.create": true,
				"goals.view":   true,
				"goals.update": "own",
				"payments.create": true,
				"profile.view": "own",
				"profile.update": "own",
			},
		},
		{
			Name:        "admin",
			Description: "Administrator role with full access",
			Permissions: map[string]interface{}{
				"*": true, // Full access
			},
		},
		{
			Name:        "moderator",
			Description: "Moderator role for community management",
			Permissions: map[string]interface{}{
				"goals.view":   true,
				"goals.moderate": true,
				"proofs.verify": true,
				"users.view": true,
			},
		},
	}

	for _, role := range defaultRoles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Role doesn't exist, create it
				if err := db.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create role %s: %w", role.Name, err)
				}
				log.Printf("Created default role: %s", role.Name)
			} else {
				return fmt.Errorf("failed to check role %s: %w", role.Name, err)
			}
		}
	}

	return nil
}

// SetupDatabase initializes the database with migrations and default data
func SetupDatabase(cfg Config) (*gorm.DB, error) {
	db, err := NewGormDB(cfg)
	if err != nil {
		return nil, err
	}

	// Run auto-migration
	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	// Create default roles
	if err := CreateDefaultRoles(db); err != nil {
		return nil, err
	}

	return db, nil
}