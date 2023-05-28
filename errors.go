package env

import (
	"errors"
	"fmt"
)

// ErrNotStructPtr is returned if you pass something that is not a pointer to a
// Struct to Parse
var ErrNotStructPtr = errors.New("env: expected a pointer to a Struct")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct {
	KeyName   string
	FieldName string
	TypeName  string
	Value     string
	Err       error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("env: assigning '%s' to '%s': converting '%s' to type '%s'. details: %s", e.KeyName, e.FieldName, e.Value, e.TypeName, e.Err)
}
