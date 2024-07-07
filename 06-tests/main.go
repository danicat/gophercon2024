package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	filename := "div.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatalln(err)
	}

	astutil.Apply(file, nil, changeBinaryOperator)

	out, err := os.Create("mutations/" + filename)
	if err != nil {
		log.Fatalln(err)
	}

	printer.Fprint(out, fset, file)

	results, err := runTests(".")
	if err != nil {
		log.Fatalln(err)
	}

	for _, test := range results {
		if test.Test != "" && (test.Action == "pass" || test.Action == "fail") {
			fmt.Printf("%20v: %v\n", test.Test, test.Action)
		}
	}
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

func runTests(pkgDir string) ([]testEvent, error) {
	res, err := exec.Command("go", "test", "--json", pkgDir).CombinedOutput()
	if err != nil {
		log.Println(err)
	}

	return parseTestOutput(res)
}

type testEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func parseTestOutput(output []byte) ([]testEvent, error) {
	var result []testEvent
	list := "[" + strings.ReplaceAll(string(output[:len(output)-1]), "\n", ",") + "]"
	err := json.Unmarshal([]byte(list), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}