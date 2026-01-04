package webhook

import (
	"errors"
	"testing"
	"time"
)

func TestValidateTimestamp(t *testing.T) {
	tolerance := 3 * time.Minute
	timestamp := time.Now().Add(-2 * time.Minute)
	whk := New(NewOptions{SignatureTolerance: tolerance})
	if err := whk.ValidateTimestamp(timestamp.Unix()); err != nil {
		t.Fatalf("expected no error, got %v, timestamp: %s", err, timestamp.Format(time.RFC3339))
	}
}

func TestValidateTimestamp_TooOld(t *testing.T) {
	tolerance := 3 * time.Minute
	timestamp := time.Now().Add(-4 * time.Minute)
	whk := New(NewOptions{SignatureTolerance: tolerance})
	err := whk.ValidateTimestamp(timestamp.Unix())
	if err == nil {
		t.Fatalf("expected error, got nil, timestamp: %s", timestamp.Format(time.RFC3339))
	}
	if !errors.Is(err, ErrTimestampTooOld) {
		t.Fatalf("expected error, got %v, timestamp: %s", err, timestamp.Format(time.RFC3339))
	}
}

func TestValidateTimestamp_InTheFuture(t *testing.T) {
	tolerance := 3 * time.Minute
	timestamp := time.Now().Add(4 * time.Minute)
	whk := New(NewOptions{SignatureTolerance: tolerance})
	err := whk.ValidateTimestamp(timestamp.Unix())
	if err == nil {
		t.Fatalf("expected error, got nil, timestamp: %s", timestamp.Format(time.RFC3339))
	}
	if !errors.Is(err, ErrTimestampInTheFuture) {
		t.Fatalf("expected error, got %v, timestamp: %s", err, timestamp.Format(time.RFC3339))
	}
}

func TestValidateNonce(t *testing.T) {
	whk := New(NewOptions{Storage: NewMockStorage()})
	err := whk.ValidateNonce("nonce")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGenerateCanonicalString(t *testing.T) {
	whk := &whk{}
	canonicalString := whk.generateCanonicalString(1, "nonce", []byte("payload"))
	expected := "1\nnonce\npayload"
	if canonicalString != expected {
		t.Fatalf("expected %s, got %s", expected, canonicalString)
	}
}

func TestValidateSignature(t *testing.T) {
	whk := New(NewOptions{
		Signer:             NewMockSigner(),
		Storage:            NewMockStorage(),
		SignatureTolerance: 1 * time.Minute,
		SignatureSecret:    "secret",
	})
	err := whk.ValidateSignature(time.Now().Unix(), "nonce", []byte("payload"), "signature")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
