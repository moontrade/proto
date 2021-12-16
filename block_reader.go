package wap

import (
	"github.com/moontrade/nogc"
	"unsafe"
)

type BlockReader interface {
	Index() int

	Count() int

	Current() (unsafe.Pointer, int)

	First() (unsafe.Pointer, int)

	Last() (unsafe.Pointer, int)

	Prev() (unsafe.Pointer, int)

	Next() (unsafe.Pointer, int)
}

// FixedBlockReader reads fixed sized records
type FixedBlockReader struct {
	header *BlockHeader
	index  int
	record int
}

func NewFixedReader(header *BlockHeader) FixedBlockReader {
	return FixedBlockReader{
		header: header,
		index:  -1,
		record: int(header.record),
	}
}

func (fr *FixedBlockReader) Index() int {
	return fr.index
}

func (fr *FixedBlockReader) Count() int {
	return int(fr.header.count)
}

func (fr *FixedBlockReader) Current() (unsafe.Pointer, int) {
	if fr.index < 0 {
		return nil, 0
	}
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record), fr.record
}

func (fr *FixedBlockReader) CurrentPtr() unsafe.Pointer {
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record)
}

func (fr *FixedBlockReader) First() (unsafe.Pointer, int) {
	if fr.index == 0 {
		return nil, 0
	}
	fr.index = 0
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record), fr.record
}

func (fr *FixedBlockReader) FirstPtr() unsafe.Pointer {
	if fr.index == 0 {
		return nil
	}
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record)
}

func (fr *FixedBlockReader) Last() (unsafe.Pointer, int) {
	if fr.header.count == 0 {
		return nil, 0
	}
	return unsafe.Add(unsafe.Pointer(&fr.header.data), int(fr.header.count-1)*fr.record), fr.record
}

func (fr *FixedBlockReader) LastPtr() unsafe.Pointer {
	if fr.header.count == 0 {
		return nil
	}
	return unsafe.Add(unsafe.Pointer(&fr.header.data), int(fr.header.count-1)*fr.record)
}

func (fr *FixedBlockReader) Prev() (unsafe.Pointer, int) {
	if fr.index == 0 {
		return nil, 0
	}
	fr.index--
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record), fr.record
}

func (fr *FixedBlockReader) PrevPtr() unsafe.Pointer {
	if fr.index == 0 {
		return nil
	}
	fr.index--
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record)
}

func (fr *FixedBlockReader) Next() (unsafe.Pointer, int) {
	if fr.index+1 >= int(fr.header.count) {
		return nil, 0
	}
	fr.index++
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record), fr.record
}

func (fr *FixedBlockReader) NextPtr() unsafe.Pointer {
	if fr.index+1 >= int(fr.header.count) {
		return nil
	}
	fr.index++
	return unsafe.Add(unsafe.Pointer(&fr.header.data), fr.index*fr.record)
}

// FlexBlockReader reads variable length records. Each record has its uint16 size prepended and appended.
// uint16 | data | uint16 | uint16 | data | uint16
type FlexBlockReader struct {
	header *BlockHeader
	index  int
	offset int
	size   int
	ptr    nogc.Pointer
}

func (fr *FlexBlockReader) Index() int {
	return fr.index
}

func (fr *FlexBlockReader) Count() int {
	return int(fr.header.count)
}

func (fr *FlexBlockReader) Current() (unsafe.Pointer, int) {
	if fr.size == 0 {
		return nil, 0
	}
	return fr.ptr.Pointer(fr.offset).Unsafe(), fr.size
}

func (fr *FlexBlockReader) First() (unsafe.Pointer, int) {
	fr.index = 0
	if fr.header.count == 0 {
		return nil, 0
	}
	fr.size = int(fr.ptr.UInt16LE(0))
	if fr.size == 0 {
		return nil, 0
	}
	fr.offset = 2
	if fr.size+fr.offset+2 > int(fr.header.size) {
		fr.offset = 0
		fr.size = 0
		fr.index = 0
		return nil, 0
	}
	return fr.ptr.Pointer(fr.offset).Unsafe(), fr.size
}

func (fr *FlexBlockReader) Last() (unsafe.Pointer, int) {
	if fr.header.count == 0 {
		return nil, 0
	}
	fr.index = int(fr.header.count - 1)
	fr.offset = int(fr.header.size) - 2
	if fr.offset < 0 {
		return nil, 0
	}
	fr.size = int(fr.ptr.UInt16LE(fr.offset))
	if fr.size == 0 {
		return nil, 0
	}
	fr.offset -= fr.size
	return fr.ptr.Pointer(fr.offset).Unsafe(), fr.size
}

func (fr *FlexBlockReader) Prev() (unsafe.Pointer, int) {
	if fr.index == 0 {
		return nil, 0
	}
	fr.index--
	fr.offset -= 2
	fr.size = int(fr.ptr.UInt16LE(fr.offset))
	if fr.size == 0 {
		return nil, 0
	}
	fr.offset -= fr.size
	if fr.offset < 0 {
		return nil, 0
	}
	return fr.ptr.Pointer(fr.offset).Unsafe(), fr.size
}

func (fr *FlexBlockReader) Next() (unsafe.Pointer, int) {
	if fr.index+1 >= int(fr.header.count) {
		return nil, 0
	}
	fr.index++
	if fr.size > 0 {
		fr.offset += fr.size + 2
	}
	fr.size = int(fr.ptr.UInt16LE(fr.offset))
	if fr.size == 0 {
		return nil, 0
	}
	fr.offset += 2
	if fr.offset+fr.size > int(fr.header.size) {
		return nil, 0
	}
	return fr.ptr.Pointer(fr.offset).Unsafe(), fr.size
}
