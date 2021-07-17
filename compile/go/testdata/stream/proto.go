// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package stream

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type MessageType byte

const (
	MessageType_Record = MessageType(1)

	MessageType_Records = MessageType(2)

	MessageType_Block = MessageType(3)

	MessageType_EOS = MessageType(4)

	MessageType_EOSWaiting = MessageType(5)

	MessageType_Savepoint = MessageType(6)

	MessageType_Starting = MessageType(7)

	MessageType_Started = MessageType(8)

	MessageType_Stopped = MessageType(9)

	MessageType_SyncStarted = MessageType(10)

	MessageType_SyncProgress = MessageType(11)

	MessageType_SyncStopped = MessageType(12)
)

type StreamKind byte

const (
	StreamKind_Log = StreamKind(0)

	StreamKind_TimeSeries = StreamKind(1)

	StreamKind_Table = StreamKind(2)
)

type SchemaKind byte

const (
	SchemaKind_Bytes = SchemaKind(0)

	SchemaKind_Wasmbuf = SchemaKind(1)

	SchemaKind_Protobuf = SchemaKind(2)

	SchemaKind_Flatbuffers = SchemaKind(3)

	SchemaKind_Json = SchemaKind(4)

	SchemaKind_Msgpack = SchemaKind(5)
)

type Compression byte

const (
	Compression_None = Compression(0)

	Compression_LZ4 = Compression(1)
)

type SyncStoppedReason byte

const (
	SyncStoppedReason_Success = SyncStoppedReason(1)

	SyncStoppedReason_Error = SyncStoppedReason(2)
)

type StopReason byte

const (
	// Stream is composed from another stream or external datasource and it stopped
	StopReason_Source = StopReason(1)

	// Stream has been paused
	StopReason_Paused = StopReason(2)

	// Stream is being migrated to a new writer
	StopReason_Migrate = StopReason(3)

	// Stream has stopped unexpectedly
	StopReason_Unexpected = StopReason(4)
)

type RecordID struct {
	streamID int64
	blockID  int64
	id       int64
}

func (s *RecordID) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *RecordID) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["streamID"] = s.StreamID()
	m["blockID"] = s.BlockID()
	m["id"] = s.Id()
	return m
}

func (s *RecordID) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *RecordID) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *RecordID) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *RecordID) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *RecordID) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*RecordID)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *RecordID) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*RecordID)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *RecordID) Clone() *RecordID {
	v := &RecordID{}
	*v = *s
	return v
}
func (s *RecordID) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *RecordID) Mut() *RecordIDMut {
	return (*RecordIDMut)(unsafe.Pointer(s))
}
func (s *RecordID) StreamID() int64 {
	return s.streamID
}
func (s *RecordID) BlockID() int64 {
	return s.blockID
}
func (s *RecordID) Id() int64 {
	return s.id
}

type RecordIDMut struct {
	RecordID
}

func (s *RecordIDMut) Clone() *RecordIDMut {
	v := &RecordIDMut{}
	*v = *s
	return v
}
func (s *RecordIDMut) Freeze() *RecordID {
	return (*RecordID)(unsafe.Pointer(s))
}
func (s *RecordIDMut) SetStreamID(v int64) *RecordIDMut {
	s.streamID = v
	return s
}
func (s *RecordIDMut) SetBlockID(v int64) *RecordIDMut {
	s.blockID = v
	return s
}
func (s *RecordIDMut) SetId(v int64) *RecordIDMut {
	s.id = v
	return s
}

// Block
type Block struct {
	header BlockHeader
	data   Bytes65160
}

func (s *Block) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Block) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["header"] = s.Header().MarshalMap(nil)
	m["data"] = s.Data()
	return m
}

func (s *Block) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[65256]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 65256 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Block) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[65256]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Block) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[65256]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Block) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[65256]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Block) Read(b []byte) (n int, err error) {
	if len(b) < 65256 {
		return -1, io.ErrShortBuffer
	}
	v := (*Block)(unsafe.Pointer(&b[0]))
	*v = *s
	return 65256, nil
}
func (s *Block) UnmarshalBinary(b []byte) error {
	if len(b) < 65256 {
		return io.ErrShortBuffer
	}
	v := (*Block)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Block) Clone() *Block {
	v := &Block{}
	*v = *s
	return v
}
func (s *Block) Bytes() []byte {
	return (*(*[65256]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Block) Mut() *BlockMut {
	return (*BlockMut)(unsafe.Pointer(s))
}
func (s *Block) Header() *BlockHeader {
	return &s.header
}
func (s *Block) Data() Bytes65160 {
	return s.data
}

// Block
type BlockMut struct {
	Block
}

func (s *BlockMut) Clone() *BlockMut {
	v := &BlockMut{}
	*v = *s
	return v
}
func (s *BlockMut) Freeze() *Block {
	return (*Block)(unsafe.Pointer(s))
}
func (s *BlockMut) Header() *BlockHeaderMut {
	return s.header.Mut()
}
func (s *BlockMut) SetHeader(v *BlockHeader) *BlockMut {
	s.header = *v
	return s
}
func (s *BlockMut) SetData(v Bytes65160) *BlockMut {
	s.data = v
	return s
}

// EOS = End of Stream
// The reader is caught up on the stream and is NOT subscribed
// to new records.
type EOS struct {
	recordID  RecordID
	timestamp int64
}

func (s *EOS) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *EOS) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	return m
}

func (s *EOS) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 32 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *EOS) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[32]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *EOS) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *EOS) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *EOS) Read(b []byte) (n int, err error) {
	if len(b) < 32 {
		return -1, io.ErrShortBuffer
	}
	v := (*EOS)(unsafe.Pointer(&b[0]))
	*v = *s
	return 32, nil
}
func (s *EOS) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*EOS)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *EOS) Clone() *EOS {
	v := &EOS{}
	*v = *s
	return v
}
func (s *EOS) Bytes() []byte {
	return (*(*[32]byte)(unsafe.Pointer(s)))[0:]
}
func (s *EOS) Mut() *EOSMut {
	return (*EOSMut)(unsafe.Pointer(s))
}
func (s *EOS) RecordID() *RecordID {
	return &s.recordID
}
func (s *EOS) Timestamp() int64 {
	return s.timestamp
}

// EOS = End of Stream
// The reader is caught up on the stream and is NOT subscribed
// to new records.
type EOSMut struct {
	EOS
}

func (s *EOSMut) Clone() *EOSMut {
	v := &EOSMut{}
	*v = *s
	return v
}
func (s *EOSMut) Freeze() *EOS {
	return (*EOS)(unsafe.Pointer(s))
}
func (s *EOSMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *EOSMut) SetRecordID(v *RecordID) *EOSMut {
	s.recordID = *v
	return s
}
func (s *EOSMut) SetTimestamp(v int64) *EOSMut {
	s.timestamp = v
	return s
}

type Stopped struct {
	recordID  RecordID
	timestamp int64
	reason    StopReason
	_         [7]byte // Padding
	starts    int64
}

func (s *Stopped) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Stopped) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	m["reason"] = s.Reason()
	m["starts"] = s.Starts()
	return m
}

