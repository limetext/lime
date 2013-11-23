// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package textmate

import (
	"github.com/quarnster/parser"
	"github.com/quarnster/util/text"
	"lime/backend/render"
	"lime/backend/textmate"
	"sync"
)

type (
	TextmateSyntax struct {
		textmate.LanguageParser
	}
)

func (ts *TextmateSyntax) PrepareForRendering(text.Region) []render.ViewRegions {
	ts.Lock()
	defer ts.Unlock()
	return ts.flatten([]render.ViewRegions{}, "", ts.RootNode())
}
