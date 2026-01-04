package memory

import (
	"sync"
	"time"

	"github.com/guiferpa/guidger/domain/ledger"
	"github.com/guiferpa/guidger/domain/webhook"
	"github.com/guiferpa/guidger/mathutil"
)

type MemoryStorage struct {
	mu            sync.RWMutex
	ledger        map[string]map[string]*mathutil.Big
	webhookNonces map[string]time.Time
}

func (s *MemoryStorage) UpdateLedger(user string, asset string, amount string) error {
	// Parse the amount string to mathutil.Big
	amountBig, err := mathutil.NewBig(amount)
	if err != nil {
		return ledger.ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize user ledger if it doesn't exist
	if s.ledger[user] == nil {
		s.ledger[user] = make(map[string]*mathutil.Big)
	}

	// Initialize asset balance if it doesn't exist
	if s.ledger[user][asset] == nil {
		zeroBig, err := mathutil.NewBig("0")
		if err != nil {
			return ledger.ErrInvalidAmount
		}
		s.ledger[user][asset] = zeroBig
	}

	// Add the amount to the current balance
	s.ledger[user][asset] = s.ledger[user][asset].Add(amountBig)

	return nil
}

func (s *MemoryStorage) GetLedgerBalanceByUser(user string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userLedger, exists := s.ledger[user]
	if !exists || len(userLedger) == 0 {
		return make(map[string]string), nil
	}

	balances := make(map[string]string, len(userLedger))
	for asset, amount := range userLedger {
		balances[asset] = amount.String()
	}

	return balances, nil
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
		ledger:        make(map[string]map[string]*mathutil.Big),
		webhookNonces: make(map[string]time.Time),
	}
}