func (s *Stopped) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[48]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 48 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Stopped) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[48]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Stopped) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Stopped) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Stopped) Read(b []byte) (n int, err error) {
	if len(b) < 48 {
		return -1, io.ErrShortBuffer
	}
	v := (*Stopped)(unsafe.Pointer(&b[0]))
	*v = *s
	return 48, nil
}
func (s *Stopped) UnmarshalBinary(b []byte) error {
	if len(b) < 48 {
		return io.ErrShortBuffer
	}
	v := (*Stopped)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Stopped) Clone() *Stopped {
	v := &Stopped{}
	*v = *s
	return v
}
func (s *Stopped) Bytes() []byte {
	return (*(*[48]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Stopped) Mut() *StoppedMut {
	return (*StoppedMut)(unsafe.Pointer(s))
}
func (s *Stopped) RecordID() *RecordID {
	return &s.recordID
}
func (s *Stopped) Timestamp() int64 {
	return s.timestamp
}
func (s *Stopped) Reason() StopReason {
	return s.reason
}
func (s *Stopped) Starts() int64 {
	return s.starts
}

type StoppedMut struct {
	Stopped
}

func (s *StoppedMut) Clone() *StoppedMut {
	v := &StoppedMut{}
	*v = *s
	return v
}
func (s *StoppedMut) Freeze() *Stopped {
	return (*Stopped)(unsafe.Pointer(s))
}
func (s *StoppedMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *StoppedMut) SetRecordID(v *RecordID) *StoppedMut {
	s.recordID = *v
	return s
}
func (s *StoppedMut) SetTimestamp(v int64) *StoppedMut {
	s.timestamp = v
	return s
}
func (s *StoppedMut) SetReason(v StopReason) *StoppedMut {
	s.reason = v
	return s
}
func (s *StoppedMut) SetStarts(v int64) *StoppedMut {
	s.starts = v
	return s
}

type RecordsHeader struct {
	header RecordHeader
	count  uint16
	record uint16
	_      [4]byte // Padding
}

func (s *RecordsHeader) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *RecordsHeader) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["header"] = s.Header().MarshalMap(nil)
	m["count"] = s.Count()
	m["record"] = s.Record()
	return m
}

func (s *RecordsHeader) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[72]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 72 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *RecordsHeader) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[72]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *RecordsHeader) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[72]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *RecordsHeader) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[72]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *RecordsHeader) Read(b []byte) (n int, err error) {
	if len(b) < 72 {
		return -1, io.ErrShortBuffer
	}
	v := (*RecordsHeader)(unsafe.Pointer(&b[0]))
	*v = *s
	return 72, nil
}
func (s *RecordsHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 72 {
		return io.ErrShortBuffer
	}
	v := (*RecordsHeader)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *RecordsHeader) Clone() *RecordsHeader {
	v := &RecordsHeader{}
	*v = *s
	return v
}
func (s *RecordsHeader) Bytes() []byte {
	return (*(*[72]byte)(unsafe.Pointer(s)))[0:]
}
func (s *RecordsHeader) Mut() *RecordsHeaderMut {
	return (*RecordsHeaderMut)(unsafe.Pointer(s))
}
func (s *RecordsHeader) Header() *RecordHeader {
	return &s.header
}
func (s *RecordsHeader) Count() uint16 {
	return s.count
}
func (s *RecordsHeader) Record() uint16 {
	return s.record
}

type RecordsHeaderMut struct {
	RecordsHeader
}

func (s *RecordsHeaderMut) Clone() *RecordsHeaderMut {
	v := &RecordsHeaderMut{}
	*v = *s
	return v
}
func (s *RecordsHeaderMut) Freeze() *RecordsHeader {
	return (*RecordsHeader)(unsafe.Pointer(s))
}
func (s *RecordsHeaderMut) Header() *RecordHeaderMut {
	return s.header.Mut()
}
func (s *RecordsHeaderMut) SetHeader(v *RecordHeader) *RecordsHeaderMut {
	s.header = *v
	return s
}
func (s *RecordsHeaderMut) SetCount(v uint16) *RecordsHeaderMut {
	s.count = v
	return s
}
func (s *RecordsHeaderMut) SetRecord(v uint16) *RecordsHeaderMut {
	s.record = v
	return s
}

type StreamStats struct {
	storage  Stats
	appender Stats
}

func (s *StreamStats) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *StreamStats) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["storage"] = s.Storage().MarshalMap(nil)
	m["appender"] = s.Appender().MarshalMap(nil)
	return m
}

func (s *StreamStats) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[48]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 48 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *StreamStats) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[48]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *StreamStats) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *StreamStats) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *StreamStats) Read(b []byte) (n int, err error) {
	if len(b) < 48 {
		return -1, io.ErrShortBuffer
	}
	v := (*StreamStats)(unsafe.Pointer(&b[0]))
	*v = *s
	return 48, nil
}
func (s *StreamStats) UnmarshalBinary(b []byte) error {
	if len(b) < 48 {
		return io.ErrShortBuffer
	}
	v := (*StreamStats)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *StreamStats) Clone() *StreamStats {
	v := &StreamStats{}
	*v = *s
	return v
}
func (s *StreamStats) Bytes() []byte {
	return (*(*[48]byte)(unsafe.Pointer(s)))[0:]
}
func (s *StreamStats) Mut() *StreamStatsMut {
	return (*StreamStatsMut)(unsafe.Pointer(s))
}
func (s *StreamStats) Storage() *Stats {
	return &s.storage
}
func (s *StreamStats) Appender() *Stats {
	return &s.appender
}

type StreamStatsMut struct {
	StreamStats
}

func (s *StreamStatsMut) Clone() *StreamStatsMut {
	v := &StreamStatsMut{}
	*v = *s
	return v
}
func (s *StreamStatsMut) Freeze() *StreamStats {
	return (*StreamStats)(unsafe.Pointer(s))
}
func (s *StreamStatsMut) Storage() *StatsMut {
	return s.storage.Mut()
}
func (s *StreamStatsMut) SetStorage(v *Stats) *StreamStatsMut {
	s.storage = *v
	return s
}
func (s *StreamStatsMut) Appender() *StatsMut {
	return s.appender.Mut()
}
func (s *StreamStatsMut) SetAppender(v *Stats) *StreamStatsMut {
	s.appender = *v
	return s
}

type SyncProgress struct {
	recordID  RecordID
	timestamp int64
	started   int64
	count     int64
	remaining int64
}

func (s *SyncProgress) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *SyncProgress) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	m["started"] = s.Started()
	m["count"] = s.Count()
	m["remaining"] = s.Remaining()
	return m
}

