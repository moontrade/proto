//go:build 386 || amd64 || arm || arm64 || ppc64le || mips64le || mipsle || riscv64 || wasm
// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm

package schema2

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type RecordLayout byte

const (
	RecordLayout_Aligned = RecordLayout(0)
	RecordLayout_Compact = RecordLayout(1)
)

type BlockLayout byte

const (
	BlockLayout_Row    = BlockLayout(1)
	BlockLayout_Column = BlockLayout(2)
)

type Format byte

const (
	Format_Raw      = Format(0)
	Format_WAP      = Format(1)
	Format_Json     = Format(2)
	Format_Protobuf = Format(3)
)

type StreamKind byte

const (
	StreamKind_Log    = StreamKind(0)
	StreamKind_Series = StreamKind(1)
	StreamKind_Table  = StreamKind(2)
)

type Kind byte

const (
	Kind_Unknown      = Kind(0)
	Kind_Bool         = Kind(1)
	Kind_Byte         = Kind(2)
	Kind_Int8         = Kind(3)
	Kind_UInt8        = Kind(4)
	Kind_Int16        = Kind(5)
	Kind_UInt16       = Kind(6)
	Kind_Int32        = Kind(7)
	Kind_UInt32       = Kind(8)
	Kind_Int64        = Kind(9)
	Kind_UInt64       = Kind(10)
	Kind_Float32      = Kind(11)
	Kind_Float64      = Kind(12)
	Kind_String       = Kind(13)
	Kind_Bytes        = Kind(14)
	Kind_RecordHeader = Kind(20)
	Kind_BlockHeader  = Kind(21)
	Kind_Enum         = Kind(30)
	Kind_Record       = Kind(40)
	Kind_Struct       = Kind(41)
	Kind_List         = Kind(50)
	Kind_LinkedList   = Kind(51)
	Kind_Map          = Kind(52)
	Kind_LinkedMap    = Kind(53)
)

type SyncStoppedReason byte

const (
	SyncStoppedReason_Success = SyncStoppedReason(1)
	SyncStoppedReason_Error   = SyncStoppedReason(2)
)

type BlockSize uint16

const (
	BlockSize_B1kb  = BlockSize(1024)
	BlockSize_B2kb  = BlockSize(2048)
	BlockSize_B4kb  = BlockSize(4096)
	BlockSize_B8kb  = BlockSize(8192)
	BlockSize_B16kb = BlockSize(16384)
	BlockSize_B32kb = BlockSize(32768)
	BlockSize_B64kb = BlockSize(65535)
)

type MessageType byte

const (
	MessageType_Record       = MessageType(1)
	MessageType_Records      = MessageType(2)
	MessageType_Block        = MessageType(3)
	MessageType_EOS          = MessageType(4)
	MessageType_EOSWaiting   = MessageType(5)
	MessageType_Savepoint    = MessageType(6)
	MessageType_Starting     = MessageType(7)
	MessageType_Started      = MessageType(8)
	MessageType_Stopped      = MessageType(9)
	MessageType_SyncStarted  = MessageType(10)
	MessageType_SyncProgress = MessageType(11)
	MessageType_SyncStopped  = MessageType(12)
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

type Encoding byte

const (
	Encoding_None   = Encoding(0)
	Encoding_LZ4    = Encoding(1)
	Encoding_ZSTD   = Encoding(2)
	Encoding_Brotli = Encoding(3)
	Encoding_Gzip   = Encoding(4)
)

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

type Stream struct {
	id        int64
	created   int64
	accountID int64
	duration  int64
	record    uint16
	_         [6]byte // Padding
	name      String32
	kind      StreamKind
	format    Format
	blockSize BlockSize
	_         [4]byte // Padding
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
	m["format"] = s.Format()
	m["blockSize"] = s.BlockSize()
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
func (s *Stream) Record() uint16 {
	return s.record
}
func (s *Stream) Name() *String32 {
	return &s.name
}
func (s *Stream) Kind() StreamKind {
	return s.kind
}
func (s *Stream) Format() Format {
	return s.format
}
func (s *Stream) BlockSize() BlockSize {
	return s.blockSize
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
func (s *StreamMut) SetRecord(v uint16) *StreamMut {
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
func (s *StreamMut) SetFormat(v Format) *StreamMut {
	s.format = v
	return s
}
func (s *StreamMut) SetBlockSize(v BlockSize) *StreamMut {
	s.blockSize = v
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

type UnionOption struct {
	name String40
	kind Kind
	_    [7]byte // Padding
	id   String40
}

func (s *UnionOption) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *UnionOption) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["name"] = s.Name()
	m["kind"] = s.Kind()
	m["id"] = s.Id()
	return m
}

func (s *UnionOption) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[88]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 88 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *UnionOption) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[88]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *UnionOption) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[88]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *UnionOption) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[88]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *UnionOption) Read(b []byte) (n int, err error) {
	if len(b) < 88 {
		return -1, io.ErrShortBuffer
	}
	v := (*UnionOption)(unsafe.Pointer(&b[0]))
	*v = *s
	return 88, nil
}
func (s *UnionOption) UnmarshalBinary(b []byte) error {
	if len(b) < 88 {
		return io.ErrShortBuffer
	}
	v := (*UnionOption)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *UnionOption) Clone() *UnionOption {
	v := &UnionOption{}
	*v = *s
	return v
}
func (s *UnionOption) Bytes() []byte {
	return (*(*[88]byte)(unsafe.Pointer(s)))[0:]
}
func (s *UnionOption) Mut() *UnionOptionMut {
	return (*UnionOptionMut)(unsafe.Pointer(s))
}
func (s *UnionOption) Name() *String40 {
	return &s.name
}
func (s *UnionOption) Kind() Kind {
	return s.kind
}
func (s *UnionOption) Id() *String40 {
	return &s.id
}

type UnionOptionMut struct {
	UnionOption
}

func (s *UnionOptionMut) Clone() *UnionOptionMut {
	v := &UnionOptionMut{}
	*v = *s
	return v
}
func (s *UnionOptionMut) Freeze() *UnionOption {
	return (*UnionOption)(unsafe.Pointer(s))
}
func (s *UnionOptionMut) Name() *String40Mut {
	return s.name.Mut()
}
func (s *UnionOptionMut) SetName(v *String40) *UnionOptionMut {
	s.name = *v
	return s
}
func (s *UnionOptionMut) SetKind(v Kind) *UnionOptionMut {
	s.kind = v
	return s
}
func (s *UnionOptionMut) Id() *String40Mut {
	return s.id.Mut()
}
func (s *UnionOptionMut) SetId(v *String40) *UnionOptionMut {
	s.id = *v
	return s
}

type Line struct {
	number int32
	begin  int32
	end    int32
	_      [4]byte // Padding
}

func (s *Line) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Line) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["number"] = s.Number()
	m["begin"] = s.Begin()
	m["end"] = s.End()
	return m
}

