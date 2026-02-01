package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the payments service
type Config struct {
	// Service Configuration
	ServiceName string
	ServicePort string
	Environment string
	Version     string

	// Paystack Configuration
	PaystackSecretKey    string
	PaystackPublicKey    string
	PaystackWebhookSecret string
	PaystackBaseURL      string

	// MongoDB Configuration
	MongoDBURI      string
	MongoDBDatabase string

	// RabbitMQ Configuration
	RabbitMQURL string

	// Redis Configuration
	RedisURL string

	// JWT Configuration
	JWTSecret string

	// Datadog Configuration
	DatadogAPIKey       string
	DatadogSite         string
	DatadogAgentHost    string
	DatadogTracePort    string
	DatadogEnv          string
	DatadogVersion      string

	// Webhook Configuration
	WebhookTimeoutSeconds int
	WebhookMaxRetries     int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables")
	}

	config := &Config{
		// Service Configuration
		ServiceName: getEnv("DD_SERVICE", "payments-service"),
		ServicePort: getEnv("PAYMENTS_SERVICE_PORT", "8081"),
		Environment: getEnv("PAYMENTS_SERVICE_ENV", "development"),
		Version:     getEnv("DD_VERSION", "1.0.0"),

		// Paystack Configuration
		PaystackSecretKey:     getEnv("PAYSTACK_SECRET_KEY", ""),
		PaystackPublicKey:     getEnv("PAYSTACK_PUBLIC_KEY", ""),
		PaystackWebhookSecret: getEnv("PAYSTACK_WEBHOOK_SECRET", ""),
		PaystackBaseURL:       getEnv("PAYSTACK_BASE_URL", "https://api.paystack.co"),

		// MongoDB Configuration
		MongoDBURI:      getEnv("PAYMENTS_MONGODB_URI", getEnv("MONGODB_URI", "")),
		MongoDBDatabase: getEnv("PAYMENTS_MONGODB_DATABASE", "payments_db"),

		// RabbitMQ Configuration
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),

		// Redis Configuration
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),

		// JWT Configuration
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Datadog Configuration
		DatadogAPIKey:    getEnv("DD_API_KEY", ""),
		DatadogSite:      getEnv("DD_SITE", "us5.datadoghq.com"),
		DatadogAgentHost: getEnv("DD_AGENT_HOST", "localhost"),
		DatadogTracePort: getEnv("DD_TRACE_AGENT_PORT", "8126"),
		DatadogEnv:       getEnv("DD_ENV", "dev"),
		DatadogVersion:   getEnv("DD_VERSION", "1.0.0"),

		// Webhook Configuration
		WebhookTimeoutSeconds: getEnvAsInt("WEBHOOK_TIMEOUT_SECONDS", 30),
		WebhookMaxRetries:     getEnvAsInt("WEBHOOK_MAX_RETRIES", 3),
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.PaystackSecretKey == "" {
		return fmt.Errorf("PAYSTACK_SECRET_KEY is required")
	}
	if c.PaystackPublicKey == "" {
		return fmt.Errorf("PAYSTACK_PUBLIC_KEY is required")
	}
	if c.MongoDBURI == "" {
		return fmt.Errorf("MONGODB_URI is required")
	}
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.JWTSecret == "" {
		log.Printf("Warning: JWT_SECRET is not set")
	}

	return nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as int with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return fallback
}