func (s *SyncProgress) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[56]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 56 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *SyncProgress) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[56]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *SyncProgress) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[56]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *SyncProgress) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[56]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *SyncProgress) Read(b []byte) (n int, err error) {
	if len(b) < 56 {
		return -1, io.ErrShortBuffer
	}
	v := (*SyncProgress)(unsafe.Pointer(&b[0]))
	*v = *s
	return 56, nil
}
func (s *SyncProgress) UnmarshalBinary(b []byte) error {
	if len(b) < 56 {
		return io.ErrShortBuffer
	}
	v := (*SyncProgress)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *SyncProgress) Clone() *SyncProgress {
	v := &SyncProgress{}
	*v = *s
	return v
}
func (s *SyncProgress) Bytes() []byte {
	return (*(*[56]byte)(unsafe.Pointer(s)))[0:]
}
func (s *SyncProgress) Mut() *SyncProgressMut {
	return (*SyncProgressMut)(unsafe.Pointer(s))
}
func (s *SyncProgress) RecordID() *RecordID {
	return &s.recordID
}
func (s *SyncProgress) Timestamp() int64 {
	return s.timestamp
}
func (s *SyncProgress) Started() int64 {
	return s.started
}
func (s *SyncProgress) Count() int64 {
	return s.count
}
func (s *SyncProgress) Remaining() int64 {
	return s.remaining
}

type SyncProgressMut struct {
	SyncProgress
}

func (s *SyncProgressMut) Clone() *SyncProgressMut {
	v := &SyncProgressMut{}
	*v = *s
	return v
}
func (s *SyncProgressMut) Freeze() *SyncProgress {
	return (*SyncProgress)(unsafe.Pointer(s))
}
func (s *SyncProgressMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *SyncProgressMut) SetRecordID(v *RecordID) *SyncProgressMut {
	s.recordID = *v
	return s
}
func (s *SyncProgressMut) SetTimestamp(v int64) *SyncProgressMut {
	s.timestamp = v
	return s
}
func (s *SyncProgressMut) SetStarted(v int64) *SyncProgressMut {
	s.started = v
	return s
}
func (s *SyncProgressMut) SetCount(v int64) *SyncProgressMut {
	s.count = v
	return s
}
func (s *SyncProgressMut) SetRemaining(v int64) *SyncProgressMut {
	s.remaining = v
	return s
}

type Stream struct {
	id        int64
	created   int64
	accountID int64
	duration  int64
	record    int32
	_         [4]byte // Padding
	name      String32
	kind      StreamKind
	schema    SchemaKind
	realTime  bool
	_         [5]byte // Padding
}

func (s *Stream) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Stream) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id()
	m["created"] = s.Created()
	m["accountID"] = s.AccountID()
	m["duration"] = s.Duration()
	m["record"] = s.Record()
	m["name"] = s.Name()
	m["kind"] = s.Kind()
	m["schema"] = s.Schema()
	m["realTime"] = s.RealTime()
	return m
}

func (s *Stream) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[80]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 80 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Stream) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[80]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Stream) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[80]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Stream) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[80]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Stream) Read(b []byte) (n int, err error) {
	if len(b) < 80 {
		return -1, io.ErrShortBuffer
	}
	v := (*Stream)(unsafe.Pointer(&b[0]))
	*v = *s
	return 80, nil
}
func (s *Stream) UnmarshalBinary(b []byte) error {
	if len(b) < 80 {
		return io.ErrShortBuffer
	}
	v := (*Stream)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Stream) Clone() *Stream {
	v := &Stream{}
	*v = *s
	return v
}
func (s *Stream) Bytes() []byte {
	return (*(*[80]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Stream) Mut() *StreamMut {
	return (*StreamMut)(unsafe.Pointer(s))
}
func (s *Stream) Id() int64 {
	return s.id
}
func (s *Stream) Created() int64 {
	return s.created
}
func (s *Stream) AccountID() int64 {
	return s.accountID
}
func (s *Stream) Duration() int64 {
	return s.duration
}
func (s *Stream) Record() int32 {
	return s.record
}
func (s *Stream) Name() *String32 {
	return &s.name
}
func (s *Stream) Kind() StreamKind {
	return s.kind
}
func (s *Stream) Schema() SchemaKind {
	return s.schema
}
func (s *Stream) RealTime() bool {
	return s.realTime
}

type StreamMut struct {
	Stream
}

func (s *StreamMut) Clone() *StreamMut {
	v := &StreamMut{}
	*v = *s
	return v
}
func (s *StreamMut) Freeze() *Stream {
	return (*Stream)(unsafe.Pointer(s))
}
func (s *StreamMut) SetId(v int64) *StreamMut {
	s.id = v
	return s
}
func (s *StreamMut) SetCreated(v int64) *StreamMut {
	s.created = v
	return s
}
func (s *StreamMut) SetAccountID(v int64) *StreamMut {
	s.accountID = v
	return s
}
func (s *StreamMut) SetDuration(v int64) *StreamMut {
	s.duration = v
	return s
}
func (s *StreamMut) SetRecord(v int32) *StreamMut {
	s.record = v
	return s
}
func (s *StreamMut) Name() *String32Mut {
	return s.name.Mut()
}
func (s *StreamMut) SetName(v *String32) *StreamMut {
	s.name = *v
	return s
}
func (s *StreamMut) SetKind(v StreamKind) *StreamMut {
	s.kind = v
	return s
}
func (s *StreamMut) SetSchema(v SchemaKind) *StreamMut {
	s.schema = v
	return s
}
func (s *StreamMut) SetRealTime(v bool) *StreamMut {
	s.realTime = v
	return s
}

// EOSWaiting = End of Stream Waiting for next record.
// The reader is caught up on the stream and is subscribed
// to new records.
type EOSWaiting struct {
	recordID  RecordID
	timestamp int64
}

func (s *EOSWaiting) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *EOSWaiting) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	return m
}

func (s *EOSWaiting) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 32 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *EOSWaiting) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[32]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *EOSWaiting) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *EOSWaiting) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *EOSWaiting) Read(b []byte) (n int, err error) {
	if len(b) < 32 {
		return -1, io.ErrShortBuffer
	}
	v := (*EOSWaiting)(unsafe.Pointer(&b[0]))
	*v = *s
	return 32, nil
}
func (s *EOSWaiting) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*EOSWaiting)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *EOSWaiting) Clone() *EOSWaiting {
	v := &EOSWaiting{}
	*v = *s
	return v
}
func (s *EOSWaiting) Bytes() []byte {
	return (*(*[32]byte)(unsafe.Pointer(s)))[0:]
}
func (s *EOSWaiting) Mut() *EOSWaitingMut {
	return (*EOSWaitingMut)(unsafe.Pointer(s))
}
func (s *EOSWaiting) RecordID() *RecordID {
	return &s.recordID
}
func (s *EOSWaiting) Timestamp() int64 {
	return s.timestamp
}

// EOSWaiting = End of Stream Waiting for next record.
// The reader is caught up on the stream and is subscribed
// to new records.
type EOSWaitingMut struct {
	EOSWaiting
}

func (s *EOSWaitingMut) Clone() *EOSWaitingMut {
	v := &EOSWaitingMut{}
	*v = *s
	return v
}
func (s *EOSWaitingMut) Freeze() *EOSWaiting {
	return (*EOSWaiting)(unsafe.Pointer(s))
}
func (s *EOSWaitingMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *EOSWaitingMut) SetRecordID(v *RecordID) *EOSWaitingMut {
	s.recordID = *v
	return s
}
func (s *EOSWaitingMut) SetTimestamp(v int64) *EOSWaitingMut {
	s.timestamp = v
	return s
}

type Stats struct {
	size   int64
	count  int64
	blocks int64
}

