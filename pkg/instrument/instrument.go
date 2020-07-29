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

package instrument

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

	"golang.org/x/tools/go/ast/astutil"
)

//  returns the source code representation of a AST file
func astToByte(fset *token.FileSet, f *ast.File) []byte {
	// return the source code as a string
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)
	return buf.Bytes()
}

// writes necessary counters for instrumentation using w
// suffix is the suffix string that will be added to the end of variables and functions
func addCounters(w io.Writer, content []byte, suffix string) {

	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		panic(err)
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

	var currentIndex = 0

	// Writes the instrumented code using w io.Writer
	// Insert set instructions to the functions
	// f() {
	// 	funcCover_hash.Count[funcNumber] = 1;
	// 	...
	// }
	// Also insert defer retrieve_coverage_data_hash() to the main
	// func main {
	// 	...
	//	defer retrieve_coverage_data_hash()
	// }
	for i := 0; i < contentLength; i++ {
		if currentIndex < len(events) && i == events[currentIndex] {
			fmt.Fprintf(w, "\n\tfuncCover_%s.Count[%v] = 1;", suffix, currentIndex)
			currentIndex++
		}
		if i == mainRbrace {
			fmt.Fprintf(w, "\n\tdefer retrieve_coverage_data_%s()\n", suffix)
		}
		fmt.Fprintf(w, string(content[i]))
	}

}

// writes the declaration of funcCover variable and necessery functions
// to the end of the file using go templates
func declCover(w io.Writer, src, suffix, out string, period time.Duration, funcCover FuncCover) {

	funcTemplate, err := template.New("cover functions and variables").Parse(declTmpl)

	if err != nil {
		panic(err)
	}

	usePeriod := 0

	if period > 0 {
		usePeriod = 1
	}

	var declParams = struct {
		Src        string
		Suffix     string
		Output     string
		UsePeriod  int
		Period     string
		FuncCount  int
		FuncBlocks []FuncCoverBlock
	}{src, suffix, out, usePeriod, fmt.Sprint(period), len(funcCover.FuncBlocks), funcCover.FuncBlocks}

	err = funcTemplate.Execute(w, declParams)

	if err != nil {
		panic(err)
	}
}

// Annotate instruments given content (source code) with given parameters and writes it using w
func Annotate(w io.Writer, content []byte, src, suffix, outputFile string, period time.Duration) {

	funcCover := SaveFuncs(content)

	fset := token.NewFileSet() // positions are relative to fset

	parsedFile, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Ensure necessary imports are inserted
	imports := []string{"bufio", "fmt", "os", "time"}
	for _, impr := range imports {
		astutil.AddImport(fset, parsedFile, impr)
	}

	content = astToByte(fset, parsedFile)

	// Add necessary counters to the functions and defer that calls retrieve_coverage_data_$hash() to the main
	// Write the instrumented source with w
	addCounters(w, content, suffix)

	// Write necessary functions variables and an init function that calls periodical_retrieve_$hash() with w
	declCover(w, src, suffix, outputFile, period, funcCover)
}

var declTmpl = `
var funcCover_{{.Suffix}} = struct {
	Count     [{{.FuncCount}}]uint32
	Line      [{{.FuncCount}}]uint32
	Name      [{{.FuncCount}}]string
} {
	Line: [{{.FuncCount}}]uint32{ {{range .FuncBlocks}}
		{{.Line}},{{end}}
	},
	Name: [{{.FuncCount}}]string{ {{range .FuncBlocks}}
		"{{.Name}}",{{end}}
	},
}
{{ if eq .UsePeriod 1 }}
func init() {
    go periodical_retrieve_{{.Suffix}}()
}

func periodical_retrieve_{{.Suffix}}() {

	period, err := time.ParseDuration("{{.Period}}")

	if err != nil {
    	panic(err)
	}

    ticker := time.NewTicker(period)

    for _ = range ticker.C {
        retrieve_coverage_data_{{.Suffix}}()
    } 

}
{{ end }}
func retrieve_coverage_data_{{.Suffix}}() {

	fd, err := os.Create("{{.Output}}")
	if err != nil {
    	panic(err)
	}
	
	w := bufio.NewWriter(fd)

	defer func() { 
    	w.Flush()
    	fd.Close()
	}()

	fmt.Fprintf(w, "funccover: %s\n", "{{.Src}}")
	
	for i, count := range funcCover_{{.Suffix}}.Count {
		fmt.Fprintf(w, "%s:%d:%d\n",funcCover_{{.Suffix}}.Name[i], funcCover_{{.Suffix}}.Line[i], count)
	}   
}`
