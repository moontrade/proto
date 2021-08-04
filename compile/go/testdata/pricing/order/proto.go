// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package order

import (
	"fmt"
	"github.com/moontrade/proto/compile/go/testdata/pricing"
	"io"
	"reflect"
	"unsafe"
)

type Order struct {
	candle pricing.Candle
	_      [1]byte // Padding
}

func (s *Order) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Order) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["candle"] = s.Candle()
	return m
}

func (s *Order) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[1]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 1 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Order) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[1]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Order) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[1]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Order) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[1]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Order) Read(b []byte) (n int, err error) {
	if len(b) < 1 {
		return -1, io.ErrShortBuffer
	}
	v := (*Order)(unsafe.Pointer(&b[0]))
	*v = *s
	return 1, nil
}
func (s *Order) UnmarshalBinary(b []byte) error {
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	v := (*Order)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Order) Clone() *Order {
	v := &Order{}
	*v = *s
	return v
}
func (s *Order) Bytes() []byte {
	return (*(*[1]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Order) Mut() *OrderMut {
	return (*OrderMut)(unsafe.Pointer(s))
}
func (s *Order) Candle() pricing.Candle {
	return s.candle
}

type OrderMut struct {
	Order
}

func (s *OrderMut) Clone() *OrderMut {
	v := &OrderMut{}
	*v = *s
	return v
}
func (s *OrderMut) Freeze() *Order {
	return (*Order)(unsafe.Pointer(s))
}
func (s *OrderMut) SetCandle(v pricing.Candle) *OrderMut {
	s.candle = v
	return s
}
func init() {
	{
		var b [2]byte
		v := uint16(1)
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		if *(*uint16)(unsafe.Pointer(&b[0])) != 1 {
			panic("BigEndian not supported")
		}
	}
	type b struct {
		n    string
		o, s uintptr
	}
	a := func(x interface{}, y interface{}, s uintptr, z []b) {
		t := reflect.TypeOf(x)
		r := reflect.TypeOf(y)
		if t.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", t.Name(), t.Size(), s))
		}
		if r.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", r.Name(), r.Size(), s))
		}
		if t.NumField() != len(z) {
			panic(fmt.Sprintf("%s field count = %d: expected %d", t.Name(), t.NumField(), len(z)))
		}
		for i, e := range z {
			f := t.Field(i)
			if f.Offset != e.o {
				panic(fmt.Sprintf("%s.%s offset = %d, expected = %d", t.Name(), f.Name, f.Offset, e.o))
			}
			if f.Type.Size() != e.s {
				panic(fmt.Sprintf("%s.%s size = %d, expected = %d", t.Name(), f.Name, f.Type.Size(), e.s))
			}
			if f.Name != e.n {
				panic(fmt.Sprintf("%s.%s expected field: %s", t.Name(), f.Name, e.n))
			}
		}
	}

	a(Order{}, OrderMut{}, 1, []b{
		{"candle", 0, 0},
		{"_", 0, 1},
	})

}
