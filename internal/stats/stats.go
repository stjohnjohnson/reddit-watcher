package stats

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Handler represents all information stored by the application
type Handler struct {
	startTime time.Time
	lastTime  time.Time
	data      map[string]int64
}

// Interface is the stats public functions
type Interface interface {
	GetAll() map[string]string
	Increment(string)
}

func prettyPrint(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}
	d := time.Since(t)

	mins := math.Trunc(d.Minutes())
	days := math.Trunc(d.Hours() / 24)
	hours := math.Trunc(math.Mod(d.Hours(), 24))

	resp := []string{}
	if days != 0 {
		resp = append(resp, fmt.Sprintf("%vd", days))
	}
	if hours != 0 {
		resp = append(resp, fmt.Sprintf("%vh", hours))
	}
	if mins != 0 {
		resp = append(resp, fmt.Sprintf("%vm", mins))
	}
	if days == 0 && hours == 0 && mins == 0 {
		return "Just Now"
	}

	return strings.Join(resp, " ")
}

// GetAll provides a map of stat name and value
// It additionally returns current uptime
func (ud *Handler) GetAll() map[string]string {
	data := make(map[string]string)

	data["uptime"] = prettyPrint(ud.startTime)
	data["last post"] = prettyPrint(ud.lastTime)

	for keyword, hits := range ud.data {
		data[keyword] = fmt.Sprintf("%d hits", hits)
	}

	return data
}

// Increment keeps track of the number of occurrences per flag
func (ud *Handler) Increment(flag string) {
	_, ok := ud.data[flag]
	if !ok {
		ud.data[flag] = 0
	}

	ud.data[flag]++
	ud.lastTime = time.Now()
}

// New creates a new stats object
func New() *Handler {
	handler := &Handler{
		startTime: time.Now(),
		data:      make(map[string]int64),
	}

	return handler
}
