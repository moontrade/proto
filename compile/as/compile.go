package as

import (
	"errors"
	"fmt"
	"sort"

	//. "github.com/moontrade/proto/compile"
	. "github.com/moontrade/proto/schema"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func NewCompiler(schema *Schema, config *ASConfig) (*Compiler, error) {
	return &Compiler{
		schema:   schema,
		config:   config,
		packages: make(map[string]*asPackage),
	}, nil
}

func (c *Compiler) Compile() error {
	// Convert into Go specific model
	packages := make(map[string]*asPackage)
	var err error
	for k, v := range c.schema.Files {
		packages[k], err = c.createPackage(v, 0)
		if err != nil {
			return err
		}
	}

	path := c.schema.Config.Path
	err = filepath.Walk(path, c.walkClear)
	for _, f := range packages {
		b := NewBuilder()
		if err := c.writeFile(f, b); err != nil {
			return err
		}

		// Compute dir path
		dir := filepath.Join(path, f.path)

		// mkdir
		_ = os.MkdirAll(dir, 0755)

		// Create new "moon.go"
		out, err := os.Create(filepath.Join(dir, TSFileName))
		if err != nil {
			return err
		}
		_, err = out.WriteString(b.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) walkClear(path string, info fs.FileInfo, err error) error {
	if strings.HasSuffix(path, TSFileNameSuffix) {
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

func (c *Compiler) primitive(file *asPackage, t *Type, goName string) *asType {
	return &asType{
		pkg:       file,
		t:         t,
		name:      goName,
		mut:       goName,
		primitive: true,
	}
}

func (c *Compiler) stringType(file *asPackage, t *Type) *asType {
	name := fmt.Sprintf("String%d", t.Len)
	mut := fmt.Sprintf("%sMut", name)
	return &asType{
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

func (c *Compiler) toTSPath(f *File) string {
	packageParts := strings.Split(f.Package, ".")
	return joinWithSlash(c.config.Package, joinWithSlash(packageParts...))
}

func (c *Compiler) createPackage(file *File, level int) (*asPackage, error) {
	if level >= 20 {
		return nil, errors.New("cyclic package dependency")
	}
	if existing := c.packages[file.Package]; existing != nil {
		return existing, nil
	}
	packageParts := strings.Split(file.Package, ".")
	path := filepath.Join(packageParts...)
	pkg := &asPackage{
		file:        file,
		path:        path,
		dir:         filepath.Join(c.schema.Config.Path, path),
		packageName: packageParts[len(packageParts)-1],
		byType:      make(map[*Type]*asType),
		importMap:   make(map[string]*asImport),
		types:       make(map[string]*asType),
		structs:     make(map[string]*asType),
		strings:     make(map[string]*asType),
		enums:       make(map[string]*asType),
		lists:       make(map[string]*asType),
		unions:      make(map[string]*asType),
		names:       make(map[string]struct{}),
	}

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
		_, err := c.resolve(pkg, value, 0)
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

	c.packages[file.Package] = pkg

	return pkg, nil
}

func (c *Compiler) addImport(imports map[string]*asImport, path, alias string) *asImport {
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
	imp := &asImport{
		path:     path,
		alias:    alias,
		useAlias: useAlias,
	}
	imports[path] = imp
	return imp
}

//func relativePath(from, to string) string {
//	fromParts := strings.Split(from, "/")
//	toParts := strings.Split(to, "/")
//	i := 0
//	for {
//		if i >= len(toParts) ||  {
//			break
//		}
//		i++
//	}
//
//
//	return "../" + strings.Join(toParts[i:], "/")
//}

func (c *Compiler) writeFile(file *asPackage, b *Builder) error {
	//b.W("package %s\n", file.packageName)

	if len(file.importMap) > 0 {
		for _, imp := range file.importMap {
			path, err := filepath.Rel(file.path, imp.path)
			if err != nil {
				return err
			}
			b.W("import * as %s from '%s'", imp.alias, path)
		}
		b.W("")
	}

	for _, enum := range file.enums {
		if err := c.genEnum(file, enum, b); err != nil {
			return err
		}
	}

	init := NewBuilder()
	//init.W("func init() {")
	//
	//init.W(`    {
	//	var b [2]byte
	//    v := uint16(1)
	//    b[0] = byte(v)
	//    b[1] = byte(v >> 8)
	//	if *(*uint16)(unsafe.Pointer(&b[0])) != 1 {
	//		panic("BigEndian detected... compiled for LittleEndian only!!!")
	//	}
	//}`)
	//_, _ = init.Write([]byte(`    to := reflect.TypeOf
	//type sf struct {
	//    n string
	//    o uintptr
	//    s uintptr
	//}
	//ss := func(tt interface{}, mtt interface{}, s uintptr, fl []sf) {
	//    t := to(tt)
	//    mt := to(mtt)
	//    if t.Size() != s {
	//        panic(fmt.Sprintf("sizeof %%s = %%d, expected = %%d", t.Name(), t.Size(), s))
	//    }
	//    if mt.Size() != s {
	//        panic(fmt.Sprintf("sizeof %%s = %%d, expected = %%d", mt.Name(), mt.Size(), s))
	//    }
	//    if t.NumField() != len(fl) {
	//        panic(fmt.Sprintf("%%s field count = %%d: expected %%d", t.Name(), t.NumField(), len(fl)))
	//    }
	//    for i, ef := range fl {
	//        f := t.Field(i)
	//        if f.Offset != ef.o {
	//            panic(fmt.Sprintf("%%s.%%s offset = %%d, expected = %%d", t.Name(), f.Name, f.Offset, ef.o))
	//        }
	//        if f.Type.Size() != ef.s {
	//            panic(fmt.Sprintf("%%s.%%s size = %%d, expected = %%d", t.Name(), f.Name, f.Type.Size(), ef.s))
	//        }
	//        if f.Name != ef.n {
	//            panic(fmt.Sprintf("%%s.%%s expected field: %%s", t.Name(), f.Name, ef.n))
	//        }
	//    }
	//}`))
	//init.WriteLine("")
	//init.WriteLine("")

	for _, st := range file.structs {
		if err := c.genStruct(file, st, false, b); err != nil {
			return err
		}
		if err := c.genStruct(file, st, true, b); err != nil {
			return err
		}

		//init.W("    ss(%s{}, %s{}, %d, []sf{", st.name, st.mut, st.t.Size)
		//if st.t.HeaderSize > 0 {
		//	init.W("        {\"%s\", %d, %d},", headerFieldName, 0, st.t.HeaderSize)
		//}
		//for _, field := range st.st.fields {
		//	init.W("        {\"%s\", %d, %d},", field.private, field.field.Offset, field.t.t.Size)
		//}
		//init.W("    })")
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
		if err := c.genString(str, false, b); err != nil {
			return err
		}
		if err := c.genString(str, true, b); err != nil {
			return err
		}
	}

	for _, list := range file.lists {
		if err := c.genList(list, false, b); err != nil {
			return err
		}
		if err := c.genList(list, true, b); err != nil {
			return err
		}
	}

	//init.W("}\n")

	b.W(init.String())

	return nil
}

func (c *asPackage) uniqueName(n string) string {
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

func (c *Compiler) resolve(pkg *asPackage, t *Type, level int) (*asType, error) {
	if existing := pkg.byType[t]; existing != nil {
		return existing, nil
	}
	if t.Import != nil {
		imp := c.addImport(
			pkg.importMap,
			c.toTSPath(t.Import.Parent),
			t.Import.Alias,
		)

		importPkg, err := c.createPackage(t.Import.Parent, level+1)
		if err != nil {
			return nil, err
		}

		base := t.Base()
		importedType := importPkg.byType[base]
		if importedType != nil {

		}

		// Handle imported type
		return &asType{
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
			ret := &asType{}
			*ret = *existing
			ret.t = t
			pkg.byType[t] = ret
			return existing, nil
		}
		//_ = c.addImport(pkg.importMap, "fmt", "")
		//_ = c.addImport(pkg.importMap, "io", "")
		//_ = c.addImport(pkg.importMap, "reflect", "")
		//_ = c.addImport(pkg.importMap, "unsafe", "")

		fields := make([]*asField, 0, len(t.Struct.Fields))
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
			fields = append(fields, &asField{
				field:     field,
				isPointer: c.isPointerType(field.Type),
				public:    fieldName,
				private:   structField,
				t:         fieldType,
			})
			pkg.byType[field.Type] = fieldType
		}
		gt := &asType{
			pkg:       pkg,
			t:         t,
			name:      structName,
			primitive: false,
			st: &asStruct{
				st:     t.Struct,
				fields: fields,
			},
		}
		gt.mut = pkg.uniqueName(fmt.Sprintf("%sMut", gt.name))
		pkg.types[gt.name] = gt
		if pkg.structs == nil {
			pkg.structs = make(map[string]*asType)
		}
		pkg.structs[gt.name] = gt
		pkg.byType[t] = gt
		return gt, nil

	case KindEnum:
		enumName := Capitalize(t.Enum.Name)
		if existing := pkg.enums[enumName]; existing != nil {
			ret := &asType{}
			*ret = *existing
			ret.t = t
			pkg.byType[t] = existing
			return existing, nil
		}
		value, err := c.resolve(pkg, t.Element, level+1)
		if err != nil {
			return nil, err
		}
		enum := &asEnum{
			enum:    t.Enum,
			value:   value,
			options: make([]*asEnumOption, 0, len(t.Enum.Options)),
		}
		for _, option := range t.Enum.Options {
			o := &asEnumOption{
				option: option,
				name:   c.enumOptionName(option),
			}
			enum.options = append(enum.options, o)
		}
		gt := &asType{
			pkg:  pkg,
			t:    t,
			name: enumName,
			enum: enum,
		}
		pkg.types[gt.name] = gt
		if pkg.enums == nil {
			pkg.enums = make(map[string]*asType)
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
		gt := &asType{
			pkg:       pkg,
			t:         t,
			name:      Capitalize(t.Name),
			primitive: false,
			list: &asList{
				element: element,
			},
		}
		gt.mut = pkg.uniqueName(fmt.Sprintf("%sMut", gt.name))
		pkg.types[gt.name] = gt
		if pkg.lists == nil {
			pkg.lists = make(map[string]*asType)
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
		return c.primitive(pkg, t, "u8"), nil
	case KindInt8:
		return c.primitive(pkg, t, "i8"), nil
	case KindInt16:
		return c.primitive(pkg, t, "i16"), nil
	case KindUInt16:
		return c.primitive(pkg, t, "u16"), nil
	case KindInt32:
		return c.primitive(pkg, t, "i32"), nil
	case KindUInt32:
		return c.primitive(pkg, t, "u32"), nil
	case KindInt64:
		return c.primitive(pkg, t, "i64"), nil
	case KindUInt64:
		return c.primitive(pkg, t, "u64"), nil
	case KindFloat32:
		return c.primitive(pkg, t, "f32"), nil
	case KindFloat64:
		return c.primitive(pkg, t, "f64"), nil
	case KindString:
		gt := c.stringType(pkg, t)
		pkg.strings[gt.name] = gt
		pkg.byType[t] = gt
		return gt, nil

	case KindPad:
		return &asType{
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

func (c *Compiler) writeComments(prefix string, b *Builder, comments []string) {
	if len(comments) == 0 {
		return
	}
	for _, comment := range comments {
		b.W("%s//%s", prefix, comment)
	}
}

func (c *Compiler) fieldName(f string) string {
	f = Uncapitalize(f)
	switch f {
	case "free":
		return "free_"
	case "alloc":
		return "alloc_"
	case "dealloc":
		return "dealloc_"
	case "wrap":
		return "wrap_"
	case "unsafe":
		return "unsafe_"
	case "mut":
		return "mut_"
	case "freeze":
		return "freeze_"
	case "sizeof":
		return "sizeof_"
	case "toString":
		return "toString_"
	}
	return f
}

func (c *Compiler) enumOptionName(option *EnumOption) string {
	return option.Name
}

func (c *Compiler) unionOptionName(option *UnionOption) string {
	return fmt.Sprintf("%s_%s", Capitalize(option.Union.Name), Capitalize(option.Name))
}

func (c *Compiler) asTypeName(t *Type) string {
	switch t.Kind {
	case KindBool:
		return "bool"
	case KindByte:
		return "u8"
	case KindInt8:
		return "i8"
	case KindInt16:
		return "i16"
	case KindUInt16:
		return "u16"
	case KindInt32:
		return "i32"
	case KindUInt32:
		return "u32"
	case KindInt64:
		return "i64"
	case KindUInt64:
		return "u64"
	case KindFloat32:
		return "f32"
	case KindFloat64:
		return "f64"
	case KindString:
		return "string"
	case KindEnum:
		return Capitalize(t.Enum.Name)
	case KindStruct:
		return Capitalize(t.Struct.Name)
	case KindUnion:
		return Capitalize(t.Union.Name)
	case KindList:
		return fmt.Sprintf("[%d]%s", t.Len, c.asTypeName(t.Element))
	case KindPad:
		return fmt.Sprintf("[%d]byte", t.Size)
	}
	return "unknown"
}

func (c *Compiler) genEnum(file *asPackage, t *asType, b *Builder) error {
	b.W("export namespace %s {", t.name)
	for _, option := range t.enum.options {
		c.writeComments("    ", b, option.option.Comments)
		b.W("    export const %s:%s = %d", option.name, t.name, option.option.Value)
		//if i < len(t.enum.options)-1 {
		//	b.W("    export const %s:%s = %d,", option.name, t.enum.value.name, option.option.Value)
		//} else {
		//
		//}
	}
	b.W("}")
	c.writeComments("    ", b, t.t.Comments)
	b.W("export type %s = %s\n", t.name, t.enum.value.name)
	return nil
}

func (c *Compiler) writeFieldGetter(mut bool, b *Builder, st *asType, f *asField) {
	fieldName := f.public
	typeName := f.t.name
	if mut && len(f.t.mut) > 0 {
		typeName = f.t.mut
	}

	getBuffer := "changetype<usize>(this)"
	c.writeComments("    ", b, f.field.Type.Comments)

	W := b.W
	if f.field.Type.Optional {
		/*
			func Set(b, flag Bits) Bits    { return b | flag }
			func Clear(b, flag Bits) Bits  { return b &^ flag }
			func Toggle(b, flag Bits) Bits { return b ^ flag }
			func Has(b, flag Bits) bool    { return b&flag != 0 }
		*/
		W("    @inline get %s(): %s | null {", fieldName, typeName)
		W("        let flag = load<u8>(%s+%d)&%d", getBuffer, f.field.OptOffset, f.field.OptMask)
		W("        if (flag == 0) {")
		W("            return null")
		W("        }")

		switch f.field.Type.Kind {
		case KindBool:
			//W("    return s[%d] != 0", f.Offset)
			//W("    return (*%s)(unsafe.Pointer(&%s[%d]))", typeName, getBuffer, f.field.Offset)
			W("        return load<u8>(%s+%d) != 0", getBuffer, f.field.Offset)
		default:
			if f.isPointer {
				W("        return changetype<%s>(%s+%d)", typeName, getBuffer, f.field.Offset)
			} else {
				W("        return load<%s>(%s+%d)", typeName, getBuffer, f.field.Offset)
			}
		}

		W("    }")
	} else {
		if f.isPointer {
			W("    @inline get %s(): %s {", fieldName, f.t.name)
			W("        return changetype<%s>(%s+%d)", f.t.name, getBuffer, f.field.Offset)
			//W("    return (*%s)(unsafe.Pointer(&%s[%d]))", typeName, getBuffer, f.field.Offset)
		} else {
			W("    @inline get %s(): %s {", fieldName, typeName)
			switch f.field.Type.Kind {
			case KindBool:
				W("        return load<u8>(%s+%d) != 0", getBuffer, f.field.Offset)
			default:
				W("        return load<%s>(%s+%d)", typeName, getBuffer, f.field.Offset)
			}
		}
		W("    }")
	}
}

func (c *Compiler) writeFieldSetter(mut bool, b *Builder, st *asType, f *asField) {
	//embeddedStructName := st.name
	fieldName := f.public
	c.writeComments("    ", b, f.field.Type.Comments)

	getBuffer := "changetype<usize>(this)"
	name := f.t.name
	if mut {
		name = f.t.mut
	}

	W := b.W
	if f.field.Type.Optional {
		/*
			func Set(b, flag Bits) Bits    { return b | flag }
			func Clear(b, flag Bits) Bits  { return b &^ flag }
			func Toggle(b, flag Bits) Bits { return b ^ flag }
			func Has(b, flag Bits) bool    { return b&flag != 0 }
		*/

		W("    set %s(v: %s | null) {", fieldName, name)
		W("        if (v == null) {")
		W("            memory.fill(%s+%d, 0, %d)", getBuffer, f.field.Offset, f.t.t.Size)
		//W("        %s[%d] = %s[%d] &^ %d", getBuffer, f.field.OptOffset, getBuffer, f.field.OptOffset, f.field.OptMask)
		//W("        *(*%s)(unsafe.Pointer(&%s[%d])) = %s{}", typeName, getBuffer, f.field.Offset, typeName)
		W("            return")
		W("        }")
		W("        store<u8>(%s+%d, load<u8>(%s+%d) | %d)", getBuffer, f.field.OptOffset, getBuffer, f.field.OptOffset, f.field.OptMask)
		//W("        %s[%d] = %s[%d] | %d", getBuffer, f.field.OptOffset, getBuffer, f.field.OptOffset, f.field.OptMask)

		switch f.field.Type.Kind {
		case KindBool:
			W("        store<u8>(%s+%d, v ? 1 : 0)", getBuffer, f.field.Offset)

		default:
			if f.isPointer {
				W("        memory.copy(%s+%d, changetype<usize>(v), %d)", getBuffer, f.field.Offset, f.t.t.Size)
			} else {
				W("        store<%s>(%s+%d, v)", f.t.name, getBuffer, f.field.Offset)
			}
			//W("    *(*%s)(unsafe.Pointer(&%s[%d])) = *v", typeName, getBuffer, f.field.Offset)
		}
		W("    }\n")
	} else {
		if f.isPointer {
			W("    set %s(v: %s) {", fieldName, f.t.name)
			W("        memory.copy(%s+%d, changetype(v), %d)", getBuffer, f.field.Offset, f.t.t.Size)
		} else {
			W("    set %s(v: %s) {", fieldName, f.t.name)

			switch f.field.Type.Kind {
			case KindBool:
				W("        store<u8>(%s+%d, v ? 1 : 0)", getBuffer, f.field.Offset)

			default:
				W("        store<%s>(%s+%d, v)", f.t.name, getBuffer, f.field.Offset)
			}

			//W("    return s")
		}
		W("    }\n")
	}
}

func (c *Compiler) genStruct(file *asPackage, t *asType, mut bool, b *Builder) error {
	st := t.st
	c.writeComments("", b, t.t.Comments)
	W := b.W

	name := t.name
	if mut {
		name = t.mut
	}
	W("@unmanaged")
	W("export class %s {", name)
	for i := 0; i < t.t.Size; i += 8 {
		W("    private _%d: u64", i)
	}

	//W("    private constructor(public dataStart: usize, public isOwned: boolean = false) {}\n")

	W("    @inline static get sizeof(): usize {")
	W("        return %d", t.t.Size)
	W("    }\n")

	//W("    @inline static get alloc(): %s {", name)
	//W("        return new %s(heap.alloc(%d), true)", name, t.t.Size)
	//W("    }\n")
	//
	//W("    @inline dealloc(): void {")
	//W("        if (this.dataStart == 0 || !this.isOwned) return")
	//W("        heap.free(this.dataStart)")
	//W("        this.dataStart = 0")
	//W("    }\n")

	//W("    @inline static unsafe(ptr: usize, owned: boolean = false): %s {", name)
	//W("        return new %s(ptr, owned)", name)
	//W("    }\n")
	//
	//W("    @inline static wrap(buf: ArrayBuffer): %s {", name)
	//W("        return new %s(changetype<usize>(buf), false)", name)
	//W("    }\n")

	if mut {
		W("    @inline get freeze(): %s {", t.name)
		W("        return changetype<%s>(changetype<usize>(this))", t.name)
		//W("        return %s.unsafe(this.dataStart, this.isOwned)", t.name)
		W("    }\n")
	} else {
		W("    @inline get mut(): %s {", t.mut)
		W("        return changetype<%s>(changetype<usize>(this))", t.mut)
		W("    }\n")
	}
	// Getters
	for _, field := range st.fields {
		if field.t.t.Kind == KindPad {
			continue
		}

		c.writeFieldGetter(mut, b, t, field)
		if mut {
			c.writeFieldSetter(mut, b, t, field)
		}
	}

	W("}\n")

	return nil
}

func (c *Compiler) genList(t *asType, mut bool, b *Builder) error {
	if t.list == nil {
		return errors.New("type is not a list")
	}
	W := b.W

	name := t.name
	if mut {
		name = t.mut
	}

	W("@unmanaged export class %s {", name)

	W("    private constructor(public dataStart: usize, public isOwned: boolean = false) {}\n")

	W("    @inline static get sizeof(): usize {")
	W("        return %d", t.t.Size)
	W("    }\n")

	W("    @inline static get alloc(): %s {", name)
	W("        return new %s(heap.alloc(%d), true)", name, t.t.Size)
	W("    }\n")

	W("    @inline dealloc(): void {")
	W("        if (this.dataStart == 0 || !this.isOwned) return")
	W("        heap.free(this.dataStart)")
	W("        this.dataStart = 0")
	W("    }\n")

	W("    @inline static unsafe(ptr: usize, owned: boolean = false): %s {", name)
	W("        return new %s(ptr, owned)", name)
	W("    }\n")

	W("    @inline static wrap(buf: ArrayBuffer): %s {", name)
	W("        return new %s(changetype<usize>(buf), false)", name)
	W("    }\n")

	if mut {
		W("    @inline get freeze(): %s {", t.name)
		W("        return %s.unsafe(this.dataStart, this.isOwned)", t.name)
		W("    }\n")
	} else {
		W("    @inline get mut(): %s {", t.mut)
		W("        return %s.unsafe(this.dataStart, this.isOwned)", t.mut)
		W("    }\n")
	}

	W("    @inline @operator('[]') get(i: i32): %s {", t.list.element.name)
	W("        if (i < 0 || i >= this.length) {")
	W("             throw new RangeError()")
	W("        }")
	if c.isPointerType(t.list.element.t) {
		W("        return %s.unsafe(this.dataStart+(i * %d), false)", t.list.element.name, t.list.element.t.Size)
	} else {
		W("        return load<%s>(this.dataStart + (i * %d))", t.list.element.name, t.list.element.t.Size)
	}
	W("    }")

	W("    @inline get length(): i32 {")
	if t.t.Len < 256 {
		W("        return <i32>load<u8>(this.dataStart + %d)", t.t.Size-1)
	} else {
		W("        return <i32>load<u16>(this.dataStart + %d)", t.t.Size-2)
	}
	W("    }")

	W("    @inline get cap(): i32 {")
	W("        return %d", t.t.Len)
	W("    }")

	if mut {
		W("    @inline @operator('[]=') set(i: i32, v: %s): void {", t.list.element.name)
		W("        if (i < 0 || i >= this.length) {")
		W("             return")
		W("        }")
		if c.isPointerType(t.list.element.t) {
			W("        memory.copy(this.dataStart+(i * %d), v.dataStart, %d)", t.list.element.t.Size, t.list.element.t.Size)
		} else {
			W("        store<%s>(this.dataStart + (i * %d), v)", t.list.element.name, t.list.element.t.Size)
		}
		W("    }")

		//W("func (s *%s) setLen(l int) {", t.mut)
		//if t.t.HeaderSize == 1 {
		//	W("    s.l = byte(l)")
		//} else {
		//	//W("    s[%d] = byte(l)", t.t.Size-2)
		//	//W("    s[%d] = byte(l >> 8)", t.t.Size-1)
		//	//W("    *(*uint16)(unsafe.Pointer(&s[%d])) = uint16(l)", t.t.Size-2)
		//	W("    s.l = uint16(l)")
		//}
		//W("}")

		//if c.isPointerType(t.list.element.t) {
		//	W("func (s *%s) Push(v *%s) bool {", t.mut, t.list.element.name)
		//	W("    l := s.Len()")
		//	W("    if l == %d {", t.t.Len)
		//	W("        return false")
		//	W("    }")
		//	W("    s.b[l] = *v")
		//	//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = *v", t.list.element.name, t.t.ItemSize)
		//	W("    s.setLen(l+1)")
		//	W("    return true")
		//	W("}")
		//} else {
		//	W("func (s *%s) Push(v %s) bool {", t.mut, t.list.element.name)
		//	W("    l := s.Len()")
		//	W("    if l == %d {", t.t.Len)
		//	W("        return false")
		//	W("    }")
		//	W("    s.b[l] = v")
		//	//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = *v", t.list.element.name, t.t.ItemSize)
		//	W("    s.setLen(l+1)")
		//	W("    return true")
		//	W("}")
		//}
		//
		//W("// Removes the last item")
		//W("func (s *%s) Pop(v *%s) bool {", t.mut, t.list.element.name)
		//W("    l := s.Len()")
		//W("    if l == 0 {")
		//W("        return false")
		//W("    }")
		//W("    l -= 1")
		//W("    if v != nil {")
		//W("        *v = s.b[l]")
		////W("        *v = *(*%s)(unsafe.Pointer(&s[l * %d]))", t.list.element.name, t.t.ItemSize)
		//W("    }")
		//// Clear last element.
		//if t.list.element.primitive {
		//	W("    s.b[l] = 0")
		//	//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = 0", t.list.element.name, t.t.ItemSize)
		//} else {
		//	W("    s.b[l] = %s{}", t.list.element.name)
		//	//W("    *(*%s)(unsafe.Pointer(&s[l * %d])) = %s{}", t.list.element.name, t.t.ItemSize, t.list.element.name)
		//}
		//W("    s.setLen(l)")
		//W("    return true")
		//W("}")
		//
		//W("// Removes the first item")
		//W("func (s *%s) Shift(v *%s) bool {", t.mut, t.list.element.name)
		//W("    l := s.Len()")
		//W("    if l == 0 {")
		//W("        return false")
		//W("    }")
		//W("    if v != nil {")
		//W("        *v = s.b[0]")
		//W("    }")
		//W("    if l > 1 {")
		//// Shift bytes over
		//W("        copy(s.b[0:], s.b[1:s.l])")
		////W("        copy(s[0:], s[%d:l*%d])", t.t.ItemSize, t.t.ItemSize)
		//W("    }")
		//W("    l -= 1")
		//// Clear last element
		//if t.list.element.primitive {
		//	W("    s.b[l] = 0")
		//	//W("    *(*%s)(unsafe.Pointer(&s[l*%d])) = 0", t.list.element.name, t.t.ItemSize)
		//} else {
		//	W("    s.b[l] = %s{}", t.list.element.name)
		//	//W("    *(*%s)(unsafe.Pointer(&s[l*%d])) = %s{}", t.list.element.name, t.t.ItemSize, t.list.element.name)
		//}
		//W("    s.setLen(l)")
		//W("    return true")
		//W("}")

		W("    @inline clear(): void {")
		W("        memory.fill(this.dataStart, 0, %d)", t.t.Size)
		W("    }")

	} else {

		//W("func (s *%s) Cap() int {", t.name)
		//W("    return %d", t.t.Len)
		//W("}")
		//
		//if t.list.element.st != nil {
		//	W("func (s *%s) MarshalMap(m []map[string]interface{}) []map[string]interface{} {", t.name)
		//	W("    if m == nil {")
		//	W("        m = make([]map[string]interface{}, 0, s.Len())")
		//	W("    }")
		//	W("    for _, v := range s.Unsafe() {")
		//	W("        m = append(m, v.MarshalMap(nil))")
		//	W("    }")
		//	W("    return m")
		//	W("}")
		//}
		//
		//W("func (s *%s) CopyTo(v []%s) []%s {", t.name, t.list.element.name, t.list.element.name)
		//W("    return append(v, s.Unsafe()...)")
		//W("}")
		//
		//W("func (s *%s) Unsafe() []%s {", t.name, t.list.element.name)
		//W("    return s.b[0:s.Len()]")
		////W("    return *(*[]%s)(unsafe.Pointer(&reflect.SliceHeader{", t.list.element.name)
		////W("        Data: uintptr(unsafe.Pointer(&s.b[0])),")
		////W("        Len: s.Len(),")
		////W("        Cap: s.Cap(),")
		////W("    }))")
		//W("}")
		//
		//W("func (s *%s) Bytes() []byte {", t.name)
		//W("    return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{")
		//W("        Data: uintptr(unsafe.Pointer(&s.b[0])),")
		//if t.t.ItemSize > 1 {
		//	W("        Len: s.Len() * %d,", t.t.ItemSize)
		//} else {
		//	W("        Len: s.Len(),")
		//}
		//W("        Cap: %d,", t.t.ItemSize*t.t.Len)
		//W("    }))")
		//W("}")
	}

	//W("func (s *%s) Hash() uint32 {", t.name)
	//W("    return bu.Hash32Bytes(s[0:s.Len()*%d])", t.t.ItemSize)
	//W("}")
	//
	//W("func (s *%s) Hash64() uint64 {", t.name)
	//W("    return bu.Hash64Bytes(s[0:s.Len()*%d])", t.t.ItemSize)
	//W("}")
	W("}")

	return nil
}

func (c *Compiler) genUnion(t *asType, b *Builder) error {
	return nil
}

func (c *Compiler) genString(t *asType, mut bool, b *Builder) error {
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

		W("func (s *String%d) set(v string) {", size)
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

		W("func (s *String%d) Len() int {", size)
		if sizeBytes == 1 {
			W("    return int(s[%d])", sizeIndex)
		} else if sizeBytes == 2 {
			W("    return int(*(*uint16)(unsafe.Pointer(&s[%d]))", sizeIndex)
			//W("    return int(uint16(s[%d]) | uint16(s[%d]) << 8)", sizeIndex, sizeIndex+1)
		}
		W("}")

		W("func (s *String%d) Cap() int {", size)
		W("    return %d", sizeIndex)
		W("}")

		W("func (s *String%d) Unsafe() string {", size)
		W("    return *(*string)(unsafe.Pointer(&reflect.StringHeader{")
		W("        Data: uintptr(unsafe.Pointer(&s[0])),")
		if sizeBytes == 1 {
			W("        Len: int(s[%d]),", sizeIndex)
		} else {
			W("        Len: int(*(*uint16)(unsafe.Pointer(&s[%d])),", sizeIndex)
			//W("        Len: int(uint16(s[%d]) | uint16(s[%d]) << 8),", sizeIndex, sizeIndex+1)
		}
		W("    }))")
		W("}")

		W("func (s *String%d) String() string {", size)
		if sizeBytes == 1 {
			W("    return string(s[0:s[%d]])", sizeIndex)
		} else {
			W("    return string(s[0:s.Len()])")
		}
		W("}")

		W("func (s *String%d) Bytes() []byte {", size)
		W("    return s[0:s.Len()]")
		W("}")

		W("func (s *String%d) Clone() *String%d {", size, size)
		W("    v := String%d{}", size)
		W("    copy(s[0:], v[0:])")
		W("    return &v")
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

func (a *Compiler) asType(t *Type) string {
	switch t.Kind {
	case KindByte:
		return "u8"
	case KindInt8:
		return "i8"
	case KindInt16:
		return "i16"
	case KindUInt16:
		return "u16"
	case KindInt32:
		return "i32"
	case KindUInt32:
		return "u32"
	case KindInt64:
		return "i64"
	case KindUInt64:
		return "u64"
	case KindFloat32:
		return "f32"
	case KindFloat64:
		return "f64"
	case KindString:
		return "string"
	}
	return ""
}
