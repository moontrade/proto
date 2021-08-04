// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package charting

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type Plot struct {
	_ [1]byte // Padding
}

func (s *Plot) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Plot) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	return m
}

func (s *Plot) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[1]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 1 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Plot) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[1]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Plot) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[1]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Plot) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[1]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Plot) Read(b []byte) (n int, err error) {
	if len(b) < 1 {
		return -1, io.ErrShortBuffer
	}
	v := (*Plot)(unsafe.Pointer(&b[0]))
	*v = *s
	return 1, nil
}
func (s *Plot) UnmarshalBinary(b []byte) error {
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	v := (*Plot)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Plot) Clone() *Plot {
	v := &Plot{}
	*v = *s
	return v
}
func (s *Plot) Bytes() []byte {
	return (*(*[1]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Plot) Mut() *PlotMut {
	return (*PlotMut)(unsafe.Pointer(s))
}

type PlotMut struct {
	Plot
}

func (s *PlotMut) Clone() *PlotMut {
	v := &PlotMut{}
	*v = *s
	return v
}
func (s *PlotMut) Freeze() *Plot {
	return (*Plot)(unsafe.Pointer(s))
}

type Fill struct {
	_ [1]byte // Padding
}

func (s *Fill) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Fill) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	return m
}

func (s *Fill) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[1]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 1 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Fill) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[1]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Fill) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[1]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Fill) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[1]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Fill) Read(b []byte) (n int, err error) {
	if len(b) < 1 {
		return -1, io.ErrShortBuffer
	}
	v := (*Fill)(unsafe.Pointer(&b[0]))
	*v = *s
	return 1, nil
}
func (s *Fill) UnmarshalBinary(b []byte) error {
	if len(b) < 1 {
		return io.ErrShortBuffer
	}
	v := (*Fill)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Fill) Clone() *Fill {
	v := &Fill{}
	*v = *s
	return v
}
func (s *Fill) Bytes() []byte {
	return (*(*[1]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Fill) Mut() *FillMut {
	return (*FillMut)(unsafe.Pointer(s))
}

type FillMut struct {
	Fill
}

func (s *FillMut) Clone() *FillMut {
	v := &FillMut{}
	*v = *s
	return v
}
func (s *FillMut) Freeze() *Fill {
	return (*Fill)(unsafe.Pointer(s))
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

	a(Plot{}, PlotMut{}, 1, []b{
		{"_", 0, 1},
	})
	a(Fill{}, FillMut{}, 1, []b{
		{"_", 0, 1},
	})

}
