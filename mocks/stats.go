package mocks

// Stats is mocked
type Stats struct {
	MockGetAll    func() map[string]string
	MockIncrement func() error
}

// GetAll is mocked
func (s *Stats) GetAll() map[string]string {
	if s.MockGetAll != nil {
		return s.MockGetAll()
	}
	return nil
}

// Increment is mocked
func (s *Stats) Increment(string) error {
	if s.MockIncrement != nil {
		return s.MockIncrement()
	}
	return nil
}
