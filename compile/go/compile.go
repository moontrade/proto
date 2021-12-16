package _go

import (
	"encoding/binary"
	"errors"
	"fmt"
	. "github.com/moontrade/proto/schema"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func NewCompiler(schema *Schema, config *Config) (*Compiler, error) {
	return &Compiler{
		schema:   schema,
		config:   config,
		packages: make(map[string]*goPackage),
	}, nil
}

func (c *Compiler) Compile() error {
	// Convert into Go specific model
	packages := make(map[string]*goPackage)
	var (
		err    error
		output string
		info   os.FileInfo
	)
	for k, v := range c.schema.Files {
		packages[k], err = c.createPackage(v, 0)
		if err != nil {
			return err
		}
	}

	//path := c.config.Output
	//path, err = filepath.Abs(path)
	//output := RelativePath(path, c.config.Output)
	//output := RelativePath(path, c.config.Output)
	output, err = filepath.Abs(c.config.Output)
	if err != nil {
		return err
	}
	_ = os.MkdirAll(output, 0755)
	info, err = os.Stat(output)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("output directory '%s' is not a directory", output)
	}
	//output := filepath.Join(path, c.config.Output)
	createFile := func(f *goPackage, order binary.ByteOrder) error {
		b := NewBuilder()
		if err := c.writeFile(f, b, order); err != nil {
			return err
		}

		// Compute dir path
		dir := filepath.Join(output, f.path)

		// mkdir
		_ = os.MkdirAll(dir, 0755)

		name, err := filepath.Abs(filepath.Join(dir, goFileName(order)))
		if err != nil {
			return err
		}
		// Create new "moon.go"
		out, err := os.Create(name)
		if err != nil {
			return err
		}
		_, err = out.WriteString(b.String())
		if err != nil {
			return err
		}

		if err = out.Sync(); err != nil {
			return err
		}
		if err = out.Close(); err != nil {
			return err
		}

		if !c.config.NoGoFmt {
			err = exec.Command("gofmt", "-w", name).Run()
			if err != nil {

			}
		}
		return nil
	}
	err = filepath.Walk(output, c.walkClear)
	for _, f := range packages {
		if err = createFile(f, binary.LittleEndian); err != nil {
			return err
		}
		if c.config.BigEndian {
			if err = createFile(f, binary.BigEndian); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Compiler) walkClear(path string, info fs.FileInfo, err error) error {
	if strings.HasSuffix(path, "proto_be.go") {
		_ = os.Remove(path)
	} else if strings.HasSuffix(path, "proto.go") {
		_ = os.Remove(path)
	}
	return nil
}

func (c *Compiler) goName(n string, types map[string]*Type) string {
	f := Capitalize(n)
	//ctorName := fmt.Sprintf("New%s", f)
	//if types[ctorName] != nil {
	//	f = f + "_"
	//}
	//if types["Unmarshal" + f] != nil {
	//	f = f + "_"
	//}
	return f
}

func (c *Compiler) primitive(file *goPackage, t *Type, goName string) *goType {
	return &goType{
		pkg:       file,
		t:         t,
		name:      goName,
		mut:       goName,
		primitive: true,
	}
}

func (c *Compiler) stringType(file *goPackage, t *Type) *goType {
	prefix := ""
	if t.Kind == KindString {
		prefix = "String"
	} else {
		prefix = "Bytes"
	}
	name := fmt.Sprintf("%s%d", prefix, t.Len)
	mut := fmt.Sprintf("%sMut", name)
	return &goType{
		pkg:       file,
		t:         t,
		name:      name,
		mut:       mut,
		primitive: true,
	}
}

func joinWithSlash(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return filepath.Clean(strings.Join(elem[i:], "/"))
		}
	}
	return ""
}

func (c *Compiler) toGoPath(f *File) string {
	return filepath.Join(c.config.Package, f.Dir)
}

const (
	PathSeparator = string(os.PathSeparator)
)

func (c *Compiler) createPackage(file *File, level int) (*goPackage, error) {
	if level >= 20 {
		return nil, fmt.Errorf("cyclic package dependency: %s", file.Path)
	}
	if existing := c.packages[file.Dir]; existing != nil {
		return existing, nil
	}
	//packageParts := strings.Split(file.Path, PathSeparator)
	path := file.Dir //filepath.Join(packageParts...)
	pkg := &goPackage{
		file:        file,
		path:        path,
		dir:         filepath.Join(c.config.Output, path),
		packageName: file.Package,
		byType:      make(map[*Type]*goType),
		importMap:   make(map[string]*goImport),
		types:       make(map[string]*goType),
		structs:     make(map[string]*goType),
		strings:     make(map[string]*goType),
		enums:       make(map[string]*goType),
		lists:       make(map[string]*goType),
		unions:      make(map[string]*goType),
		names:       make(map[string]struct{}),
	}
	c.packages[file.Dir] = pkg

	for k, t := range file.Types {
		n := Capitalize(k)
		pkg.names[n] = struct{}{}
		if t.Enum != nil {
			for _, option := range t.Enum.Options {
				pkg.names[c.enumOptionName(option)] = struct{}{}
			}
		} else if t.Struct != nil {
			pkg.names["Reinterpret"+n] = struct{}{}
			pkg.names["Unmarshal"+n] = struct{}{}
		}
	}
	for _, value := range file.Types {
		_, err := c.resolve(pkg, value, level+1)
		if err != nil {
			return nil, err
		}
	}

	imps := make([]string, 0, len(pkg.importMap))
	for k := range pkg.importMap {
		imps = append(imps, k)
	}
	sort.Strings(imps)

	for _, key := range imps {
		imp := pkg.importMap[key]
		if imp.useAlias && len(imp.alias) > 0 {
			pkg.imports = append(pkg.imports, fmt.Sprintf("%s \"%s\"", imp.alias, key))
		} else {
			pkg.imports = append(pkg.imports, fmt.Sprintf("\"%s\"", imp.path))
		}
	}

	return pkg, nil
}

func (c *Compiler) addImport(imports map[string]*goImport, path, alias string) *goImport {
	if existing := imports[path]; existing != nil {
		return existing
	}
	lastIndex := strings.LastIndexByte(path, '/')
	useAlias := true
	if lastIndex > -1 {
		name := path[lastIndex+1:]
		if name == alias {
			useAlias = false
		}
	} else {
		alias = path
		useAlias = false
	}

	if len(alias) > 0 {
		for existing := imports[alias]; existing != nil; {
			alias = "_" + alias
			useAlias = true
		}
	}
	imp := &goImport{
		path:     path,
		alias:    alias,
		useAlias: useAlias,
	}
	imports[path] = imp
	return imp
}

func (c *Compiler) littleEndianFlags(b *Builder) {
	b.W("// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm")
	b.W("")
}

func (c *Compiler) bigEndianFlags(b *Builder) {
	b.W("// +build ppc64 s390x mips mips64")
	b.W("")
}

