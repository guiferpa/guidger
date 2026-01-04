package memory

import (
	"math/big"
	"sync"
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
)

type MemoryStorage struct {
	mu            sync.RWMutex
	ledger        map[string]map[string]*big.Rat
	webhookNonces map[string]time.Time
}

func (s *MemoryStorage) UpdateLedger(user string, asset string, amount string) error {
	// Parse the amount string to big.Rat
	amountRat := new(big.Rat)
	if _, ok := amountRat.SetString(amount); !ok {
		return ledger.ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize user ledger if it doesn't exist
	if s.ledger[user] == nil {
		s.ledger[user] = make(map[string]*big.Rat)
	}

	// Initialize asset balance if it doesn't exist
	if s.ledger[user][asset] == nil {
		s.ledger[user][asset] = new(big.Rat)
	}

	// Add the amount to the current balance
	s.ledger[user][asset].Add(s.ledger[user][asset], amountRat)

	return nil
}

func (s *MemoryStorage) UseWebhookSignatureNonce(nonce string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	v, ok := s.webhookNonces[nonce]
	if !ok {
		s.webhookNonces[nonce] = now.Add(ttl)
		return nil
	}
	if now.After(v) {
		s.webhookNonces[nonce] = now.Add(ttl)
		return nil
	}
	return webhook.ErrNonceAlreadyUsed
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		ledger:        make(map[string]map[string]*big.Rat),
		webhookNonces: make(map[string]time.Time),
	}
}
