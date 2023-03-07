package error

import (
	"errors"
	"reflect"

	"github.com/ikonglong/op-status"
)

type OpError struct {
	cause  error
	status *opstatus.Status
}

func NewWithStatus(status opstatus.Status) *OpError {
	return &OpError{
		status: &status,
	}
}

func NewWithStatusAndCause(status opstatus.Status, cause error) *OpError {
	return &OpError{
		status: &status,
		cause:  cause,
	}
}

func (e *OpError) Status() *opstatus.Status {
	return e.status
}

func (e *OpError) Cause() error {
	return e.cause
}

func (e *OpError) Error() string {
	// todo
	return ""
}

// StatusFromErrChain finds the first OpError from the causal chain of given error.
// If one is found, return its status. Otherwise, return nil
func StatusFromErrChain(err error) *opstatus.Status {
	if IsNil(err) {
		return nil
	}
	cause := err
	for !IsNil(cause) {
		if match, opErr := AsOpError(cause); match {
			return opErr.Status()
		}
		cause = errors.Unwrap(cause)
	}
	return nil
}

// AsOpError finds the first error in given error chain that is of type opError,
// and if one is found, sets target to that error value and returns true. Otherwise,
// it returns false.
func AsOpError(err error) (bool, *OpError) {
	var opErr OpError
	return errors.As(err, &opErr), &opErr
}

// IsNil tells if given err is nil. If the value of given interface variable is nil
// or the value stored into the second word of given interface value is nil, return true.
// Otherwise, return false.
func IsNil(err error) bool {
	if err == nil {
		return true
	}

	// For case: if the second word of given interface value is nil, `err == nil` is false
	ifaceVal := reflect.ValueOf(err)
	switch ifaceVal.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return ifaceVal.IsNil()
	}
	return false
}
