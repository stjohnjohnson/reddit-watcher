package mocks

// Stats is mocked
type Stats struct {
	MockGetAll    func() map[string]string
	MockIncrement func()
}

// GetAll is mocked
func (s *Stats) GetAll() map[string]string {
	if s.MockGetAll != nil {
		return s.MockGetAll()
	}
	return nil
}

// Increment is mocked
func (s *Stats) Increment(string) {
	if s.MockIncrement != nil {
		s.MockIncrement()
	}
}
