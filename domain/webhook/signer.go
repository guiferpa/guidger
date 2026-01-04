package webhook

type Signer interface {
	GenerateSignature(canonicalString, secret string) []byte
	ValidateSignature(providedSignature, signature string) bool
}

type MockSigner struct {
	GenerateSignatureFunc func(canonicalString, secret string) []byte
	ValidateSignatureFunc func(providedSignature, signature string) bool
}

func (m *MockSigner) GenerateSignature(canonicalString, secret string) []byte {
	return m.GenerateSignatureFunc(canonicalString, secret)
}

func (m *MockSigner) ValidateSignature(providedSignature, signature string) bool {
	return m.ValidateSignatureFunc(providedSignature, signature)
}

func NewMockSigner() *MockSigner {
	return &MockSigner{
		GenerateSignatureFunc: func(canonicalString, secret string) []byte {
			return []byte("")
		},
		ValidateSignatureFunc: func(providedSignature, signature string) bool {
			return true
		},
	}
}
