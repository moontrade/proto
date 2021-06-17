// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package stream

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type Stream struct {
	id      ID
	first   int64
	records int64
	pages   int64
	size    int64
	xsize   int64
}

func (s *Stream) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Stream) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id().MarshalMap(nil)
	m["first"] = s.First()
	m["records"] = s.Records()
	m["pages"] = s.Pages()
	m["size"] = s.Size()
	m["xsize"] = s.Xsize()
	return m
}

func (s *Stream) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[56]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 56 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Stream) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[56]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Stream) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[56]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Stream) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[56]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Stream) Read(b []byte) (n int, err error) {
	if len(b) < 56 {
		return -1, io.ErrShortBuffer
	}
	v := (*Stream)(unsafe.Pointer(&b[0]))
	*v = *s
	return 56, nil
}
func (s *Stream) UnmarshalBinary(b []byte) error {
	if len(b) < 56 {
		return io.ErrShortBuffer
	}
	v := (*Stream)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Stream) Clone() *Stream {
	v := &Stream{}
	*v = *s
	return v
}
func (s *Stream) Bytes() []byte {
	return (*(*[56]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Stream) Mut() *StreamMut {
	return (*StreamMut)(unsafe.Pointer(s))
}
func (s *Stream) Id() *ID {
	return &s.id
}
func (s *Stream) First() int64 {
	return s.first
}
func (s *Stream) Records() int64 {
	return s.records
}
func (s *Stream) Pages() int64 {
	return s.pages
}
func (s *Stream) Size() int64 {
	return s.size
}
func (s *Stream) Xsize() int64 {
	return s.xsize
}

type StreamMut struct {
	Stream
}

func (s *StreamMut) Clone() *StreamMut {
	v := &StreamMut{}
	*v = *s
	return v
}
func (s *StreamMut) Freeze() *Stream {
	return (*Stream)(unsafe.Pointer(s))
}
func (s *StreamMut) Id() *IDMut {
	return s.id.Mut()
}
func (s *StreamMut) SetId(v *ID) *StreamMut {
	s.id = *v
	return s
}
func (s *StreamMut) SetFirst(v int64) *StreamMut {
	s.first = v
	return s
}
func (s *StreamMut) SetRecords(v int64) *StreamMut {
	s.records = v
	return s
}
func (s *StreamMut) SetPages(v int64) *StreamMut {
	s.pages = v
	return s
}
func (s *StreamMut) SetSize(v int64) *StreamMut {
	s.size = v
	return s
}
func (s *StreamMut) SetXsize(v int64) *StreamMut {
	s.xsize = v
	return s
}

type PageID struct {
	stream ID
	id     int64
}

func (s *PageID) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *PageID) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["stream"] = s.Stream().MarshalMap(nil)
	m["id"] = s.Id()
	return m
}

func (s *PageID) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *PageID) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *PageID) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *PageID) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *PageID) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*PageID)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *PageID) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*PageID)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *PageID) Clone() *PageID {
	v := &PageID{}
	*v = *s
	return v
}
func (s *PageID) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *PageID) Mut() *PageIDMut {
	return (*PageIDMut)(unsafe.Pointer(s))
}
func (s *PageID) Stream() *ID {
	return &s.stream
}
func (s *PageID) Id() int64 {
	return s.id
}

type PageIDMut struct {
	PageID
}

func (s *PageIDMut) Clone() *PageIDMut {
	v := &PageIDMut{}
	*v = *s
	return v
}
func (s *PageIDMut) Freeze() *PageID {
	return (*PageID)(unsafe.Pointer(s))
}
func (s *PageIDMut) Stream() *IDMut {
	return s.stream.Mut()
}
func (s *PageIDMut) SetStream(v *ID) *PageIDMut {
	s.stream = *v
	return s
}
func (s *PageIDMut) SetId(v int64) *PageIDMut {
	s.id = v
	return s
}

//
type Page struct {
	first    int64
	duration int64
	last     int64
	count    uint16
	record   uint16
	size     uint16
	xsize    uint16
	pad      String32
	data     String65472
}

func (s *Page) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Page) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["first"] = s.First()
	m["duration"] = s.Duration()
	m["last"] = s.Last()
	m["count"] = s.Count()
	m["record"] = s.Record()
	m["size"] = s.Size()
	m["xsize"] = s.Xsize()
	m["pad"] = s.Pad()
	m["data"] = s.Data()
	return m
}

