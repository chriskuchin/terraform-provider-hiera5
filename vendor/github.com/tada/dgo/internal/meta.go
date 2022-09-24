package internal

import (
	"reflect"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

var reflectTypeType = reflect.TypeOf((*dgo.Type)(nil)).Elem()

// metaType is the Type returned by a Type
type metaType struct {
	tp dgo.Type
}

// DefaultMetaType is the unconstrained meta type
var DefaultMetaType = &metaType{tp: DefaultAnyType}

// MetaType creates the meta type for the given type
func MetaType(t dgo.Type) dgo.Meta {
	return &metaType{t}
}

func (t *metaType) Type() dgo.Type {
	if t.tp == nil {
		return t // type of meta type is meta type
	}
	return &metaType{nil} // Short circuit meta chain
}

func (t *metaType) Assignable(ot dgo.Type) bool {
	if mt, ok := ot.(*metaType); ok {
		if t.tp == nil {
			// Only MetaTypeType is assignable to MetaTypeType
			return mt.tp == nil
		}
		return t.tp.Equals(mt.tp)
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t *metaType) Describes() dgo.Type {
	return t.tp
}

func (t *metaType) Equals(v interface{}) bool {
	if mt, ok := v.(*metaType); ok {
		if t.tp == nil {
			return mt.tp == nil
		}
		return t.tp.Equals(mt.tp)
	}
	return false
}

func (t *metaType) HashCode() dgo.Hash {
	h := dgo.Hash(dgo.TiMeta) * 1321
	if t.tp != nil {
		h += t.tp.HashCode()
	}
	return h
}

func (t *metaType) Instance(v interface{}) bool {
	if ot, ok := v.(dgo.Type); ok {
		if t.tp == nil {
			// MetaTypeType
			_, ok = ot.(*metaType)
			return ok
		}
		return t.tp.Assignable(ot)
	}
	return false
}

func (t *metaType) New(arg dgo.Value) dgo.Value {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`type`, 1, 1)
		arg = args.Get(0)
	}
	var tv dgo.Type
	if s, ok := arg.(dgo.String); ok {
		tv = AsType(Parse(s.GoString()))
	} else {
		tv = AsType(arg)
	}
	if !t.Instance(tv) {
		panic(catch.Error(IllegalAssignment(t, tv)))
	}
	return tv
}

func (t *metaType) Operator() dgo.TypeOp {
	return dgo.OpMeta
}

func (t *metaType) Operand() dgo.Type {
	return t.tp
}

func (t *metaType) ReflectType() reflect.Type {
	return reflectTypeType
}

func (t *metaType) Resolve(ap dgo.AliasAdder) {
	tp := t.tp
	t.tp = DefaultAnyType
	t.tp = ap.Replace(tp).(dgo.Type)
}

func (t *metaType) String() string {
	return TypeString(t)
}

func (t *metaType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMeta
}
