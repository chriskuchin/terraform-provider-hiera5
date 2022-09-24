package catch

import (
	"errors"
	"fmt"
)

// special error that is recovered by the Do() function returning its original cause.
type errorWithCause struct {
	// Cause is the original error
	cause error
}

// Error returns a new special error object that wraps the given cause. The returned error can be used as an
// argument to the IsError() and Cause() functions.
//
// One argument
//
// The argument can be either an error or a string. An error will used verbatim whereas a string will be
// converted using the function errors.New().
//
// Multiple arguments
//
// The first argument must be a string. If that is not the case, an illegal argument panic is raised. The
// string and the remaining arguments are passed to the function fmt.Errorf() to create an error.
func Error(cause ...interface{}) error {
	var err error
	if nc := len(cause); nc > 0 {
		switch a1 := cause[0].(type) {
		case error:
			if nc == 1 {
				err = &errorWithCause{cause: a1}
			}
		case string:
			if nc == 1 {
				err = &errorWithCause{cause: errors.New(a1)}
			} else {
				err = &errorWithCause{cause: fmt.Errorf(a1, cause[1:]...)}
			}
		}
	}
	if err != nil {
		return err
	}
	panic(errors.New("catch.Error(): illegal argument"))
}

// Error returns the result of calling Error() on the contained cause
func (e *errorWithCause) Error() string {
	return e.cause.Error()
}

// IsError returns true if, and only if, the argument is an error produced by the Error function.
func IsError(e interface{}) bool {
	_, ok := e.(*errorWithCause)
	return ok
}

// Cause returns the underlying cause of the error provided that the argument is an error created by the Error function
// or nil if that is not the case.
func Cause(e interface{}) error {
	if e, ok := e.(*errorWithCause); ok {
		return e.cause
	}
	return nil
}
