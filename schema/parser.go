package schema

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	errNotFound = errors.New("not found")
)

type Parser struct {
	current   rune
	lineCount int
	mark      int
	index     int
	content   string
	comments  []string
	file      *File
}

type Expression string

func Parse(path, content string) (*File, error) {
	p := &Parser{
		content: content,
		file: &File{
			Package: PackageName(path),
			Path:    path,
			Types:   make(map[string]*Type),
		},
	}
	return p.Parse()
}

func ParseBytes(path string, b []byte) (*File, error) {
	return Parse(path, string(b))
}

// Returns the next line to parse
func (p *Parser) nextLine() (string, error) {
	mark := p.index
	p.mark = p.index
	for p.index < len(p.content) {
		switch p.content[p.index] {
		case '\n':
			s := p.content[mark:p.index]
			p.index++
			if len(s) > 0 && s[len(s)-1] == '\r' {
				s = s[0 : len(s)-1]
			}
			p.lineCount++
			return s, nil
		}
		p.index++
	}

	line := p.content[mark:]
	if len(line) == 0 {
		return "", io.EOF
	}
	p.lineCount++
	return line, nil
}

func (p *Parser) error(msg string, args ...interface{}) error {
	return fmt.Errorf("%s:%d %s", p.file.Path, p.lineCount, fmt.Sprintf(msg, args...))
}

func (p *Parser) Parse() (*File, error) {
	f := p.file

	var comments []string
	for {
		line, err := p.nextLine()
		if err != nil {
			if err == io.EOF {
				_ = f.resolve()
				return f, nil
			}
			return nil, err
		}
		mark := 0

	loop:
		for i, c := range line {
			switch c {
			// Ignore whitespace
			case ' ', '\t', '\r':
				line = line[i+1:]

			// comment
			case '/':
				line = line[mark:]
				if len(line) == 1 || line[1] != '/' {
					return nil, p.error("expected '/' after first '/'")
				}
				comments = append(comments, line[2:])
				break loop

			// package
			//case 'p':
			//	err := p.parsePackage(line, comments)
			//	if err != nil {
			//		return nil, err
			//	}
			//	comments = nil
			//	break loop

			// import
			case 'i':
				err := p.parseImports(line, comments)
				if err != nil {
					return nil, err
				}
				break loop

			// const
			case 'c':
				cst, err := p.parseConst(line, comments)
				if err != nil {
					return nil, err
				}
				if f.Types == nil {
					f.Types = make(map[string]*Type)
				}
				if existing := f.Types[cst.Name]; existing != nil {
					return nil, p.error(
						fmt.Sprintf("Name '%s' already used on line %d", cst.Name, existing.Line))
				}
				comments = nil
				f.Consts = append(f.Consts, cst)
				f.Types[cst.Name] = cst.Type
				break loop

			// enum
			case 'e':
				enum, err := p.parseEnum(line, comments)
				if err != nil {
					return nil, err
				}

				if f.Types == nil {
					f.Types = make(map[string]*Type)
				}
				if existing := f.Types[enum.Name]; existing != nil {
					return nil, p.error(
						fmt.Sprintf("name '%s' already used on line %d", enum.Name, existing.Line))
				}
				comments = nil
				f.Enums = append(f.Enums, enum)
				f.Types[enum.Name] = enum.Type
				break loop

			// union
			case 'u':
				union, err := p.parseUnion(line, comments)
				if err != nil {
					return nil, err
				}

				if f.Types == nil {
					f.Types = make(map[string]*Type)
				}
				if existing := f.Types[union.Name]; existing != nil {
					return nil, p.error(
						fmt.Sprintf("name '%s' already used on line %d", union.Name, existing.Line))
				}
				comments = nil
				f.Unions = append(f.Unions, union)
				f.Types[union.Name] = union.Type
				break loop

			// struct
			case 's':
				st, err := p.parseStruct(line, comments)
				if err != nil {
					return nil, err
				}

				// Set optionals
				st.setOptionals()

				if f.Types == nil {
					f.Types = make(map[string]*Type)
				}
				if existing := f.Types[st.Name]; existing != nil {
					return nil, p.error(
						fmt.Sprintf("name '%s' already used on line %d", st.Name, existing.Line))
				}
				comments = nil
				f.Structs = append(f.Structs, st)
				f.Types[st.Name] = st.Type

				break loop

			default:
				return nil, p.error(fmt.Sprintf("invalid syntax '%s'", line))
			}
		}
	}
}

