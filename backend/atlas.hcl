// Atlas configuration for GoFund database migrations

// Define the environment for development
env "dev" {
  // Database URL for development
  url = "postgres://postgres:postgres@localhost:5432/gofund_dev?sslmode=disable"
  
  // Migration directory
  migration {
    dir = "file://migrations"
  }
  
  // Schema definition
  src = "file://schema.hcl"
}

// Define the environment for production
env "prod" {
  // Database URL from environment variable
  url = env("DATABASE_URL")
  
  // Migration directory
  migration {
    dir = "file://migrations"
  }
  
  // Schema definition
  src = "file://schema.hcl"
}

// Define the environment for testing
env "test" {
  // Database URL for testing
  url = "postgres://postgres:postgres@localhost:5432/gofund_test?sslmode=disable"
  
  // Migration directory
  migration {
    dir = "file://migrations"
  }
  
  // Schema definition
  src = "file://schema.hcl"
}