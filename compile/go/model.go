package _go

import (
	"encoding/binary"
	. "github.com/moontrade/proto/schema"
)

const headerFieldName = "_h_"

func goFileName(order binary.ByteOrder) string {
	if order == binary.BigEndian {
		return "proto_be.go"
	}
	return "proto.go"
}

// Config provides configuration for the Compiler
type Config struct {
	BigEndian bool
	Fluent    bool
	NoGoFmt   bool

	Package string
	// By default a separate Mutable type is created for each struct
	// This flag sets that only one struct that is both read-only and mutable is generated
	Mutable       bool
	MultipleFiles bool
	Output        string
}

// Compiler generates Go code for a supplied Schema
type Compiler struct {
	schema   *Schema
	config   *Config
	packages map[string]*goPackage
}

type goType struct {
	pkg       *goPackage
	t         *Type
	name      string    // type name
	mut       string    // Mutable type name
	primitive bool      // golang primitive
	imp       *goImport // associated import
	cst       *goConst  // const
	enum      *goEnum   // enum
	st        *goStruct // struct
	list      *goList   // list
}

type goList struct {
	element   *goType
	name      string
	sliceName string
}

type goPackage struct {
	file        *File
	importAlias string
	packageName string
	dir         string
	path        string
	imports     []string
	aliasCount  int
	byType      map[*Type]*goType
	importMap   map[string]*goImport
	types       map[string]*goType
	lists       map[string]*goType
	strings     map[string]*goType
	structs     map[string]*goType
	enums       map[string]*goType
	unions      map[string]*goType
	names       map[string]struct{}
}

type goImport struct {
	path     string
	alias    string
	useAlias bool
}

type goFile struct {
	imports []*goImport
	types   []*goType
}

type goTypeAlias struct {
	t    *goType
	name string
}

type goConst struct {
	cst  *Const
	name string
}

type goEnum struct {
	enum    *Enum
	value   *goType
	options []*goEnumOption
}

type goEnumOption struct {
	option *EnumOption
	name   string
}

type goStruct struct {
	st     *Struct
	fields []*goField
}

type goField struct {
	field     *StructField
	isPointer bool
	public    string // Name of public accessor
	private   string // Name of field if declared inside a struct
	t         *goType
}
