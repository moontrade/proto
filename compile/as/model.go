package as

import . "github.com/moontrade/proto/schema"

const TSFileName = "index.ts"
const TSFileNameSuffix = "_buf.ts"

// Configuration for the AssemblyScript/TypeScript code generator
type ASConfig struct {
	BigEndianSafe bool
	Fluent        bool

	Package       string
	GlobalPackage string
	// By default a separate Mutable type is created for each struct
	// This flag sets that only one struct that is both read-only and mutable is generated
	Mutable       bool
	MultipleFiles bool
	Output        string
}

// Generates Go code
type Compiler struct {
	schema   *Schema
	config   *ASConfig
	packages map[string]*asPackage
}

type asType struct {
	pkg       *asPackage
	t         *Type
	name      string
	mut       string
	primitive bool
	imp       *asImport
	cst       *asConst
	enum      *asEnum
	st        *asStruct
	list      *asList
}

type asList struct {
	element   *asType
	name      string
	sliceName string
}

type asPackage struct {
	file *File

	importAlias string
	packageName string
	dir         string
	path        string
	imports     []string
	aliasCount  int

	byType    map[*Type]*asType
	importMap map[string]*asImport
	types     map[string]*asType
	lists     map[string]*asType
	strings   map[string]*asType
	structs   map[string]*asType
	enums     map[string]*asType
	unions    map[string]*asType
	names     map[string]struct{}
}

type asImport struct {
	path     string
	alias    string
	useAlias bool
}

type asFile struct {
	imports []*asImport
	types   []*asType
}

type asTypeAlias struct {
	t    *asType
	name string
}

type asConst struct {
	cst  *Const
	name string
}

type asEnum struct {
	enum    *Enum
	value   *asType
	options []*asEnumOption
}

type asEnumOption struct {
	option *EnumOption
	name   string
}

type asStruct struct {
	st     *Struct
	fields []*asField
}

type asField struct {
	field     *StructField
	isPointer bool
	public    string // Name of public accessor
	private   string // Name of field if declared inside a struct
	t         *asType
}
