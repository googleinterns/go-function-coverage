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
	"fmt"
	"functions"
)

func main() {
	help()
	start()
}

func help() {
	fmt.Println("This is an example program to test new coverage functionality for Bazel")
	fmt.Println("Please print an integer in range [0-9].")
	fmt.Println("Program will call f$number function in the functions package.")
	fmt.Println("You can enter as many numbers you want.")
	fmt.Println("Example handler will save the data to coverage.out file every 500ms and when the program exits.")
	fmt.Println("Enter -1 to exit the program.")
}

func start() {
	for {
		fmt.Print("Enter a number:")
		var num int
		fmt.Scanf("%d", &num)
		switch num {
		case 0:
			functions.F0()
		case 1:
			functions.F1()
		case 2:
			functions.F2()
		case 3:
			functions.F3()
		case 4:
			functions.F4()
		case 5:
			functions.F5()
		case 6:
			functions.F6()
		case 7:
			functions.F7()
		case 8:
			functions.F8()
		case 9:
			functions.F9()
		case -1:
			return
		default:
		}
	}
}