func (s *Line) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[16]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 16 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Line) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[16]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Line) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Line) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[16]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Line) Read(b []byte) (n int, err error) {
	if len(b) < 16 {
		return -1, io.ErrShortBuffer
	}
	v := (*Line)(unsafe.Pointer(&b[0]))
	*v = *s
	return 16, nil
}
func (s *Line) UnmarshalBinary(b []byte) error {
	if len(b) < 16 {
		return io.ErrShortBuffer
	}
	v := (*Line)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Line) Clone() *Line {
	v := &Line{}
	*v = *s
	return v
}
func (s *Line) Bytes() []byte {
	return (*(*[16]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Line) Mut() *LineMut {
	return (*LineMut)(unsafe.Pointer(s))
}
func (s *Line) Number() int32 {
	return s.number
}
func (s *Line) Begin() int32 {
	return s.begin
}
func (s *Line) End() int32 {
	return s.end
}

type LineMut struct {
	Line
}

func (s *LineMut) Clone() *LineMut {
	v := &LineMut{}
	*v = *s
	return v
}
func (s *LineMut) Freeze() *Line {
	return (*Line)(unsafe.Pointer(s))
}
func (s *LineMut) SetNumber(v int32) *LineMut {
	s.number = v
	return s
}
func (s *LineMut) SetBegin(v int32) *LineMut {
	s.begin = v
	return s
}
func (s *LineMut) SetEnd(v int32) *LineMut {
	s.end = v
	return s
}

type Import struct {
	id    int32
	_     [4]byte // Padding
	line  Line
	path  String128
	name  String32
	alias String32
}

func (s *Import) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Import) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id()
	m["line"] = s.Line().MarshalMap(nil)
	m["path"] = s.Path()
	m["name"] = s.Name()
	m["alias"] = s.Alias()
	return m
}

func (s *Import) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[216]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 216 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Import) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[216]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Import) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[216]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Import) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[216]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Import) Read(b []byte) (n int, err error) {
	if len(b) < 216 {
		return -1, io.ErrShortBuffer
	}
	v := (*Import)(unsafe.Pointer(&b[0]))
	*v = *s
	return 216, nil
}
func (s *Import) UnmarshalBinary(b []byte) error {
	if len(b) < 216 {
		return io.ErrShortBuffer
	}
	v := (*Import)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Import) Clone() *Import {
	v := &Import{}
	*v = *s
	return v
}
func (s *Import) Bytes() []byte {
	return (*(*[216]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Import) Mut() *ImportMut {
	return (*ImportMut)(unsafe.Pointer(s))
}
func (s *Import) Id() int32 {
	return s.id
}
func (s *Import) Line() *Line {
	return &s.line
}
func (s *Import) Path() *String128 {
	return &s.path
}
func (s *Import) Name() *String32 {
	return &s.name
}
func (s *Import) Alias() *String32 {
	return &s.alias
}

type ImportMut struct {
	Import
}

func (s *ImportMut) Clone() *ImportMut {
	v := &ImportMut{}
	*v = *s
	return v
}
func (s *ImportMut) Freeze() *Import {
	return (*Import)(unsafe.Pointer(s))
}
func (s *ImportMut) SetId(v int32) *ImportMut {
	s.id = v
	return s
}
func (s *ImportMut) Line() *LineMut {
	return s.line.Mut()
}
func (s *ImportMut) SetLine(v *Line) *ImportMut {
	s.line = *v
	return s
}
func (s *ImportMut) Path() *String128Mut {
	return s.path.Mut()
}
func (s *ImportMut) SetPath(v *String128) *ImportMut {
	s.path = *v
	return s
}
func (s *ImportMut) Name() *String32Mut {
	return s.name.Mut()
}
func (s *ImportMut) SetName(v *String32) *ImportMut {
	s.name = *v
	return s
}
func (s *ImportMut) Alias() *String32Mut {
	return s.alias.Mut()
}
func (s *ImportMut) SetAlias(v *String32) *ImportMut {
	s.alias = *v
	return s
}

type Imports struct {
	id   int32
	_    [4]byte // Padding
	line Line
	list Import16List
}

func (s *Imports) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Imports) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id()
	m["line"] = s.Line().MarshalMap(nil)
	m["list"] = s.List().CopyTo(nil)
	return m
}

func (s *Imports) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[3488]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 3488 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Imports) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[3488]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Imports) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[3488]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Imports) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[3488]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Imports) Read(b []byte) (n int, err error) {
	if len(b) < 3488 {
		return -1, io.ErrShortBuffer
	}
	v := (*Imports)(unsafe.Pointer(&b[0]))
	*v = *s
	return 3488, nil
}
func (s *Imports) UnmarshalBinary(b []byte) error {
	if len(b) < 3488 {
		return io.ErrShortBuffer
	}
	v := (*Imports)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Imports) Clone() *Imports {
	v := &Imports{}
	*v = *s
	return v
}
func (s *Imports) Bytes() []byte {
	return (*(*[3488]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Imports) Mut() *ImportsMut {
	return (*ImportsMut)(unsafe.Pointer(s))
}
func (s *Imports) Id() int32 {
	return s.id
}
func (s *Imports) Line() *Line {
	return &s.line
}
func (s *Imports) List() *Import16List {
	return &s.list
}

type ImportsMut struct {
	Imports
}

func (s *ImportsMut) Clone() *ImportsMut {
	v := &ImportsMut{}
	*v = *s
	return v
}
func (s *ImportsMut) Freeze() *Imports {
	return (*Imports)(unsafe.Pointer(s))
}
func (s *ImportsMut) SetId(v int32) *ImportsMut {
	s.id = v
	return s
}
func (s *ImportsMut) Line() *LineMut {
	return s.line.Mut()
}
func (s *ImportsMut) SetLine(v *Line) *ImportsMut {
	s.line = *v
	return s
}
func (s *ImportsMut) List() *Import16ListMut {
	return s.list.Mut()
}
func (s *ImportsMut) SetList(v *Import16List) *ImportsMut {
	s.list = *v
	return s
}

// BlockHeader
type BlockHeader struct {
	streamID  int64
	id        int64
	headID    int64
	headMin   int64
	headStart int64
	blocks    int64
	records   int64
	storage   uint64
	storageU  uint64
	created   int64
	completed int64
	start     int64
	end       int64
	min       int64
	max       int64
	count     uint16
	size      uint16
	sizeU     uint16
	sizeX     uint16
	record    uint16
	blockSize BlockSize
	encoding  Encoding
	kind      StreamKind
	format    Format
	_         [1]byte // Padding
}

func (s *BlockHeader) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *BlockHeader) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["streamID"] = s.StreamID()
	m["id"] = s.Id()
	m["headID"] = s.HeadID()
	m["headMin"] = s.HeadMin()
	m["headStart"] = s.HeadStart()
	m["blocks"] = s.Blocks()
	m["records"] = s.Records()
	m["storage"] = s.Storage()
	m["storageU"] = s.StorageU()
	m["created"] = s.Created()
	m["completed"] = s.Completed()
	m["start"] = s.Start()
	m["end"] = s.End()
	m["min"] = s.Min()
	m["max"] = s.Max()
	m["count"] = s.Count()
	m["size"] = s.Size()
	m["sizeU"] = s.SizeU()
	m["sizeX"] = s.SizeX()
	m["record"] = s.Record()
	m["blockSize"] = s.BlockSize()
	m["encoding"] = s.Encoding()
	m["kind"] = s.Kind()
	m["format"] = s.Format()
	return m
}

