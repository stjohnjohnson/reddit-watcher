package stats

import (
	"errors"
	"fmt"
	"time"

	"github.com/matryer/persist"
)

// Handler represents all information stored by the application
type Handler struct {
	startTime time.Time
	path      string
	data      map[string]int64
}

// GetAll provides a map of stat name and value
// It additionally returns current uptime
func (ud *Handler) GetAll() map[string]string {
	data := make(map[string]string)

	data["uptime"] = time.Since(ud.startTime).String()

	for keyword, hits := range ud.data {
		data[keyword] = fmt.Sprintf("%d hits", hits)
	}

	return data
}

// Increment keeps track of the number of occurrences per flag
func (ud *Handler) Increment(flag string) error {
	_, ok := ud.data[flag]
	if !ok {
		ud.data[flag] = 0
	}

	ud.data[flag]++

	err := persist.Save(ud.path, ud)
	if err != nil {
		return errors.New(fmt.Sprintf("Save failed: %v", err))
	}

	return nil
}

// Save persists the user data and stats to disk
func (ud *Handler) Save() error {
	err := persist.Save(fmt.Sprintf("%s/stats.json", ud.path), ud.data)
	if err != nil {
		return errors.New(fmt.Sprintf("Save Stats failed: %v", err))
	}

	return nil
}

// Load recovers the user data and stats from disk
func Load(path string) (*Handler, error) {
	var data map[string]int64

	err := persist.Load(fmt.Sprintf("%s/stats.json", path), &data)
	if err != nil {
		data = make(map[string]int64)
	}

	handler := &Handler{
		startTime: time.Now(),
		data:      data,
		path:      path,
	}

	return handler, err
}
