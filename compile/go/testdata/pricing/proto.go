// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package pricing

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type Aggressor int32

const (
	Aggressor_Buy = Aggressor(0)

	Aggressor_Sell = Aggressor(1)
)

// Candlestick
type Candle struct {
	open  float64
	high  float64
	low   float64
	close float64
}

func (s *Candle) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Candle) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["open"] = s.Open()
	m["high"] = s.High()
	m["low"] = s.Low()
	m["close"] = s.Close()
	return m
}

func (s *Candle) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 32 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Candle) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[32]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Candle) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Candle) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Candle) Read(b []byte) (n int, err error) {
	if len(b) < 32 {
		return -1, io.ErrShortBuffer
	}
	v := (*Candle)(unsafe.Pointer(&b[0]))
	*v = *s
	return 32, nil
}
func (s *Candle) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*Candle)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Candle) Clone() *Candle {
	v := &Candle{}
	*v = *s
	return v
}
func (s *Candle) Bytes() []byte {
	return (*(*[32]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Candle) Mut() *CandleMut {
	return (*CandleMut)(unsafe.Pointer(s))
}
func (s *Candle) Open() float64 {
	return s.open
}
func (s *Candle) High() float64 {
	return s.high
}
func (s *Candle) Low() float64 {
	return s.low
}
func (s *Candle) Close() float64 {
	return s.close
}

// Candlestick
type CandleMut struct {
	Candle
}

func (s *CandleMut) Clone() *CandleMut {
	v := &CandleMut{}
	*v = *s
	return v
}
func (s *CandleMut) Freeze() *Candle {
	return (*Candle)(unsafe.Pointer(s))
}
func (s *CandleMut) SetOpen(v float64) *CandleMut {
	s.open = v
	return s
}
func (s *CandleMut) SetHigh(v float64) *CandleMut {
	s.high = v
	return s
}
func (s *CandleMut) SetLow(v float64) *CandleMut {
	s.low = v
	return s
}
func (s *CandleMut) SetClose(v float64) *CandleMut {
	s.close = v
	return s
}

type Ticks struct {
	total int64
	up    int64
	down  int64
}

func (s *Ticks) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Ticks) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["total"] = s.Total()
	m["up"] = s.Up()
	m["down"] = s.Down()
	return m
}

func (s *Ticks) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Ticks) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Ticks) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Ticks) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Ticks) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*Ticks)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *Ticks) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*Ticks)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Ticks) Clone() *Ticks {
	v := &Ticks{}
	*v = *s
	return v
}
func (s *Ticks) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Ticks) Mut() *TicksMut {
	return (*TicksMut)(unsafe.Pointer(s))
}
func (s *Ticks) Total() int64 {
	return s.total
}
func (s *Ticks) Up() int64 {
	return s.up
}
func (s *Ticks) Down() int64 {
	return s.down
}

type TicksMut struct {
	Ticks
}

func (s *TicksMut) Clone() *TicksMut {
	v := &TicksMut{}
	*v = *s
	return v
}
func (s *TicksMut) Freeze() *Ticks {
	return (*Ticks)(unsafe.Pointer(s))
}
func (s *TicksMut) SetTotal(v int64) *TicksMut {
	s.total = v
	return s
}
func (s *TicksMut) SetUp(v int64) *TicksMut {
	s.up = v
	return s
}
func (s *TicksMut) SetDown(v int64) *TicksMut {
	s.down = v
	return s
}

type Volume struct {
	total float64
	buy   float64
	sell  float64
}

func (s *Volume) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Volume) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["total"] = s.Total()
	m["buy"] = s.Buy()
	m["sell"] = s.Sell()
	return m
}

func (s *Volume) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Volume) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Volume) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Volume) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Volume) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*Volume)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *Volume) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*Volume)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Volume) Clone() *Volume {
	v := &Volume{}
	*v = *s
	return v
}
func (s *Volume) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Volume) Mut() *VolumeMut {
	return (*VolumeMut)(unsafe.Pointer(s))
}
func (s *Volume) Total() float64 {
	return s.total
}
func (s *Volume) Buy() float64 {
	return s.buy
}
func (s *Volume) Sell() float64 {
	return s.sell
}

type VolumeMut struct {
	Volume
}

