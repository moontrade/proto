package schema

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrUnresolved = errors.New("unresolved")
)

const (
	maxDepth   = 10
	FileSuffix = ".wb"
)

type Config struct {
	Path string
}

type Schema struct {
	Config  *Config
	RootDir string
	Files   map[string]*File
	Strings map[string][]*Type
	Lists   map[string][]*Type

	byPackage map[string]*File
}

func (c *Schema) load(path string, count int) (*File, error) {
	if count > 10 {
		return nil, fmt.Errorf("dependency cycle: %s", path)
	}
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		c.Files[path] = &File{
			Err: err,
		}
		return nil, err
	}

	if existing := c.Files[path]; existing != nil {
		return existing, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		c.Files[path] = &File{
			Err: err,
		}
		return nil, err
	}

	dir := filepath.Dir(path)
	if existing := c.byPackage[dir]; existing != nil {
		return existing, fmt.Errorf("only one '%s' file per directory: '%s' has '%s' and '%s'", FileSuffix, dir, existing.Path, path)
	}

	content := string(buf)
	file, err := Parse(path, content)
	if err != nil {
		if file == nil {
			return nil, err
		}
		file.Err = err
	}
	if file == nil {
		return nil, errors.New("nil file")
	}

	file.Path = path
	c.byPackage[dir] = file
	c.Files[path] = file

	if file.Strings != nil {
		for k, l := range file.Strings {
			existing := c.Strings[k]
			existing = append(existing, l...)
			c.Strings[k] = existing
		}
	}

	if file.GlobalLists != nil {
		for k, l := range file.GlobalLists {
			existing := c.Lists[k]
			existing = append(existing, l...)
			c.Lists[k] = existing
		}
	}

	if len(file.Imports) > 0 {
		for _, imps := range file.Imports {
			for _, imp := range imps.List {
				p := RelativePath(file.Path, imp.Path)
				if len(p) == 0 {
					file.Err = fmt.Errorf("import '%s' could not be resolved", imp.Path)
					return nil, err
				}
				imp.Parent, err = c.load(p, count+1)
				if err != nil {
					file.Err = fmt.Errorf("import '%s' could not be resolved: %s", imp.Path, err.Error())
					return nil, err
				}
				imp.Path = p
				//imp.Parent = file
			}
		}
		//for _, imps := range file.Imports {
		//	for _, imp := range imps.List {
		//		if imp.Parent == nil {
		//			file.Err = fmt.Errorf("import '%s' could not be resolved", imp.Path)
		//			return nil, err
		//		}
		//		err = imp.File.resolve()
		//		if err != nil {
		//			file.Err = fmt.Errorf("import '%s' resolve error: %s", imp.Path, err.Error())
		//			return nil, err
		//		}
		//	}
		//}
	}

	return file, nil
}

func NewSchema(config *Config) (*Schema, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}
	if len(config.Path) == 0 {
		return nil, errors.New("no Files or directories to process")
	}

	root, err := os.Stat(config.Path)
	if err != nil {
		return nil, err
	}
	rootPath := root.Name()

	c := &Schema{
		Config:    config,
		Files:     make(map[string]*File),
		Strings:   make(map[string][]*Type),
		Lists:     make(map[string][]*Type),
		byPackage: make(map[string]*File),
	}

	walk := func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, FileSuffix) {
			return nil
		}
		_, err = c.load(path, 0)
		return err
	}

	// Load schema files
	if root.IsDir() {
		if err := filepath.Walk(config.Path, walk); err != nil {
			return nil, err
		}
	} else if strings.HasSuffix(rootPath, FileSuffix) {
		if err := walk(rootPath, root, nil); err != nil {
			return nil, err
		}
	}

	for _, file := range c.Files {
		_ = file.resolve()
	}
	for _, file := range c.Files {
		if err = file.resolve(); err != nil {
			return nil, err
		}
	}

	// Find root package directory
	maxSeparatorsCount := 0
	rootDir := ""
	for _, file := range c.Files {
		cnt := strings.Count(file.Path, string(os.PathSeparator))
		if len(rootDir) == 0 {
			rootDir = file.Path
			maxSeparatorsCount = cnt
			continue
		}
		if maxSeparatorsCount < cnt {
			maxSeparatorsCount = cnt
			rootDir = file.Path
		}
	}
	// Find common ancestor directory
	for {
		rootDir = filepath.Dir(rootDir)
		if len(rootDir) == 0 || rootDir == "/" {
			break
		}

		cnt := 0
		for _, file := range c.Files {
			if strings.Index(file.Path, rootDir) == 0 {
				cnt++
			}
		}
		if cnt == len(c.Files) {
			break
		}
	}

	c.RootDir = rootDir

	for _, file := range c.Files {
		file.Package = filepath.Dir(file.Path[len(c.RootDir)+1:])
	}

	return c, nil
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

