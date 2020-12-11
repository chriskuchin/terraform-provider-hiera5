// Package yaml contains the Marshal and Unmarshal functions
package yaml

import (
	"fmt"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
	y3 "gopkg.in/yaml.v3"
)

// Marshal decodes the YAML representation of the given bytes into a dgo.Value
func Marshal(v dgo.Value) (bytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			if ye, ok := r.(yamlError); ok {
				err = ye.error
			} else {
				panic(r)
			}
		}
	}()
	bytes, err = y3.Marshal(yamlEncodeValue(v))
	return
}

func yamlEncodeValue(v dgo.Value) (nv *y3.Node) {
	switch v := v.(type) {
	case dgo.Array:
		nv = encodeArray(v)
	case dgo.Binary:
		nv = encodeBinary(v)
	case dgo.Boolean:
		nv = encodeBoolean(v)
	case dgo.Float:
		nv = encodeFloat(v)
	case dgo.Integer:
		nv = encodeInteger(v)
	case dgo.Struct:
		nv = encodeStruct(v)
	case dgo.Map:
		nv = encodeMap(v)
	case dgo.Native:
		nv = encodeNative(v)
	case dgo.Nil:
		nv = encodeNil()
	case dgo.String:
		nv = encodeString(v)
	case dgo.Time:
		nv = encodeTime(v)
	case dgo.Type:
		nv = encodeType(v)
	default:
		panic(yamlError{fmt.Errorf(`unable to marshal into value of type %v`, v.Type())})
	}
	return
}

func encodeArray(v dgo.Array) *y3.Node {
	s := make([]*y3.Node, v.Len())
	v.EachWithIndex(func(e dgo.Value, i int) {
		s[i] = yamlEncodeValue(e)
	})
	return &y3.Node{Kind: y3.SequenceNode, Tag: `!!seq`, Content: s}
}

func encodeBinary(v dgo.Binary) *y3.Node {
	return &y3.Node{Kind: y3.ScalarNode, Tag: `!!binary`, Value: v.String()}
}

func encodeBoolean(v dgo.Boolean) *y3.Node {
	return &y3.Node{Kind: y3.ScalarNode, Tag: `!!bool`, Value: v.String()}
}

func encodeFloat(v dgo.Float) *y3.Node {
	return &y3.Node{Kind: y3.ScalarNode, Tag: `!!float`, Value: v.String()}
}

func encodeInteger(v dgo.Integer) *y3.Node {
	return &y3.Node{Kind: y3.ScalarNode, Tag: `!!int`, Value: v.String()}
}

// encodeMap returns a *yaml.Node that represents the given map.
func encodeMap(v dgo.Map) *y3.Node {
	s := make([]*y3.Node, v.Len()*2)
	i := 0
	v.EachEntry(func(e dgo.MapEntry) {
		s[i] = yamlEncodeValue(e.Key())
		i++
		s[i] = yamlEncodeValue(e.Value())
		i++
	})
	return &y3.Node{Kind: y3.MappingNode, Tag: `!!map`, Content: s}
}

func encodeNative(n dgo.Native) *y3.Node {
	iv := n.GoValue()
	if ym, ok := iv.(y3.Marshaler); ok {
		yv, err := ym.MarshalYAML()
		if err != nil {
			panic(yamlError{err})
		}
		if n, ok := yv.(*y3.Node); ok {
			return n
		}
		return yamlEncodeValue(vf.Value(yv))
	}
	panic(yamlError{fmt.Errorf(`unable to marshal into value of type %T`, iv)})
}

func encodeNil() *y3.Node {
	return &y3.Node{Kind: y3.ScalarNode, Tag: `!!null`, Value: `null`}
}

func encodeString(v dgo.String) *y3.Node {
	n := &y3.Node{}
	n.SetString(v.GoString())
	return n
}

func encodeStruct(v dgo.Struct) *y3.Node {
	// A bit wasteful but this is currently the only way to create a yaml.Node
	// from a struct
	n := &y3.Node{}
	b, err := y3.Marshal(v.GoStruct())
	if err == nil {
		err = y3.Unmarshal(b, n)
	}
	if err != nil {
		panic(yamlError{err})
	}
	// n is the document node at this point.
	return n.Content[0]
}

func encodeTime(v dgo.Time) *y3.Node {
	return &y3.Node{
		Kind:  y3.ScalarNode,
		Tag:   `!!timestamp`,
		Value: v.GoTime().Format(time.RFC3339Nano),
		Style: y3.TaggedStyle}
}

func encodeType(t dgo.Type) *y3.Node {
	return &y3.Node{Tag: `!puppet.com,2019:dgo/type`, Kind: y3.ScalarNode, Value: t.String()}
}