func (c *Compiler) writeFile(file *goPackage, b *Builder, order binary.ByteOrder) error {
	W := b.W

	switch order {
	case binary.LittleEndian:
		c.littleEndianFlags(b)
	case binary.BigEndian:
		c.bigEndianFlags(b)
	}

	W("package %s\n", file.packageName)

	if len(file.imports) > 0 {
		W("import (")
		for _, imp := range file.imports {
			W("    %s", imp)
		}
		W(")\n")
	}

	for _, enum := range file.enums {
		if err := c.genEnum(file, enum, b); err != nil {
			return err
		}
	}

	init := NewBuilder()
	init.W("func init() {")

	if len(file.structs) > 0 {
		init.W(`    {
		var b [2]byte
        v := uint16(1)
        b[0] = byte(v)
        b[1] = byte(v >> 8)
		if *(*uint16)(unsafe.Pointer(&b[0])) != 1 {
			panic("BigEndian not supported")
		}
	}`)
		_, _ = init.Write([]byte(`    type b struct {
        n string
        o, s uintptr
    }
    a := func(x interface{}, y interface{}, s uintptr, z []b) {
        t := reflect.TypeOf(x)
        r := reflect.TypeOf(y)
        if t.Size() != s {
            panic(fmt.Sprintf("sizeof %%s = %%d, expected = %%d", t.Name(), t.Size(), s))
        }
        if r.Size() != s {
            panic(fmt.Sprintf("sizeof %%s = %%d, expected = %%d", r.Name(), r.Size(), s))
        }
        if t.NumField() != len(z) {
            panic(fmt.Sprintf("%%s field count = %%d: expected %%d", t.Name(), t.NumField(), len(z)))
        }
        for i, e := range z {
            f := t.Field(i)
            if f.Offset != e.o {
                panic(fmt.Sprintf("%%s.%%s offset = %%d, expected = %%d", t.Name(), f.Name, f.Offset, e.o))
            }
            if f.Type.Size() != e.s {
                panic(fmt.Sprintf("%%s.%%s size = %%d, expected = %%d", t.Name(), f.Name, f.Type.Size(), e.s))
            }
            if f.Name != e.n {
                panic(fmt.Sprintf("%%s.%%s expected field: %%s", t.Name(), f.Name, e.n))
            }
        }
    }`))
		init.W("")
		init.W("")

		for _, st := range file.structs {
			if err := c.genStruct(file, st, false, b, order); err != nil {
				return err
			}
			if err := c.genStruct(file, st, true, b, order); err != nil {
				return err
			}

			init.W("    a(%s{}, %s{}, %d, []b{", st.name, st.mut, st.t.Size)
			if st.t.HeaderSize > 0 {
				init.W("        {\"%s\", %d, %d},", headerFieldName, 0, st.t.HeaderSize)
			}
			for _, field := range st.st.fields {
				init.W("        {\"%s\", %d, %d},", field.private, field.field.Offset, field.t.t.Size)
			}
			init.W("    })")
		}
	}

	//for _, st := range file.structs {
	//	if err := c.genStructBytes(file, st, false, b); err != nil {
	//		return err
	//	}
	//	if err := c.genStructBytes(file, st, true, b); err != nil {
	//		return err
	//	}
	//}

	for _, str := range file.strings {
		if err := c.genString(str, false, b, order); err != nil {
			return err
		}
		if err := c.genString(str, true, b, order); err != nil {
			return err
		}
	}

	for _, list := range file.lists {
		if err := c.genArrayList(list, false, b, order); err != nil {
			return err
		}
		if err := c.genArrayList(list, true, b, order); err != nil {
			return err
		}
	}

	initStr := init.String()
	if initStr != "func init() {\n" {
		W(initStr)
		W("}\n")
	}

	return nil
}

func (c *goPackage) uniqueName(n string) string {
	for {
		if _, ok := c.names[n]; ok {
			n = n + "_"
		} else {
			break
		}
	}
	c.names[n] = struct{}{}
	return n
}

func (c *Compiler) resolve(pkg *goPackage, t *Type, level int) (*goType, error) {
	if existing := pkg.byType[t]; existing != nil {
		return existing, nil
	}
	if t.Import != nil {
		imp := c.addImport(
			pkg.importMap,
			c.toGoPath(t.Import.File),
			t.Import.Alias,
		)

		importPkg, err := c.createPackage(t.Import.File, level+1)
		if err != nil {
			return nil, err
		}

		base := t.Base()
		importedType := importPkg.types[base.Name]
		if importedType == nil {
			return nil, fmt.Errorf("'%s' could not resolve type '%s' in '%s'", t.File.Path, base.Name, t.File.Path)
		}

		// Handle imported type
		return &goType{
			pkg:       pkg,
			t:         t,
			name:      fmt.Sprintf("%s.%s", imp.alias, importedType.name),
			mut:       fmt.Sprintf("%s.%s", imp.alias, importedType.mut),
			primitive: false,
			imp:       imp,
			cst:       nil,
			enum:      nil,
			st:        nil,
			list:      nil,
		}, nil
	}

	switch t.Kind {
	case KindStruct:
		structName := Capitalize(t.Struct.Name)
		if existing := pkg.structs[structName]; existing != nil {
			ret := &goType{}
			*ret = *existing
			ret.t = t
			pkg.byType[t] = ret
			return existing, nil
		}
		_ = c.addImport(pkg.importMap, "fmt", "")
		_ = c.addImport(pkg.importMap, "io", "")
		_ = c.addImport(pkg.importMap, "reflect", "")
		_ = c.addImport(pkg.importMap, "unsafe", "")

		fields := make([]*goField, 0, len(t.Struct.Fields))
		names := make(map[string]struct{})
		createFieldName := func(n string) string {
			for _, ok := names[n]; ok; {
				n = n + "_"
			}
			names[n] = struct{}{}
			return n
		}
		for _, field := range t.Struct.Fields {
			fieldType, err := c.resolve(pkg, field.Type, level+1)
			if err != nil {
				return nil, err
			}
			fieldName := ""
			if fieldType.t.Kind == KindPad {
				fieldName = "_"
			} else {
				fieldName = createFieldName(c.fieldName(field.Name))
			}
			if field.Type.Kind == KindPad {
				fieldName = "_"
			}
			structField := Uncapitalize(fieldName)
			fields = append(fields, &goField{
				field:     field,
				isPointer: c.isPointerType(field.Type),
				public:    fieldName,
				private:   structField,
				t:         fieldType,
			})
			pkg.byType[field.Type] = fieldType
		}
		gt := &goType{
			pkg:       pkg,
			t:         t,
			name:      structName,
			primitive: false,
			st: &goStruct{
				st:     t.Struct,
				fields: fields,
			},
		}
		gt.mut = pkg.uniqueName(fmt.Sprintf("%sMut", gt.name))
		pkg.types[gt.name] = gt
		if pkg.structs == nil {
			pkg.structs = make(map[string]*goType)
		}
		pkg.structs[gt.name] = gt
		pkg.byType[t] = gt
		return gt, nil

	case KindEnum:
		enumName := Capitalize(t.Enum.Name)
		if existing := pkg.enums[enumName]; existing != nil {
			ret := &goType{}
			*ret = *existing
			ret.t = t
			pkg.byType[t] = existing
			return existing, nil
		}
		value, err := c.resolve(pkg, t.Element, level+1)
		if err != nil {
			return nil, err
		}
		enum := &goEnum{
			enum:    t.Enum,
			value:   value,
			options: make([]*goEnumOption, 0, len(t.Enum.Options)),
		}
		for _, option := range t.Enum.Options {
			o := &goEnumOption{
				option: option,
				name:   c.enumOptionName(option),
			}
			enum.options = append(enum.options, o)
		}
		gt := &goType{
			pkg:  pkg,
			t:    t,
			name: enumName,
			enum: enum,
		}
		pkg.types[gt.name] = gt
		if pkg.enums == nil {
			pkg.enums = make(map[string]*goType)
		}
		pkg.enums[gt.name] = gt
		pkg.byType[t] = gt
		return gt, nil

	case KindList:
		_ = c.addImport(pkg.importMap, "fmt", "")
		_ = c.addImport(pkg.importMap, "reflect", "")
		_ = c.addImport(pkg.importMap, "unsafe", "")
		element, err := c.resolve(pkg, t.Element, level+1)
		if err != nil {
			return nil, err
		}
		gt := &goType{
			pkg:       pkg,
			t:         t,
			name:      Capitalize(t.Name),
			primitive: false,
			list: &goList{
				element: element,
			},
		}
		gt.mut = pkg.uniqueName(fmt.Sprintf("%sMut", gt.name))
		pkg.types[gt.name] = gt
		if pkg.lists == nil {
			pkg.lists = make(map[string]*goType)
		}
		pkg.lists[gt.name] = gt
		pkg.byType[t] = gt
		return gt, nil

	case KindUnion:
		// TODO:
		return nil, fmt.Errorf("unions not supported yet: %s:%d %s", t.File.Path, t.Line.Number, t.Name)

	case KindBool:
		return c.primitive(pkg, t, "bool"), nil
	case KindByte:
		return c.primitive(pkg, t, "byte"), nil
	case KindInt8:
		return c.primitive(pkg, t, "int8"), nil
	case KindInt16:
		return c.primitive(pkg, t, "int16"), nil
	case KindUInt16:
		return c.primitive(pkg, t, "uint16"), nil
	case KindInt32:
		return c.primitive(pkg, t, "int32"), nil
	case KindUInt32:
		return c.primitive(pkg, t, "uint32"), nil
	case KindInt64:
		return c.primitive(pkg, t, "int64"), nil
	case KindUInt64:
		return c.primitive(pkg, t, "uint64"), nil
	case KindFloat32:
		return c.primitive(pkg, t, "float32"), nil
	case KindFloat64:
		return c.primitive(pkg, t, "float64"), nil
	case KindString, KindBytes:
		gt := c.stringType(pkg, t)
		pkg.strings[gt.name] = gt
		pkg.byType[t] = gt
		return gt, nil

	case KindPad:
		return &goType{
			pkg:       pkg,
			t:         t,
			name:      "_",
			mut:       "_",
			primitive: true,
			imp:       nil,
			cst:       nil,
			enum:      nil,
			st:        nil,
			list:      nil,
		}, nil
	}
	return nil, fmt.Errorf("type not supported yet: %s:%d %s", t.File.Path, t.Line.Number, t.Name)
}

