package stats

import (
	"regexp"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	obj := New()

	data := obj.GetAll()
	if len(data) != 2 {
		t.Errorf("Expected 2 item, got %+v", data)
	}
	if _, ok := data["uptime"]; !ok {
		t.Errorf("Expected uptime, got %+v", data)
	}
	if _, ok := data["last post"]; !ok {
		t.Errorf("Expected last post, got %+v", data)
	}
}

func TestLongDate(t *testing.T) {
	obj := &Handler{
		lastTime: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
	}

	data := obj.GetAll()
	if ok, _ := regexp.MatchString(`\d+d \d+h \d+m`, data["last post"]); !ok {
		t.Errorf("Expected long last post, got %+v", data["last post"])
	}
}

func TestReadWrite(t *testing.T) {
	obj := New()

	obj.Increment("foo")
	obj.Increment("bar")
	obj.Increment("bar")

	data := obj.GetAll()
	if len(data) != 4 {
		t.Errorf("Expected 4 item, got %+v", data)
	}

	if data["foo"] != "1 hits" {
		t.Errorf("Expected foo=1, got %v", data["foo"])
	}

	if data["bar"] != "2 hits" {
		t.Errorf("Expected bar=2, got %v", data["bar"])
	}
}
