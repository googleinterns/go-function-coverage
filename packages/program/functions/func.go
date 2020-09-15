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

// Package functions implements a package to collect the function coverage data
package functions

import "fmt"

// F0 print "F0" is executed
func F0() {
	fmt.Println("F0 is executed")
}

// F1 prints "F0 is executed"
func F1() {
	fmt.Println("F1 is executed")
}

// F2 prints "F2 is executed"
func F2() {
	fmt.Println("F2 is executed")
}

// F3 prints "F3 is executed"
func F3() {
	fmt.Println("F3 is executed")
}

// F4 prints "F4 is executed"
func F4() {
	fmt.Println("F4 is executed")
}

// F5 prints "F5 is executed"
func F5() {
	fmt.Println("F5 is executed")
}