func (s *VolumeMut) Clone() *VolumeMut {
	v := &VolumeMut{}
	*v = *s
	return v
}
func (s *VolumeMut) Freeze() *Volume {
	return (*Volume)(unsafe.Pointer(s))
}
func (s *VolumeMut) SetTotal(v float64) *VolumeMut {
	s.total = v
	return s
}
func (s *VolumeMut) SetBuy(v float64) *VolumeMut {
	s.buy = v
	return s
}
func (s *VolumeMut) SetSell(v float64) *VolumeMut {
	s.sell = v
	return s
}

type Trades struct {
	count int64
	min   int64
	max   int64
}

func (s *Trades) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Trades) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["count"] = s.Count()
	m["min"] = s.Min()
	m["max"] = s.Max()
	return m
}

func (s *Trades) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Trades) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Trades) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Trades) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Trades) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*Trades)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *Trades) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*Trades)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Trades) Clone() *Trades {
	v := &Trades{}
	*v = *s
	return v
}
func (s *Trades) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Trades) Mut() *TradesMut {
	return (*TradesMut)(unsafe.Pointer(s))
}
func (s *Trades) Count() int64 {
	return s.count
}
func (s *Trades) Min() int64 {
	return s.min
}
func (s *Trades) Max() int64 {
	return s.max
}

type TradesMut struct {
	Trades
}

func (s *TradesMut) Clone() *TradesMut {
	v := &TradesMut{}
	*v = *s
	return v
}
func (s *TradesMut) Freeze() *Trades {
	return (*Trades)(unsafe.Pointer(s))
}
func (s *TradesMut) SetCount(v int64) *TradesMut {
	s.count = v
	return s
}
func (s *TradesMut) SetMin(v int64) *TradesMut {
	s.min = v
	return s
}
func (s *TradesMut) SetMax(v int64) *TradesMut {
	s.max = v
	return s
}

type Bar struct {
	start    int64
	end      int64
	price    Candle
	bid      Candle
	ask      Candle
	spread   Spread
	ticks    Ticks
	volume   Volume
	interest Volume
}

func (s *Bar) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Bar) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["start"] = s.Start()
	m["end"] = s.End()
	m["price"] = s.Price().MarshalMap(nil)
	m["bid"] = s.Bid().MarshalMap(nil)
	m["ask"] = s.Ask().MarshalMap(nil)
	m["spread"] = s.Spread().MarshalMap(nil)
	m["ticks"] = s.Ticks().MarshalMap(nil)
	m["volume"] = s.Volume().MarshalMap(nil)
	m["interest"] = s.Interest().MarshalMap(nil)
	return m
}

