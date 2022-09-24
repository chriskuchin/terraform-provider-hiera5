// Package pio contains io related functions that instead of returning an error, panics with
// a catch.Error which has a Cause (the original error). The catch.Do function can then be used to
// recover the panic and return the original error
package pio
