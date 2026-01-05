package webhook

import (
	"errors"
	"time"
)

var (
	// ErrTimestampTooOld is returned when a request timestamp is beyond the allowed tolerance
	ErrTimestampTooOld = errors.New("timestamp is too old")

	// ErrTimestampInTheFuture is returned when a request timestamp is in the future
	ErrTimestampInTheFuture = errors.New("timestamp is in the future")
)

// Webhook validates signed webhook requests
type Webhook interface {
	// ValidateTimestamp checks if the timestamp is within the allowed tolerance window
	ValidateTimestamp(timestamp int64) error

	// ValidateSignature validates the HMAC SHA256 signature of a webhook request
	// It also validates timestamp and nonce as part of the signature validation
	ValidateSignature(timestamp int64, nonce string, payload []byte, signature string) error

	// ValidateNonce checks if a nonce has been used before (replay attack prevention)
	ValidateNonce(nonce string) error
}

type whk struct {
	signatureTolerance time.Duration
	signatureSecret    string
	storage            Storage
	signer             Signer
}

func (w *whk) ValidateTimestamp(timestamp int64) error {
	unix := time.Unix(timestamp, 0)
	now := time.Now()
	diff := now.Sub(unix)
	if diff < 0 {
		return ErrTimestampInTheFuture
	}
	if diff > w.signatureTolerance {
		return ErrTimestampTooOld
	}
	return nil
}

func (w *whk) ValidateSignature(timestamp int64, nonce string, payload []byte, signature string) error {
	if err := w.ValidateTimestamp(timestamp); err != nil {
		return err
	}
	if err := w.ValidateNonce(nonce); err != nil {
		return err
	}
	generatedSignature := w.signer.GenerateSignature(timestamp, nonce, payload, w.signatureSecret)
	if !w.signer.ValidateSignature(generatedSignature, signature) {
		return ErrInvalidSignature
	}
	return nil
}

func (w *whk) ValidateNonce(nonce string) error {
	return w.storage.UseWebhookSignatureNonce(nonce, w.signatureTolerance)
}

// NewOptions contains configuration for creating a new Webhook validator
type NewOptions struct {
	// SignatureTolerance is the maximum age of a request timestamp (e.g., 5 * time.Minute)
	SignatureTolerance time.Duration

	// SignatureSecret is the secret key used for HMAC SHA256 signature validation
	SignatureSecret string

	// Storage is used for nonce tracking to prevent replay attacks
	Storage Storage

	// Signer is used to generate and validate HMAC signatures
	Signer Signer
}

// New creates a new Webhook validator with the provided configuration
func New(opts NewOptions) Webhook {
	return &whk{
		signatureTolerance: opts.SignatureTolerance,
		signatureSecret:    opts.SignatureSecret,
		storage:            opts.Storage,
		signer:             opts.Signer,
	}
}