func (s *BlockHeader) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[136]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 136 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *BlockHeader) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[136]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *BlockHeader) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[136]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *BlockHeader) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[136]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *BlockHeader) Read(b []byte) (n int, err error) {
	if len(b) < 136 {
		return -1, io.ErrShortBuffer
	}
	v := (*BlockHeader)(unsafe.Pointer(&b[0]))
	*v = *s
	return 136, nil
}
func (s *BlockHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 136 {
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
	return (*(*[136]byte)(unsafe.Pointer(s)))[0:]
}
func (s *BlockHeader) Mut() *BlockHeaderMut {
	return (*BlockHeaderMut)(unsafe.Pointer(s))
}
func (s *BlockHeader) StreamID() int64 {
	return s.streamID
}
func (s *BlockHeader) Id() int64 {
	return s.id
}
func (s *BlockHeader) HeadID() int64 {
	return s.headID
}
func (s *BlockHeader) HeadMin() int64 {
	return s.headMin
}
func (s *BlockHeader) HeadStart() int64 {
	return s.headStart
}
func (s *BlockHeader) Blocks() int64 {
	return s.blocks
}
func (s *BlockHeader) Records() int64 {
	return s.records
}
func (s *BlockHeader) Storage() uint64 {
	return s.storage
}
func (s *BlockHeader) StorageU() uint64 {
	return s.storageU
}
func (s *BlockHeader) Created() int64 {
	return s.created
}
func (s *BlockHeader) Completed() int64 {
	return s.completed
}
func (s *BlockHeader) Start() int64 {
	return s.start
}
func (s *BlockHeader) End() int64 {
	return s.end
}
func (s *BlockHeader) Min() int64 {
	return s.min
}
func (s *BlockHeader) Max() int64 {
	return s.max
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
func (s *BlockHeader) BlockSize() BlockSize {
	return s.blockSize
}
func (s *BlockHeader) Encoding() Encoding {
	return s.encoding
}
func (s *BlockHeader) Kind() StreamKind {
	return s.kind
}
func (s *BlockHeader) Format() Format {
	return s.format
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
func (s *BlockHeaderMut) SetStreamID(v int64) *BlockHeaderMut {
	s.streamID = v
	return s
}
func (s *BlockHeaderMut) SetId(v int64) *BlockHeaderMut {
	s.id = v
	return s
}
func (s *BlockHeaderMut) SetHeadID(v int64) *BlockHeaderMut {
	s.headID = v
	return s
}
func (s *BlockHeaderMut) SetHeadMin(v int64) *BlockHeaderMut {
	s.headMin = v
	return s
}
func (s *BlockHeaderMut) SetHeadStart(v int64) *BlockHeaderMut {
	s.headStart = v
	return s
}
func (s *BlockHeaderMut) SetBlocks(v int64) *BlockHeaderMut {
	s.blocks = v
	return s
}
func (s *BlockHeaderMut) SetRecords(v int64) *BlockHeaderMut {
	s.records = v
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
func (s *BlockHeaderMut) SetCreated(v int64) *BlockHeaderMut {
	s.created = v
	return s
}
func (s *BlockHeaderMut) SetCompleted(v int64) *BlockHeaderMut {
	s.completed = v
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
func (s *BlockHeaderMut) SetMin(v int64) *BlockHeaderMut {
	s.min = v
	return s
}
func (s *BlockHeaderMut) SetMax(v int64) *BlockHeaderMut {
	s.max = v
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
func (s *BlockHeaderMut) SetBlockSize(v BlockSize) *BlockHeaderMut {
	s.blockSize = v
	return s
}
func (s *BlockHeaderMut) SetEncoding(v Encoding) *BlockHeaderMut {
	s.encoding = v
	return s
}
func (s *BlockHeaderMut) SetKind(v StreamKind) *BlockHeaderMut {
	s.kind = v
	return s
}
func (s *BlockHeaderMut) SetFormat(v Format) *BlockHeaderMut {
	s.format = v
	return s
}

type Record struct {
	name   String40
	fields Field64List
}

func (s *Record) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Record) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["name"] = s.Name()
	m["fields"] = s.Fields().CopyTo(nil)
	return m
}

func (s *Record) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[4144]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 4144 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Record) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[4144]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Record) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[4144]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Record) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[4144]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Record) Read(b []byte) (n int, err error) {
	if len(b) < 4144 {
		return -1, io.ErrShortBuffer
	}
	v := (*Record)(unsafe.Pointer(&b[0]))
	*v = *s
	return 4144, nil
}
func (s *Record) UnmarshalBinary(b []byte) error {
	if len(b) < 4144 {
		return io.ErrShortBuffer
	}
	v := (*Record)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Record) Clone() *Record {
	v := &Record{}
	*v = *s
	return v
}
func (s *Record) Bytes() []byte {
	return (*(*[4144]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Record) Mut() *RecordMut {
	return (*RecordMut)(unsafe.Pointer(s))
}
func (s *Record) Name() *String40 {
	return &s.name
}
func (s *Record) Fields() *Field64List {
	return &s.fields
}

type RecordMut struct {
	Record
}

func (s *RecordMut) Clone() *RecordMut {
	v := &RecordMut{}
	*v = *s
	return v
}
func (s *RecordMut) Freeze() *Record {
	return (*Record)(unsafe.Pointer(s))
}
func (s *RecordMut) Name() *String40Mut {
	return s.name.Mut()
}
func (s *RecordMut) SetName(v *String40) *RecordMut {
	s.name = *v
	return s
}
func (s *RecordMut) Fields() *Field64ListMut {
	return s.fields.Mut()
}
func (s *RecordMut) SetFields(v *Field64List) *RecordMut {
	s.fields = *v
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

type Enum struct {
	name    String40
	options EnumOption16List
}

func (s *Enum) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Enum) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["name"] = s.Name()
	m["options"] = s.Options().CopyTo(nil)
	return m
}

func (s *Enum) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[2352]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 2352 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Enum) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[2352]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Enum) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[2352]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Enum) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[2352]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Enum) Read(b []byte) (n int, err error) {
	if len(b) < 2352 {
		return -1, io.ErrShortBuffer
	}
	v := (*Enum)(unsafe.Pointer(&b[0]))
	*v = *s
	return 2352, nil
}
func (s *Enum) UnmarshalBinary(b []byte) error {
	if len(b) < 2352 {
		return io.ErrShortBuffer
	}
	v := (*Enum)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Enum) Clone() *Enum {
	v := &Enum{}
	*v = *s
	return v
}
func (s *Enum) Bytes() []byte {
	return (*(*[2352]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Enum) Mut() *EnumMut {
	return (*EnumMut)(unsafe.Pointer(s))
}
func (s *Enum) Name() *String40 {
	return &s.name
}
func (s *Enum) Options() *EnumOption16List {
	return &s.options
}

type EnumMut struct {
	Enum
}

func (s *EnumMut) Clone() *EnumMut {
	v := &EnumMut{}
	*v = *s
	return v
}
func (s *EnumMut) Freeze() *Enum {
	return (*Enum)(unsafe.Pointer(s))
}
func (s *EnumMut) Name() *String40Mut {
	return s.name.Mut()
}
func (s *EnumMut) SetName(v *String40) *EnumMut {
	s.name = *v
	return s
}
func (s *EnumMut) Options() *EnumOption16ListMut {
	return s.options.Mut()
}
func (s *EnumMut) SetOptions(v *EnumOption16List) *EnumMut {
	s.options = *v
	return s
}

type RecordHeader struct {
	id        RecordID
	prevID    int64
	timestamp int64
	start     int64
	end       int64
	seq       uint16
	size      uint16
	sizeU     uint16
	encoding  Encoding
	pad       bool
}

func (s *RecordHeader) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *RecordHeader) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["id"] = s.Id().MarshalMap(nil)
	m["prevID"] = s.PrevID()
	m["timestamp"] = s.Timestamp()
	m["start"] = s.Start()
	m["end"] = s.End()
	m["seq"] = s.Seq()
	m["size"] = s.Size()
	m["sizeU"] = s.SizeU()
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
func (s *RecordHeader) Id() *RecordID {
	return &s.id
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
func (s *RecordHeader) Size() uint16 {
	return s.size
}
func (s *RecordHeader) SizeU() uint16 {
	return s.sizeU
}
func (s *RecordHeader) Encoding() Encoding {
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
func (s *RecordHeaderMut) Id() *RecordIDMut {
	return s.id.Mut()
}
func (s *RecordHeaderMut) SetId(v *RecordID) *RecordHeaderMut {
	s.id = *v
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
func (s *RecordHeaderMut) SetSize(v uint16) *RecordHeaderMut {
	s.size = v
	return s
}
func (s *RecordHeaderMut) SetSizeU(v uint16) *RecordHeaderMut {
	s.sizeU = v
	return s
}
func (s *RecordHeaderMut) SetEncoding(v Encoding) *RecordHeaderMut {
	s.encoding = v
	return s
}
func (s *RecordHeaderMut) SetPad(v bool) *RecordHeaderMut {
	s.pad = v
	return s
}

type Struct struct {
	name   String40
	fields Field64List
}

func (s *Struct) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Struct) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["name"] = s.Name()
	m["fields"] = s.Fields().CopyTo(nil)
	return m
}

func (s *Struct) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[4144]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 4144 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Struct) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[4144]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Struct) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[4144]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Struct) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[4144]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Struct) Read(b []byte) (n int, err error) {
	if len(b) < 4144 {
		return -1, io.ErrShortBuffer
	}
	v := (*Struct)(unsafe.Pointer(&b[0]))
	*v = *s
	return 4144, nil
}
func (s *Struct) UnmarshalBinary(b []byte) error {
	if len(b) < 4144 {
		return io.ErrShortBuffer
	}
	v := (*Struct)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Struct) Clone() *Struct {
	v := &Struct{}
	*v = *s
	return v
}
func (s *Struct) Bytes() []byte {
	return (*(*[4144]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Struct) Mut() *StructMut {
	return (*StructMut)(unsafe.Pointer(s))
}
func (s *Struct) Name() *String40 {
	return &s.name
}
func (s *Struct) Fields() *Field64List {
	return &s.fields
}

type StructMut struct {
	Struct
}

func (s *StructMut) Clone() *StructMut {
	v := &StructMut{}
	*v = *s
	return v
}
func (s *StructMut) Freeze() *Struct {
	return (*Struct)(unsafe.Pointer(s))
}
func (s *StructMut) Name() *String40Mut {
	return s.name.Mut()
}
func (s *StructMut) SetName(v *String40) *StructMut {
	s.name = *v
	return s
}
func (s *StructMut) Fields() *Field64ListMut {
	return s.fields.Mut()
}
func (s *StructMut) SetFields(v *Field64List) *StructMut {
	s.fields = *v
	return s
}

type Union struct {
	name    String40
	options UnionOption16List
}

func (s *Union) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Union) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["name"] = s.Name()
	m["options"] = s.Options().CopyTo(nil)
	return m
}

