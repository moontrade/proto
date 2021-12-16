package wap

import (
	"errors"
	"github.com/moontrade/nogc"
	"unsafe"
)

var (
	ErrOutOfMemory = errors.New("out of memory")
)

type Mutable struct {
	*Builder
	offset int32 // The underlying buffer may reallocate so add offset to root pointer each time.
	//length int32
}

func (mut *Mutable) IsRoot() bool           { return mut.offset == 0 }
func (mut Mutable) IsEmpty() bool           { return mut.ref == nil }
func (mut *Mutable) IsNil() bool            { return mut.ref == nil }
func (mut *Mutable) Unsafe() unsafe.Pointer { return unsafe.Add(mut.ptr.Unsafe(), mut.offset) }

func (b *Mutable) Free(vp *VPointer) {
	if vp == nil || *vp == 0 {
		return
	}
	length := vp.Len()
	b.trash += length
	*vp = 0
}

//func (b *Mutable) Mut(p unsafe.Pointer) Mutable {
//	offset := int32(b.OffsetOf(p))
//	return Mutable{b.Builder, offset}
//}

func (b *Mutable) Slice(p unsafe.Pointer) Mutable {
	return Mutable{b.Builder, int32(nogc.Pointer(p) - b.ptr)}
}

//func (b *Mutable) SliceRef(vp *VPointer, size int32) Mutable {
//	if vp == nil || *vp == 0 {
//		return Mutable{}
//	}
//	m := b.SliceAlloc(vp, size)
//	return m
//}

// SliceAlloc gets the sub Buffer pointed by VPointer allocating if VPointer is 0. It guarantees
// the sub Buffer is allocated.
func (b *Mutable) SliceAlloc(vp *VPointer, size int32) Mutable {
	if *vp == 0 {
		offset := b.VPointerOffset(vp)
		vp, _, _ = b.alloc(vp, size)
		if *vp == 0 {
			return Mutable{}
		}
		return Mutable{b.Builder, offset + int32(*vp)}
	} else {
		return Mutable{b.Builder, b.VPointerOffset(vp) + int32(*vp)}
	}
}

type _string struct {
	Data unsafe.Pointer
	Len  int
}

func (b *Mutable) WStr(existing *VPointer, value string) {
	val := *(*_string)(unsafe.Pointer(&value))
	b.writeString(existing, val.Data, int32(len(value)))
}

func (b *Mutable) WriteBytes(existing *VPointer, value []byte) {
	b.writeBytes(existing, unsafe.Pointer(&value[0]), int32(len(value)), int32(len(value)))
}

type Builder struct {
	ptr      nogc.Pointer
	len      int32
	cap      int32
	ref      unsafe.Pointer // GC needs a ref to keep alive
	trash    int32
	maxTrash int32 // Maximum amount of trash
	manual   bool
}

func (b *Builder) Len() int {
	return int(b.len)
}

func (b *Builder) Get() *Builder {
	return b
}

type BuilderProvider interface {
	Get() *Builder
}

type BuilderPointer struct {
	*Builder
}

func (b *BuilderPointer) Close() error {
	bb := b.Builder
	if bb == nil {
		return nil
	}
	if b.ptr != 0 {
		nogc.Free(b.ptr)
		b.ptr = 0
	}
	nogc.Free(nogc.Pointer(unsafe.Pointer(bb)))
	b.Builder = nil
	return nil
}

