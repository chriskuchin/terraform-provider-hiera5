package internal

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

type (
	// sizedStringType is a size constrained String type. It only represents all strings whose
	// length is within the inclusive min, max range
	sizedStringType struct {
		min int
		max int
	}

	// defaultStringType represents an string without constraints
	defaultStringType int

	// defaultDgoStringType represents strings that conform to Dgo type syntax
	defaultDgoStringType int

	// exactStringType only represents its own string
	exactStringType struct {
		exactType
		value *hstring
	}

	ciStringType struct {
		exactStringType
	}

	// patternType constrains its instances to those that match the regexp pattern
	patternType struct {
		*regexp.Regexp
	}

	// hstring is a string that caches the hash value when it is computed
	hstring struct {
		s string
		h int
	}
)

var formatPattern = PatternType(regexp.MustCompile(`\A%([\s\[+#0{<(|-]*)([1-9][0-9]*)?(?:\.([0-9]+))?([a-zA-Z])\z`))

func newString(t dgo.Type, arg dgo.Value) dgo.String {
	var s dgo.String
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`string`, 1, 2)
		if args.Len() == 2 {
			var v interface{}
			FromValue(args.Get(0), &v)
			s = String(fmt.Sprintf(args.Arg(`string`, 1, formatPattern).(dgo.String).GoString(), v))
		} else {
			arg = args.Get(0)
		}
	}

	if s == nil {
		switch arg := arg.(type) {
		case dgo.String:
			s = arg
		default:
			s = String(arg.String())
		}
	}

	if !t.Instance(s) {
		panic(IllegalAssignment(t, s))
	}
	return s
}

func (t defaultDgoStringType) String() string {
	return TypeString(t)
}

func (t defaultDgoStringType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultDgoStringType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultDgoStringType) HashCode() int {
	return int(dgo.TiDgoString)
}

func (t defaultDgoStringType) Assignable(other dgo.Type) bool {
	if t == other {
		return true
	}
	if ot, ok := other.(*exactStringType); ok {
		return t.Instance(ot.value)
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultDgoStringType) Instance(value interface{}) (ok bool) {
	var s string
	var ov *hstring
	ov, ok = value.(*hstring)
	if ok {
		s = ov.s
	} else {
		s, ok = value.(string)
	}
	if ok {
		ok = len(s) > 0
		if ok {
			defer func() {
				if recover() != nil {
					ok = false
				}
			}()
			Parse(s)
		}
	}
	return
}

func (t defaultDgoStringType) New(arg dgo.Value) dgo.Value {
	return newString(t, arg)
}

func (t defaultDgoStringType) ReflectType() reflect.Type {
	return reflectStringType
}

func (t defaultDgoStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiDgoString
}

func (t defaultDgoStringType) Max() int {
	return math.MaxInt64
}

func (t defaultDgoStringType) Min() int {
	return 1
}

func (t defaultDgoStringType) Unbounded() bool {
	return false
}

// DefaultStringType is the unconstrained String type
const DefaultStringType = defaultStringType(0)

// DefaultDgoStringType is the unconstrained Dgo String type
const DefaultDgoStringType = defaultDgoStringType(0)

var reflectStringType = reflect.TypeOf(` `)

// EnumType returns a StringType that represents all of the given strings
func EnumType(strings []string) dgo.Type {
	switch len(strings) {
	case 0:
		return &notType{DefaultAnyType}
	case 1:
		return makeHString(strings[0]).Type()
	}
	ts := make([]dgo.Value, len(strings))
	for i := range strings {
		ts[i] = makeHString(strings[i]).Type()
	}
	return &anyOfType{slice: ts, frozen: true}
}

// CiEnumType returns a StringType that represents all strings that are equal to one of the given strings
// under Unicode case-folding.
func CiEnumType(strings []string) dgo.Type {
	switch len(strings) {
	case 0:
		return &notType{DefaultAnyType}
	case 1:
		return CiStringType(strings[0])
	}
	ts := make([]dgo.Value, len(strings))
	for i := range strings {
		ts[i] = CiStringType(strings[i])
	}
	return &anyOfType{slice: ts, frozen: true}
}

// String returns the dgo.String for the given string
func String(s string) dgo.String {
	return makeHString(s)
}

func (t defaultStringType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultStringType, defaultDgoStringType, *exactStringType, *ciStringType, *sizedStringType, *patternType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultStringType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultStringType) HashCode() int {
	return int(dgo.TiString)
}

func (t defaultStringType) Instance(value interface{}) bool {
	if _, ok := value.(*hstring); ok {
		return true
	}
	if _, ok := value.(string); ok {
		return true
	}
	return false
}

func (t defaultStringType) Max() int {
	return math.MaxInt64
}

func (t defaultStringType) Min() int {
	return 0
}

func (t defaultStringType) New(arg dgo.Value) dgo.Value {
	return newString(t, arg)
}

func (t defaultStringType) String() string {
	return TypeString(t)
}

func (t defaultStringType) ReflectType() reflect.Type {
	return reflectStringType
}

func (t defaultStringType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiString
}

func (t defaultStringType) Unbounded() bool {
	return true
}

func (t *exactStringType) Generic() dgo.Type {
	return DefaultStringType
}

func (t *exactStringType) Max() int {
	return len(t.value.s)
}

func (t *exactStringType) Min() int {
	return len(t.value.s)
}

func (t *exactStringType) New(arg dgo.Value) dgo.Value {
	return newString(t, arg)
}

func (t *exactStringType) ReflectType() reflect.Type {
	return reflectStringType
}

func (t *exactStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStringExact
}

func (t *exactStringType) Unbounded() bool {
	return false
}

func (t *exactStringType) ExactValue() dgo.Value {
	return t.value
}

// CiStringType returns a StringType that is constrained to strings that are equal to the given string under
// Unicode case-folding.
func CiStringType(si interface{}) dgo.StringType {
	var s string
	if ov, ok := si.(*hstring); ok {
		s = ov.s
	} else {
		s = si.(string)
	}
	et := &ciStringType{exactStringType: exactStringType{value: makeHString(strings.ToLower(s))}}
	et.ExactType = et
	return et
}

func (t *ciStringType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*exactStringType); ok {
		return strings.EqualFold(t.value.s, ot.value.s)
	}
	return t.exactType.Assignable(other)
}

