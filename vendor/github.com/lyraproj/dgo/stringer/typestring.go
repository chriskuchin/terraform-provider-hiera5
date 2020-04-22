package stringer

import (
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/lyraproj/dgo/internal"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

const (
	commaPrio = iota
	orPrio
	xorPrio
	andPrio
	typePrio
)

type typeToString func(typ dgo.Type, prio int)

type typeBuilder struct {
	io.Writer
	complexTypes map[dgo.TypeIdentifier]typeToString
	aliasMap     dgo.AliasMap
	seen         []dgo.Value
}

// TypeString produces a string with the go-like syntax for the given type.
func TypeString(typ dgo.Type) string {
	return TypeStringWithAliasMap(typ, internal.DefaultAliases())
}

// TypeStringWithAliasMap produces a string with the go-like syntax for the given type.
func TypeStringWithAliasMap(typ dgo.Type, am dgo.AliasMap) string {
	s := strings.Builder{}
	newTypeBuilder(&s, am).buildTypeString(typ, 0)
	return s.String()
}

func (sb *typeBuilder) anyOf(typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `|`, orPrio)
}

func (sb *typeBuilder) oneOf(typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `^`, xorPrio)
}

func (sb *typeBuilder) allOf(typ dgo.Type, prio int) {
	sb.writeTernary(typ, typeAsType, prio, `&`, andPrio)
}

func (sb *typeBuilder) allOfValue(typ dgo.Type, prio int) {
	sb.writeTernary(typ, valueAsType, prio, `&`, andPrio)
}

func (sb *typeBuilder) array(typ dgo.Type, _ int) {
	at := typ.(dgo.ArrayType)
	if at.Unbounded() {
		util.WriteString(sb, `[]`)
	} else {
		util.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(at.Min()), int64(at.Max()))
		util.WriteByte(sb, ']')
	}
	sb.buildTypeString(at.ElementType(), typePrio)
}

func (sb *typeBuilder) arrayExact(typ dgo.Type, _ int) {
	util.WriteByte(sb, '{')
	sb.joinValueTypes(typ.(dgo.ExactType).ExactValue().(dgo.Iterable), `,`, commaPrio)
	util.WriteByte(sb, '}')
}

func (sb *typeBuilder) binary(typ dgo.Type, _ int) {
	st := typ.(dgo.BinaryType)
	util.WriteString(sb, `binary`)
	if !st.Unbounded() {
		util.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(st.Min()), int64(st.Max()))
		util.WriteByte(sb, ']')
	}
}

func (sb *typeBuilder) binaryExact(typ dgo.Type, _ int) {
	util.WriteString(sb, `binary "`)
	util.WriteString(sb, typ.(dgo.ExactType).ExactValue().(dgo.Binary).Encode())
	util.WriteByte(sb, '"')
}

func (sb *typeBuilder) exactValue(typ dgo.Type, _ int) {
	util.WriteString(sb, typ.(dgo.ExactType).ExactValue().String())
}

func (sb *typeBuilder) tuple(typ dgo.Type, _ int) {
	sb.writeTupleArgs(typ.(dgo.TupleType), '{', '}')
}

func (sb *typeBuilder) _map(typ dgo.Type, _ int) {
	at := typ.(dgo.MapType)
	util.WriteString(sb, `map[`)
	sb.buildTypeString(at.KeyType(), commaPrio)
	if !at.Unbounded() {
		util.WriteByte(sb, ',')
		sb.writeSizeBoundaries(int64(at.Min()), int64(at.Max()))
	}
	util.WriteByte(sb, ']')
	sb.buildTypeString(at.ValueType(), typePrio)
}

func (sb *typeBuilder) mapExact(typ dgo.Type, _ int) {
	util.WriteByte(sb, '{')
	sb.joinValueTypes(typ.(dgo.ExactType).ExactValue().(dgo.Map), `,`, commaPrio)
	util.WriteByte(sb, '}')
}

func (sb *typeBuilder) _struct(typ dgo.Type, _ int) {
	util.WriteByte(sb, '{')
	st := typ.(dgo.StructMapType)
	sb.joinStructMapEntries(st)
	if st.Additional() {
		if st.Len() > 0 {
			util.WriteByte(sb, ',')
		}
		util.WriteString(sb, `...`)
	}
	util.WriteByte(sb, '}')
}

func (sb *typeBuilder) mapEntryExact(typ dgo.Type, _ int) {
	me := typ.(dgo.ExactType).ExactValue().(dgo.MapEntry)
	sb.buildTypeString(me.Key().Type(), commaPrio)
	util.WriteByte(sb, ':')
	sb.buildTypeString(me.Value().Type(), commaPrio)
}

func (sb *typeBuilder) floatRange(typ dgo.Type, _ int) {
	st := typ.(dgo.FloatType)
	sb.writeFloatRange(st.Min(), st.Max(), st.Inclusive())
}