func (s *Union) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[1456]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 1456 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Union) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[1456]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Union) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[1456]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Union) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[1456]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Union) Read(b []byte) (n int, err error) {
	if len(b) < 1456 {
		return -1, io.ErrShortBuffer
	}
	v := (*Union)(unsafe.Pointer(&b[0]))
	*v = *s
	return 1456, nil
}
func (s *Union) UnmarshalBinary(b []byte) error {
	if len(b) < 1456 {
		return io.ErrShortBuffer
	}
	v := (*Union)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Union) Clone() *Union {
	v := &Union{}
	*v = *s
	return v
}
func (s *Union) Bytes() []byte {
	return (*(*[1456]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Union) Mut() *UnionMut {
	return (*UnionMut)(unsafe.Pointer(s))
}
func (s *Union) Name() *String40 {
	return &s.name
}
func (s *Union) Options() *UnionOption16List {
	return &s.options
}

type UnionMut struct {
	Union
}

func (s *UnionMut) Clone() *UnionMut {
	v := &UnionMut{}
	*v = *s
	return v
}
func (s *UnionMut) Freeze() *Union {
	return (*Union)(unsafe.Pointer(s))
}
func (s *UnionMut) Name() *String40Mut {
	return s.name.Mut()
}
func (s *UnionMut) SetName(v *String40) *UnionMut {
	s.name = *v
	return s
}
func (s *UnionMut) Options() *UnionOption16ListMut {
	return s.options.Mut()
}
func (s *UnionMut) SetOptions(v *UnionOption16List) *UnionMut {
	s.options = *v
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

type SyncStopped struct {
	progress SyncProgress
	message  String64
	reason   SyncStoppedReason
	_        [7]byte // Padding
}

func (s *SyncStopped) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *SyncStopped) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["progress"] = s.Progress().MarshalMap(nil)
	m["message"] = s.Message()
	m["reason"] = s.Reason()
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
func (s *SyncStopped) Message() *String64 {
	return &s.message
}
func (s *SyncStopped) Reason() SyncStoppedReason {
	return s.reason
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
func (s *SyncStoppedMut) Message() *String64Mut {
	return s.message.Mut()
}
func (s *SyncStoppedMut) SetMessage(v *String64) *SyncStoppedMut {
	s.message = *v
	return s
}
func (s *SyncStoppedMut) SetReason(v SyncStoppedReason) *SyncStoppedMut {
	s.reason = v
	return s
}

type EnumOption struct {
	index    int32
	_        [4]byte // Padding
	name     String64
	value    int64
	valueStr String64
}

func (s *EnumOption) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *EnumOption) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["index"] = s.Index()
	m["name"] = s.Name()
	m["value"] = s.Value()
	m["valueStr"] = s.ValueStr()
	return m
}

func (s *EnumOption) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[144]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 144 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *EnumOption) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[144]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *EnumOption) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[144]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *EnumOption) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[144]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *EnumOption) Read(b []byte) (n int, err error) {
	if len(b) < 144 {
		return -1, io.ErrShortBuffer
	}
	v := (*EnumOption)(unsafe.Pointer(&b[0]))
	*v = *s
	return 144, nil
}
func (s *EnumOption) UnmarshalBinary(b []byte) error {
	if len(b) < 144 {
		return io.ErrShortBuffer
	}
	v := (*EnumOption)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *EnumOption) Clone() *EnumOption {
	v := &EnumOption{}
	*v = *s
	return v
}
func (s *EnumOption) Bytes() []byte {
	return (*(*[144]byte)(unsafe.Pointer(s)))[0:]
}
func (s *EnumOption) Mut() *EnumOptionMut {
	return (*EnumOptionMut)(unsafe.Pointer(s))
}
func (s *EnumOption) Index() int32 {
	return s.index
}
func (s *EnumOption) Name() *String64 {
	return &s.name
}
func (s *EnumOption) Value() int64 {
	return s.value
}
func (s *EnumOption) ValueStr() *String64 {
	return &s.valueStr
}

type EnumOptionMut struct {
	EnumOption
}

func (s *EnumOptionMut) Clone() *EnumOptionMut {
	v := &EnumOptionMut{}
	*v = *s
	return v
}
func (s *EnumOptionMut) Freeze() *EnumOption {
	return (*EnumOption)(unsafe.Pointer(s))
}
func (s *EnumOptionMut) SetIndex(v int32) *EnumOptionMut {
	s.index = v
	return s
}
func (s *EnumOptionMut) Name() *String64Mut {
	return s.name.Mut()
}
func (s *EnumOptionMut) SetName(v *String64) *EnumOptionMut {
	s.name = *v
	return s
}
func (s *EnumOptionMut) SetValue(v int64) *EnumOptionMut {
	s.value = v
	return s
}
func (s *EnumOptionMut) ValueStr() *String64Mut {
	return s.valueStr.Mut()
}
func (s *EnumOptionMut) SetValueStr(v *String64) *EnumOptionMut {
	s.valueStr = *v
	return s
}

type Schema struct {
	imports Imports16List
	records Record128List
}

func (s *Schema) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Schema) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["imports"] = s.Imports().CopyTo(nil)
	m["records"] = s.Records().CopyTo(nil)
	return m
}

