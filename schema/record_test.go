package schema

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestRecord(t *testing.T) {
	b := make([]byte, 16)
	p := PointerOf(b)
	p.SetInt16(10)
	fmt.Println(p.Int16())
	fmt.Println(unsafe.Sizeof(Pointer{}))
}

func BenchmarkNewPointerString(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			s := fmt.Sprintf("%d", i)
			PointerOfString(s)
		}
	})
	b.Run("bytes", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			s := []byte(fmt.Sprintf("%d", i))
			PointerOf(s)
		}
	})
}

type Ptr interface {
	Int16(vp VPointer) int16
}

type PtrMut interface {
	Ptr

	SetInt16(vp VPointer) int16
}

type NativePointer struct {
	ptr uintptr
	len int
}

type Pointer struct {
	ptr unsafe.Pointer
	len int
}

func (p Pointer) Grow(by int) Pointer {
	dst := make([]byte, p.len+by)
	src := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr),
		Len:  p.len,
		Cap:  p.len,
	}))
	copy(dst, src)
	return PointerOf(dst)
}

func Alloc(size int) Pointer {
	b := make([]byte, size)
	return Pointer{
		ptr: unsafe.Pointer(&b[0]),
		len: len(b),
	}
}

func PointerOf(b []byte) Pointer {
	return Pointer{
		ptr: unsafe.Pointer(&b[0]),
		len: len(b),
	}
}

func PointerOfString(s string) Pointer {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return Pointer{
		ptr: unsafe.Pointer(h.Data),
		len: len(s),
	}
}

func (p Pointer) IsNil() bool {
	return uintptr(p.ptr) == 0
}

func (p Pointer) IsEmpty() bool {
	return uintptr(p.ptr) == 0 || p.len == 0
}

func (p Pointer) Bounds(size int) bool {
	return uintptr(p.ptr) == 0 || p.len < size
}

func (p Pointer) String() string {
	if p.IsEmpty() {
		return EmptyString
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(p.ptr),
		Len:  int(p.len),
	}))
}

func (p Pointer) UInt16() uint16 {
	if p.Bounds(2) {
		return 0
	}
	return *(*uint16)(p.ptr)
}

func (p Pointer) SetInt16(value int16) {
	if p.Bounds(2) {
		return
	}
	*(*int16)(p.ptr) = value
}

func (p Pointer) Int16() int16 {
	if p.Bounds(2) {
		return 0
	}
	return *(*int16)(p.ptr)
}

func (p Pointer) At(offset VPointer) Pointer {
	if p.IsNil() {
		return Pointer{}
	}
	if int(offset.offset+offset.size) > p.len {
		return Pointer{}
	}
	return Pointer{
		ptr: unsafe.Pointer(uintptr(p.ptr) + uintptr(offset.offset)),
		len: int(offset.size),
	}
}

type Order struct {
	Pointer
}

const (
	EmptyString = ""
)

type VPointer struct {
	offset int32
	size   int32
}