func (s *Page) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[65536]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 65536 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Page) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[65536]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Page) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[65536]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Page) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[65536]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Page) Read(b []byte) (n int, err error) {
	if len(b) < 65536 {
		return -1, io.ErrShortBuffer
	}
	v := (*Page)(unsafe.Pointer(&b[0]))
	*v = *s
	return 65536, nil
}
func (s *Page) UnmarshalBinary(b []byte) error {
	if len(b) < 65536 {
		return io.ErrShortBuffer
	}
	v := (*Page)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Page) Clone() *Page {
	v := &Page{}
	*v = *s
	return v
}
func (s *Page) Bytes() []byte {
	return (*(*[65536]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Page) Mut() *PageMut {
	return (*PageMut)(unsafe.Pointer(s))
}
func (s *Page) First() int64 {
	return s.first
}
func (s *Page) Duration() int64 {
	return s.duration
}
func (s *Page) Last() int64 {
	return s.last
}
func (s *Page) Count() uint16 {
	return s.count
}
func (s *Page) Record() uint16 {
	return s.record
}
func (s *Page) Size() uint16 {
	return s.size
}
func (s *Page) Xsize() uint16 {
	return s.xsize
}
func (s *Page) Pad() *String32 {
	return &s.pad
}
func (s *Page) Data() *String65472 {
	return &s.data
}

//
type PageMut struct {
	Page
}

func (s *PageMut) Clone() *PageMut {
	v := &PageMut{}
	*v = *s
	return v
}
func (s *PageMut) Freeze() *Page {
	return (*Page)(unsafe.Pointer(s))
}
func (s *PageMut) SetFirst(v int64) *PageMut {
	s.first = v
	return s
}
func (s *PageMut) SetDuration(v int64) *PageMut {
	s.duration = v
	return s
}
func (s *PageMut) SetLast(v int64) *PageMut {
	s.last = v
	return s
}
func (s *PageMut) SetCount(v uint16) *PageMut {
	s.count = v
	return s
}
func (s *PageMut) SetRecord(v uint16) *PageMut {
	s.record = v
	return s
}
func (s *PageMut) SetSize(v uint16) *PageMut {
	s.size = v
	return s
}
func (s *PageMut) SetXsize(v uint16) *PageMut {
	s.xsize = v
	return s
}
func (s *PageMut) Pad() *String32Mut {
	return s.pad.Mut()
}
func (s *PageMut) SetPad(v *String32) *PageMut {
	s.pad = *v
	return s
}
func (s *PageMut) Data() *String65472Mut {
	return s.data.Mut()
}
func (s *PageMut) SetData(v *String65472) *PageMut {
	s.data = *v
	return s
}

type ID struct {
	partition  uint16
	record     uint16
	resolution uint32
	sequence   int64
}

func (s *ID) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *ID) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["partition"] = s.Partition()
	m["record"] = s.Record()
	m["resolution"] = s.Resolution()
	m["sequence"] = s.Sequence()
	return m
}

func (s *ID) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[16]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 16 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *ID) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[16]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *ID) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *ID) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *ID) Read(b []byte) (n int, err error) {
	if len(b) < 16 {
		return -1, io.ErrShortBuffer
	}
	v := (*ID)(unsafe.Pointer(&b[0]))
	*v = *s
	return 16, nil
}
func (s *ID) UnmarshalBinary(b []byte) error {
	if len(b) < 16 {
		return io.ErrShortBuffer
	}
	v := (*ID)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *ID) Clone() *ID {
	v := &ID{}
	*v = *s
	return v
}
func (s *ID) Bytes() []byte {
	return (*(*[16]byte)(unsafe.Pointer(s)))[0:]
}
func (s *ID) Mut() *IDMut {
	return (*IDMut)(unsafe.Pointer(s))
}
func (s *ID) Partition() uint16 {
	return s.partition
}
func (s *ID) Record() uint16 {
	return s.record
}
func (s *ID) Resolution() uint32 {
	return s.resolution
}
func (s *ID) Sequence() int64 {
	return s.sequence
}

type IDMut struct {
	ID
}

