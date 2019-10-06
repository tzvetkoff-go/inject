package inject

import (
	"fmt"
	"reflect"
)

// Error constants.
const (
	ErrorNotPointerToStruct    = iota // Object passed to inject is not a struct pointer.
	ErrorNoProviderForName            // Provider for name missing,
	ErrorCannotSetPrivateField        // Cannot set private fields.
)

// Error represents an error returned by Inject().
type Error struct {
	Code      int
	Message   string
	Reference interface{}
}

// Error implements the error interfrace.
func (i *Error) Error() string {
	return i.Message
}

// NewErrorNotPointerToStruct returns a new error with Code = ErrorNotPointerToStruct.
func NewErrorNotPointerToStruct(reference interface{}) *Error {
	return &Error{
		Code:      ErrorNotPointerToStruct,
		Message:   fmt.Sprintf("object of type `%s` is not a struct pointer", reflect.TypeOf(reference).String()),
		Reference: reference,
	}
}

// IsErrorNotPointerToStruct tests if the error has Code = ErrorNotPointerToStruct.
func IsErrorNotPointerToStruct(err error) bool {
	if err, ok := err.(*Error); ok {
		return err.Code == ErrorNotPointerToStruct
	}

	return false
}

// NewErrorNoProviderForName returns a new error with Code = ErrorNoProviderForName.
func NewErrorNoProviderForName(typeName string, reference interface{}) *Error {
	return &Error{
		Code:      ErrorNoProviderForName,
		Message:   fmt.Sprintf("no provider provider returning `%s`", typeName),
		Reference: reference,
	}
}

// IsErrorNoProviderForName tests if the error has Code = ErrorNoProviderForName.
func IsErrorNoProviderForName(err error) bool {
	if err, ok := err.(*Error); ok {
		return err.Code == ErrorNoProviderForName
	}

	return false
}

// NewErrorCannotSetPrivateField returns a new error with Code = ErrorCannotSetPrivateField.
func NewErrorCannotSetPrivateField(fieldName string, reference interface{}) *Error {
	return &Error{
		Code:      ErrorCannotSetPrivateField,
		Message:   fmt.Sprintf("cannot set private field `%s`", fieldName),
		Reference: reference,
	}
}

// IsErrorCannotSetPrivateField tests if the error has Code = ErrorCannotSetPrivateField.
func IsErrorCannotSetPrivateField(err error) bool {
	if err, ok := err.(*Error); ok {
		return err.Code == ErrorCannotSetPrivateField
	}

	return false
}
