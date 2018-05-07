package stats

import (
	"testing"
)

func TestLoad(t *testing.T) {
	obj, _ := Load("/tmp/foo-bar")

	data := obj.GetAll()
	if len(data) != 1 {
		t.Errorf("Expected 1 item, got %+v", data)
	}
	if _, ok := data["uptime"]; !ok {
		t.Errorf("Expected uptime, got %+v", data)
	}
}

func TestReadWriteee(t *testing.T) {
	obj, _ := Load("/tmp/foo-bar")

	obj.Increment("foo")
	obj.Increment("bar")
	obj.Increment("bar")

	data := obj.GetAll()
	if len(data) != 3 {
		t.Errorf("Expected 3 item, got %+v", data)
	}

	if data["foo"] != "1 hits" {
		t.Errorf("Expected foo=1, got %v", data["foo"])
	}

	if data["bar"] != "2 hits" {
		t.Errorf("Expected bar=2, got %v", data["bar"])
	}
}
