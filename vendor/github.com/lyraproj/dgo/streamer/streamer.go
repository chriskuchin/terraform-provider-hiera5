// Package streamer contains the logic to convert dgo values into sequences of data and vice versa.
package streamer

import (
	"fmt"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
)

const (
	// NoDedup effectively prevents any attempt to do data deduplication
	NoDedup = DedupLevel(iota)

	// NoKeyDedup prevents that keys in maps are deduplicated. This will only affect consumers that
	// can handle keys other than strings. This is the default deduplication level
	NoKeyDedup

	// MaxDedup will cause deduplication of both keys and values
	MaxDedup
)

type (
	// DedupLevel controls the level of deduplication that will occur during serialization
	DedupLevel int

	// Options controls some aspects of the Streamer.
	Options struct {
		DedupLevel DedupLevel
		Dialect    Dialect
		RichData   bool
	}

	// Streamer is a re-entrant fully configured serializer that streams the given
	// value to the given consumer.
	Streamer interface {
		// Stream the given RichData value to a series of values streamed to the
		// given consumer.
		Stream(value dgo.Value, consumer Consumer)
	}

	rdSerializer struct {
		Options
		aliasMap dgo.AliasMap
	}

	context struct {
		config     *rdSerializer
		values     map[interface{}]int
		refIndex   int
		dedupLevel DedupLevel
		consumer   Consumer
	}
)

// DefaultOptions returns the default options for the streamer. The returned value is a private copy that can be
// modified by the caller before it is passed on to a streamer.
func DefaultOptions() *Options {
	return &Options{
		DedupLevel: NoKeyDedup,
		Dialect:    DgoDialect(),
		RichData:   true}
}

// New returns a new Streamer
func New(ctx dgo.AliasMap, options *Options) Streamer {
	if options == nil {
		options = DefaultOptions()
	}
	return &rdSerializer{Options: *options, aliasMap: ctx}
}

func (t *rdSerializer) Stream(value dgo.Value, consumer Consumer) {
	c := context{config: t, values: make(map[interface{}]int, 31), refIndex: 0, consumer: consumer, dedupLevel: t.DedupLevel}
	if c.dedupLevel >= MaxDedup && !consumer.CanDoComplexKeys() {
		c.dedupLevel = NoKeyDedup
	}
	c.emitData(value)
}

func (sc *context) emitDataNoDedup(value dgo.Value) {
	dl := sc.dedupLevel
	sc.dedupLevel = NoDedup
	sc.emitData(value)
	sc.dedupLevel = dl
}

func (sc *context) emitData(value dgo.Value) {
	if value == nil || value == vf.Nil {
		sc.addData(vf.Nil)
		return
	}

	switch value := value.(type) {
	case dgo.Integer, dgo.Float, dgo.Boolean:
		// Never dedup
		sc.addData(value)
	case dgo.String:
		sc.emitString(value)
	case dgo.Map:
		sc.emitMap(value)
	case dgo.Array:
		sc.emitArray(value)
	case dgo.Sensitive:
		sc.emitSensitive(value)
	case dgo.Binary:
		sc.emitBinary(value)
	case dgo.Time:
		sc.emitTime(value)
	case dgo.Type:
		sc.emitType(value)
	default:
		if sc.config.RichData {
			if nt, ok := value.Type().(dgo.NamedType); ok {
				sc.emitNamed(nt, value)
				break
			}
		}
		panic(sc.unknownSerialization(value))
	}
}

func (sc *context) emitNamed(t dgo.NamedType, value dgo.Value) {
	sc.process(value, func() {
		v := t.ExtractInitArg(value)
		d := sc.config.Dialect
		sc.addMap(2, func() {
			sc.addData(d.TypeKey())
			sc.addData(vf.String(t.Name()))
			if vm, ok := v.(dgo.Map); ok {
				vm.EachEntry(func(e dgo.MapEntry) {
					sc.addData(e.Key())
					sc.emitData(e.Value())
				})
			} else {
				sc.addData(d.ValueKey())
				sc.emitData(v)
			}
		})
	})
}

func (sc *context) unknownSerialization(value dgo.Value) error {
	return fmt.Errorf(`unable to serialize value of type %s`, value.Type())
}

