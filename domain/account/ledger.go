package account

type Ledger interface {
	Apply(user string, asset string, amount int64) error
}
