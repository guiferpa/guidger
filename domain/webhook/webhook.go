package webhook

import (
	"errors"
	"time"
)

var (
	ErrTimestampTooOld      = errors.New("timestamp is too old")
	ErrTimestampInTheFuture = errors.New("timestamp is in the future")
)

type Webhook interface {
	ValidateTimestamp(timestamp int64) error
	ValidateSignature(signature string) error
	ValidateNonce() error
}

type whk struct {
	signatureTolerance time.Duration
	storage            Storage
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

func (w *whk) ValidateSignature(signature string) error {
	return nil
}

func (w *whk) ValidateNonce() error {
	return nil
}

type NewOptions struct {
	SignatureTolerance time.Duration
	Storage            Storage
}

func New(opts NewOptions) Webhook {
	return &whk{signatureTolerance: opts.SignatureTolerance, storage: opts.Storage}
}
