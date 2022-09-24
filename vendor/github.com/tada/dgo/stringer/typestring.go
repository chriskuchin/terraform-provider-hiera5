package stringer

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tada/catch/pio"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
	"github.com/tada/dgo/util"
)

const (
	commaPrio = iota
	orPrio
	xorPrio
	andPrio
	typePrio
)

type typeToString func(sb *typeBuilder, typ dgo.Type, prio int)

type typeBuilder struct {
	io.Writer
	aliasMap dgo.AliasMap
	seen     []dgo.Value
}

func newTypeBuilder(w io.Writer, am dgo.AliasMap) *typeBuilder {
	return &typeBuilder{Writer: w, aliasMap: am}
}

// TypeString produces a string with the go-like syntax for the given type.
func TypeString(typ dgo.Type) string {
	return TypeStringWithAliasMap(typ, internal.DefaultAliases())
}

// TypeStringOn produces a string with the go-like syntax for the given type onto the given io.Writer.
func TypeStringOn(typ dgo.Type, w io.Writer) {
	newTypeBuilder(w, internal.DefaultAliases()).buildTypeString(typ, 0)
}

// TypeStringWithAliasMap produces a string with the go-like syntax for the given type.
func TypeStringWithAliasMap(typ dgo.Type, am dgo.AliasMap) string {
	s := strings.Builder{}
	newTypeBuilder(&s, am).buildTypeString(typ, 0)
	return s.String()
}

func anyOf(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `|`, orPrio)
}

func oneOf(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `^`, xorPrio)
}

func allOf(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `&`, andPrio)
}

func allOfValue(sb *typeBuilder, typ dgo.Type, prio int) {
	sb.writeTernary(typ, valueAsType, prio, `&`, andPrio)
}

func array(sb *typeBuilder, typ dgo.Type, _ int) {
	at := typ.(dgo.ArrayType)
	if at.Unbounded() {
		pio.WriteString(sb, `[]`)
	} else {
		pio.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(at.Min()), int64(at.Max()))
		pio.WriteByte(sb, ']')
	}
	sb.buildTypeString(at.ElementType(), typePrio)
}

func arrayExact(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteByte(sb, '{')
	sb.joinValueTypes(typ.(dgo.Iterable), `,`, commaPrio)
	pio.WriteByte(sb, '}')
}

func binary(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.BinaryType)
	pio.WriteString(sb, `binary`)
	if !st.Unbounded() {
		pio.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(st.Min()), int64(st.Max()))
		pio.WriteByte(sb, ']')
	}
}

func stringValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteQuotedString(sb, typ.(dgo.String).GoString())
}

func bigFloatValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, `big `)
	pio.WriteString(sb, typ.(dgo.BigFloat).GoBigFloat().String())
}

func bigIntValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, typ.(dgo.BigInt).GoBigInt().String())
}

func binaryValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, `binary `)
	pio.WriteQuotedString(sb, typ.(dgo.Binary).Encode())
}

func booleanValue(sb *typeBuilder, typ dgo.Type, _ int) {
	s := `false`
	if typ.(dgo.Boolean).GoBool() {
		s = `true`
	}
	pio.WriteString(sb, s)
}

func errorValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, `error `)
	pio.WriteQuotedString(sb, typ.(error).Error())
}

func exactValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, typ.(dgo.ExactType).ExactValue().String())
}

func floatValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, util.Ftoa(typ.(dgo.Float).GoFloat()))
}

func intValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, strconv.FormatInt(typ.(dgo.Integer).GoInt(), 10))
}

func nativeValue(sb *typeBuilder, typ dgo.Type, _ int) {
	rv := typ.(dgo.Native).ReflectValue()
	if rv.CanInterface() {
		iv := rv.Interface()
		if s, ok := iv.(fmt.Stringer); ok {
			pio.WriteString(sb, s.String())
		} else {
			util.Fprintf(sb, "%#v", iv)
		}
		return
	}
	sb.writeNativeType(rv.Type())
}

func regexpValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, `regexp `)
	pio.WriteQuotedString(sb, typ.(dgo.Regexp).GoRegexp().String())
}

func timeValue(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteString(sb, `time `)
	pio.WriteQuotedString(sb, typ.(dgo.Time).GoTime().Format(time.RFC3339Nano))
}

func tuple(sb *typeBuilder, typ dgo.Type, _ int) {
	sb.writeTupleArgs(typ.(dgo.TupleType), '{', '}')
}

func _map(sb *typeBuilder, typ dgo.Type, _ int) {
	at := typ.(dgo.MapType)
	pio.WriteString(sb, `map[`)
	sb.buildTypeString(at.KeyType(), commaPrio)
	if !at.Unbounded() {
		pio.WriteByte(sb, ',')
		sb.writeSizeBoundaries(int64(at.Min()), int64(at.Max()))
	}
	pio.WriteByte(sb, ']')
	sb.buildTypeString(at.ValueType(), typePrio)
}

func mapExact(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteByte(sb, '{')
	sb.joinValueTypes(typ.(dgo.Map), `,`, commaPrio)
	pio.WriteByte(sb, '}')
}

func _struct(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteByte(sb, '{')
	st := typ.(dgo.StructMapType)
	sb.joinStructMapEntries(st)
	if st.Additional() {
		if st.Len() > 0 {
			pio.WriteByte(sb, ',')
		}
		pio.WriteString(sb, `...`)
	}
	pio.WriteByte(sb, '}')
}

func mapEntryExact(sb *typeBuilder, typ dgo.Type, _ int) {
	me := typ.(dgo.MapEntry)
	sb.buildTypeString(typeAsType(me.Key()), commaPrio)
	pio.WriteByte(sb, ':')
	sb.buildTypeString(typeAsType(me.Value()), commaPrio)
}

func floatRange(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.FloatType)
	_, big := st.Min().(dgo.BigFloat)
	if !big {
		_, big = st.Max().(dgo.BigFloat)
	}
	if big {
		pio.WriteString(sb, `big `)
		if st.Min() != nil {
			pio.WriteString(sb, st.Min().(dgo.BigFloat).GoBigFloat().String())
		}
		sb.writeRangeDots(st.Inclusive())
		if st.Max() != nil {
			pio.WriteString(sb, st.Max().(dgo.BigFloat).GoBigFloat().String())
		}
	} else {
		sb.writeRange(st.Min(), st.Max(), st.Inclusive())
	}
}

func integerRange(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.IntegerType)
	sb.writeRange(st.Min(), st.Max(), st.Inclusive())
}

func sensitive(sb *typeBuilder, typ dgo.Type, prio int) {
	pio.WriteString(sb, `sensitive`)
	if op := typ.(dgo.UnaryType).Operand(); internal.DefaultAnyType != op {
		pio.WriteByte(sb, '[')
		sb.buildTypeString(op, prio)
		pio.WriteByte(sb, ']')
	}
}

func stringPattern(sb *typeBuilder, typ dgo.Type, _ int) {
	internal.RegexpSlashQuote(sb, typ.(dgo.PatternType).GoRegexp().String())
}

func stringSized(sb *typeBuilder, typ dgo.Type, _ int) {
	st := typ.(dgo.StringType)
	pio.WriteString(sb, `string`)
	if !st.Unbounded() {
		pio.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(st.Min()), int64(st.Max()))
		pio.WriteByte(sb, ']')
	}
}

func ciString(sb *typeBuilder, typ dgo.Type, _ int) {
	pio.WriteByte(sb, '~')
	pio.WriteString(sb, strconv.Quote(typ.(dgo.String).GoString()))
}

func native(sb *typeBuilder, typ dgo.Type, _ int) {
	sb.writeNativeType(typ.(dgo.NativeType).GoType())
}

func not(sb *typeBuilder, typ dgo.Type, _ int) {
	nt := typ.(dgo.UnaryType)
	pio.WriteByte(sb, '!')
	sb.buildTypeString(nt.Operand(), typePrio)
}

