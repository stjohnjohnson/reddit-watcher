package data

import (
	"fmt"
	"strings"

	"github.com/matryer/persist"
)

// Keywords keeps track of the criteria and number of hits
type Keywords map[string]int

// Handler represents all information stored by the application
type Handler struct {
	keyMap   map[string][]int64
	userMap  map[int64]Keywords
	keywords []string
	path     string
}

// Interface is the stats public functions
type Interface interface {
	Get(int64) Keywords
	GetByKeyword(string) []int64
	GetKeywords() []string
	Sync()
	Add(int64, string) error
	Exists(int64, string) bool
	Remove(int64, string) error
	Increment(int64, string) error
}

// Get returns the map of Keywords per user ID
func (ud *Handler) Get(id int64) Keywords {
	_, ok := ud.userMap[id]
	if !ok {
		ud.userMap[id] = make(Keywords)
	}

	return ud.userMap[id]
}

// GetByKeyword returns a list of user IDs for a given keyword
func (ud *Handler) GetByKeyword(keyword string) []int64 {
	ids, ok := ud.keyMap[strings.ToLower(keyword)]
	if !ok {
		ids = []int64{}
	}

	return ids
}

// GetKeywords returns the list of all keywords being searched for
func (ud *Handler) GetKeywords() []string {
	return ud.keywords
}

// Sync updates the temporary variables keyMap and keywords
func (ud *Handler) Sync() {
	keyMap := make(map[string][]int64)
	for id, keys := range ud.userMap {
		for key := range keys {
			keyMap[key] = append(keyMap[key], id)
		}
	}
	ud.keyMap = keyMap

	keywords := make([]string, len(keyMap))
	i := 0
	for key := range keyMap {
		keywords[i] = key
		i++
	}
	ud.keywords = keywords
}

// Add watches a keyword for a given user ID
func (ud *Handler) Add(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)] = 0
	ud.Sync()

	return ud.save()
}

// Exists checks if something is being watched
func (ud *Handler) Exists(id int64, keyword string) bool {
	_, ok := ud.Get(id)[strings.ToLower(keyword)]
	return ok
}

// Remove no longer watches a keyword for a given user ID
func (ud *Handler) Remove(id int64, keyword string) error {
	delete(ud.Get(id), strings.ToLower(keyword))
	ud.Sync()

	return ud.save()
}

// Increment bumps the hit counter on a given keyword and user ID
func (ud *Handler) Increment(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)]++

	return ud.save()
}

// save persists the user data and stats to disk
func (ud *Handler) save() error {
	err := persist.Save(fmt.Sprintf("%s.json", ud.path), ud.userMap)
	if err != nil {
		return fmt.Errorf("save data failed: %v", err)
	}

	return nil
}

// Load recovers the user data and stats from disk
func Load(path string) (*Handler, error) {
	var userMap map[int64]Keywords

	err := persist.Load(fmt.Sprintf("%s.json", path), &userMap)
	if err != nil {
		userMap = make(map[int64]Keywords)
	}

	appData := &Handler{
		userMap: userMap,
		path:    path,
	}
	appData.Sync()

	return appData, err
}
