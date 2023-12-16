package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strings"
)

// enumarr parses enums from a file and generates a file with enum array.
type enumarr struct {
	Parsed Parsed

	// config
	TypeName   string
	ExportVar  bool
	ExportFunc bool
	Files      []string
	Output     string

	// calculated
	VarName  string
	FuncName string
}

type Parsed struct {
	Pkg   string
	Names []string
}

func (e *enumarr) Run() error {
	for _, fileName := range e.Files {
		if err := e.parse(fileName); err != nil {
			return err
		}
	}

	e.calculateNames()

	err := e.generate()
	if err != nil {
		return err
	}

	return nil
}

func (e *enumarr) parse(fileName string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.File:
			e.Parsed.Pkg = t.Name.Name
			return true
		case *ast.GenDecl:
			if t.Tok != token.CONST {
				return true
			}

			curType := ""

			for _, spec := range t.Specs {
				vspec, ok := spec.(*ast.ValueSpec)
				if !ok {
					// const is always ValueSpec
					continue
				}

				// find out type of constant
				// (or multiple constants: `const a, b Letter = "a", "b"`)

				if vspec.Type == nil && len(vspec.Values) > 0 {
					// prev decl's type is inferred only if value not specified
					curType = ""
				}

				if vspec.Type != nil {
					ident, ok := vspec.Type.(*ast.Ident)
					if !ok {
						continue
					}

					curType = ident.Name
				}

				if curType != e.TypeName {
					// not interested in this type
					continue
				}

				// get the enum names!
				for _, name := range vspec.Names {
					e.Parsed.Names = append(e.Parsed.Names, name.Name)
				}

				if len(vspec.Names) > 1 {
					//
					curType = ""
				}
			}

			return false
		default:
			// we are interested only in constants that are declared on package level
			return false

		}
	})

	return nil
}

func (e *enumarr) calculateNames() {
	exported := strings.ToUpper(e.TypeName[:1]) + e.TypeName[1:]

	if e.ExportFunc {
		// method is exported so no conflicts
		e.FuncName = exported + "All"
	}

	if e.ExportVar {
		e.VarName = exported + "Array"
	} else {
		e.VarName = "_" + e.TypeName + "Array"
	}
}

func (e *enumarr) generate() error {
	tmpl, err := template.New("arr.tmpl").Parse(enumArrTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(e.Output)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, e)
	if err != nil {
		return err
	}

	return nil
}
