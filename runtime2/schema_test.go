package runtime2

import (
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"
)

func Test_Schema(t *testing.T) {
	s := &Schema{
		Records: []Record{
			{
				Name: "Bar",
				Fields: []Field{
					Int64Field("id", 0),
					Int64Field("start", 8),
					StringField("name", 16, 4),
					ListField("errors", 24, ListElement(StringElement())),
				},
			},
		},
	}
	_ = s
	fmt.Println(s)
	fmt.Println(toJson(s))
}

func toJson(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type IPrice interface {
	Open() float64
	High() float64
	Low() float64
	Close() float64
}

type BlockHeader struct {
	count uint16
	size  uint16
}

type IBlock interface {
	Header() *BlockHeader

	RecordAt(index int)
}

type IPriceBlock interface {
	Open() float64
	High() float64
	Low() float64
	Close() float64
}

type PriceBlock struct {
	BlockHeader
	data [16200]byte
}

//func (self *PriceBlock) Open(begin, end int) []float64 {
//
//}
//
//func (self *PriceBlock) Get(index int) IPrice {
//	return SOAPrice{}
//}

type PriceSOAReader struct {
	index       int
	count       int
	openOffset  int
	highOffset  int
	closeOffset int
	block       []byte
	records     []PriceSOA
}

func (self *PriceSOAReader) Open() []float64 {
	return *(*[]float64)(unsafe.Pointer(&struct {
		Data uintptr
		Len  int
		Cap  int
	}{
		Data: uintptr(unsafe.Pointer(&self.block[0])),
		Len:  self.count,
		Cap:  self.count,
	}))
}

func (self *PriceSOAReader) High() float64 {
	return *(*float64)(unsafe.Pointer(&self.block[self.highOffset+(self.index*8)]))
}

type PriceSOA struct {
	data        []byte
	index       int
	openOffset  int
	highOffset  int
	lowOffset   int
	closeOffset int
}

func (self *PriceSOA) Open() float64 {
	return *(*float64)(unsafe.Pointer(&self.data[self.index*8]))
}

func (self *PriceSOA) High() float64 {
	return *(*float64)(unsafe.Pointer(&self.data[self.highOffset+(self.index*8)]))
}

func (self *PriceSOA) Low() float64 {
	return *(*float64)(unsafe.Pointer(&self.data[self.lowOffset+(self.index*8)]))
}

func (self *PriceSOA) Close() float64 {
	return *(*float64)(unsafe.Pointer(&self.data[self.closeOffset+(self.index*8)]))
}

type PriceAOSReader struct {
	index   int
	count   int
	record  int
	current IPrice
	block   []byte
}

func (self *PriceAOSReader) Open() float64 {
	return self.current.Open()
}

func (self *PriceAOSReader) Next() bool {
	self.index++
	if self.index >= self.count {
		return false
	}

	return true
}
