package data

import (
	"testing"
)

func TestSave(t *testing.T) {
	obj, _ := Load("/tmp/foo")
	err := obj.Add(1, "foo")

	if err != nil {
		t.Errorf("Expected no error, got %+v", err)
	}
}

func TestExists(t *testing.T) {
	obj, _ := Load("/tmp/foo/bar")

	exists := obj.Exists(1, "foo")
	if exists {
		t.Errorf("Expected not to exist, got %+v", exists)
	}

	obj.Add(1, "foo")
	exists = obj.Exists(1, "foo")
	if !exists {
		t.Errorf("Expected to exist, got %+v", exists)
	}

	obj.Remove(1, "foo")
	exists = obj.Exists(1, "foo")
	if exists {
		t.Errorf("Expected not to exist, got %+v", exists)
	}
}

func TestKeywords(t *testing.T) {
	obj, _ := Load("/tmp/foo/bar")

	obj.Add(1, "foo")
	obj.Add(2, "foo")
	obj.Add(2, "bar")
	obj.Increment(1, "foo")
	obj.Increment(1, "foo")

	keywordMap := obj.Get(1)
	if keywordMap["foo"] != 2 {
		t.Errorf("Expected keywordMap to be 2, got %+v", keywordMap)
	}

	ids := obj.GetByKeyword("foo")
	if len(ids) != 2 {
		t.Errorf("Expected ids to be 2, got %+v", ids)
	}

	ids = obj.GetByKeyword("bar")
	if len(ids) != 1 {
		t.Errorf("Expected ids to be 1, got %+v", ids)
	}

	ids = obj.GetByKeyword("baz")
	if len(ids) != 0 {
		t.Errorf("Expected ids to be 0, got %+v", ids)
	}

	keywords := obj.GetKeywords()
	if len(keywords) != 2 {
		t.Errorf("Expected ids to be 2, got %+v", ids)
	}
}
