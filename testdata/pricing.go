package schema

import (
	. "github.com/moontrade/proto"
	"io"
	"unsafe"
)

type Price struct {
	open, high, low, close float64
}

func (p *Price) Open() float64  { return p.open }
func (p *Price) High() float64  { return p.high }
func (p *Price) Low() float64   { return p.low }
func (p *Price) Close() float64 { return p.close }

type PriceBlock struct {
	ColumnBlock
}

func (self *PriceBlock) Reader() PriceReader {
	return PriceReader{
		NewFixedReader(&self.BlockHeader),
	}
}

type PriceBlockMut struct {
	ColumnBlockMut
}

func (self *PriceBlockMut) append(blockSize BlockSize, price *Price) error {
	if price == nil {
		return io.ErrShortBuffer
	}
	return self.Append(blockSize, unsafe.Pointer(price), unsafe.Sizeof(Price{}))
}

type PriceBlock4 struct {
	PriceBlock
	Data [3900]byte
}

func (self *PriceBlock4) Mut() *PriceBlock4Mut {
	return (*PriceBlock4Mut)(unsafe.Pointer(self))
}

type PriceBlock4Mut struct {
	PriceBlockMut
	Data [3900]byte
}

func (self *PriceBlock4Mut) Freeze() *PriceBlock4 {
	return (*PriceBlock4)(unsafe.Pointer(self))
}

func (self *PriceBlock4Mut) Append(price *Price) error {
	return self.append(BlockSize4KB, price)
}

func (self *PriceBlockMut) Compress(b []byte) []byte {
	return b
}

func (self *PriceBlock4Mut) ToRows(into *PriceBlock4Mut) *PriceBlock4Mut {
	if self.Layout() == BlockLayoutRow {
		return self
	}
	if into == self || into == nil {
		into = &PriceBlock4Mut{}
	}
	into.PriceBlockMut = self.PriceBlockMut
	into.PriceBlockMut.SetLayout(BlockLayoutRow)
	for i := uintptr(0); i < uintptr(self.Count()); i++ {
		*(*Price)(unsafe.Pointer(&into.Data[i*unsafe.Sizeof(Price{})])) = Price{
			open:  *(*float64)(self.Item(121, 0, 8, i)),  // open
			high:  *(*float64)(self.Item(121, 8, 8, i)),  // high
			low:   *(*float64)(self.Item(121, 16, 8, i)), // low
			close: *(*float64)(self.Item(121, 24, 8, i)), // close
		}
	}

	return into
}

func (self *PriceBlock4) ToColumns(into *PriceBlock4) *PriceBlock4 {
	if self.Layout() == BlockLayoutColumn {
		return self
	}
	if into == self || into == nil {
		into = &PriceBlock4{}
	}
	return nil
}

type PriceReader struct {
	FixedBlockReader
}

func (r *PriceReader) Current() *Price {
	return (*Price)(r.FixedBlockReader.CurrentPtr())
}

func (r *PriceReader) First() *Price {
	return (*Price)(r.FixedBlockReader.FirstPtr())
}

func (r *PriceReader) Last() *Price {
	return (*Price)(r.FixedBlockReader.LastPtr())
}

func (r *PriceReader) Prev() *Price {
	return (*Price)(r.FixedBlockReader.PrevPtr())
}

func (r *PriceReader) Next() *Price {
	return (*Price)(r.FixedBlockReader.NextPtr())
}
