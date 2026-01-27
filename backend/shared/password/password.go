package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Config holds password hashing configuration
type Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultConfig returns a secure default configuration
func DefaultConfig() *Config {
	return &Config{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,         // 3 iterations
		Parallelism: 2,         // 2 threads
		SaltLength:  16,        // 16 bytes salt
		KeyLength:   32,        // 32 bytes key
	}
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Generate random salt
	salt := make([]byte, config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generate hash
	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// Encode to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	config, salt, hash, err := parseEncodedHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate hash from provided password
	otherHash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// Compare hashes using constant time comparison
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// parseEncodedHash parses an encoded hash string
func parseEncodedHash(encodedHash string) (*Config, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported algorithm")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version")
	}

	config := &Config{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Iterations, &config.Parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.KeyLength = uint32(len(hash))

	return config, salt, hash, nil
}