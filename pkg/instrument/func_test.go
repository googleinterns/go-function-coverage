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
	"testing"
)

func TestSaveFunc(t *testing.T) {
	for i, content := range sampleCodesFunc {
		temp := SaveFuncs(content)
		for j := range temp.FuncBlocks {
			if len(temp.FuncBlocks) != len(sampleWanted[i].FuncBlocks) {
				t.Errorf("number of functions do not match")
			}
			if temp.FuncBlocks[j].Name != sampleWanted[i].FuncBlocks[j].Name {
				t.Errorf("function names do not match")
			}
			if temp.FuncBlocks[j].Line != sampleWanted[i].FuncBlocks[j].Line {
				t.Errorf("function lines do not match")
			}
		}
	}
}

func TestSaveFuncBroken(t *testing.T) {
	for _, content := range sampleBrokenCodes {
		if len(SaveFuncs(content).FuncBlocks) != 0 {
			t.Errorf("there should be no functions found")
		}
	}
}

var sampleCodesFunc = [][]byte{[]byte(`package main

func f1() {}
func f2() {}

func main() {
	f1()
	f2()
}`), []byte(`package main

func main() {
}`), []byte(`package main

func init() {
}

func main() {
}`)}

var sampleWanted = []FuncCover{
	FuncCover{
		FuncBlocks: []FuncCoverBlock{
			FuncCoverBlock{
				Name: "f1",
				Line: 3,
			},
			FuncCoverBlock{
				Name: "f2",
				Line: 4,
			},
			FuncCoverBlock{
				Name: "main",
				Line: 6,
			},
		},
	},
	FuncCover{
		FuncBlocks: []FuncCoverBlock{
			FuncCoverBlock{
				Name: "main",
				Line: 3,
			},
		},
	},
	FuncCover{
		FuncBlocks: []FuncCoverBlock{
			FuncCoverBlock{
				Name: "init",
				Line: 3,
			},
			FuncCoverBlock{
				Name: "main",
				Line: 6,
			},
		},
	},
}

var sampleBrokenCodes = [][]byte{[]byte(`package main
no func this is a bad line () {}
func main() {
}`), []byte(`hello world`)}