func (c *Compiler) isPointerType(t *Type) bool {
	switch t.Kind {
	case KindString, KindStruct, KindUnion, KindList, KindMap:
		return true
	default:
		return false
	}
}

func (c *Compiler) genComments(b *Builder, comments []string) {
	if len(comments) == 0 {
		return
	}
	for _, comment := range comments {
		b.W("//%s", comment)
	}
}

func (c *Compiler) fieldName(f string) string {
	f = Capitalize(f)
	switch f {
	case "Mut":
		return "Mut_"
	case "Write":
		return "Write_"
	case "WriteTo":
		return "WriteTo_"
	case "Read":
		return "Read_"
	case "ReadFrom":
		return "ReadFrom_"
	case "MarshalBinary":
		return "MarshalBinary_"
	case "UnmarshalBinary":
		return "UnmarshalBinary_"
	case "MarshalBinaryTo":
		return "MarshalBinaryTo_"
	case "Reset":
		return "Reset_"
	case "MarshalTo":
		return "MarshalTo_"
	case "Bytes":
		return "Bytes_"
	case "Copy":
		return "Copy_"
	case "MarshalMap":
		return "MarshalMap_"
	case "String":
		return "String_"
	}
	return f
}

func (c *Compiler) enumOptionName(option *EnumOption) string {
	return fmt.Sprintf("%s_%s", Capitalize(option.Enum.Name), option.Name)
}

func (c *Compiler) unionOptionName(option *UnionOption) string {
	return fmt.Sprintf("%s_%s", Capitalize(option.Union.Name), Capitalize(option.Name))
}

func (c *Compiler) goTypeName(t *Type) string {
	switch t.Kind {
	case KindBool:
		return "bool"
	case KindByte:
		return "byte"
	case KindInt8:
		return "int8"
	case KindInt16:
		return "int16"
	case KindUInt16:
		return "uint16"
	case KindInt32:
		return "int32"
	case KindUInt32:
		return "uint32"
	case KindInt64:
		return "int64"
	case KindUInt64:
		return "uint64"
	case KindFloat32:
		return "float32"
	case KindFloat64:
		return "float64"
	case KindString, KindBytes:
		return "string"
	case KindEnum:
		return Capitalize(t.Enum.Name)
	case KindStruct:
		return Capitalize(t.Struct.Name)
	case KindUnion:
		return Capitalize(t.Union.Name)
	case KindList:
		return fmt.Sprintf("[%d]%s", t.Len, c.goTypeName(t.Element))
	case KindPad:
		return fmt.Sprintf("[%d]byte", t.Size)
	}
	return "unknown"
}

func (c *Compiler) genEnum(file *goPackage, t *goType, b *Builder) error {
	c.genComments(b, t.t.Comments)
	b.W("type %s %s\n", t.name, t.enum.value.name)

	b.W("const (")
	for _, option := range t.enum.options {
		if len(option.option.Comments) > 0 {
			c.genComments(b, option.option.Comments)
		}
		b.W("    %s = %s(%d)", option.name, t.name, option.option.Value)
		//if i < len(t.enum.options)-1 {
		//	b.W("")
		//}
	}
	b.W(")\n")
	return nil
}

