package ledger

type Storage interface {
	UpdateLedger(user string, asset string, amount string) error
}
