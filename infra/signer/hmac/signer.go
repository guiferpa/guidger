package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
)

type HMACSigner struct{}

func (s *HMACSigner) GenerateSignature(canonicalString, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonicalString))
	return mac.Sum(nil)
}

func (s *HMACSigner) ValidateSignature(providedSignature, signature string) bool {
	return hmac.Equal([]byte(providedSignature), []byte(signature))
}

func NewHMACSigner() *HMACSigner {
	return &HMACSigner{}
}
