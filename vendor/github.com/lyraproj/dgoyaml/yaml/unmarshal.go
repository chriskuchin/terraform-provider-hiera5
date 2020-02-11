package yaml

import (
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/vf"
	y3 "gopkg.in/yaml.v3"
)

// yamlError is used internally to panic with errors from the yaml package. It is recovered and returned at the top
type yamlError struct {
	error
}

// Unmarshal decodes the YAML representation of the given bytes into a dgo.Value
func Unmarshal(b []byte) (val dgo.Value, err error) {
	var n y3.Node
	if err = y3.Unmarshal(b, &n); err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			if ye, ok := r.(yamlError); ok {
				err = ye.error
			} else {
				panic(r)
			}
		}
	}()
	val = decodeValue(&n)
	return
}

func decodeScalar(n *y3.Node) dgo.Value {
	var v dgo.Value
	switch n.Tag {
	case `!!null`:
		v = vf.Nil
	case `!!bool`:
		var x bool
		_ = n.Decode(&x)
		v = vf.Boolean(x)
	case `!!int`:
		var x int64
		_ = n.Decode(&x)
		v = vf.Integer(x)
	case `!!float`:
		var x float64
		_ = n.Decode(&x)
		v = vf.Float(x)
	case `!!str`:
		v = vf.String(n.Value)
	case `!!timestamp`:
		var x time.Time
		if err := n.Decode(&x); err != nil {
			panic(yamlError{err})
		}
		v = vf.Time(x)
	case `!!binary`:
		v = vf.BinaryFromString(n.Value)
	case `!puppet.com,2019:dgo/type`:
		v = tf.Parse(n.Value)
	default:
		var x interface{}
		if err := n.Decode(&x); err != nil {
			panic(yamlError{err})
		}
		v = vf.Value(x)
	}
	return v
}

func decodeValue(n *y3.Node) dgo.Value {
	var v dgo.Value
	switch n.Kind {
	case y3.DocumentNode:
		v = decodeValue(n.Content[0])
	case y3.SequenceNode:
		v = decodeArray(n)
	case y3.MappingNode:
		v = decodeMap(n)
	default:
		v = decodeScalar(n)
	}
	return v
}

func decodeArray(n *y3.Node) dgo.Array {
	ms := n.Content
	es := make([]dgo.Value, len(ms))
	for i, me := range ms {
		es[i] = decodeValue(me)
	}
	return vf.WrapSlice(es)
}

func decodeMap(n *y3.Node) dgo.Map {
	ms := n.Content
	top := len(ms)
	m := vf.MapWithCapacity(top/8*6, nil)
	for i := 0; i < top; i += 2 {
		m.Put(decodeValue(ms[i]), decodeValue(ms[i+1]))
	}
	return m
}
