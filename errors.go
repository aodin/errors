package errors

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

// Error is an error structure with meta and field-specific errors
type Error struct {
	Code   int               `json:"-"`
	Meta   []string          `json:"meta"`
	Fields map[string]string `json:"fields"`
}

// Error must implement the error interface
var _ error = Error{}

// AddMeta appends a meta error
func (e *Error) AddMeta(msg string, args ...interface{}) {
	e.Meta = append(e.Meta, fmt.Sprintf(msg, args...))
}

// Error implements the error interface
func (e Error) Error() string {
	output := e.Meta
	for field, err := range e.Fields {
		output = append(output, fmt.Sprintf("%s (%s)", err, field))
	}
	if e.Code == 0 {
		return fmt.Sprintf("%s", strings.Join(output, "; "))
	}
	return fmt.Sprintf("%d: %s", e.Code, strings.Join(output, "; "))
}

// Exists returns true if there are either meta or field errors
func (e Error) Exists() bool {
	return len(e.Meta) > 0 || len(e.Fields) > 0
}

// IsEmpty return false if there are no meta or field errors
func (e Error) IsEmpty() bool {
	return !e.Exists()
}

// InField returns true if the given field has an error
func (e Error) InField(field string) bool {
	_, exists := e.Fields[field]
	return exists
}

// MarshalXML implements a custom marshaler because Errors has a map
func (er Error) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}

// SetField sets the error message for the field. Mutiple calls will overwrite
// previous messages for that field.
func (e *Error) SetField(field, msg string, args ...interface{}) {
	e.Fields[field] = fmt.Sprintf(msg, args...)
}

// BadRequest creates an error with a 400 status code
func BadRequest() *Error {
	errs := New()
	errs.Code = http.StatusBadRequest
	return errs
}

// New creates a new empty error
func New() *Error {
	return &Error{
		Meta:   make([]string, 0), // Make JSON meta [] instead of null
		Fields: make(map[string]string),
	}
}

// Meta returns an error with a pre-set meta error
func Meta(code int, msg string, args ...interface{}) *Error {
	errs := New()
	errs.Code = code
	errs.Meta = []string{fmt.Sprintf(msg, args...)}
	return errs
}