func (c *Compiler) genStructBytes(file *goPackage, t *goType, mut bool, b *Builder, order binary.ByteOrder) error {
	//_ = c.genStruct(file, t, mut, b)
	st := t.st
	c.genComments(b, t.t.Comments)
	name := t.name
	if mut {
		name = t.mut
	}

	if mut {
		b.W("type %s struct {", name)
		b.W("    %s", t.name)
		b.W("}")
	} else {
		b.W("type %s [%d]byte\n", name, t.t.Size)
	}

	if file.types["New"+name] != nil {
		b.W("func New%s_() *%s {", name, name)
	} else {
		b.W("func New%s() *%s {", name, name)
	}
	b.W("    return &%s{}", name)
	b.W("}\n")

	//b.W("func (s *%s) Hash() uint32 {", name)
	//b.W("    return bu.Hash32Bytes(s[0:])")
	//b.W("}")
	//
	//b.W("func (s *%s) Hash64() uint64 {", name)
	//b.W("    return bu.Hash64Bytes(s[0:])")
	//b.W("}")

	b.W("func (s *%s) Copy(v *%s) {", name, name)
	b.W("    if v == nil {")
	b.W("        *s = %s{}", name)
	b.W("    } else {")
	b.W("        *s = *v")
	b.W("    }")
	b.W("}\n")

	b.W("func (s *%s) Clone() *%s {", name, name)
	b.W("    v := &%s{}", name)
	b.W("    *v = *s")
	b.W("    return v")
	b.W("}\n")

	if strings.HasSuffix(t.mut, "_") {
		t.mut = t.mut[0 : len(t.mut)-2]
	}

	if !mut {
		b.W("func (s *%s) String() string {", name)
		b.W("    return fmt.Sprintf(\"%%v\", s.MarshalMap(nil))")
		b.W("}\n")

		b.W("func (s *%s) MarshalMap(m map[string]interface{}) map[string]interface{} {", name)
		b.W("    if m == nil {")
		b.W("        m = make(map[string]interface{})")
		b.W("    }")
		for _, field := range st.fields {
			if field.t.t.Kind == KindPad {
				continue
			}
			fieldName := Capitalize(field.public)
			switch field.field.Type.Kind {
			case KindStruct:
				if field.field.Type.Optional {
					b.W("    {")
					b.W("        v := s.%s()", fieldName)
					b.W("        if v == nil {")
					b.W("            m[\"%s\"] = nil", field.field.Name)
					b.W("        } else {")
					b.W("            m[\"%s\"] = v.MarshalMap(nil)", field.field.Name)
					b.W("        }")
					b.W("    }")
				} else {
					b.W("    m[\"%s\"] = s.%s().MarshalMap(nil)", field.field.Name, fieldName)
				}
			case KindList:
				if field.field.Type.Optional {
					b.W("    {")
					b.W("        v := s.%s()", fieldName)
					b.W("        if v == nil {")
					b.W("            m[\"%s\"] = nil", field.field.Name)
					b.W("        } else {")
					if field.field.Type.Element.Kind == KindStruct || field.field.Type.Element.Kind == KindUnion {
						b.W("            m[\"%s\"] = v.MarshalMap(nil)", field.field.Name)
					} else {
						b.W("            m[\"%s\"] = s.%s().CopyTo(nil)", field.field.Name, fieldName)
					}
					b.W("        }")
					b.W("    }")
				} else {
					b.W("    m[\"%s\"] = s.%s().CopyTo(nil)", field.field.Name, fieldName)
				}
			default:
				if field.field.Type.Optional {
					b.W("    {")
					b.W("        v := s.%s()", fieldName)
					b.W("        if v == nil {")
					b.W("            m[\"%s\"] = nil", field.field.Name)
					b.W("        } else {")
					b.W("            m[\"%s\"] = *v", field.field.Name)
					b.W("        }")
					b.W("    }")
				} else {
					b.W("    m[\"%s\"] = s.%s()", field.field.Name, fieldName)
				}
			}
		}
		b.W("    return m")
		b.W("}\n")

		/*
			func (s *Position) MarshalBinary() (data []byte, err error) {
			    return s[0:], nil
			}

			func (s *Position) UnmarshalBinary(b []byte) error {
			    v, e := UnmarshalPosition(b)
			    if e != nil {
			        return e
			    }
			    *s = *v
			    return nil
			}
		*/

		b.W("func (s *%s) Read(r io.Reader) error {", name)
		b.W("    n, err := io.ReadFull(r, s[0:])")
		b.W("    if err != nil {")
		b.W("        return err")
		b.W("    }")
		b.W("    if n != %d {", t.t.Size)
		b.W("        return io.ErrShortBuffer")
		b.W("    }")
		b.W("    return nil")
		b.W("}")

		b.W("func (s *%s) Write(w io.Writer) (n int, err error) {", name)
		b.W("    return w.Write(s[0:])")
		b.W("}")

		b.W("func (s *%s) MarshalBinaryTo(b []byte) []byte {", name)
		b.W("    return append(b, s[0:]...)")
		b.W("}")

		b.W("func (s *%s) MarshalBinary() ([]byte, error) {", name)
		b.W("    var v []byte")
		b.W("    return append(v, s[0:]...), nil")
		b.W("}")

		b.W("func (s *%s) UnmarshalBinary(b []byte) error {", name)
		b.W("    if len(b) < %d {", t.t.Size)
		b.W("        return io.ErrShortBuffer")
		b.W("    }")
		b.W("    v := (*%s)(unsafe.Pointer(&b[0]))", name)
		b.W("    *s = *v")
		b.W("    return nil")
		b.W("}")

		t.mut = fmt.Sprintf("%sMut", t.name)

		b.W("func (s *%s) Mut() *%s {", t.name, t.mut)
		b.W("    return (*%s)(unsafe.Pointer(&s[0]))", t.mut)
		b.W("}")
	} else {
		b.W("func (s *%s) Freeze() *%s {", t.mut, t.name)
		b.W("    return (*%s)(unsafe.Pointer(&s.%s[0]))", t.name, t.name)
		b.W("}")
	}

	if file.types["Reinterpret"+name] != nil {
		b.W("func Reinterpret%s_(b []byte) (*%s, error) {", name, name)
	} else {
		b.W("func Reinterpret%s(b []byte) (*%s, error) {", name, name)
	}
	b.W("    if len(b) < %d {", t.t.Size)
	b.W("        return nil, io.ErrShortBuffer")
	b.W("    }")
	b.W("    return (*%s)(unsafe.Pointer(&b[0])), nil", name)
	b.W("}\n")

	for _, field := range st.fields {
		if field.t.t.Kind == KindPad {
			continue
		}
		if mut {
			c.genStructFieldSetter(b, t, field, order)
			c.genStructFieldGetter(mut, b, t, field, order)
		} else {
			c.genStructFieldGetter(mut, b, t, field, order)
		}
	}
	return nil
}

