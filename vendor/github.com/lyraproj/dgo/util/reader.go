package util

import (
	"errors"
	"unicode/utf8"
)

// StringReader is a helper for reading runes of a String. Its typical use is from a tokenizer
type StringReader struct {
	p int
	l int
	c int
	s string
}

// NewStringReader returns a new reader initialized to read the given string
func NewStringReader(s string) *StringReader {
	return &StringReader{s: s}
}

// Next returns the next rune in the string and advances the position, or returns 0 if the end as been reached
func (r *StringReader) Next() rune {
	if r.p >= len(r.s) {
		if r.p == len(r.s) {
			r.p++
			r.c++
		}
		return 0
	}
	c := rune(r.s[r.p])
	if c < utf8.RuneSelf {
		r.p++
		if c == '\n' {
			r.l++
			r.c = 0
		}
		r.c++
	} else {
		var size int
		c, size = utf8.DecodeRuneInString(r.s[r.p:])
		if c == utf8.RuneError {
			panic(errors.New("unicode error"))
		}
		r.p += size
		r.c++
	}
	return c
}

// Peek returns the next rune in the string or 0 if the end as been reached. The position is not affected
func (r *StringReader) Peek() rune {
	if r.p >= len(r.s) {
		return 0
	}
	c := rune(r.s[r.p])
	if c >= utf8.RuneSelf {
		c, _ = utf8.DecodeRuneInString(r.s[r.p:])
		if c == utf8.RuneError {
			panic(errors.New("unicode error"))
		}
	}
	return c
}

// Peek2 returns the second next rune in the string or 0 if the end as been reached. The position is not affected
func (r *StringReader) Peek2() rune {
	if r.p >= len(r.s) {
		return 0
	}
	c := rune(r.s[r.p])
	var sz int
	if c >= utf8.RuneSelf {
		c, sz = utf8.DecodeRuneInString(r.s[r.p:])
		if c == utf8.RuneError {
			panic(errors.New("unicode error"))
		}
	} else {
		sz = 1
	}

	np := r.p + sz
	if np >= len(r.s) {
		return 0
	}

	c = rune(r.s[np])
	if c >= utf8.RuneSelf {
		c, _ = utf8.DecodeRuneInString(r.s[np:])
		if c == utf8.RuneError {
			panic(errors.New("unicode error"))
		}
	}
	return c
}

// Column returns the column of the line at the current position. First column is 1.
func (r *StringReader) Column() int {
	return r.c + 1
}

// Line returns the line of the current position. First line is 1
func (r *StringReader) Line() int {
	return r.l + 1
}

// Pos returns the the current position
func (r *StringReader) Pos() int {
	return r.p
}

// Rewind resets the current position to zero so that Next and Peek will return the first character of the string
func (r *StringReader) Rewind() {
	r.p = 0
	r.l = 0
	r.c = 0
}