func (s *Stats) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Stats) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["size"] = s.Size()
	m["count"] = s.Count()
	m["blocks"] = s.Blocks()
	return m
}

func (s *Stats) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[24]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 24 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Stats) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[24]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Stats) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Stats) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[24]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Stats) Read(b []byte) (n int, err error) {
	if len(b) < 24 {
		return -1, io.ErrShortBuffer
	}
	v := (*Stats)(unsafe.Pointer(&b[0]))
	*v = *s
	return 24, nil
}
func (s *Stats) UnmarshalBinary(b []byte) error {
	if len(b) < 24 {
		return io.ErrShortBuffer
	}
	v := (*Stats)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Stats) Clone() *Stats {
	v := &Stats{}
	*v = *s
	return v
}
func (s *Stats) Bytes() []byte {
	return (*(*[24]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Stats) Mut() *StatsMut {
	return (*StatsMut)(unsafe.Pointer(s))
}
func (s *Stats) Size() int64 {
	return s.size
}
func (s *Stats) Count() int64 {
	return s.count
}
func (s *Stats) Blocks() int64 {
	return s.blocks
}

type StatsMut struct {
	Stats
}

func (s *StatsMut) Clone() *StatsMut {
	v := &StatsMut{}
	*v = *s
	return v
}
func (s *StatsMut) Freeze() *Stats {
	return (*Stats)(unsafe.Pointer(s))
}
func (s *StatsMut) SetSize(v int64) *StatsMut {
	s.size = v
	return s
}
func (s *StatsMut) SetCount(v int64) *StatsMut {
	s.count = v
	return s
}
func (s *StatsMut) SetBlocks(v int64) *StatsMut {
	s.blocks = v
	return s
}

type AccountStats struct {
	id       int64
	storage  Stats
	appender Stats
	streams  int64
}

func (s *AccountStats) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *AccountStats) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id()
	m["storage"] = s.Storage().MarshalMap(nil)
	m["appender"] = s.Appender().MarshalMap(nil)
	m["streams"] = s.Streams()
	return m
}

func (s *AccountStats) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[64]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 64 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *AccountStats) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[64]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *AccountStats) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[64]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *AccountStats) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[64]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *AccountStats) Read(b []byte) (n int, err error) {
	if len(b) < 64 {
		return -1, io.ErrShortBuffer
	}
	v := (*AccountStats)(unsafe.Pointer(&b[0]))
	*v = *s
	return 64, nil
}
func (s *AccountStats) UnmarshalBinary(b []byte) error {
	if len(b) < 64 {
		return io.ErrShortBuffer
	}
	v := (*AccountStats)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *AccountStats) Clone() *AccountStats {
	v := &AccountStats{}
	*v = *s
	return v
}
func (s *AccountStats) Bytes() []byte {
	return (*(*[64]byte)(unsafe.Pointer(s)))[0:]
}
func (s *AccountStats) Mut() *AccountStatsMut {
	return (*AccountStatsMut)(unsafe.Pointer(s))
}
func (s *AccountStats) Id() int64 {
	return s.id
}
func (s *AccountStats) Storage() *Stats {
	return &s.storage
}
func (s *AccountStats) Appender() *Stats {
	return &s.appender
}
func (s *AccountStats) Streams() int64 {
	return s.streams
}

type AccountStatsMut struct {
	AccountStats
}

func (s *AccountStatsMut) Clone() *AccountStatsMut {
	v := &AccountStatsMut{}
	*v = *s
	return v
}
func (s *AccountStatsMut) Freeze() *AccountStats {
	return (*AccountStats)(unsafe.Pointer(s))
}
func (s *AccountStatsMut) SetId(v int64) *AccountStatsMut {
	s.id = v
	return s
}
func (s *AccountStatsMut) Storage() *StatsMut {
	return s.storage.Mut()
}
func (s *AccountStatsMut) SetStorage(v *Stats) *AccountStatsMut {
	s.storage = *v
	return s
}
func (s *AccountStatsMut) Appender() *StatsMut {
	return s.appender.Mut()
}
func (s *AccountStatsMut) SetAppender(v *Stats) *AccountStatsMut {
	s.appender = *v
	return s
}
func (s *AccountStatsMut) SetStreams(v int64) *AccountStatsMut {
	s.streams = v
	return s
}

// BlockID represents a globally unique ID of a single page of a single stream.
// String representation
type BlockID struct {
	streamID int64
	id       int64
}

func (s *BlockID) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *BlockID) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["streamID"] = s.StreamID()
	m["id"] = s.Id()
	return m
}

func (s *BlockID) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[16]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 16 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *BlockID) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[16]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *BlockID) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *BlockID) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *BlockID) Read(b []byte) (n int, err error) {
	if len(b) < 16 {
		return -1, io.ErrShortBuffer
	}
	v := (*BlockID)(unsafe.Pointer(&b[0]))
	*v = *s
	return 16, nil
}
func (s *BlockID) UnmarshalBinary(b []byte) error {
	if len(b) < 16 {
		return io.ErrShortBuffer
	}
	v := (*BlockID)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *BlockID) Clone() *BlockID {
	v := &BlockID{}
	*v = *s
	return v
}
func (s *BlockID) Bytes() []byte {
	return (*(*[16]byte)(unsafe.Pointer(s)))[0:]
}
func (s *BlockID) Mut() *BlockIDMut {
	return (*BlockIDMut)(unsafe.Pointer(s))
}
func (s *BlockID) StreamID() int64 {
	return s.streamID
}
func (s *BlockID) Id() int64 {
	return s.id
}

// BlockID represents a globally unique ID of a single page of a single stream.
// String representation
type BlockIDMut struct {
	BlockID
}

func (s *BlockIDMut) Clone() *BlockIDMut {
	v := &BlockIDMut{}
	*v = *s
	return v
}
func (s *BlockIDMut) Freeze() *BlockID {
	return (*BlockID)(unsafe.Pointer(s))
}
func (s *BlockIDMut) SetStreamID(v int64) *BlockIDMut {
	s.streamID = v
	return s
}
func (s *BlockIDMut) SetId(v int64) *BlockIDMut {
	s.id = v
	return s
}

type Started struct {
	recordID  RecordID
	timestamp int64
	writerID  int64
	stops     int64
}

func (s *Started) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Started) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	m["writerID"] = s.WriterID()
	m["stops"] = s.Stops()
	return m
}

