package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/lyraproj/dgo/internal"

	"github.com/lyraproj/dgo/util"
)

const (
	end = iota
	integer
	float
	stringLiteral
	regexpLiteral
	identifier
	dotdot
	dotdotdot
)

const (
	digit = 1 << iota
	hexDigit
	letter
	lcLetter
	ucLetter
	idCharStart
	idChar
	exprEnd
)

var charTypes = [128]uint8{}

// This init() function ensures that the charTypes map is initialized with correct bits.
func init() {
	i := int('$')
	charTypes[i] = idCharStart | idChar
	i = int('_')
	charTypes[i] = idCharStart | idChar
	for i = '0'; i <= '9'; i++ {
		charTypes[i] = digit | hexDigit | idChar
	}
	for i = 'A'; i <= 'F'; i++ {
		charTypes[i] = hexDigit | letter | ucLetter | idCharStart | idChar
	}
	for i = 'G'; i <= 'Z'; i++ {
		charTypes[i] = letter | ucLetter | idCharStart | idChar
	}
	for i = 'a'; i <= 'f'; i++ {
		charTypes[i] = hexDigit | letter | lcLetter | idCharStart | idChar
	}
	for i = 'g'; i <= 'z'; i++ {
		charTypes[i] = letter | lcLetter | idCharStart | idChar
	}
	for _, i = range []int{0, ')', '}', ']', ',', ':', '?', '|', '&', '^', '.'} {
		charTypes[i] = exprEnd
	}
}

// Token is what the lexer produces for the parser to consume
type Token struct {
	Value string
	Type  int
}

func tokenString(t *Token) (s string) {
	if t == nil || t.Type == end {
		return "EOT"
	}
	switch t.Type {
	case identifier, integer, float, dotdot, dotdotdot:
		s = t.Value
	case regexpLiteral:
		sb := &strings.Builder{}
		internal.RegexpSlashQuote(sb, t.Value)
		s = sb.String()
	case stringLiteral:
		s = strconv.Quote(t.Value)
	default:
		s = fmt.Sprintf(`'%c'`, rune(t.Type))
	}
	return
}

func badToken(r rune) error {
	if r == 0 {
		return errors.New(`unexpected end`)
	}
	return fmt.Errorf("unexpected character '%c'", r)
}

func nextToken(sr *util.StringReader) *Token {
	for {
		var t *Token
		r := sr.Next()
		switch r {
		case 0:
			t = &Token{``, end}
		case ' ', '\t', '\n':
			continue
		case '`':
			t = &Token{consumeRawString(sr), stringLiteral}
		case '"':
			t = &Token{ConsumeString(sr, r), stringLiteral}
		case '/':
			t = &Token{ConsumeRegexp(sr), regexpLiteral}
		case '.':
			if sr.Peek() == '.' {
				sr.Next()
				if sr.Peek() == '.' {
					sr.Next()
					t = &Token{`...`, dotdotdot}
				} else {
					t = &Token{`..`, dotdot}
				}
			} else {
				t = &Token{Type: int(r)}
			}
		case '-', '+':
			n := sr.Next()
			if !IsDigit(n) {
				panic(badToken(n))
			}
			buf := bytes.NewBufferString(``)
			if r == '-' {
				util.WriteRune(buf, r)
			}
			tkn := ConsumeNumber(sr, n, buf, integer)
			t = &Token{buf.String(), tkn}
		default:
			t = buildToken(r, sr)
		}
		return t
	}
}

func buildToken(r rune, sr *util.StringReader) *Token {
	switch {
	case IsDigit(r):
		buf := bytes.NewBufferString(``)
		tkn := ConsumeNumber(sr, r, buf, integer)
		return &Token{buf.String(), tkn}
	case IsIdentifierStart(r):
		buf := bytes.NewBufferString(``)
		consumeIdentifier(sr, r, buf)
		return &Token{buf.String(), identifier}
	default:
		return &Token{Type: int(r)}
	}
}

func consumeUnsignedInteger(sr *util.StringReader, buf io.Writer) {
	for {
		r := sr.Peek()
		switch {
		case r == '.' || IsLetter(r):
			panic(badToken(r))
		case IsDigit(r):
			sr.Next()
			util.WriteRune(buf, r)
		default:
			return
		}
	}
}

func isExpressionEnd(r rune) bool {
	return r < 128 && (charTypes[r]&exprEnd) != 0
}

// IsDigit returns true if the given rune is an ASCII digit
func IsDigit(r rune) bool {
	return r < 128 && (charTypes[r]&digit) != 0
}

// IsLetter returns true if the given rune is an ASCII letter
func IsLetter(r rune) bool {
	return r < 128 && (charTypes[r]&letter) != 0
}

// IsLetterOrDigit returns true if the given rune is an ASCII letter or digit
func IsLetterOrDigit(r rune) bool {
	return r < 128 && (charTypes[r]&(letter|digit)) != 0
}

