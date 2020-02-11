package dgo

// An Indenter helps building strings where all newlines are supposed to be followed by
// a sequence of zero or many spaces that reflect an indent level.
type Indenter interface {
	// Append appends a string to the internal buffer without checking for newlines
	Append(string)

	// AppendRune appends a rune to the internal buffer without checking for newlines
	AppendRune(rune)

	// AppendValue appends the string form of the given value to the internal buffer. Indentation
	// will be recursive if the value implements Indentable
	AppendValue(interface{})

	// AppendIndented is like Append but replaces all occurrences of newline with an indented newline
	AppendIndented(string)

	// AppendBool writes the string "true" or "false" to the internal buffer
	AppendBool(bool)

	// AppendInt writes the result of calling strconf.Itoa() in the given argument
	AppendInt(int)

	// Indent returns a new indenter instance that shares the same buffer but has an
	// indent level that is increased by one.
	Indent() Indenter

	// Indenting returns true if this indenter has an indent string with a length > 0
	Indenting() bool

	// Len returns the current number of bytes that has been appended to the indenter
	Len() int

	// Level returns the indent level for the indenter
	Level() int

	// NewLine writes a newline followed by the current indent after trimming trailing whitespaces
	NewLine()

	// Printf formats according to a format specifier and writes to the internal buffer.
	Printf(string, ...interface{})

	// Reset resets the internal buffer. It does not reset the indent
	Reset()

	// String returns the current string that has been built using the indenter. Trailing whitespaces
	// are deleted from all lines.
	String() string

	// Write appends a slice of bytes to the internal buffer without checking for newlines
	Write(p []byte) (n int, err error)

	// WriteString appends a string to the internal buffer without checking for newlines
	WriteString(s string) (n int, err error)
}

// An Indentable can create build a string representation of itself using an indenter
type Indentable interface {
	// AppendTo appends a string representation of the Node to the indenter
	AppendTo(w Indenter)
}