func (s *Started) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[48]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 48 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Started) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[48]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Started) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Started) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[48]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Started) Read(b []byte) (n int, err error) {
	if len(b) < 48 {
		return -1, io.ErrShortBuffer
	}
	v := (*Started)(unsafe.Pointer(&b[0]))
	*v = *s
	return 48, nil
}
func (s *Started) UnmarshalBinary(b []byte) error {
	if len(b) < 48 {
		return io.ErrShortBuffer
	}
	v := (*Started)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Started) Clone() *Started {
	v := &Started{}
	*v = *s
	return v
}
func (s *Started) Bytes() []byte {
	return (*(*[48]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Started) Mut() *StartedMut {
	return (*StartedMut)(unsafe.Pointer(s))
}
func (s *Started) RecordID() *RecordID {
	return &s.recordID
}
func (s *Started) Timestamp() int64 {
	return s.timestamp
}
func (s *Started) WriterID() int64 {
	return s.writerID
}
func (s *Started) Stops() int64 {
	return s.stops
}

type StartedMut struct {
	Started
}

func (s *StartedMut) Clone() *StartedMut {
	v := &StartedMut{}
	*v = *s
	return v
}
func (s *StartedMut) Freeze() *Started {
	return (*Started)(unsafe.Pointer(s))
}
func (s *StartedMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *StartedMut) SetRecordID(v *RecordID) *StartedMut {
	s.recordID = *v
	return s
}
func (s *StartedMut) SetTimestamp(v int64) *StartedMut {
	s.timestamp = v
	return s
}
func (s *StartedMut) SetWriterID(v int64) *StartedMut {
	s.writerID = v
	return s
}
func (s *StartedMut) SetStops(v int64) *StartedMut {
	s.stops = v
	return s
}

type SyncStopped struct {
	progress SyncProgress
	reason   SyncStoppedReason
	_        [7]byte // Padding
	message  String64
}

func (s *SyncStopped) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *SyncStopped) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["progress"] = s.Progress().MarshalMap(nil)
	m["reason"] = s.Reason()
	m["message"] = s.Message()
	return m
}

func (s *SyncStopped) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[128]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 128 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *SyncStopped) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[128]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *SyncStopped) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[128]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *SyncStopped) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[128]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *SyncStopped) Read(b []byte) (n int, err error) {
	if len(b) < 128 {
		return -1, io.ErrShortBuffer
	}
	v := (*SyncStopped)(unsafe.Pointer(&b[0]))
	*v = *s
	return 128, nil
}
func (s *SyncStopped) UnmarshalBinary(b []byte) error {
	if len(b) < 128 {
		return io.ErrShortBuffer
	}
	v := (*SyncStopped)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *SyncStopped) Clone() *SyncStopped {
	v := &SyncStopped{}
	*v = *s
	return v
}
func (s *SyncStopped) Bytes() []byte {
	return (*(*[128]byte)(unsafe.Pointer(s)))[0:]
}
func (s *SyncStopped) Mut() *SyncStoppedMut {
	return (*SyncStoppedMut)(unsafe.Pointer(s))
}
func (s *SyncStopped) Progress() *SyncProgress {
	return &s.progress
}
func (s *SyncStopped) Reason() SyncStoppedReason {
	return s.reason
}
func (s *SyncStopped) Message() *String64 {
	return &s.message
}

type SyncStoppedMut struct {
	SyncStopped
}

func (s *SyncStoppedMut) Clone() *SyncStoppedMut {
	v := &SyncStoppedMut{}
	*v = *s
	return v
}
func (s *SyncStoppedMut) Freeze() *SyncStopped {
	return (*SyncStopped)(unsafe.Pointer(s))
}
func (s *SyncStoppedMut) Progress() *SyncProgressMut {
	return s.progress.Mut()
}
func (s *SyncStoppedMut) SetProgress(v *SyncProgress) *SyncStoppedMut {
	s.progress = *v
	return s
}
func (s *SyncStoppedMut) SetReason(v SyncStoppedReason) *SyncStoppedMut {
	s.reason = v
	return s
}
func (s *SyncStoppedMut) Message() *String64Mut {
	return s.message.Mut()
}
func (s *SyncStoppedMut) SetMessage(v *String64) *SyncStoppedMut {
	s.message = *v
	return s
}

type Starting struct {
	recordID  RecordID
	timestamp int64
	writerID  int64
}

func (s *Starting) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Starting) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	m["writerID"] = s.WriterID()
	return m
}

func (s *Starting) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[40]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 40 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Starting) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[40]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Starting) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[40]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Starting) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[40]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Starting) Read(b []byte) (n int, err error) {
	if len(b) < 40 {
		return -1, io.ErrShortBuffer
	}
	v := (*Starting)(unsafe.Pointer(&b[0]))
	*v = *s
	return 40, nil
}
func (s *Starting) UnmarshalBinary(b []byte) error {
	if len(b) < 40 {
		return io.ErrShortBuffer
	}
	v := (*Starting)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Starting) Clone() *Starting {
	v := &Starting{}
	*v = *s
	return v
}
func (s *Starting) Bytes() []byte {
	return (*(*[40]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Starting) Mut() *StartingMut {
	return (*StartingMut)(unsafe.Pointer(s))
}
func (s *Starting) RecordID() *RecordID {
	return &s.recordID
}
func (s *Starting) Timestamp() int64 {
	return s.timestamp
}
func (s *Starting) WriterID() int64 {
	return s.writerID
}

type StartingMut struct {
	Starting
}

func (s *StartingMut) Clone() *StartingMut {
	v := &StartingMut{}
	*v = *s
	return v
}
func (s *StartingMut) Freeze() *Starting {
	return (*Starting)(unsafe.Pointer(s))
}
func (s *StartingMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *StartingMut) SetRecordID(v *RecordID) *StartingMut {
	s.recordID = *v
	return s
}
func (s *StartingMut) SetTimestamp(v int64) *StartingMut {
	s.timestamp = v
	return s
}
func (s *StartingMut) SetWriterID(v int64) *StartingMut {
	s.writerID = v
	return s
}

// BlockHeader
type BlockHeader struct {
	id        BlockID
	created   int64
	completed int64
	min       int64
	max       int64
	start     int64
	end       int64
	storage   uint64
	storageU  uint64
	count     uint16
	size      uint16
	sizeU     uint16
	sizeX     uint16
	record    uint16
	encoding  Compression
	kind      StreamKind
	schema    SchemaKind
	realTime  bool
	_         [2]byte // Padding
}

func (s *BlockHeader) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *BlockHeader) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id().MarshalMap(nil)
	m["created"] = s.Created()
	m["completed"] = s.Completed()
	m["min"] = s.Min()
	m["max"] = s.Max()
	m["start"] = s.Start()
	m["end"] = s.End()
	m["storage"] = s.Storage()
	m["storageU"] = s.StorageU()
	m["count"] = s.Count()
	m["size"] = s.Size()
	m["sizeU"] = s.SizeU()
	m["sizeX"] = s.SizeX()
	m["record"] = s.Record()
	m["encoding"] = s.Encoding()
	m["kind"] = s.Kind()
	m["schema"] = s.Schema()
	m["realTime"] = s.RealTime()
	return m
}

