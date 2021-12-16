package runtime2

import (
	"github.com/moontrade/nogc"
)

type Kind byte

const (
	KindUnknown     = Kind(0) // Possibly resolved later
	KindBool        = Kind(1)
	KindByte        = Kind(2)
	KindInt8        = Kind(3)
	KindInt16       = Kind(4)
	KindUInt16      = Kind(5)
	KindInt32       = Kind(6)
	KindUInt32      = Kind(7)
	KindInt64       = Kind(8)
	KindUInt64      = Kind(9)
	KindFloat32     = Kind(10)
	KindFloat64     = Kind(11)
	KindString      = Kind(12)
	KindStringFixed = Kind(13)
	KindBytes       = Kind(14)
	KindFixed       = Kind(15)
	KindTimeRange   = Kind(16)
	KindEnum        = Kind(20) // User-defined enum
	KindRecord      = Kind(30) // User-defined record
	KindStruct      = Kind(31) // User-defined fixed sized record
	KindList        = Kind(40)
	KindMap         = Kind(50)
	KindUnion       = Kind(60)  // User-defined union
	KindPad         = Kind(100) // struct alignment padding
)

const (
	VPointerSize = int32(2)
)

func Int64Field(name string, offset int32) Field {
	return Field{Name: name, Offset: offset, Size: 8, Kind: KindInt64}
}
func Int64FieldNamed(name, short string, offset int32) Field {
	return Field{Name: name, Offset: offset, Size: 8, Kind: KindInt64}
}
func UInt64Field(name string, offset int32) Field {
	return Field{Name: name, Offset: offset, Size: 8, Kind: KindUInt64}
}
func StringField(name string, offset int32, pointerSize int32) Field {
	return Field{Name: name, Offset: offset, Size: pointerSize, Kind: KindString}
}
func StringElement() Field {
	return Field{Size: 8, Kind: KindString}
}
func StringFixedField(name string, offset, max int32) Field {
	return Field{Name: name, Offset: offset, Size: max, Kind: KindStringFixed}
}
func StringFixedElement(name string, offset, max int32) Field {
	return Field{Offset: 0, Size: max, Kind: KindString}
}
func ListField(name string, offset int32, element Field) Field {
	return Field{
		Name:   name,
		Offset: offset,
		Size:   8,
		Kind:   KindList,
		List: &List{
			Element: element,
		},
	}
}
func ListElement(element Field) Field {
	return Field{
		Size: 8,
		Kind: KindList,
		List: &List{
			Element: element,
		},
	}
}
func ListFixedField(name string, offset, max int32, element Field) Field {
	return Field{
		Name:   name,
		Offset: offset, Size: max * element.Size, Kind: KindList,
		List: &List{
			Element: element,
		},
	}
}
func ListFixedElement(max int32, element Field) Field {
	return Field{
		Size: max * element.Size, Kind: KindList,
		List: &List{
			Element: element,
		},
	}
}

type LogRecordHeader interface {
	RecordID() int64

	Timestamp() int64

	Size() int
}

type SeriesRecordHeader interface {
	LogRecordHeader

	Start() int64

	End() int64
}

type RecordLayout int32

const (
	RecordLayoutCompact      RecordLayout = 0
	RecordLayoutCacheAligned RecordLayout = 1
)

type Record struct {
	Name         string            `json:"name"`
	Comments     []string          `json:"comments"`
	Fields       []Field           `json:"fields"`
	Header       []int             `json:"header,omitempty"`
	HeaderFields []*Field          `json:"-"`
	FieldsMap    map[string]*Field `json:"-"`
	Version      int64             `json:"version"`
	Size         int32             `json:"size"`
	Layout       RecordLayout      `json:"layout,omitempty"`
	Flex         bool              `json:"flex"` // Does Record have any variable fields?
}

func (r *Record) Field(name string) *Field {
	m := r.FieldsMap
	if m == nil {
		m = make(map[string]*Field)
	}
	if len(m) == 0 {
		for i := 0; i < len(r.Fields); i++ {
			field := &r.Fields[i]
			m[field.Name] = field
			if len(field.CompactName) > 0 {
				m[field.CompactName] = field
			}
		}
		r.FieldsMap = m
	}
	return m[name]
}

