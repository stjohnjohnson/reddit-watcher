package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"

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