func (s *Bar) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[208]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 208 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Bar) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[208]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Bar) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[208]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Bar) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[208]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Bar) Read(b []byte) (n int, err error) {
	if len(b) < 208 {
		return -1, io.ErrShortBuffer
	}
	v := (*Bar)(unsafe.Pointer(&b[0]))
	*v = *s
	return 208, nil
}
func (s *Bar) UnmarshalBinary(b []byte) error {
	if len(b) < 208 {
		return io.ErrShortBuffer
	}
	v := (*Bar)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Bar) Clone() *Bar {
	v := &Bar{}
	*v = *s
	return v
}
func (s *Bar) Bytes() []byte {
	return (*(*[208]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Bar) Mut() *BarMut {
	return (*BarMut)(unsafe.Pointer(s))
}
func (s *Bar) Start() int64 {
	return s.start
}
func (s *Bar) End() int64 {
	return s.end
}
func (s *Bar) Price() *Candle {
	return &s.price
}
func (s *Bar) Bid() *Candle {
	return &s.bid
}
func (s *Bar) Ask() *Candle {
	return &s.ask
}
func (s *Bar) Spread() *Spread {
	return &s.spread
}
func (s *Bar) Ticks() *Ticks {
	return &s.ticks
}
func (s *Bar) Volume() *Volume {
	return &s.volume
}
func (s *Bar) Interest() *Volume {
	return &s.interest
}

type BarMut struct {
	Bar
}

func (s *BarMut) Clone() *BarMut {
	v := &BarMut{}
	*v = *s
	return v
}
func (s *BarMut) Freeze() *Bar {
	return (*Bar)(unsafe.Pointer(s))
}
func (s *BarMut) SetStart(v int64) *BarMut {
	s.start = v
	return s
}
func (s *BarMut) SetEnd(v int64) *BarMut {
	s.end = v
	return s
}
func (s *BarMut) Price() *CandleMut {
	return s.price.Mut()
}
func (s *BarMut) SetPrice(v *Candle) *BarMut {
	s.price = *v
	return s
}
func (s *BarMut) Bid() *CandleMut {
	return s.bid.Mut()
}
func (s *BarMut) SetBid(v *Candle) *BarMut {
	s.bid = *v
	return s
}
func (s *BarMut) Ask() *CandleMut {
	return s.ask.Mut()
}
func (s *BarMut) SetAsk(v *Candle) *BarMut {
	s.ask = *v
	return s
}
func (s *BarMut) Spread() *SpreadMut {
	return s.spread.Mut()
}
func (s *BarMut) SetSpread(v *Spread) *BarMut {
	s.spread = *v
	return s
}
func (s *BarMut) Ticks() *TicksMut {
	return s.ticks.Mut()
}
func (s *BarMut) SetTicks(v *Ticks) *BarMut {
	s.ticks = *v
	return s
}
func (s *BarMut) Volume() *VolumeMut {
	return s.volume.Mut()
}
func (s *BarMut) SetVolume(v *Volume) *BarMut {
	s.volume = *v
	return s
}
func (s *BarMut) Interest() *VolumeMut {
	return s.interest.Mut()
}
func (s *BarMut) SetInterest(v *Volume) *BarMut {
	s.interest = *v
	return s
}

type Trade struct {
	id        int64
	price     float64
	quantity  float64
	aggressor Aggressor
	_         [4]byte // Padding
}

func (s *Trade) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Trade) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id()
	m["price"] = s.Price()
	m["quantity"] = s.Quantity()
	m["aggressor"] = s.Aggressor()
	return m
}

func (s *Trade) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 32 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Trade) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[32]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Trade) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Trade) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Trade) Read(b []byte) (n int, err error) {
	if len(b) < 32 {
		return -1, io.ErrShortBuffer
	}
	v := (*Trade)(unsafe.Pointer(&b[0]))
	*v = *s
	return 32, nil
}
func (s *Trade) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*Trade)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Trade) Clone() *Trade {
	v := &Trade{}
	*v = *s
	return v
}
func (s *Trade) Bytes() []byte {
	return (*(*[32]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Trade) Mut() *TradeMut {
	return (*TradeMut)(unsafe.Pointer(s))
}
func (s *Trade) Id() int64 {
	return s.id
}
func (s *Trade) Price() float64 {
	return s.price
}
func (s *Trade) Quantity() float64 {
	return s.quantity
}
func (s *Trade) Aggressor() Aggressor {
	return s.aggressor
}

type TradeMut struct {
	Trade
}

func (s *TradeMut) Clone() *TradeMut {
	v := &TradeMut{}
	*v = *s
	return v
}
func (s *TradeMut) Freeze() *Trade {
	return (*Trade)(unsafe.Pointer(s))
}
func (s *TradeMut) SetId(v int64) *TradeMut {
	s.id = v
	return s
}
func (s *TradeMut) SetPrice(v float64) *TradeMut {
	s.price = v
	return s
}
func (s *TradeMut) SetQuantity(v float64) *TradeMut {
	s.quantity = v
	return s
}
func (s *TradeMut) SetAggressor(v Aggressor) *TradeMut {
	s.aggressor = v
	return s
}

type Spread struct {
	low  float64
	high float64
	avg  float64
}

func (s *Spread) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Spread) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["low"] = s.Low()
	m["high"] = s.High()
	m["avg"] = s.Avg()
	return m
}

func (s *Spread) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Spread) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Spread) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Spread) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Spread) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*Spread)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *Spread) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*Spread)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Spread) Clone() *Spread {
	v := &Spread{}
	*v = *s
	return v
}
func (s *Spread) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Spread) Mut() *SpreadMut {
	return (*SpreadMut)(unsafe.Pointer(s))
}
func (s *Spread) Low() float64 {
	return s.low
}
func (s *Spread) High() float64 {
	return s.high
}
func (s *Spread) Avg() float64 {
	return s.avg
}

type SpreadMut struct {
	Spread
}