func (p *Parser) parseImport(line string, comments []string) (*Import, error) {
	type stateCode int
	const (
		StateStart stateCode = iota
		StatePath
		StatePathAfter
		StateAlias
		StateAfterAlias
		StateComment
	)

	state := StateStart
	mark := 0
	//last := byte(0)

	name := ""
	path := ""
	alias := ""

loop:
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch state {
		case StateStart:
			switch c {
			case ' ', '\t', '\r':
				mark = i
				continue

			case '"':
				state = StatePath
				//last = c
				mark = i + 1

			default:
				if !IsLetter(c) && c != '_' {
					return nil, p.error(fmt.Sprintf("import package name must start with a '\"' or a letter: '%s'", string(c)))
				}
				state = StateAlias
				//last = c
				mark = i
			}

		case StatePath:
			if c == '"' {
				path = line[mark:i]
				name = PackageName(path)
				state = StatePathAfter
				mark = i + 1
			}
			//switch c {
			//case ' ', '\t', '\r':
			//	name = line[mark:i]
			//	state = StatePathAfter
			//	mark = i + 1
			//case '/':
			//	name = line[mark:i]
			//	state = StateComment
			//	mark = i
			//case '.':
			//	if last == '.' {
			//		return nil, p.error(fmt.Sprintf("import package name must start with a letter not '%s'", string(c)))
			//	}
			//default:
			//	if last == '.' {
			//		if !IsLetter(c) {
			//			return nil, p.error(fmt.Sprintf("import package name must start with a letter at each level: cannot start with %s", string(c)))
			//		}
			//	} else if !IsNumeral(c) && !IsLetter(c) && c != '_' {
			//		return nil, p.error(fmt.Sprintf("invalid import package character '%s'", string(c)))
			//	}
			//}
			//last = c

		case StatePathAfter:
			switch c {
			case ' ', '\t', '\r':
				mark = i + 1

			case '/':
				state = StateComment
				mark = i

			default:
				if c != '_' && !IsLetter(c) {
					return nil, p.error(fmt.Sprintf("invalid import package alias starting character '%s'", string(c)))
				}
				state = StateAlias
				mark = i
			}

		case StateAlias:
			switch c {
			case ' ', '\t', '\r':
				alias = line[mark:i]
				mark = i + 1
				state = StateAfterAlias

			case '/':
				alias = line[mark:i]
				state = StateComment
				mark = i

			default:
				if c != '_' && !IsLetter(c) && !IsNumeral(c) {
					return nil, p.error(fmt.Sprintf("invalid import package alias character '%s'", string(c)))
				}
			}

		case StateAfterAlias:
			switch c {
			case ' ', '\t', '\r':
			case '"':
				state = StatePath
				mark = i + 1
			default:
				return nil, p.error(fmt.Sprintf(
					"invalid import package character after alias when expecting '\"': '%s'", string(c)))
			}

		case StateComment:
			switch c {
			case '/':
				comments = append(comments, line[i+1:])
				break loop
			default:
				return nil, p.error("expected comment")
			}
		}
	}

	switch state {
	case StatePath:
		name = line[mark:]
	case StateAlias:
		alias = line[mark:]
	case StateComment:
		return nil, p.error("expected comment")
	}

	if len(alias) == 0 {
		alias = name
	}
	return &Import{
		File: p.file,
		Line: Line{
			Number: p.lineCount,
			Begin:  p.mark,
			End:    p.index,
		},
		Name:     name,
		Path:     path,
		Alias:    alias,
		Comments: comments,
	}, nil
}

func (p *Parser) parseImports(line string, comments []string) error {
	if len(line) < 7 || line[1:7] != "mport " {
		return p.error("expected 'import' keyword")
	}

	line = strings.TrimSpace(line[7:])
	if len(line) == 0 {
		return p.error("expected import declaration")
	}

	imports := &Imports{
		Line: Line{
			Number: p.lineCount,
			Begin:  p.mark,
			End:    p.index,
		},
	}
	p.file.Imports = append(p.file.Imports, imports)

	addImport := func(line string, comments []string) error {
		imp, err := p.parseImport(line, comments)
		if err != nil {
			return err
		}
		imp.Imports = imports
		if p.file.ImportMap == nil {
			p.file.ImportMap = make(map[string]*Import)
		}
		if p.file.ImportMap[imp.Name] != nil {
			return p.error("import already declared: %s", imp.Name)
		}
		if existing := p.file.ImportMap[imp.Alias]; existing != nil {
			return p.error("import alias '%s' clashes with import on line %d", imp.Alias, existing.Line)
		}
		p.file.ImportMap[imp.Name] = imp
		p.file.ImportMap[imp.Alias] = imp
		imports.List = append(imports.List, imp)

		return nil
	}

	if line != "(" {
		return addImport(line, comments)
	}

	imports.Comments = comments

	type stateCode int
	const (
		StateStart stateCode = iota
		StateComment
	)

	comments = nil
	importCount := 0
	var err error
	for {
		line, err = p.nextLine()
		if err != nil {
			return err
		}

		if len(line) == 0 {
			continue
		}

		state := StateStart
		mark := 0

	loop:
		for i, c := range line {
			switch state {
			case StateStart:
				switch c {
				case ' ', '\t', '\r':
					mark = i + 1
					continue
				case '/':
					state = StateComment
					mark = i + 1
				case ')':
					if importCount == 0 {
						return p.error("empty import declaration")
					}
					return nil
				default:
					if err := addImport(line[mark:], comments); err != nil {
						return err
					}
					importCount++
					comments = nil
					break loop
				}

			case StateComment:
				switch c {
				case '/':
					comments = append(comments, line[i+1:])
					break loop
				default:
					return p.error("expected comment")
				}
			}
		}
	}
}

