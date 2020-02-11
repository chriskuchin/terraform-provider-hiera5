package internal

import (
	"io"
	"reflect"
	"regexp"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// regexpType represents an regexp type without constraints
	regexpType int

	exactRegexpType struct {
		exactType
		value *regexpVal
	}

	regexpVal regexp.Regexp
)

// DefaultRegexpType is the unconstrained Regexp type
const DefaultRegexpType = regexpType(0)

var reflectRegexpType = reflect.TypeOf(&regexp.Regexp{})

func (t regexpType) Assignable(ot dgo.Type) bool {
	switch ot.(type) {
	case regexpType, *exactRegexpType:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t regexpType) Equals(v interface{}) bool {
	return t == v
}

func (t regexpType) HashCode() int {
	return int(dgo.TiRegexp)
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
	return &metaType{t}
}

func (t regexpType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexp
}

func (t *exactRegexpType) Generic() dgo.Type {
	return DefaultRegexpType
}

func (t *exactRegexpType) IsInstance(v *regexp.Regexp) bool {
	return t.value.String() == v.String()
}

func (t *exactRegexpType) ReflectType() reflect.Type {
	return reflectRegexpType
}

func (t *exactRegexpType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiRegexpExact
}

func (t *exactRegexpType) ExactValue() dgo.Value {
	return t.value
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

func (v *regexpVal) HashCode() int {
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
	return (*regexp.Regexp)(v).String()
}

func (v *regexpVal) Type() dgo.Type {
	et := &exactRegexpType{value: v}
	et.ExactType = et
	return et
}

// RegexpSlashQuote converts the given string into a slash delimited string with internal slashes escaped
// and writes it on the given builder.
func RegexpSlashQuote(sb io.Writer, str string) {
	util.WriteByte(sb, '/')
	for _, c := range str {
		switch c {
		case '\t':
			util.WriteString(sb, `\t`)
		case '\n':
			util.WriteString(sb, `\n`)
		case '\r':
			util.WriteString(sb, `\r`)
		case '/':
			util.WriteString(sb, `\/`)
		case '\\':
			util.WriteString(sb, `\\`)
		default:
			if c < 0x20 {
				util.Fprintf(sb, `\u{%X}`, c)
			} else {
				util.WriteRune(sb, c)
			}
		}
	}
	util.WriteByte(sb, '/')
}
