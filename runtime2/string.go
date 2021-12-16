package runtime2

import (
	"io"
	"unsafe"
)

type String32 [32]byte

func NewString32(s string) *String32 {
	v := String32{}
	v.set(s)
	return &v
}
func (s *String32) set(v string) {
	copy(s[0:31], v)
	c := 31
	l := len(v)
	if l > c {
		s[31] = byte(c)
	} else {
		s[31] = byte(l)
	}
}
func (s *String32) Len() int {
	return int(s[31])
}
func (s *String32) Cap() int {
	return 31
}
func (s *String32) StringClone() string {
	b := s[0:s.Len()]
	return string(b)
}
func (s *String32) String() string {
	b := s[0:s.Len()]
	return *(*string)(unsafe.Pointer(&b))
}
func (s *String32) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String32) Clone() *String32 {
	v := String32{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String32) Mut() *String32Mut {
	return *(**String32Mut)(unsafe.Pointer(&s))
}
func (s *String32) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 32 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String32) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[32]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String32) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String32) MarshalBinary() ([]byte, error) {
	return s[0:s.Len()], nil
}
func (s *String32) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*String32)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String32Mut struct {
	String32
}

func (s *String32Mut) Set(v string) {
	s.set(v)
}
