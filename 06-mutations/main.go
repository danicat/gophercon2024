package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	filename := "div.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatalln(err)
	}

	astutil.Apply(file, nil, reverseIfCond)

	out, err := os.Create("mutations/" + filename)
	if err != nil {
		log.Fatalln(err)
	}

	printer.Fprint(out, fset, file)
}

func reverseIfCond(c *astutil.Cursor) bool {
	n := c.Node()
	switch x := n.(type) {
	case *ast.IfStmt:
		x.Cond = &ast.UnaryExpr{
			Op: token.NOT,
			X:  x.Cond,
		}
		return false
	}
	return true
}

func changeBinaryOperator(c *astutil.Cursor) bool {
	n := c.Node()
	switch x := n.(type) {
	case *ast.BinaryExpr:
		if x.Op == token.QUO {
			x.Op = token.MUL
			return false
		}
	}
	return true
}
