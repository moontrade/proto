package schema

import (
	"errors"
	"fmt"
)

type Schema struct {
	//Root     *Dir
	Files    map[string]*File
	Errors   []error
	resolved bool
}

type VirtualDir struct {
	Name  string        `json:"name"`
	Path  string        `json:"path"`
	Dirs  []VirtualDir  `json:"dirs"`
	Files []VirtualFile `json:"files"`
}

type VirtualFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Data []byte `json:"data"`
}

type Dir struct {
	Parent *Dir
	Name   string
	Path   string
	Files  map[string]*File
	Dirs   map[string]*Dir
}

//type FileData struct {
//	Dir    *Dir
//	Name   string
//	Path   string
//	Data   []byte
//	Err    error
//	Parsed *File
//}

//func LoadVirtualDir(root *VirtualDir) (*Schema, error) {
//	if root == nil {
//		return nil, errors.New("nil")
//	}
//
//	schema := &Schema{
//		Root: &Dir{
//			Name:  root.Name,
//			Path:  root.Path,
//			Dirs:  make(map[string]*Dir),
//			Files: make(map[string]*File),
//		},
//		Files: make(map[string]*File),
//	}
//
//	addDir := func(parent *Dir) {
//
//	}
//}
//
//func LoadVirtualFiles(root []VirtualFile) (*Schema, error) {
//	if root == nil {
//		return nil, errors.New("nil")
//	}
//
//	schema := &Schema{
//		Root: &Dir{
//			Dirs:  make(map[string]*Dir),
//			Files: make(map[string]*File),
//		},
//		Files: make(map[string]*File),
//	}
//
//	addDir := func(parent *Dir) {
//
//	}
//}

// LoadFromFS loads a schema from the filesystem optionally resolving.
func LoadFromFS(dirOrFile string, resolve bool) (*Schema, error) {
	schema, err := loadFromFS(dirOrFile)
	if err != nil {
		return nil, err
	}
	if !resolve {
		return schema, nil
	}
	if err = schema.Resolve(); err != nil {
		return nil, err
	}
	return schema, nil
}

func (pa *Schema) Resolve() error {
	if len(pa.Errors) > 0 {
		return pa.Errors[0]
	}
	if len(pa.Files) == 0 {
		return errors.New("no files")
	}
	if pa.resolved {
		return nil
	}

	sorted := make([]*File, 0, len(pa.Files))

	pa.resolved = true
	// Resolve imports
	for _, f := range pa.Files {
		if f != nil && len(f.Imports) > 0 {
		OUTER:
			for _, imps := range f.Imports {
				for _, imp := range imps.List {
					p := RelativePath(f.Path, imp.Path)
					if len(p) == 0 {
						f.Err = fmt.Errorf("import '%s' could not be resolved", imp.Path)
						pa.Errors = append(pa.Errors, f.Err)
						break OUTER
					}

					importFile := pa.Files[p]
					if importFile == nil {
						f.Err = fmt.Errorf("import '%s' could not be resolved in file: %s", imp.Path, f.Path)
						pa.Errors = append(pa.Errors, f.Err)
						break OUTER
					}

					sorted = append(sorted, importFile)

					imp.File = importFile
					imp.Path = p
					imp.Parent = f
				}
			}

			sorted = append(sorted, f)
		}
	}

	// Resolve types
	for _, f := range sorted {
		if f != nil {
			if err := f.resolve(); err != nil {
				pa.Errors = append(pa.Errors, err)
				f.Err = err
			}
		}
	}

	if len(pa.Errors) > 0 {
		return pa.Errors[0]
	}

	return nil
}