func (t *ciStringType) Instance(v interface{}) bool {
	if ov, ok := v.(*hstring); ok {
		return strings.EqualFold(t.value.s, ov.s)
	}
	if ov, ok := v.(string); ok {
		return strings.EqualFold(t.value.s, ov)
	}
	return false
}

func (t *ciStringType) New(arg dgo.Value) dgo.Value {
	return newString(t, arg)
}

func (t *ciStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiCiString
}

// PatternType returns a StringType that is constrained to strings that match the given
// regular expression pattern
func PatternType(pattern *regexp.Regexp) dgo.Type {
	return &patternType{Regexp: pattern}
}

// StringType returns a new dgo.StringType. It can be called with two optional integer arguments denoting
// the min and max length of the string. If only one integer is given, it represents the min length.
//
// The method can also be called with one string parameter. The returned type will then match that exact
// string and nothing else.
func StringType(args []interface{}) dgo.StringType {
	switch len(args) {
	case 0:
		return DefaultStringType
	case 1:
		switch a0 := Value(args[0]).(type) {
		case dgo.String:
			return a0.Type().(dgo.StringType) // Exact string type
		case dgo.Integer:
			return SizedStringType(int(a0.GoInt()), math.MaxInt64)
		}
		panic(illegalArgument(`StringType`, `Integer or String`, args, 0))
	case 2:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			var a1 dgo.Integer
			if a1, ok = Value(args[1]).(dgo.Integer); ok {
				return SizedStringType(int(a0.GoInt()), int(a1.GoInt()))
			}
			panic(illegalArgument(`StringType`, `Integer`, args, 1))
		}
		panic(illegalArgument(`StringType`, `Integer`, args, 0))
	}
	panic(illegalArgumentCount(`StringType`, 0, 2, len(args)))
}

func (t *patternType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case *exactStringType:
		return t.IsInstance(ot.value.s)
	case *patternType:
		return t.rxString() == ot.rxString()
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *patternType) Equals(v interface{}) bool {
	if ov, ok := v.(*patternType); ok {
		return t.rxString() == ov.rxString()
	}
	return false
}

func (t *patternType) Generic() dgo.Type {
	return DefaultStringType
}

func (t *patternType) HashCode() int {
	return util.StringHash(t.rxString())
}

func (t *patternType) Instance(v interface{}) bool {
	if sv, ok := v.(*hstring); ok {
		return t.MatchString(sv.s)
	}
	if sv, ok := v.(string); ok {
		return t.MatchString(sv)
	}
	return false
}