func (c *Compiler) genStructFieldGetter(mut bool, b *Builder, st *goType, f *goField, order binary.ByteOrder) {
	goStructName := st.name
	if mut {
		goStructName = st.mut
	}
	fieldName := f.public
	typeName := f.t.name
	if mut && len(f.t.mut) > 0 {

		typeName = f.t.mut
	}
	//if strings.HasSuffix(typeName, "_") {
	//	typeName = typeName[0:len(typeName)-1]
	//}

	getBuffer := "s"
	if mut {
		getBuffer = fmt.Sprintf("s.%s", st.name)
	}
	c.genComments(b, f.field.Type.Comments)

	if f.field.Type.Optional {
		/*
			func Set(b, flag Bits) Bits    { return b | flag }
			func Clear(b, flag Bits) Bits  { return b &^ flag }
			func Toggle(b, flag Bits) Bits { return b ^ flag }
			func Has(b, flag Bits) bool    { return b&flag != 0 }
		*/
		b.W("func (s *%s) %s() *%s {", goStructName, fieldName, typeName)
		b.W("    if %s[%d]&%d == 0 {", getBuffer, f.field.OptOffset, f.field.OptMask)
		b.W("        return nil")
		b.W("    }")

		switch f.field.Type.Kind {
		case KindBool:
			//b.W("    return s[%d] != 0", f.Offset)
			b.W("    return (*%s)(unsafe.Pointer(&%s[%d]))", typeName, getBuffer, f.field.Offset)
		default:
			b.W("    return (*%s)(unsafe.Pointer(&%s[%d]))", typeName, getBuffer, f.field.Offset)
		}

		b.W("}\n")
	} else {
		if f.isPointer {
			b.W("func (s *%s) %s() *%s {", goStructName, fieldName, typeName)
			b.W("    return (*%s)(unsafe.Pointer(&%s[%d]))", typeName, getBuffer, f.field.Offset)
		} else {
			b.W("func (s *%s) %s() %s {", goStructName, fieldName, typeName)

			switch f.field.Type.Kind {
			case KindBool:
				b.W("    return %s[%d] != 0", getBuffer, f.field.Offset)
			default:
				b.W("    return *(*%s)(unsafe.Pointer(&%s[%d]))", typeName, getBuffer, f.field.Offset)
			}
		}
		b.W("}\n")
	}
}

func (c *Compiler) genStructFieldSetter(b *Builder, st *goType, f *goField, order binary.ByteOrder) {
	//embeddedStructName := st.name
	goStructName := st.mut
	fieldName := f.public
	typeName := f.t.name
	c.genComments(b, f.field.Type.Comments)

	getBuffer := fmt.Sprintf("s.%s", st.name)

	W := b.W
	if f.field.Type.Optional {
		/*
			func Set(b, flag Bits) Bits    { return b | flag }
			func Clear(b, flag Bits) Bits  { return b &^ flag }
			func Toggle(b, flag Bits) Bits { return b ^ flag }
			func Has(b, flag Bits) bool    { return b&flag != 0 }
		*/
		W("func (s *%s) Set%s(v *%s) *%s {", goStructName, fieldName, typeName, goStructName)
		W("    if v == nil {")
		W("        %s[%d] = %s[%d] &^ %d", getBuffer, f.field.OptOffset, getBuffer, f.field.OptOffset, f.field.OptMask)
		W("        *(*%s)(unsafe.Pointer(&%s[%d])) = %s{}", typeName, getBuffer, f.field.Offset, typeName)
		W("        return s")
		W("    }")
		W("    %s[%d] = %s[%d] | %d", getBuffer, f.field.OptOffset, getBuffer, f.field.OptOffset, f.field.OptMask)

		switch f.field.Type.Kind {
		case KindBool:
			W("    if v {")
			W("        %s[%d] = 1", getBuffer, f.field.Offset)
			W("    } else {")
			W("        %s[%d] = 0", getBuffer, f.field.Offset)
			W("    }")

		default:
			W("    *(*%s)(unsafe.Pointer(&%s[%d])) = *v", typeName, getBuffer, f.field.Offset)
		}

		W("    return s")
		W("}\n")
	} else {
		if f.isPointer {
			W("func (s *%s) Set%s(v *%s) *%s {", goStructName, fieldName, typeName, goStructName)
			W("    if v == nil {")
			W("        v = &%s{}", typeName)
			W("    }")
			W("    *(*%s)(unsafe.Pointer(&%s[%d])) = *v", typeName, getBuffer, f.field.Offset)
			W("    return s")
		} else {
			W("func (s *%s) Set%s(v %s) *%s {", goStructName, fieldName, typeName, goStructName)

			switch f.field.Type.Kind {
			case KindBool:
				W("    if v {")
				W("        %s[%d] = 1", getBuffer, f.field.Offset)
				W("    } else {")
				W("        %s[%d] = 0", getBuffer, f.field.Offset)
				W("    }")

			default:
				W("    *(*%s)(unsafe.Pointer(&%s[%d])) = v", typeName, getBuffer, f.field.Offset)
			}

			W("    return s")
		}
		W("}\n")
	}
}

