package wap

import (
	"github.com/moontrade/nogc"
	"io"
	"unsafe"
)

type _slice struct {
	Data uintptr
	Len  int
	Cap  int
}

type RecordHeader struct {
	id        int64
	timestamp int64
	start     int64
	end       int64
}

type BlockHeader struct {
	streamID  uint64
	id        uint64
	headID    uint64
	headStart int64
	blocks    uint64
	records   uint64
	storage   uint64
	storageU  uint64
	created   int64
	completed int64
	start     int64
	end       int64
	min       uint64
	max       uint64
	count     uint16
	size      uint16
	sizeU     uint16
	sizeX     uint16
	record    uint16
	block     BlockSize
	encoding  Encoding
	layout    BlockLayout
	kind      StreamKind
	format    Format
	data      struct{}
}

type BlockHeaderMut struct {
	BlockHeader
}

type FixedBlockHeaderMut struct {
	BlockHeaderMut
}

func (self *BlockHeader) Offset(offset uintptr) unsafe.Pointer {
	return nogc.Pointer(uintptr(unsafe.Pointer(&self.data)) + offset).Unsafe()
}

func (self *FixedBlockHeaderMut) Append(blockSize BlockSize, data unsafe.Pointer, size uintptr) error {
	self.block = blockSize
	if uintptr(self.size)+size > uintptr(self.block) {
		return io.ErrShortWrite
	}
	dst := unsafe.Add(unsafe.Pointer(&self.data), self.size)
	self.count += 1
	self.size += uint16(size)
	self.sizeU = self.size
	self.record = uint16(size)
	nogc.Copy(dst, data, size)
	return nil
}

func (self *BlockHeaderMut) Append(blockSize BlockSize, data unsafe.Pointer, size uintptr) error {
	self.block = blockSize
	if uintptr(self.size)+size+4 > uintptr(self.block) {
		return io.ErrShortWrite
	}
	dst := unsafe.Add(unsafe.Pointer(&self.data), self.size)
	self.count += 1
	self.size += uint16(size) + 4
	self.sizeU = self.size
	*(*uint16)(dst) = uint16(size)
	nogc.Copy(unsafe.Add(dst, 2), data, size)
	*(*uint16)(unsafe.Add(dst, 2+size)) = uint16(size)
	return nil
}

func (self *BlockHeader) Count() uint16 {
	return self.count
}

func (self *BlockHeaderMut) AddCount(value uint16) {
	self.count += value
}

func (self *BlockHeaderMut) SetCount(count uint16) {
	self.count = count
}

func (self *BlockHeader) Size() uint16 {
	return self.size
}

func (self *BlockHeaderMut) SetSize(value uint16) {
	self.size = value
}

func (self *BlockHeaderMut) SetRecord(value uint16) {
	self.record = value
}

func (self *BlockHeader) SizeU() uint16 {
	return self.sizeU
}

func (self *BlockHeaderMut) SetSizeU(value uint16) {
	self.sizeU = value
}

func (self *BlockHeader) BlockSize() BlockSize {
	return self.block
}

func (self *BlockHeaderMut) SetBlockSize(value BlockSize) {
	self.block = value
}

func (self *BlockHeader) Layout() BlockLayout {
	return self.layout
}

func (self *BlockHeaderMut) SetLayout(value BlockLayout) {
	self.layout = value
}

type FixedBlock struct {
	BlockHeader
}

func (h *FixedBlock) Record(index uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(&h.data)) + (index * uintptr(h.record)))
}

func (h *FixedBlock) Reader() *FixedBlockReader {
	return &FixedBlockReader{
		header: &h.BlockHeader,
		index:  0,
		record: int(h.record),
	}
}

type FixedBlockMut struct {
	BlockHeader
}

func (h *FixedBlockMut) Record(index uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(&h.data)) + (index * uintptr(h.record)))
}

type ColumnBlock struct {
	FixedBlock
}

