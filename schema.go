package wap

import (
	"errors"
	"math"
)

var (
	ErrColumnPointer = errors.New("columns cannot be pointer type")
)

type Kind byte

const (
	KindUnknown       = Kind(0) // Possibly resolved later
	KindBool          = Kind(1)
	KindByte          = Kind(2)
	KindInt8          = Kind(3)
	KindInt16         = Kind(4)
	KindUInt16        = Kind(5)
	KindInt32         = Kind(6)
	KindUInt32        = Kind(7)
	KindInt64         = Kind(8)
	KindUInt64        = Kind(9)
	KindFloat32       = Kind(10)
	KindFloat64       = Kind(11)
	KindString        = Kind(12)
	KindStringInline  = Kind(13)
	KindBytes         = Kind(14)
	KindBytesInline   = Kind(15)
	KindRecordHeader  = Kind(17)
	KindBlockHeader   = Kind(18)
	KindRef           = Kind(25)
	KindEnum          = Kind(30) // User-defined enum
	KindRecord        = Kind(40) // User-defined variable sized record
	KindStruct        = Kind(41) // User-defined fixed sized record
	KindList          = Kind(50)
	KindListInline    = Kind(51)
	KindMap           = Kind(60)
	KindMapInline     = Kind(61)
	KindUnion         = Kind(70) // User-defined union
	KindUnionUntagged = Kind(71) // User-defined union
	KindStream        = Kind(80)
	KindPad           = Kind(100) // struct alignment padding
)

type BlockSize uint16

const (
	BlockSize1KB  BlockSize = 1024
	BlockSize2KB  BlockSize = 2048
	BlockSize4KB  BlockSize = 4096
	BlockSize8KB  BlockSize = 8192
	BlockSize16KB BlockSize = 16384
	BlockSize32KB BlockSize = 32768
	BlockSize64KB BlockSize = math.MaxUint16
)

type Format byte

const (
	FormatRaw      Format = 0
	FormatJson     Format = 1
	FormatMoon     Format = 2
	FormatProtobuf Format = 3
	FormatMsgpack  Format = 4
)

type Encoding byte

const (
	EncodingNone   Encoding = 0
	EncodingLZ4    Encoding = 1
	EncodingZSTD   Encoding = 2
	EncodingBrotli Encoding = 3
	EncodingGzip   Encoding = 4
)

type UnionKind byte

const (
	UnionKindTagged   UnionKind = 0
	UnionKindUntagged UnionKind = 1
)

type StreamKind byte

const (
	StreamKindLog    StreamKind = 0
	StreamKindSeries StreamKind = 1
)

type BlockLayout byte

const (
	BlockLayoutRow    BlockLayout = 0
	BlockLayoutColumn BlockLayout = 1
)

const (
	VPointerSize = int32(4)
)

func NewField(name, short string, dataType Type) Field {
	return Field{name: name, kind: dataType}
}

func NumberType(kind Kind) Type {
	switch kind {
	case KindInt8, KindByte:
		return Type{kind: kind, size: 1}
	case KindInt16, KindUInt16:
		return Type{kind: kind, size: 2}
	case KindInt32, KindUInt32, KindFloat32:
		return Type{kind: kind, size: 4}
	case KindInt64, KindUInt64, KindFloat64:
		return Type{kind: kind, size: 8}
	}
	return Type{kind: kind}
}

func StringType() Type {
	return Type{kind: KindString, size: VPointerSize}
}
func StringInlineType(size int32) Type {
	return Type{kind: KindStringInline, size: size}
}
func BytesType() Type {
	return Type{kind: KindBytes, size: VPointerSize}
}
func BytesInlineType(size int32) Type {
	return Type{kind: KindBytesInline, size: size}
}

func ListType(element Type) Type {
	return Type{kind: KindList, size: VPointerSize, list: &List{
		element: element,
	}}
}

type Line struct {
	number int32
	begin  int32
	end    int32
}

type RecordLayout int32

const (
	RecordLayoutCompact RecordLayout = 0
	RecordLayoutAligned RecordLayout = 1
)

type Record struct {
	line           Line
	name           string
	comments       []string
	fields         []Field
	fieldsMap      map[string]*Field
	version        int64
	size           int32
	offset         int32
	versionOffset  int32
	versionSize    int32
	optionals      int32
	optionalOffset int32
	optionalSize   int32
	unions         int32
	unionOffset    int32
	unionSize      int32
	align          int32
	layout         RecordLayout
	fixed          bool
}

