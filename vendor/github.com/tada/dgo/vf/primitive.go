// Package vf (Value Factory) contains all factory methods for creating values
package vf

import (
	"math/big"
	"regexp"
	"time"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
)

// True is the dgo.Value representation of true
const True = internal.True

// False is the dgo.Value representation of false
const False = internal.False

// Nil is the dgo.Value representation of nil
const Nil = internal.Nil

// Boolean returns a Boolean that represents the given bool
func Boolean(v bool) dgo.Boolean {
	if v {
		return True
	}
	return False
}

// BigInt returns the given value as a dgo.BigInt
func BigInt(value *big.Int) dgo.BigInt {
	return internal.BigInt(value)
}

// BigFloat returns the given value as a dgo.BigFloat
func BigFloat(value *big.Float) dgo.BigFloat {
	return internal.BigFloat(value)
}

// Integer returns the given value as a dgo.Integer
func Integer(value int64) dgo.Integer {
	return internal.Integer(value)
}

// Float returns the given value as a dgo.Float
func Float(value float64) dgo.Float {
	return internal.Float(value)
}

// String returns the given string as a dgo.String
func String(string string) dgo.String {
	return internal.String(string)
}

// Time returns the given timestamp as a dgo.Time
func Time(ts time.Time) dgo.Time {
	return internal.Time(ts)
}

// TimeFromString returns the given time string as a dgo.Time. The string must conform to
// the time.RFC3339 or time.RFC3339Nano format. The function will panic if the given string
// cannot be parsed.
func TimeFromString(s string) dgo.Time {
	return internal.TimeFromString(s)
}

// Regexp returns the given regexp as a dgo.Regexp
func Regexp(rx *regexp.Regexp) dgo.Regexp {
	return internal.Regexp(rx)
}
