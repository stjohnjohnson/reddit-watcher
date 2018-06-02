package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/stjohnjohnson/reddit-watcher/internal/data"
	"github.com/stjohnjohnson/reddit-watcher/internal/matcher"
	"github.com/stjohnjohnson/reddit-watcher/mocks"
	"github.com/turnage/graw/reddit"
)

func TestBadPost(t *testing.T) {
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
	}

	expected := fmt.Errorf("unable to parse title: not parsable: ")
	err := obj.incomingPost(&reddit.Post{})

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Expected %q, got %q", expected, err)
	}
}

func TestBadType(t *testing.T) {
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		data:   make(map[string]data.Interface),
	}

	expected := fmt.Errorf("unknown type: buying")
	err := obj.incomingPost(&reddit.Post{
		Title: "[US-CA] [H] Money [W] Tada68",
	})

	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Expected %q, got %q", expected, err)
	}
}

func TestNoMatch(t *testing.T) {
	data := make(map[string]data.Interface)
	data[matcher.Buying] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"banana"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				t.Errorf("Unexpected call to SendMessage %d, %s", i, s)
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title: "[US-CA] [H] Money [W] Tada68",
	})

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
}

func TestBadRegion(t *testing.T) {
	data := make(map[string]data.Interface)
	data[matcher.Buying] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"tada68"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				t.Errorf("Unexpected call to SendMessage %d, %s", i, s)
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title: "[EU-GB] [H] Money [W] Tada68",
	})

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
}

func TestBadMessage(t *testing.T) {
	var actual string
	data := make(map[string]data.Interface)
	data[matcher.Buying] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"tada68"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = s
				return fmt.Errorf("failed to chat")
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title:     "[US-CA] [H] Money [W] Tada68",
		Permalink: "/r/foo",
		URL:       "https://r.com/r/foobar",
	})

	// @TODO check the logger output
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "[US-CA] [H] Money [W] <b>Tada68</b> [<a href=\"https://r.com/r/foobar\">web</a>] [<a href=\"https://git.io/vhZZN#/r/foo\">app</a>] <i>(matched buying tada68)</i>"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestBadHitIncrement(t *testing.T) {
	var actual string
	data := make(map[string]data.Interface)
	data[matcher.Buying] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"tada68"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
		MockIncrement: func(int64, string) error {
			return fmt.Errorf("increment failed")
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = s
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title:     "[US-CA] [H] Money [W] Tada68",
		Permalink: "/r/foo",
		URL:       "https://r.com/r/foobar",
	})

	// @TODO check the logger output
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "[US-CA] [H] Money [W] <b>Tada68</b> [<a href=\"https://r.com/r/foobar\">web</a>] [<a href=\"https://git.io/vhZZN#/r/foo\">app</a>] <i>(matched buying tada68)</i>"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestHit(t *testing.T) {
	var actual string
	data := make(map[string]data.Interface)
	data[matcher.Buying] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"tada68"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
		MockIncrement: func(int64, string) error {
			return nil
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = s
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title:     "[US-CA] [H] Money [W] Tada68",
		Permalink: "/r/foo",
		URL:       "https://r.com/r/foobar",
	})

	// @TODO check the logger output
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "[US-CA] [H] Money [W] <b>Tada68</b> [<a href=\"https://r.com/r/foobar\">web</a>] [<a href=\"https://git.io/vhZZN#/r/foo\">app</a>] <i>(matched buying tada68)</i>"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestHitStar(t *testing.T) {
	var actual string
	data := make(map[string]data.Interface)
	data[matcher.Buying] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"*"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
		MockIncrement: func(int64, string) error {
			return nil
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = s
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title:     "[US-CA] [H] Money** [W] Tada68",
		Permalink: "/r/foo",
		URL:       "https://r.com/r/foobar",
	})

	// @TODO check the logger output
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "<b>[US-CA]</b> <b>[H]</b> Money** <b>[W]</b> Tada68 [<a href=\"https://r.com/r/foobar\">web</a>] [<a href=\"https://git.io/vhZZN#/r/foo\">app</a>] <i>(matched buying *)</i>"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestHitStarVendor(t *testing.T) {
	var actual string
	data := make(map[string]data.Interface)
	data[matcher.Vendor] = &mocks.Data{
		MockGetKeywords: func() []string {
			return []string{"*"}
		},
		MockGetByKeyword: func(s string) []int64 {
			return []int64{1}
		},
		MockIncrement: func(int64, string) error {
			return nil
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		stats:  &mocks.Stats{},
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = s
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingPost(&reddit.Post{
		Title:     "[vendor] ** Tada68",
		Permalink: "/r/foo",
		URL:       "https://r.com/r/foobar",
	})

	// @TODO check the logger output
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "<b>[vendor]</b> ** Tada68 [<a href=\"https://r.com/r/foobar\">web</a>] [<a href=\"https://git.io/vhZZN#/r/foo\">app</a>] <i>(matched vendor *)</i>"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}