func (r *Record) Line() Line {
	return r.line
}

func (r *Record) SetLine(line Line) {
	r.line = line
}

func (r *Record) Name() string {
	return r.name
}

func (r *Record) SetName(name string) {
	r.name = name
}

func (r *Record) Comments() []string {
	return r.comments
}

func (r *Record) SetComments(comments []string) {
	r.comments = comments
}

func (r *Record) Fields() []Field {
	return r.fields
}

func (r *Record) SetFields(fields []Field) {
	r.fields = fields
}

func (r *Record) FieldsMap() map[string]*Field {
	return r.fieldsMap
}

func (r *Record) SetFieldsMap(fieldsMap map[string]*Field) {
	r.fieldsMap = fieldsMap
}

func (r *Record) Version() int64 {
	return r.version
}

func (r *Record) SetVersion(version int64) {
	r.version = version
}

func (r *Record) Size() int32 {
	return r.size
}

func (r *Record) SetSize(size int32) {
	r.size = size
}

func (r *Record) Offset() int32 {
	return r.offset
}

func (r *Record) SetOffset(offset int32) {
	r.offset = offset
}

func (r *Record) VersionOffset() int32 {
	return r.versionOffset
}

func (r *Record) SetVersionOffset(versionOffset int32) {
	r.versionOffset = versionOffset
}

func (r *Record) VersionSize() int32 {
	return r.versionSize
}

func (r *Record) SetVersionSize(versionSize int32) {
	r.versionSize = versionSize
}

func (r *Record) Optionals() int32 {
	return r.optionals
}

func (r *Record) SetOptionals(optionals int32) {
	r.optionals = optionals
}

func (r *Record) OptionalOffset() int32 {
	return r.optionalOffset
}

func (r *Record) SetOptionalOffset(optionalOffset int32) {
	r.optionalOffset = optionalOffset
}

func (r *Record) OptionalSize() int32 {
	return r.optionalSize
}

func (r *Record) SetOptionalSize(optionalSize int32) {
	r.optionalSize = optionalSize
}

func (r *Record) Align() int32 {
	return r.align
}

func (r *Record) SetAlign(align int32) {
	r.align = align
}

func (r *Record) Layout() RecordLayout {
	return r.layout
}

func (r *Record) SetLayout(layout RecordLayout) {
	r.layout = layout
}

func (r *Record) Fixed() bool {
	return r.fixed
}

func (r *Record) SetFixed(fixed bool) {
	r.fixed = fixed
}

func (r *Record) Field(name string) *Field {
	m := r.fieldsMap
	if m == nil {
		m = make(map[string]*Field)
	}
	if len(m) == 0 {
		for i := 0; i < len(r.fields); i++ {
			field := &r.fields[i]
			m[field.name] = field
			if len(field.short) > 0 {
				m[field.short] = field
			}
		}
		r.fieldsMap = m
	}
	return m[name]
}

func (r *Record) FieldAt(index int) *Field {
	if index < 0 || index >= len(r.Fields()) {
		return nil
	}
	return &r.Fields()[index]
}

func (r *Record) LayoutColumns(offset int32) error {
	for i := 0; i < len(r.Fields()); i++ {
		field := &r.Fields()[i]
		if field.IsPointer() {
			return ErrColumnPointer
		}
		field.SetColumnOffset(offset)

		switch field.Kind().Kind() {
		case KindRecord:
			if field.Kind().Record() != nil {
				if err := field.Kind().Record().LayoutColumns(offset); err != nil {
					return err
				}
			}
		}

		offset += field.Align()
	}
	return nil
}

func (r *Record) ToColumns(columns []*Field) []*Field {
	if columns == nil {
		columns = make([]*Field, 0, 16)
	}
	for i := 0; i < len(r.Fields()); i++ {
		field := &r.Fields()[i]
		switch field.Kind().Kind() {
		case KindRecord:
			if field.Kind().Record() != nil {
				field.Kind().Record().ToColumns(columns)
			}
		}
	}
	return columns
}

// List represents an Array like structure.
// Header: | LEN 2 bytes | [LEN]Element
type List struct {
	element Type
	fixed   int
}

func (l *List) Element() *Type {
	return &l.element
}

func (l *List) SetElement(element Type) {
	l.element = element
}

func (l *List) Fixed() int {
	return l.fixed
}

func (l *List) SetFixed(fixed int) {
	l.fixed = fixed
}

