package webhook

import "errors"

var ErrInvalidSignature = errors.New("invalid signature")

type Signer interface {
	GenerateSignature(timestamp int64, nonce string, payload []byte, secret string) string
	ValidateSignature(providedSignature, signature string) bool
}

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
