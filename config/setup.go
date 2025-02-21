package config

import (
	"os"
	"sync"
	"github.com/joho/godotenv"
)

// Config represents the configuration singleton
type Config struct {
	initialized bool
}

var (
	instance *Config
	once     sync.Once
)

// GetInstance returns the singleton instance of Config
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.initialize()
	})
	return instance
}

func (c *Config) initialize() {
	if c.initialized {
		return
	}
	
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		// Handle error if .env file doesn't exist or can't be loaded
		// For now, we'll just print the error
		println("Error loading .env file:", err)
	}

	c.initialized = true
}

// GetEnv is a helper method to get environment variables
func (c *Config) GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

    
        