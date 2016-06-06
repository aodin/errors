/*
Package errors provides a robust errors type which implements the built-in error interface. It includes the following:

* Code, an int field for integer error codes such as HTTP status

* Meta, a []string field for high-level errors

* Fields, a map[string]string field for named errors

It supports both JSON and XML marshaling.
*/
package errors

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

// Error is an error structure with meta and field-specific errors
type Error struct {
	Code   int               `json:"code,omitempty"`
	Meta   []string          `json:"meta"`
	Fields map[string]string `json:"fields"`
}

// Error must implement the error interface
var _ error = Error{}

// AddMeta appends a meta error
func (er *Error) AddMeta(msg string, args ...interface{}) {
	er.Meta = append(er.Meta, fmt.Sprintf(msg, args...))
}

// Add is an alias for AddMeta
func (er *Error) Add(msg string, args ...interface{}) {
	er.AddMeta(msg, args...)
}

// Error implements the error interface
func (er Error) Error() string {
	output := er.Meta
	for field, err := range er.Fields {
		output = append(output, fmt.Sprintf("%s (%s)", err, field))
	}
	if er.Code == 0 {
		return fmt.Sprintf("%s", strings.Join(output, "; "))
	}
	return fmt.Sprintf("%d: %s", er.Code, strings.Join(output, "; "))
}

// Exists returns true if there are either meta or field errors
func (er Error) Exists() bool {
	return len(er.Meta) > 0 || len(er.Fields) > 0
}

// IsEmpty return false if there are no meta or field errors
func (er Error) IsEmpty() bool {
	return !er.Exists()
}

// InField returns true if the given field has an error
func (er Error) InField(field string) bool {
	_, exists := er.Fields[field]
	return exists
}

type field struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type metas struct {
	Meta []string
}

type fields struct {
	Field []field
}

// MarshalXML implements a custom marshaler because Errors has a map
func (er Error) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var f fields
	for key, value := range er.Fields {
		f.Field = append(f.Field, field{
			XMLName: xml.Name{Local: key},
			Value:   value,
		})
	}
	anon := struct {
		Code   int `xml:",omitempty"`
		Metas  metas
		Fields fields
	}{
		Code:   er.Code,
		Metas:  metas{Meta: er.Meta},
		Fields: f,
	}
	return e.EncodeElement(anon, start)
}

// SetField sets the error message for the field. Mutiple calls will overwrite
// previous messages for that field.
func (er *Error) SetField(field, msg string, args ...interface{}) {
	er.Fields[field] = fmt.Sprintf(msg, args...)
}

// Set is an alias for SetField
func (er *Error) Set(field, msg string, args ...interface{}) {
	er.SetField(field, msg, args...)
}

// BadRequest creates an error with a 400 status code
func BadRequest() *Error {
	errs := New()
	errs.Code = http.StatusBadRequest
	return errs
}

// Message creates a new *Error with the given message
func Message(msg string, args ...interface{}) *Error {
	return Meta(0, msg, args...)
}

// Meta returns an error with a pre-set meta error
func Meta(code int, msg string, args ...interface{}) *Error {
	errs := New()
	errs.Code = code
	errs.Meta = []string{fmt.Sprintf(msg, args...)}
	return errs
}

// New creates a new empty error
func New() *Error {
	return &Error{
		Meta:   make([]string, 0), // Make JSON meta [] instead of null
		Fields: make(map[string]string),
	}
}
