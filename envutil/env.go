package envutil

import (
	"errors"
	"os"
	"strconv"
)

var (
	// ErrRequiredEnvironmentVariable is returned when a required environment variable is not set
	ErrRequiredEnvironmentVariable = errors.New("required environment variable is not set")

	// ErrFailedToParseInteger is returned when an environment variable cannot be parsed as an integer
	ErrFailedToParseInteger = errors.New("failed to parse integer")
)

// GetString retrieves an environment variable as a string
// Returns the default value if the variable is not set or empty
func GetString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetStringRequired retrieves an environment variable as a string
// Returns an error if the variable is not set or empty
func GetStringRequired(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", ErrRequiredEnvironmentVariable
	}
	return value, nil
}

// GetInt retrieves an environment variable as an integer
// Returns the default value if the variable is not set, empty, or cannot be parsed
// Returns an error if the value cannot be parsed as an integer
func GetInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, ErrFailedToParseInteger
	}
	return parsed, nil
}
