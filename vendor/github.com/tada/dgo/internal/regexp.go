package internal

import (
	"fmt"
	"io"
	"reflect"
	"regexp"

	"github.com/tada/catch/pio"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/util"
)

type (
	// regexpType represents an regexp type without constraints
	regexpType int

	regexpVal regexp.Regexp
)

// DefaultRegexpType is the unconstrained Regexp type
const DefaultRegexpType = regexpType(0)

var reflectRegexpType = reflect.TypeOf(&regexp.Regexp{})

func (t regexpType) Assignable(ot dgo.Type) bool {
	switch ot.(type) {
	case *regexpVal, regexpType:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t regexpType) Equals(v interface{}) bool {
	return t == v
}

func (t regexpType) HashCode() dgo.Hash {
	return dgo.Hash(dgo.TiRegexp)
}

func (t regexpType) Instance(v interface{}) bool {
	_, ok := v.(*regexpVal)
	if !ok {
		_, ok = v.(*regexp.Regexp)
	}
	return ok
}

func (t regexpType) IsInstance(v *regexp.Regexp) bool {
	return true
}

func (t regexpType) ReflectType() reflect.Type {
	return reflectRegexpType
}

func (t regexpType) String() string {
	return TypeString(t)
}

func (t regexpType) Type() dgo.Type {
	return MetaType(t)
}

func (t regexpType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexp
}

func (v *regexpVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *regexpVal) Format(s fmt.State, format rune) {
	doFormat((*regexp.Regexp)(v), s, format)
}

func (v *regexpVal) Generic() dgo.Type {
	return DefaultRegexpType
}

func (v *regexpVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *regexpVal) IsInstance(ov *regexp.Regexp) bool {
	return v.GoRegexp().String() == ov.String()
}

func (v *regexpVal) ReflectType() reflect.Type {
	return reflectRegexpType
}

func (v *regexpVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexpExact
}

// Regexp returns the given regexp as a dgo.Regexp
func Regexp(rx *regexp.Regexp) dgo.Regexp {
	return (*regexpVal)(rx)
}

func (v *regexpVal) GoRegexp() *regexp.Regexp {
	return (*regexp.Regexp)(v)
}

func (v *regexpVal) Equals(other interface{}) bool {
	if ot, ok := other.(*regexpVal); ok {
		return (*regexp.Regexp)(v).String() == (*regexp.Regexp)(ot).String()
	}
	if ot, ok := other.(*regexp.Regexp); ok {
		return (*regexp.Regexp)(v).String() == (ot).String()
	}
	return false
}

func (v *regexpVal) HashCode() dgo.Hash {
	return util.StringHash((*regexp.Regexp)(v).String())
}

func (v *regexpVal) ReflectTo(value reflect.Value) {
	rv := reflect.ValueOf((*regexp.Regexp)(v))
	k := value.Kind()
	if !(k == reflect.Ptr || k == reflect.Interface) {
		rv = rv.Elem()
	}
	value.Set(rv)
}

func (v *regexpVal) String() string {
	return TypeString(v)
}

func (v *regexpVal) Type() dgo.Type {
	return v
}

// RegexpSlashQuote converts the given string into a slash delimited string with internal slashes escaped
// and writes it on the given builder.
func RegexpSlashQuote(sb io.Writer, str string) {
	pio.WriteByte(sb, '/')
	for _, c := range str {
		switch c {
		case '\t':
			pio.WriteString(sb, `\t`)
		case '\n':
			pio.WriteString(sb, `\n`)
		case '\r':
			pio.WriteString(sb, `\r`)
		case '/':
			pio.WriteString(sb, `\/`)
		case '\\':
			pio.WriteString(sb, `\\`)
		default:
			if c < 0x20 {
				util.Fprintf(sb, `\u{%X}`, c)
			} else {
				pio.WriteRune(sb, c)
			}
		}
	}
	pio.WriteByte(sb, '/')
}
