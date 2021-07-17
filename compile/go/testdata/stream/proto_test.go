package stream

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

type VarHeader struct {
	ID   int64
	ID2  int64
	Size uint16
}

func TestPage_Bytes(t *testing.T) {
	printLayout(VarHeader{ID: 1, ID2: 2, Size: 3})
	fmt.Println()
	printLayout(Block{})
	printLayout(BlockHeader{})

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, 10)
	fmt.Println(binary.LittleEndian.Uint64(b))
}

type streamID struct {
	sequence  [20]byte
	partition int64
}

func printLayout(t interface{}) {
	// First ask Go to give us some information about the MyData type
	typ := reflect.TypeOf(t)
	fmt.Printf("Struct is %d bytes long\n", typ.Size())
	// We can run through the fields in the structure in order
	n := typ.NumField()
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		fmt.Printf("%s at offset %v, size=%d, align=%d\n",
			field.Name, field.Offset, field.Type.Size(),
			field.Type.Align())
	}
}
