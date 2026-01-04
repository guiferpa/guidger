package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
)

func TestUseWebhookSignatureNonce_NotAlreadyUsed(t *testing.T) {
	sm := NewMemoryStorage()
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
	sm := NewMemoryStorage()
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

func TestUpdateLedger(t *testing.T) {
	sm := NewMemoryStorage()
	err := sm.UpdateLedger("user", "asset", "1.0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateLedger_InvalidAmount(t *testing.T) {
	sm := NewMemoryStorage()
	err := sm.UpdateLedger("user", "asset", "invalid")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ledger.ErrInvalidAmount) {
		t.Fatalf("expected error, got %v, want %v", err, ledger.ErrInvalidAmount)
	}
}

func TestGetLedgerBalanceByUser_ETH(t *testing.T) {
	sm := NewMemoryStorage()
	err := sm.UpdateLedger("user", "asset", "1.0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	balances, err := sm.GetLedgerBalanceByUser("user")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(balances) != 1 {
		t.Fatalf("expected 1 balance, got %d", len(balances))
	}
	if balances["asset"] != "1.0" {
		t.Fatalf("expected balance of 1.0, got %s", balances["asset"])
	}
}
