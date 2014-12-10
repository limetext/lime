// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// The parser package defines interfaces responsible for creating
// an Abstract Syntax Tree like structure of a text document.
//
// It should then be possible to query this structure for
// the name and extend of the various code scopes defined within it.
//
// TODO:
// It should be possible to hook in for example libclang,
// go/ast and other "proper" code parsers. Do these interfaces
// make sense for those or should they be changed?

package parser

import (
	"bytes"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/text"
	"github.com/quarnster/parser"
	"sort"
	"sync"
)

type (
	// The Parser interface is responsible for creating
	// a parser.Node structure of a given text data.
	Parser interface {
		Parse() (*parser.Node, error)
	}

	// The SyntaxHighlighter interface is responsible for
	// identifying the extent and name of code scopes given
	// a position in the code buffer this specific SyntaxHighlighter
	// is responsible for.
	//
	// It's expected that the syntax highlighter monkey patches its existing
	// scope data rather than performing a full reparse when the underlying
	// buffer changes.
	//
	// This is because a full reparse, for which the Parser interface is responsible,
	// will be going on in parallel in a separate thread and the "monkey patch"
	// will allow some accuracy in the meantime until the Parse operation has finished.
	SyntaxHighlighter interface {
		// Adjust is called when the underlying text buffer changes at "position"
		// with a change of "delta" characters either being inserted or removed.
		//
		// See note above regarding "monkey patching".
		Adjust(position, delta int)

		// Returns the Region of the inner most Scope extent which contains "point".
		//
		// This method can be called a lot by plugins, and should therefore be as
		// fast as possible.
		ScopeExtent(point int) text.Region

		// Returns the full concatenated nested scope name of the scope(s) containing "point".
		//
		// This method can be called a lot by plugins, and should therefore be as
		// fast as possible.
		ScopeName(point int) string

		// Flatten creates a map where the key is the concatenated nested scope names
		// and the key is the render.ViewRegions associated with that key.
		//
		// This function is only called once by the View, which merges
		// the regions into its own region map and adjusts them as appropriate.
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

// Creates a new default implementation of SyntaxHighlighter operating
// on the AST created by  "p"'s Parse().
func NewSyntaxHighlighter(p Parser) (SyntaxHighlighter, error) {
	if rn, err := p.Parse(); err != nil {
		return nil, err
	} else {
		return &nodeHighlighter{rootNode: rn}, nil
	}
}

// Given a text region, returns the innermost node covering that region.
// Side-effects: Writes to nh.lastScopeBuf...
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

// Caches the full concatenated nested scope name and the innermost node that covers "point".
// TODO: multiple cursors being in different scopes is harsh on the cache...
func (nh *nodeHighlighter) updateScope(point int) {
	if nh.rootNode == nil {
		return
	}

	search := text.Region{A: point, B: point + 1}
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
		return text.Region{A: r.A, B: r.B}
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
			reg.Regions.Add(text.Region{A: cur.A, B: c.Range.A})
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
		reg.Regions.Add(text.Region{A: cur.A, B: cur.B})
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