// Map represents a HashMap data structure using a robin-hood algorithm
// Header: | LEN 2 bytes | SIZE 2 bytes | List<MapEntry>
// Item: | KEY | VALUE | Distance (2 bytes)
type Map struct {
	key   Type
	value Type
	size  int32
}

func (m *Map) Key() Type {
	return m.key
}

func (m *Map) SetKey(key Type) {
	m.key = key
}

func (m *Map) Value() Type {
	return m.value
}

func (m *Map) SetValue(value Type) {
	m.value = value
}

func (m *Map) Size() int32 {
	return m.size
}

func (m *Map) SetSize(size int32) {
	m.size = size
}

// Union represents a C-like union or a protobuf oneOf
type Union struct {
	name   string
	fields []Field
	size   int32
	align  int32
	kind   UnionKind
}

func (union *Union) Name() string {
	return union.name
}

func (union *Union) SetName(name string) {
	union.name = name
}

func (union *Union) Fields() []Field {
	return union.fields
}

func (union *Union) SetFields(fields []Field) {
	union.fields = fields
}

func (union *Union) Size() int32 {
	return union.size
}

func (union *Union) SetSize(size int32) {
	union.size = size
}

func (union *Union) Align() int32 {
	return union.align
}

func (union *Union) SetAlign(align int32) {
	union.align = align
}

func (union *Union) Kind() UnionKind {
	return union.kind
}

func (union *Union) SetKind(kind UnionKind) {
	union.kind = kind
}

type Type struct {
	line   Line
	size   int32
	align  int32
	record *Record
	list   *List
	map_   *Map
	union  *Union
	enum   *Enum
	stream *Stream
	kind   Kind
}

func (d *Type) Line() Line {
	return d.line
}

func (d *Type) SetLine(line Line) {
	d.line = line
}

func (d *Type) Kind() Kind {
	return d.kind
}
func (d *Type) SetKind(v Kind) *Type {
	d.kind = v
	return d
}
func (d *Type) Size() int32 {
	return d.size
}
func (d *Type) SetSize(v int32) *Type {
	d.size = v
	return d
}
func (d *Type) Align() int32 {
	return d.align
}
func (d *Type) SetAlign(v int32) *Type {
	d.align = v
	return d
}
func (d *Type) Record() *Record {
	return d.record
}
func (d *Type) SetRecord(v *Record) *Type {
	d.kind = KindRecord
	d.record = v
	return d
}
func (d *Type) List() *List {
	return d.list
}
func (d *Type) SetList(v *List) *Type {
	d.kind = KindList
	d.list = v
	return d
}
func (d *Type) Map() *Map {
	return d.map_
}
func (d *Type) SetMap(v *Map) *Type {
	d.kind = KindMap
	d.map_ = v
	return d
}
func (d *Type) Union() *Union {
	return d.union
}
func (d *Type) SetUnion(v *Union) *Type {
	d.kind = KindUnion
	d.union = v
	return d
}
func (d *Type) Stream() *Stream {
	return d.stream
}
func (d *Type) SetStream(v *Stream) *Type {
	d.kind = KindStream
	d.stream = v
	return d
}
func (d *Type) Enum() *Enum {
	return d.enum
}
func (d *Type) SetEnum(v *Enum) *Type {
	d.kind = KindEnum
	d.enum = v
	return d
}

// Field represents a field on a Record, List Element, or Map Key/Value
type Field struct {
	name         string
	short        string
	comments     []string
	kind         Type
	offset       int32
	columnOffset int32
	align        int32
	number       uint16
	optional     bool
	pointer      bool
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) SetName(name string) {
	f.name = name
}

func (f *Field) Short() string {
	return f.short
}

func (f *Field) SetShort(short string) {
	f.short = short
}

func (f *Field) Comments() []string {
	return f.comments
}

func (f *Field) SetComments(comments []string) {
	f.comments = comments
}

func (f *Field) Kind() *Type {
	return &f.kind
}

func (f *Field) SetKind(kind Type) {
	f.kind = kind
}

func (f *Field) Offset() int32 {
	return f.offset
}

func (f *Field) SetOffset(offset int32) {
	f.offset = offset
}

func (f *Field) ColumnOffset() int32 {
	return f.columnOffset
}

func (f *Field) SetColumnOffset(columnOffset int32) {
	f.columnOffset = columnOffset
}

func (f *Field) Align() int32 {
	return f.align
}