func (s *Schema) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[586256]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 586256 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Schema) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[586256]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Schema) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[586256]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Schema) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[586256]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Schema) Read(b []byte) (n int, err error) {
	if len(b) < 586256 {
		return -1, io.ErrShortBuffer
	}
	v := (*Schema)(unsafe.Pointer(&b[0]))
	*v = *s
	return 586256, nil
}
func (s *Schema) UnmarshalBinary(b []byte) error {
	if len(b) < 586256 {
		return io.ErrShortBuffer
	}
	v := (*Schema)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Schema) Clone() *Schema {
	v := &Schema{}
	*v = *s
	return v
}
func (s *Schema) Bytes() []byte {
	return (*(*[586256]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Schema) Mut() *SchemaMut {
	return (*SchemaMut)(unsafe.Pointer(s))
}
func (s *Schema) Imports() *Imports16List {
	return &s.imports
}
func (s *Schema) Records() *Record128List {
	return &s.records
}

type SchemaMut struct {
	Schema
}

func (s *SchemaMut) Clone() *SchemaMut {
	v := &SchemaMut{}
	*v = *s
	return v
}
func (s *SchemaMut) Freeze() *Schema {
	return (*Schema)(unsafe.Pointer(s))
}
func (s *SchemaMut) Imports() *Imports16ListMut {
	return s.imports.Mut()
}
func (s *SchemaMut) SetImports(v *Imports16List) *SchemaMut {
	s.imports = *v
	return s
}
func (s *SchemaMut) Records() *Record128ListMut {
	return s.records.Mut()
}
func (s *SchemaMut) SetRecords(v *Record128List) *SchemaMut {
	s.records = *v
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

type Field struct {
	name       String40
	compact    String8
	offset     uint16
	rootOffset uint16
	size       uint16
	align      uint16
	number     uint16
	kind       Kind
	isOptional bool
	isPointer  bool
	_          [3]byte // Padding
}

func (s *Field) String() string {
	return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Field) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["name"] = s.Name()
	m["compact"] = s.Compact()
	m["offset"] = s.Offset()
	m["rootOffset"] = s.RootOffset()
	m["size"] = s.Size()
	m["align"] = s.Align()
	m["number"] = s.Number()
	m["kind"] = s.Kind()
	m["isOptional"] = s.IsOptional()
	m["isPointer"] = s.IsPointer()
	return m
}

func (s *Field) ReadFrom(r io.Reader) (int64, error) {
	n, err := io.ReadFull(r, (*(*[64]byte)(unsafe.Pointer(s)))[0:])
	if err != nil {
		return int64(n), err
	}
	if n != 64 {
		return int64(n), io.ErrShortBuffer
	}
	return int64(n), nil
}
func (s *Field) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write((*(*[64]byte)(unsafe.Pointer(s)))[0:])
	return int64(n), err
}
func (s *Field) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[64]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Field) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[64]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Field) Read(b []byte) (n int, err error) {
	if len(b) < 64 {
		return -1, io.ErrShortBuffer
	}
	v := (*Field)(unsafe.Pointer(&b[0]))
	*v = *s
	return 64, nil
}
func (s *Field) UnmarshalBinary(b []byte) error {
	if len(b) < 64 {
		return io.ErrShortBuffer
	}
	v := (*Field)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Field) Clone() *Field {
	v := &Field{}
	*v = *s
	return v
}
func (s *Field) Bytes() []byte {
	return (*(*[64]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Field) Mut() *FieldMut {
	return (*FieldMut)(unsafe.Pointer(s))
}
func (s *Field) Name() *String40 {
	return &s.name
}
func (s *Field) Compact() *String8 {
	return &s.compact
}
func (s *Field) Offset() uint16 {
	return s.offset
}
func (s *Field) RootOffset() uint16 {
	return s.rootOffset
}
func (s *Field) Size() uint16 {
	return s.size
}
func (s *Field) Align() uint16 {
	return s.align
}
func (s *Field) Number() uint16 {
	return s.number
}
func (s *Field) Kind() Kind {
	return s.kind
}
func (s *Field) IsOptional() bool {
	return s.isOptional
}
func (s *Field) IsPointer() bool {
	return s.isPointer
}

type FieldMut struct {
	Field
}

func (s *FieldMut) Clone() *FieldMut {
	v := &FieldMut{}
	*v = *s
	return v
}
func (s *FieldMut) Freeze() *Field {
	return (*Field)(unsafe.Pointer(s))
}
func (s *FieldMut) Name() *String40Mut {
	return s.name.Mut()
}
func (s *FieldMut) SetName(v *String40) *FieldMut {
	s.name = *v
	return s
}
func (s *FieldMut) Compact() *String8Mut {
	return s.compact.Mut()
}
func (s *FieldMut) SetCompact(v *String8) *FieldMut {
	s.compact = *v
	return s
}
func (s *FieldMut) SetOffset(v uint16) *FieldMut {
	s.offset = v
	return s
}
func (s *FieldMut) SetRootOffset(v uint16) *FieldMut {
	s.rootOffset = v
	return s
}
func (s *FieldMut) SetSize(v uint16) *FieldMut {
	s.size = v
	return s
}
func (s *FieldMut) SetAlign(v uint16) *FieldMut {
	s.align = v
	return s
}
func (s *FieldMut) SetNumber(v uint16) *FieldMut {
	s.number = v
	return s
}
func (s *FieldMut) SetKind(v Kind) *FieldMut {
	s.kind = v
	return s
}
func (s *FieldMut) SetIsOptional(v bool) *FieldMut {
	s.isOptional = v
	return s
}
func (s *FieldMut) SetIsPointer(v bool) *FieldMut {
	s.isPointer = v
	return s
}

type Stopped struct {
	recordID  RecordID
	timestamp int64
	starts    int64
	reason    StopReason
	_         [7]byte // Padding
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
	m["starts"] = s.Starts()
	m["reason"] = s.Reason()
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
func (s *Stopped) Starts() int64 {
	return s.starts
}
func (s *Stopped) Reason() StopReason {
	return s.reason
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
func (s *StoppedMut) SetStarts(v int64) *StoppedMut {
	s.starts = v
	return s
}
func (s *StoppedMut) SetReason(v StopReason) *StoppedMut {
	s.reason = v
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
func (s *String32) StringClone() string {
	b := s[0:s.Len()]
	return string(b)
}
func (s *String32) String() string {
	b := s[0:s.Len()]
	return *(*string)(unsafe.Pointer(&b))
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
	return s[0:s.Len()], nil
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

type String40 [40]byte

func NewString40(s string) *String40 {
	v := String40{}
	v.set(s)
	return &v
}
func (s *String40) set(v string) {
	copy(s[0:39], v)
	c := 39
	l := len(v)
	if l > c {
		s[39] = byte(c)
	} else {
		s[39] = byte(l)
	}
}
func (s *String40) Len() int {
	return int(s[39])
}
func (s *String40) Cap() int {
	return 39
}
func (s *String40) StringClone() string {
	b := s[0:s.Len()]
	return string(b)
}
func (s *String40) String() string {
	b := s[0:s.Len()]
	return *(*string)(unsafe.Pointer(&b))
}
func (s *String40) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String40) Clone() *String40 {
	v := String40{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String40) Mut() *String40Mut {
	return *(**String40Mut)(unsafe.Pointer(&s))
}
func (s *String40) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[40]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 40 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String40) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[40]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String40) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[40]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String40) MarshalBinary() ([]byte, error) {
	return s[0:s.Len()], nil
}
func (s *String40) UnmarshalBinary(b []byte) error {
	if len(b) < 40 {
		return io.ErrShortBuffer
	}
	v := (*String40)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String40Mut struct {
	String40
}

func (s *String40Mut) Set(v string) {
	s.set(v)
}

type String8 [8]byte

func NewString8(s string) *String8 {
	v := String8{}
	v.set(s)
	return &v
}
func (s *String8) set(v string) {
	copy(s[0:7], v)
	c := 7
	l := len(v)
	if l > c {
		s[7] = byte(c)
	} else {
		s[7] = byte(l)
	}
}
func (s *String8) Len() int {
	return int(s[7])
}
func (s *String8) Cap() int {
	return 7
}
func (s *String8) StringClone() string {
	b := s[0:s.Len()]
	return string(b)
}
func (s *String8) String() string {
	b := s[0:s.Len()]
	return *(*string)(unsafe.Pointer(&b))
}
func (s *String8) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String8) Clone() *String8 {
	v := String8{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String8) Mut() *String8Mut {
	return *(**String8Mut)(unsafe.Pointer(&s))
}
func (s *String8) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[8]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 8 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String8) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[8]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String8) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[8]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String8) MarshalBinary() ([]byte, error) {
	return s[0:s.Len()], nil
}
func (s *String8) UnmarshalBinary(b []byte) error {
	if len(b) < 8 {
		return io.ErrShortBuffer
	}
	v := (*String8)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String8Mut struct {
	String8
}

func (s *String8Mut) Set(v string) {
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
func (s *String64) StringClone() string {
	b := s[0:s.Len()]
	return string(b)
}
func (s *String64) String() string {
	b := s[0:s.Len()]
	return *(*string)(unsafe.Pointer(&b))
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
	return s[0:s.Len()], nil
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

type String128 [128]byte

func NewString128(s string) *String128 {
	v := String128{}
	v.set(s)
	return &v
}
func (s *String128) set(v string) {
	copy(s[0:127], v)
	c := 127
	l := len(v)
	if l > c {
		s[127] = byte(c)
	} else {
		s[127] = byte(l)
	}
}
func (s *String128) Len() int {
	return int(s[127])
}
func (s *String128) Cap() int {
	return 127
}
func (s *String128) StringClone() string {
	b := s[0:s.Len()]
	return string(b)
}
func (s *String128) String() string {
	b := s[0:s.Len()]
	return *(*string)(unsafe.Pointer(&b))
}
func (s *String128) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *String128) Clone() *String128 {
	v := String128{}
	copy(s[0:], v[0:])
	return &v
}
func (s *String128) Mut() *String128Mut {
	return *(**String128Mut)(unsafe.Pointer(&s))
}
func (s *String128) ReadFrom(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[128]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 128 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *String128) WriteTo(w io.Writer) (n int, err error) {
	return w.Write((*(*[128]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *String128) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[128]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *String128) MarshalBinary() ([]byte, error) {
	return s[0:s.Len()], nil
}
func (s *String128) UnmarshalBinary(b []byte) error {
	if len(b) < 128 {
		return io.ErrShortBuffer
	}
	v := (*String128)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type String128Mut struct {
	String128
}

func (s *String128Mut) Set(v string) {
	s.set(v)
}

type Field64List struct {
	b [64]Field
	_ [7]byte // Padding
	l byte
}

func (s *Field64List) Get(i int) *Field {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return &s.b[i]
}
func (s *Field64List) Len() int {
	return int(s.l)
}
func (s *Field64List) Cap() int {
	return 64
}
func (s *Field64List) MarshalMap(m []map[string]interface{}) []map[string]interface{} {
	if m == nil {
		m = make([]map[string]interface{}, 0, s.Len())
	}
	for _, v := range s.Unsafe() {
		m = append(m, v.MarshalMap(nil))
	}
	return m
}
func (s *Field64List) CopyTo(v []Field) []Field {
	return append(v, s.Unsafe()...)
}
func (s *Field64List) Unsafe() []Field {
	return s.b[0:s.Len()]
}
func (s *Field64List) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&s.b[0])),
		Len:  s.Len() * 64,
		Cap:  4096,
	}))
}
func (s *Field64List) Mut() *Field64ListMut {
	return *(**Field64ListMut)(unsafe.Pointer(&s))
}
func (s *Field64List) Read(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[4104]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 4104 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *Field64List) Write(w io.Writer) (n int, err error) {
	return w.Write((*(*[4104]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *Field64List) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[4104]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *Field64List) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[4104]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *Field64List) UnmarshalBinary(b []byte) error {
	if len(b) < 4104 {
		return io.ErrShortBuffer
	}
	v := (*Field64List)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type Field64ListMut struct {
	Field64List
}

func (s *Field64ListMut) setLen(l int) {
	s.l = byte(l)
}
func (s *Field64ListMut) Push(v *Field) bool {
	l := s.Len()
	if l == 64 {
		return false
	}
	s.b[l] = *v
	s.setLen(l + 1)
	return true
}

// Removes the last item
func (s *Field64ListMut) Pop(v *Field) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	l -= 1
	if v != nil {
		*v = s.b[l]
	}
	s.b[l] = Field{}
	s.setLen(l)
	return true
}

// Removes the first item
func (s *Field64ListMut) Shift(v *Field) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	if v != nil {
		*v = s.b[0]
	}
	if l > 1 {
		copy(s.b[0:], s.b[1:s.l])
	}
	l -= 1
	s.b[l] = Field{}
	s.setLen(l)
	return true
}
func (s *Field64ListMut) Clear() {
	s.b = [64]Field{}
	s.l = 0
}