func AllocBuilder() BuilderPointer {
	p := nogc.AllocZeroed(unsafe.Sizeof(Builder{}))
	b := (*Builder)(p.Unsafe())
	b.manual = true
	return BuilderPointer{b}
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Finish() unsafe.Pointer {
	if b == nil || b.ptr == 0 {
		return nil
	}
	if b.ref == nil {
		p := b.ptr
		b.ptr = 0
		b.ref = nil
		return p.Unsafe()
	} else {
		buf := b.ref
		b.ptr = 0
		b.ref = nil
		return buf
	}
}

func (b *Builder) FinishPointer() unsafe.Pointer {
	if b == nil {
		return nil
	}
	if b.ref != nil {
		return nil
	}
	p := b.ptr
	b.ptr = 0
	b.ref = nil
	return p.Unsafe()
}

func (b *Builder) New(length, flex int32) Mutable {
	var (
		c uintptr
	)
	if b.manual {
		panic("manually allocated builder cannot GC allocate")
	}
	b.Finish()
	if flex == 0 {
		c = uintptr(length)
		b.ref = gcAllocZeroed(c)
	} else {
		c = uintptr(length + flex)
		b.ref = gcAlloc(c)
		nogc.Zero(b.ref, uintptr(length))
	}
	b.ptr = nogc.Pointer(b.ref)
	b.cap = int32(c)
	b.len = length
	return Mutable{b, 0}
}

func (b *Builder) Alloc(length, flex int32) Mutable {
	var (
		c uintptr
	)
	b.Finish()
	b.ref = nil
	b.ptr, c = nogc.AllocZeroedCap(uintptr(length + flex))
	b.cap = int32(c)
	b.len = length
	return Mutable{b, 0}
}

func (b *Builder) OffsetOf(p unsafe.Pointer) int64 {
	return int64(nogc.Pointer(p) - b.ptr)
}

func (b *Builder) FieldOffset(field *VPointer) int64 {
	return int64(nogc.Pointer(unsafe.Pointer(field)) - b.ptr)
}

func (b *Builder) newCap(needed int32) int32 {
	newCap := b.cap * 2
	if newCap < needed {
		newCap = needed
	}
	return newCap
}

func (b *Builder) VPointerOffset(existing *VPointer) int32 {
	return int32(nogc.Pointer(unsafe.Pointer(existing)) - b.ptr)
}

func (b *Builder) alloc(vp *VPointer, size int32) (*VPointer, nogc.Pointer, bool) {
	var (
		offset   = int32(nogc.Pointer(unsafe.Pointer(vp)) - b.ptr)
		value    = VPointer(b.len - offset)
		newLen   = size + b.len
		extended = false
	)
	if newLen > b.cap {
		c := uintptr(b.newCap(newLen))
		if b.ref == nil {
			b.ptr, c = nogc.ReallocCap(b.ptr, c)
			if b.ptr == 0 {
				panic(ErrOutOfMemory)
			}
		} else {
			b.ref = gcAlloc(c)
			if b.ref == nil {
				panic(ErrOutOfMemory)
			}
			nogc.Copy(b.ref, b.ptr.Unsafe(), uintptr(b.len))
			b.ptr = nogc.Pointer(b.ref)
		}
		b.cap = int32(c)
		extended = true
		vp = (*VPointer)((b.ptr + nogc.Pointer(offset)).Unsafe())
		*vp = value
	} else {
		*vp = value
	}
	r := b.ptr + nogc.Pointer(b.len)
	b.len = newLen
	return vp, r, extended
}

func (b *Builder) deleteSlice(existing *VPointer, size int32) {
	if existing == nil || *existing == 0 {
		return
	}
	b.trash += size
}

func (b *Builder) writeSlice(vp *VPointer, data unsafe.Pointer, size int32) bool {
	if *vp > 0 {
		nogc.Copy(vp.Deref(), data, uintptr(size))
		return false
	}

	vp, ptr, extended := b.alloc(vp, size)
	if *vp == 0 || b.ptr == 0 {
		panic(ErrOutOfMemory)
	}
	ptr.SetInt32LE(0, size)
	nogc.Copy((ptr + 4).Unsafe(), data, uintptr(size))
	return extended
}

func (b *Builder) deleteString(existing *VPointer) {
	if existing == nil || *existing == 0 {
		return
	}
	ptr := existing.deref()
	b.trash += ptr.Int32LE(0) + 4
}

func (b *Builder) writeString(vp *VPointer, data unsafe.Pointer, size int32) bool {
	if *vp > 0 {
		ptr := vp.deref()
		prevSize := ptr.Int32LE(0)
		if prevSize >= size+4 {
			b.trash += prevSize - size
			ptr.SetInt32LE(0, size)
			nogc.Copy((ptr + 4).Unsafe(), data, uintptr(size))
			return false
		}
		b.trash += prevSize + 4
	}

	vp, ptr, extended := b.alloc(vp, size+4)
	if *vp == 0 || b.ptr == 0 {
		panic(ErrOutOfMemory)
	}
	ptr.SetInt32LE(0, size)
	nogc.Copy((ptr + 4).Unsafe(), data, uintptr(size))
	return extended
}

func (b *Builder) deleteBytes(vp *VPointer) {
	if vp == nil || *vp == 0 {
		return
	}
	ptr := vp.deref()
	b.trash += ptr.Int32LE(0) + 8
}

func (b *Builder) writeBytes(existing *VPointer, data unsafe.Pointer, length, size int32) bool {
	if *existing > 0 {
		ptr := existing.deref()
		prevSize := ptr.Int32LE(0)
		if prevSize >= size+8 {
			b.trash += prevSize - size
			ptr.SetInt32LE(0, size)
			b.ptr.SetInt32LE(4, length)
			nogc.Copy((ptr + 8).Unsafe(), data, uintptr(length))
			return false
		}
		b.trash += prevSize + 8
	}

	existing, ptr, extended := b.alloc(existing, size+8)
	if *existing == 0 || b.ptr == 0 {
		panic(ErrOutOfMemory)
	}
	ptr.SetInt32LE(0, size)
	ptr.SetInt32LE(4, length)
	nogc.Copy((ptr + 8).Unsafe(), data, uintptr(length))
	return extended
}