func (s *IDMut) Clone() *IDMut {
	v := &IDMut{}
	*v = *s
	return v
}
func (s *IDMut) Freeze() *ID {
	return (*ID)(unsafe.Pointer(s))
}
func (s *IDMut) SetPartition(v uint16) *IDMut {
	s.partition = v
	return s
}
func (s *IDMut) SetRecord(v uint16) *IDMut {
	s.record = v
	return s
}
func (s *IDMut) SetResolution(v uint32) *IDMut {
	s.resolution = v
	return s
}
func (s *IDMut) SetSequence(v int64) *IDMut {
	s.sequence = v
	return s
}

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
func (s *String32) Unsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&s[0])),
		Len:  int(s[31]),
	}))
}
func (s *String32) String() string {
	return string(s[0:s[31]])
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
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(&s)))[0:]...), nil
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

type String65472 [65472]byte

func NewString65472(s string) *String65472 {
	v := String65472{}
	v.set(s)
	return &v
}
func (s *String65472) set(v string) {
	copy(s[0:65470], v)
	c := 65470
	l := len(v)
	if l > c {
		s[65470] = byte(c)
		s[65471] = byte(c >> 8)
	} else {
		s[65470] = byte(l)
		s[65471] = byte(l >> 8)
	}
}
func (s *String65472) Len() int {
	return int(*(*uint16)(unsafe.Pointer(&s[65470])))
}
func (s *String65472) Cap() int {
	return 65470
}
func (s *String65472) Unsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&s[0])),
		Len:  int(*(*uint16)(unsafe.Pointer(&s[65470]))),
	}))
}
func (s *String65472) String() string {
	return string(s[0:s.Len()])
}
func (s *String65472) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String65472) Clone() *String65472 {
	v := String65472{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String65472) Mut() *String65472Mut {
	return *(**String65472Mut)(unsafe.Pointer(&s))
}
func (s *String65472) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[65472]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 65472 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String65472) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[65472]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String65472) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[65472]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String65472) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[65472]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *String65472) UnmarshalBinary(b []byte) error {
	if len(b) < 65472 {
		return io.ErrShortBuffer
	}
	v := (*String65472)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String65472Mut struct {
	String65472
}

func (s *String65472Mut) Set(v string) {
	s.set(v)
}
func init() {
	{
		var b [2]byte
		v := uint16(1)
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		if *(*uint16)(unsafe.Pointer(&b[0])) != 1 {
			panic("BigEndian detected... compiled for LittleEndian only!!!")
		}
	}
	to := reflect.TypeOf
	type sf struct {
		n string
		o uintptr
		s uintptr
	}
	ss := func(tt interface{}, mtt interface{}, s uintptr, fl []sf) {
		t := to(tt)
		mt := to(mtt)
		if t.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", t.Name(), t.Size(), s))
		}
		if mt.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", mt.Name(), mt.Size(), s))
		}
		if t.NumField() != len(fl) {
			panic(fmt.Sprintf("%s field count = %d: expected %d", t.Name(), t.NumField(), len(fl)))
		}
		for i, ef := range fl {
			f := t.Field(i)
			if f.Offset != ef.o {
				panic(fmt.Sprintf("%s.%s offset = %d, expected = %d", t.Name(), f.Name, f.Offset, ef.o))
			}
			if f.Type.Size() != ef.s {
				panic(fmt.Sprintf("%s.%s size = %d, expected = %d", t.Name(), f.Name, f.Type.Size(), ef.s))
			}
			if f.Name != ef.n {
				panic(fmt.Sprintf("%s.%s expected field: %s", t.Name(), f.Name, ef.n))
			}
		}
	}

	ss(Stream{}, StreamMut{}, 56, []sf{
		{"id", 0, 16},
		{"first", 16, 8},
		{"records", 24, 8},
		{"pages", 32, 8},
		{"size", 40, 8},
		{"xsize", 48, 8},
	})
	ss(PageID{}, PageIDMut{}, 24, []sf{
		{"stream", 0, 16},
		{"id", 16, 8},
	})
	ss(Page{}, PageMut{}, 65536, []sf{
		{"first", 0, 8},
		{"duration", 8, 8},
		{"last", 16, 8},
		{"count", 24, 2},
		{"record", 26, 2},
		{"size", 28, 2},
		{"xsize", 30, 2},
		{"pad", 32, 32},
		{"data", 64, 65472},
	})
	ss(ID{}, IDMut{}, 16, []sf{
		{"partition", 0, 2},
		{"record", 2, 2},
		{"resolution", 4, 4},
		{"sequence", 8, 8},
	})

}