func (s *SpreadMut) Clone() *SpreadMut {
	v := &SpreadMut{}
	*v = *s
	return v
}
func (s *SpreadMut) Freeze() *Spread {
	return (*Spread)(unsafe.Pointer(s))
}
func (s *SpreadMut) SetLow(v float64) *SpreadMut {
	s.low = v
	return s
}
func (s *SpreadMut) SetHigh(v float64) *SpreadMut {
	s.high = v
	return s
}
func (s *SpreadMut) SetAvg(v float64) *SpreadMut {
	s.avg = v
	return s
}

// Greeks are financial measures of the sensitivity of an option’s price to its
// underlying determining parameters, such as volatility or the price of the underlying
// asset. The Greeks are utilized in the analysis of an options portfolio and in sensitivity
// analysis of an option or portfolio of options. The measures are considered essential by
// many investors for making informed decisions in options trading.
//
// Delta, Gamma, Vega, Theta, and Rho are the key option Greeks. However, there are many other
// option Greeks that can be derived from those mentioned above.
type Greeks struct {
	iv    float64
	delta float64
	gamma float64
	vega  float64
	theta float64
	rho   float64
}

func (s *Greeks) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Greeks) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["iv"] = s.Iv()
	m["delta"] = s.Delta()
	m["gamma"] = s.Gamma()
	m["vega"] = s.Vega()
	m["theta"] = s.Theta()
	m["rho"] = s.Rho()
	return m
}

func (s *Greeks) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[48]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 48 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Greeks) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[48]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Greeks) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Greeks) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Greeks) Read(b []byte) (n int, err error) {
	if len(b) < 48 {
		return -1, io.ErrShortBuffer
	}
	v := (*Greeks)(unsafe.Pointer(&b[0]))
	*v = *s
	return 48, nil
}
func (s *Greeks) UnmarshalBinary(b []byte) error {
	if len(b) < 48 {
		return io.ErrShortBuffer
	}
	v := (*Greeks)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Greeks) Clone() *Greeks {
	v := &Greeks{}
	*v = *s
	return v
}
func (s *Greeks) Bytes() []byte {
	return (*(*[48]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Greeks) Mut() *GreeksMut {
	return (*GreeksMut)(unsafe.Pointer(s))
}
func (s *Greeks) Iv() float64 {
	return s.iv
}
func (s *Greeks) Delta() float64 {
	return s.delta
}
func (s *Greeks) Gamma() float64 {
	return s.gamma
}
func (s *Greeks) Vega() float64 {
	return s.vega
}
func (s *Greeks) Theta() float64 {
	return s.theta
}
func (s *Greeks) Rho() float64 {
	return s.rho
}

// Greeks are financial measures of the sensitivity of an option’s price to its
// underlying determining parameters, such as volatility or the price of the underlying
// asset. The Greeks are utilized in the analysis of an options portfolio and in sensitivity
// analysis of an option or portfolio of options. The measures are considered essential by
// many investors for making informed decisions in options trading.
//
// Delta, Gamma, Vega, Theta, and Rho are the key option Greeks. However, there are many other
// option Greeks that can be derived from those mentioned above.
type GreeksMut struct {
	Greeks
}

func (s *GreeksMut) Clone() *GreeksMut {
	v := &GreeksMut{}
	*v = *s
	return v
}
func (s *GreeksMut) Freeze() *Greeks {
	return (*Greeks)(unsafe.Pointer(s))
}
func (s *GreeksMut) SetIv(v float64) *GreeksMut {
	s.iv = v
	return s
}
func (s *GreeksMut) SetDelta(v float64) *GreeksMut {
	s.delta = v
	return s
}
func (s *GreeksMut) SetGamma(v float64) *GreeksMut {
	s.gamma = v
	return s
}
func (s *GreeksMut) SetVega(v float64) *GreeksMut {
	s.vega = v
	return s
}
func (s *GreeksMut) SetTheta(v float64) *GreeksMut {
	s.theta = v
	return s
}
func (s *GreeksMut) SetRho(v float64) *GreeksMut {
	s.rho = v
	return s
}

type OptionBar struct {
	start    int64
	end      int64
	price    Candle
	bid      Candle
	ask      Candle
	spread   Spread
	ticks    Ticks
	volume   Volume
	interest Volume
	greeks   Greeks
}

func (s *OptionBar) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *OptionBar) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["start"] = s.Start()
	m["end"] = s.End()
	m["price"] = s.Price().MarshalMap(nil)
	m["bid"] = s.Bid().MarshalMap(nil)
	m["ask"] = s.Ask().MarshalMap(nil)
	m["spread"] = s.Spread().MarshalMap(nil)
	m["ticks"] = s.Ticks().MarshalMap(nil)
	m["volume"] = s.Volume().MarshalMap(nil)
	m["interest"] = s.Interest().MarshalMap(nil)
	m["greeks"] = s.Greeks().MarshalMap(nil)
	return m
}

