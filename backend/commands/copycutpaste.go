// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	"github.com/quarnster/util/text"
	"sort"
	"strings"
)

type (
	CopyCommand struct {
		DefaultCommand
	}

	CutCommand struct {
		DefaultCommand
	}

	PasteCommand struct {
		DefaultCommand
	}
)

func getRegions(v *View, cut bool) *text.RegionSet {
	rs := &text.RegionSet{}
	rs.AddAll(v.Sel().Regions())
	sort.Sort(rs)

	he, ae := rs.HasEmpty(), !rs.HasNonEmpty() || cut
	for _, r := range rs.Regions() {
		if ae && r.Empty() {
			rs.Add(v.Buffer().FullLineR(r))
		} else if he && r.Empty() {
			rs.Substract(r)
		}
	}

	return rs
}

func getSelSubstrs(v *View, rs *text.RegionSet) []string {
	var add, s1 string
	s := make([]string, len(rs.Regions()))
	for i, r := range rs.Regions() {
		add = ""
		s1 = v.Buffer().Substr(r)
		if !v.Sel().HasNonEmpty() && !strings.HasSuffix(s1, "\n") {
			add = "\n"
		}
		s[i] = s1 + add
	}
	return s
}

func (c *CopyCommand) Run(v *View, e *Edit) error {
	rs := getRegions(v, false)
	s := getSelSubstrs(v, rs)

	GetEditor().SetClipboard(strings.Join(s, "\n"))

	return nil
}

func (c *CutCommand) Run(v *View, e *Edit) error {
	s := getSelSubstrs(v, getRegions(v, false))

	rs := getRegions(v, true)
	sort.Sort(sort.Reverse(rs))
	for _, r := range rs.Regions() {
		v.Erase(e, r)
	}

	GetEditor().SetClipboard(strings.Join(s, "\n"))

	return nil
}

func (c *PasteCommand) Run(v *View, e *Edit) error {
	// TODO: Paste the entire line on the line before the cursor if a
	//		 line was autocopied.

	ed := GetEditor()

	rs := &text.RegionSet{}
	rs.AddAll(v.Sel().Regions())
	sort.Sort(sort.Reverse(rs))
	for _, r := range rs.Regions() {
		v.Replace(e, r, ed.GetClipboard())
	}

	return nil
}

func init() {
	register([]Command{
		&CopyCommand{},
		&CutCommand{},
		&PasteCommand{},
	})
}