func (f *Field) SetAlign(align int32) {
	f.align = align
}

func (f *Field) Number() uint16 {
	return f.number
}

func (f *Field) SetNumber(number uint16) {
	f.number = number
}

func (f *Field) Optional() bool {
	return f.optional
}

func (f *Field) SetOptional(optional bool) {
	f.optional = optional
}

func (f *Field) Pointer() bool {
	return f.pointer
}

func (f *Field) SetPointer(pointer bool) {
	f.pointer = pointer
}

func (f *Field) Column(count int) (offset, size int) {
	return int(f.columnOffset) * count, int(f.Align()) * count
}

func (f *Field) IsPointer() bool {
	switch f.Kind().Kind() {
	case KindRecord:
		return true
	case KindString, KindBytes:
		return true
	case KindList:
		l := f.Kind().List()
		if l == nil {
			return false
		}
		return l.Fixed() > 0
	case KindMap:
		return true
	}
	return false
}

type Enum struct {
	name    string
	options []EnumOption
	kind    Kind
}

func (e *Enum) Name() string {
	return e.name
}

func (e *Enum) SetName(name string) {
	e.name = name
}

func (e *Enum) Options() []EnumOption {
	return e.options
}

func (e *Enum) SetOptions(options []EnumOption) {
	e.options = options
}

func (e *Enum) Kind() Kind {
	return e.kind
}

func (e *Enum) SetKind(kind Kind) {
	e.kind = kind
}

type EnumOption struct {
	name   string
	value  int64
	valueU uint64
}

func (e *EnumOption) Name() string {
	return e.name
}

func (e *EnumOption) SetName(name string) {
	e.name = name
}

func (e *EnumOption) Value() int64 {
	return e.value
}

func (e *EnumOption) SetValue(value int64) {
	e.value = value
}

func (e *EnumOption) ValueU() uint64 {
	return e.valueU
}

func (e *EnumOption) SetValueU(valueU uint64) {
	e.valueU = valueU
}

// Stream is a special type for Streaming. It's a list of Records that fits inside
// 1KB, 2KB, 4KB, 8KB, 16KB, 32KB or 64KB blocks.
type Stream struct {
	name   string
	kind   StreamKind
	record *Record
	layout BlockLayout // Row (Arrays of Structs) or Column (Struct of Arrays *fixed only)
}

func (s *Stream) Name() string {
	return s.name
}

func (s *Stream) SetName(name string) {
	s.name = name
}

func (s *Stream) Kind() StreamKind {
	return s.kind
}

func (s *Stream) SetKind(kind StreamKind) {
	s.kind = kind
}

func (s *Stream) Record() *Record {
	return s.record
}

func (s *Stream) SetRecord(record *Record) {
	s.record = record
}

func (s *Stream) Layout() BlockLayout {
	return s.layout
}

func (s *Stream) SetLayout(layout BlockLayout) {
	s.layout = layout
}

//type Block struct {
//	Record         *Record `json:"record"`
//	ID             string  `json:"id"`
//	IDField        *Field  `json:"-"`
//	Timestamp      string  `json:"timestamp"`
//	TimestampField *Field  `json:"-"`
//}

type BlockRecord struct {
}

type Schema struct {
	fqn       string
	name      string
	alias     string
	types     []Type
	typesMap  map[string]*Type
	imports   []Import
	importMap map[string]*Import
}

func (s *Schema) Fqn() string {
	return s.fqn
}

func (s *Schema) SetFqn(fqn string) {
	s.fqn = fqn
}

func (s *Schema) Name() string {
	return s.name
}

func (s *Schema) SetName(name string) {
	s.name = name
}

func (s *Schema) Alias() string {
	return s.alias
}

func (s *Schema) SetAlias(alias string) {
	s.alias = alias
}

func (s *Schema) Types() []Type {
	return s.types
}

func (s *Schema) SetTypes(types []Type) {
	s.types = types
}

func (s *Schema) TypesMap() map[string]*Type {
	return s.typesMap
}

func (s *Schema) SetTypesMap(typesMap map[string]*Type) {
	s.typesMap = typesMap
}

func (s *Schema) Layout() {
	for i := 0; i < len(s.types); i++ {
		t := &s.types[i]
		if t.record != nil {
			t.record.DoLayout()
		}
	}
}

type ImportGroup struct {
	comments []string
	imports  []*Import
}

type Import struct {
	path     string
	alias    string
	comments []string
	schema   Schema
}
