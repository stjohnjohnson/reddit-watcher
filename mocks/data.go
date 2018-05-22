package mocks

import "github.com/stjohnjohnson/reddit-watcher/internal/data"

// Data is mocked
type Data struct {
	MockGet          func(int64) data.Keywords
	MockGetByKeyword func(string) []int64
	MockGetKeywords  func() []string
	MockSync         func()
	MockAdd          func(int64, string) error
	MockExists       func(int64, string) bool
	MockRemove       func(int64, string) error
	MockIncrement    func(int64, string) error
}

// Get is mocked
func (m *Data) Get(i int64) data.Keywords {
	if m.MockGet != nil {
		return m.MockGet(i)
	}
	return nil
}

// GetByKeyword is mocked
func (m *Data) GetByKeyword(s string) []int64 {
	if m.MockGetByKeyword != nil {
		return m.MockGetByKeyword(s)
	}
	return nil
}

// GetKeywords is mocked
func (m *Data) GetKeywords() []string {
	if m.MockGetKeywords != nil {
		return m.MockGetKeywords()
	}
	return nil
}

// Sync is mocked
func (m *Data) Sync() {
	if m.MockSync != nil {
		m.MockSync()
	}
}

// Add is mocked
func (m *Data) Add(i int64, s string) error {
	if m.MockAdd != nil {
		return m.MockAdd(i, s)
	}
	return nil
}

// Exists is mocked
func (m *Data) Exists(i int64, s string) bool {
	if m.MockExists != nil {
		return m.MockExists(i, s)
	}
	return false
}

// Remove is mocked
func (m *Data) Remove(i int64, s string) error {
	if m.MockRemove != nil {
		return m.MockRemove(i, s)
	}
	return nil
}

// Increment is mocked
func (m *Data) Increment(i int64, s string) error {
	if m.MockIncrement != nil {
		return m.MockIncrement(i, s)
	}
	return nil
}
