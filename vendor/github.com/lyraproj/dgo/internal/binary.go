package internal

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"unicode/utf8"

	"github.com/lyraproj/dgo/dgo"
)

type (
	binaryType struct {
		min int
		max int
	}

	exactBinaryType struct {
		exactType
		value *binary
	}

	binary struct {
		bytes  []byte
		frozen bool
	}
)

// DefaultBinaryType is the unconstrained Binary type
var DefaultBinaryType = &binaryType{0, math.MaxInt64}

// BinaryType returns a new dgo.BinaryType. It can be called with two optional integer arguments denoting
// the min and max length of the binary. If only one integer is given, it represents the min length.
func BinaryType(args ...interface{}) dgo.BinaryType {
	switch len(args) {
	case 0:
		return DefaultBinaryType
	case 1:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			return SizedBinaryType(int(a0.GoInt()), math.MaxInt64)
		}
		panic(illegalArgument(`BinaryType`, `Integer`, args, 0))
	case 2:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			var a1 dgo.Integer
			if a1, ok = Value(args[1]).(dgo.Integer); ok {
				return SizedBinaryType(int(a0.GoInt()), int(a1.GoInt()))
			}
			panic(illegalArgument(`BinaryType`, `Integer`, args, 1))
		}
		panic(illegalArgument(`BinaryType`, `Integer`, args, 0))
	}
	panic(illegalArgumentCount(`BinaryType`, 0, 2, len(args)))
}

// SizedBinaryType returns a BinaryType that is constrained to binaries whose size is within the
// inclusive range given by min and max.
func SizedBinaryType(min, max int) dgo.BinaryType {
	if min < 0 {
		min = 0
	}
	if max < min {
		tmp := max
		max = min
		min = tmp
	}
	if min == 0 && max == math.MaxInt64 {
		return DefaultBinaryType
	}
	return &binaryType{min: min, max: max}
}

func (t *binaryType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(dgo.BinaryType); ok {
		return t.min <= ot.Min() && t.max >= ot.Max()
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *binaryType) Equals(other interface{}) bool {
	if ob, ok := other.(*binaryType); ok {
		return *t == *ob
	}
	return false
}

func (t *binaryType) HashCode() int {
	h := int(dgo.TiBinary)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	return h
}

func (t *binaryType) Instance(value interface{}) bool {
	if ov, ok := value.(*binary); ok {
		return t.IsInstance(ov.bytes)
	}
	if ov, ok := value.([]byte); ok {
		return t.IsInstance(ov)
	}
	return false
}

func (t *binaryType) IsInstance(v []byte) bool {
	l := len(v)
	return t.min <= l && l <= t.max
}

func (t *binaryType) Max() int {
	return t.max
}

func (t *binaryType) Min() int {
	return t.min
}

func (t *binaryType) New(arg dgo.Value) dgo.Value {
	return newBinary(t, arg)
}

var reflectBinaryType = reflect.TypeOf([]byte{})

func (t *binaryType) ReflectType() reflect.Type {
	return reflectBinaryType
}

func (t *binaryType) String() string {
	return TypeString(t)
}

func (t *binaryType) Type() dgo.Type {
	return &metaType{t}
}

func (t *binaryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBinary
}

func (t *binaryType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

func (v *exactBinaryType) IsInstance(b []byte) bool {
	return bytes.Equal(v.value.bytes, b)
}

func (v *exactBinaryType) Max() int {
	return len(v.value.bytes)
}

func (v *exactBinaryType) Min() int {
	return len(v.value.bytes)
}

func (v *exactBinaryType) New(arg dgo.Value) dgo.Value {
	return newBinary(v, arg)
}

func (v *exactBinaryType) ReflectType() reflect.Type {
	return reflectBinaryType
}

func (v *exactBinaryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBinaryExact
}

func (v *exactBinaryType) Unbounded() bool {
	return false
}

func (v *exactBinaryType) ExactValue() dgo.Value {
	return v.value
}

var encType = EnumType([]string{`%B`, `%b`, `%u`, `%s`, `%r`})

func newBinary(t dgo.Type, arg dgo.Value) dgo.Value {
	enc := `%B`
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`binary`, 1, 2)
		if args.Len() == 2 {
			arg = args.Arg(`binary`, 0, DefaultStringType)
			enc = args.Arg(`binary`, 1, encType).(dgo.String).GoString()
		} else {
			arg = args.Get(0)
		}
	}
	var b dgo.Value
	switch arg := arg.(type) {
	case dgo.Binary:
		b = arg
	case dgo.Array:
		bs := make([]byte, arg.Len())
		bt := primitivePTypes[reflect.Uint8]
		arg.EachWithIndex(func(v dgo.Value, i int) {
			if !bt.Instance(v) {
				panic(IllegalAssignment(bt, v))
			}
			bs[i] = byte(v.(intVal))
		})
		b = Binary(bs, true)
	case dgo.String:
		b = BinaryFromEncoded(arg.GoString(), enc)
	default:
		panic(illegalArgument(`binary`, `binary, string, or array`, []interface{}{arg}, 0))
	}
	if !t.Instance(b) {
		panic(IllegalAssignment(t, b))
	}
	return b
}

