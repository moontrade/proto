package schema

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	ErrUnresolved = errors.New("unresolved")
)

const (
	maxDepth   = 10
	FileSuffix = ".moon"
)

type File struct {
	Dir          string
	Name         string
	Path         string
	Package      string
	Hash         uint64
	Err          error
	contentBytes []byte
	Content      string
	Imports      []*Imports
	Consts       []*Const
	Structs      []*Struct
	Messages     []*Message
	Enums        []*Enum
	Unions       []*Union
	Lists        map[string][]*Type
	Types        map[string]*Type
	ImportMap    map[string]*Import
	Strings      map[string][]*Type
}

type Imports struct {
	Line     Line
	Comments []string
	List     []*Import
}

type Import struct {
	Imports  *Imports
	Parent   *File
	Path     string
	Name     string
	Alias    string
	File     *File
	Comments []string
	Line     Line
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
	case KindBytes:
		return fmt.Sprintf("Bytes%d", t.Len)
	case KindStruct, KindEnum, KindUnion, KindUnknown:
		return strings.ReplaceAll(Capitalize(t.Name), ".", "_")
	case KindList:
		return fmt.Sprintf("%s%dList", f.createTypeName(t.Element, cycle+1), t.Len)
	case KindMap:
		return f.uniqueName(fmt.Sprintf("%s%sMap", f.createTypeName(t.Element, cycle+1), f.createTypeName(t.Value, cycle+1)))
	}
	panic("unknown")
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
	case KindString, KindBytes:
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

		t.Name = f.createTypeName(t, cycle+1)
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
			// Is type imported?
			if field.Type.Import != nil {
				if field.Type.Import.File != nil {
					if err := field.Type.Import.File.resolve(); err != nil {
						return err
					}
					if err := field.Type.Import.File.resolveType(field.Type, cycle+1); err != nil {
						return err
					}
				}
			} else if err := field.Type.File.resolveType(field.Type, cycle+1); err != nil {
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

		fields := make([]*StructField, 0, len(t.Struct.Fields))
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
				fields = append(fields, &StructField{
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
			t.Struct.Fields = append(t.Struct.Fields, &StructField{
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

	case KindMessage:

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
