// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

// Checks if element exists in a slice
func exist(paths []string, name string) bool {
	for _, p := range paths {
		if p == name {
			return true
		}
	}
	return false
}

// Removes an element from slice
func remove(slice []string, name string) []string {
	for i, el := range slice {
		if el == name {
			slice[i], slice = slice[len(slice)-1], slice[:len(slice)-1]
			break
		}
	}
	return slice
}
