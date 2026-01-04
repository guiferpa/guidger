package webhook

import (
	"errors"
	"time"
)

var (
	ErrNonceAlreadyUsed = errors.New("nonce already used")
)

type Storage interface {
	UseWebhookSignatureNonce(nonce string, ttl time.Time) error
}
