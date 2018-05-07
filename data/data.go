package data

import (
	"errors"
	"fmt"
	"strings"

	"github.com/matryer/persist"
)

// Keywords keeps track of the criteria and number of hits
type Keywords map[string]int

// AppData represents all information stored by the application
type AppData struct {
	keyMap   map[string][]int64
	userMap  map[int64]Keywords
	keywords []string
	path     string
}

// Hits provides a map of stat name and value
// It captures:
//  - # of channels/users asking for reply
//  - # of keywords it is looking for
func (ud *AppData) Hits() map[string]string {
	stats := make(map[string]string)

	stats["subscribers"] = fmt.Sprintf("%d channels", len(ud.userMap))
	stats["watchlist"] = fmt.Sprintf("%d items", len(ud.keyMap))

	return stats
}

// Get returns the map of Keywords per user ID
func (ud *AppData) Get(id int64) Keywords {
	_, ok := ud.userMap[id]
	if !ok {
		ud.userMap[id] = make(Keywords)
	}

	return ud.userMap[id]
}

// GetByKeyword returns a list of user IDs for a given keyword
func (ud *AppData) GetByKeyword(keyword string) []int64 {
	ids, ok := ud.keyMap[strings.ToLower(keyword)]
	if !ok {
		ids = []int64{}
	}

	return ids
}

// GetKeywords returns the list of all keywords being searched for
func (ud *AppData) GetKeywords() []string {
	return ud.keywords
}

// Sync updates the temporary variables keyMap and keywords
func (ud *AppData) Sync() {
	keyMap := make(map[string][]int64)
	for id, keys := range ud.userMap {
		for key, _ := range keys {
			keyMap[key] = append(keyMap[key], id)
		}
	}
	ud.keyMap = keyMap

	keywords := make([]string, len(keyMap))
	i := 0
	for key, _ := range keyMap {
		keywords[i] = key
		i++
	}
	ud.keywords = keywords
}

// Add watches a keyword for a given user ID
func (ud *AppData) Add(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)] = 0
	ud.Sync()

	return ud.Save()
}

// Remove no longer watches a keyword for a given user ID
func (ud *AppData) Remove(id int64, keyword string) error {
	delete(ud.Get(id), strings.ToLower(keyword))
	ud.Sync()

	return ud.Save()
}

// Increment bumps the hit counter on a given keyword and user ID
func (ud *AppData) Increment(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)] += 1

	return ud.Save()
}

// Save persists the user data and stats to disk
func (ud *AppData) Save() error {
	err := persist.Save(fmt.Sprintf("%s/user-map.json", ud.path), ud.userMap)
	if err != nil {
		return errors.New(fmt.Sprintf("Save UserMap failed: %v", err))
	}

	return nil
}

// Load recovers the user data and stats from disk
func Load(path string) (*AppData, error) {
	var userMap map[int64]Keywords

	err := persist.Load(fmt.Sprintf("%s/user-map.json", path), &userMap)
	if err != nil {
		userMap = make(map[int64]Keywords)
	}

	appData := &AppData{
		userMap: userMap,
		path:    path,
	}
	appData.Sync()

	return appData, err
}