func (s *BlockHeader) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[96]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 96 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *BlockHeader) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[96]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *BlockHeader) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[96]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *BlockHeader) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[96]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *BlockHeader) Read(b []byte) (n int, err error) {
	if len(b) < 96 {
		return -1, io.ErrShortBuffer
	}
	v := (*BlockHeader)(unsafe.Pointer(&b[0]))
	*v = *s
	return 96, nil
}
func (s *BlockHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 96 {
		return io.ErrShortBuffer
	}
	v := (*BlockHeader)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *BlockHeader) Clone() *BlockHeader {
	v := &BlockHeader{}
	*v = *s
	return v
}
func (s *BlockHeader) Bytes() []byte {
	return (*(*[96]byte)(unsafe.Pointer(s)))[0:]
}
func (s *BlockHeader) Mut() *BlockHeaderMut {
	return (*BlockHeaderMut)(unsafe.Pointer(s))
}
func (s *BlockHeader) Id() *BlockID {
	return &s.id
}
func (s *BlockHeader) Created() int64 {
	return s.created
}
func (s *BlockHeader) Completed() int64 {
	return s.completed
}
func (s *BlockHeader) Min() int64 {
	return s.min
}
func (s *BlockHeader) Max() int64 {
	return s.max
}
func (s *BlockHeader) Start() int64 {
	return s.start
}
func (s *BlockHeader) End() int64 {
	return s.end
}
func (s *BlockHeader) Storage() uint64 {
	return s.storage
}
func (s *BlockHeader) StorageU() uint64 {
	return s.storageU
}
func (s *BlockHeader) Count() uint16 {
	return s.count
}
func (s *BlockHeader) Size() uint16 {
	return s.size
}
func (s *BlockHeader) SizeU() uint16 {
	return s.sizeU
}
func (s *BlockHeader) SizeX() uint16 {
	return s.sizeX
}
func (s *BlockHeader) Record() uint16 {
	return s.record
}
func (s *BlockHeader) Encoding() Compression {
	return s.encoding
}
func (s *BlockHeader) Kind() StreamKind {
	return s.kind
}
func (s *BlockHeader) Schema() SchemaKind {
	return s.schema
}
func (s *BlockHeader) RealTime() bool {
	return s.realTime
}

// BlockHeader
type BlockHeaderMut struct {
	BlockHeader
}

func (s *BlockHeaderMut) Clone() *BlockHeaderMut {
	v := &BlockHeaderMut{}
	*v = *s
	return v
}
func (s *BlockHeaderMut) Freeze() *BlockHeader {
	return (*BlockHeader)(unsafe.Pointer(s))
}
func (s *BlockHeaderMut) Id() *BlockIDMut {
	return s.id.Mut()
}
func (s *BlockHeaderMut) SetId(v *BlockID) *BlockHeaderMut {
	s.id = *v
	return s
}
func (s *BlockHeaderMut) SetCreated(v int64) *BlockHeaderMut {
	s.created = v
	return s
}
func (s *BlockHeaderMut) SetCompleted(v int64) *BlockHeaderMut {
	s.completed = v
	return s
}
func (s *BlockHeaderMut) SetMin(v int64) *BlockHeaderMut {
	s.min = v
	return s
}
func (s *BlockHeaderMut) SetMax(v int64) *BlockHeaderMut {
	s.max = v
	return s
}
func (s *BlockHeaderMut) SetStart(v int64) *BlockHeaderMut {
	s.start = v
	return s
}
func (s *BlockHeaderMut) SetEnd(v int64) *BlockHeaderMut {
	s.end = v
	return s
}
func (s *BlockHeaderMut) SetStorage(v uint64) *BlockHeaderMut {
	s.storage = v
	return s
}
func (s *BlockHeaderMut) SetStorageU(v uint64) *BlockHeaderMut {
	s.storageU = v
	return s
}
func (s *BlockHeaderMut) SetCount(v uint16) *BlockHeaderMut {
	s.count = v
	return s
}
func (s *BlockHeaderMut) SetSize(v uint16) *BlockHeaderMut {
	s.size = v
	return s
}
func (s *BlockHeaderMut) SetSizeU(v uint16) *BlockHeaderMut {
	s.sizeU = v
	return s
}
func (s *BlockHeaderMut) SetSizeX(v uint16) *BlockHeaderMut {
	s.sizeX = v
	return s
}
func (s *BlockHeaderMut) SetRecord(v uint16) *BlockHeaderMut {
	s.record = v
	return s
}
func (s *BlockHeaderMut) SetEncoding(v Compression) *BlockHeaderMut {
	s.encoding = v
	return s
}
func (s *BlockHeaderMut) SetKind(v StreamKind) *BlockHeaderMut {
	s.kind = v
	return s
}
func (s *BlockHeaderMut) SetSchema(v SchemaKind) *BlockHeaderMut {
	s.schema = v
	return s
}
func (s *BlockHeaderMut) SetRealTime(v bool) *BlockHeaderMut {
	s.realTime = v
	return s
}

type RecordHeader struct {
	blockID   BlockID
	id        int64
	prevID    int64
	timestamp int64
	start     int64
	end       int64
	seq       uint16
	sizeU     uint16
	size      uint16
	encoding  Compression
	pad       bool
}

func (s *RecordHeader) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *RecordHeader) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["blockID"] = s.BlockID().MarshalMap(nil)
	m["id"] = s.Id()
	m["prevID"] = s.PrevID()
	m["timestamp"] = s.Timestamp()
	m["start"] = s.Start()
	m["end"] = s.End()
	m["seq"] = s.Seq()
	m["sizeU"] = s.SizeU()
	m["size"] = s.Size()
	m["encoding"] = s.Encoding()
	m["pad"] = s.Pad()
	return m
}

func (s *RecordHeader) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[64]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 64 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *RecordHeader) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[64]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *RecordHeader) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[64]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *RecordHeader) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[64]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *RecordHeader) Read(b []byte) (n int, err error) {
	if len(b) < 64 {
		return -1, io.ErrShortBuffer
	}
	v := (*RecordHeader)(unsafe.Pointer(&b[0]))
	*v = *s
	return 64, nil
}
func (s *RecordHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 64 {
		return io.ErrShortBuffer
	}
	v := (*RecordHeader)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *RecordHeader) Clone() *RecordHeader {
	v := &RecordHeader{}
	*v = *s
	return v
}
func (s *RecordHeader) Bytes() []byte {
	return (*(*[64]byte)(unsafe.Pointer(s)))[0:]
}
func (s *RecordHeader) Mut() *RecordHeaderMut {
	return (*RecordHeaderMut)(unsafe.Pointer(s))
}
func (s *RecordHeader) BlockID() *BlockID {
	return &s.blockID
}
func (s *RecordHeader) Id() int64 {
	return s.id
}
func (s *RecordHeader) PrevID() int64 {
	return s.prevID
}
func (s *RecordHeader) Timestamp() int64 {
	return s.timestamp
}
func (s *RecordHeader) Start() int64 {
	return s.start
}
func (s *RecordHeader) End() int64 {
	return s.end
}
func (s *RecordHeader) Seq() uint16 {
	return s.seq
}
func (s *RecordHeader) SizeU() uint16 {
	return s.sizeU
}
func (s *RecordHeader) Size() uint16 {
	return s.size
}
func (s *RecordHeader) Encoding() Compression {
	return s.encoding
}
func (s *RecordHeader) Pad() bool {
	return s.pad
}

type RecordHeaderMut struct {
	RecordHeader
}

