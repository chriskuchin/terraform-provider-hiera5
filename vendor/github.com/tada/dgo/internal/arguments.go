package internal

import (
	"github.com/tada/dgo/dgo"
)

type arguments struct {
	arrayFrozen
}

// Arguments returns an immutable Arguments instance that represents the given slice
func Arguments(values []interface{}) dgo.Arguments {
	return ArgumentsFromArray(Values(values))
}

// ArgumentsFromArray returns an Arguments instance backed by the given array
func ArgumentsFromArray(values dgo.Array) dgo.Arguments {
	a := values.Copy(true).(*arrayFrozen)
	return &arguments{*a}
}

func (a *arguments) AssertSize(funcName string, min, max int) {
	l := a.Len()
	if min > l || l > max {
		panic(illegalArgumentCount(funcName, min, max, l))
	}
}

func (a *arguments) Arg(funcName string, n int, typ dgo.Type) dgo.Value {
	v := a.Get(n)
	if typ.Instance(v) {
		return v
	}
	panic(illegalArgument(funcName, typ, a.InterfaceSlice(), n))
}

func (a *arguments) Equals(other interface{}) bool {
	oa, ok := other.(*arguments)
	return ok && a.array.Equals(&oa.array)
}
