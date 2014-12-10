// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package render

import (
	"github.com/limetext/text"
	"reflect"
	"testing"
)

func TestViewRegionMapCull(t *testing.T) {
	tests := []struct {
		regionsA []text.Region
		regionsB []text.Region
		cull     text.Region
		expA     []text.Region
		expB     []text.Region
	}{
		{
			[]text.Region{},
			[]text.Region{},
			text.Region{},
			nil,
			nil,
		},
		{
			[]text.Region{{A: 100, B: 200}},
			[]text.Region{},
			text.Region{},
			nil,
			nil,
		},
		{
			[]text.Region{},
			[]text.Region{},
			text.Region{A: 0, B: 50},
			nil,
			nil,
		},
		{
			[]text.Region{{A: 100, B: 200}},
			[]text.Region{},
			text.Region{A: 0, B: 50},
			nil,
			nil,
		},
		{
			[]text.Region{{A: 100, B: 200}},
			[]text.Region{{A: 300, B: 400}},
			text.Region{A: 0, B: 50},
			nil,
			nil,
		},
		{
			[]text.Region{{A: 100, B: 200}},
			[]text.Region{{A: 300, B: 400}},
			text.Region{A: 0, B: 150},
			[]text.Region{{A: 100, B: 150}},
			nil,
		},
		{
			[]text.Region{{A: 100, B: 200}},
			[]text.Region{{A: 300, B: 400}},
			text.Region{A: 150, B: 350},
			[]text.Region{{A: 150, B: 200}},
			[]text.Region{{A: 300, B: 350}},
		},
	}

	for i, test := range tests {
		vrm := ViewRegionMap{}

		vrA := ViewRegions{}
		vrA.Regions.AddAll(test.regionsA)
		vrm["A"] = vrA

		vrB := ViewRegions{}
		vrB.Regions.AddAll(test.regionsB)
		vrm["B"] = vrB

		vrm.Cull(test.cull)

		if test.expA != nil {
			rA := vrA.Regions.Regions()

			if !reflect.DeepEqual(rA, test.expA) {
				t.Errorf("Test %d: Expected %s, but got %s", i, test.expA, rA)
			}
		} else if _, ok := vrm["A"]; ok {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.expA, vrm["A"])
		}

		if test.expB != nil {
			rB := vrB.Regions.Regions()

			if !reflect.DeepEqual(rB, test.expB) {
				t.Errorf("Test %d: Expected %s, but got %s", i, test.expB, rB)
			}
		} else if _, ok := vrm["B"]; ok {
			t.Errorf("Test %d: Expected %s, but got %s", i, test.expB, vrm["B"])
		}
	}
}

func TestViewRegionsCull(t *testing.T) {
	tests := []struct {
		regions []text.Region
		cull    text.Region
		exp     []text.Region
	}{
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 0, B: 50},
			[]text.Region{},
		},
		{
			[]text.Region{{A: 100, B: 100}},
			text.Region{A: 100, B: 100},
			[]text.Region{},
		},
		{
			[]text.Region{{A: 100, B: 100}},
			text.Region{A: 95, B: 105},
			[]text.Region{},
		},
		{
			[]text.Region{{A: 100, B: 100}},
			text.Region{A: 95, B: 100},
			[]text.Region{},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 150, B: 150},
			[]text.Region{},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 90, B: 100},
			[]text.Region{},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 100, B: 150},
			[]text.Region{{A: 100, B: 150}},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 150, B: 175},
			[]text.Region{{A: 150, B: 175}},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 0, B: 150},
			[]text.Region{{A: 100, B: 150}},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 150, B: 250},
			[]text.Region{{A: 150, B: 200}},
		},
		{
			[]text.Region{{A: 100, B: 200}},
			text.Region{A: 0, B: 250},
			[]text.Region{{A: 100, B: 200}},
		},
		{
			[]text.Region{{A: 100, B: 200}, {A: 300, B: 400}},
			text.Region{A: 0, B: 500},
			[]text.Region{{A: 100, B: 200}, {A: 300, B: 400}},
		},
		{
			[]text.Region{{A: 100, B: 200}, {A: 300, B: 400}},
			text.Region{A: 150, B: 350},
			[]text.Region{{A: 150, B: 200}, {A: 300, B: 350}},
		},
		{
			[]text.Region{{A: 100, B: 200}, {A: 300, B: 400}},
			text.Region{A: 150, B: 250},
			[]text.Region{{A: 150, B: 200}},
		},
		{
			[]text.Region{{A: 100, B: 200}, {A: 300, B: 400}},
			text.Region{A: 250, B: 350},
			[]text.Region{{A: 300, B: 350}},
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
