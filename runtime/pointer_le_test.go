package runtime

import (
	"fmt"
	"testing"
	"unsafe"
)

type ListPointer struct {
	Len int16
}

func TestAlloc(t *testing.T) {
	fmt.Println(unsafe.Sizeof(ListPointer{}))
	p := GetPointerMut(128)
	p.SetBytes(8, []byte("hello"))

	fmt.Println(p.Substr(8, 5))

	b := p.Bytes()
	fmt.Println(b)
}
