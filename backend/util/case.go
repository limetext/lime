// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"regexp"
	"strings"
)

func PascalCaseToSnakeCase(in string) string {
	casere := regexp.MustCompile(`([A-Z])`)

	first := true
	return casere.ReplaceAllStringFunc(in, func(in string) string {
		if first {
			first = false
			return strings.ToLower(in)
		}
		return "_" + strings.ToLower(in)
	})

}
