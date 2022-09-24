package util

import (
	"fmt"
	"io"
)

// Fprintf is like fmt.Fprintf but it panics in case of error instead of returning it.
func Fprintf(b io.Writer, format string, args ...interface{}) int {
	n, err := fmt.Fprintf(b, format, args...)
	if err != nil {
		panic(err)
	}
	return n
}

// Fprintln is like fmt.Fprintln but it panics in case of error instead of returning it.
func Fprintln(b io.Writer, args ...interface{}) int {
	n, err := fmt.Fprintln(b, args...)
	if err != nil {
		panic(err)
	}
	return n
}