// Binary creates a new Binary based on the given slice. If frozen is true, the
// binary will be immutable and contain a copy of the slice, otherwise the slice
// is simply wrapped and modifications to its elements will also modify the binary.
func Binary(bs []byte, frozen bool) dgo.Binary {
	if frozen {
		c := make([]byte, len(bs))
		copy(c, bs)
		bs = c
	}
	return &binary{bytes: bs, frozen: frozen}
}

// BinaryFromEncoded creates a new Binary from the given string and encoding. Enocding can be one of:
//
// `%b`: base64.StdEncoding
//
// `%u`: base64.URLEncoding
//
// `%B`: base64.StdEncoding.Strict()
//
// `%s`: check using utf8.ValidString(str), then cast to []byte
//
// `%r`: cast to []byte
func BinaryFromEncoded(str, enc string) dgo.Binary {
	var bs []byte
	var err error

	switch enc {
	case `%b`:
		bs, err = base64.StdEncoding.DecodeString(str)
	case `%u`:
		bs, err = base64.URLEncoding.DecodeString(str)
	case `%B`:
		bs, err = base64.StdEncoding.Strict().DecodeString(str)
	case `%s`:
		if !utf8.ValidString(str) {
			panic(illegalArgument(`binary`, `valid utf8 string`, []interface{}{str, enc}, 0))
		}
		bs = []byte(str)
	case `%r`:
		bs = []byte(str)
	default:
		panic(illegalArgument(`binary`, `one of the supported format specifiers %B, %b, %s, %r, %u`, []interface{}{str, enc}, 1))
	}
	if err != nil {
		panic(err)
	}
	return &binary{bytes: bs, frozen: true}
}

// BinaryFromData creates a new frozen Binary based on data read from the given io.Reader.
func BinaryFromData(data io.Reader) dgo.Binary {
	bs, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err)
	}
	return &binary{bytes: bs, frozen: true}
}

func (v *binary) Copy(frozen bool) dgo.Binary {
	if frozen && v.frozen {
		return v
	}
	cp := make([]byte, len(v.bytes))
	copy(cp, v.bytes)
	return &binary{bytes: cp, frozen: frozen}
}

func (v *binary) CompareTo(other interface{}) (int, bool) {
	var b []byte
	var ob *binary
	var ok bool
	if ob, ok = other.(*binary); ok {
		if v == ob {
			return 0, true
		}
		b = ob.bytes
	} else {
		b, ok = other.([]byte)
		if !ok {
			if other == nil || other == Nil {
				return 1, true
			}
			return 0, false
		}
	}
	a := v.bytes
	top := len(a)
	max := len(b)
	r := 0
	if top < max {
		r = -1
		max = top
	} else if top > max {
		r = 1
	}
	for i := 0; i < max; i++ {
		c := int(a[i]) - int(b[i])
		if c != 0 {
			if c > 0 {
				r = 1
			} else {
				r = -1
			}
			break
		}
	}
	return r, ok
}

func (v *binary) Encode() string {
	return base64.StdEncoding.Strict().EncodeToString(v.bytes)
}

func (v *binary) Equals(other interface{}) bool {
	if ot, ok := other.(*binary); ok {
		return bytes.Equal(v.bytes, ot.bytes)
	}
	if ot, ok := other.([]byte); ok {
		return bytes.Equal(v.bytes, ot)
	}
	return false
}

func (v *binary) Freeze() {
	if !v.frozen {
		bs := v.bytes
		v.bytes = make([]byte, len(bs))
		copy(v.bytes, bs)
		v.frozen = true
	}
}

func (v *binary) Frozen() bool {
	return v.frozen
}

func (v *binary) FrozenCopy() dgo.Value {
	if !v.frozen {
		cs := make([]byte, len(v.bytes))
		copy(cs, v.bytes)
		return &binary{bytes: cs, frozen: true}
	}
	return v
}

func (v *binary) GoBytes() []byte {
	if v.frozen {
		c := make([]byte, len(v.bytes))
		copy(c, v.bytes)
		return c
	}
	return v.bytes
}

func (v *binary) HashCode() int {
	return bytesHash(v.bytes)
}

func (v *binary) ReflectTo(value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		x := reflect.New(reflectBinaryType)
		x.Elem().SetBytes(v.GoBytes())
		value.Set(x)
	case reflect.Slice:
		value.SetBytes(v.GoBytes())
	default:
		value.Set(reflect.ValueOf(v.GoBytes()))
	}
}

func (v *binary) String() string {
	return base64.StdEncoding.Strict().EncodeToString(v.bytes)
}

func (v *binary) Type() dgo.Type {
	et := &exactBinaryType{value: v}
	et.ExactType = et
	return et
}

func bytesHash(s []byte) int {
	h := 1
	for i := range s {
		h = 31*h + int(s[i])
	}
	return h
}
