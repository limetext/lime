// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"strings"
)

// Naive algorithm from http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func mDiff(av, bv []string) (ret []string) {
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
	var inner func(i, j int, context int)
	inner = func(i, j int, context int) {
		if i > 0 && j > 0 && av[i-1] == bv[j-1] {
			i--
			j--
			inner(i, j, context-1)
			if context > 0 {
				ret = append(ret, "  "+bv[j])
			}
		} else if j > 0 && (i == 0 || matrix[i*pitch+j-1] >= matrix[(i-1)*pitch+j]) {
			inner(i, j-1, 3)
			ret = append(ret, "+ "+bv[j-1])
		} else if i > 0 && (j == 0 || matrix[i*pitch+j-1] < matrix[(i-1)*pitch+j]) {
			inner(i-1, j, 3)
			ret = append(ret, "- "+av[i-1])
		}
	}
	inner(len(av), len(bv), 0)
	return
}

func Diff(a, b string) string {
	if a == b {
		return ""
	}
	a = strings.Replace(a, "\r", "", -1)
	b = strings.Replace(b, "\r", "", -1)
	if a == b {
		return ""
	}
	return strings.Join(mDiff(strings.Split(a, "\n"), strings.Split(b, "\n")), "\n")
}