func isGlobal(t *Type, cycle int) bool {
	if t == nil {
		return false
	}
	switch t.Kind {
	case KindStruct, KindEnum, KindUnion, KindUnknown:
		return false
	case KindList:
		return isGlobal(t.Element, cycle+1)
	case KindMap:
		return isGlobal(t.Element, cycle+1) && isGlobal(t.Value, cycle+1)
	default:
		return true
	}
}

func (f *File) uniqueName(n string) string {
	for {
		if f.Types[n] == nil {
			return n
		}
		n = n + "_"
	}
}

func (f *File) createTypeName(t *Type, cycle int) string {
	if cycle >= maxDepth {
		panic("cyclic")
	}
	switch t.Kind {
	case KindBool:
		return "Bool"
	case KindByte:
		return "U8"
	case KindInt8:
		return "I8"
	case KindInt16:
		return "I16"
	case KindUInt16:
		return "U16"
	case KindInt32:
		return "I32"
	case KindUInt32:
		return "U32"
	case KindInt64:
		return "I64"
	case KindUInt64:
		return "U64"
	case KindFloat32:
		return "F32"
	case KindFloat64:
		return "F64"
	case KindString:
		return fmt.Sprintf("String%d", t.Len)
	case KindStruct, KindEnum, KindUnion, KindUnknown:
		return strings.ReplaceAll(Capitalize(t.Name), ".", "_")
	case KindList:
		return fmt.Sprintf("%s%dList", f.createTypeName(t.Element, cycle+1), t.Len)
	case KindMap:
		return f.uniqueName(fmt.Sprintf("%s%sMap", f.createTypeName(t.Element, cycle+1), f.createTypeName(t.Value, cycle+1)))
	}
	panic("unknown")
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

func (f *File) resolveType(t *Type, cycle int) error {
	if t == nil {
		return errNotFound
	}
	if t.Resolved {
		return nil
	}
	if cycle >= maxDepth {
		return fmt.Errorf("%s:%d cyclic dependency: %s", t.File.Path, t.Line.Number, t.Name)
	}

	switch t.Kind {
	case KindByte, KindInt8, KindBool:
		t.Name = f.createTypeName(t, 0)
		t.Size = 1
		t.Resolved = true
	case KindInt16, KindUInt16:
		t.Name = f.createTypeName(t, 0)
		t.Size = 2
		t.Resolved = true
	case KindInt32, KindUInt32, KindFloat32:
		t.Name = f.createTypeName(t, 0)
		t.Size = 4
		t.Resolved = true
	case KindInt64, KindUInt64, KindFloat64:
		t.Name = f.createTypeName(t, 0)
		t.Size = 8
		t.Resolved = true
	case KindString:
		t.Name = f.createTypeName(t, 0)
		t.Size = t.Len
		t.Resolved = true
		if f.Strings == nil {
			f.Strings = make(map[string][]*Type)
		}
		f.Strings[t.Name] = append(f.Strings[t.Name], t)
		//t.Size = Align(t.Size)

	case KindPad:
		t.Resolved = true
		return nil

	case KindList:
		if t.Element == nil {
			return errNotFound
		}
		err := t.Element.File.resolveType(t.Element, cycle+1)
		if err != nil {
			return err
		}
		if t.Len == 0 {
			return fmt.Errorf("%s:%d lists must specify a length greater than 0", f.Path, t.Line.Number)
		}

		t.Name = f.createTypeName(t, 0)
		t.Size = t.Element.Size * t.Len
		t.ItemSize = t.Element.Size
		if t.Len <= 255 {
			t.Size += 1
			t.HeaderSize = 1
		} else if t.Len <= math.MaxUint16 {
			t.Size += 2
			t.HeaderSize = 2
		} else {
			return fmt.Errorf("%s:%d lists cannot have more than %d elements", f.Path, t.Line.Number, math.MaxUint16)
		}
		aligned := Align(t)
		if aligned > t.Size {
			t.Padding = aligned - t.Size
			t.Size = aligned
		}
		t.HeaderOffset = t.Size - t.HeaderSize - 1
		//if isGlobal(t, 0) {
		//	if f.GlobalLists == nil {
		//		f.GlobalLists = make(map[string][]*Type)
		//	}
		//	f.GlobalLists[t.Name] = append(f.GlobalLists[t.Name], t)
		//} else {
		//	if f.Lists == nil {
		//		f.Lists = make(map[string][]*Type)
		//	}
		//	f.Lists[t.Name] = append(f.Lists[t.Name], t)
		//}
		if f.Lists == nil {
			f.Lists = make(map[string][]*Type)
		}
		f.Lists[t.Name] = append(f.Lists[t.Name], t)
		t.Resolved = true

	case KindMap:
		if t.Element == nil {
			return errNotFound
		}
		if err := t.Element.File.resolveType(t.Element, cycle+1); err != nil {
			return err
		}
		if t.Value == nil {
			return errNotFound
		}
		if err := t.Value.File.resolveType(t.Value, cycle+1); err != nil {
			return err
		}
		t.Name = f.createTypeName(t, 0)
		t.HeaderSize = MapHeaderSize
		t.ItemSize = MapItemHeaderSize + t.Element.Size + t.Value.Size
		t.Size = MapHeaderSize + (t.Len * t.ItemSize)
		t.Size = Align(t)
		t.Resolved = true

	case KindEnum:
		if t.Enum == nil {
			return fmt.Errorf("%s:%d invalid state: type of enum had a nil *Enum", f.Path, t.Line.Number)
		}
		if t.Element == nil {
			return fmt.Errorf("%s:%d enum '%s' did not specify a type: (e.g. enum Code : byte)",
				f.Path, t.Line, t.Name)
		}
		if err := t.File.resolveType(t.Element, cycle+1); err != nil {
			return err
		}
		t.Size = t.Element.Size
		if t.Init == nil {
			t.Resolved = true
			return nil
		}

		if err := resolveEnumInit(t, t.Init); err != nil {
			return err
		}
		t.Resolved = true
		return nil

	case KindStruct:
		if t.Struct == nil {
			return fmt.Errorf("%s:%d invalid state: type of struct had a nil struct", f.Path, t.Line.Number)
		}
		for _, field := range t.Struct.Fields {
			if err := field.Type.File.resolveType(field.Type, cycle+1); err != nil {
				return err
			}
		}

		t.Struct.setOptionals()
		t.Padding = 0
		if len(t.Struct.Optionals) == 0 {
			t.Padding = 0
		} else if len(t.Struct.Optionals) < 8 {
			t.Padding = 1
		} else {
			t.Padding = len(t.Struct.Optionals) / 8
			if len(t.Struct.Optionals)%8 > 0 {
				t.Padding++
			}
		}

		fields := make([]*Field, 0, len(t.Struct.Fields))
		//if t.Padding > 0 {
		//	fields = append(fields, &Field{
		//		Number: -1,
		//		Struct: t.Struct,
		//		Name:   "",
		//		Type: &Type{
		//			File:     f,
		//			Kind:     KindPad,
		//			Resolved: true,
		//			Size:     t.Padding,
		//			Struct:   t.Struct,
		//		},
		//		Offset:    0,
		//		OptOffset: 0,
		//		OptMask:   0,
		//	})
		//}

		t.HeaderSize = t.Padding
		t.Size = t.Padding
		offset := t.Padding
		for _, field := range t.Struct.Fields {
			if field.Type.Kind == KindPad {
				offset += field.Type.Size
				t.Size += field.Type.Size
				fields = append(fields, field)
				continue
			}
			alignTo := FieldAlign(field.Type.Size)

			pad := 0
			diff := offset % alignTo
			if diff > 0 {
				pad = alignTo - diff
			}

			// Add padding?
			if pad > 0 {
				t.Padding += pad
				fields = append(fields, &Field{
					Number: -1,
					Struct: t.Struct,
					Name:   "",
					Type: &Type{
						Line:     field.Type.Line,
						File:     f,
						Kind:     KindPad,
						Resolved: true,
						Size:     pad,
						Struct:   t.Struct,
					},
					Offset:    offset,
					OptOffset: 0,
					OptMask:   0,
				})

				offset += pad
				t.Size += pad
			}

			field.Offset = offset
			offset += field.Type.Size
			t.Size += field.Type.Size

			fields = append(fields, field)
		}
		t.Struct.Fields = fields
		aligned := Align(t)
		if aligned > t.Size {
			pad := aligned - t.Size
			t.Struct.Fields = append(t.Struct.Fields, &Field{
				Number: -1,
				Struct: t.Struct,
				Name:   "",
				Type: &Type{
					File:     f,
					Kind:     KindPad,
					Resolved: true,
					Size:     pad,
					Struct:   t.Struct,
				},
				Offset:    offset,
				OptOffset: 0,
				OptMask:   0,
			})
			t.Padding += pad
			t.Size = aligned
		}
		t.Resolved = true

	case KindUnion:
		if t.Union == nil {
			return fmt.Errorf("%s:%d invalid state: type of union had a nil *Union", f.Path, t.Line.Number)
		}
		if len(t.Union.Options) > 255 {
			return fmt.Errorf("%s:%d unions can have a maximum of 255 options: %d were declared",
				f.Path, t.Line, len(t.Union.Options))
		}
		t.Size = 0
		for _, option := range t.Union.Options {
			if err := option.Type.File.resolveType(option.Type, cycle+1); err != nil {
				return err
			}
			if t.Size < option.Type.Size {
				t.Size = option.Type.Size
			}
		}
		t.Resolved = true

	case KindUnknown:
		var found *Type
		if t.Import != nil {
			// Look for import
			imp := f.ImportMap[t.Import.Name]
			if imp == nil {
				return fmt.Errorf("could not find import for package: %s", t.Import.Name)
			}
			if imp.Parent == nil {
				return errNotFound
			}
			found = imp.Parent.Types[t.Name]
			if found == nil {
				return fmt.Errorf("%s:%d type not found: %s.%s", t.File.Path, t.Line.Number, imp.Alias, t.Name)
			}
		} else {
			found = f.Types[t.Name]
			if found == nil {
				return fmt.Errorf("%s:%d type not found: %s", t.File.Path, t.Line.Number, t.Name)
			}
		}

		if err := found.File.resolveType(found, cycle+1); err != nil {
			return err
		}

		// Copy
		optional := t.Optional
		init := t.Init
		imp := t.Import
		file := t.File
		*t = *found
		t.File = file
		t.Optional = optional
		t.Import = imp

		if init != nil {
			switch t.Kind {
			case KindEnum:
				if err := resolveEnumInit(t, init); err != nil {
					return err
				}
			}
		}
		t.Resolved = false
		if err := f.resolveType(t, cycle+1); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func resolveEnumInit(t *Type, init interface{}) error {
	if t == nil || init == nil || t.Kind != KindEnum {
		return nil
	}
	enum := t.Enum
	if enum == nil {
		return nil
	}

	var option *EnumOption
	var name string
	switch v := init.(type) {
	case Nil:
		t.Init = v
	case *EnumOption:
		option = v
	case Expression:
		name = string(v)
	case string:
		name = v
	}
	if option != nil {
		t.Init = option
		return nil
	}
	option = enum.GetOption(name)
	if option != nil {
		return nil
	}
	return fmt.Errorf("%s:%d invalid enum option: %s:%d %s does not have an option named: %s",
		t.File.Path, t.Line, t.Enum.Type.File.Path, t.Enum.Type.Line, t.Enum.Name, name)
}

func (f *File) resolve() error {
	var errs []error
	for _, t := range f.Types {
		if err := f.resolveType(t, 0); err != nil {
			errs = append(errs, err)
		} else {
			//delete(f.resolveTypes, t)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	// Consts
	//	// Resolve value literals as consts
	//	for resolveType := range f.resolveTypes {
	//		field := resolveType.Field
	//		if field == nil {
	//			continue
	//		}
	//		switch v := resolveType.Init.(type) {
	//		case Expression:
	//			t := f.Types[string(v)]
	//			if t == nil {
	//				return errors.New(fmt.Sprintf("field '%s' on line %d type '%s' not found",
	//					field.Name, field.Line, string(v)))
	//			}
	//			if t.Const == nil {
	//				return errors.New(fmt.Sprintf("field '%s' on line %d id '%s' is not a const",
	//					field.Name, field.Line, string(v)))
	//			}
	//			if field.Type.Kind != t.Const.Type.Kind {
	//				return errors.New(fmt.Sprintf("field '%s' on line %d const '%s' type mismatch: %s <> %s",
	//					field.Name, field.Line, string(v), field.Type.Name, t.Const.Type.Name))
	//			}
	//			field.Type.Init = t.Const
	//		}
	//	}
	//}

	//if len(f.resolveTypes) > 0 {
	//	return ErrUnresolved
	//}
	return nil
}
