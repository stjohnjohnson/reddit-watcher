package data

import (
	"errors"
	"fmt"
	"strings"

	"github.com/matryer/persist"
)

type Keywords map[string]int
type UserData struct {
	keyMap   map[string][]int64
	userMap  map[int64]Keywords
	keywords []string
	path     string
	stats    map[string]int64
}

func (ud *UserData) GetStats() map[string]string {
	stats := make(map[string]string)

	stats["subscribers"] = fmt.Sprintf("%d channels", len(ud.userMap))
	stats["watchlist"] = fmt.Sprintf("%d items", len(ud.keyMap))

	for keyword, hits := range ud.stats {
		stats[keyword] = fmt.Sprintf("%d hits", hits)
	}

	return stats
}

func (ud *UserData) IncrementStat(flag string) error {
	_, ok := ud.stats[flag]
	if !ok {
		ud.stats[flag] = 0
	}

	ud.stats[flag]++

	err := persist.Save(ud.path, ud)
	if err != nil {
		return errors.New(fmt.Sprintf("Save failed: %v", err))
	}

	return nil
}

func (ud *UserData) Get(id int64) Keywords {
	_, ok := ud.userMap[id]
	if !ok {
		ud.userMap[id] = make(Keywords)
	}

	return ud.userMap[id]
}

func (ud *UserData) GetByKeyword(keyword string) []int64 {
	ids, ok := ud.keyMap[strings.ToLower(keyword)]
	if !ok {
		ids = []int64{}
	}

	return ids
}

func (ud *UserData) GetKeywords() []string {
	return ud.keywords
}

func (ud *UserData) Sync() {
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

func (ud *UserData) Add(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)] = 0
	ud.Sync()

	return ud.Save()
}

func (ud *UserData) Remove(id int64, keyword string) error {
	delete(ud.Get(id), strings.ToLower(keyword))
	ud.Sync()

	return ud.Save()
}

func (ud *UserData) Increment(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)] += 1

	return ud.Save()
}

func (ud *UserData) Save() error {
	err := persist.Save(fmt.Sprintf("%s/user-map.json", ud.path), ud.userMap)
	if err != nil {
		return errors.New(fmt.Sprintf("Save UserMap failed: %v", err))
	}

	err = persist.Save(fmt.Sprintf("%s/stats.json", ud.path), ud.stats)
	if err != nil {
		return errors.New(fmt.Sprintf("Save Stats failed: %v", err))
	}

	return nil
}

func Load(path string) (*UserData, error) {
	var userMap map[int64]Keywords
	var stats map[string]int64

	err := persist.Load(fmt.Sprintf("%s/user-map.json", path), &userMap)
	if err != nil {
		userMap = make(map[int64]Keywords)
	}

	err = persist.Load(fmt.Sprintf("%s/stats.json", path), &stats)
	if err != nil {
		stats = make(map[string]int64)
	}

	userData := &UserData{
		userMap: userMap,
		stats:   stats,
		path:    path,
	}
	userData.Sync()

	return userData, err
}
