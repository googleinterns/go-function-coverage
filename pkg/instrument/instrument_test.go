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
	"go/parser"
	"go/token"
	"testing"
	"time"
)

func TestAstToByte(t *testing.T) {
	for i, content := range sampleCodes {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, "", content.src, parser.ParseComments)
		if err != nil {
			t.Errorf("could not parse the source")
		}
		temp := astToByte(fset, parsedFile)
		if bytes.Equal(temp, content.src) == false {
			t.Errorf("source code sample does not match the astToByte output %d", i)
		}
	}
}

func TestDeclCover(t *testing.T) {

	for i := range sampleDecl {
		buf := new(bytes.Buffer)
		period, err := time.ParseDuration(sampleDecl[i].period)
		if err != nil {
			t.Errorf("could not parse the duration")
		}
		declCover(buf, "main.go", sampleDecl[i].suffix, sampleDecl[i].out, period, FuncCover{sampleDecl[i].funcBlocks})
		fmt.Print("++++")
		fmt.Print(buf)
		fmt.Print("++++")
		if string(sampleDecl[i].res) != buf.String() {
			t.Errorf("decleration statements do not match %d", i)
		}
	}
}

func TestAddCounters(t *testing.T) {

	for i := range sampleCodes {
		buf := new(bytes.Buffer)
		addCounters(buf, sampleCodes[i].src, sampleCodes[i].suffix)
		if string(sampleCodes[i].res) != buf.String() {
			t.Errorf("counter statements do not match %d", i)
		}
	}
}

var sampleDecl = []struct {
	suffix     string
	period     string
	out        string
	funcBlocks []FuncCoverBlock
	res        []byte
}{
	{"100001", "-500ms", "output.out", []FuncCoverBlock{
		FuncCoverBlock{
			Name: "main",
			Line: 6,
		},
	}, []byte(`
var funcCover_100001 = struct {
	Count     [1]uint32
	Line      [1]uint32
	Name      [1]string
} {
	Line: [1]uint32{ 
		6,
	},
	Name: [1]string{ 
		"main",
	},
}

func retrieve_coverage_data_100001() {

	fd, err := os.Create("output.out")
	if err != nil {
    	panic(err)
	}
	
	w := bufio.NewWriter(fd)

	defer func() { 
    	w.Flush()
    	fd.Close()
	}()

	fmt.Fprintf(w, "funccover: %s\n", "main.go")
	
	for i, count := range funcCover_100001.Count {
		fmt.Fprintf(w, "%s:%d:%d\n",funcCover_100001.Name[i], funcCover_100001.Line[i], count)
	}   
}`)},
	{"122340", "2s", "nope.txt", []FuncCoverBlock{
		FuncCoverBlock{
			Name: "main",
			Line: 6,
		},
	}, []byte(`
var funcCover_122340 = struct {
	Count     [1]uint32
	Line      [1]uint32
	Name      [1]string
} {
	Line: [1]uint32{ 
		6,
	},
	Name: [1]string{ 
		"main",
	},
}

func init() {
    go periodical_retrieve_122340()
}

func periodical_retrieve_122340() {

	period, err := time.ParseDuration("2s")

	if err != nil {
    	panic(err)
	}

    ticker := time.NewTicker(period)

    for _ = range ticker.C {
        retrieve_coverage_data_122340()
    } 

}

func retrieve_coverage_data_122340() {

	fd, err := os.Create("nope.txt")
	if err != nil {
    	panic(err)
	}
	
	w := bufio.NewWriter(fd)

	defer func() { 
    	w.Flush()
    	fd.Close()
	}()

	fmt.Fprintf(w, "funccover: %s\n", "main.go")
	
	for i, count := range funcCover_122340.Count {
		fmt.Fprintf(w, "%s:%d:%d\n",funcCover_122340.Name[i], funcCover_122340.Line[i], count)
	}   
}`)},
	{"", "0ns", "hello.txt", []FuncCoverBlock{}, []byte(`
var funcCover_ = struct {
	Count     [0]uint32
	Line      [0]uint32
	Name      [0]string
} {
	Line: [0]uint32{ 
	},
	Name: [0]string{ 
	},
}

func retrieve_coverage_data_() {

	fd, err := os.Create("hello.txt")
	if err != nil {
    	panic(err)
	}
	
	w := bufio.NewWriter(fd)

	defer func() { 
    	w.Flush()
    	fd.Close()
	}()

	fmt.Fprintf(w, "funccover: %s\n", "main.go")
	
	for i, count := range funcCover_.Count {
		fmt.Fprintf(w, "%s:%d:%d\n",funcCover_.Name[i], funcCover_.Line[i], count)
	}   
}`)},
}

var sampleCodes = []struct {
	src    []byte
	suffix string
	res    []byte
}{
	{[]byte(`package main

func f1()	{}

func main() {
	f1()
}
`), "999829", []byte(`package main

func f1()	{
	funcCover_999829.Count[0] = 1;}

func main() {
	funcCover_999829.Count[1] = 1;
	f1()

	defer retrieve_coverage_data_999829()
}
`)},
	{[]byte(`package main

func main() {
}
`), "677261", []byte(`package main

func main() {
	funcCover_677261.Count[0] = 1;

	defer retrieve_coverage_data_677261()
}
`)},
	{[]byte(`package main

func f() {

}
`), "101010", []byte(`package main

func f() {
	funcCover_101010.Count[0] = 1;

}
`)},
}
