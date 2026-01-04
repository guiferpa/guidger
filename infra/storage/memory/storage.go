package memory

type StorageMemory struct {
	state map[string]map[string]int64
}

func (s *StorageMemory) Apply(user string, asset string, amount int64) error {
	s.state[user][asset] += amount
	return nil
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{state: make(map[string]map[string]int64)}
}
