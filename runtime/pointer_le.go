//go:build 386 || amd64 || arm || arm64 || ppc64le || mips64le || mipsle || riscv64 || wasm || tinygo.wasm
// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm tinygo.wasm

package runtime

import (
	"github.com/moontrade/proto/runtime/internal/pmath"
	"math"
	"reflect"
	"unsafe"
)

const (
	EmptyString = ""
)

// Pointer is a fat pointer
type Pointer struct {
	ptr unsafe.Pointer
	len uint16
	cap uint16
}

func (p Pointer) Unsafe() unsafe.Pointer {
	return p.ptr
}

func (p Pointer) Len() int {
	return int(p.len)
}

func NewPointer(size int) Pointer {
	return Wrap(make([]byte, size))
}

func NewPointerMut(size int) PointerMut {
	return WrapMut(make([]byte, size))
}

func Wrap(b []byte) Pointer {
	return Pointer{
		ptr: unsafe.Pointer(&b[0]),
		len: uint16(len(b)),
	}
}

func WrapMut(b []byte) PointerMut {
	return PointerMut{Pointer{
		ptr: unsafe.Pointer(&b[0]),
		len: uint16(len(b)),
	}}
}

func WrapString(s string) Pointer {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return Pointer{
		ptr: unsafe.Pointer(h.Data),
		len: uint16(len(s)),
	}
}

func WrapStringMut(s string) PointerMut {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return PointerMut{Pointer{
		ptr: unsafe.Pointer(h.Data),
		len: uint16(len(s)),
	}}
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

func (p Pointer) Bytes() []byte {
	if p.IsNil() {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr),
		Len:  int(p.len),
		Cap:  int(p.len),
	}))
}

func (p Pointer) Grow(extra int) int {
	if extra <= 0 {
		return -1
		//panic(fmt.Errorf("Pointer.Grow supplied negative extra: %d", extra))
	}
	newLen := pmath.CeilToPowerOfTwo(int(p.len) + extra)
	if newLen > math.MaxUint16 {
		return -1
		//panic(fmt.Sprintf("pointer exceeds max size of 64k: %d", newLen))
	}
	dst := GetBytes(pmath.CeilToPowerOfTwo(int(p.len) + extra))

	src := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr),
		Len:  int(p.len),
		Cap:  int(p.len),
	}))
	copy(dst, src)
	p.ptr = unsafe.Pointer(&dst[0])
	p.cap = uint16(len(dst))
	return int(p.len)
}

func (p Pointer) IsNil() bool {
	return uintptr(p.ptr) == 0
}

func (p Pointer) IsEmpty() bool {
	return uintptr(p.ptr) == 0 || p.len == 0
}

func (p Pointer) CheckBounds(offset int) bool {
	return uintptr(p.ptr) == 0 || int(p.len) < offset
}

