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
	println("sizeof<Pointer>", unsafe.Sizeof(Pointer{}))
	p := GetPointerMut(128)
	p.SetBytes(8, []byte("hello"))

	fmt.Println(p.Substr(8, 5))

	b := p.Bytes()
	fmt.Println(b)
}
