package ledger

// Storage provides methods for persisting and retrieving ledger data
type Storage interface {
	// UpdateLedger updates a user's balance by adding the specified amount to the given asset
	// The amount must be a valid decimal string (e.g., "1.5", "0.0001")
	// Returns ErrInvalidAmount if the amount format is invalid
	UpdateLedger(user string, asset string, amount string) error

	// GetLedgerBalanceByUser retrieves all balances for a specific user
	// Returns a map of asset names to balance strings (decimal representation)
	// Returns an empty map if the user has no balances
	GetLedgerBalanceByUser(user string) (map[string]string, error)
}
