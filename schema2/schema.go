package schema2

type ID struct {
	kind      Kind
	record    *SRecord
	structure *SStruct
	enum      *SEnum
	union     *SUnion
}

type SRecord struct {
	Record
	Fields []SField
}

type SStruct struct {
	Struct
	Fields []SField
}

type SField struct {
	Field
}

type SEnum struct {
	Enum
	Options []SEnumOption
}

type SEnumOption struct{}

type SUnion struct {
	Union
	Options []UnionOption
}

type SUnionOption struct {
	UnionOption
}

type SSchema struct {
	Schema
	types   map[string]ID
	records map[string]*SRecord
}
