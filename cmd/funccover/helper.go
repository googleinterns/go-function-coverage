//Copyright 2020 Google LLC

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"text/template"
	"time"
)

// returns the source code representation of a AST file
func astToByte(fset *token.FileSet, f *ast.File) []byte {
	// return the source code as a string
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)
	return buf.Bytes()
}

// writes necessary set instrucions for instrumentation to cover variable
// suffix is the string that will be added to the end of the cover variable
// filename is the argument that passed into coverage collection functions (can be changed later)
func addCounters(w io.Writer, content []byte, suffix, fileName string, currentIndex int) (int, error) {

	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	var contentLength = len(content)
	var events []int
	var mainRbrace = -1

	// Iterate over the functions to find the positions to insert instructions
	// Save the positions to insert to the events array
	// mainRbrace will be used to insert defer to the end of main
	// Positions are changed due to imports but saved information in funcCover will be the same with the source code
	for _, decl := range parsedFile.Decls {
		switch t := decl.(type) {
		// Function Decleration
		case *ast.FuncDecl:
			events = append(events, int(t.Body.Lbrace))
			if t.Name.Name == "main" {
				mainRbrace = int(t.Body.Rbrace) - 1
			}
		}
	}

	// Writes the instrumented code using w io.Writer
	// Insert set instructions to the functions
	// f() {
	// 	cover_hash.Counts[funcNumber] = 1;
	// 	...
	// }
	// Also inserts defercover_hash.Collect(args) to the main
	// func main {
	// 	...
	//	defer cover_hash.Collect(args)
	// }

	eventIndex := 0

	for i := 0; i < contentLength; i++ {
		if eventIndex < len(events) && i == events[eventIndex] {
			fmt.Fprintf(w, "\ncover_%s.Counts[%v] = true;", suffix, currentIndex)
			currentIndex++
			eventIndex++
		}
		if i == mainRbrace {
			fmt.Fprintf(w, "\n\tdefer cover_%s.Collect(\"%s\")\n", suffix, fileName)
		}
		fmt.Fprintf(w, "%s", string(content[i]))
	}

	return currentIndex, nil
}

// writes the declaration of cover variable to the end of the main source file using go templates
func declCover(w io.Writer, suffix string, fileName string, period time.Duration, funcCover FuncCover) {

	funcTemplate, err := template.New("cover functions and variables").Parse(declTmpl)

	if err != nil {
		panic(err)
	}

	usePeriod := 0

	if period > 0 {
		usePeriod = 1
	}

	var declParams = struct {
		Suffix     string
		UsePeriod  int
		Period     string
		FuncCount  int
		FuncBlocks []FuncCoverBlock
		FileName   string
	}{suffix, usePeriod, fmt.Sprint(period), len(funcCover.FuncBlocks), funcCover.FuncBlocks, fileName}

	err = funcTemplate.Execute(w, declParams)

	if err != nil {
		panic(err)
	}
}

var declTmpl = `
var cover_{{.Suffix}} = covcollect.Cover {
	Lines: []uint32{ {{range .FuncBlocks}}
		{{.Line}},{{end}}
	},
	Names: []string{ {{range .FuncBlocks}}
		"{{.Name}}",{{end}}
	},
	Counts: []bool{ {{range .FuncBlocks}}
		false,{{end}}
	},
}
{{ if eq .UsePeriod 1 }}
func init() {
    go cover_{{.Suffix}}.PeriodicalCollect("{{.Period}}",  "{{.FileName}}")
}{{end}}`