func (t *patternType) IsInstance(v string) bool {
	return t.MatchString(v)
}

func (t *patternType) Max() int {
	return math.MaxInt64
}

func (t *patternType) Min() int {
	return 0
}

func (t *patternType) New(arg dgo.Value) dgo.Value {
	return newString(t, arg)
}

func (t *patternType) ReflectType() reflect.Type {
	return reflectStringType
}

func (t *patternType) rxString() string {
	return (t.Regexp).String()
}

func (t *patternType) Type() dgo.Type {
	return &metaType{t}
}

func (t *patternType) String() string {
	return TypeString(t)
}

func (t *patternType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStringPattern
}

func (t *patternType) Unbounded() bool {
	return true
}

func (t *patternType) ExactValue() dgo.Value {
	return (*regexpVal)(t.Regexp)
}

// SizedStringType returns a StringType that is constrained to strings whose length is within the
// inclusive range given by min and max.
func SizedStringType(min, max int) dgo.StringType {
	if min < 0 {
		min = 0
	}
	if max < min {
		tmp := max
		max = min
		min = tmp
	}
	if min == 0 && max == math.MaxInt64 {
		return DefaultStringType
	}
	return &sizedStringType{min: min, max: max}
}

func (t *sizedStringType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultDgoStringType:
		return t.min <= 1
	case *exactStringType:
		return t.IsInstance(ot.value.s)
	case *ciStringType:
		return t.IsInstance(ot.value.s)
	case *sizedStringType:
		return t.min <= ot.min && t.max >= ot.max
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *sizedStringType) Equals(v interface{}) bool {
	if ob, ok := v.(*sizedStringType); ok {
		return *t == *ob
	}
	return false
}

func (t *sizedStringType) HashCode() int {
	h := int(dgo.TiStringSized)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	return h
}

func (t *sizedStringType) Instance(v interface{}) bool {
	if sv, ok := v.(*hstring); ok {
		return t.IsInstance(sv.s)
	}
	if sv, ok := v.(string); ok {
		return t.IsInstance(sv)
	}
	return false
}

func (t *sizedStringType) IsInstance(v string) bool {
	l := len(v)
	return t.min <= l && l <= t.max
}

func (t *sizedStringType) Max() int {
	return t.max
}

func (t *sizedStringType) Min() int {
	return t.min
}

func (t *sizedStringType) New(arg dgo.Value) dgo.Value {
	return newString(t, arg)
}

func (t *sizedStringType) ReflectType() reflect.Type {
	return reflectStringType
}

func (t *sizedStringType) String() string {
	return TypeString(t)
}

func (t *sizedStringType) Type() dgo.Type {
	return &metaType{t}
}

func (t *sizedStringType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStringSized
}

func (t *sizedStringType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

func makeHString(s string) *hstring {
	return &hstring{s: s}
}

func (v *hstring) AppendTo(w dgo.Indenter) {
	w.Append(strconv.Quote(v.s))
}

func (v *hstring) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	switch ov := other.(type) {
	case *hstring:
		switch {
		case v.s > ov.s:
			r = 1
		case v.s < ov.s:
			r = -1
		default:
			r = 0
		}
	case string:
		switch {
		case v.s > ov:
			r = 1
		case v.s < ov:
			r = -1
		default:
			r = 0
		}
	case nilValue:
		r = 1
	default:
		ok = false
	}
	return
}

func (v *hstring) Equals(other interface{}) bool {
	// comparison for *hstring must be first here or the HashMap will get a penalty. It
	// must always use *hstring to get the hash code
	if ov, ok := other.(*hstring); ok {
		return v.s == ov.s
	}
	if s, ok := other.(string); ok {
		return v.s == s
	}
	return false
}

func (v *hstring) GoString() string {
	return v.s
}

func (v *hstring) HashCode() int {
	if v.h == 0 {
		v.h = util.StringHash(v.s)
	}
	return v.h
}

func (v *hstring) ReflectTo(value reflect.Value) {
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(v.s))
	case reflect.Ptr:
		value.Set(reflect.ValueOf(&v.s))
	default:
		value.SetString(v.s)
	}
}

func (v *hstring) String() string {
	return v.s
}

func (v *hstring) Type() dgo.Type {
	et := &exactStringType{value: v}
	et.ExactType = et
	return et
}
