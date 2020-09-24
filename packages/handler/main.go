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

// Package covcollect implements a package to collect the function coverage data
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/coverdata"
)

// Function collect writes the data in coverdata to "coverage.out" file
var collect func() = func() {
	fd, err := os.Create("coverage.out")
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(fd)
	defer func() {
		w.Flush()
		fd.Close()
	}()

	fmt.Fprintf(w, "cover mode: %s\n", coverdata.Cover.Mode)

	for key, element := range coverdata.Cover.Blocks {
		for i, block := range element {
			fmt.Fprintf(w, "%s:%d:%d:%d:%d:%d:%d\n", key, block.Line0, block.Col0, block.Line1,
				block.Col1, block.Stmts, coverdata.Cover.Counters[key][i])
		}
	}
}

// Initializes the collection
func init() {
	fmt.Println("-- Coverage Collection is started --")
	initSignalHandler()
	go periodicalCollect()
}

// Initializes the signal handler, collect() will be called when program is interrupted
func initSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		collect()
		os.Exit(0)
	}()
}

// Function periodicalCollect calls the collect function every 500ms
func periodicalCollect() {
	duration := 5000 * time.Millisecond
	ticker := time.NewTicker(duration)

	for range ticker.C {
		collect()
	}
}
