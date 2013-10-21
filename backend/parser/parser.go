package parser

import (
	"bytes"
	"github.com/quarnster/parser"
	"github.com/quarnster/util/text"
	"sort"
	"sync"
)

type (
	Parser interface {
		Parse() (*parser.Node, error)
	}

	NamedRegion struct {
		text.Region
		Name string
	}

	SyntaxHighlighter interface {
		ScopeExtent(point int) text.Region
		ScopeName(point int) string
		Flatten(viewport text.Region) []NamedRegion
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
	} else {
		nh.lastScopeNode = nil
		nh.lastScopeBuf.Reset()
		nh.lastScopeNode = nh.findScope(search, nh.rootNode)
		nh.lastScopeName = nh.lastScopeBuf.String()
	}
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

func (nh *nodeHighlighter) flatten(in []NamedRegion, scopename string, node *parser.Node) []NamedRegion {
	scopename += " " + node.Name
	cur := node.Range

	for _, c := range node.Children {
		if cur.A < c.Range.A {
			var vr NamedRegion
			vr.A, vr.B = cur.A, c.Range.A
			vr.Name = scopename
			in = append(in, vr)
			cur.A = c.Range.B
		}
		in = nh.flatten(in, scopename, c)
	}
	if cur.A != cur.B {
		var vr NamedRegion
		vr.A, vr.B = cur.A, cur.B
		vr.Name = scopename
		in = append(in, vr)
	}
	return in
}

func (nh *nodeHighlighter) Flatten(viewport text.Region) []NamedRegion {
	return nh.flatten(nil, "", nh.rootNode)
}
