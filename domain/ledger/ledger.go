package ledger

import (
	"errors"
	"math/big"
	"strings"
)

var (
	ErrInvalidAmount = errors.New("invalid amount format")
)

type Ledger interface {
	Apply(user string, asset string, amount string) error
}

type ldgr struct {
	storage Storage
}

func (l *ldgr) Apply(user string, asset string, amount string) error {
	// Validate and parse the decimal amount
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return ErrInvalidAmount
	}

	// Parse as big.Rat for exact decimal precision
	rat := new(big.Rat)
	if _, ok := rat.SetString(amount); !ok {
		return ErrInvalidAmount
	}

	return l.storage.UpdateLedger(user, asset, amount)
}

type NewOptions struct {
	Storage Storage
}

func New(opts NewOptions) Ledger {
	return &ldgr{storage: opts.Storage}
}