func (c *Compiler) genStruct(file *goPackage, t *goType, mut bool, b *Builder, order binary.ByteOrder) error {
	st := t.st
	c.genComments(b, t.t.Comments)
	W := b.W

	headerName := ""
	if t.t.HeaderSize > 0 {
		headerName = headerFieldName
	}

	if mut {
		//headerName = fmt.Sprintf("%s.%s", t.name, headerName)
		W("type %s struct {", t.mut)
		W("    %s", t.name)
		W("}")

		W("func (s *%s) Clone() *%s {", t.mut, t.mut)
		W("    v := &%s{}", t.mut)
		W("    *v = *s")
		W("    return v")
		W("}")

		W("func (s *%s) Freeze() *%s {", t.mut, t.name)
		W("    return (*%s)(unsafe.Pointer(s))", t.name)
		W("}")
	} else {
		W("type %s struct {", t.name)

		longestName := len(headerName)

		for _, field := range st.fields {
			if len(field.private) > longestName {
				longestName = len(field.private)
			}
		}
		longestName += 1

		if t.t.HeaderSize > 0 {
			W("    %s [%d]byte // Header", PadEnd(headerName, longestName), t.t.HeaderSize)
		}
		for _, field := range st.fields {
			if field.t.t.Kind == KindPad {
				W("    %s [%d]byte // Padding", PadEnd("_", longestName), field.t.t.Size)
			} else {
				W("    %s %s", PadEnd(field.private, longestName), field.t.name)
			}
		}
		W("}")

		W("func (s *%s) String() string {", t.name)
		W("    return fmt.Sprintf(\"%%v\", s.MarshalMap(nil))")
		W("}\n")

		W("func (s *%s) MarshalMap(m map[string]interface{}) map[string]interface{} {", t.name)
		W("    if m == nil {")
		W("        m = make(map[string]interface{})")
		W("    }")
		for _, field := range st.fields {
			if field.t.t.Kind == KindPad {
				continue
			}
			fieldName := Capitalize(field.public)
			switch field.field.Type.Kind {
			case KindStruct:
				if field.field.Type.Optional {
					W("    {")
					W("        v := s.%s()", fieldName)
					W("        if v == nil {")
					W("            m[\"%s\"] = nil", field.field.Name)
					W("        } else {")
					W("            m[\"%s\"] = v.MarshalMap(nil)", field.field.Name)
					W("        }")
					W("    }")
				} else {
					W("    m[\"%s\"] = s.%s().MarshalMap(nil)", field.field.Name, fieldName)
				}
			case KindList:
				if field.field.Type.Optional {
					W("    {")
					W("        v := s.%s()", fieldName)
					W("        if v == nil {")
					W("            m[\"%s\"] = nil", field.field.Name)
					W("        } else {")
					if field.field.Type.Element.Kind == KindStruct || field.field.Type.Element.Kind == KindUnion {
						W("            m[\"%s\"] = v.MarshalMap(nil)", field.field.Name)
					} else {
						W("            m[\"%s\"] = s.%s().CopyTo(nil)", field.field.Name, fieldName)
					}
					W("        }")
					W("    }")
				} else {
					W("    m[\"%s\"] = s.%s().CopyTo(nil)", field.field.Name, fieldName)
				}
			default:
				if field.field.Type.Optional {
					W("    {")
					W("        v := s.%s()", fieldName)
					W("        if v == nil {")
					W("            m[\"%s\"] = nil", field.field.Name)
					W("        } else {")
					W("            m[\"%s\"] = *v", field.field.Name)
					W("        }")
					W("    }")
				} else {
					W("    m[\"%s\"] = s.%s()", field.field.Name, fieldName)
				}
			}
		}
		W("    return m")
		W("}\n")

		/*
			func (s *Position) MarshalBinary() (data []byte, err error) {
			    return s[0:], nil
			}

			func (s *Position) UnmarshalBinary(b []byte) error {
			    v, e := UnmarshalPosition(b)
			    if e != nil {
			        return e
			    }
			    *s = *v
			    return nil
			}
		*/

		W("func (s *%s) ReadFrom(r io.Reader) (int64, error) {", t.name)
		W("    n, err := io.ReadFull(r, (*(*[%d]byte)(unsafe.Pointer(s)))[0:])", t.t.Size)
		W("    if err != nil {")
		W("        return int64(n), err")
		W("    }")
		W("    if n != %d {", t.t.Size)
		W("        return int64(n), io.ErrShortBuffer")
		W("    }")
		W("    return int64(n), nil")
		W("}")

		W("func (s *%s) WriteTo(w io.Writer) (int64, error) {", t.name)
		W("    n, err := w.Write((*(*[%d]byte)(unsafe.Pointer(s)))[0:])", t.t.Size)
		W("    return int64(n), err")
		W("}")

		W("func (s *%s) MarshalBinaryTo(b []byte) []byte {", t.name)
		W("    return append(b, (*(*[%d]byte)(unsafe.Pointer(s)))[0:]...)", t.t.Size)
		W("}")

		W("func (s *%s) MarshalBinary() ([]byte, error) {", t.name)
		W("    var v []byte")
		W("    return append(v, (*(*[%d]byte)(unsafe.Pointer(s)))[0:]...), nil", t.t.Size)
		W("}")

		W("func (s *%s) Read(b []byte) (n int, err error) {", t.name)
		W("    if len(b) < %d {", t.t.Size)
		W("        return -1, io.ErrShortBuffer")
		W("    }")
		W("    v := (*%s)(unsafe.Pointer(&b[0]))", t.name)
		W("    *v = *s")
		W("    return %d, nil", t.t.Size)
		W("}")

		W("func (s *%s) UnmarshalBinary(b []byte) error {", t.name)
		W("    if len(b) < %d {", t.t.Size)
		W("        return io.ErrShortBuffer")
		W("    }")
		W("    v := (*%s)(unsafe.Pointer(&b[0]))", t.name)
		W("    *s = *v")
		W("    return nil")
		W("}")

		W("func (s *%s) Clone() *%s {", t.name, t.name)
		W("    v := &%s{}", t.name)
		W("    *v = *s")
		W("    return v")
		W("}")

		W("func (s *%s) Bytes() []byte {", t.name)
		W("    return (*(*[%d]byte)(unsafe.Pointer(s)))[0:]", t.t.Size)
		W("}")

		W("func (s *%s) Mut() *%s {", t.name, t.mut)
		W("    return (*%s)(unsafe.Pointer(s))", t.mut)
		W("}")
	}

	// Getters
	for _, field := range st.fields {
		if field.t.t.Kind == KindPad {
			continue
		}

		if field.field.Type.Optional {
			if mut {
				if field.t.name != field.t.mut {
					W("func (s *%s) %s() *%s {", t.mut, field.public, field.t.mut)
					W("    if s.%s[%d]&%d == 0 {", headerName, field.field.OptOffset, field.field.OptMask)
					W("        return nil")
					W("    }")
					W("    return s.%s.Mut()", field.private)
					W("}")
				}

				W("func (s *%s) Set%s(v *%s) *%s {", t.mut, field.public, field.t.name, t.mut)
				W("    if v == nil {")
				W("        s.%s[%d] = s.%s[%d] &^ %d",
					headerName,
					field.field.OptOffset,
					headerName,
					field.field.OptOffset,
					field.field.OptMask,
				)
				W("        return s")
				W("    }")
				W("    s.%s = *v", field.private)
				W("    return s")
				W("}")
			} else {
				W("func (s *%s) %s() *%s {", t.name, field.public, field.t.name)
				W("    if s.%s[%d]&%d == 0 {", headerName, field.field.OptOffset, field.field.OptMask)
				W("        return nil")
				W("    }")
				W("    return &s.%s", field.private)
				W("}")
			}
		} else {
			if mut {
				if field.isPointer {
					if field.t.name != field.t.mut {
						if strings.HasSuffix(field.t.mut, "_") {
							field.t.mut = field.t.mut[0 : len(field.t.mut)-1]
						}
						W("func (s *%s) %s() *%s {", t.mut, field.public, field.t.mut)
						W("    return s.%s.Mut()", field.private)
						W("}")
					}

					W("func (s *%s) Set%s(v *%s) *%s {", t.mut, field.public, field.t.name, t.mut)
					W("    s.%s = *v", field.private)
					W("    return s")
					W("}")
				} else {
					W("func (s *%s) Set%s(v %s) *%s {", t.mut, field.public, field.t.name, t.mut)
					W("    s.%s = v", field.private)
					W("    return s")
					W("}")
				}
			} else {
				if field.isPointer || field.t.t.Kind == KindBytes {
					W("func (s *%s) %s() *%s {", t.name, field.public, field.t.name)
					W("    return &s.%s", field.private)
					W("}")
				} else {
					W("func (s *%s) %s() %s {", t.name, field.public, field.t.name)
					W("    return s.%s", field.private)
					W("}")
				}
			}
		}
	}

	return nil
}

