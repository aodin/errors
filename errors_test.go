package errors

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	err := New()
	if err.Exists() {
		t.Error("empty Errors should not exist")
	}

	err.Add("Hello %s", "World")
	if err.IsEmpty() {
		t.Error("Errors with a meta should not be empty")
	}

	err.Set("whatever", "failure")
	if !err.InField("whatever") {
		t.Error("InField() should return true when an error has been set")
	}

	metaErr := Meta(500, "I am an error")
	if !strings.HasPrefix(metaErr.Error(), "500:") {
		t.Error("meta Errors should have code as the first part of its output")
	}

	uncodedErr := New()
	uncodedErr.Set("what", "huh")
	if !strings.HasPrefix(uncodedErr.Error(), "huh") {
		t.Error("uncoded Errors should begin with an error")
	}

	msg := Message("Hello %d", 1)
	if msg == nil {
		t.Fatal("Message() should not create a nil error")
	}
	if msg.Error() != "Hello 1" {
		t.Error(
			`unexpected error from Message(): want "Hello 1", have "%s"`,
			msg.Error(),
		)
	}
}

func TestErrors_MarshalXML(t *testing.T) {
	err := BadRequest()
	err.AddMeta("Not Found")
	err.AddMeta("Unknown format")
	err.SetField("UUID", "Invalid UUID format")
	err.SetField("ID", "Missing ID")
	out, xmlErr := xml.Marshal(err)
	if xmlErr != nil {
		t.Fatalf("unexpected error during xml.Marshal: %s", xmlErr)
	}

	// Map iteration is non-deterministic
	if len(string(out)) != 163 {
		t.Errorf("unexpected output length: %d != 163", len(out))
	}
	prefix := `<Error><Code>400</Code><Metas><Meta>Not Found</Meta><Meta>Unknown format</Meta></Metas><Fields>`
	if !strings.HasPrefix(string(out), prefix) {
		t.Errorf("unexpected start to XML output: %s", out)
	}
}
