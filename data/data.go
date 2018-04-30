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

	err := persist.Save(ud.path, ud.userMap)
	if err != nil {
		return errors.New(fmt.Sprintf("Save failed: %v", err))
	}

	return nil
}

func (ud *UserData) Remove(id int64, keyword string) error {
	delete(ud.Get(id), strings.ToLower(keyword))
	ud.Sync()

	err := persist.Save(ud.path, ud.userMap)
	if err != nil {
		return errors.New(fmt.Sprintf("Save failed: %v", err))
	}

	return nil
}

func (ud *UserData) Increment(id int64, keyword string) error {
	ud.Get(id)[strings.ToLower(keyword)] += 1

	err := persist.Save(ud.path, ud.userMap)
	if err != nil {
		return errors.New(fmt.Sprintf("Save failed: %v", err))
	}

	return nil
}

func Load(path string) (*UserData, error) {
	var userMap map[int64]Keywords

	err := persist.Load(path, &userMap)
	if err != nil {
		userMap = make(map[int64]Keywords)
	}

	obj := &UserData{
		userMap: userMap,
		path:    path,
	}
	obj.Sync()

	return obj, err
}
