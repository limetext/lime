// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"strings"
)

// Naive algorithm from http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func mDiff(av, bv []string, context int) (ret []string) {
	matrix := make([]int, (len(av)+1)*(len(bv)+1))
	pitch := (len(bv) + 1)
	for i, a := range av {
		mp := (i+1)*pitch + 1

		for _, b := range bv {
			if a == b {
				matrix[mp] = matrix[mp-1-pitch] + 1
			} else if matrix[mp-1] > matrix[mp-pitch] {
				matrix[mp] = matrix[mp-1]
			} else {
				matrix[mp] = matrix[mp-pitch]
			}
			mp++
		}
	}

	innerContext := func(i, count int) {
		if count%2 == 0 {
			for index := 0; index != count; index += 2 {
				ret = append(ret, "  "+av[i])
				ret = append(ret, "  "+av[i+1])
				i += 2
			}
		} else {
			for index := 0; index < count; index++ {
				ret = append(ret, "  "+av[i])
				i++
			}
		}
	}

	var inner func(i, j, k, iLast, contextLast int)
	inner = func(i, j, k, iLast, contextLast int) {
		changed := false
		if i > 0 && j > 0 && av[i-1] == bv[j-1] {
			c := contextLast
			if k > 0 {
				c = i - 1
			}

			inner(i-1, j-1, k-1, iLast, c)

			// add context before the change
			if k > 0 {
				ret = append(ret, "  "+av[i-1])
			}
		} else if j > 0 && (i == 0 || matrix[i*pitch+j-1] >= matrix[(i-1)*pitch+j]) {
			changed = true
			inner(i, j-1, context, i, contextLast)
			ret = append(ret, "+ "+bv[j-1])
		} else if i > 0 && (j == 0 || matrix[i*pitch+j-1] < matrix[(i-1)*pitch+j]) {
			changed = true
			inner(i-1, j, context, i-1, contextLast)
			ret = append(ret, "- "+av[i-1])
		}

		if changed {
			// add context after the change
			l := iLast
			if l > contextLast {
				l = contextLast
			}

			m := l - i
			if m > 0 {
				if m > context {
					m = context
				}
				innerContext(i, m)
			}
		}
	}

	inner(len(av), len(bv), 0, len(av), len(av))
	return
}

// Diff returns the difference between two strings.
func Diff(a, b string) string {
	split := func(element *string) []string {
		return strings.Split(*element, "\n")
	}
	a = strings.Replace(a, "\r\n", "\n", -1)
	b = strings.Replace(b, "\r\n", "\n", -1)
	if a != b {
		return strings.Join(mDiff(split(&a), split(&b), 3), "\n")
	} else {
		return ""
	}
}
