package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/tools/go/ast/astutil"
)

func AutoImmutableColumns(pkgPath string, immutableColumns []string) error {
	fmt.Println("change mutable field set")
	fset := token.NewFileSet()
	return filepath.Walk(
		pkgPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".go" {
				// Parse the file and create an AST
				file, err := parser.ParseFile(
					fset, path, nil, parser.ParseComments,
				)
				if err != nil {
					return err
				}
				// Rewrite ast tree
				immutableFields(file, immutableColumns)
				outputFile, err := os.Create(path)
				if err != nil {
					return err
				}
				defer outputFile.Close()

				// Write the ast tree as code to file
				err = format.Node(outputFile, fset, file)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func immutableFields(file *ast.File, immutableFields []string) {
	astutil.Apply(
		file, nil, func(c *astutil.Cursor) bool {
			if i, ok := c.Node().(*ast.Ident); ok {
				if i.Name != "mutableColumns" {
					return true
				}
				if i.Obj != nil && i.Obj.Kind == ast.Var {
					if vs, ok := i.Obj.Decl.(*ast.ValueSpec); ok {
						if len(vs.Values) != 1 {
							return true
						}
						list, ok := vs.Values[0].(*ast.CompositeLit)
						if !ok {
							return true
						}
						newList := make([]ast.Expr, 0)
						for _, v := range list.Elts {
							if v, ok := v.(*ast.Ident); ok && !slices.Contains(immutableFields, v.Name) {
								newList = append(newList, v)
							}
						}
						list.Elts = newList
					}
				}
			}
			return true
		},
	)
}
