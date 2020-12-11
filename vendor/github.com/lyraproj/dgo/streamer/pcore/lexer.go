package pcore

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lyraproj/dgo/parser"
	"github.com/lyraproj/dgo/util"
)

const (
	end = iota
	integer
	float
	stringLiteral
	regexpLiteral
	identifier
	name
	rocket
)

func tokenTypeString(t int) (s string) {
	switch t {
	case end:
		s = "end"
	case name:
		s = "name"
	case identifier:
		s = "identifier"
	case integer:
		s = "integer"
	case float:
		s = "float"
	case regexpLiteral:
		s = "regexp"
	case stringLiteral:
		s = "string"
	case rocket:
		s = "rocket"
	default:
		s = string(rune(t))
	}
	return
}

func tokenString(t *parser.Token) string {
	return fmt.Sprintf("%s: '%s'", tokenTypeString(t.Type), t.Value)
}

func badToken(r rune) error {
	return fmt.Errorf("unexpected character '%c'", r)
}

func nextToken(sr *util.StringReader) (t *parser.Token) {
	for {
		r := sr.Next()
		if r == 0 {
			return &parser.Token{Type: end}
		}

		switch r {
		case ' ', '\t', '\n':
			continue
		case '#':
			consumeLineComment(sr)
			continue
		case '\'', '"':
			t = &parser.Token{Value: parser.ConsumeString(sr, r), Type: stringLiteral}
		case '/':
			t = &parser.Token{Value: parser.ConsumeRegexp(sr), Type: regexpLiteral}
		case '=':
			if sr.Peek() == '>' {
				sr.Next()
				t = &parser.Token{Value: `=>`, Type: rocket}
			} else {
				t = &parser.Token{Type: int(r)}
			}
		case '-', '+':
			n := sr.Next()
			if n < '0' || n > '9' {
				panic(badToken(r))
			}
			buf := bytes.NewBufferString(string(r))
			tkn := parser.ConsumeNumber(sr, n, buf, integer)
			t = &parser.Token{Value: buf.String(), Type: tkn}
		default:
			t = buildToken(r, sr)
		}
		break
	}
	return t
}

func buildToken(r rune, sr *util.StringReader) *parser.Token {
	switch {
	case parser.IsDigit(r):
		buf := bytes.NewBufferString(``)
		tkn := parser.ConsumeNumber(sr, r, buf, integer)
		return &parser.Token{Value: buf.String(), Type: tkn}
	case parser.IsUpperCase(r):
		buf := bytes.NewBufferString(``)
		consumeTypeName(sr, r, buf)
		return &parser.Token{Value: buf.String(), Type: name}
	case parser.IsLowerCase(r):
		buf := bytes.NewBufferString(``)
		consumeIdentifier(sr, r, buf)
		return &parser.Token{Value: buf.String(), Type: identifier}
	default:
		return &parser.Token{Type: int(r)}
	}
}

func consumeLineComment(sr *util.StringReader) {
	for {
		switch sr.Next() {
		case 0, '\n':
			return
		}
	}
}

func consumeIdentifier(sr *util.StringReader, start rune, buf io.Writer) {
	util.WriteRune(buf, start)
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		case ':':
			sr.Next()
			util.WriteRune(buf, r)
			r = sr.Next()
			if r == ':' {
				util.WriteRune(buf, r)
				r = sr.Next()
				if r == '_' || parser.IsLowerCase(r) {
					util.WriteRune(buf, r)
					continue
				}
			}
			panic(badToken(r))
		default:
			if r == '_' || parser.IsLetterOrDigit(r) {
				sr.Next()
				util.WriteRune(buf, r)
				continue
			}
			return
		}
	}
}

func consumeTypeName(sr *util.StringReader, start rune, buf io.Writer) {
	util.WriteRune(buf, start)
	for {
		r := sr.Peek()
		switch r {
		case 0:
			return
		case ':':
			sr.Next()
			util.WriteRune(buf, r)
			r = sr.Next()
			if r == ':' {
				util.WriteRune(buf, r)
				r = sr.Next()
				if parser.IsUpperCase(r) {
					util.WriteRune(buf, r)
					continue
				}
			}
			panic(badToken(r))
		default:
			if r == '_' || parser.IsLetterOrDigit(r) {
				sr.Next()
				util.WriteRune(buf, r)
				continue
			}
			return
		}
	}
}
