package memory

import (
	"time"

	"github.com/guiferpa/guidger/domain/webhook"
)

type StorageMemory struct {
	ledger        map[string]map[string]int64
	webhookNonces map[string]time.Time
}

func (s *StorageMemory) UpdateLedger(user string, asset string, amount int64) error {
	s.ledger[user][asset] += amount
	return nil
}

func (s *StorageMemory) UseWebhookSignatureNonce(nonce string, ttl time.Duration) error {
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

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		ledger:        make(map[string]map[string]int64),
		webhookNonces: make(map[string]time.Time),
	}
}
