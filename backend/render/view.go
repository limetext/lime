// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package render

import (
	"github.com/limetext/lime/backend/util"
	"github.com/limetext/text"
)

const (
	DRAW_EMPTY              ViewRegionFlags = (1 << iota) // Draw a vertical line for an empty (zero area) region
	HIDE_ON_MINIMAP                                       // Don't draw this region in the minimap
	DRAW_EMPTY_AS_OVERWRITE                               // Rather than a vertical line, draw empty regions as a horizontal one
	DRAW_NO_FILL                                          // Don't draw the filling of the region
	DRAW_NO_OUTLINE                                       // Don't draw the outline of the region
	DRAW_SOLID_UNDERLINE                                  // Draw a solid underline under the whole region
	DRAW_STIPPLED_UNDERLINE                               // Draw a stippled underline under the whole region
	DRAW_SQUIGGLY_UNDERLINE                               // Draw a squiggly underline under the whole region
	PERSISTENT                                            // Region is saved with the session
	HIDDEN                                                // Region is not rendered
	SELECTION                                             // This Region is part of selected text
	HIGHLIGHT                                             // This Region is part of highlighted text
	DRAW_TEXT                                             // The actual text contained in the region should be rendered
	DEFAULT                 ViewRegionFlags = 0           // No flags at all, only draw the region itself and not the text
)

type (
	// A set of ViewRegions associated by a string identifier key.
	// The name of the key itself has no special meaning other than for
	// being able set, retrieve and update *your* set of ViewRegions.
	ViewRegionMap map[string]ViewRegions

	// Flags used to hint at how the region should be rendered.
	ViewRegionFlags int

	// The ViewRegions object contains information
	// related to the rendering of a specific RegionSet
	// and can be set both by a https://godoc.org/github.com/limetext/lime/backend/parser#SyntaxHighlighter
	// and from plugins via https://godoc.org/github.com/limetext/lime/backend#View.AddRegions.
	//
	// Turning this information into a concrete https://godoc.org/github.com/limetext/lime/backend/render#Flavour
	// is the job of the https://godoc.org/github.com/limetext/lime/backend/render#ColourScheme interface.
	ViewRegions struct {
		// The Regions this ViewRegions object is relevant to.
		Regions text.RegionSet
		// The scope identifier is used to determine colour and other style options.
		Scope string
		// Gutter icon (displayed next to line numbers) URI.
		Icon string
		// Flags used to hint at how the region should be rendered.
		Flags ViewRegionFlags
	}
)

// Calls Cull on each ViewRegions object contained in the map,
// removing all entries that are outside of the viewport.
func (vrm *ViewRegionMap) Cull(viewport text.Region) {
	pe := util.Prof.Enter("render.vrm.Cull")
	defer pe.Exit()
	rm := []string{}
	for k, v := range *vrm {
		v.Cull(viewport)
		if v.Regions.Len() == 0 {
			rm = append(rm, k)
		} else {
			(*vrm)[k] = v
		}
	}
	for _, r := range rm {
		delete(*vrm, r)
	}
}

// Removes any regions that are outside of the given viewport,
// and clips the regions that are intersecting it so that
// all regions remaining are fully contained inside of the viewport.
func (vr *ViewRegions) Cull(viewport text.Region) {
	pe := util.Prof.Enter("render.vr.Cull")
	defer pe.Exit()
	nr := []text.Region{}
	for _, r := range vr.Regions.Regions() {
		if viewport.Intersects(r) {
			in := viewport.Intersection(r)
			if in.Size() != 0 {
				nr = append(nr, in)
			}
		}
	}
	vr.Regions.Clear()
	vr.Regions.AddAll(nr)
}

// Creates a copy of this ViewRegions object.
func (vr *ViewRegions) Clone() ViewRegions {
	ret := ViewRegions{Scope: vr.Scope, Icon: vr.Icon, Flags: vr.Flags}
	ret.Regions.AddAll(vr.Regions.Regions())
	return ret
}
