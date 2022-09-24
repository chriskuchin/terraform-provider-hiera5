package util

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tada/catch/pio"
	"github.com/tada/dgo/dgo"
)

type indenter struct {
	b *bytes.Buffer
	i int
	s string
}

// ToString will produce an unindented string from an Indentable
func ToString(ia dgo.Indentable) string {
	i := NewIndenter(``)
	ia.AppendTo(i)
	return i.String()
}

// ToIndentedString will produce a string from an Indentable using an indenter initialized
// with a two space indentation.
func ToIndentedString(ia dgo.Indentable) string {
	i := NewIndenter(`  `)
	ia.AppendTo(i)
	return i.String()
}

// NewIndenter creates a new indenter for indent level zero using the given string to perform
// one level of indentation. An empty string will yield unindented output
func NewIndenter(indent string) dgo.Indenter {
	return &indenter{b: &bytes.Buffer{}, i: 0, s: indent}
}

func (i *indenter) Len() int {
	return i.b.Len()
}

func (i *indenter) Level() int {
	return i.i
}

func (i *indenter) Reset() {
	i.b.Reset()
}

func (i *indenter) String() string {
	n := bytes.NewBuffer(make([]byte, 0, i.b.Len()))
	wb := &bytes.Buffer{}
	for {
		r, _, err := i.b.ReadRune()
		if err == io.EOF {
			break
		}
		if r == ' ' || r == '\t' {
			// Defer whitespace output
			pio.WriteByte(wb, byte(r))
			continue
		}
		if r == '\n' {
			// Truncate trailing space
			wb.Reset()
		} else if wb.Len() > 0 {
			_, _ = n.Write(wb.Bytes())
			wb.Reset()
		}
		pio.WriteRune(n, r)
	}
	return n.String()
}

func (i *indenter) WriteString(s string) (n int, err error) {
	return i.b.WriteString(s)
}

func (i *indenter) Write(p []byte) (n int, err error) {
	return i.b.Write(p)
}

func (i *indenter) AppendRune(r rune) {
	pio.WriteRune(i.b, r)
}

func (i *indenter) Append(s string) {
	pio.WriteString(i.b, s)
}

func (i *indenter) AppendValue(v interface{}) {
	switch v := v.(type) {
	case dgo.Indentable:
		v.AppendTo(i)
	case fmt.Stringer:
		pio.WriteString(i.b, v.String())
	default:
		Fprintf(i.b, "%#v", v)
	}
}

func (i *indenter) AppendIndented(s string) {
	for ni := strings.IndexByte(s, '\n'); ni >= 0; ni = strings.IndexByte(s, '\n') {
		if ni > 0 {
			pio.WriteString(i.b, s[:ni])
		}
		i.NewLine()
		ni++
		if ni >= len(s) {
			return
		}
		s = s[ni:]
	}
	if len(s) > 0 {
		pio.WriteString(i.b, s)
	}
}

func (i *indenter) AppendBool(b bool) {
	var s string
	if b {
		s = `true`
	} else {
		s = `false`
	}
	pio.WriteString(i.b, s)
}

func (i *indenter) AppendInt(b int) {
	pio.WriteString(i.b, strconv.Itoa(b))
}

func (i *indenter) Indent() dgo.Indenter {
	c := *i
	c.i++
	return &c
}

func (i *indenter) Indenting() bool {
	return len(i.s) > 0
}

func (i *indenter) Printf(s string, args ...interface{}) {
	Fprintf(i.b, s, args...)
}

func (i *indenter) NewLine() {
	if len(i.s) > 0 {
		pio.WriteByte(i.b, '\n')
		for n := 0; n < i.i; n++ {
			pio.WriteString(i.b, i.s)
		}
	}
}
