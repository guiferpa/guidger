package ledger

type Storage interface {
	UpdateLedger(user string, asset string, amount string) error
	GetLedgerBalanceByUser(user string) (map[string]string, error)
}