func (sb *typeBuilder) integerRange(typ dgo.Type, _ int) {
	st := typ.(dgo.IntegerType)
	sb.writeIntRange(st.Min(), st.Max(), st.Inclusive())
}

func (sb *typeBuilder) regexpExact(typ dgo.Type, _ int) {
	util.WriteString(sb, typ.TypeIdentifier().String())
	util.WriteByte(sb, '[')
	util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).ExactValue().String()))
	util.WriteByte(sb, ']')
}

func (sb *typeBuilder) timeExact(typ dgo.Type, _ int) {
	util.WriteString(sb, typ.TypeIdentifier().String())
	util.WriteByte(sb, '[')
	util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).ExactValue().String()))
	util.WriteByte(sb, ']')
}

func (sb *typeBuilder) sensitive(typ dgo.Type, prio int) {
	util.WriteString(sb, `sensitive`)
	if op := typ.(dgo.UnaryType).Operand(); internal.DefaultAnyType != op {
		util.WriteByte(sb, '[')
		sb.buildTypeString(op, prio)
		util.WriteByte(sb, ']')
	}
}

func (sb *typeBuilder) stringExact(typ dgo.Type, _ int) {
	util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).ExactValue().String()))
}

func (sb *typeBuilder) stringPattern(typ dgo.Type, _ int) {
	internal.RegexpSlashQuote(sb, typ.(dgo.ExactType).ExactValue().String())
}

func (sb *typeBuilder) stringSized(typ dgo.Type, _ int) {
	st := typ.(dgo.StringType)
	util.WriteString(sb, `string`)
	if !st.Unbounded() {
		util.WriteByte(sb, '[')
		sb.writeSizeBoundaries(int64(st.Min()), int64(st.Max()))
		util.WriteByte(sb, ']')
	}
}

func (sb *typeBuilder) ciString(typ dgo.Type, _ int) {
	util.WriteByte(sb, '~')
	util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).ExactValue().String()))
}

func (sb *typeBuilder) native(typ dgo.Type, _ int) {
	rt := typ.(dgo.NativeType).GoType()
	util.WriteString(sb, `native`)
	if rt != nil {
		util.WriteByte(sb, '[')
		util.WriteString(sb, strconv.Quote(rt.String()))
		util.WriteByte(sb, ']')
	}
}

func (sb *typeBuilder) not(typ dgo.Type, _ int) {
	nt := typ.(dgo.UnaryType)
	util.WriteByte(sb, '!')
	sb.buildTypeString(nt.Operand(), typePrio)
}

func (sb *typeBuilder) meta(typ dgo.Type, prio int) {
	nt := typ.(dgo.UnaryType)
	util.WriteString(sb, `type`)
	if op := nt.Operand(); internal.DefaultAnyType != op {
		if op == nil {
			util.WriteString(sb, `[type]`)
		} else {
			util.WriteByte(sb, '[')
			sb.buildTypeString(op, prio)
			util.WriteByte(sb, ']')
		}
	}
}

func (sb *typeBuilder) function(typ dgo.Type, prio int) {
	ft := typ.(dgo.FunctionType)
	util.WriteString(sb, `func`)
	sb.writeTupleArgs(ft.In(), '(', ')')
	out := ft.Out()
	if out.Len() > 0 {
		util.WriteByte(sb, ' ')
		if out.Len() == 1 && !out.Variadic() {
			sb.buildTypeString(out.Element(0), prio)
		} else {
			sb.writeTupleArgs(ft.Out(), '(', ')')
		}
	}
}

func (sb *typeBuilder) errorExact(typ dgo.Type, _ int) {
	util.WriteString(sb, `error[`)
	util.WriteString(sb, strconv.Quote(typ.(dgo.ExactType).ExactValue().String()))
	util.WriteByte(sb, ']')
}

func (sb *typeBuilder) named(typ dgo.Type, _ int) {
	nt := typ.(dgo.NamedType)
	util.WriteString(sb, nt.Name())
	if params := nt.Parameters(); params != nil {
		util.WriteByte(sb, '[')
		sb.joinValueTypes(params, `,`, commaPrio)
		util.WriteByte(sb, ']')
	}
}

