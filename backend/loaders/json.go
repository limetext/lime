// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package loaders

import (
	sj "encoding/json"
	"errors"
	"fmt"
	"github.com/limetext/lime/backend/loaders/json"
	. "github.com/limetext/text"
)

func LoadJSON(data []byte, intf interface{}) error {
	var (
		b   = NewBuffer()
		p   json.JSON
		set RegionSet
	)
	defer b.Close()
	str := string(data)
	b.Insert(0, str)

	// Lime works with rune indices, but the parser works with byte indices
	// so we need to create a lookup table to map these correctly.
	lut := make([]int, len(str))
	runeidx := 0
	for i := range str {
		lut[i] = runeidx
		runeidx++
	}

	if !p.Parse(str) {
		return fmt.Errorf("%s, %s", p.Error(), p.RootNode())
	}
	root := p.RootNode()
	for _, child := range root.Children {
		switch child.Name {
		case "BlockComment", "LineComment", "EndOfFile", "JunkComma":
			if child.Range.End() < len(lut) {
				set.Add(Region{lut[child.Range.Begin()], lut[child.Range.End()]})
			}
		default:
			return errors.New("Unhandled node: " + child.Name)
		}
	}
	b.AddObserver(set)
	i := 0
	for {
		l := set.Len()
		if i >= l {
			break
		}
		r := set.Get(i)
		b.Erase(r.Begin(), r.Size())
		if l2 := set.Len(); l2 != l {
			continue
		}
		i++
	}
	// TODO(q): Map any line/column errors to the actual file's line/column
	return sj.Unmarshal([]byte(b.Substr(Region{0, b.Size()})), intf)
}

// BufferObserver

func (rs *RegionSet) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	rs.Adjust(region_removed.B, region_removed.A - region_removed.B)
}

func (rs *RegionSet) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	rs.Adjust(region_inserted.A, region_inserted.B - region_inserted.A)
}

// End of Buffer Observer