func (s *RecordHeaderMut) Clone() *RecordHeaderMut {
	v := &RecordHeaderMut{}
	*v = *s
	return v
}
func (s *RecordHeaderMut) Freeze() *RecordHeader {
	return (*RecordHeader)(unsafe.Pointer(s))
}
func (s *RecordHeaderMut) BlockID() *BlockIDMut {
	return s.blockID.Mut()
}
func (s *RecordHeaderMut) SetBlockID(v *BlockID) *RecordHeaderMut {
	s.blockID = *v
	return s
}
func (s *RecordHeaderMut) SetId(v int64) *RecordHeaderMut {
	s.id = v
	return s
}
func (s *RecordHeaderMut) SetPrevID(v int64) *RecordHeaderMut {
	s.prevID = v
	return s
}
func (s *RecordHeaderMut) SetTimestamp(v int64) *RecordHeaderMut {
	s.timestamp = v
	return s
}
func (s *RecordHeaderMut) SetStart(v int64) *RecordHeaderMut {
	s.start = v
	return s
}
func (s *RecordHeaderMut) SetEnd(v int64) *RecordHeaderMut {
	s.end = v
	return s
}
func (s *RecordHeaderMut) SetSeq(v uint16) *RecordHeaderMut {
	s.seq = v
	return s
}
func (s *RecordHeaderMut) SetSizeU(v uint16) *RecordHeaderMut {
	s.sizeU = v
	return s
}
func (s *RecordHeaderMut) SetSize(v uint16) *RecordHeaderMut {
	s.size = v
	return s
}
func (s *RecordHeaderMut) SetEncoding(v Compression) *RecordHeaderMut {
	s.encoding = v
	return s
}
func (s *RecordHeaderMut) SetPad(v bool) *RecordHeaderMut {
	s.pad = v
	return s
}

type SyncStarted struct {
	recordID  RecordID
	timestamp int64
}

func (s *SyncStarted) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *SyncStarted) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	return m
}

func (s *SyncStarted) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 32 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *SyncStarted) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[32]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *SyncStarted) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *SyncStarted) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *SyncStarted) Read(b []byte) (n int, err error) {
	if len(b) < 32 {
		return -1, io.ErrShortBuffer
	}
	v := (*SyncStarted)(unsafe.Pointer(&b[0]))
	*v = *s
	return 32, nil
}
func (s *SyncStarted) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*SyncStarted)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *SyncStarted) Clone() *SyncStarted {
	v := &SyncStarted{}
	*v = *s
	return v
}
func (s *SyncStarted) Bytes() []byte {
	return (*(*[32]byte)(unsafe.Pointer(s)))[0:]
}
func (s *SyncStarted) Mut() *SyncStartedMut {
	return (*SyncStartedMut)(unsafe.Pointer(s))
}
func (s *SyncStarted) RecordID() *RecordID {
	return &s.recordID
}
func (s *SyncStarted) Timestamp() int64 {
	return s.timestamp
}

type SyncStartedMut struct {
	SyncStarted
}

func (s *SyncStartedMut) Clone() *SyncStartedMut {
	v := &SyncStartedMut{}
	*v = *s
	return v
}
func (s *SyncStartedMut) Freeze() *SyncStarted {
	return (*SyncStarted)(unsafe.Pointer(s))
}
func (s *SyncStartedMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *SyncStartedMut) SetRecordID(v *RecordID) *SyncStartedMut {
	s.recordID = *v
	return s
}
func (s *SyncStartedMut) SetTimestamp(v int64) *SyncStartedMut {
	s.timestamp = v
	return s
}

type Savepoint struct {
	recordID  RecordID
	timestamp int64
	duration  int64
}

func (s *Savepoint) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Savepoint) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["recordID"] = s.RecordID().MarshalMap(nil)
	m["timestamp"] = s.Timestamp()
	m["duration"] = s.Duration()
	return m
}

