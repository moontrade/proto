package wap

import (
	"unsafe"
)

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

func gcAlloc(size uintptr) unsafe.Pointer {
	return mallocgc(size, nil, false)
	//b := make([]byte, size)
	//return unsafe.Pointer(&b[0])
}

func gcAllocZeroed(size uintptr) unsafe.Pointer {
	return mallocgc(size, nil, true)
}
