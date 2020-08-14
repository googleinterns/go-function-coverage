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
package covcollect

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Cover is coverage struct that keeps the coverage data during runtime
type Cover struct {
	Names  []string
	Lines  []uint32
	Counts []bool
}

// InsertFuncs inserts functions of a source code to c
func (c Cover) InsertFuncs(names []string, lines []uint32, counts []bool) {
	c.Names = append(c.Names, names...)
	c.Lines = append(c.Lines, lines...)
	c.Counts = append(c.Counts, counts...)
}

// PeriodicalCollect periodically calls the Collect function with args
func (c Cover) PeriodicalCollect(period string, args ...string) {

	duration, err := time.ParseDuration(period)
	if err != nil || duration <= 0 {
		return
	}

	ticker := time.NewTicker(duration)

	for range ticker.C {
		c.Collect(args...)
	}
}

// Collect writes the data in c using args given
func (c Cover) Collect(args ...string) {

	fd, err := os.Create(args[0])
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(fd)

	defer func() {
		w.Flush()
		fd.Close()
	}()

	for i, count := range c.Counts {
		fmt.Fprintf(w, "%s:%d:%t\n", c.Names[i], c.Lines[i], count)
	}
}
