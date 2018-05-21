package mocks

import "github.com/stjohnjohnson/reddit-watcher/chatter"

// Chatter is mocked
type Chatter struct {
	MockStart       func() (chatter.Channel, error)
	MockSendMessage func(int64, string) error
}

// GetAll is mocked
func (m *Chatter) GetAll() map[string]string {
	return nil
}

// Start is mocked
func (m *Chatter) Start() (chatter.Channel, error) {
	if m.MockStart != nil {
		return m.MockStart()
	}
	return nil, nil
}

// SendMessage is mocked
func (m *Chatter) SendMessage(i int64, s string) error {
	if m.MockSendMessage != nil {
		return m.MockSendMessage(i, s)
	}
	return nil
}

// Save is mocked
func (m *Chatter) Save() error {
	return nil
}
