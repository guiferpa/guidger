package webhook

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimestampTooOld      = errors.New("timestamp is too old")
	ErrTimestampInTheFuture = errors.New("timestamp is in the future")
	ErrInvalidSignature     = errors.New("invalid signature")
)

type Webhook interface {
	ValidateTimestamp(timestamp int64) error
	ValidateSignature(timestamp int64, nonce string, payload []byte, signature string) error
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

func (w *whk) generateCanonicalString(timestamp int64, nonce string, payload []byte) string {
	return fmt.Sprintf("%d.%s.%s", timestamp, nonce, string(payload))
}

func (w *whk) ValidateSignature(timestamp int64, nonce string, payload []byte, signature string) error {
	if err := w.ValidateTimestamp(timestamp); err != nil {
		return err
	}
	if err := w.ValidateNonce(nonce); err != nil {
		return err
	}
	canonicalString := w.generateCanonicalString(timestamp, nonce, payload)
	generatedSignature := w.signer.GenerateSignature(canonicalString, w.signatureSecret)
	if !w.signer.ValidateSignature(string(generatedSignature), signature) {
		return ErrInvalidSignature
	}
	return nil
}

func (w *whk) ValidateNonce(nonce string) error {
	return w.storage.UseWebhookSignatureNonce(nonce, w.signatureTolerance)
}

type NewOptions struct {
	SignatureTolerance time.Duration
	SignatureSecret    string
	Storage            Storage
	Signer             Signer
}

func New(opts NewOptions) Webhook {
	return &whk{
		signatureTolerance: opts.SignatureTolerance,
		signatureSecret:    opts.SignatureSecret,
		storage:            opts.Storage,
		signer:             opts.Signer,
	}
}
