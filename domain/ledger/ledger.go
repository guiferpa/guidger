package ledger

type Ledger interface {
	Apply(user string, asset string, amount int64) error
}

type ldgr struct {
	storage Storage
}

func (l *ldgr) Apply(user string, asset string, amount int64) error {
	return nil
}

type NewOptions struct {
	Storage Storage
}

func New(opts NewOptions) Ledger {
	return &ldgr{storage: opts.Storage}
}
