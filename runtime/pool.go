package runtime

import (
	"github.com/moontrade/proto/runtime/internal/pmath"
	"sync"
)

var (
	pool1 = &sync.Pool{New: func() interface{} {
		return make([]byte, 1)
	}}
	pool2 = &sync.Pool{New: func() interface{} {
		return make([]byte, 2)
	}}
	pool4 = &sync.Pool{New: func() interface{} {
		return make([]byte, 4)
	}}
	pool8 = &sync.Pool{New: func() interface{} {
		return make([]byte, 8)
	}}
	pool12 = &sync.Pool{New: func() interface{} {
		return make([]byte, 12)
	}}
	pool16 = &sync.Pool{New: func() interface{} {
		return make([]byte, 16)
	}}
	pool24 = &sync.Pool{New: func() interface{} {
		return make([]byte, 24)
	}}
	pool32 = &sync.Pool{New: func() interface{} {
		return make([]byte, 32)
	}}
	pool40 = &sync.Pool{New: func() interface{} {
		return make([]byte, 40)
	}}
	pool48 = &sync.Pool{New: func() interface{} {
		return make([]byte, 48)
	}}
	pool56 = &sync.Pool{New: func() interface{} {
		return make([]byte, 56)
	}}
	pool64 = &sync.Pool{New: func() interface{} {
		return make([]byte, 64)
	}}
	pool72 = &sync.Pool{New: func() interface{} {
		return make([]byte, 72)
	}}
	pool96 = &sync.Pool{New: func() interface{} {
		return make([]byte, 96)
	}}
	pool128 = &sync.Pool{New: func() interface{} {
		return make([]byte, 128)
	}}
	pool192 = &sync.Pool{New: func() interface{} {
		return make([]byte, 192)
	}}
	pool256 = &sync.Pool{New: func() interface{} {
		return make([]byte, 256)
	}}
	pool384 = &sync.Pool{New: func() interface{} {
		return make([]byte, 384)
	}}
	pool512 = &sync.Pool{New: func() interface{} {
		return make([]byte, 512)
	}}
	pool768 = &sync.Pool{New: func() interface{} {
		return make([]byte, 768)
	}}
	pool1024 = &sync.Pool{New: func() interface{} {
		return make([]byte, 1024)
	}}
	pool2048 = &sync.Pool{New: func() interface{} {
		return make([]byte, 2048)
	}}
	pool4096 = &sync.Pool{New: func() interface{} {
		return make([]byte, 4096)
	}}
	pool8192 = &sync.Pool{New: func() interface{} {
		return make([]byte, 8192)
	}}
	pool16384 = &sync.Pool{New: func() interface{} {
		return make([]byte, 16384)
	}}
	pool32768 = &sync.Pool{New: func() interface{} {
		return make([]byte, 32768)
	}}
	pool65536 = &sync.Pool{New: func() interface{} {
		return make([]byte, 65536)
	}}
)

func GetPointer(n int) Pointer {
	return Wrap(GetBytes(n))
}

func GetPointerMut(n int) PointerMut {
	return WrapMut(GetBytes(n))
}

func PutPointer(p Pointer) {
	if p.IsNil() {
		return
	}
	PutBytes(p.Bytes())
}

func PutPointerMut(p PointerMut) {
	if p.IsNil() {
		return
	}
	PutBytes(p.Bytes())
}

func GetBytes(n int) []byte {
	v := pmath.CeilToPowerOfTwo(n)
	switch v {
	case 0, 1:
		return pool1.Get().([]byte)[:n]
	case 2:
		return pool2.Get().([]byte)[:n]
	case 4:
		return pool4.Get().([]byte)[:n]
	case 8:
		return pool8.Get().([]byte)[:n]
	case 16:
		return pool16.Get().([]byte)[:n]
	case 24:
		return pool24.Get().([]byte)[:n]
	case 32:
		return pool32.Get().([]byte)[:n]
	case 64:
		switch {
		case n < 41:
			return pool40.Get().([]byte)[:n]
		case n < 49:
			return pool48.Get().([]byte)[:n]
		case n < 57:
			return pool56.Get().([]byte)[:n]
		}
		return pool64.Get().([]byte)[:n]
	case 128:
		switch {
		case n < 73:
			return pool72.Get().([]byte)[:n]
		case n < 97:
			return pool96.Get().([]byte)[:n]
		}
		return pool128.Get().([]byte)[:n]
	case 256:
		switch {
		case n < 193:
			return pool192.Get().([]byte)[:n]
		}
		return pool256.Get().([]byte)[:n]
	case 512:
		if n <= 384 {
			return pool384.Get().([]byte)
		}
		return pool512.Get().([]byte)[:n]
	case 1024:
		if n <= 768 {
			return pool768.Get().([]byte)[:n]
		}
		return pool1024.Get().([]byte)[:n]
	case 2048:
		return pool2048.Get().([]byte)[:n]
	case 4096:
		return pool4096.Get().([]byte)[:n]
	case 8192:
		return pool8192.Get().([]byte)[:n]
	case 16384:
		return pool16384.Get().([]byte)[:n]
	case 32768:
		return pool32768.Get().([]byte)[:n]
	case 65536:
		return pool65536.Get().([]byte)[:n]
	}

	return make([]byte, n)
}

func PutBytes(b []byte) {
	switch cap(b) {
	case 1:
		pool1.Put(b)
	case 2:
		pool2.Put(b)
	case 4:
		pool4.Put(b)
	case 8:
		pool8.Put(b)
	case 12:
		pool12.Put(b)
	case 16:
		pool16.Put(b)
	case 24:
		pool24.Put(b)
	case 32:
		pool32.Put(b)
	case 40:
		pool40.Put(b)
	case 48:
		pool48.Put(b)
	case 56:
		pool56.Put(b)
	case 64:
		pool64.Put(b)
	case 72:
		pool72.Put(b)
	case 96:
		pool96.Put(b)
	case 128:
		pool128.Put(b)
	case 192:
		pool192.Put(b)
	case 256:
		pool256.Put(b)
	case 384:
		pool384.Put(b)
	case 512:
		pool512.Put(b)
	case 768:
		pool768.Put(b)
	case 1024:
		pool1024.Put(b)
	case 2048:
		pool2048.Put(b)
	case 4096:
		pool4096.Put(b)
	case 8192:
		pool8192.Put(b)
	case 16384:
		pool16384.Put(b)
	case 32768:
		pool32768.Put(b)
	case 65536:
		pool65536.Put(b)
	}
}