func (sc *context) process(value interface{}, doer dgo.Doer) {
	if sc.dedupLevel == NoDedup {
		doer()
		return
	}

	if ref, ok := sc.values[value]; ok {
		sc.consumer.AddRef(ref)
	} else {
		sc.values[value] = sc.refIndex
		doer()
	}
}

func (sc *context) addData(v dgo.Value) {
	sc.refIndex++
	sc.consumer.Add(v)
}

func (sc *context) emitString(value dgo.String) {
	// Dedup only if length exceeds stringThreshold
	str := value.GoString()
	if len(str) >= sc.consumer.StringDedupThreshold() {
		sc.process(str, func() {
			sc.addData(value)
		})
	} else {
		sc.addData(value)
	}
}

func (sc *context) emitType(typ dgo.Type) {
	sc.process(typ, func() {
		sc.addMap(2, func() {
			d := sc.config.Dialect
			sc.addData(d.TypeKey())
			if am := sc.config.aliasMap; am != nil {
				if tn := am.GetName(typ); tn != nil {
					sc.addData(d.AliasTypeName())
					sc.addData(d.ValueKey())
					sc.emitData(vf.Values(tn, typ.String()))
				}
			} else {
				sc.addData(vf.String(typ.String()))
			}
		})
	})
}

func (sc *context) addArray(len int, doer dgo.Doer) {
	sc.refIndex++
	sc.consumer.AddArray(len, doer)
}

func (sc *context) addMap(len int, doer dgo.Doer) {
	sc.refIndex++
	sc.consumer.AddMap(len, doer)
}

func (sc *context) emitArray(value dgo.Array) {
	sc.process(value, func() {
		sc.addArray(value.Len(), func() {
			value.EachWithIndex(func(elem dgo.Value, idx int) {
				sc.emitData(elem)
			})
		})
	})
}

func (sc *context) emitTime(value dgo.Time) {
	sc.process(value, func() {
		if sc.consumer.CanDoTime() {
			sc.addData(value)
		} else {
			if !sc.config.RichData {
				panic(sc.unknownSerialization(value))
			}
			sc.addMap(2, func() {
				d := sc.config.Dialect
				sc.addData(d.TypeKey())
				sc.addData(d.TimeTypeName())
				sc.addData(d.ValueKey())
				sc.emitData(vf.String(value.String()))
			})
		}
	})
}

func (sc *context) emitBinary(value dgo.Binary) {
	sc.process(value, func() {
		if sc.consumer.CanDoBinary() {
			sc.addData(value)
		} else {
			if !sc.config.RichData {
				panic(sc.unknownSerialization(value))
			}
			sc.addMap(2, func() {
				d := sc.config.Dialect
				sc.addData(d.TypeKey())
				sc.addData(d.BinaryTypeName())
				sc.addData(d.ValueKey())
				sc.emitData(vf.String(value.String()))
			})
		}
	})
}

func (sc *context) emitMap(value dgo.Map) {
	if sc.consumer.CanDoComplexKeys() || value.StringKeys() {
		sc.process(value, func() {
			sc.addMap(value.Len(), func() {
				value.EachEntry(func(e dgo.MapEntry) {
					if sc.dedupLevel == NoKeyDedup {
						sc.emitDataNoDedup(e.Key())
					} else {
						sc.emitData(e.Key())
					}
					sc.emitData(e.Value())
				})
			})
		})
		return
	}
	if !sc.config.RichData {
		panic(sc.unknownSerialization(value))
	}
	sc.process(value, func() {
		sc.addMap(2, func() {
			d := sc.config.Dialect
			sc.addData(d.TypeKey())
			sc.addData(d.MapTypeName())
			sc.addData(d.ValueKey())
			sc.addArray(value.Len()*2, func() {
				value.EachEntry(func(e dgo.MapEntry) {
					sc.emitDataNoDedup(e.Key())
					sc.emitData(e.Value())
				})
			})
		})
	})
}

func (sc *context) emitSensitive(value dgo.Sensitive) {
	sc.process(value, func() {
		if !sc.config.RichData {
			panic(sc.unknownSerialization(value))
		}
		sc.addMap(2, func() {
			d := sc.config.Dialect
			sc.addData(d.TypeKey())
			sc.addData(d.SensitiveTypeName())
			sc.addData(d.ValueKey())
			sc.emitData(value.Unwrap())
		})
	})
}
