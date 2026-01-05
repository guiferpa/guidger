package webhook

import (
	"errors"
	"time"
)

var (
	// ErrNonceAlreadyUsed is returned when a nonce has already been used (replay attack detected)
	ErrNonceAlreadyUsed = errors.New("nonce already used")
)

// Storage provides methods for tracking webhook nonces to prevent replay attacks
type Storage interface {
	// UseWebhookSignatureNonce marks a nonce as used and validates it hasn't been used before
	// The nonce will be stored for the duration specified by ttl
	// Returns ErrNonceAlreadyUsed if the nonce has already been used
	UseWebhookSignatureNonce(nonce string, ttl time.Duration) error
}

// MockStorage is a mock implementation of Storage for testing purposes
type MockStorage struct {
	UseWebhookSignatureNonceFunc func(nonce string, ttl time.Duration) error
}

func (m *MockStorage) UseWebhookSignatureNonce(nonce string, ttl time.Duration) error {
	return m.UseWebhookSignatureNonceFunc(nonce, ttl)
}

// NewMockStorage creates a new MockStorage with default implementations
func NewMockStorage() *MockStorage {
	return &MockStorage{
		UseWebhookSignatureNonceFunc: func(nonce string, ttl time.Duration) error {
			return nil
		},
	}
}
