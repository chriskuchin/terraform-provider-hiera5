package pio

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	"github.com/tada/catch"
)

// Write writes the bytes b to the Writer, returning its length.
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func Write(w io.Writer, b []byte) int {
	n, err := w.Write(b)
	if err != nil {
		panic(catch.Error(err))
	}
	return n
}

// WriteString writes the bytes of s to the Writer, returning its length.
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func WriteString(w io.Writer, s string) int {
	return Write(w, []byte(s))
}

// WriteByte writes the byte r to the Writer.
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func WriteByte(w io.Writer, b byte) {
	Write(w, []byte{b})
}

// WriteRune writes the UTF-8 encoding of Unicode code point r to the Writer, returning its length.
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func WriteRune(w io.Writer, r rune) int {
	if r < utf8.RuneSelf {
		WriteByte(w, byte(r))
		return 1
	}
	b := make([]byte, utf8.UTFMax)
	n := utf8.EncodeRune(b, r)
	Write(w, b[:n])
	return n
}

// WriteBool writes the string "true" or "false" onto the stream
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func WriteBool(w io.Writer, b bool) {
	s := "false"
	if b {
		s = "true"
	}
	WriteString(w, s)
}

// WriteInt writes decimal string representation of the given integer onto the stream
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func WriteInt(w io.Writer, i int64) {
	WriteString(w, strconv.FormatInt(i, 10))
}

// WriteFloat writes the "%g" string representation of the given integer onto the stream
//
// Any error from the io.Writer will result in a panic(catch.Error(err))
func WriteFloat(w io.Writer, f float64) {
	WriteString(w, strconv.FormatFloat(f, 'g', -1, 64))
}

// WriteQuotedString writes the given string enclosed in double quotes on the given Writer and panics
// if an error occurs.
func WriteQuotedString(b io.Writer, s string) {
	WriteByte(b, '"')
	for _, c := range s {
		switch c {
		case '\t':
			WriteString(b, `\t`)
		case '\n':
			WriteString(b, `\n`)
		case '\r':
			WriteString(b, `\r`)
		case '"':
			WriteString(b, `\"`)
		case '\\':
			WriteString(b, `\\`)
		case '\a':
			WriteString(b, `\a`)
		case '\b':
			WriteString(b, `\b`)
		case '\f':
			WriteString(b, `\f`)
		case '\v':
			WriteString(b, `\v`)
		default:
			if c < 0x20 {
				_, _ = fmt.Fprintf(b, `\x%.2X`, c)
			} else {
				WriteRune(b, c)
			}
		}
	}
	WriteByte(b, '"')
}
