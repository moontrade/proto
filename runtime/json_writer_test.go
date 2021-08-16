package runtime

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestJsonWriter(t *testing.T) {
	w := JsonWriter{
		W: Buffer{
			b: GetBytes(1024),
		},
	}

	w.RawByte('{')
	w.RawByte('"')
	w.RawString("id")
	w.RawByte('"')
	w.RawByte(':')
	w.Int64(10)
	w.RawByte('}')

	r := JsonLexer{Data: w.W.Take()}
	fmt.Println(*(*string)(unsafe.Pointer(&r.Data)))

	r.Delim('{')
}
