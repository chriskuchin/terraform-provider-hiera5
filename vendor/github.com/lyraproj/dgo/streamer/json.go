package streamer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
)

const (
	firstInArray = iota
	firstInObject
	afterElement
	afterValue
	afterKey
)

// MarshalJSON returns the JSON encoding for the given dgo.Value
func MarshalJSON(v interface{}, dialect Dialect) []byte {
	b := bytes.Buffer{}
	opts := DefaultOptions()
	if dialect != nil {
		opts.Dialect = dialect
	}
	New(nil, opts).Stream(vf.Value(v), JSON(&b))
	return b.Bytes()
}

// UnmarshalJSON decodes the JSON representation of the given bytes into a dgo.Value. The order of entries
// in an object is retained in its corresponding dgo.Map and rich data constructs such as Sensitive and Timestamp are
// converted.
func UnmarshalJSON(b []byte, dialect Dialect) dgo.Value {
	// Using an explicit decoder enables setting the UseNumber() attribute which in turn
	// will allow the vf.Value() method to perform the actual decoding of that number and
	// turn it into an int64 or a float64 depending on the if the string representation can
	// be parsed into an integer or not.
	je := json.NewDecoder(bytes.NewReader(b))
	je.UseNumber()

	opts := DefaultOptions()
	if dialect != nil {
		opts.Dialect = dialect
	}
	vc := DataDecoder(nil, opts.Dialect)

	j := &jsonDecoder{consumer: vc, refKey: opts.Dialect.RefKey().GoString(), decoder: je}
	j.decode()
	return vc.Value()
}

// jsonDecoder decodes a json stream into a dgo.Value. It retains the order of maps and
// resolves references.
type jsonDecoder struct {
	consumer Consumer
	refKey   string
	decoder  *json.Decoder
	pbToken  json.Token
}

func (j *jsonDecoder) decode() {
	j.decodeElem(json.Delim(0))
}

func (j *jsonDecoder) decodeElem(end json.Delim) bool {
	switch t := j.nextToken().(type) {
	case json.Delim:
		if t == end {
			return false
		}
		j.decodeCollection(t)
	case string:
		j.consumer.Add(vf.String(t))
	case json.Number:
		if i, err := t.Int64(); err == nil {
			j.consumer.Add(vf.Integer(i))
		} else {
			f, _ := t.Float64()
			j.consumer.Add(vf.Float(f))
		}
	case bool:
		j.consumer.Add(vf.Boolean(t))
	default:
		j.consumer.Add(vf.Nil)
	}
	return true
}

func (j *jsonDecoder) decodeCollection(delim json.Delim) {
	if delim == json.Delim('{') {
		k := j.nextToken()
		if k != j.refKey {
			j.pbToken = k
			j.consumer.AddMap(0, j.decodeMap)
		} else {
			if n, ok := j.nextToken().(json.Number); ok {
				if ri, err := n.Int64(); err == nil {
					j.consumer.AddRef(int(ri))
					if j.nextToken() != json.Delim('}') {
						panic(fmt.Errorf(`expected end of object after "%s": %d`, j.refKey, ri))
					}
					return
				}
			}
			panic(fmt.Errorf(`expected integer after key "%s"`, j.refKey))
		}
	} else {
		j.consumer.AddArray(0, j.decodeArray)
	}
}

func (j *jsonDecoder) decodeArray() {
	for j.decodeElem(json.Delim(']')) {
	}
}

func (j *jsonDecoder) decodeMap() {
	for j.decodeElem(json.Delim('}')) {
	}
}

func (j *jsonDecoder) nextToken() (t json.Token) {
	if j.pbToken != nil {
		t = j.pbToken
		j.pbToken = nil
	} else {
		var err error
		t, err = j.decoder.Token()
		if err != nil {
			if io.EOF == err {
				err = io.ErrUnexpectedEOF
			}
			panic(err)
		}
	}
	return
}

// JSON creates a new Consumer encode everything into JSON
func JSON(out io.Writer) Consumer {
	return &jsonEncoder{out: out, state: firstInArray, dialect: DgoDialect()}
}

type jsonEncoder struct {
	out     io.Writer
	dialect Dialect
	state   int
}

func (j *jsonEncoder) AddArray(_ int, doer dgo.Doer) {
	j.delimit(func() {
		j.state = firstInArray
		assertOk(j.out.Write([]byte{'['}))
		doer()
		assertOk(j.out.Write([]byte{']'}))
		j.state = afterElement
	})
}

func (j *jsonEncoder) AddMap(_ int, doer dgo.Doer) {
	j.delimit(func() {
		assertOk(j.out.Write([]byte{'{'}))
		j.state = firstInObject
		doer()
		assertOk(j.out.Write([]byte{'}'}))
		j.state = afterElement
	})
}

func (j *jsonEncoder) Add(element dgo.Value) {
	j.delimit(func() {
		j.write(element)
	})
}

func (j *jsonEncoder) AddRef(ref int) {
	j.delimit(func() {
		assertOk(fmt.Fprintf(j.out, `{"%s":%d}`, j.dialect.RefKey(), ref))
	})
}

func (j *jsonEncoder) CanDoBinary() bool {
	return false
}

func (j *jsonEncoder) CanDoComplexKeys() bool {
	return false
}

func (j *jsonEncoder) CanDoTime() bool {
	return false
}

func (j *jsonEncoder) StringDedupThreshold() int {
	return 20
}

func (j *jsonEncoder) delimit(doer dgo.Doer) {
	switch j.state {
	case firstInArray:
		doer()
		j.state = afterElement
	case firstInObject:
		doer()
		j.state = afterKey
	case afterKey:
		assertOk(j.out.Write([]byte{':'}))
		doer()
		j.state = afterValue
	case afterValue:
		assertOk(j.out.Write([]byte{','}))
		doer()
		j.state = afterKey
	default: // Element
		assertOk(j.out.Write([]byte{','}))
		doer()
	}
}

func (j *jsonEncoder) write(e dgo.Value) {
	var v []byte
	var err error
	switch e := e.(type) {
	case dgo.String:
		v, err = json.Marshal(e.GoString())
	case dgo.Float:
		v, err = json.Marshal(e.GoFloat())
	case dgo.Integer:
		v, err = json.Marshal(e.GoInt())
	case dgo.Boolean:
		v, err = json.Marshal(e.GoBool())
	default:
		v = []byte(`null`)
	}
	assertOk(0, err)
	assertOk(j.out.Write(v))
}

func assertOk(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
