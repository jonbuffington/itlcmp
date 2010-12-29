//
// Copyright 2010 by Jon Buffington. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
)


type statistics struct {
	tracks int
	exts   map[string]int
	files  int
	kinds  map[string]int
}

func newStatistics() *statistics {
	return &statistics{
		exts:   make(map[string]int),
		files:  0,
		kinds:  make(map[string]int),
		tracks: 0,
	}
}

func (s statistics) print() {
	printMap := func(m map[string]int, singular string, plural string) {
		for k, v := range m {
			var format string
			if v == 1 {
				format = singular
			} else {
				format = plural
			}
			fmt.Printf("\t"+format, v, k)
		}
	}

	fmt.Printf("\nFound %d total tracks:\n", s.tracks)
	printMap(s.kinds, "%d %s.\n", "%d %ss.\n")
	fmt.Printf("\nFound %d total files:\n", s.files)
	printMap(s.exts, "%d %s file.\n", "%d %s files.\n")
}
