package wap

import (
	"github.com/moontrade/nogc"
	"unsafe"
)

type VPointer int32

func (vp *VPointer) Deref() unsafe.Pointer {
	return unsafe.Add(unsafe.Pointer(vp), *vp)
}

func (vp *VPointer) deref() nogc.Pointer {
	return nogc.Pointer(uintptr(int64(uintptr(unsafe.Pointer(vp))) + int64(*vp)))
}

func Str(p *VPointer) string {
	if p == nil || *p == 0 {
		return ""
	}
	ptr := nogc.Pointer(p.Deref())
	l := int(ptr.Int32LE(0))
	if l == 0 {
		return ""
	}
	return ptr.String(4, l)
}

func Bytes(p *VPointer) []byte {
	if p == nil || *p == 0 {
		return nil
	}
	ptr := nogc.Pointer(p.Deref())
	l := int(ptr.Int32LE(0))
	if l == 0 {
		return nil
	}
	return ptr.Bytes(4, l, l)
}

func Slice(p *VPointer) unsafe.Pointer {
	if p == nil || *p == 0 {
		return nil
	}
	return nogc.Pointer(p.Deref()).Unsafe()
}

func Slab(p *VPointer) (unsafe.Pointer, int32) {
	if p == nil || *p == 0 {
		return nil, 0
	}
	ptr := nogc.Pointer(p.Deref())
	l := int(ptr.Int32LE(0))
	if l == 0 {
		return ptr.Unsafe(), 0
	}
	return ptr.Unsafe(), int32(l)
}

func (p *VPointer) Len() int32 {
	if p == nil || *p == 0 {
		return 0
	}
	return nogc.Pointer(p.Deref()).Int32LE(0)
}

func (p *VPointer) SetLen(value int32) {
	nogc.Pointer(p.Deref()).SetInt32LE(0, value)
}

func (p *VPointer) Str() string {
	if p == nil || *p == 0 {
		return ""
	}
	ptr := nogc.Pointer(p.Deref())
	l := int(ptr.Int32LE(0))
	if l == 0 {
		return ""
	}
	return ptr.String(4, l)
}

func (p *VPointer) Bytes() []byte {
	if p == nil || *p == 0 {
		return nil
	}
	ptr := nogc.Pointer(p.Deref())
	l := int(ptr.Int32LE(0))
	if l == 0 {
		return nil
	}
	return ptr.Bytes(4, l, l)
}

func (p *VPointer) Slice() unsafe.Pointer {
	if p == nil || *p == 0 {
		return nil
	}
	return nogc.Pointer(p.Deref()).Unsafe()
}

func (p *VPointer) Slab() (unsafe.Pointer, int32) {
	if p == nil || *p == 0 {
		return nil, 0
	}
	ptr := nogc.Pointer(p.Deref())
	l := int(ptr.Int32LE(0))
	if l == 0 {
		return ptr.Unsafe(), 0
	}
	return ptr.Unsafe(), int32(l)
}