type UnionOption16List struct {
	b [16]UnionOption
	_ [7]byte // Padding
	l byte
}

func (s *UnionOption16List) Get(i int) *UnionOption {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return &s.b[i]
}
func (s *UnionOption16List) Len() int {
	return int(s.l)
}
func (s *UnionOption16List) Cap() int {
	return 16
}
func (s *UnionOption16List) MarshalMap(m []map[string]interface{}) []map[string]interface{} {
	if m == nil {
		m = make([]map[string]interface{}, 0, s.Len())
	}
	for _, v := range s.Unsafe() {
		m = append(m, v.MarshalMap(nil))
	}
	return m
}
func (s *UnionOption16List) CopyTo(v []UnionOption) []UnionOption {
	return append(v, s.Unsafe()...)
}
func (s *UnionOption16List) Unsafe() []UnionOption {
	return s.b[0:s.Len()]
}
func (s *UnionOption16List) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&s.b[0])),
		Len:  s.Len() * 88,
		Cap:  1408,
	}))
}
func (s *UnionOption16List) Mut() *UnionOption16ListMut {
	return *(**UnionOption16ListMut)(unsafe.Pointer(&s))
}
func (s *UnionOption16List) Read(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[1416]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 1416 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *UnionOption16List) Write(w io.Writer) (n int, err error) {
	return w.Write((*(*[1416]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *UnionOption16List) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[1416]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *UnionOption16List) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[1416]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *UnionOption16List) UnmarshalBinary(b []byte) error {
	if len(b) < 1416 {
		return io.ErrShortBuffer
	}
	v := (*UnionOption16List)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type UnionOption16ListMut struct {
	UnionOption16List
}

func (s *UnionOption16ListMut) setLen(l int) {
	s.l = byte(l)
}
func (s *UnionOption16ListMut) Push(v *UnionOption) bool {
	l := s.Len()
	if l == 16 {
		return false
	}
	s.b[l] = *v
	s.setLen(l + 1)
	return true
}

// Removes the last item
func (s *UnionOption16ListMut) Pop(v *UnionOption) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	l -= 1
	if v != nil {
		*v = s.b[l]
	}
	s.b[l] = UnionOption{}
	s.setLen(l)
	return true
}

// Removes the first item
func (s *UnionOption16ListMut) Shift(v *UnionOption) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	if v != nil {
		*v = s.b[0]
	}
	if l > 1 {
		copy(s.b[0:], s.b[1:s.l])
	}
	l -= 1
	s.b[l] = UnionOption{}
	s.setLen(l)
	return true
}
func (s *UnionOption16ListMut) Clear() {
	s.b = [16]UnionOption{}
	s.l = 0
}

type Import16List struct {
	b [16]Import
	_ [7]byte // Padding
	l byte
}

