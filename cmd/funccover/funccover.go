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
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/muratekici/go-function-coverage/pkg/instrument"
)

const usageMessage = "" + `usage: funccover [instrumentation flags] [arguments...]:
funccover generates an instrumented source code for function coverage
Generated source code can be built or ran normally to get the coverage data
Coverage data will be written to a file periodically while binary is running also when main ends

Currently funccover only works for single source file, source file path shall be given as an argument
`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	funcCoverPeriod    = flag.Duration("period", 0, "period of the data collection ex.500ms\nif not given no periodical collection")
	instrumentedSource = flag.String("dst", "", "file for instrumented source\nif not given instrumented_$source.go")
	outputFile         = flag.String("o", "cover.out", "file for coverage output")
)

var functionCounter int
var uniqueHash string

func main() {

	flag.Usage = usage
	flag.Parse()

	// Usage information when no arguments.
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
	}

	err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, `For usage information, run 'funccover' with no arguments and flags`)
		os.Exit(2)
	}

	var src = flag.Args()[0]

	if *instrumentedSource == "" {
		*instrumentedSource = "instrumented_" + src
	}

	fd, err := os.Create(*instrumentedSource)
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(fd)

	defer func() {
		w.Flush()
		fd.Close()
	}()

	fmt.Fprintf(w, "//line %s:1\n", src)

	// Declare the unique suffix for functions and variables
	sum := sha256.Sum256([]byte(src))
	uniqueHash = fmt.Sprintf("%x", sum[:6])

	instrument.Annotate(w, content, src, uniqueHash, *outputFile, *funcCoverPeriod)
}

// parseFlags performs validations.
func parseFlags() error {

	if *funcCoverPeriod < 0 {
		return fmt.Errorf("-period: %s is not a valid period", *funcCoverPeriod)
	}

	if *outputFile == *instrumentedSource {
		return fmt.Errorf("coverage output file and instrumented source file can not have the same name: %s", *instrumentedSource)
	}

	if *outputFile == "" {
		return fmt.Errorf("output file name can not be empty")
	}

	if flag.NArg() == 0 {
		return fmt.Errorf("missing source file")
	}

	if flag.NArg() == 1 {
		return nil
	}

	return fmt.Errorf("too many arguments")
}