func (h *ColumnBlock) Slice(max, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(&_slice{
		Data: uintptr(unsafe.Pointer(&h.data)) + (offset * max),
		Len:  int(h.count),
		Cap:  int(max),
	})
}

func (h *ColumnBlock) Item(max, offset, size, index uintptr) unsafe.Pointer {
	return unsafe.Pointer(
		uintptr(unsafe.Pointer(&h.data)) +
			(offset * max) +
			(size * index))
}

//func (h *ColumnBlock) I8(max, offset, size, index uintptr) int8 {
//	return *(*int8)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) I8Column(max, offset uintptr) []int8 {
//	return *(*[]int8)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) U8(max, offset, size, index uintptr) uint8 {
//	return *(*uint8)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) U8Column(max, offset uintptr) []uint8 {
//	return *(*[]uint8)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) I16(max, offset, size, index uintptr) int16 {
//	return *(*int16)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) I16Column(max, offset uintptr) []int16 {
//	return *(*[]int16)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) U16(max, offset, size, index uintptr) uint16 {
//	return *(*uint16)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) U16Column(max, offset uintptr) []uint16 {
//	return *(*[]uint16)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) I32(max, offset, size, index uintptr) int32 {
//	return *(*int32)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) I32Column(max, offset uintptr) []int32 {
//	return *(*[]int32)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) U32(max, offset, size, index uintptr) uint32 {
//	return *(*uint32)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) U32Column(max, offset uintptr) []uint32 {
//	return *(*[]uint32)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) I64(max, offset, size, index uintptr) int64 {
//	return *(*int64)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) I64Column(max, offset uintptr) []int64 {
//	return *(*[]int64)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) U64(max, offset, size, index uintptr) uint64 {
//	return *(*uint64)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) U64Column(max, offset uintptr) []uint64 {
//	return *(*[]uint64)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) F32(max, offset, size, index uintptr) float32 {
//	return *(*float32)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) F32Column(max, offset uintptr) []float32 {
//	return *(*[]float32)(h.Slice(max, offset))
//}
//
//func (h *ColumnBlock) F64(max, offset, size, index uintptr) float64 {
//	return *(*float64)(h.Item(max, offset, size, index))
//}
//
//func (h *ColumnBlock) F64Column(max, offset uintptr) []float64 {
//	return *(*[]float64)(h.Slice(max, offset))
//}

type ColumnBlockMut struct {
	FixedBlockHeaderMut
}

func (h *ColumnBlockMut) Slice(max, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(&_slice{
		Data: uintptr(unsafe.Pointer(&h.data)) + (offset * max),
		Len:  int(h.count),
		Cap:  int(max),
	})
}

func (h *ColumnBlockMut) Item(max, offset, size, index uintptr) unsafe.Pointer {
	return unsafe.Pointer(
		uintptr(unsafe.Pointer(&h.data)) +
			(offset * max) +
			(size * index))
}

//func (h *ColumnBlockMut) SetI8(max, offset, size, index uintptr, value int8) {
//	*(*int8)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetU8(max, offset, size, index uintptr, value uint8) {
//	*(*uint8)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetI16(max, offset, size, index uintptr, value int16) {
//	*(*int16)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetU16(max, offset, size, index uintptr, value uint16) {
//	*(*uint16)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetI32(max, offset, size, index uintptr, value int32) {
//	*(*int32)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetU32(max, offset, size, index uintptr, value uint32) {
//	*(*uint32)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetI64(max, offset, size, index uintptr, value int64) {
//	*(*int64)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetU64(max, offset, size, index uintptr, value uint64) {
//	*(*uint64)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetF32(max, offset, size, index uintptr, value float32) {
//	*(*float32)(h.Item(max, offset, size, index)) = value
//}
//
//func (h *ColumnBlockMut) SetF64(max, offset, size, index uintptr, value float64) {
//	*(*float64)(h.Item(max, offset, size, index)) = value
//}
