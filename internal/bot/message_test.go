package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/stjohnjohnson/reddit-watcher/internal/data"
	"github.com/stjohnjohnson/reddit-watcher/internal/matcher"
	"github.com/stjohnjohnson/reddit-watcher/mocks"
)

func TestMessageBad(t *testing.T) {
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
	}

	err := obj.incomingMessage(1, "Foo")

	// @TODO check the logger output
	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
}

func TestMessageInvalid(t *testing.T) {
	var actual string
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = fmt.Sprintf("%d/%s", i, s)
				return nil
			},
		},
	}

	err := obj.incomingMessage(1, "/invalid")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "1/That command doesn't look like anything to me."
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestMessageHelp(t *testing.T) {
	var actual string
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = fmt.Sprintf("%d/%s", i, s)
				return nil
			},
		},
	}

	err := obj.incomingMessage(1, "/help")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "1/Hi, I'm"
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("Expected %q to start with %q", actual, expected)
	}
}

func TestMessageStart(t *testing.T) {
	var actual string
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = fmt.Sprintf("%d/%s", i, s)
				return nil
			},
		},
	}

	err := obj.incomingMessage(1, "/start")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "1/Hi, I'm"
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("Expected %q to start with %q", actual, expected)
	}
}

func TestMessageUnsubscribeFail(t *testing.T) {
	var actual []string
	data := make(map[string]data.Interface)
	data[matcher.Selling] = &mocks.Data{
		MockExists: func(int64, string) bool {
			return true
		},
		MockAdd: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("add/%d/%s", i, s))
			return nil
		},
		MockRemove: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("rm/%d/%s", i, s))
			return fmt.Errorf("Bad")
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = append(actual, fmt.Sprintf("msg/%d/%s", i, s))
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingMessage(1, "/selling foo")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := []string{
		"rm/1/foo",
		"msg/1/I'm no longer watching for <b>selling</b> posts that match <b>foo</b>",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageUnsubscribe(t *testing.T) {
	var actual []string
	data := make(map[string]data.Interface)
	data[matcher.Selling] = &mocks.Data{
		MockExists: func(int64, string) bool {
			return true
		},
		MockAdd: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("add/%d/%s", i, s))
			return nil
		},
		MockRemove: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("rm/%d/%s", i, s))
			return nil
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = append(actual, fmt.Sprintf("msg/%d/%s", i, s))
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingMessage(1, "/selling foo")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := []string{
		"rm/1/foo",
		"msg/1/I'm no longer watching for <b>selling</b> posts that match <b>foo</b>",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageSubscribeAll(t *testing.T) {
	var actual []string
	data := make(map[string]data.Interface)
	data[matcher.Selling] = &mocks.Data{
		MockExists: func(int64, string) bool {
			return false
		},
		MockAdd: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("add/%d/%s", i, s))
			return nil
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = append(actual, fmt.Sprintf("msg/%d/%s", i, s))
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingMessage(1, "/selling")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := []string{
		"add/1/*",
		"msg/1/Okay, I'm going to watch for <b>selling</b> posts that match <b>*</b>",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageSubscribe(t *testing.T) {
	var actual []string
	data := make(map[string]data.Interface)
	data[matcher.Selling] = &mocks.Data{
		MockExists: func(int64, string) bool {
			return false
		},
		MockAdd: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("add/%d/%s", i, s))
			return nil
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = append(actual, fmt.Sprintf("msg/%d/%s", i, s))
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingMessage(1, "/selling foo")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := []string{
		"add/1/foo",
		"msg/1/Okay, I'm going to watch for <b>selling</b> posts that match <b>foo</b>",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageSubscribeFail(t *testing.T) {
	var actual []string
	data := make(map[string]data.Interface)
	data[matcher.Selling] = &mocks.Data{
		MockExists: func(int64, string) bool {
			return false
		},
		MockAdd: func(i int64, s string) error {
			actual = append(actual, fmt.Sprintf("add/%d/%s", i, s))
			return fmt.Errorf("fail")
		},
	}
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = append(actual, fmt.Sprintf("msg/%d/%s", i, s))
				return nil
			},
		},
		data: data,
	}

	err := obj.incomingMessage(1, "/selling foo")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := []string{
		"add/1/foo",
		"msg/1/Okay, I'm going to watch for <b>selling</b> posts that match <b>foo</b>",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageWatchlistEmpty(t *testing.T) {
	d := make(map[string]data.Interface)
	for _, t := range matcher.Types {
		d[t] = &mocks.Data{
			MockGet: func(int64) data.Keywords {
				return make(data.Keywords)
			},
		}
	}

	var actual string
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = fmt.Sprintf("%d/%s", i, s)
				return nil
			},
		},
		data: d,
	}

	err := obj.incomingMessage(1, "/items")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "1/There are no items on your watch list"
	if actual != expected {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageWatchlist(t *testing.T) {
	d := make(map[string]data.Interface)
	for _, t := range matcher.Types {
		d[t] = &mocks.Data{
			MockGet: func(int64) data.Keywords {
				v := make(data.Keywords)
				v["foo"] = 1
				return v
			},
		}
	}

	var actual string
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = fmt.Sprintf("%d/%s", i, s)
				return nil
			},
		},
		data: d,
	}

	err := obj.incomingMessage(1, "/items")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "1/These are your current watch items:\n<b>BUYING:</b>\n - foo <i>(1 hits)</i>\n\n<b>SELLING:</b>\n - foo <i>(1 hits)</i>\n\n<b>VENDOR:</b>\n - foo <i>(1 hits)</i>\n\n<b>ARTISAN:</b>\n - foo <i>(1 hits)</i>\n\n<b>GROUPBUY:</b>\n - foo <i>(1 hits)</i>\n\n<b>INTERESTCHECK:</b>\n - foo <i>(1 hits)</i>\n\n<b>GIVEAWAY:</b>\n - foo <i>(1 hits)</i>\n"
	if actual != expected {
		t.Errorf("Expected %q to equal %q", actual, expected)
	}
}

func TestMessageStats(t *testing.T) {
	var actual string
	obj := &Handler{
		logger: log.New(ioutil.Discard, "", 0),
		chat: &mocks.Chatter{
			MockSendMessage: func(i int64, s string) error {
				actual = fmt.Sprintf("%d/%s", i, s)
				return nil
			},
		},
		stats: &mocks.Stats{
			MockGetAll: func() map[string]string {
				data := make(map[string]string)
				data["foo"] = "bar"
				return data
			},
		},
	}

	err := obj.incomingMessage(1, "/stats")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("Expected nil, got %q", err)
	}
	expected := "1/<b>Interesting Statistics:</b>\n - foo <i>(bar)</i>"
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("Expected %q to start with %q", actual, expected)
	}
}
