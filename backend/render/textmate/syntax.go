package textmate

import (
	"github.com/quarnster/parser"
	"lime/backend/primitives"
	"lime/backend/render"
	"lime/backend/textmate"
)

type (
	TextmateSyntax struct {
		textmate.LanguageParser
	}
)

func (ts *TextmateSyntax) flatten(in []render.ViewRegions, scopename string, node *parser.Node) []render.ViewRegions {
	scopename += " " + node.Name
	cur := node.Range

	for _, c := range node.Children {
		if cur.Start < c.Range.Start {
			var vr render.ViewRegions
			vr.Regions.Add(primitives.Region{cur.Start, c.Range.Start})
			vr.Scope = scopename
			in = append(in, vr)
			cur.Start = c.Range.End
		}
		in = ts.flatten(in, scopename, c)
	}
	if cur.Start != cur.End {
		var vr render.ViewRegions
		vr.Regions.Add(primitives.Region{cur.Start, cur.End})
		vr.Scope = scopename
		in = append(in, vr)
	}
	return in
}

func (ts *TextmateSyntax) Process(buf primitives.Buffer) []render.ViewRegions {
	ts.Parse(buf.Substr(primitives.Region{0, buf.Size()}))
	return ts.flatten([]render.ViewRegions{}, "", ts.RootNode())
}
