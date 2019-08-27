package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	taskName = flag.String("name", "", "Run task name")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of `go run main.go -name {TASK_NAME}`\n")
	flag.PrintDefaults()
}

type Generator struct {
	buf bytes.Buffer
	pkg *Package
}

type Package struct {
	dir   string
	name  string
	files []File
}

type File struct {
	name    string
	astFile *ast.File
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	if len(*taskName) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	g := newGenerator()
	src, err := g.Format()
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("gen.go", src, 0644); err != nil {
		panic(err)
	}
}

func newGenerator() *Generator {
	g := &Generator{
		pkg: &Package{
			dir:  *taskName,
			name: "main",
			files: []File{
				{
					name:    "main.go",
					astFile: basicAST,
				},
			},
		},
	}
	g.Printf("// Code generated by gen-task; DO NOT EDIT\n")
	g.Printf("\n")
	return g
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) Format() ([]byte, error) {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		return nil, err
	}
	return src, nil
}

var basicAST = &ast.File{
	Name: ast.NewIdent("main"),
	Decls: []ast.Decl{
		&ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote("fmt"),
					},
				},
			},
		},
		&ast.FuncDecl{
			Name: ast.NewIdent("main"),
			Type: &ast.FuncType{},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fmt"),
								Sel: ast.NewIdent("Println"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: strconv.Quote("===> RUN TASK"),
								},
							},
						},
					},
				},
			},
		},
	},
}