func (s *OptionBar) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[256]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 256 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *OptionBar) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[256]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *OptionBar) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[256]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *OptionBar) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[256]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *OptionBar) Read(b []byte) (n int, err error) {
	if len(b) < 256 {
		return -1, io.ErrShortBuffer
	}
	v := (*OptionBar)(unsafe.Pointer(&b[0]))
	*v = *s
	return 256, nil
}
func (s *OptionBar) UnmarshalBinary(b []byte) error {
	if len(b) < 256 {
		return io.ErrShortBuffer
	}
	v := (*OptionBar)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *OptionBar) Clone() *OptionBar {
	v := &OptionBar{}
	*v = *s
	return v
}
func (s *OptionBar) Bytes() []byte {
	return (*(*[256]byte)(unsafe.Pointer(s)))[0:]
}
func (s *OptionBar) Mut() *OptionBarMut {
	return (*OptionBarMut)(unsafe.Pointer(s))
}
func (s *OptionBar) Start() int64 {
	return s.start
}
func (s *OptionBar) End() int64 {
	return s.end
}
func (s *OptionBar) Price() *Candle {
	return &s.price
}
func (s *OptionBar) Bid() *Candle {
	return &s.bid
}
func (s *OptionBar) Ask() *Candle {
	return &s.ask
}
func (s *OptionBar) Spread() *Spread {
	return &s.spread
}
func (s *OptionBar) Ticks() *Ticks {
	return &s.ticks
}
func (s *OptionBar) Volume() *Volume {
	return &s.volume
}
func (s *OptionBar) Interest() *Volume {
	return &s.interest
}
func (s *OptionBar) Greeks() *Greeks {
	return &s.greeks
}

type OptionBarMut struct {
	OptionBar
}

func (s *OptionBarMut) Clone() *OptionBarMut {
	v := &OptionBarMut{}
	*v = *s
	return v
}
func (s *OptionBarMut) Freeze() *OptionBar {
	return (*OptionBar)(unsafe.Pointer(s))
}
func (s *OptionBarMut) SetStart(v int64) *OptionBarMut {
	s.start = v
	return s
}
func (s *OptionBarMut) SetEnd(v int64) *OptionBarMut {
	s.end = v
	return s
}
func (s *OptionBarMut) Price() *CandleMut {
	return s.price.Mut()
}
func (s *OptionBarMut) SetPrice(v *Candle) *OptionBarMut {
	s.price = *v
	return s
}
func (s *OptionBarMut) Bid() *CandleMut {
	return s.bid.Mut()
}
func (s *OptionBarMut) SetBid(v *Candle) *OptionBarMut {
	s.bid = *v
	return s
}
func (s *OptionBarMut) Ask() *CandleMut {
	return s.ask.Mut()
}
func (s *OptionBarMut) SetAsk(v *Candle) *OptionBarMut {
	s.ask = *v
	return s
}
func (s *OptionBarMut) Spread() *SpreadMut {
	return s.spread.Mut()
}
func (s *OptionBarMut) SetSpread(v *Spread) *OptionBarMut {
	s.spread = *v
	return s
}
func (s *OptionBarMut) Ticks() *TicksMut {
	return s.ticks.Mut()
}
func (s *OptionBarMut) SetTicks(v *Ticks) *OptionBarMut {
	s.ticks = *v
	return s
}
func (s *OptionBarMut) Volume() *VolumeMut {
	return s.volume.Mut()
}
func (s *OptionBarMut) SetVolume(v *Volume) *OptionBarMut {
	s.volume = *v
	return s
}
func (s *OptionBarMut) Interest() *VolumeMut {
	return s.interest.Mut()
}
func (s *OptionBarMut) SetInterest(v *Volume) *OptionBarMut {
	s.interest = *v
	return s
}
func (s *OptionBarMut) Greeks() *GreeksMut {
	return s.greeks.Mut()
}
func (s *OptionBarMut) SetGreeks(v *Greeks) *OptionBarMut {
	s.greeks = *v
	return s
}