func parseAliasAndName(s string) (alias string, name string, err error) {
	type stateCode int
	const (
		StateBegin stateCode = iota
		StateDotOrEnd
		StateAfterDot
		StateName
	)

	state := StateBegin
	mark := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch state {
		case StateBegin:
			if IsWhitespace(c) {
				mark = i + 1
				continue
			}
			if IsLetter(c) || c == '_' {
				state = StateDotOrEnd
				continue
			}

			return "", "", errors.New("expected letter or _ as first character")

		case StateDotOrEnd:
			if c == '.' {
				state = StateAfterDot
				alias = s[mark:i]
				mark = i + 1
				continue
			}
			if !IsLetter(c) && c != '_' && !IsNumeral(c) {
				return "", "", fmt.Errorf("invalid type name character: %s", string(c))
			}

		case StateAfterDot:
			if !IsLetter(c) && c != '_' {
				return "", "", fmt.Errorf("invalid type name character: %s", string(c))
			}
			state = StateName

		case StateName:
			if !IsLetter(c) && !IsNumeral(c) && c != '_' {
				return "", "", fmt.Errorf("invalid type name character: %s", string(c))
			}
		}
	}

	switch state {
	case StateDotOrEnd:
		return alias, s[mark:], nil

	case StateAfterDot:
		return "", "", errors.New("expected type name after '.'")
	}

	return alias, s[mark:], nil
}

