package main

import (
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	filename := "div.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		log.Fatalln(err)
	}

	out, err := os.Create("mutations/" + filename)
	if err != nil {
		log.Fatalln(err)
	}

	printer.Fprint(out, fset, file)
}
