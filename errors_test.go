package errors

import (
	"testing"
)

func TestErrors(t *testing.T) {
	err := New()
	if err.Exists() {
		t.Error("empty Errors should not exist")
	}

	err.AddMeta("Hello %s", "World")
	if err.IsEmpty() {
		t.Error("Errors with a meta should not be empty")
	}

	err.SetField("whatever", "failure")
	if !err.InField("whatever") {
		t.Error("InField() should return true when an error has been set")
	}

	metaErr := Meta(500, "I am an error")
	if metaErr.Error() == "" {
		t.Error("meta Errors should have Error() output")
	}
}
