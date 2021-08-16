package schema

// Struct represent a fixed sized memory layout similar to how structs memory layout
// in languages such as Go, Rust, C/C++, etc. Optionally structs can be compact which
// removes all padding which favors memory size vs CPU cache aligning. For variable
// memory sizes, use Record.
type Struct struct {
	Name      string
	Type      *Type
	Fields    []*StructField
	FieldMap  map[string]*StructField
	Optionals []*StructField
	Version   int64
	Compact   bool
}

type StructField struct {
	Number    int
	Struct    *Struct
	Name      string
	Type      *Type
	Offset    int
	OptOffset int
	OptMask   byte
}

func (st *Struct) setOptionals() {
	if st.Optionals != nil {
		return
	}
	// Setup optionals
	for _, field := range st.Fields {
		if field.Type.Optional {
			if len(st.Optionals) == 0 {
				field.OptOffset = 0
				field.OptMask = bitSlot0
			} else {
				field.OptOffset = len(st.Optionals) / 8
				switch len(st.Optionals) % 8 {
				case 0:
					field.OptMask = bitSlot0
				case 1:
					field.OptMask = bitSlot1
				case 2:
					field.OptMask = bitSlot2
				case 3:
					field.OptMask = bitSlot3
				case 4:
					field.OptMask = bitSlot4
				case 5:
					field.OptMask = bitSlot5
				case 6:
					field.OptMask = bitSlot6
				case 7:
					field.OptMask = bitSlot7
				}
			}
			st.Optionals = append(st.Optionals, field)
		}
	}

	if st.Optionals == nil {
		st.Optionals = emptyOptionals
	}
}
