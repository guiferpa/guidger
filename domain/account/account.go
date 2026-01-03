package account

type Account interface {
	Apply(user string, asset string, amount int64) error
}

type acc struct {
	ledger Ledger
}

func (a *acc) Apply(user string, asset string, amount int64) error {
	return nil
}

func New(ledger Ledger) Account {
	return &acc{ledger}
}
