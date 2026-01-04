package memory

import (
	"time"

	"github.com/guiferpa/guidger/domain/webhook"
)

type StorageMemory struct {
	ledger        map[string]map[string]int64
	webhookNonces map[string]time.Duration
}

func (s *StorageMemory) UpdateLedger(user string, asset string, amount int64) error {
	s.ledger[user][asset] += amount
	return nil
}

func (s *StorageMemory) UseWebhookSignatureNonce(nonce string) error {
	if _, ok := s.webhookNonces[nonce]; ok {
		return webhook.ErrNonceAlreadyUsed
	}
	return nil
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{ledger: make(map[string]map[string]int64)}
}
