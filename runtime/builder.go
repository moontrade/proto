package runtime

import (
	"unsafe"
)

type Builder struct {
	PointerMut
	i int
}

func NewBuilder(size int) *Builder {
	return &Builder{
		PointerMut: GetPointerMut(size),
	}
}

func (b *Builder) append(value string) int {
	if b.len-b.i < len(value) {
		b.PointerMut = b.Grow(len(value))
	}
	i := b.i
	b.i += len(value)
	return i
}

func (b *Builder) Ensure(length int) *Builder {
	if b.len < length {
		b.PointerMut = b.Grow(length)
	}
	return b
}

func (b *Builder) Reserve(length int) int {
	if b.len-b.i < length {
		b.PointerMut = b.Grow(length)
	}
	i := b.i
	b.i += length
	return i
}

func (b *Builder) AppendStringFixed(max int, value string) *Builder {
	b.Ensure(b.i + max)
	i := b.i
	b.i += max
	b.SetStringFixed(i, max, value)
	return b
}

//func (b *Builder) AppendBytesFixed(max int, value []byte) *Builder {
//	b.Ensure(b.i + max)
//	if len(value) > max {
//		value = value[:max]
//	}
//
//	i := b.i
//	b.i += max
//	b.SetBytes(i, max, value)
//	return b
//}

func (b *Builder) AppendString(field int, value string) *Builder {
	offset := b.append(value)
	b.SetInt32(field, int32(offset))
	//b.SetVPointerUnsafe(field, offset, len(value))
	return b
}

func (b *Builder) AppendBytes(field int, value []byte) *Builder {
	offset := b.append(*(*string)(unsafe.Pointer(&value)))
	_ = offset
	//b.SetVPointerUnsafe(field, offset, len(value))
	return b
}

func (b *Builder) AppendPointer(field int, value Pointer) *Builder {
	offset := b.append(value.String())
	_ = offset
	//b.SetVPointerUnsafe(field, offset, value.len)
	return b
}

type Buffer struct {
	b []byte
	i int
}

func (w *Buffer) Take() []byte {
	if w.b == nil {
		return nil
	}
	result := w.b[0:w.i]
	w.b = nil
	w.i = 0
	return result
}

func (w *Buffer) Size() int {
	return w.i
}

func (w *Buffer) Ensure(size int) {
	if w.i+size > cap(w.b) {
		// Get pooled
		next := GetBytes(w.i + size)
		// Copy to new buf
		copy(next, w.b)
		if w.b != nil {
			// Release old back to pool
			PutBytes(w.b[:cap(w.b)])
		}
		// Set to new buf
		w.b = next
	}
}

func (w *Buffer) AppendByte(value byte) {
	w.Ensure(w.i + 1)
	w.b[w.i] = value
	w.i++
}

func (w *Buffer) AppendString(value string) {
	w.Ensure(w.i + len(value))
	copy(w.b[w.i:], value)
	w.i += len(value)
}

func (w *Buffer) AppendBytes(value []byte) {
	w.Ensure(w.i + len(value))
	copy(w.b[w.i:], value)
	w.i += len(value)
}

type ListBuilder struct {
	ptr         PointerMut
	elementSize int
}
