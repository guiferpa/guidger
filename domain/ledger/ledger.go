package ledger

import (
	"errors"
)

var (
	ErrInvalidAmount = errors.New("invalid amount format")
)

type Ledger interface {
	Apply(user string, asset string, amount string) error
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

type NewOptions struct {
	Storage Storage
}

func New(opts NewOptions) Ledger {
	return &ldgr{storage: opts.Storage}
}
