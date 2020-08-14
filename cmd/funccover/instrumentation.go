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

// This code implements a source file instrumentation function
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/tools/go/ast/astutil"
)

// Instrumenter is the interface for instrumentation
type Instrumenter interface {
	AddFile(path string) error
	Instrument() map[string][]byte
	WriteInstrumented(dir string, instrumented map[string][]byte)
}

// Keeps the data necessary for instrumentation
type packageInstrumentation struct {
	fileName    string
	suffix      string
	period      time.Duration
	fset        *token.FileSet
	parsedFiles map[string][]byte
	dir         string
}

// Adds a file to instrument to h as source code
func (h *packageInstrumentation) AddFile(src string) error {

	if h.fset == nil {
		h.fset = token.NewFileSet()
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}

	if h.parsedFiles == nil {
		h.parsedFiles = make(map[string][]byte)
	}

	h.parsedFiles[src] = content
	return nil
}

// Instrument instruments all the sources in h and returns a map
func (h *packageInstrumentation) Instrument() (map[string][]byte, error) {

	var funcCover = FuncCover{}
	var mainFile = ""
	var instrumented = make(map[string][]byte)

	for src, content := range h.parsedFiles {
		temp, flag, err := SaveFuncs(src, content)
		if err != nil {
			return nil, err
		}
		funcCover.FuncBlocks = append(funcCover.FuncBlocks, temp...)
		if flag == true {
			mainFile = src
		}
	}

	var currentIndex = 0

	for src, content := range h.parsedFiles {

		buf := new(bytes.Buffer)

		index, err := addCounters(buf, content, h.suffix, h.fileName, currentIndex)
		currentIndex = index

		if err != nil {
			return nil, err
		}

		if mainFile != src {
			instrumented[src] = buf.Bytes()
			continue
		}

		parsedFile, err := parser.ParseFile(h.fset, "", buf, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		// Ensure necessary imports are inserted
		imports := []string{"github.com/muratekici/go-function-coverage/pkg/covcollect"}
		for _, impr := range imports {
			astutil.AddImport(h.fset, parsedFile, impr)
		}

		importedContent := astToByte(h.fset, parsedFile)
		buf = new(bytes.Buffer)

		fmt.Fprintf(buf, "%s", importedContent)

		// Write necessary functions variables and an init function that calls periodical_retrieve_$hash() with w
		declCover(buf, h.suffix, h.fileName, h.period, funcCover)

		instrumented[src] = buf.Bytes()
	}

	return instrumented, nil
}

// WriteInstrumented writes the instrumented sources to dir with same names
func (h *packageInstrumentation) WriteInstrumented(instrumented map[string][]byte) error {

	path := h.dir
	os.MkdirAll(path, os.ModePerm)

	for src, content := range instrumented {
		if path != "" {
			filePath := filepath.Join(path, filepath.Base(src))
			fd, err := os.Create(filePath)
			if err != nil {
				return err
			}
			w := bufio.NewWriter(fd)
			fmt.Fprintf(w, "// %s\n%s", src, content)
			w.Flush()
			fd.Close()
		} else {
			fmt.Printf("// %s\n%s", src, content)
		}

	}
	return nil
}
