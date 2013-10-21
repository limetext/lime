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
