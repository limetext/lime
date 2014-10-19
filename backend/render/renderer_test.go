// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package render

import (
	"github.com/limetext/text"
	"reflect"
	"testing"
)

type (
	DummyColourScheme struct {
		ColourScheme
	}
)

var (
	flavourDef = Flavour{Background: Colour{0, 0, 0, 0}}
	flavourA   = Flavour{Background: Colour{1, 0, 0, 0}}
	flavourB   = Flavour{Background: Colour{0, 1, 0, 0}}
	flavourC   = Flavour{Background: Colour{0, 0, 1, 0}}
)

func (cs DummyColourScheme) Spice(vr *ViewRegions) (f Flavour) {
	switch vr.Scope {
	case "A":
		f = flavourA
	case "B":
		f = flavourB
	case "C":
		f = flavourC
	default:
		f = flavourDef
	}

	return
}

func TestTransform(t *testing.T) {
	tests := []struct {
		vrmap    ViewRegionMap
		rmap     map[string][]text.Region
		viewport text.Region
		expflav  []Flavour
		notflav  []Flavour
	}{
		{
			ViewRegionMap{"A": {Scope: "A"}},
			map[string][]text.Region{"A": {{0, 1}}},
			text.Region{0, 1},
			[]Flavour{flavourA},
			[]Flavour{flavourDef},
		},
		{
			ViewRegionMap{"A": {Scope: "A"}, "B": {Scope: "B"}},
			map[string][]text.Region{"A": {{0, 1}}},
			text.Region{0, 1},
			[]Flavour{flavourA},
			[]Flavour{flavourB},
		},
		{
			ViewRegionMap{"A": {Scope: "A"}, "B": {Scope: "B"}},
			map[string][]text.Region{"B": {{0, 1}}},
			text.Region{0, 1},
			[]Flavour{flavourB},
			[]Flavour{flavourA},
		},
		{
			ViewRegionMap{"A": {Scope: "A"}, "B": {Scope: "B"}},
			map[string][]text.Region{"A": {{0, 1}}, "B": {{0, 1}}},
			text.Region{0, 1},
			[]Flavour{flavourA, flavourB},
			[]Flavour{},
		},
		{
			ViewRegionMap{"A": {Scope: "A"}, "B": {Scope: "B"}},
			map[string][]text.Region{"A": {{10, 11}}, "B": {{0, 1}}},
			text.Region{0, 1},
			[]Flavour{flavourB},
			[]Flavour{flavourA},
		},
		{
			ViewRegionMap{"A": {Scope: "A"}, "B": {Scope: "B"}},
			map[string][]text.Region{"A": {{10, 11}}, "B": {{20, 21}}},
			text.Region{0, 1},
			[]Flavour{},
			[]Flavour{flavourA, flavourB},
		},
	}

	for i, test := range tests {
		for k, reg := range test.rmap {
			if vr, ok := test.vrmap[k]; ok {
				vr.Regions.AddAll(reg)
				test.vrmap[k] = vr
			}
		}

		rec := Transform(DummyColourScheme{}, test.vrmap, test.viewport)

		for _, f := range test.expflav {
			if _, ok := rec[f]; !ok {
				t.Errorf("Test %d: Expected Recipe to contain %#v, but it didn't", i, f)
			}
		}

		for _, f := range test.notflav {
			if _, ok := rec[f]; ok {
				t.Errorf("Test %d: Expected Recipe to not contain %#v, but it did", i, f)
			}
		}
	}
}

func TestRecipeTranscribe(t *testing.T) {
	tests := []struct {
		rec  Recipe
		rmap map[Flavour][]text.Region
		exp  TranscribedRecipe
	}{
		{
			Recipe{flavourA: text.RegionSet{}},
			map[Flavour][]text.Region{flavourA: {{0, 1}}},
			TranscribedRecipe{{flavourA, text.Region{0, 1}}},
		},
		{
			Recipe{flavourA: text.RegionSet{}},
			map[Flavour][]text.Region{flavourA: {{0, 1}, {10, 11}}},
			TranscribedRecipe{{flavourA, text.Region{0, 1}}, {flavourA, text.Region{10, 11}}},
		},
		{
			Recipe{flavourA: text.RegionSet{}, flavourB: text.RegionSet{}},
			map[Flavour][]text.Region{flavourA: {{0, 1}}},
			TranscribedRecipe{{flavourA, text.Region{0, 1}}},
		},
		{
			Recipe{flavourA: text.RegionSet{}, flavourB: text.RegionSet{}},
			map[Flavour][]text.Region{flavourA: {{0, 1}}, flavourB: {{10, 11}}},
			TranscribedRecipe{{flavourA, text.Region{0, 1}}, {flavourB, text.Region{10, 11}}},
		},
		{
			Recipe{flavourA: text.RegionSet{}, flavourB: text.RegionSet{}},
			map[Flavour][]text.Region{flavourA: {{0, 1}}, flavourB: {{0, 1}}},
			TranscribedRecipe{{flavourA, text.Region{0, 1}}, {flavourB, text.Region{0, 1}}},
		},
	}

	for i, test := range tests {
		for k, reg := range test.rmap {
			if rs, ok := test.rec[k]; ok {
				rs.AddAll(reg)
				test.rec[k] = rs
			}
		}

		tr := test.rec.Transcribe()
		if !reflect.DeepEqual(tr, test.exp) {
			t.Errorf("Test %d: Expected %#v, but got %#v", i, test.exp, tr)
		}
	}
}
