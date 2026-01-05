package ledger

import (
	"errors"
)

var (
	// ErrInvalidAmount is returned when an amount string cannot be parsed as a valid decimal number
	ErrInvalidAmount = errors.New("invalid amount format")
)

// Ledger handles balance updates and queries for users
type Ledger interface {
	// Apply updates a user's balance by adding the specified amount to the given asset
	// The amount must be a valid decimal string (e.g., "1.5", "0.0001")
	Apply(user string, asset string, amount string) error

	// GetBalanceByUser returns all balances for a specific user
	// Returns a map of asset names to balance strings (decimal representation)
	GetBalanceByUser(user string) (map[string]string, error)
}

type ldgr struct {
	storage Storage
}

func (l *ldgr) Apply(user string, asset string, amount string) error {
	return l.storage.UpdateLedger(user, asset, amount)
}

func (l *ldgr) GetBalanceByUser(user string) (map[string]string, error) {
	return l.storage.GetLedgerBalanceByUser(user)
}

// NewOptions contains configuration for creating a new Ledger instance
type NewOptions struct {
	Storage Storage
}

// New creates a new Ledger instance with the provided storage
func New(opts NewOptions) Ledger {
	return &ldgr{storage: opts.Storage}
}