func (s *Savepoint) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[40]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 40 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Savepoint) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[40]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Savepoint) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[40]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Savepoint) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[40]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Savepoint) Read(b []byte) (n int, err error) {
	if len(b) < 40 {
		return -1, io.ErrShortBuffer
	}
	v := (*Savepoint)(unsafe.Pointer(&b[0]))
	*v = *s
	return 40, nil
}
func (s *Savepoint) UnmarshalBinary(b []byte) error {
	if len(b) < 40 {
		return io.ErrShortBuffer
	}
	v := (*Savepoint)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Savepoint) Clone() *Savepoint {
	v := &Savepoint{}
	*v = *s
	return v
}
func (s *Savepoint) Bytes() []byte {
	return (*(*[40]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Savepoint) Mut() *SavepointMut {
	return (*SavepointMut)(unsafe.Pointer(s))
}
func (s *Savepoint) RecordID() *RecordID {
	return &s.recordID
}
func (s *Savepoint) Timestamp() int64 {
	return s.timestamp
}
func (s *Savepoint) Duration() int64 {
	return s.duration
}

type SavepointMut struct {
	Savepoint
}

func (s *SavepointMut) Clone() *SavepointMut {
	v := &SavepointMut{}
	*v = *s
	return v
}
func (s *SavepointMut) Freeze() *Savepoint {
	return (*Savepoint)(unsafe.Pointer(s))
}
func (s *SavepointMut) RecordID() *RecordIDMut {
	return s.recordID.Mut()
}
func (s *SavepointMut) SetRecordID(v *RecordID) *SavepointMut {
	s.recordID = *v
	return s
}
func (s *SavepointMut) SetTimestamp(v int64) *SavepointMut {
	s.timestamp = v
	return s
}
func (s *SavepointMut) SetDuration(v int64) *SavepointMut {
	s.duration = v
	return s
}

type String32 [32]byte

func NewString32(s string) *String32 {
	v := String32{}
	v.set(s)
	return &v
}
func (s *String32) set(v string) {
	copy(s[0:31], v)
	c := 31
	l := len(v)
	if l > c {
		s[31] = byte(c)
	} else {
		s[31] = byte(l)
	}
}
func (s *String32) Len() int {
	return int(s[31])
}
func (s *String32) Cap() int {
	return 31
}
func (s *String32) Unsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&s[0])),
		Len:  int(s[31]),
	}))
}
func (s *String32) String() string {
	return string(s[0:s[31]])
}
func (s *String32) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String32) Clone() *String32 {
	v := String32{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String32) Mut() *String32Mut {
	return *(**String32Mut)(unsafe.Pointer(&s))
}
func (s *String32) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[32]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 32 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String32) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[32]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String32) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[32]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String32) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[32]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *String32) UnmarshalBinary(b []byte) error {
	if len(b) < 32 {
		return io.ErrShortBuffer
	}
	v := (*String32)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String32Mut struct {
	String32
}

func (s *String32Mut) Set(v string) {
	s.set(v)
}

type Bytes65160 [65160]byte

func NewBytes65160(s string) *Bytes65160 {
	v := Bytes65160{}
	v.set(s)
	return &v
}
func (s *Bytes65160) set(v string) {
	copy(s[0:], v)
}
func (s *Bytes65160) Len() int {
	return 65160
}
func (s *Bytes65160) Cap() int {
	return 65160
}
func (s *Bytes65160) Unsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&s[0])),
		Len:  int(*(*uint16)(unsafe.Pointer(&s[65158]))),
	}))
}
func (s *Bytes65160) String() string {
	return string(s[0:s.Len()])
}
func (s *Bytes65160) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *Bytes65160) Clone() *Bytes65160 {
	v := Bytes65160{}
	copy(s[0:], v[0:])
	return &v
}
func (s *Bytes65160) Mut() *Bytes65160Mut {
	return *(**Bytes65160Mut)(unsafe.Pointer(&s))
}
func (s *Bytes65160) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[65160]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 65160 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *Bytes65160) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[65160]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *Bytes65160) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[65160]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *Bytes65160) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[65160]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *Bytes65160) UnmarshalBinary(b []byte) error {
	if len(b) < 65160 {
		return io.ErrShortBuffer
	}
	v := (*Bytes65160)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type Bytes65160Mut struct {
	Bytes65160
}

func (s *Bytes65160Mut) Set(v string) {
	s.set(v)
}

type String64 [64]byte

func NewString64(s string) *String64 {
	v := String64{}
	v.set(s)
	return &v
}
func (s *String64) set(v string) {
	copy(s[0:63], v)
	c := 63
	l := len(v)
	if l > c {
		s[63] = byte(c)
	} else {
		s[63] = byte(l)
	}
}
func (s *String64) Len() int {
	return int(s[63])
}
func (s *String64) Cap() int {
	return 63
}
func (s *String64) Unsafe() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&s[0])),
		Len:  int(s[63]),
	}))
}
func (s *String64) String() string {
	return string(s[0:s[63]])
}
func (s *String64) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String64) Clone() *String64 {
	v := String64{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String64) Mut() *String64Mut {
	return *(**String64Mut)(unsafe.Pointer(&s))
}
func (s *String64) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[64]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 64 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String64) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[64]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String64) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[64]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String64) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[64]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *String64) UnmarshalBinary(b []byte) error {
	if len(b) < 64 {
		return io.ErrShortBuffer
	}
	v := (*String64)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String64Mut struct {
	String64
}

func (s *String64Mut) Set(v string) {
	s.set(v)
}
func init() {
	{
		var b [2]byte
		v := uint16(1)
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		if *(*uint16)(unsafe.Pointer(&b[0])) != 1 {
			panic("BigEndian detected... compiled for LittleEndian only!!!")
		}
	}
	to := reflect.TypeOf
	type sf struct {
		n string
		o uintptr
		s uintptr
	}
	ss := func(tt interface{}, mtt interface{}, s uintptr, fl []sf) {
		t := to(tt)
		mt := to(mtt)
		if t.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", t.Name(), t.Size(), s))
		}
		if mt.Size() != s {
			panic(fmt.Sprintf("sizeof %s = %d, expected = %d", mt.Name(), mt.Size(), s))
		}
		if t.NumField() != len(fl) {
			panic(fmt.Sprintf("%s field count = %d: expected %d", t.Name(), t.NumField(), len(fl)))
		}
		for i, ef := range fl {
			f := t.Field(i)
			if f.Offset != ef.o {
				panic(fmt.Sprintf("%s.%s offset = %d, expected = %d", t.Name(), f.Name, f.Offset, ef.o))
			}
			if f.Type.Size() != ef.s {
				panic(fmt.Sprintf("%s.%s size = %d, expected = %d", t.Name(), f.Name, f.Type.Size(), ef.s))
			}
			if f.Name != ef.n {
				panic(fmt.Sprintf("%s.%s expected field: %s", t.Name(), f.Name, ef.n))
			}
		}
	}

	ss(RecordID{}, RecordIDMut{}, 24, []sf{
		{"streamID", 0, 8},
		{"blockID", 8, 8},
		{"id", 16, 8},
	})
	ss(Block{}, BlockMut{}, 65256, []sf{
		{"header", 0, 96},
		{"data", 96, 65160},
	})
	ss(EOS{}, EOSMut{}, 32, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
	})
	ss(Stopped{}, StoppedMut{}, 48, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"reason", 32, 1},
		{"_", 33, 7},
		{"starts", 40, 8},
	})
	ss(RecordsHeader{}, RecordsHeaderMut{}, 72, []sf{
		{"header", 0, 64},
		{"count", 64, 2},
		{"record", 66, 2},
		{"_", 68, 4},
	})
	ss(StreamStats{}, StreamStatsMut{}, 48, []sf{
		{"storage", 0, 24},
		{"appender", 24, 24},
	})
	ss(SyncProgress{}, SyncProgressMut{}, 56, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"started", 32, 8},
		{"count", 40, 8},
		{"remaining", 48, 8},
	})
	ss(Stream{}, StreamMut{}, 80, []sf{
		{"id", 0, 8},
		{"created", 8, 8},
		{"accountID", 16, 8},
		{"duration", 24, 8},
		{"record", 32, 4},
		{"_", 36, 4},
		{"name", 40, 32},
		{"kind", 72, 1},
		{"schema", 73, 1},
		{"realTime", 74, 1},
		{"_", 75, 5},
	})
	ss(EOSWaiting{}, EOSWaitingMut{}, 32, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
	})
	ss(Stats{}, StatsMut{}, 24, []sf{
		{"size", 0, 8},
		{"count", 8, 8},
		{"blocks", 16, 8},
	})
	ss(AccountStats{}, AccountStatsMut{}, 64, []sf{
		{"id", 0, 8},
		{"storage", 8, 24},
		{"appender", 32, 24},
		{"streams", 56, 8},
	})
	ss(BlockID{}, BlockIDMut{}, 16, []sf{
		{"streamID", 0, 8},
		{"id", 8, 8},
	})
	ss(Started{}, StartedMut{}, 48, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"writerID", 32, 8},
		{"stops", 40, 8},
	})
	ss(SyncStopped{}, SyncStoppedMut{}, 128, []sf{
		{"progress", 0, 56},
		{"reason", 56, 1},
		{"_", 57, 7},
		{"message", 64, 64},
	})
	ss(Starting{}, StartingMut{}, 40, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"writerID", 32, 8},
	})
	ss(BlockHeader{}, BlockHeaderMut{}, 96, []sf{
		{"id", 0, 16},
		{"created", 16, 8},
		{"completed", 24, 8},
		{"min", 32, 8},
		{"max", 40, 8},
		{"start", 48, 8},
		{"end", 56, 8},
		{"storage", 64, 8},
		{"storageU", 72, 8},
		{"count", 80, 2},
		{"size", 82, 2},
		{"sizeU", 84, 2},
		{"sizeX", 86, 2},
		{"record", 88, 2},
		{"encoding", 90, 1},
		{"kind", 91, 1},
		{"schema", 92, 1},
		{"realTime", 93, 1},
		{"_", 94, 2},
	})
	ss(RecordHeader{}, RecordHeaderMut{}, 64, []sf{
		{"blockID", 0, 16},
		{"id", 16, 8},
		{"prevID", 24, 8},
		{"timestamp", 32, 8},
		{"start", 40, 8},
		{"end", 48, 8},
		{"seq", 56, 2},
		{"sizeU", 58, 2},
		{"size", 60, 2},
		{"encoding", 62, 1},
		{"pad", 63, 1},
	})
	ss(SyncStarted{}, SyncStartedMut{}, 32, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
	})
	ss(Savepoint{}, SavepointMut{}, 40, []sf{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"duration", 32, 8},
	})

}
