package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/guiferpa/guidger/domain/webhook"
)

func TestUseWebhookSignatureNonce_NotAlreadyUsed(t *testing.T) {
	sm := NewStorageMemory()
	err := sm.UseWebhookSignatureNonce("nonce", 1*time.Second)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	time.Sleep(2 * time.Second)
	err = sm.UseWebhookSignatureNonce("nonce", 1*time.Minute)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUseWebhookSignatureNonce_AlreadyUsed(t *testing.T) {
	sm := NewStorageMemory()
	err := sm.UseWebhookSignatureNonce("nonce", 1*time.Minute)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = sm.UseWebhookSignatureNonce("nonce", 1*time.Minute)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, webhook.ErrNonceAlreadyUsed) {
		t.Fatalf("expected error, got %v, want %v", err, webhook.ErrNonceAlreadyUsed)
	}
}