func (c *Compiler) genArrayList(t *goType, mut bool, b *Builder, order binary.ByteOrder) error {
	if t.list == nil {
		return errors.New("type is not a list")
	}
	W := b.W

	if mut {
		if strings.HasSuffix(t.mut, "_") {
			t.mut = t.mut[0 : len(t.mut)-1]
		}
		W("type %s struct {", t.mut)
		W("    %s", t.name)
		W("}")

		W("func (s *%s) setLen(l int) {", t.mut)
		if t.t.HeaderSize == 1 {
			W("    s.l = byte(l)")
		} else {
			//W("    s[%d] = byte(l)", t.t.Size-2)
			//W("    s[%d] = byte(l >> 8)", t.t.Size-1)
			//W("    *(*uint16)(unsafe.Pointer(&s[%d])) = uint16(l)", t.t.Size-2)
			W("    s.l = uint16(l)")
		}
		W("}")

		if c.isPointerType(t.list.element.t) {
			W("func (s *%s) Push(v *%s) bool {", t.mut, t.list.element.name)
			W("    l := s.Len()")
			W("    if l == %d {", t.t.Len)
			W("        return false")
			W("    }")
			W("    s.b[l] = *v")
			//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = *v", t.list.element.name, t.t.ItemSize)
			W("    s.setLen(l+1)")
			W("    return true")
			W("}")
		} else {
			W("func (s *%s) Push(v %s) bool {", t.mut, t.list.element.name)
			W("    l := s.Len()")
			W("    if l == %d {", t.t.Len)
			W("        return false")
			W("    }")
			W("    s.b[l] = v")
			//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = *v", t.list.element.name, t.t.ItemSize)
			W("    s.setLen(l+1)")
			W("    return true")
			W("}")
		}

		W("// Removes the last item")
		W("func (s *%s) Pop(v *%s) bool {", t.mut, t.list.element.name)
		W("    l := s.Len()")
		W("    if l == 0 {")
		W("        return false")
		W("    }")
		W("    l -= 1")
		W("    if v != nil {")
		W("        *v = s.b[l]")
		//W("        *v = *(*%s)(unsafe.Pointer(&s[l * %d]))", t.list.element.name, t.t.ItemSize)
		W("    }")
		// Clear last element.
		if t.list.element.primitive {
			W("    s.b[l] = 0")
			//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = 0", t.list.element.name, t.t.ItemSize)
		} else {
			W("    s.b[l] = %s{}", t.list.element.name)
			//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = %s{}", t.list.element.name, t.t.ItemSize, t.list.element.name)
		}
		W("    s.setLen(l)")
		W("    return true")
		W("}")

		W("// Removes the first item")
		W("func (s *%s) Shift(v *%s) bool {", t.mut, t.list.element.name)
		W("    l := s.Len()")
		W("    if l == 0 {")
		W("        return false")
		W("    }")
		W("    if v != nil {")
		W("        *v = s.b[0]")
		W("    }")
		W("    if l > 1 {")
		// Shift bytes over
		W("        copy(s.b[0:], s.b[1:s.l])")
		//W("        copy(s[0:], s[%d:l*%d])", t.t.ItemSize, t.t.ItemSize)
		W("    }")
		W("    l -= 1")
		// Clear last element
		if t.list.element.primitive {
			W("    s.b[l] = 0")
			//W("    *(*%s)(unsafe.Pointer(&s[l*%d])) = 0", t.list.element.name, t.t.ItemSize)
		} else {
			W("    s.b[l] = %s{}", t.list.element.name)
			//W("    *(*%s)(unsafe.Pointer(&s[l*%d])) = %s{}", t.list.element.name, t.t.ItemSize, t.list.element.name)
		}
		W("    s.setLen(l)")
		W("    return true")
		W("}")

		W("func (s *%s) Clear() {", t.mut)
		W("    s.b = [%d]%s{}", t.t.Len, t.list.element.name)
		W("    s.l = 0")
		W("}")
	} else {
		W("type %s struct {", t.name)
		W("    b [%d]%s", t.t.Len, t.list.element.name)
		if t.t.Padding > 0 {
			W("    _ [%d]byte // Padding", t.t.Padding)
		}
		if t.t.HeaderSize == 1 {
			W("    l byte")
		} else if t.t.HeaderSize == 2 {
			W("    l uint16")
		}
		W("}")
		//W("type %s [%d]byte\n", t.name, t.t.Size)
		W("func (s *%s) Get(i int) *%s {", t.name, t.list.element.name)
		W("    if i < 0 || i >= s.Len() {")
		W("        return nil")
		W("    }")
		W("    return &s.b[i]")
		W("}")

		W("func (s *%s) Len() int {", t.name)
		W("    return int(s.l)")
		//if t.t.Len < 256 {
		//	W("    return int(s[%d])", t.t.Size-1)
		//} else {
		//	W("    return int(*(*uint16)(unsafe.Pointer(&s[%d])))", t.t.Size-2)
		//	//W("    return int(uint16(s[%d]) | uint16(s[%d]) << 8)", t.t.Size-2, t.t.Size-1)
		//	//W("    return int(*(*uint16)(unsafe.Pointer(&s[%d])))", t.t.Size-2)
		//}
		W("}")

		W("func (s *%s) Cap() int {", t.name)
		W("    return %d", t.t.Len)
		W("}")

		if t.list.element.st != nil {
			W("func (s *%s) MarshalMap(m []map[string]interface{}) []map[string]interface{} {", t.name)
			W("    if m == nil {")
			W("        m = make([]map[string]interface{}, 0, s.Len())")
			W("    }")
			W("    for _, v := range s.Unsafe() {")
			W("        m = append(m, v.MarshalMap(nil))")
			W("    }")
			W("    return m")
			W("}")
		}

		W("func (s *%s) CopyTo(v []%s) []%s {", t.name, t.list.element.name, t.list.element.name)
		W("    return append(v, s.Unsafe()...)")
		W("}")

		W("func (s *%s) Unsafe() []%s {", t.name, t.list.element.name)
		W("    return s.b[0:s.Len()]")
		//W("    return *(*[]%s)(unsafe.Pointer(&reflect.SliceHeader{", t.list.element.name)
		//W("        Data: uintptr(unsafe.Pointer(&s.b[0])),")
		//W("        Len: s.Len(),")
		//W("        Cap: s.Cap(),")
		//W("    }))")
		W("}")

		W("func (s *%s) Bytes() []byte {", t.name)
		W("    return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{")
		W("        Data: uintptr(unsafe.Pointer(&s.b[0])),")
		if t.t.ItemSize > 1 {
			W("        Len: s.Len() * %d,", t.t.ItemSize)
		} else {
			W("        Len: s.Len(),")
		}
		W("        Cap: %d,", t.t.ItemSize*t.t.Len)
		W("    }))")
		W("}")

		W("func (s *%s) Mut() *%s {", t.name, t.mut)
		W("    return *(**%s)(unsafe.Pointer(&s))", t.mut)
		W("}")

		W("func (s *%s) Read(r io.Reader) error {", t.name)
		W("    n, err := io.ReadFull(r, (*(*[%d]byte)(unsafe.Pointer(&s)))[0:])", t.t.Size)
		W("    if err != nil {")
		W("        return err")
		W("    }")
		W("    if n != %d {", t.t.Size)
		W("        return io.ErrShortBuffer")
		W("    }")
		W("    return nil")
		W("}")

		W("func (s *%s) Write(w io.Writer) (n int, err error) {", t.name)
		W("    return w.Write((*(*[%d]byte)(unsafe.Pointer(&s)))[0:])", t.t.Size)
		W("}")

		W("func (s *%s) MarshalBinaryTo(b []byte) []byte {", t.name)
		W("    return append(b, (*(*[%d]byte)(unsafe.Pointer(&s)))[0:]...)", t.t.Size)
		W("}")

		W("func (s *%s) MarshalBinary() ([]byte, error) {", t.name)
		W("    var v []byte")
		W("    return append(v, (*(*[%d]byte)(unsafe.Pointer(&s)))[0:]...), nil", t.t.Size)
		W("}")

		W("func (s *%s) UnmarshalBinary(b []byte) error {", t.name)
		W("    if len(b) < %d {", t.t.Size)
		W("        return io.ErrShortBuffer")
		W("    }")
		W("    v := (*%s)(unsafe.Pointer(&b[0]))", t.name)
		W("    *s = *v")
		W("    return nil")
		W("}")
	}

	//W("func (s *%s) Hash() uint32 {", t.name)
	//W("    return bu.Hash32Bytes(s[0:s.Len()*%d])", t.t.ItemSize)
	//W("}")
	//
	//W("func (s *%s) Hash64() uint64 {", t.name)
	//W("    return bu.Hash64Bytes(s[0:s.Len()*%d])", t.t.ItemSize)
	//W("}")

	return nil
}

