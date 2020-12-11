package vf

import (
	"io"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Binary creates a new Binary based on the given slice. If frozen is true, the
// binary will be immutable and contain a copy of the slice, otherwise the slice
// is simply wrapped and modifications to its elements will also modify the binary.
func Binary(bs []byte, frozen bool) dgo.Binary {
	return internal.Binary(bs, frozen)
}

// BinaryFromString creates a new Binary from the given string using strict UTF8 encoding
func BinaryFromString(s string) dgo.Binary {
	return internal.BinaryFromEncoded(s, `%B`)
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
func BinaryFromEncoded(s, enc string) dgo.Binary {
	return internal.BinaryFromEncoded(s, enc)
}

// BinaryFromData creates a new frozen Binary based on data read from the given io.Reader.
func BinaryFromData(data io.Reader) dgo.Binary {
	return internal.BinaryFromData(data)
}
