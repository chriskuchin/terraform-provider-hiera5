package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"unicode/utf8"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	binaryType struct {
		sizeRange
	}

	// binary must be a struct because the []byte slice isn't hashable.
	binary struct {
		bytes []byte
	}

	binaryFrozen struct {
		binary
	}
)

// DefaultBinaryType is the unconstrained Binary type
var DefaultBinaryType = &binaryType{sizeRange{0, dgo.UnboundedSize}}

// BinaryType returns a new dgo.BinaryType. It can be called with two optional integer arguments denoting
// the min and max length of the binary. If only one integer is given, it represents the min length.
func BinaryType(args ...interface{}) dgo.BinaryType {
	switch len(args) {
	case 0:
		return DefaultBinaryType
	case 1:
		if a0, ok := Value(args[0]).(dgo.Integer); ok {
			return SizedBinaryType(int(a0.GoInt()), dgo.UnboundedSize)
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
	if min == 0 && max == dgo.UnboundedSize {
		return DefaultBinaryType
	}
	return &binaryType{sizeRange: sizeRange{min: uint32(min), max: uint32(max)}}
}

func (t *binaryType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(dgo.BinaryType); ok {
		return int(t.min) <= ot.Min() && int(t.max) >= ot.Max()
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *binaryType) Equals(other interface{}) bool {
	if ob, ok := other.(*binaryType); ok {
		return *t == *ob
	}
	return false
}

func (t *binaryType) HashCode() dgo.Hash {
	return t.sizeRangeHash(dgo.TiBinary)
}

func (t *binaryType) Instance(value interface{}) bool {
	yes := false
	switch ov := value.(type) {
	case *binary:
		yes = t.IsInstance(ov.bytes)
	case *binaryFrozen:
		yes = t.IsInstance(ov.bytes)
	case []byte:
		yes = t.IsInstance(ov)
	}
	return yes
}

func (t *binaryType) IsInstance(v []byte) bool {
	return t.inRange(len(v))
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
	return MetaType(t)
}

func (t *binaryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBinary
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
				panic(catch.Error(IllegalAssignment(bt, v)))
			}
			bs[i] = byte(v.(intVal))
		})
		b = &binaryFrozen{binary{bytes: bs}}
	case dgo.String:
		b = BinaryFromEncoded(arg.GoString(), enc)
	default:
		panic(illegalArgument(`binary`, `binary, string, or array`, []interface{}{arg}, 0))
	}
	if !t.Instance(b) {
		panic(catch.Error(IllegalAssignment(t, b)))
	}
	return b
}

// Binary creates a new Binary based on the given slice. If frozen is true, the
// binary will be immutable and contain a copy of the slice, otherwise the slice
// is simply wrapped and modifications to its elements will also modify the binary.
func Binary(bs []byte, frozen bool) dgo.Binary {
	if frozen {
		return &binaryFrozen{binary{bytes: bytesCopy(bs)}}
	}
	return &binary{bytes: bs}
}

// BinaryFromEncoded creates a new Binary from the given string and encoding. Encoding can be one of:
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
		panic(catch.Error(err))
	}
	return &binaryFrozen{binary{bytes: bs}}
}

// BinaryFromData creates a new frozen Binary based on data read from the given io.Reader.
func BinaryFromData(data io.Reader) dgo.Binary {
	bs, err := ioutil.ReadAll(data)
	if err != nil {
		panic(catch.Error(err))
	}
	return &binaryFrozen{binary{bytes: bs}}
}

func (v *binary) Copy(frozen bool) dgo.Binary {
	cp := bytesCopy(v.bytes)
	if frozen {
		return &binaryFrozen{binary{bytes: cp}}
	}
	return &binary{bytes: cp}
}

func (v *binaryFrozen) Copy(frozen bool) dgo.Binary {
	if frozen {
		return v
	}
	return &binary{bytes: bytesCopy(v.bytes)}
}

func (v *binary) CompareTo(other interface{}) (int, bool) {
	r := 0
	ok := true
	switch ov := other.(type) {
	case nil, nilValue:
		r = 1
	case *binary:
		r = bytes.Compare(v.bytes, ov.bytes)
	case *binaryFrozen:
		r = bytes.Compare(v.bytes, ov.bytes)
	case []byte:
		r = bytes.Compare(v.bytes, ov)
	default:
		ok = false
	}
	return r, ok
}

func (v *binary) Encode() string {
	return base64.StdEncoding.Strict().EncodeToString(v.bytes)
}

func (v *binary) Equals(other interface{}) bool {
	yes := false
	switch ov := other.(type) {
	case *binary:
		yes = bytes.Equal(v.bytes, ov.bytes)
	case *binaryFrozen:
		yes = bytes.Equal(v.bytes, ov.bytes)
	case []byte:
		yes = bytes.Equal(v.bytes, ov)
	}
	return yes
}

func (v *binary) Format(s fmt.State, format rune) {
	doFormat(v.bytes, s, format)
}

func (v *binary) Frozen() bool {
	return false
}

func (v *binary) FrozenCopy() dgo.Value {
	return v.Copy(true)
}

func (v *binary) ThawedCopy() dgo.Value {
	return v.Copy(false)
}

func (v *binary) GoBytes() []byte {
	return v.bytes
}

func (v *binary) HashCode() dgo.Hash {
	return bytesHash(v.bytes)
}

func (v *binary) ReflectTo(value reflect.Value) {
	setReflected(value, v.bytes)
}

func setReflected(value reflect.Value, bs []byte) {
	switch value.Kind() {
	case reflect.Ptr:
		x := reflect.New(reflectBinaryType)
		x.Elem().SetBytes(bs)
		value.Set(x)
	case reflect.Slice:
		value.SetBytes(bs)
	default:
		value.Set(reflect.ValueOf(bs))
	}
}

func (v *binary) String() string {
	return TypeString(v)
}

func (v *binary) Type() dgo.Type {
	return v
}

func (v *binary) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *binary) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *binary) IsInstance(b []byte) bool {
	return bytes.Equal(v.bytes, b)
}

func (v *binary) Max() int {
	return len(v.bytes)
}

func (v *binary) Min() int {
	return len(v.bytes)
}

func (v *binary) New(arg dgo.Value) dgo.Value {
	return newBinary(v, arg)
}

func (v *binary) ReflectType() reflect.Type {
	return reflectBinaryType
}

func (v *binary) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBinaryExact
}

func (v *binary) Unbounded() bool {
	return false
}

func (v *binaryFrozen) Frozen() bool {
	return true
}

func (v *binaryFrozen) FrozenCopy() dgo.Value {
	return v
}

func (v *binaryFrozen) GoBytes() []byte {
	return bytesCopy(v.bytes)
}

func (v *binaryFrozen) ReflectTo(value reflect.Value) {
	setReflected(value, bytesCopy(v.bytes))
}

func bytesCopy(bs []byte) []byte {
	c := make([]byte, len(bs))
	copy(c, bs)
	return c
}

func bytesHash(s []byte) dgo.Hash {
	h := dgo.Hash(1)
	for i := range s {
		h = 31*h + dgo.Hash(s[i])
	}
	return h
}