// IsHex returns true if the given rune is an ASCII digit, 'a' - 'f' or 'A' - 'F'
func IsHex(r rune) bool {
	return r < 128 && (charTypes[r]&hexDigit) != 0
}

// IsIdentifier returns true if the given rune is an ASCII letter, a digit, '$' or '_'
func IsIdentifier(r rune) bool {
	return r < 128 && (charTypes[r]&idChar) != 0
}

// IsIdentifierStart returns true if the given rune is an ASCII letter, '$' or '_'
func IsIdentifierStart(r rune) bool {
	return r < 128 && (charTypes[r]&idCharStart) != 0
}

// IsUpperCase returns true if the given rune is an uppercase ASCII letter
func IsUpperCase(r rune) bool {
	return r < 128 && (charTypes[r]&ucLetter) != 0
}

// IsLowerCase returns true if the given rune is an lowercase ASCII letter
func IsLowerCase(r rune) bool {
	return r < 128 && (charTypes[r]&lcLetter) != 0
}

func consumeExponent(sr *util.StringReader, buf io.Writer) {
	for {
		r := sr.Next()
		switch r {
		case 0:
			panic(errors.New("unexpected end"))
		case '+', '-':
			util.WriteRune(buf, r)
			r = sr.Next()
			fallthrough
		default:
			if IsDigit(r) {
				util.WriteRune(buf, r)
				consumeUnsignedInteger(sr, buf)
				return
			}
			panic(badToken(r))
		}
	}
}

func consumeHexInteger(sr *util.StringReader, buf io.Writer) {
	for IsHex(sr.Peek()) {
		util.WriteRune(buf, sr.Next())
	}
}

// ConsumeNumber consumes the current number into the given Writer and returns the consumed token type.
func ConsumeNumber(sr *util.StringReader, start rune, buf io.Writer, t int) int {
	util.WriteRune(buf, start)
	firstZero := t != float && start == '0'

	for r := sr.Peek(); r != 0; r = sr.Peek() {
		switch r {
		case '0':
			sr.Next()
			util.WriteRune(buf, r)
		case 'e', 'E':
			sr.Next()
			util.WriteRune(buf, r)
			consumeExponent(sr, buf)
			return float
		case 'x', 'X':
			if firstZero {
				sr.Next()
				util.WriteRune(buf, r)
				r = sr.Next()
				if IsHex(r) {
					util.WriteRune(buf, r)
					consumeHexInteger(sr, buf)
					return t
				}
			}
			panic(badToken(r))
		case '.':
			if sr.Peek2() == '.' {
				return t
			}
			if t != float {
				sr.Next()
				util.WriteRune(buf, r)
				r = sr.Next()
				if IsDigit(r) {
					return ConsumeNumber(sr, r, buf, float)
				}
			}
			panic(badToken(r))
		default:
			if !IsDigit(r) {
				return t
			}
			sr.Next()
			util.WriteRune(buf, r)
		}
	}
	return t
}

// ConsumeRegexp consumes the current regexp up to the ending '/' character, taking escaped
// escapes and ends into account.
func ConsumeRegexp(sr *util.StringReader) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		switch r {
		case '/':
			return buf.String()
		case '\\':
			r = sr.Next()
			switch r {
			case 0:
				panic(errors.New("unterminated regexp"))
			case '/': // Escape is removed
			default:
				util.WriteRune(buf, '\\')
			}
			util.WriteRune(buf, r)
		case 0, '\n':
			panic(errors.New("unterminated regexp"))
		default:
			util.WriteRune(buf, r)
		}
	}
}

// ConsumeString consumes the current string up to the given end character while taking
// escaped nl, cr, tab, escape, and end character into account.
func ConsumeString(sr *util.StringReader, end rune) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		if r == end {
			return buf.String()
		}
		switch r {
		case 0:
			panic(errors.New("unterminated string"))
		case '\\':
			consumeEscape(sr.Next(), buf, end)
		case '\n':
			panic(errors.New("unterminated string"))
		default:
			util.WriteRune(buf, r)
		}
	}
}

func consumeEscape(r rune, buf io.Writer, end rune) {
	switch r {
	case 0:
		panic(errors.New("unterminated string"))
	case 'n':
		r = '\n'
	case 'r':
		r = '\r'
	case 't':
		r = '\t'
	case '\\':
	default:
		if r != end {
			panic(fmt.Errorf("illegal escape '\\%c'", r))
		}
	}
	util.WriteRune(buf, r)
}

func consumeRawString(sr *util.StringReader) string {
	buf := bytes.NewBufferString(``)
	for {
		r := sr.Next()
		if r == '`' {
			return buf.String()
		}
		if r == 0 {
			panic(errors.New("unterminated string"))
		}
		util.WriteRune(buf, r)
	}
}

func consumeIdentifier(sr *util.StringReader, start rune, buf io.Writer) {
	util.WriteRune(buf, start)
	for IsIdentifier(sr.Peek()) {
		util.WriteRune(buf, sr.Next())
	}
}
