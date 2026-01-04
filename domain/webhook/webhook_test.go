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

func TestValidateSignature(t *testing.T) {
	whk := New(NewOptions{})
	err := whk.ValidateSignature("signature")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateNonce(t *testing.T) {
	whk := New(NewOptions{})
	err := whk.ValidateNonce()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
