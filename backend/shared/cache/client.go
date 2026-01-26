package cache

// Client represents a Redis client interface
type Client interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl int) error
	Delete(key string) error
	Exists(key string) (bool, error)
}
