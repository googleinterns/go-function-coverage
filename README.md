**This is not an officially supported Google product.**

# Go Function Coverage

The project aims to collect Go function-level coverage with low overhead for any
binary.

For context, the existing coverage in Go works only for tests, and only have the
line-coverage, which can be too inefficient to run in a production environment.

## Source Code Headers

Every file containing source code must include copyright and license
information. This includes any JS/CSS files that you might be serving out to
browsers. (This is to help well-intentioned people avoid accidental copying that
doesn't comply with the license.)

Apache header:

```
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Overview

Go Function Coverage tool 'funccover' support generating instrumented source code files
so that running binary automatically collects the coverage data for functions.
    
'funccover' inserts a global function coverage variable that will keep the coverage data for functions 
to the given source code. Then it inserts necessary instructions to the top of each function 
(basic assignment instruction to global coverage variable). This way when a function starts executing, 
global coverage variable will keep the information that this function started execution some time. 

We also have to save that coverage information somewhere. Initially 'funccover' tool writes coverage information
to a file (RPC will be more useful in the future). Currently 'funccover' tool inserts 2 functions to the given
source code, one writes coverage data to a file, other calls it periodically. Period must be given as a flag to the tool.
Tool also inserts a defer call to the main to write coverage data after main function ends. So it is more general. 

## Quickstart



```bash
# Get the module from Github and install it into your $GOPATH/bin/
$ go get github.com/googleinterns/go-function-coverage/...
```
- If you add your _$GOPATH/bin_ into your _$PATH_ ([instructions](
https://github.com/golang/go/wiki/GOPATH)) you can run 'funccover' directly by writing 'funccover' to the terminal. 

## How To Use It

```bash
$ funccover [flags] [arguments...]
```

Currently 'funccover' only works for single source files, you have to give the name of that file as an argument.  

### Flags

'funccover' has 3 flags. Each flag tells 'funccover' how it should instrument the source code. 

#### -period duration

This flag represents the period of the data collection, if it is not given periodical collection will be disabled. 

```bash
$ funccover -period=500ms source.go
```

#### -dst string

This flag sets the destination file name for the instrumented source code (default "instrumented_$source.go").

```bash
$ funccover -dst=instrumented_source.go source.go
```

#### -o string

This flag sets the coverage output file name (default "cover.out").

```bash
$ funccover -period=1s -dst=instrumented.go -o=function_coverage.out source.go
```

### Example Usage

You have a source file named src.go that should get 2 arguments to run normally. You want to get the function coverage data for it to a file named cover.txt and since it is a long running code you want to get the coverage data every 1 minutes.

```bash
$ funccover -period=1m -o=cover.txt src.go
$ go build instrumented_src.go
$ ./instrumented_src argument1 argument2
```

After you build the instrumented binary, you can run the binary normally (same way you run the binary for src.go) and coverage data will be written to cover.txt in following format:

```
funccover: src.go
name1:line1:coverage1
name2:line2:coverage2
name3:line3:coverage3
...
```
Here name is the name of the function. Line is the starting line number of the corresponding function. Lastly coverage is the coverage data for that function, **1** if it is called, **0** otherwise. 

