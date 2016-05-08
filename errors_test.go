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
}

func TestErrors_MarshalXML(t *testing.T) {
	er := BadRequest()
	er.AddMeta("Not Found")
	er.AddMeta("Unknown format")
	er.SetField("UUID", "Invalid UUID format")
	er.SetField("ID", "Missing ID")
	out, err := xml.Marshal(er)
	if err != nil {
		t.Fatalf("unexpected err during xml.Marshal: %s", err)
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
