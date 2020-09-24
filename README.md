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

Compared to the current coverage implementations in Bazel for Go, our tools work for any binary and not for just unit tests.

Our aim is to integrate our coverage tools to Blaze as well. We are proposing a new instrumentation mode for Go with a new coverage mode build configuration for Go rules. Also a new option ```--embed_library``` for Bazel, which will be useful for our purposes and will be useful for other use cases as well. 
 

Go Function Coverage consists of two parts. 

* Instrumenting sources to keep function execution data.
* Linking a collection library to the final binary that will write/upload the coverage data.

This repository implements an example program and two different [handler libraries](#handler) demonstrate the power of our coverage tools.

## Quickstart

Install our version of [Bazel](https://github.com/muratekici/bazel).

You can install it from source. Follow [this](https://docs.bazel.build/versions/master/install-compile-source.html) instructions. 

### With The Included Examples

Only thing you need to do is build the example program with "func" coverage mode. You need to embed example handler for "func" mode to the binary target as well. 
```bash
$ bazel build --collect_code_coverage --embed_library=//packages/func-handler:example-handler //packages/program:example-program --@io_bazel_rules_go//go/config:covermode=func
```
This will generate an executable ```bazel-bin/packages/program/example-program_/example-program```.
When you run the binary, it will ask you to enter numbers in range [0-9]. For each number you enter program will execute ```F$number``` function.

You can also use already existing coverage modes with covermode option. This repository implements a handler for already existing modes as well which you can find [here](https://blog.golang.org/cover). 

```bash
$ bazel build --collect_code_coverage --embed_library=//packages/handler:example-handler //packages/program:example-program --@io_bazel_rules_go//go/config:covermode=[count | set | atomic] 
```

Coverage data will be saved to ```coverage.out``` for both examples.

### With Customization

* Implement your [handler library](#handler)
   * Handler library must implement a main package without the main() function.
   * ```--embed_library``` option will embed this library to the ```go_binary``` target. 
   * It's ```init()``` function will start executing.
   * Instrumentation mode ```"func"``` defines an exit hook, ```FuncCoverExitHook``` which points to an empty function with no paramteres. You can set it to your custom function and it will be called automatically just before program exits.
   * Please take a look at the example handlers for [func](../master/packages/func-handler/main.go) and for [set, count and atomic](../master/packages/handler/main.go).
   
* To be able to use new coverage mode and new options, you must load [modified rules_go repository](https://github.com/muratekici/rules_go) instead of official one in your WORKSPACE file. Please take look at the [WORKSPACE](../master/WORKSPACE) in this repository.

* Build your program with ```"func"``` coverage and your handler. 

```bash
$ bazel build --collect_code_coverage --embed_library=//:your_handler //:your-program --@io_bazel_rules_go//go/config:covermode=func
```

### Concepts

Go Function Coverage tools consists of 2 parts, instrumentation and a handler. To collect production coverage data you need to have both.

#### Instrumentation

Our instrumentation is built inside the rules_go. To be able to use it you need to use [modified rules_go](https://github.com/muratekici/rules_go). You can set the coverage mode using new ```--@io_bazel_rules_go//go/config:covermode=mode``` build configuration. 

Other existing modes are 

* set: did each statement run?
* count: how many times did each statement run?
* atomic: like count, but for parellel programs.

#### Handler

Handler is a package that implements reporting functionality. It will be [embedded](https://github.com/bazelbuild/rules_go/blob/master/go/core.rst#embedding) to the main package using new ```--embed_library``` option in build time. To be able to use our new ```--embed_library``` option, you need to install our modified version of [bazel](#Quickstart).  

Its ```init()``` function will be invoked when program starts executing. You can call your go routines that will report the coverage data inside it. Coverage data will be saved inside the ```coverdata``` package during runtime. Please look at the [examples](../master/packages/) for more clarification. 