func parseStringLen(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

func (p *Parser) parseType(line string, comments []string) (t *Type, err error) {
	mark := 0
	count := 0

	t = &Type{
		File:     p.file,
		Comments: comments,
		Line: Line{
			Number: p.lineCount,
			Begin:  p.mark,
			End:    p.index,
		},
	}

	type stateCode int
	const (
		StateName stateCode = iota
		StateSize
		StateAfterName
		StateMapKeyword
		StateComment
		StateValue
		StateValueLiteral
		StateMapLiteral
		StateListLiteral
		StateTrueLiteral
		StateTrueLiteralEnd
		StateFalseLiteral
		StateFalseLiteralEnd
		StateStringLiteral
		StateNumberLiteral
		StateNilLiteral
		StateNilLiteralEnd
		StateMaybeComment
		StateValueLiteralMaybeComment
	)

	state := StateName

	parseName := func(name string) error {
		if len(name) == 0 {
			return p.error("type Name required")
		}

		if !t.Optional {
			t.Optional = name[0] == '?'
			if t.Optional {
				name = name[1:]
			}
			if name[len(name)-1] == '?' {
				if t.Optional {
					return p.error("optional '?' declared in front already")
				}
				t.Optional = true
				name = name[0 : len(name)-1]
			}
		}
		if len(name) == 0 {
			return p.error("type Name required")
		}

		t.Name = name
		if StartsWith(name, "string") || StartsWith(name, "bytes") {
			var (
				kind   Kind
				length int
				err    error
			)

			if StartsWith(name, "string") {
				kind = KindString
				length, err = parseStringLen(name[6:])
			} else {
				kind = KindBytes
				length, err = parseStringLen(name[5:])
			}

			if err != nil {
				return p.error("invalid string declaration: %s", err.Error())
			}

			switch t.Kind {
			case KindList:
				if t.Optional {
					return p.error("list elements or map keys cannot be optional")
				}
				t.Element = &Type{
					File: p.file,
					Name: name,
					Kind: kind,
					Len:  length,
				}

			case KindMap:
				t.Value = &Type{
					File: p.file,
					Name: name,
					Kind: kind,
					Len:  length,
				}

			default:
				t.Name = name
				t.Kind = kind
				t.Len = length
			}
			state = StateAfterName
		} else {
			kind := KindOf(name)

			var imp *Import
			if kind == KindUnknown {
				var err error
				var alias string
				alias, name, err = parseAliasAndName(name)
				if err != nil {
					return err
				}
				if len(alias) > 0 {
					imp = p.file.ImportMap[alias]
					if imp == nil {
						return p.error("no import for alias: %s", alias)
					}
				}
			}

			switch t.Kind {
			case KindUnknown:
				t.Name = name
				t.Kind = kind
				t.Import = imp

			case KindList:
				t.Element = &Type{
					File:   p.file,
					Name:   name,
					Kind:   kind,
					Import: imp,
				}
			case KindMap:
				t.Value = &Type{
					File:   p.file,
					Name:   name,
					Kind:   kind,
					Import: imp,
				}
			default:
				t.Name = name
				t.Kind = kind
				t.Import = imp
			}
			state = StateAfterName
		}
		return nil
	}

	for i, c := range line {
		switch state {
		// Type
		case StateName:
			switch c {
			case '?':
				if t.Optional {
					return nil, p.error("optional token '?' already specified")
				}
				mark = i + 1
				t.Optional = true

			case '[':
				mark = i + 1
				state = StateSize

			case ' ', '\t', '\n':
				if count > 0 {
					name := line[mark:i]
					if err := parseName(name); err != nil {
						return nil, err
					}
					mark = i + 1
					count = 0
				} else {
					// Skip whitespace
					mark = i + 1
				}
			default:
				count++
			}

		// Array declaration '['
		case StateSize:
			if t.Len != 0 {
				return nil, p.error("size already defined")
			}
			switch c {
			case ' ', '\t', '\n':
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case ']':
				length, err := strconv.Atoi(strings.TrimSpace(line[mark:i]))
				if err != nil {
					return nil, p.error(fmt.Sprintf("invalid list length value %s", err.Error()))
				}
				t.Len = length
				if t.Len <= 0 {
					return nil, p.error(fmt.Sprintf("invalid list length %d", length))
				}
				t.Kind = KindList
				mark = i + 1
				state = StateName
			}

		// Equals
		case StateAfterName:
			switch c {
			case '=':
				mark = i + 1
				state = StateValue
				count = 0

			case '/':
				state = StateComment
				mark = i
				count = 0

			case '-':
				if t.Kind == KindMap {
					return nil, p.error("maps cannot be nested")
				}
				if t.Kind != KindList {
					return nil, p.error("maps require max size '[8] int64 -> int64")
				}
				state = StateMapKeyword
				mark = i
				count = 0

			case ' ', '\t', '\n', '\r':
				mark = i + 1

			default:
				return nil, p.error(fmt.Sprintf("invalid character '%s' expected '=' or whitespace or '/'", string(c)))
			}

		case StateComment:
			if c != '/' {
				return nil, p.error("expected second comment character '/'")
			}
			t.Comments = append(t.Comments, line[i+1:])
			return t, nil

		case StateMapKeyword:
			if c != '>' {
				return nil, p.error("expected map keyword '->'")
			}
			t.Kind = KindMap
			mark = i + 1
			count = 0
			state = StateName

		// Value
		case StateValue:
			switch c {
			case ' ', '\t', '\n', '\r':
				mark = i + 1
				continue

			case 'n':
				state = StateNilLiteral
				mark = i
				count = 0

			case '[':
				if t.Kind != KindList {
					return nil, p.error("list value literals can only be used on list types")
				}
				return nil, p.error("list literals are not supported yet")

			case '{':
				if t.Kind != KindMap && t.Kind != KindStruct {
					return nil, p.error("map value literals can only be used on map types")
				}
				return nil, p.error("map literals are not supported yet")

			case '\'', '"':
				if t.Kind != KindString {
					return nil, p.error("string literals can only be used on string types")
				}
				state = StateStringLiteral
				mark = i
				count = 0

			case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				state = StateNumberLiteral
				mark = i
				count = 0

			case 't':
				state = StateTrueLiteral
				mark = i
				count = 0

			case 'f':
				state = StateFalseLiteral
				mark = i
				count = 0

			case '/':
				return nil, p.error("expected value declaration not comment '//'")

			default:
				state = StateValueLiteral
				mark = i
				count = 0
			}

		case StateTrueLiteral:
			switch c {
			case 'r':
				if count != 0 {
					return nil, p.error("expected 'true' literal")
				}
				count++
			case 'u':
				if count != 1 {
					return nil, p.error("expected 'true' literal")
				}
				count++
			case 'e':
				if count != 2 {
					return nil, p.error("expected 'true' literal")
				}
				count = 0
				t.Init = true
				state = StateTrueLiteralEnd
			}

		case StateTrueLiteralEnd:
			switch c {
			case '/':
				if t.Kind != KindBool {
					return nil, p.error("only bool excepts value of 'true'")
				}
				t.Init = true
				state = StateComment
			case ' ', '\t', '\r':
				if t.Kind != KindBool {
					return nil, p.error("only bool excepts value of 'true'")
				}
				t.Init = true
				state = StateMaybeComment

			default:
				return nil, p.error("expected 'nil'")
			}

		case StateFalseLiteral:
			switch c {
			case 'a':
				if count != 0 {
					return nil, p.error("expected 'false' literal")
				}
				count++
			case 'l':
				if count != 1 {
					return nil, p.error("expected 'false' literal")
				}
				count++
			case 's':
				if count != 2 {
					return nil, p.error("expected 'false' literal")
				}
				count++
			case 'e':
				if count != 3 {
					return nil, p.error("expected 'false' literal")
				}
				count = 0
				state = StateFalseLiteralEnd
			}

		case StateFalseLiteralEnd:
			switch c {
			case '/':
				if t.Kind != KindBool {
					return nil, p.error("only bool excepts value of 'false'")
				}
				t.Init = false
				state = StateComment
			case ' ', '\t', '\r':
				if t.Kind != KindBool {
					return nil, p.error("only bool excepts value of 'false'")
				}
				t.Init = false
				state = StateMaybeComment

			default:
				return nil, p.error("expected 'nil'")
			}

		case StateNilLiteral:
			switch c {
			case 'i':
				if count != 0 {
					return nil, p.error("invalid value literal")
				}
				count = 1

			case 'l':
				if count != 1 {
					return nil, p.error("invalid value literal")
				}
				count = 2
				state = StateNilLiteralEnd
			}

		case StateNilLiteralEnd:
			switch c {
			case '/':
				if !t.Optional {
					return nil, p.error("only optional values accept nil values")
				}
				t.Init = Nil{}
				state = StateComment
			case ' ', '\t', '\r':
				t.Init = Nil{}
				state = StateMaybeComment

			default:
				return nil, p.error("expected 'nil'")
			}

		case StateStringLiteral:
			switch c {
			case '\'', '"':
				init := line[mark : i+1]
				if init[0] != init[len(init)-1] {
					return nil, p.error("mismatching begin and end quotation marks")
				}
				t.Init = init[1 : len(init)-1]
				state = StateMaybeComment
				mark = i
			}

		case StateNumberLiteral:
			switch c {
			case ' ', '\t', '\r':
				literal := line[mark:i]
				if strings.Index(literal, ".") > -1 {
					switch t.Kind {
					case KindFloat32, KindFloat64:
						val, err := strconv.ParseFloat(literal, 64)
						if err != nil {
							return nil, p.error(fmt.Sprintf("failed to parse float literal: %s", err.Error()))
						}
						t.Init = val

					default:
						return nil, p.error("type cannot be set to a float")
					}
				} else {
					switch t.Kind {
					case KindByte, KindUInt16, KindUInt32, KindUInt64:
						val, err := strconv.ParseUint(literal, 10, 64)
						if err != nil {
							return nil, p.error(fmt.Sprintf("failed to parse unsigned integer literal: %s", err.Error()))
						}
						t.Init = val

					case KindInt8, KindInt16, KindInt32, KindInt64:
						val, err := strconv.ParseInt(literal, 10, 64)
						if err != nil {
							return nil, p.error(fmt.Sprintf("failed to parse integer literal: %s", err.Error()))
						}
						t.Init = val

					default:
						return nil, p.error("type cannot be set to an integer")
					}
				}
			}

		case StateMaybeComment:
			switch c {
			case ' ', '\t', '\r':
				continue
			case '/':
				state = StateComment
				mark = i
			}

		case StateValueLiteral:
			switch c {
			case '/':
				state = StateValueLiteralMaybeComment
			}

		case StateValueLiteralMaybeComment:
			switch c {
			case '/':
				t.Init = Expression(line[mark:i])
				state = StateComment
				mark = i + 1
				count = 0
			default:
				state = StateValueLiteral
			}
		}
	}

	switch state {
	case StateName:
		if err := parseName(line[mark:]); err != nil {
			return nil, err
		}
		return t, nil

	case StateTrueLiteral:
		return nil, p.error("expected 'true'")

	case StateTrueLiteralEnd:
		t.Init = true

	case StateFalseLiteral:
		return nil, p.error("expected 'false'")

	case StateFalseLiteralEnd:
		t.Init = false

	case StateNilLiteral:
		return nil, p.error("expected 'nil'")

	case StateNilLiteralEnd:
		if !t.Optional {
			return nil, p.error("only optional types can be set nil")
		}
		t.Init = Nil{}

	case StateStringLiteral:
		return nil, p.error("incomplete string literal missing end '\"'")

	case StateNumberLiteral:
		literal := line[mark:]
		t.Init = literal
		if strings.Index(literal, ".") > -1 {
			switch t.Kind {
			case KindFloat32, KindFloat64:
				val, err := strconv.ParseFloat(literal, 64)
				if err != nil {
					return nil, p.error(fmt.Sprintf("failed to parse float literal: %s", err.Error()))
				}
				t.Init = val

			default:
				return nil, p.error("type cannot be set to a float")
			}
		} else {
			var err error
			t.Init, err = ParseInt(t.Kind, literal)
			if err != nil {
				return nil, p.error(fmt.Sprintf("failed to parse integer literal: %s", err.Error()))
			}
		}

	case StateValueLiteral, StateValueLiteralMaybeComment:
		t.Init = Expression(line[mark:])

	}

	return t, nil
}

func (p *Parser) parsePackage(line string, comments []string) error {
	if len(line) < 7 || line[1:7] != "ackage" {
		return p.error("expected 'package' keyword")
	}
	line = line[7:]
	if len(line) == 0 || !IsWhitespace(line[0]) {
		return p.error("expected a whitespace after 'package'")
	}
	line = line[1:]

	p.file.Package = strings.TrimSpace(line)
	p.comments = comments

	return nil
}

func (p *Parser) parseConst(line string, comments []string) (*Const, error) {
	if len(line) < 5 || line[1:5] != "onst" {
		return nil, p.error("expected 'const' keyword")
	}
	line = line[5:]
	if len(line) == 0 || !IsWhitespace(line[0]) {
		return nil, p.error("expected a whitespace after 'const'")
	}
	line = line[1:]
	if len(line) == 0 {
		return nil, p.error("expected const declaration")
	}

	cst := &Const{}
	state := 0
	mark := 0

	for i := 0; i < len(line); i++ {
		c := line[i]

		switch state {
		case 0:
			// Read until space
			switch c {
			case ' ', '\t', '\n':
			default:
				mark = i
				state = 1
			}

			//
		case 1:
			// Read until space
			switch c {
			case ' ', '\t', '\n':
				cst.Name = line[mark:i]
				mark = i + 1
				state = 2
			default:
			}

		case 2:
			// Read until space
			switch c {
			case ' ', '\t', '\n':
				mark = i + 1
			default:
				t, err := p.parseType(line[mark:], comments)
				if err != nil {
					return nil, p.error(err.Error())
				}
				if t.Optional {
					return nil, p.error("const cannot be optional")
				}
				if t.Init == nil {
					return nil, p.error("const must have a value")
				}
				cst.Type = t
				cst.Type.Const = cst
				return cst, nil
			}
		}
	}
	return nil, nil
}

func (p *Parser) parseStruct(line string, comments []string) (*Struct, error) {
	if len(line) < 6 || line[1:6] != "truct" {
		return nil, p.error("expected 'struct' keyword")
	}
	line = line[6:]

	if len(line) == 0 || !IsWhitespace(line[0]) {
		return nil, p.error("expected a whitespace after 'struct'")
	}
	line = line[1:]

	type stateCode int
	const (
		StateName stateCode = iota
		StateNumberOrName
		StateNumber
		StateNumberAfter
		StateCurlyBrace
		StateComment
		StateEnd
	)
	state := StateName

	mark := 0
	count := 0
	st := &Struct{
		Type: &Type{
			Line: Line{
				Number: p.lineCount,
				Begin:  p.mark,
				End:    p.index,
			},
			File:     p.file,
			Kind:     KindStruct,
			Comments: comments,
		},
	}
	st.Type.Struct = st

	for i, c := range line {
		switch state {
		case StateName:
			switch c {
			case ' ', '\t', '\r':
				if count == 0 {
					mark = i + 1
					continue
				}
				st.Name = line[mark:i]
				st.Type.Name = st.Name
				state = StateCurlyBrace
				mark = i + 1
				count = 0

			case '{':
				if count == 0 {
					return nil, p.error("expected struct Name")
				}
				state = StateEnd

			default:
				count++
			}

		case StateCurlyBrace:
			switch c {
			case ' ', '\t', '\r':
				// Skip whitespace
			case '{':
				state = StateEnd

			default:
				return nil, p.error("expected '{'")
			}

		case StateEnd:
			return nil, p.error("expected EOL")
		}

	}

	comments = nil
	var err error

	for {
		line, err = p.nextLine()
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '/':
			if len(line) == 1 || line[1] != '/' {
				return nil, p.error("comment expected")
			}
			comments = append(comments, line[2:])

		case '[', '@':
			return nil, p.error("attributes not supported yet")

		case '}':
			return st, nil

		default:
			state = StateNumberOrName
			mark = 0
			count = 0
			field := &Field{
				Struct: st,
			}
		loop:
			for i := 0; i < len(line); i++ {
				c := line[i]
				switch state {
				case StateNumberOrName:
					switch c {
					case ' ', '\t', '\r':
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					default:
						if !IsLetter(c) {
							return nil, p.error("expected field number or name")
						}
						state = StateName
						mark = i
					}

				case StateNumber:
					switch c {
					case ' ', '\t', '\r':
						var num uint64
						num, err = strconv.ParseUint(line[mark:i], 10, 64)
						if err != nil {
							return nil, p.error("invalid field number '%s': %s", line[mark:i], err.Error())
						}
						field.Number = int(num)
						state = StateName
						mark = i + 1
					}

				case StateNumberAfter:
					switch c {
					case ' ', '\t', '\r':
						mark = i + 1
						continue

					default:
						if !IsLetter(c) && c != '_' {
							return nil, p.error("invalid first character for field name: %s", string(c))
						}
						mark = i
						state = StateName
					}

				case StateName:
					switch c {
					case ' ', '\t', '\r':
						field.Name = line[mark:i]
						t, err := p.parseType(line[i:], comments)
						if err != nil {
							return nil, err
						}

						field.Type = t
						st.Fields = append(st.Fields, field)
						t.Field = field

						state = StateEnd
						comments = nil
						break loop
					}

				case StateComment:
					if c != '/' {
						return nil, p.error("expected comment")
					}
					comments = append(comments, line[i+2:])
					break loop
				}
			}
		}
	}
}