func (s *Import16List) Get(i int) *Import {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return &s.b[i]
}
func (s *Import16List) Len() int {
	return int(s.l)
}
func (s *Import16List) Cap() int {
	return 16
}
func (s *Import16List) MarshalMap(m []map[string]interface{}) []map[string]interface{} {
	if m == nil {
		m = make([]map[string]interface{}, 0, s.Len())
	}
	for _, v := range s.Unsafe() {
		m = append(m, v.MarshalMap(nil))
	}
	return m
}
func (s *Import16List) CopyTo(v []Import) []Import {
	return append(v, s.Unsafe()...)
}
func (s *Import16List) Unsafe() []Import {
	return s.b[0:s.Len()]
}
func (s *Import16List) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&s.b[0])),
		Len:  s.Len() * 216,
		Cap:  3456,
	}))
}
func (s *Import16List) Mut() *Import16ListMut {
	return *(**Import16ListMut)(unsafe.Pointer(&s))
}
func (s *Import16List) Read(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[3464]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 3464 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *Import16List) Write(w io.Writer) (n int, err error) {
	return w.Write((*(*[3464]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *Import16List) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[3464]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *Import16List) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[3464]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *Import16List) UnmarshalBinary(b []byte) error {
	if len(b) < 3464 {
		return io.ErrShortBuffer
	}
	v := (*Import16List)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type Import16ListMut struct {
	Import16List
}

func (s *Import16ListMut) setLen(l int) {
	s.l = byte(l)
}
func (s *Import16ListMut) Push(v *Import) bool {
	l := s.Len()
	if l == 16 {
		return false
	}
	s.b[l] = *v
	s.setLen(l + 1)
	return true
}

// Removes the last item
func (s *Import16ListMut) Pop(v *Import) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	l -= 1
	if v != nil {
		*v = s.b[l]
	}
	s.b[l] = Import{}
	s.setLen(l)
	return true
}

// Removes the first item
func (s *Import16ListMut) Shift(v *Import) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	if v != nil {
		*v = s.b[0]
	}
	if l > 1 {
		copy(s.b[0:], s.b[1:s.l])
	}
	l -= 1
	s.b[l] = Import{}
	s.setLen(l)
	return true
}
func (s *Import16ListMut) Clear() {
	s.b = [16]Import{}
	s.l = 0
}

type Imports16List struct {
	b [16]Imports
	_ [7]byte // Padding
	l byte
}

func (s *Imports16List) Get(i int) *Imports {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return &s.b[i]
}
func (s *Imports16List) Len() int {
	return int(s.l)
}
func (s *Imports16List) Cap() int {
	return 16
}
func (s *Imports16List) MarshalMap(m []map[string]interface{}) []map[string]interface{} {
	if m == nil {
		m = make([]map[string]interface{}, 0, s.Len())
	}
	for _, v := range s.Unsafe() {
		m = append(m, v.MarshalMap(nil))
	}
	return m
}
func (s *Imports16List) CopyTo(v []Imports) []Imports {
	return append(v, s.Unsafe()...)
}
func (s *Imports16List) Unsafe() []Imports {
	return s.b[0:s.Len()]
}
func (s *Imports16List) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&s.b[0])),
		Len:  s.Len() * 3488,
		Cap:  55808,
	}))
}
func (s *Imports16List) Mut() *Imports16ListMut {
	return *(**Imports16ListMut)(unsafe.Pointer(&s))
}
func (s *Imports16List) Read(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[55816]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 55816 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *Imports16List) Write(w io.Writer) (n int, err error) {
	return w.Write((*(*[55816]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *Imports16List) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[55816]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *Imports16List) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[55816]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *Imports16List) UnmarshalBinary(b []byte) error {
	if len(b) < 55816 {
		return io.ErrShortBuffer
	}
	v := (*Imports16List)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type Imports16ListMut struct {
	Imports16List
}

func (s *Imports16ListMut) setLen(l int) {
	s.l = byte(l)
}
func (s *Imports16ListMut) Push(v *Imports) bool {
	l := s.Len()
	if l == 16 {
		return false
	}
	s.b[l] = *v
	s.setLen(l + 1)
	return true
}

// Removes the last item
func (s *Imports16ListMut) Pop(v *Imports) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	l -= 1
	if v != nil {
		*v = s.b[l]
	}
	s.b[l] = Imports{}
	s.setLen(l)
	return true
}

// Removes the first item
func (s *Imports16ListMut) Shift(v *Imports) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	if v != nil {
		*v = s.b[0]
	}
	if l > 1 {
		copy(s.b[0:], s.b[1:s.l])
	}
	l -= 1
	s.b[l] = Imports{}
	s.setLen(l)
	return true
}
func (s *Imports16ListMut) Clear() {
	s.b = [16]Imports{}
	s.l = 0
}

type Record128List struct {
	b [128]Record
	_ [7]byte // Padding
	l byte
}

func (s *Record128List) Get(i int) *Record {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return &s.b[i]
}
func (s *Record128List) Len() int {
	return int(s.l)
}
func (s *Record128List) Cap() int {
	return 128
}
func (s *Record128List) MarshalMap(m []map[string]interface{}) []map[string]interface{} {
	if m == nil {
		m = make([]map[string]interface{}, 0, s.Len())
	}
	for _, v := range s.Unsafe() {
		m = append(m, v.MarshalMap(nil))
	}
	return m
}
func (s *Record128List) CopyTo(v []Record) []Record {
	return append(v, s.Unsafe()...)
}
func (s *Record128List) Unsafe() []Record {
	return s.b[0:s.Len()]
}
func (s *Record128List) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&s.b[0])),
		Len:  s.Len() * 4144,
		Cap:  530432,
	}))
}
func (s *Record128List) Mut() *Record128ListMut {
	return *(**Record128ListMut)(unsafe.Pointer(&s))
}
func (s *Record128List) Read(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[530440]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 530440 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *Record128List) Write(w io.Writer) (n int, err error) {
	return w.Write((*(*[530440]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *Record128List) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[530440]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *Record128List) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[530440]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *Record128List) UnmarshalBinary(b []byte) error {
	if len(b) < 530440 {
		return io.ErrShortBuffer
	}
	v := (*Record128List)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type Record128ListMut struct {
	Record128List
}

func (s *Record128ListMut) setLen(l int) {
	s.l = byte(l)
}
func (s *Record128ListMut) Push(v *Record) bool {
	l := s.Len()
	if l == 128 {
		return false
	}
	s.b[l] = *v
	s.setLen(l + 1)
	return true
}

// Removes the last item
func (s *Record128ListMut) Pop(v *Record) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	l -= 1
	if v != nil {
		*v = s.b[l]
	}
	s.b[l] = Record{}
	s.setLen(l)
	return true
}

// Removes the first item
func (s *Record128ListMut) Shift(v *Record) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	if v != nil {
		*v = s.b[0]
	}
	if l > 1 {
		copy(s.b[0:], s.b[1:s.l])
	}
	l -= 1
	s.b[l] = Record{}
	s.setLen(l)
	return true
}
func (s *Record128ListMut) Clear() {
	s.b = [128]Record{}
	s.l = 0
}

type EnumOption16List struct {
	b [16]EnumOption
	_ [7]byte // Padding
	l byte
}

