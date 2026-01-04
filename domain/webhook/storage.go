package webhook

import (
	"errors"
	"time"
)

var (
	ErrNonceAlreadyUsed = errors.New("nonce already used")
)

type Storage interface {
	UseWebhookSignatureNonce(nonce string, ttl time.Duration) error
}

type MockStorage struct {
	UseWebhookSignatureNonceFunc func(nonce string, ttl time.Duration) error
}

func (m *MockStorage) UseWebhookSignatureNonce(nonce string, ttl time.Duration) error {
	return m.UseWebhookSignatureNonceFunc(nonce, ttl)
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		UseWebhookSignatureNonceFunc: func(nonce string, ttl time.Duration) error {
			return nil
		},
	}
}
