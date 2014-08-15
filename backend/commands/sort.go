// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"sort"
	"strings"
)

type (
	// The SortLinesCommand sorts all lines
	// intersecting a selection region
	SortLinesCommand struct {
		DefaultCommand
		CaseSensitive    bool
		Reverse          bool
		RemoveDuplicates bool
	}

	// The SortSelectionCommand sorts contents
	// of each selection region with respect to
	// each other
	SortSelectionCommand struct {
		DefaultCommand
		CaseSensitive    bool
		Reverse          bool
		RemoveDuplicates bool
	}

	// Helper type to sort Regions by theirs positions
	regionSorter []Region

	// Helper struct to sort strings
	textSorter struct {
		texts         []string
		caseSensitive bool
		reverse       bool
	}
)

// regionSorter implements sort.Interface
func (s regionSorter) Len() int { return len(s) }
func (s regionSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s regionSorter) Less(i, j int) bool {
	return s[i].Begin() < s[j].Begin()
}

// stringSorter implements sort.Interface
func (s textSorter) Len() int { return len(s.texts) }
func (s textSorter) Swap(i, j int) {
	s.texts[i], s.texts[j] = s.texts[j], s.texts[i]
}
func (s textSorter) Less(i, j int) bool {
	textA := s.texts[i]
	textB := s.texts[j]

	if !s.caseSensitive {
		textA = strings.ToLower(textA)
		textB = strings.ToLower(textB)
	}

	if s.reverse {
		return textA > textB
	}
	return textA < textB
}

func (c *SortLinesCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	buf := v.Buffer()

	// Used as a set of int
	sorted_rows := make(map[int]bool)

	regions := []Region{}
	texts := []string{}
	for i := 0; i < sel.Len(); i++ {
		// Get regions containing each line.
		for _, r := range buf.Lines(sel.Get(i)) {
			if ok := sorted_rows[r.Begin()]; !ok {
				sorted_rows[r.Begin()] = true
				regions = append(regions, r)
				texts = append(texts, buf.Substr(r))
			}
		}
	}

	sort.Sort(textSorter{
		texts:         texts,
		caseSensitive: c.CaseSensitive,
		reverse:       c.Reverse,
	})
	if c.RemoveDuplicates {
		texts = removeDuplicates(c.CaseSensitive, texts)
	}

	sort.Sort(regionSorter(regions))

	offset := 0
	for i, r := range regions {
		r = Region{r.A + offset, r.B + offset}
		if i < len(texts) {
			v.Replace(e, r, texts[i])
			offset += len(texts[i]) - r.Size()
		} else {
			// Erase the line and its ending
			fullLine := buf.FullLineR(r)
			v.Erase(e, fullLine)
			offset -= fullLine.Size()
		}
	}

	return nil
}

func (c *SortSelectionCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	buf := v.Buffer()

	regions := make([]Region, sel.Len())
	texts := make([]string, sel.Len())
	for i := 0; i < sel.Len(); i++ {
		regions[i] = sel.Get(i)
		texts[i] = buf.Substr(regions[i])
	}

	sort.Sort(textSorter{
		texts:         texts,
		caseSensitive: c.CaseSensitive,
		reverse:       c.Reverse,
	})
	if c.RemoveDuplicates {
		texts = removeDuplicates(c.CaseSensitive, texts)
	}

	sort.Sort(regionSorter(regions))

	offset := 0
	for i, r := range regions {
		r = Region{r.A + offset, r.B + offset}
		if i < len(texts) {
			v.Replace(e, r, texts[i])
			offset += len(texts[i]) - r.Size()
		} else {
			v.Erase(e, r)
			offset -= r.Size()
		}
	}

	return nil
}

// Remove duplicate ones from a sorted slice of string
func removeDuplicates(caseSensitive bool, xs []string) []string {
	var i, j int
	for j < len(xs) {
		var accept bool
		if i > 0 {
			prev := xs[i-1]
			curr := xs[j]
			if !caseSensitive {
				prev = strings.ToLower(prev)
				curr = strings.ToLower(curr)
			}
			accept = (prev != curr)
		} else {
			accept = true
		}
		if accept {
			xs[i] = xs[j]
			i++
			j++
		} else {
			j++
		}
	}
	return xs[:i]
}

func init() {
	register([]cmd{
		{"sort_lines", &SortLinesCommand{}},
		{"sort_selection", &SortSelectionCommand{}},
	})
}
