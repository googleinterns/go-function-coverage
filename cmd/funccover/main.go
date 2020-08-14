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
	"crypto/sha256"
	"flag"
	"fmt"
	_ "io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"strings"
)

const usageMessage = "" + `usage: go compile -preprocess=[funccover [funccover args]] [args]:
Funccover generates an instrumented source code for function coverage that imports covcollect package
Generated source code can be built or ran normally to get the coverage data
Coverage data will be written to a file periodically while binary is running also when main ends
`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	funcCoverPeriod = flag.Duration("period", 0, "period of the data collection ex.500ms\nif not given no periodical collection")
	direc           = flag.String("dir", "", "directory for instrumented package\nif not given stdout")
	outputFile      = flag.String("o", "cover.out", "file for coverage output")
)

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

	args := flag.Args()

	err = instrumentPackage(getGoSources(args))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, `For ;usage information, run 'funccover' with no arguments and flags`)
		os.Exit(2)
	}
}

// parseFlags performs validations.
func parseFlags() error {

	if *funcCoverPeriod < 0 {
		return fmt.Errorf("-period: %s is not a valid period", *funcCoverPeriod)
	}

	if flag.NArg() == 0 {
		return fmt.Errorf("missing source files")
	}

	return nil
}

func getGoSources(args []string) []string {
	var sources []string
	for _, arg := range args {
		if strings.HasSuffix(arg, ".go") {
			sources = append(sources, arg)
		}
	}
	return sources
}

func instrumentPackage(files []string) error {

	if len(files) == 0 {
		return nil
	}

	sum := sha256.Sum256([]byte(files[0]))
	uniqueHash := fmt.Sprintf("%x", sum[:6])

	var instrumentation = packageInstrumentation{
		suffix:   uniqueHash,
		dir:      *direc,
		period:   *funcCoverPeriod,
		fileName: *outputFile,
	}

	for _, src := range files {
		res, _ := filepath.Abs(src)
		err := instrumentation.AddFile(res)
		if err != nil {
			return err
		}
	}

	instrumented, err := instrumentation.Instrument()
	if err != nil {
		return err
	}

	err = instrumentation.WriteInstrumented(instrumented)

	if err != nil {
		return err
	}

	return nil
}
