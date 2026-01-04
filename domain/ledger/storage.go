package ledger

type Storage interface {
	Apply(user string, asset string, amount int64) error
}
