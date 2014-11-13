// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

type (
	Package interface {
		// Returns the path of the package
		Name() string

		// Depending on the implemented package
		// returns useful data for python plugin is
		// python files for setting is file content
		Get() interface{}

		// Reloads package data
		Reload()
	}
)

// This is useful when we are loading new plugin or
// scanning for user settings, snippets and etc we
// will add files which their suffix contains one of
// these keywords
var types = []string{"settings", "keymap"}