func meta(sb *typeBuilder, typ dgo.Type, prio int) {
	nt := typ.(dgo.UnaryType)
	pio.WriteString(sb, `type`)
	if op := nt.Operand(); internal.DefaultAnyType != op {
		if op == nil {
			pio.WriteString(sb, `[type]`)
		} else {
			pio.WriteByte(sb, '[')
			sb.buildTypeString(op, prio)
			pio.WriteByte(sb, ']')
		}
	}
}

func function(sb *typeBuilder, typ dgo.Type, prio int) {
	ft := typ.(dgo.FunctionType)
	pio.WriteString(sb, `func`)
	sb.writeTupleArgs(ft.In(), '(', ')')
	out := ft.Out()
	if out.Len() > 0 {
		pio.WriteByte(sb, ' ')
		if out.Len() == 1 && !out.Variadic() {
			sb.buildTypeString(out.ElementTypeAt(0), prio)
		} else {
			sb.writeTupleArgs(ft.Out(), '(', ')')
		}
	}
}

func named(sb *typeBuilder, typ dgo.Type, _ int) {
	nt := typ.(dgo.NamedType)
	pio.WriteString(sb, nt.Name())
	if params := nt.Parameters(); params != nil {
		pio.WriteByte(sb, '[')
		sb.joinValueTypes(params, `,`, commaPrio)
		pio.WriteByte(sb, ']')
	}
}

var complexTypes map[dgo.TypeIdentifier]typeToString

func init() {
	complexTypes = map[dgo.TypeIdentifier]typeToString{
		dgo.TiAllOf:         allOf,
		dgo.TiAllOfValue:    allOfValue,
		dgo.TiAnyOf:         anyOf,
		dgo.TiArray:         array,
		dgo.TiArrayExact:    arrayExact,
		dgo.TiBigFloatExact: bigFloatValue,
		dgo.TiBigIntExact:   bigIntValue,
		dgo.TiBinary:        binary,
		dgo.TiBinaryExact:   binaryValue,
		dgo.TiBooleanExact:  booleanValue,
		dgo.TiCiString:      ciString,
		dgo.TiErrorExact:    errorValue,
		dgo.TiFloatExact:    floatValue,
		dgo.TiFunction:      function,
		dgo.TiIntegerExact:  intValue,
		dgo.TiIntegerRange:  integerRange,
		dgo.TiFloatRange:    floatRange,
		dgo.TiMap:           _map,
		dgo.TiMapExact:      mapExact,
		dgo.TiMapEntryExact: mapEntryExact,
		dgo.TiMeta:          meta,
		dgo.TiNamed:         named,
		dgo.TiNamedExact:    exactValue,
		dgo.TiNative:        native,
		dgo.TiNativeExact:   nativeValue,
		dgo.TiNot:           not,
		dgo.TiOneOf:         oneOf,
		dgo.TiRegexpExact:   regexpValue,
		dgo.TiSensitive:     sensitive,
		dgo.TiStringExact:   stringValue,
		dgo.TiStringPattern: stringPattern,
		dgo.TiStringSized:   stringSized,
		dgo.TiStruct:        _struct,
		dgo.TiTimeExact:     timeValue,
		dgo.TiTuple:         tuple,
	}
}

func (sb *typeBuilder) joinTypes(v dgo.Iterable, s string, prio int) {
	sb.joinX(v, typeAsType, s, prio)
}

func (sb *typeBuilder) joinValueTypes(v dgo.Iterable, s string, prio int) {
	sb.joinX(v, valueAsType, s, prio)
}

func (sb *typeBuilder) joinX(v dgo.Iterable, tc func(dgo.Value) dgo.Type, s string, prio int) {
	first := true
	v.Each(func(v dgo.Value) {
		if first {
			first = false
		} else {
			pio.WriteString(sb, s)
		}
		sb.buildTypeString(tc(v), prio)
	})
}

