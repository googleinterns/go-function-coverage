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
	"go/ast"
	"go/parser"
	"go/token"
)

// FuncCoverBlock contains tha name and line of a function
// Line contains the first line of the definition in the source code
type FuncCoverBlock struct {
	Name string // Name of the function
	Line uint32 // Line number for block start.
}

// FuncCover contains the FuncCoverBlock informations of a source code
type FuncCover struct {
	FuncBlocks []FuncCoverBlock
}

// SaveFuncs parses given source code and returns a FuncCover instance, also returns true if main function is given
func SaveFuncs(src string, content []byte) ([]FuncCoverBlock, bool, error) {

	fset := token.NewFileSet()

	parsedFile, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, false, err
	}

	var funcBlocks []FuncCoverBlock
	flag := false

	// Find function declerations to instrument and save them to funcCover
	for _, decl := range parsedFile.Decls {
		switch t := decl.(type) {
		// Function Decleration
		case *ast.FuncDecl:
			if t.Name.Name == "main" {
				flag = true
			}
			funcBlocks = append(funcBlocks, FuncCoverBlock{
				Name: src + ":" + t.Name.Name,
				Line: uint32(fset.Position(t.Pos()).Line),
			})
		}
	}

	return funcBlocks, flag, nil
}
