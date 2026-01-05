package webhook

import "errors"

// ErrInvalidSignature is returned when a signature validation fails
var ErrInvalidSignature = errors.New("invalid signature")

// Signer provides methods for generating and validating HMAC SHA256 signatures
type Signer interface {
	// GenerateSignature creates an HMAC SHA256 signature for the given parameters
	// The canonical string format is: timestamp + "\n" + nonce + "\n" + payload
	GenerateSignature(timestamp int64, nonce string, payload []byte, secret string) string

	// ValidateSignature compares two signatures using constant-time comparison
	// Returns true if the signatures match, false otherwise
	ValidateSignature(providedSignature, signature string) bool
}

// MockSigner is a mock implementation of Signer for testing purposes
type MockSigner struct {
	GenerateSignatureFunc func(timestamp int64, nonce string, payload []byte, secret string) string
	ValidateSignatureFunc func(providedSignature, signature string) bool
}

func (m *MockSigner) GenerateSignature(timestamp int64, nonce string, payload []byte, secret string) string {
	return m.GenerateSignatureFunc(timestamp, nonce, payload, secret)
}

func (m *MockSigner) ValidateSignature(providedSignature, signature string) bool {
	return m.ValidateSignatureFunc(providedSignature, signature)
}

// NewMockSigner creates a new MockSigner with default implementations
func NewMockSigner() *MockSigner {
	return &MockSigner{
		GenerateSignatureFunc: func(timestamp int64, nonce string, payload []byte, secret string) string {
			return ""
		},
		ValidateSignatureFunc: func(providedSignature, signature string) bool {
			return true
		},
	}
}