func (sb *typeBuilder) joinStructMapEntries(v dgo.StructMapType) {
	first := true
	v.EachEntryType(func(e dgo.StructMapEntry) {
		if first {
			first = false
		} else {
			pio.WriteByte(sb, ',')
		}
		sb.buildTypeString(e.Key().(dgo.Type), commaPrio)
		if !e.Required() {
			pio.WriteByte(sb, '?')
		}
		pio.WriteByte(sb, ':')
		sb.buildTypeString(e.Value().(dgo.Type), commaPrio)
	})
}

func (sb *typeBuilder) writeNativeType(rt reflect.Type) {
	pio.WriteString(sb, `native`)
	if rt != nil {
		pio.WriteByte(sb, '[')
		pio.WriteString(sb, strconv.Quote(rt.String()))
		pio.WriteByte(sb, ']')
	}
}

func (sb *typeBuilder) writeSizeBoundaries(min, max int64) {
	pio.WriteString(sb, strconv.FormatInt(min, 10))
	if max != dgo.UnboundedSize {
		pio.WriteByte(sb, ',')
		pio.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func (sb *typeBuilder) writeRangeDots(inclusive bool) {
	op := `...`
	if inclusive {
		op = `..`
	}
	pio.WriteString(sb, op)
}

func (sb *typeBuilder) writeRange(min, max dgo.Number, inclusive bool) {
	if min != nil {
		sb.buildTypeString(min.(dgo.Type), commaPrio)
	}
	sb.writeRangeDots(inclusive)
	if max != nil {
		sb.buildTypeString(max.(dgo.Type), commaPrio)
	}
}

func (sb *typeBuilder) writeTupleArgs(tt dgo.TupleType, leftSep, rightSep byte) {
	es := tt.ElementTypes()
	if tt.Variadic() {
		n := es.Len() - 1
		sep := leftSep
		for i := 0; i < n; i++ {
			pio.WriteByte(sb, sep)
			sep = ','
			sb.buildTypeString(es.Get(i).(dgo.Type), commaPrio)
		}
		pio.WriteByte(sb, sep)
		pio.WriteString(sb, `...`)
		sb.buildTypeString(es.Get(n).(dgo.Type), commaPrio)
		pio.WriteByte(sb, rightSep)
	} else {
		pio.WriteByte(sb, leftSep)
		sb.joinTypes(es, `,`, commaPrio)
		pio.WriteByte(sb, rightSep)
	}
}

func (sb *typeBuilder) writeTernary(typ dgo.Type, tc func(dgo.Value) dgo.Type, prio int, op string, opPrio int) {
	if prio >= orPrio {
		pio.WriteByte(sb, '(')
	}
	sb.joinX(typ.(dgo.TernaryType).Operands(), tc, op, opPrio)
	if prio >= orPrio {
		pio.WriteByte(sb, ')')
	}
}

func (sb *typeBuilder) buildTypeString(typ dgo.Type, prio int) {
	ti := typ.TypeIdentifier()
	if tn := sb.aliasMap.GetName(typ); tn != nil {
		pio.WriteString(sb, tn.GoString())
	} else if f, ok := complexTypes[ti]; ok {
		if util.RecursionHit(sb.seen, typ) {
			pio.WriteString(sb, `<recursive self reference to `)
			pio.WriteString(sb, ti.String())
			pio.WriteString(sb, ` type>`)
			return
		}
		os := sb.seen
		sb.seen = append(sb.seen, typ)
		f(sb, typ, prio)
		sb.seen = os
	} else {
		pio.WriteString(sb, ti.String())
	}
}

func typeAsType(v dgo.Value) dgo.Type {
	return v.(dgo.Type)
}

func valueAsType(v dgo.Value) dgo.Type {
	return v.Type()
}

// The use of an init() function is motivated by the internal package's need to call the TypeString function
// and the stringer package's need to use the internal package. The only way to break that dependency would be
// to merge  the two packages. There's however no way to use the internal package without the stringer package
// being initialized first so the circularity between them is harmless.
func init() {
	internal.TypeString = TypeString
	internal.TypeStringOn = TypeStringOn
}