func (s *EnumOption16List) Get(i int) *EnumOption {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return &s.b[i]
}
func (s *EnumOption16List) Len() int {
	return int(s.l)
}
func (s *EnumOption16List) Cap() int {
	return 16
}
func (s *EnumOption16List) MarshalMap(m []map[string]interface{}) []map[string]interface{} {
	if m == nil {
		m = make([]map[string]interface{}, 0, s.Len())
	}
	for _, v := range s.Unsafe() {
		m = append(m, v.MarshalMap(nil))
	}
	return m
}
func (s *EnumOption16List) CopyTo(v []EnumOption) []EnumOption {
	return append(v, s.Unsafe()...)
}
func (s *EnumOption16List) Unsafe() []EnumOption {
	return s.b[0:s.Len()]
}
func (s *EnumOption16List) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&s.b[0])),
		Len:  s.Len() * 144,
		Cap:  2304,
	}))
}
func (s *EnumOption16List) Mut() *EnumOption16ListMut {
	return *(**EnumOption16ListMut)(unsafe.Pointer(&s))
}
func (s *EnumOption16List) Read(r io.Reader) error {
	n, err := io.ReadFull(r, (*(*[2312]byte)(unsafe.Pointer(&s)))[0:])
	if err != nil {
		return err
	}
	if n != 2312 {
		return io.ErrShortBuffer
	}
	return nil
}
func (s *EnumOption16List) Write(w io.Writer) (n int, err error) {
	return w.Write((*(*[2312]byte)(unsafe.Pointer(&s)))[0:])
}
func (s *EnumOption16List) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[2312]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *EnumOption16List) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[2312]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *EnumOption16List) UnmarshalBinary(b []byte) error {
	if len(b) < 2312 {
		return io.ErrShortBuffer
	}
	v := (*EnumOption16List)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type EnumOption16ListMut struct {
	EnumOption16List
}

func (s *EnumOption16ListMut) setLen(l int) {
	s.l = byte(l)
}
func (s *EnumOption16ListMut) Push(v *EnumOption) bool {
	l := s.Len()
	if l == 16 {
		return false
	}
	s.b[l] = *v
	s.setLen(l + 1)
	return true
}

// Removes the last item
func (s *EnumOption16ListMut) Pop(v *EnumOption) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	l -= 1
	if v != nil {
		*v = s.b[l]
	}
	s.b[l] = EnumOption{}
	s.setLen(l)
	return true
}

// Removes the first item
func (s *EnumOption16ListMut) Shift(v *EnumOption) bool {
	l := s.Len()
	if l == 0 {
		return false
	}
	if v != nil {
		*v = s.b[0]
	}
	if l > 1 {
		copy(s.b[0:], s.b[1:s.l])
	}
	l -= 1
	s.b[l] = EnumOption{}
	s.setLen(l)
	return true
}
func (s *EnumOption16ListMut) Clear() {
	s.b = [16]EnumOption{}
	s.l = 0
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

	a(Starting{}, StartingMut{}, 40, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"writerID", 32, 8},
	})
	a(Stream{}, StreamMut{}, 80, []b{
		{"id", 0, 8},
		{"created", 8, 8},
		{"accountID", 16, 8},
		{"duration", 24, 8},
		{"record", 32, 2},
		{"_", 34, 6},
		{"name", 40, 32},
		{"kind", 72, 1},
		{"format", 73, 1},
		{"blockSize", 74, 2},
		{"_", 76, 4},
	})
	a(Started{}, StartedMut{}, 48, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"writerID", 32, 8},
		{"stops", 40, 8},
	})
	a(UnionOption{}, UnionOptionMut{}, 88, []b{
		{"name", 0, 40},
		{"kind", 40, 1},
		{"_", 41, 7},
		{"id", 48, 40},
	})
	a(Line{}, LineMut{}, 16, []b{
		{"number", 0, 4},
		{"begin", 4, 4},
		{"end", 8, 4},
		{"_", 12, 4},
	})
	a(Import{}, ImportMut{}, 216, []b{
		{"id", 0, 4},
		{"_", 4, 4},
		{"line", 8, 16},
		{"path", 24, 128},
		{"name", 152, 32},
		{"alias", 184, 32},
	})
	a(Imports{}, ImportsMut{}, 3488, []b{
		{"id", 0, 4},
		{"_", 4, 4},
		{"line", 8, 16},
		{"list", 24, 3464},
	})
	a(BlockHeader{}, BlockHeaderMut{}, 136, []b{
		{"streamID", 0, 8},
		{"id", 8, 8},
		{"headID", 16, 8},
		{"headMin", 24, 8},
		{"headStart", 32, 8},
		{"blocks", 40, 8},
		{"records", 48, 8},
		{"storage", 56, 8},
		{"storageU", 64, 8},
		{"created", 72, 8},
		{"completed", 80, 8},
		{"start", 88, 8},
		{"end", 96, 8},
		{"min", 104, 8},
		{"max", 112, 8},
		{"count", 120, 2},
		{"size", 122, 2},
		{"sizeU", 124, 2},
		{"sizeX", 126, 2},
		{"record", 128, 2},
		{"blockSize", 130, 2},
		{"encoding", 132, 1},
		{"kind", 133, 1},
		{"format", 134, 1},
		{"_", 135, 1},
	})
	a(Record{}, RecordMut{}, 4144, []b{
		{"name", 0, 40},
		{"fields", 40, 4104},
	})
	a(SyncProgress{}, SyncProgressMut{}, 56, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"started", 32, 8},
		{"count", 40, 8},
		{"remaining", 48, 8},
	})
	a(Enum{}, EnumMut{}, 2352, []b{
		{"name", 0, 40},
		{"options", 40, 2312},
	})
	a(RecordHeader{}, RecordHeaderMut{}, 64, []b{
		{"id", 0, 24},
		{"prevID", 24, 8},
		{"timestamp", 32, 8},
		{"start", 40, 8},
		{"end", 48, 8},
		{"seq", 56, 2},
		{"size", 58, 2},
		{"sizeU", 60, 2},
		{"encoding", 62, 1},
		{"pad", 63, 1},
	})
	a(Struct{}, StructMut{}, 4144, []b{
		{"name", 0, 40},
		{"fields", 40, 4104},
	})
	a(Union{}, UnionMut{}, 1456, []b{
		{"name", 0, 40},
		{"options", 40, 1416},
	})
	a(EOSWaiting{}, EOSWaitingMut{}, 32, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
	})
	a(SyncStarted{}, SyncStartedMut{}, 32, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
	})
	a(RecordsHeader{}, RecordsHeaderMut{}, 72, []b{
		{"header", 0, 64},
		{"count", 64, 2},
		{"record", 66, 2},
		{"_", 68, 4},
	})
	a(RecordID{}, RecordIDMut{}, 24, []b{
		{"streamID", 0, 8},
		{"blockID", 8, 8},
		{"id", 16, 8},
	})
	a(SyncStopped{}, SyncStoppedMut{}, 128, []b{
		{"progress", 0, 56},
		{"message", 56, 64},
		{"reason", 120, 1},
		{"_", 121, 7},
	})
	a(EnumOption{}, EnumOptionMut{}, 144, []b{
		{"index", 0, 4},
		{"_", 4, 4},
		{"name", 8, 64},
		{"value", 72, 8},
		{"valueStr", 80, 64},
	})
	a(Schema{}, SchemaMut{}, 586256, []b{
		{"imports", 0, 55816},
		{"records", 55816, 530440},
	})
	a(BlockID{}, BlockIDMut{}, 16, []b{
		{"streamID", 0, 8},
		{"id", 8, 8},
	})
	a(Savepoint{}, SavepointMut{}, 40, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"duration", 32, 8},
	})
	a(EOS{}, EOSMut{}, 32, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
	})
	a(Field{}, FieldMut{}, 64, []b{
		{"name", 0, 40},
		{"compact", 40, 8},
		{"offset", 48, 2},
		{"rootOffset", 50, 2},
		{"size", 52, 2},
		{"align", 54, 2},
		{"number", 56, 2},
		{"kind", 58, 1},
		{"isOptional", 59, 1},
		{"isPointer", 60, 1},
		{"_", 61, 3},
	})
	a(Stopped{}, StoppedMut{}, 48, []b{
		{"recordID", 0, 24},
		{"timestamp", 24, 8},
		{"starts", 32, 8},
		{"reason", 40, 1},
		{"_", 41, 7},
	})

}
