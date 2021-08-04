package schema

type Type struct {
	Line         Line
	File         *File
	Kind         Kind
	Optional     bool
	Resolved     bool
	Size         int
	Len          int // Max length if collection (list or map) or string
	HeaderSize   int
	HeaderOffset int
	Padding      int
	Element      *Type // List element or Map key type
	Value        *Type // Value type if map type
	ItemSize     int
	Import       *Import
	Name         string       // Name of type
	Comments     []string     // Comments
	Description  string       // Description are comments to the right of certain declarations
	Const        *Const       // Const if type represents a single const
	Struct       *Struct      // Struct for 'KindStruct'
	Field        *Field       // Field
	Union        *Union       // Union for 'KindUnion'
	UnionOption  *UnionOption // UnionOption if type represents a single union option
	Enum         *Enum        // Enum for 'KindEnum'
	EnumOption   *EnumOption  // EnumOption if type represents a single enum option
	Init         interface{}  // Initial value
}

func (t *Type) Base() *Type {
	switch {
	case t.Struct != nil:
		return t.Struct.Type
	case t.Enum != nil:
		return t.Enum.Type
	case t.Union != nil:
		return t.Union.Type
	}
	return t
}