type Liquidations struct {
	trades int64
	min    float64
	avg    float64
	max    float64
	buys   float64
	sells  float64
	value  float64
}

func (s *Liquidations) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Liquidations) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["trades"] = s.Trades()
	m["min"] = s.Min()
	m["avg"] = s.Avg()
	m["max"] = s.Max()
	m["buys"] = s.Buys()
	m["sells"] = s.Sells()
	m["value"] = s.Value()
	return m
}

func (s *Liquidations) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[56]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 56 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Liquidations) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[56]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Liquidations) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[56]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Liquidations) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[56]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Liquidations) Read(b []byte) (n int, err error) {
	if len(b) < 56 {
		return -1, io.ErrShortBuffer
	}
	v := (*Liquidations)(unsafe.Pointer(&b[0]))
	*v = *s
	return 56, nil
}
func (s *Liquidations) UnmarshalBinary(b []byte) error {
	if len(b) < 56 {
		return io.ErrShortBuffer
	}
	v := (*Liquidations)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Liquidations) Clone() *Liquidations {
	v := &Liquidations{}
	*v = *s
	return v
}
func (s *Liquidations) Bytes() []byte {
	return (*(*[56]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Liquidations) Mut() *LiquidationsMut {
	return (*LiquidationsMut)(unsafe.Pointer(s))
}
func (s *Liquidations) Trades() int64 {
	return s.trades
}
func (s *Liquidations) Min() float64 {
	return s.min
}
func (s *Liquidations) Avg() float64 {
	return s.avg
}
func (s *Liquidations) Max() float64 {
	return s.max
}
func (s *Liquidations) Buys() float64 {
	return s.buys
}
func (s *Liquidations) Sells() float64 {
	return s.sells
}
func (s *Liquidations) Value() float64 {
	return s.value
}

type LiquidationsMut struct {
	Liquidations
}

func (s *LiquidationsMut) Clone() *LiquidationsMut {
	v := &LiquidationsMut{}
	*v = *s
	return v
}
func (s *LiquidationsMut) Freeze() *Liquidations {
	return (*Liquidations)(unsafe.Pointer(s))
}
func (s *LiquidationsMut) SetTrades(v int64) *LiquidationsMut {
	s.trades = v
	return s
}
func (s *LiquidationsMut) SetMin(v float64) *LiquidationsMut {
	s.min = v
	return s
}
func (s *LiquidationsMut) SetAvg(v float64) *LiquidationsMut {
	s.avg = v
	return s
}
func (s *LiquidationsMut) SetMax(v float64) *LiquidationsMut {
	s.max = v
	return s
}
func (s *LiquidationsMut) SetBuys(v float64) *LiquidationsMut {
	s.buys = v
	return s
}
func (s *LiquidationsMut) SetSells(v float64) *LiquidationsMut {
	s.sells = v
	return s
}
func (s *LiquidationsMut) SetValue(v float64) *LiquidationsMut {
	s.value = v
	return s
}

type Time struct {
	start int64
	end   int64
}

func (s *Time) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Time) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["start"] = s.Start()
	m["end"] = s.End()
	return m
}