func (c *Compiler) genUnion(t *goType, b *Builder, order binary.ByteOrder) error {
	return nil
}

func (c *Compiler) genString(t *goType, mut bool, b *Builder, order binary.ByteOrder) error {
	W := b.W
	size := t.t.Len
	sizeIndex := size - 1
	sizeBytes := 1
	if size <= 256 {
		sizeBytes = 1
	} else if size < 65535 {
		sizeBytes = 2
		sizeIndex--
	}

	if mut {
		W("type %s struct {", t.mut)
		W("    %s", t.name)
		W("}")

		W("func (s *%s) Set(v string) {", t.mut)
		W("    s.set(v)")
		W("}")
	} else {
		W("type %s [%d]byte", t.name, size)

		W("func New%s(s string) *%s {", t.name, t.name)
		W("    v := %s{}", t.name)
		W("    v.set(s)")
		W("    return &v")
		W("}")

		if t.t.Kind == KindBytes {
			W("func (s *%s) set(v string) {", t.name)
			W("    copy(s[0:], v)")
			W("}")
		} else {
			W("func (s *%s) set(v string) {", t.name)
			W("    copy(s[0:%d], v)", sizeIndex)
			W("    c := %d", sizeIndex)
			W("    l := len(v)")
			W("    if l > c {")
			if sizeBytes == 1 {
				W("        s[%d] = byte(c)", sizeIndex)
			} else if sizeBytes == 2 {
				W("        s[%d] = byte(c)", sizeIndex)
				W("        s[%d] = byte(c >> 8)", sizeIndex+1)
			}
			W("    } else {")
			if sizeBytes == 1 {
				W("        s[%d] = byte(l)", sizeIndex)
			} else if sizeBytes == 2 {
				W("        s[%d] = byte(l)", sizeIndex)
				W("        s[%d] = byte(l >> 8)", sizeIndex+1)
			}
			W("    }")
			W("}")
		}

		if t.t.Kind == KindBytes {
			W("func (s *%s) Len() int {", t.name)
			W("    return %d", t.t.Len)
			W("}")

			W("func (s *%s) Cap() int {", t.name)
			W("    return %d", t.t.Len)
			W("}")
		} else {
			W("func (s *%s) Len() int {", t.name)
			if sizeBytes == 1 {
				W("    return int(s[%d])", sizeIndex)
			} else if sizeBytes == 2 {
				W("    return int(*(*uint16)(unsafe.Pointer(&s[%d])))", sizeIndex)
				//W("    return int(uint16(s[%d]) | uint16(s[%d]) << 8)", sizeIndex, sizeIndex+1)
			}
			W("}")

			W("func (s *%s) Cap() int {", t.name)
			W("    return %d", sizeIndex)
			W("}")
		}

		W("func (s *%s) StringClone() string {", t.name)
		if t.t.Kind == KindBytes {
			W("    b := s[0:%d]", t.t.Len)
			W("    return string(b)")
		} else {
			W("    b := s[0:s.Len()]")
			W("    return string(b)")
		}
		W("}")

		W("func (s *%s) String() string {", t.name)
		if t.t.Kind == KindBytes {
			W("    b := s[0:%d]", t.t.Len)
			W("    return *(*string)(unsafe.Pointer(&b))")
		} else {
			W("    b := s[0:s.Len()]")
			W("    return *(*string)(unsafe.Pointer(&b))")
		}
		W("}")

		W("func (s *%s) Bytes() []byte {", t.name)
		W("    return s[0:s.Len()]")
		W("}")

		W("func (s *%s) Clone() *%s {", t.name, t.name)
		W("    v := %s{}", t.name)
		W("    copy(s[0:], v[0:])")
		W("    return &v")
		W("}")

		W("func (s *%s) Mut() *%s {", t.name, t.mut)
		W("    return *(**%s)(unsafe.Pointer(&s))", t.mut)
		W("}")

		W("func (s *%s) ReadFrom(r io.Reader) error {", t.name)
		W("    n, err := io.ReadFull(r, (*(*[%d]byte)(unsafe.Pointer(&s)))[0:])", t.t.Size)
		W("    if err != nil {")
		W("        return err")
		W("    }")
		W("    if n != %d {", t.t.Size)
		W("        return io.ErrShortBuffer")
		W("    }")
		W("    return nil")
		W("}")

		W("func (s *%s) WriteTo(w io.Writer) (n int, err error) {", t.name)
		W("    return w.Write((*(*[%d]byte)(unsafe.Pointer(&s)))[0:])", t.t.Size)
		W("}")

		W("func (s *%s) MarshalBinaryTo(b []byte) []byte {", t.name)
		W("    return append(b, (*(*[%d]byte)(unsafe.Pointer(&s)))[0:]...)", t.t.Size)
		W("}")

		W("func (s *%s) MarshalBinary() ([]byte, error) {", t.name)
		//W("    var v []byte")
		//W("    return append(v, (*(*[%d]byte)(unsafe.Pointer(&s)))[0:]...), nil", t.t.Size)
		W("    return s[0:s.Len()], nil")
		W("}")

		W("func (s *%s) UnmarshalBinary(b []byte) error {", t.name)
		W("    if len(b) < %d {", t.t.Size)
		W("        return io.ErrShortBuffer")
		W("    }")
		W("    v := (*%s)(unsafe.Pointer(&b[0]))", t.name)
		W("    *s = *v")
		W("    return nil")
		W("}")
	}

	//W("func (s *String%d) Equal(v %sString) bool {", size, buPrefix)
	//W("    return s.Cast() == v.Cast()")
	//W("}")
	//
	//W("func (s *String%d) Hash() uint32 {", size)
	//W("    return %sHash32(s.Cast())", buPrefix)
	//W("}")
	//
	//W("func (s *String%d) Hash64() uint64 {", size)
	//if sizeBytes == 1 {
	//	W("    return %sHash64(s.Cast())", buPrefix)
	//} else {
	//	W("    return %sHash64(s.Cast())", buPrefix)
	//	//W("    return %sHash64Bytes(s[0:uint16(s[%d]) | uint16(s[%d])])", buPrefix, sizeIndex, sizeIndex+1)
	//}
	//W("}\n")
	return nil
}
