// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"strings"
)

// Naive algorithm from http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func mDiff(av, bv []string, context int) (ret []string) {
	pitch := (len(bv) + 1)
	matrix := make([]int, (len(av)+1)*pitch)
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
	// NOTE: Needs description
	innerContext := func(i, count int) {
		var index int
		switch count % 3 {
		case 0:
			break
		case 1:
			ret = append(ret, "  "+av[i])
			index, i = 1, i+1
		default:
			ret = append(ret, "  "+av[i])
			ret = append(ret, "  "+av[i+1])
			index, i = 2, i+2
		}
		for ; index != count; index += 3 {
			ret = append(ret, "  "+av[i])
			ret = append(ret, "  "+av[i+1])
			ret = append(ret, "  "+av[i+2])
			i += 3
		}
	}
	// minValue returns the minimum value of two integers.
	minValue := func(x, y int) int {
		if x > y {
			return y
		}
		return x
	}
	// addContext adds context if the context has changed.
	addContext := func(i, iLast, contextLast *int, changed *bool) {
		if !*changed {
			return
		}
		if m := minValue(*iLast, *contextLast) - *i; m > 0 {
			innerContext(*i, minValue(m, context))
		}
	}
	// NOTE: Needs description and needs to be converted into a for loop.
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
		addContext(&i, &iLast, &contextLast, &changed)
	}
	inner(len(av), len(bv), 0, len(av), len(av))
	return
}

// Diff returns the difference between two strings.
func Diff(a, b string) string {
	split := func(element string) []string {
		return strings.Split(*element, "\n")
	}
	a = strings.Replace(a, "\r\n", "\n", -1)
	b = strings.Replace(b, "\r\n", "\n", -1)
	if a != b {
		return strings.Join(mDiff(split(a), split(b), 3), "\n")
	}
	return ""
}
