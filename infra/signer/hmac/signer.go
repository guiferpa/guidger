package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
)

type SignerHMAC struct{}

func (s *SignerHMAC) GenerateSignature(canonicalString, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonicalString))
	return mac.Sum(nil)
}

func (s *SignerHMAC) ValidateSignature(providedSignature, signature string) bool {
	return hmac.Equal([]byte(providedSignature), []byte(signature))
}

func NewSignerHMAC() *SignerHMAC {
	return &SignerHMAC{}
}