func (s *Time) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[16]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 16 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Time) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[16]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Time) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Time) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Time) Read(b []byte) (n int, err error) {
	if len(b) < 16 {
		return -1, io.ErrShortBuffer
	}
	v := (*Time)(unsafe.Pointer(&b[0]))
	*v = *s
	return 16, nil
}
func (s *Time) UnmarshalBinary(b []byte) error {
	if len(b) < 16 {
		return io.ErrShortBuffer
	}
	v := (*Time)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Time) Clone() *Time {
	v := &Time{}
	*v = *s
	return v
}
func (s *Time) Bytes() []byte {
	return (*(*[16]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Time) Mut() *TimeMut {
	return (*TimeMut)(unsafe.Pointer(s))
}
func (s *Time) Start() int64 {
	return s.start
}
func (s *Time) End() int64 {
	return s.end
}

type TimeMut struct {
	Time
}

func (s *TimeMut) Clone() *TimeMut {
	v := &TimeMut{}
	*v = *s
	return v
}
func (s *TimeMut) Freeze() *Time {
	return (*Time)(unsafe.Pointer(s))
}
func (s *TimeMut) SetStart(v int64) *TimeMut {
	s.start = v
	return s
}
func (s *TimeMut) SetEnd(v int64) *TimeMut {
	s.end = v
	return s
}
func init() {
	{
		var b [2]byte
		v := uint16(1)
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		if *(*uint16)(unsafe.Pointer(&b[0])) != 1 {
			panic("BigEndian not supported")
		}
	}
	type b struct {
		n    string
		o, s uintptr
	}
	a := func(x interface{}, y interface{}, s uintptr, z []b) {
		t := reflect.TypeOf(x)
		r := reflect.TypeOf(y)
		if t.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", t.Name(), t.Size(), s))
		}
		if r.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", r.Name(), r.Size(), s))
		}
		if t.NumField() != len(z) {
			panic(fmt.Sprintf("%s field count = %d: expected %d", t.Name(), t.NumField(), len(z)))
		}
		for i, e := range z {
			f := t.Field(i)
			if f.Offset != e.o {
				panic(fmt.Sprintf("%s.%s offset = %d, expected = %d", t.Name(), f.Name, f.Offset, e.o))
			}
			if f.Type.Size() != e.s {
				panic(fmt.Sprintf("%s.%s size = %d, expected = %d", t.Name(), f.Name, f.Type.Size(), e.s))
			}
			if f.Name != e.n {
				panic(fmt.Sprintf("%s.%s expected field: %s", t.Name(), f.Name, e.n))
			}
		}
	}

	a(Candle{}, CandleMut{}, 32, []b{
		{"open", 0, 8},
		{"high", 8, 8},
		{"low", 16, 8},
		{"close", 24, 8},
	})
	a(Ticks{}, TicksMut{}, 24, []b{
		{"total", 0, 8},
		{"up", 8, 8},
		{"down", 16, 8},
	})
	a(Volume{}, VolumeMut{}, 24, []b{
		{"total", 0, 8},
		{"buy", 8, 8},
		{"sell", 16, 8},
	})
	a(Trades{}, TradesMut{}, 24, []b{
		{"count", 0, 8},
		{"min", 8, 8},
		{"max", 16, 8},
	})
	a(Bar{}, BarMut{}, 208, []b{
		{"start", 0, 8},
		{"end", 8, 8},
		{"price", 16, 32},
		{"bid", 48, 32},
		{"ask", 80, 32},
		{"spread", 112, 24},
		{"ticks", 136, 24},
		{"volume", 160, 24},
		{"interest", 184, 24},
	})
	a(Trade{}, TradeMut{}, 32, []b{
		{"id", 0, 8},
		{"price", 8, 8},
		{"quantity", 16, 8},
		{"aggressor", 24, 4},
		{"_", 28, 4},
	})
	a(Spread{}, SpreadMut{}, 24, []b{
		{"low", 0, 8},
		{"high", 8, 8},
		{"avg", 16, 8},
	})
	a(Greeks{}, GreeksMut{}, 48, []b{
		{"iv", 0, 8},
		{"delta", 8, 8},
		{"gamma", 16, 8},
		{"vega", 24, 8},
		{"theta", 32, 8},
		{"rho", 40, 8},
	})
	a(OptionBar{}, OptionBarMut{}, 256, []b{
		{"start", 0, 8},
		{"end", 8, 8},
		{"price", 16, 32},
		{"bid", 48, 32},
		{"ask", 80, 32},
		{"spread", 112, 24},
		{"ticks", 136, 24},
		{"volume", 160, 24},
		{"interest", 184, 24},
		{"greeks", 208, 48},
	})
	a(Liquidations{}, LiquidationsMut{}, 56, []b{
		{"trades", 0, 8},
		{"min", 8, 8},
		{"avg", 16, 8},
		{"max", 24, 8},
		{"buys", 32, 8},
		{"sells", 40, 8},
		{"value", 48, 8},
	})
	a(Time{}, TimeMut{}, 16, []b{
		{"start", 0, 8},
		{"end", 8, 8},
	})

}
