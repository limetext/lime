// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"github.com/quarnster/parser"
	"github.com/quarnster/util/text"
	"lime/backend/render"
	"sort"
	"sync"
)

type (
	Parser interface {
		Parse() (*parser.Node, error)
	}

	SyntaxHighlighter interface {
		Adjust(position, delta int)
		ScopeExtent(point int) text.Region
		ScopeName(point int) string
		Flatten() render.ViewRegionMap
	}

	nodeHighlighter struct {
		rootNode      *parser.Node
		lastScopeNode *parser.Node
		lastScopeBuf  bytes.Buffer
		lastScopeName string
		sync.Mutex
	}
)

func NewSyntaxHighlighter(p Parser) (SyntaxHighlighter, error) {
	if rn, err := p.Parse(); err != nil {
		return nil, err
	} else {
		return &nodeHighlighter{rootNode: rn}, nil
	}
}

func (nh *nodeHighlighter) findScope(search text.Region, node *parser.Node) *parser.Node {
	idx := sort.Search(len(node.Children), func(i int) bool {
		return node.Children[i].Range.A >= search.A || node.Children[i].Range.Covers(search)
	})
	for idx < len(node.Children) {
		c := node.Children[idx]
		if c.Range.A > search.B {
			break
		}
		if c.Range.Covers(search) {
			if node.Name != "" && node != nh.lastScopeNode {
				if nh.lastScopeBuf.Len() > 0 {
					nh.lastScopeBuf.WriteByte(' ')
				}
				nh.lastScopeBuf.WriteString(node.Name)
			}
			return nh.findScope(search, node.Children[idx])
		}
		idx++
	}
	if node != nh.lastScopeNode && node.Range.Covers(search) && node.Name != "" {
		if nh.lastScopeBuf.Len() > 0 {
			nh.lastScopeBuf.WriteByte(' ')
		}
		nh.lastScopeBuf.WriteString(node.Name)
		return node
	}
	return nil
}

func (nh *nodeHighlighter) updateScope(point int) {
	if nh.rootNode == nil {
		return
	}

	search := text.Region{point, point + 1}
	if nh.lastScopeNode != nil && nh.lastScopeNode.Range.Covers(search) {
		if len(nh.lastScopeNode.Children) != 0 {
			if no := nh.findScope(search, nh.lastScopeNode); no != nh.lastScopeNode && no != nil {
				nh.lastScopeNode = no
				nh.lastScopeName = nh.lastScopeBuf.String()
			}
		}
		return
	}
	nh.lastScopeNode = nil
	nh.lastScopeBuf.Reset()
	nh.lastScopeNode = nh.findScope(search, nh.rootNode)
	nh.lastScopeName = nh.lastScopeBuf.String()
}

func (nh *nodeHighlighter) ScopeExtent(point int) text.Region {
	nh.updateScope(point)
	if nh.lastScopeNode != nil {
		r := nh.lastScopeNode.Range
		return text.Region{r.A, r.B}
	}
	return text.Region{}
}

func (nh *nodeHighlighter) ScopeName(point int) string {
	nh.updateScope(point)
	return nh.lastScopeName
}

func (nh *nodeHighlighter) flatten(vrmap render.ViewRegionMap, scopename string, node *parser.Node) {
	scopename += " " + node.Name
	cur := node.Range

	for _, c := range node.Children {
		if cur.A <= c.Range.A {
			reg := vrmap[scopename]
			reg.Flags |= render.DRAW_TEXT
			reg.Scope = scopename
			reg.Regions.Add(text.Region{cur.A, c.Range.A})
			vrmap[scopename] = reg
		}
		cur.A = c.Range.B
		nh.flatten(vrmap, scopename, c)
	}
	// Just add the last region if it's not zero sized
	if cur.A != cur.B {
		reg := vrmap[scopename]
		reg.Flags |= render.DRAW_TEXT
		reg.Scope = scopename
		reg.Regions.Add(text.Region{cur.A, cur.B})
		vrmap[scopename] = reg
	}
}

func (nh *nodeHighlighter) Adjust(position, delta int) {
	nh.rootNode.Adjust(position, delta)
}

func (nh *nodeHighlighter) Flatten() (ret render.ViewRegionMap) {
	ret = make(render.ViewRegionMap)
	nh.flatten(ret, "lime.syntax", nh.rootNode)
	return
}
