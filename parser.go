package wap

import (
	"fmt"
	"io"
)

type Loader struct {
	schema Schema

	files []*file
}

func (p *Loader) resolve() {}

type file struct {
	path   string
	schema *Schema
}

type Parser struct {
	current   rune
	lineCount int
	mark      int
	index     int
	content   string
	comments  []string
	schema    Schema
	file      file
}

type Expression string

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
	return fmt.Errorf("%s:%d %s", p.file.path, p.lineCount, fmt.Sprintf(msg, args...))
}

func (p *Parser) Parse() (*file, error) {
	f := p.file
	_ = f

	var comments []string
	for {
		line, err := p.nextLine()
		if err != nil {
			if err == io.EOF {
				//_ = f.resolve()
				//return f, nil
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
			//case 'i':
			//	err := p.parseImports(line, comments)
			//	if err != nil {
			//		return nil, err
			//	}
			//	break loop
			//
			//// block
			//case 'b':
			//
			//
			//// const
			//case 'c':
			//	cst, err := p.parseConst(line, comments)
			//	if err != nil {
			//		return nil, err
			//	}
			//	if f.Types == nil {
			//		f.Types = make(map[string]*Type)
			//	}
			//	if existing := f.Types[cst.Name]; existing != nil {
			//		return nil, p.error(
			//			fmt.Sprintf("Name '%s' already used on line %d", cst.Name, existing.Line))
			//	}
			//	comments = nil
			//	f.Consts = append(f.Consts, cst)
			//	f.Types[cst.Name] = cst.Type
			//	break loop
			//
			//// enum
			//case 'e':
			//	enum, err := p.parseEnum(line, comments)
			//	if err != nil {
			//		return nil, err
			//	}
			//
			//	if f.Types == nil {
			//		f.Types = make(map[string]*Type)
			//	}
			//	if existing := f.Types[enum.Name]; existing != nil {
			//		return nil, p.error(
			//			fmt.Sprintf("name '%s' already used on line %d", enum.Name, existing.Line))
			//	}
			//	comments = nil
			//	f.Enums = append(f.Enums, enum)
			//	f.Types[enum.Name] = enum.Type
			//	break loop
			//
			//// union
			//case 'u':
			//	union, err := p.parseUnion(line, comments)
			//	if err != nil {
			//		return nil, err
			//	}
			//
			//	if f.Types == nil {
			//		f.Types = make(map[string]*Type)
			//	}
			//	if existing := f.Types[union.Name]; existing != nil {
			//		return nil, p.error(
			//			fmt.Sprintf("name '%s' already used on line %d", union.Name, existing.Line))
			//	}
			//	comments = nil
			//	f.Unions = append(f.Unions, union)
			//	f.Types[union.Name] = union.Type
			//	break loop
			//
			//// struct
			//case 's':
			//	st, err := p.parseStruct(line, comments)
			//	if err != nil {
			//		return nil, err
			//	}
			//
			//	// Set optionals
			//	st.setOptionals()
			//
			//	if f.Types == nil {
			//		f.Types = make(map[string]*Type)
			//	}
			//	if existing := f.Types[st.Name]; existing != nil {
			//		return nil, p.error(
			//			fmt.Sprintf("name '%s' already used on line %d", st.Name, existing.Line))
			//	}
			//	comments = nil
			//	f.Structs = append(f.Structs, st)
			//	f.Types[st.Name] = st.Type
			//
			//	break loop

			// message
			case 'm':
				return nil, p.error(fmt.Sprintf("invalid syntax '%s'", line))

			// record
			case 'r':
				return nil, p.error(fmt.Sprintf("invalid syntax '%s'", line))

			default:
				return nil, p.error(fmt.Sprintf("invalid syntax '%s'", line))
			}
		}
	}
}