func (p *Parser) parseEnum(line string, comments []string) (*Enum, error) {
	if len(line) < 4 || line[1:4] != "num" {
		return nil, p.error("expected 'enum' keyword")
	}
	line = line[4:]
	if len(line) == 0 || !IsWhitespace(line[0]) {
		return nil, p.error("expected a whitespace after 'enum'")
	}
	line = line[1:]

	type stateCode int
	const (
		StateName stateCode = iota
		StateColon
		StateTypeName
		StateCurlyBrace
		StateComment
		StateEquals
		StateValue
		StateMaybeComment
		StateEnd
	)
	state := StateName

	mark := 0
	count := 0
	enum := &Enum{}
	enum.Type = &Type{
		Line: Line{
			Number: p.lineCount,
			Begin:  p.mark,
			End:    p.index,
		},
		File:     p.file,
		Kind:     KindEnum,
		Enum:     enum,
		Comments: comments,
	}

	for i, c := range line {
		switch state {
		case StateName:
			switch c {
			case ' ', '\t', '\r':
				if count == 0 {
					mark = i + 1
					continue
				}
				enum.Name = line[mark:i]
				enum.Type.Name = enum.Name
				state = StateColon
				mark = i + 1
				count = 0

			case '{':
				if count == 0 {
					return nil, p.error("expected enum Name and type")
				}
				state = StateEnd

			default:
				count++
			}

		case StateColon:
			switch c {
			case ' ', '\t', '\r':
				continue

			case ':':
				state = StateTypeName
				mark = i + 1
				count = 0

			case '{':
				if count == 0 {
					return nil, p.error("expected enum type 'enum Name : byte'")
				}
				state = StateEnd

			default:
				return nil, p.error("expected enum type 'enum Name : byte'")
			}

		case StateTypeName:
			switch c {
			case ' ', '\t', '\r', '{':
				if count == 0 {
					mark = i + 1
					continue
				}
				name := line[mark:i]

				if StartsWith(name, "string") {
					length, err := parseStringLen(name[6:])
					if err != nil {
						return nil, p.error("invalid string declaration: %s", err.Error())
					}
					enum.Type.Size = length
					enum.Type.Element = &Type{
						File: p.file,
						Name: name,
						Kind: KindString,
					}
				} else {
					kind := KindOf(name)
					switch kind {
					case KindByte, KindInt8, KindUInt16, KindInt16, KindUInt32, KindInt32, KindUInt64, KindInt64:
						enum.Type.Size = kind.Size()
						enum.Type.Element = &Type{
							File: p.file,
							Name: name,
							Kind: kind,
						}
					default:
						return nil, p.error(fmt.Sprintf("enum must be an integer type not '%s'", name))
					}
				}
				if c == '{' {
					state = StateEnd
				} else {
					state = StateCurlyBrace
					mark = i + 1
					count = 0
				}

			default:
				count++
			}

		case StateCurlyBrace:
			switch c {
			case ' ', '\t', '\r':
				// Skip whitespace
			case '{':
				state = StateEnd

			default:
				return nil, p.error("expected '{'")
			}

		case StateEnd:
			return nil, p.error("expected EOL")
		}
	}

	comments = nil
	var err error

	for {
		line, err = p.nextLine()
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '/':
			if len(line) == 1 || line[1] != '/' {
				return nil, p.error("comment expected")
			}
			comments = append(comments, line[2:])

		case '[':
			return nil, p.error("attributes not supported yet")

		case '}':
			return enum, nil

		default:
			state = StateName
			mark = 0
			count = 0

			option := &EnumOption{
				Enum: enum,
				Line: Line{
					Number: p.lineCount,
					Begin:  p.mark,
					End:    p.index,
				},
			}
		loop:
			for i, c := range line {
				switch state {
				case StateName:
					switch c {
					case ' ', '\t', '\r':
						if count == 0 {
							mark = i + 1
							continue
						}
						option.Name = line[mark:i]
						state = StateEquals
						mark = i + 1
						count = 0

					case '=':
						state = StateValue
						mark = i + 1
						count = 0

					default:
						count++
					}

				case StateComment:
					if c != '/' {
						return nil, p.error("expected comment")
					}
					comments = append(comments, line[i+2:])
					state = StateEnd
					break loop

				case StateEquals:
					switch c {
					case ' ', '\t', '\r':

					case '=':
						state = StateValue
						mark = i + 1
						count = 0

					default:
						return nil, p.error("expected '='")
					}

				case StateValue:
					switch c {
					case ' ', '\t', '\r':
						if count == 0 {
							mark = i + 1
							continue
						}
						option.Value, err = ParseInt(enum.Type.Element.Kind, line[mark:])
						if err != nil {
							return nil, p.error("invalid integer type for option")
						}
						state = StateMaybeComment

					case '/':
						if count == 0 {
							return nil, p.error("expected an integer value for option")
						}
						option.Value, err = ParseInt(enum.Type.Element.Kind, line[mark:])
						if err != nil {
							return nil, p.error("invalid integer type for option")
						}
						state = StateComment

					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						count++

					default:
						return nil, p.error("expected an integer value")
					}

				case StateMaybeComment:
					switch c {
					case ' ', '\t', '\r':
					case '/':
						state = StateComment
					}
				}
			}

			switch state {
			case StateName:
				return nil, p.error("expected Name")
			case StateComment:
				return nil, p.error("expected comment")
			case StateEquals:
				return nil, p.error("expected '=' for option")
			case StateValue:
				if count == 0 {
					return nil, p.error("expected a value for option")
				}
				option.Comments = comments
				option.Value, err = ParseInt(enum.Type.Element.Kind, line[mark:])
				if err != nil {
					return nil, p.error("invalid integer type for option")
				}
			}

			enum.Options = append(enum.Options, option)
			comments = nil
		}
	}
}