func (p Pointer) Int8Unsafe(offset int) int8 {
	return *(*int8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int8(offset int) int8 {
	if p.CheckBounds(offset + 1) {
		return 0
	}
	return *(*int8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Byte(offset int) byte {
	if p.CheckBounds(offset + 1) {
		return 0
	}
	return *(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) ByteUnsafe(offset int) byte {
	return *(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt8(offset int) uint8 {
	if p.CheckBounds(offset + 1) {
		return 0
	}
	return *(*uint8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt8Unsafe(offset int) uint8 {
	return *(*uint8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int16(offset int) int16 {
	if p.CheckBounds(offset + 2) {
		return 0
	}
	return *(*int16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int16Unsafe(offset int) int16 {
	return *(*int16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt16(offset int) uint16 {
	if p.CheckBounds(offset + 2) {
		return 0
	}
	return *(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt16Unsafe(offset int) uint16 {
	return *(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int32(offset int) int32 {
	if p.CheckBounds(offset + 4) {
		return 0
	}
	return *(*int32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int32Unsafe(offset int) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt32(offset int) uint32 {
	if p.CheckBounds(offset + 4) {
		return 0
	}
	return *(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt32Unsafe(offset int) uint32 {
	return *(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int64(offset int) int64 {
	if p.CheckBounds(offset + 8) {
		return 0
	}
	return *(*int64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Int64Unsafe(offset int) int64 {
	return *(*int64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt64(offset int) uint64 {
	if p.CheckBounds(offset + 8) {
		return 0
	}
	return *(*uint64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) UInt64Unsafe(offset int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Float32(offset int) float32 {
	if p.CheckBounds(offset + 4) {
		return 0
	}
	return *(*float32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Float32Unsafe(offset int) float32 {
	return *(*float32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Float64(offset int) float64 {
	if p.CheckBounds(offset + 8) {
		return 0
	}
	return *(*float64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}
func (p Pointer) Float64Unsafe(offset int) float64 {
	return *(*float64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)))
}

func (p Pointer) Substr(offset, length int) string {
	return p.Slice(offset, length).String()
}

func (p Pointer) SliceBytes(offset, length int) []byte {
	return p.Slice(offset, length).Bytes()
}

func (p Pointer) Slice(offset, length int) Pointer {
	if p.IsNil() {
		return Pointer{}
	}
	if offset+length > int(p.len) {
		return Pointer{}
	}
	return Pointer{
		//ptr: unsafe.Pointer(uintptr(p.ptr) + uintptr(offset)),
		ptr: unsafe.Add(p.ptr, offset),
		len: uint16(length),
	}
}

func (p Pointer) Mut() PointerMut {
	return PointerMut{p}
}

type PointerMut struct {
	Pointer
}

func (p PointerMut) SetInt8(offset int, value int8) {
	if p.CheckBounds(offset + 1) {
		return
	}
	*(*int8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt8Unsafe(offset int, value int8) {
	*(*int8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt8(offset int, value uint8) {
	if p.CheckBounds(offset + 1) {
		return
	}
	*(*uint8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt8Unsafe(offset int, value uint8) {
	*(*uint8)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetByte(offset int, value byte) {
	if p.CheckBounds(offset + 1) {
		return
	}
	*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetByteUnsafe(offset int, value byte) {
	*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt16(offset int, value int16) {
	if p.CheckBounds(offset + 2) {
		return
	}
	*(*int16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt16Unsafe(offset int, value int16) {
	*(*int16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt16(offset int, value uint16) {
	if p.CheckBounds(offset + 2) {
		return
	}
	*(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt16Unsafe(offset int, value uint16) {
	*(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt32(offset int, value int32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	*(*int32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt32Unsafe(offset int, value int32) {
	*(*int32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt32(offset int, value uint32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	*(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt32Unsafe(offset int, value uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt64(offset int, value int64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	*(*int64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetInt64Unsafe(offset int, value int64) {
	*(*int64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt64(offset int, value uint64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	*(*uint64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetUInt64Unsafe(offset int, value uint64) {
	*(*uint64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetFloat32(offset int, value float32) {
	if p.CheckBounds(offset + 4) {
		return
	}
	*(*float32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetFloat32Unsafe(offset int, value float32) {
	*(*float32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetFloat64(offset int, value float64) {
	if p.CheckBounds(offset + 8) {
		return
	}
	*(*float64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetFloat64Unsafe(offset int, value float64) {
	*(*float64)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset))) = value
}
func (p PointerMut) SetString(offset int, value string) {
	if p.CheckBounds(offset + len(value)) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p PointerMut) SetStringUnsafe(offset int, value string) {
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p PointerMut) SetBytes(offset int, value []byte) {
	if p.CheckBounds(offset + len(value)) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p PointerMut) SetBytesUnsafe(offset int, value []byte) {
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)
}
func (p PointerMut) SetPointer(offset int, value Pointer) {
	if p.CheckBounds(offset + value.Len()) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  value.Len(),
		Cap:  value.Len(),
	}))
	copy(dst, value.String())
}
func (p PointerMut) SetPointerUnsafe(offset int, value Pointer) {
	if p.CheckBounds(offset + value.Len()) {
		return
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  value.Len(),
		Cap:  value.Len(),
	}))
	copy(dst, value.String())
}

func (p PointerMut) BytesFixed(offset, size int) []byte {
	if p.CheckBounds(offset + size) {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  size,
		Cap:  size,
	}))
}

func (p PointerMut) BytesFixedUnsafe(offset, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  size,
		Cap:  size,
	}))
}

func (p PointerMut) StringFixed(offset, max int) string {
	if p.CheckBounds(offset + max) {
		return EmptyString
	}
	sizeBytes := FixedStringLengthBytes(max)
	// Read length
	var length int
	switch sizeBytes {
	case 1:
		length = int(*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 2:
		length = int(*(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
		//length = int(*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end))) |
		//	int(*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end+1)))
	case 4:
		length = int(*(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	default:
		return EmptyString
	}

	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  length,
	}))
}

func (p PointerMut) StringFixedUnsafe(offset, max int) string {
	sizeBytes := FixedStringLengthBytes(max)
	// Read length
	var length int
	switch sizeBytes {
	case 1:
		length = int(*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 2:
		length = int(*(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 4:
		length = int(*(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	default:
		return EmptyString
	}

	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  length,
	}))
}

func (p PointerMut) StringFixedBytes(offset, max int) []byte {
	if p.CheckBounds(offset + max) {
		return nil
	}
	sizeBytes := FixedStringLengthBytes(max)
	// Read length
	var length int
	switch sizeBytes {
	case 1:
		length = int(*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 2:
		length = int(*(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 4:
		length = int(*(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	default:
		return nil
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  length,
		Cap:  length,
	}))
}

func (p PointerMut) StringFixedBytesUnsafe(offset, max int) []byte {
	sizeBytes := FixedStringLengthBytes(max)
	// Read length
	var length int
	switch sizeBytes {
	case 1:
		length = int(*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 2:
		length = int(*(*uint16)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	case 4:
		length = int(*(*uint32)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))))
	default:
		return nil
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  length,
		Cap:  length,
	}))
}

func (p PointerMut) SetStringFixed(offset, max int, value string) {
	if p.CheckBounds(offset + max) {
		return
	}
	sizeBytes := FixedStringLengthBytes(max)
	length := len(value)
	if length > max-sizeBytes {
		length = max - sizeBytes
		value = value[:length]
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)

	// Write length
	switch sizeBytes {
	case 1:
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))) = byte(length)
	case 2:
		end := uintptr(offset + max - sizeBytes)

		*(*byte)(unsafe.Add(p.ptr, end)) = byte(length)
		*(*byte)(unsafe.Add(p.ptr, end+1)) = byte(length)
		//*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end)) = byte(length)
		//*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 1)) = byte(length >> 8)
	case 3:
		end := uintptr(offset + max - sizeBytes)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end)) = byte(length)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 1)) = byte(length >> 8)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 2)) = byte(length >> 16)
	case 4:
		end := uintptr(offset + max - sizeBytes)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end)) = byte(length)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 1)) = byte(length >> 8)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 2)) = byte(length >> 16)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 3)) = byte(length >> 24)
	}
}

func (p PointerMut) SetStringFixedUnsafe(offset, max int, value string) {
	sizeBytes := FixedStringLengthBytes(max)
	length := len(value)
	if length > max-sizeBytes {
		length = max - sizeBytes
		value = value[:length]
	}
	dst := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p.ptr) + uintptr(offset),
		Len:  len(value),
		Cap:  len(value),
	}))
	copy(dst, value)

	// Write length
	switch sizeBytes {
	case 1:
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + uintptr(offset+max-sizeBytes))) = byte(length)
	case 2:
		end := uintptr(offset + max - sizeBytes)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end)) = byte(length)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 1)) = byte(length >> 8)
	case 3:
		end := uintptr(offset + max - sizeBytes)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end)) = byte(length)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 1)) = byte(length >> 8)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 2)) = byte(length >> 16)
	case 4:
		end := uintptr(offset + max - sizeBytes)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end)) = byte(length)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 1)) = byte(length >> 8)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 2)) = byte(length >> 16)
		*(*byte)(unsafe.Pointer(uintptr(p.ptr) + end + 3)) = byte(length >> 24)
	}
}

func (p PointerMut) SetSize(offset int, size int32) int {
	if size < 0 {
		return -1
	}
	if size < 255 {
		p.SetByte(offset, byte(size))
		return 1
	}
	if size < 65535 {
		p.SetUInt16(offset, uint16(size))
		return 2
	}
	p.SetInt32(offset, size)
	return 4
}

func (p PointerMut) SetVarString(offset int, value string) {

}

func (p PointerMut) append(value string) int {
	if int(p.cap-p.len) < len(value) {
		p.Grow(len(value))
	}
	if int(p.len)+len(value) > math.MaxUint16 {
		return -1
	}
	i := p.len
	p.len += uint16(len(value))
	return int(i)
}

func (p PointerMut) Ensure(length int) int {
	if int(p.cap) < length {
		return p.Grow(length)
	}
	return int(p.len)
}

func (p PointerMut) Reserve(length int) int {
	if int(p.cap-p.len) < length {
		p.Grow(length)
	}
	i := p.len
	p.len += uint16(length)
	return int(i)
}
