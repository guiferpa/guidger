package ledger

type Storage interface {
	UpdateLedger(user string, asset string, amount int64) error
}
