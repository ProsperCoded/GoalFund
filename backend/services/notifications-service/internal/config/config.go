package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the notifications service
type Config struct {
	// Server
	Port string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// RabbitMQ
	RabbitMQURL      string
	RabbitMQExchange string
	RabbitMQQueue    string

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string

	// Datadog
	DDService string
	DDEnv     string
	DDVersion string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		// Server
		Port: getEnv("PORT", "8085"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "notifications_db"),

		// RabbitMQ
		RabbitMQURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQExchange: getEnv("RABBITMQ_EXCHANGE", "gofund_events"),
		RabbitMQQueue:    getEnv("RABBITMQ_QUEUE", "notifications_queue"),

		// Email
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@gofund.com"),
		SMTPFromName: getEnv("SMTP_FROM_NAME", "GoFund"),

		// Datadog
		DDService: getEnv("DD_SERVICE", "notifications-service"),
		DDEnv:     getEnv("DD_ENV", "dev"),
		DDVersion: getEnv("DD_VERSION", "1.0.0"),
	}
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
