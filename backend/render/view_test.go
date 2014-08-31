// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package render

import (
	"github.com/quarnster/util/text"
	"reflect"
	"testing"
)

func TestViewRegionsCull(t *testing.T) {
	tests := []struct {
		regions []text.Region
		cull    text.Region
		exp     []text.Region
	}{
		{
			[]text.Region{{100, 200}},
			text.Region{0, 50},
			[]text.Region{},
		},
		{
			[]text.Region{{100, 100}},
			text.Region{100, 100},
			[]text.Region{},
		},
		{
			[]text.Region{{100, 100}},
			text.Region{95, 105},
			[]text.Region{},
		},
		{
			[]text.Region{{100, 100}},
			text.Region{95, 100},
			[]text.Region{},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{150, 150},
			[]text.Region{},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{90, 100},
			[]text.Region{},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{100, 150},
			[]text.Region{{100, 150}},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{150, 175},
			[]text.Region{{150, 175}},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{0, 150},
			[]text.Region{{100, 150}},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{150, 250},
			[]text.Region{{150, 200}},
		},
		{
			[]text.Region{{100, 200}},
			text.Region{0, 250},
			[]text.Region{{100, 200}},
		},
		{
			[]text.Region{{100, 200}, {300, 400}},
			text.Region{0, 500},
			[]text.Region{{100, 200}, {300, 400}},
		},
		{
			[]text.Region{{100, 200}, {300, 400}},
			text.Region{150, 350},
			[]text.Region{{150, 200}, {300, 350}},
		},
		{
			[]text.Region{{100, 200}, {300, 400}},
			text.Region{150, 250},
			[]text.Region{{150, 200}},
		},
		{
			[]text.Region{{100, 200}, {300, 400}},
			text.Region{250, 350},
			[]text.Region{{300, 350}},
		},
	}

	for i, test := range tests {
		vr := ViewRegions{}
		vr.Regions.AddAll(test.regions)

		vr.Cull(test.cull)

		r := vr.Regions.Regions()

		if !reflect.DeepEqual(r, test.exp) {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.exp, r)
		}
	}
}

func TestViewRegionsClone(t *testing.T) {
	vr := ViewRegions{
		Scope: "testScope",
		Icon:  "testIcon",
		Flags: 100,
	}
	vr.Regions.AddAll([]text.Region{{0, 0}, {120, 300}, {24, 34}, {45, 40}})

	c := vr.Clone()
	if !reflect.DeepEqual(c, vr) {
		t.Errorf("Expected %+v, but got %+v", vr, c)
	}
}