func newTypeBuilder(w io.Writer, am dgo.AliasMap) *typeBuilder {
	sb := &typeBuilder{Writer: w, aliasMap: am}
	sb.complexTypes = map[dgo.TypeIdentifier]typeToString{
		dgo.TiAnyOf:         sb.anyOf,
		dgo.TiOneOf:         sb.oneOf,
		dgo.TiAllOf:         sb.allOf,
		dgo.TiAllOfValue:    sb.allOfValue,
		dgo.TiArray:         sb.array,
		dgo.TiArrayExact:    sb.arrayExact,
		dgo.TiBinary:        sb.binary,
		dgo.TiBinaryExact:   sb.binaryExact,
		dgo.TiBooleanExact:  sb.exactValue,
		dgo.TiTuple:         sb.tuple,
		dgo.TiMap:           sb._map,
		dgo.TiMapExact:      sb.mapExact,
		dgo.TiMapEntryExact: sb.mapEntryExact,
		dgo.TiStruct:        sb._struct,
		dgo.TiFloatExact:    sb.exactValue,
		dgo.TiFloatRange:    sb.floatRange,
		dgo.TiIntegerExact:  sb.exactValue,
		dgo.TiIntegerRange:  sb.integerRange,
		dgo.TiRegexpExact:   sb.regexpExact,
		dgo.TiTimeExact:     sb.timeExact,
		dgo.TiSensitive:     sb.sensitive,
		dgo.TiStringExact:   sb.stringExact,
		dgo.TiStringPattern: sb.stringPattern,
		dgo.TiStringSized:   sb.stringSized,
		dgo.TiCiString:      sb.ciString,
		dgo.TiNative:        sb.native,
		dgo.TiNot:           sb.not,
		dgo.TiMeta:          sb.meta,
		dgo.TiFunction:      sb.function,
		dgo.TiErrorExact:    sb.errorExact,
		dgo.TiNamed:         sb.named,
		dgo.TiNamedExact:    sb.exactValue,
	}
	return sb
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
			util.WriteString(sb, s)
		}
		sb.buildTypeString(tc(v), prio)
	})
}

func (sb *typeBuilder) joinStructMapEntries(v dgo.StructMapType) {
	first := true
	v.Each(func(e dgo.StructMapEntry) {
		if first {
			first = false
		} else {
			util.WriteByte(sb, ',')
		}
		sb.buildTypeString(e.Key().(dgo.Type), commaPrio)
		if !e.Required() {
			util.WriteByte(sb, '?')
		}
		util.WriteByte(sb, ':')
		sb.buildTypeString(e.Value().(dgo.Type), commaPrio)
	})
}

func (sb *typeBuilder) writeSizeBoundaries(min, max int64) {
	util.WriteString(sb, strconv.FormatInt(min, 10))
	if max != math.MaxInt64 {
		util.WriteByte(sb, ',')
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func (sb *typeBuilder) writeIntRange(min, max int64, inclusive bool) {
	if min != math.MinInt64 {
		util.WriteString(sb, strconv.FormatInt(min, 10))
	}
	op := `...`
	if inclusive {
		op = `..`
	}
	util.WriteString(sb, op)
	if max != math.MaxInt64 {
		util.WriteString(sb, strconv.FormatInt(max, 10))
	}
}

func (sb *typeBuilder) writeFloatRange(min, max float64, inclusive bool) {
	if min != -math.MaxFloat64 {
		util.WriteString(sb, util.Ftoa(min))
	}
	op := `...`
	if inclusive {
		op = `..`
	}
	util.WriteString(sb, op)
	if max != math.MaxFloat64 {
		util.WriteString(sb, util.Ftoa(max))
	}
}

func (sb *typeBuilder) writeTupleArgs(tt dgo.TupleType, leftSep, rightSep byte) {
	es := tt.ElementTypes()
	if tt.Variadic() {
		n := es.Len() - 1
		sep := leftSep
		for i := 0; i < n; i++ {
			util.WriteByte(sb, sep)
			sep = ','
			sb.buildTypeString(es.Get(i).(dgo.Type), commaPrio)
		}
		util.WriteByte(sb, sep)
		util.WriteString(sb, `...`)
		sb.buildTypeString(es.Get(n).(dgo.Type), commaPrio)
		util.WriteByte(sb, rightSep)
	} else {
		util.WriteByte(sb, leftSep)
		sb.joinTypes(es, `,`, commaPrio)
		util.WriteByte(sb, rightSep)
	}
}

func (sb *typeBuilder) writeTernary(typ dgo.Type, tc func(dgo.Value) dgo.Type, prio int, op string, opPrio int) {
	if prio >= orPrio {
		util.WriteByte(sb, '(')
	}
	sb.joinX(typ.(dgo.TernaryType).Operands(), tc, op, opPrio)
	if prio >= orPrio {
		util.WriteByte(sb, ')')
	}
}

func (sb *typeBuilder) buildTypeString(typ dgo.Type, prio int) {
	if tn := sb.aliasMap.GetName(typ); tn != nil {
		util.WriteString(sb, tn.GoString())
		return
	}

	ti := typ.TypeIdentifier()
	if f, ok := sb.complexTypes[ti]; ok {
		if util.RecursionHit(sb.seen, typ) {
			util.WriteString(sb, `<recursive self reference to `)
			util.WriteString(sb, ti.String())
			util.WriteString(sb, ` type>`)
			return
		}
		os := sb.seen
		sb.seen = append(sb.seen, typ)
		f(typ, prio)
		sb.seen = os
	} else {
		util.WriteString(sb, ti.String())
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
}
