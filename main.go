package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"go/printer"
	"go/ast"
	"strings"
	"os"

	"github.com/Sirupsen/logrus"
)

const (
	ReplaceTag = "// @replace-tag "
)

func replaceTagsInFile(input string) {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing protobuf definitions
	f, err := parser.ParseFile(fset, input, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	replacedTags := false
	for _, d := range f.Decls {
		// Get General Declaration
		gDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gDecl.Specs {
			// Get Type Declaration
			tSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// Get Struct Declaration
			sDecl, ok := tSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range sDecl.Fields.List {
				if field.Doc == nil {
					continue
				}
				for _, comment := range field.Doc.List {
					if strings.Contains(comment.Text, ReplaceTag) {
						// Found a field
						commentTag := strings.TrimPrefix(comment.Text, ReplaceTag)
						commentTagKey := strings.Split(commentTag, ":")[0]
						fieldTagTokens := strings.Split(strings.Trim(field.Tag.Value, "`"), " ")
						foundFlt := ""
						for _, flt := range fieldTagTokens {
							if strings.HasPrefix(flt, commentTagKey) {
								// Found the tag
								foundFlt = flt
							}
						}
						if foundFlt == "" {
							continue
						}
						fmt.Printf("Replacing tag (%v) -> (%v) \n", foundFlt, commentTag)
						r := strings.NewReplacer(foundFlt, commentTag)
						field.Tag.Value = r.Replace(field.Tag.Value)
						replacedTags = true
					}
				}
			}
		} 
	}
	if replacedTags {
		var buf bytes.Buffer
		printer.Fprint(&buf, fset, f)
		f, err := os.Create(input)
		if err != nil {
			fmt.Println("Unable to write protobuf files: ", err)
			return
		}
		f.WriteString(buf.String())
	}
}

func main() {
	var input string

	flag.StringVar(&input, "input", "", "path to input protobuf file")
	flag.Parse()

	if len(input) == 0 {
		logrus.Errorf("input file not provided")
		return
	}
	replaceTagsInFile(input)
}
