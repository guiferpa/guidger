package memory

type LedgerMemory struct {
	state map[string]map[string]int64
}

func (l *LedgerMemory) Apply(user string, asset string, amount int64) error {
	l.state[user][asset] += amount
	return nil
}

func NewLedgerMemory() *LedgerMemory {
	return &LedgerMemory{state: make(map[string]map[string]int64)}
}
