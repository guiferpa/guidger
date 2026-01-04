package webhook

import "errors"

var (
	ErrNonceAlreadyUsed = errors.New("nonce already used")
)

type Storage interface {
	UseWebhookSignatureNonce(nonce string) error
}
