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
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			// we only care about const declarations (-> return) on the package level
			return true
		}

		curType := ""

		for _, spec := range decl.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				// const is always ValueSpec
				continue
			}

			// find out type of constant
			// (or multiple constants: `const a, b Letter = "a", "b"`)

			if vspec.Type == nil && curType == "" {
				// no type - no interest
				continue
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

			// get the enum name!
			for _, name := range vspec.Names {
				e.Parsed.Names = append(e.Parsed.Names, name.Name)
			}

			if len(vspec.Names) > 1 {
				//
				curType = ""
			}
		}

		return false
	})

	return nil
}

func (e *enumarr) calculateNames() {
	exported := strings.ToUpper(e.TypeName[:1]) + e.TypeName[1:] + "Array"

	if e.ExportFunc {
		// method is exported so no conflicts
		e.FuncName = exported
	}

	if e.ExportVar {
		e.VarName = exported
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