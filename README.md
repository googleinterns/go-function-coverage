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

The goal is to collect production Go code coverage with low overhead. Coverage scope is function granularity.
Function level run-time coverage data can be collected to:

* Detect dynamically dead code as functions never run in certain production jobs.
* Detect run-time dependencies of a test or of a production binary.

Compared to the current coverage implementations in Blaze for Go, our tools work for any binary and not for just unit tests.

Our aim is integrating our coverage to Bazel/Blaze. We are proposing a new instrumentation mode for Go and 2 new options for Bazel which will be useful for our purposes and will be useful for many other use cases as well. 

Go Function Coverage consists of two parts. 

* Instrumenting sources to keep function execution data
* Linking a library (handler) to the final binary that will write/upload the coverage data to our servers.

This repository implements an example program and a handler to demonstrate the power of our coverage.


## Quickstart

Install our version of Bazel to your workspace: https://github.com/muratekici/bazel

You can install it from source. Follow [this](https://docs.bazel.build/versions/master/install-compile-source.html) instructions. 

### With The Included Examples

* Build the example program with "func" coverage. You need to embed example handler to the binary target as well. 
```bash
$ bazel build --collect_code_coverage --code_coverage_mode=func --embed_library=//packages/handler:example-handler //packages/program:example-program
```
This will generate an executable ```bazel-bin/packages/program/example-program_/example-program```.
When you run it, it will ask you to enter numbers in range [0-9] in a line, then it will call ```F$number``` function for each number you entered.
Coverage data will be saved to ```coverage.out```.

### With Customization

* Implement your [handler] library (#handler)
   * Handler library must implement a main package without the main() function
   * --embed_library option will embed this library to the go_binary target. 
   * It's init() function will start executing.
   * Instrumentation inserts a call to LastCallForFunccoverReport function variable. You set it to run another function in init() function.
   * Please take a look at the example handler, [Handler](../master/packages/handler/main.go)
   
* To be able to use new coverage mode, you must load [modified rules_go repository](https://github.com/muratekici/rules_go) instead of official one in your WORKSPACE file. Please take look at the [WORKSPACE](../master/WORKSPACE) in this repository.

* Build your program with your handler. 

```bash
$ bazel build --collect_code_coverage --code_coverage_mode=func --embed_library=//:your_handler //:your-program
```

### Concepts

Go Function Coverage tool consists of 2 parts, instrumentation and a handler. To collect production coverage data you need to have both.

#### Instrumentation

Our instrumentation is built inside the rules_go. To be able to use it you need to install updated bazel mentioned. You can set the coverage mode using new --code_coverage_mode option. 

Other existing modes are 

* set: did each statement run?
* count: how many times did each statement run?
* atomic: like count, but counts precisely in parallel program

#### Handler

Handler is a package that implements reporting functionality. It will be [embedded](https://github.com/bazelbuild/rules_go/blob/master/go/core.rst#embedding) to the main package using new --embed_library option in build time. It's init() function will be invoked when program starts executing. You can call your go routines that will report the coverage data inside it. Coverage data will be saved inside the coverdata package during runtime. Please look at the [example](../master/packages/handler/) for more clarification. 