func (r *Record) FieldAt(index int) *Field {
	if index < 0 || index >= len(r.Fields) {
		return nil
	}
	return &r.Fields[index]
}

type Bytes struct {
	Fixed int
}

// List represents an Array like structure.
// Header: | LEN 2 bytes | [LEN]Element
type List struct {
	Element Field `json:"element"`
	Fixed   int   `json:"fixed"`
}

// Map represents a HashMap data structure using a robin-hood algorithm
// Header: | LEN 2 bytes | SIZE 2 bytes | List<MapEntry>
// Item: | KEY | VALUE | Distance (2 bytes)
type Map struct {
	Key     Field        `json:"key"`
	Value   Field        `json:"value"`
	Default nogc.Pointer `json:"-"`
}

// Union represents a C-like union or a protobuf oneOf
type Union struct {
	Options []Field
}

// Field represents a field on a Record, List Element, or Map Key/Value
type Field struct {
	Name        string   `json:"name"`
	CompactName string   `json:"compact"`
	Comments    []string `json:"comments"`
	Record      *Record  `json:"-"` // oneof
	List        *List    `json:"-"` // oneof
	Map         *Map     `json:"-"` // oneof
	Union       *Union   `json:"-"` // oneof
	Enum        *Enum    `json:"-"` // oneof
	Offset      int32    `json:"offset"`
	Size        int32    `json:"size"`  // Number of bytes
	Align       int32    `json:"align"` // Number of bytes
	Number      uint16   `json:"number"`
	Kind        Kind     `json:"kind"`
	Optional    bool     `json:"optional"`
	Pointer     bool     `json:"pointer"`
}

func (f *Field) IsPointer() bool {
	switch f.Kind {
	case KindRecord:
		r := f.Record
		if r == nil {
			return false
		}
		return r.Flex
	case KindString, KindBytes:
		return true
	case KindList:
		l := f.List
		if l == nil {
			return false
		}
		return l.Fixed > 0
	}
	return false
}

//func (f *Record) Value(p Pointer) Pointer {
//	if f.Flex {
//		offset := p.VPointer(f.Offset)
//		return p.Slice(int(offset.Pos), int(offset.Len))
//	}
//	return p.Slice(f.Offset, f.Size)
//}

func FixedStringLengthBytes(size int) int {
	switch {
	case size < 255:
		return 1
	case size < 65534:
		return 2
	default:
		return 4
	}
}

type Enum struct {
	Name    string       `json:"name"`
	Options []EnumOption `json:"options"`
	Kind    Kind         `json:"kind"`
}

type EnumOption struct {
	Name   string `json:"name"`
	Value  int64  `json:"value"`
	ValueU uint64 `json:"valueU,omitempty"`
}

type Schema struct {
	FQN        string             `json:"fqn"`
	Name       string             `json:"name"`
	Records    []Record           `json:"records"`
	RecordsMap map[string]*Record `json:"-"`
	Lists      []List             `json:"lists"`
	Streams    []Stream           `json:"streams"`
}

type Import struct {
	FQN  string `json:"fqn"`
	Name string `json:"name"`
}

type BlockLayout int32

const (
	BlockLayoutRow    BlockLayout = 0
	BlockLayoutColumn BlockLayout = 1
)

// Stream is a special type for Streaming. It's a list of Records that fits inside
// 1KB, 2KB, 4KB, 8KB, 16KB, 32KB or 64KB blocks.
type Stream struct {
	Name       string      `json:"name"`
	RecordName string      `json:"record"`
	Record     *Record     `json:"-"`
	Layout     BlockLayout `json:"layout"` // Row (Arrays of Structs) or Column (Struct of Arrays *fixed only)
}

type BlockRecord struct {
}

func TimeRangeRecord() Record {
	return Record{
		Name: "Bar",
		Fields: []Field{
			Int64FieldNamed("start", "s", 0),
			Int64FieldNamed("end", "e", 8),
		},
	}
}

func TimeRangeHeaderRecord() Record {
	return Record{
		Name: "Bar",
		Fields: []Field{
			Int64FieldNamed("id", "i", 0),
			Int64FieldNamed("start", "s", 0),
			Int64FieldNamed("end", "e", 8),
		},
	}
}

func init() {
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
}
