package schema

import "strings"

const (
	MapHeaderSize     = 4
	MapItemHeaderSize = 4
)

type Kind byte

const (
	KindUnknown = Kind(0) // Possibly resolved later
	KindBool    = Kind(1)
	KindByte    = Kind(2)
	KindInt8    = Kind(3)
	KindInt16   = Kind(4)
	KindUInt16  = Kind(5)
	KindInt32   = Kind(6)
	KindUInt32  = Kind(7)
	KindInt64   = Kind(8)
	KindUInt64  = Kind(9)
	KindFloat32 = Kind(10)
	KindFloat64 = Kind(11)
	KindString  = Kind(12)
	KindBytes   = Kind(13)
	KindStruct  = Kind(30) // User-defined structure
	KindEnum    = Kind(31) // User-defined enum
	KindUnion   = Kind(33) // User-defined union
	KindMessage = Kind(40) // User-defined message
	KindList    = Kind(50)
	KindMap     = Kind(60)
	KindPad     = Kind(100) // struct alignment padding
)

type ProtoBufKind byte

func (k Kind) Size() int {
	switch k {
	case KindBool, KindByte, KindInt8:
		return 1
	case KindInt16, KindUInt16:
		return 2
	case KindInt32, KindUInt32, KindFloat32:
		return 4
	case KindInt64, KindUInt64, KindFloat64:
		return 8
	}
	return -1
}

const (
	bitSlot0 byte = 1 << iota
	bitSlot1
	bitSlot2
	bitSlot3
	bitSlot4
	bitSlot5
	bitSlot6
	bitSlot7
)

var (
	emptyOptionals = make([]*Field, 0, 0)
)

type Nil struct{}
type ConstVal string

type Package struct {
	Name  string
	Files []*File
	Types map[string]*Type
}

type Annotation struct {
	Line  int
	Name  string
	Value interface{}
}

type Line struct {
	Number int
	Begin  int
	End    int
}

type Enum struct {
	Name      string
	Type      *Type
	Options   []*EnumOption
	optionMap map[string]*EnumOption
}

func (e *Enum) OptionMap() map[string]*EnumOption {
	if e.optionMap != nil {
		return e.optionMap
	}
	e.optionMap = make(map[string]*EnumOption)
	for _, o := range e.Options {
		e.optionMap[o.Name] = o
	}
	return e.optionMap
}

func (e *Enum) GetOption(name string) *EnumOption {
	return e.OptionMap()[name]
}

type EnumOption struct {
	Enum              *Enum
	Name              string
	Comments          []string
	Value             interface{}
	Line              Line
	Deprecated        bool
	DeprecatedMessage string
}

type Const struct {
	Name string
	Type *Type
}

type Message struct {
	Name      string
	Type      *Type
	Fields    []*MessageField
	FieldMap  map[string]*MessageField
	Optionals []*MessageField
	Version   int64
}

type MessageField struct {
	Number    int
	Owner     *Message
	Name      string
	Type      *Type
	Offset    int
	OptOffset int
	OptMask   byte
}

type Optional struct {
	Index  int
	Offset int
	Mask   byte
}

func (f *Field) IsOptional() bool {
	return f.Type != nil && f.Type.Optional
}

type List struct {
	Element *Type
}

type Map struct {
	Key      *Type
	Value    *Type
	ItemSize int
}

type Union struct {
	Name     string
	Comments []string
	Type     *Type
	Options  []*UnionOption
}

type UnionOption struct {
	Name     string
	Union    *Union
	Type     *Type
	Comments []string
}

type ImportedName struct {
	Package string
	Name    string
}

func KindOf(name string) Kind {
	if strings.Index(name, "string") == 0 {
		return KindString
	}
	if strings.Index(name, "bytes") == 0 {
		return KindBytes
	}

	switch name {
	case "i8", "int8":
		return KindInt8
	case "u8", "uint8", "byte":
		return KindByte
	case "i16", "int16", "short":
		return KindInt16
	case "u16", "uint16", "ushort":
		return KindUInt16
	case "i32", "int32", "int":
		return KindInt32
	case "u32", "uint32", "uint":
		return KindUInt32
	case "i64", "int64", "long":
		return KindInt64
	case "u64", "uint64", "ulong":
		return KindUInt64

	case "f32", "float32", "float":
		return KindFloat32
	case "f64", "float64", "double", "decimal":
		return KindFloat64

	case "bool", "boolean":
		return KindBool

	default:
		return KindUnknown
	}
}
