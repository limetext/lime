package textmate

import (
	"github.com/quarnster/parser"
	"lime/backend/primitives"
	"lime/backend/render"
	"lime/backend/textmate"
	"sync"
)

type (
	TextmateSyntax struct {
		textmate.LanguageParser
	}
)

func (ts *TextmateSyntax) PrepareForRendering(primitives.Region) []render.ViewRegions {
	ts.Lock()
	defer ts.Unlock()
	return ts.flatten([]render.ViewRegions{}, "", ts.RootNode())
}
