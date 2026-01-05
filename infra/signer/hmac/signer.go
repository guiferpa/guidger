package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type HMACSigner struct{}

func (s *HMACSigner) GenerateCanonicalString(timestamp int64, nonce string, payload []byte) string {
	return fmt.Sprintf("%d\n%s\n%s", timestamp, nonce, string(payload))
}

func (s *HMACSigner) GenerateSignature(timestamp int64, nonce string, payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(s.GenerateCanonicalString(timestamp, nonce, payload)))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *HMACSigner) ValidateSignature(providedSignature, signature string) bool {
	return hmac.Equal([]byte(providedSignature), []byte(signature))
}

func NewHMACSigner() *HMACSigner {
	return &HMACSigner{}
}