func (p *Parser) parseUnion(line string, comments []string) (*Union, error) {
	if len(line) < 5 || line[1:5] != "nion" {
		return nil, p.error("expected 'union' keyword")
	}
	line = line[5:]

	if len(line) == 0 || !IsWhitespace(line[0]) {
		return nil, p.error("expected a whitespace after 'union'")
	}
	line = line[1:]

	type stateCode int
	const (
		StateName stateCode = iota
		StateCurlyBrace
		StateComment
		StateType
		StateEnd
	)
	state := StateName

	mark := 0
	count := 0
	union := &Union{
		Comments: comments,
	}

	for i, c := range line {
		switch state {
		case StateName:
			switch c {
			case ' ', '\t', '\r':
				if count == 0 {
					mark = i + 1
					continue
				}
				union.Name = line[mark:i]
				union.Type = &Type{
					Line: Line{
						Number: p.lineCount,
						Begin:  p.mark,
						End:    p.index,
					},
					File:     p.file,
					Kind:     KindUnion,
					Name:     union.Name,
					Comments: comments,
					Union:    union,
				}
				comments = nil
				state = StateCurlyBrace
				mark = i + 1
				count = 0

			case '{':
				if count == 0 {
					return nil, p.error("expected union Name and type")
				}
				state = StateEnd

			default:
				count++
			}

		case StateCurlyBrace:
			switch c {
			case ' ', '\t', '\r':
				// Skip whitespace
			case '{':
				state = StateEnd

			default:
				return nil, p.error("expected '{'")
			}

		case StateEnd:
			return nil, p.error("expected EOL")
		}
	}

	comments = nil
	var err error

	for {
		line, err = p.nextLine()
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '/':
			if len(line) == 1 || line[1] != '/' {
				return nil, p.error("comment expected")
			}
			comments = append(comments, line[2:])

		case '[':
			return nil, p.error("attributes not supported yet")

		case '}':
			return union, nil

		default:
			state = StateName
			mark = 0
			count = 0

			option := &UnionOption{
				Union: union,
			}
		loop:
			for i, c := range line {
				switch state {
				case StateName:
					switch c {
					case ' ', '\t', '\r':
						if count == 0 {
							mark = i + 1
							continue
						}
						option.Name = line[mark:i]
						state = StateType
						mark = i + 1
						count = 0

					case '/':
						if count > 0 {
							return nil, p.error("invalid use of comment character '/'")
						}
						mark = i
						count = 0
						state = StateComment

					default:
						count++
					}

				case StateComment:
					if c != '/' {
						return nil, p.error("expected comment")
					}
					comments = append(comments, line[i+2:])
					state = StateEnd
					break loop

				case StateType:
					switch c {
					case ' ', '\t', '\r':
						continue

					default:
						t, err := p.parseType(line[i:], comments)
						if err != nil {
							return nil, p.error(fmt.Sprintf("union option on line %d has invalid type declaration: %s", p.lineCount, err.Error()))
						}
						option.Type = t
						option.Type.UnionOption = option

						if t.Init != nil {
							return nil, p.error("union option cannot have an initializer")
						}

						comments = nil
						state = StateEnd
					}
				}
			}

			switch state {
			case StateName:
				return nil, p.error("expected Name")
			case StateComment:
				return nil, p.error("expected comment")
			case StateType:
				return nil, p.error("expected a type for option")
			}

			union.Options = append(union.Options, option)
			comments = nil
		}
	}
}
