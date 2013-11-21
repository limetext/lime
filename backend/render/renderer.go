package render

import (
	"github.com/quarnster/util/text"
	"image/color"
	"sort"
)

type (
	// Colour represented by a underlying color.RGBA structure
	Colour color.RGBA

	// The Flavour struct contains the specific settings
	// used to style a particular Region.
	Flavour struct {
		Background Colour
		Foreground Colour
		Font       Font
		Flags      ViewRegionFlags
	}

	// The Recipe type groups text.RegionSets by their Flavour.
	// The idea is to allow large groups of text be rendered as
	// a single batch without any state changes inbetween the
	// batches.
	Recipe map[Flavour]text.RegionSet

	// A RenderUnit is just a Flavour and an associated Region.
	RenderUnit struct {
		Flavour Flavour
		Region  text.Region
	}

	// A TranscribedRecipe is a linear (in text.Regions) representation of a Recipe
	TranscribedRecipe []RenderUnit

	ColourScheme interface {
		// Takes a ViewRegions pointer as input and uses the data contained in it
		// to determine the Flavour it should be rendered with.
		Spice(*ViewRegions) Flavour
	}

	Renderer interface {
		// Renders the given Recipe
		Render(Recipe)
	}
)

// Transform takes a ColourScheme, a ViewRegionMap and a viewport as input.
//
// The viewport would be the text.Region of the current buffer that is visible to the user
// and any ViewRegions outside of this area are not forwarded for further processing.
//
// The remaining ViewRegions are then passed on to the ColourScheme for determining the exact Flavour
// for which that RegionSet should be styled, adding Regions of the same Flavour to the same RegionSet.
//
// Typically there are more ViewRegions available in a text buffer than there are unique Flavours in
// a ColourScheme, so this operation can be viewed as reducing the number of state changes required to
// display the text to the user.
//
// The final output, the Recipe, contains a mapping of all unique Flavours and that Flavour's
// associated RegionSet.
func Transform(scheme ColourScheme, data ViewRegionMap, viewport text.Region) Recipe {
	// TODO:
	// 	caret_blink := true
	// if b, ok := v.Settings().Get("caret_blink", true).(bool); ok {
	// 	caret_blink = b
	// }
	//
	// highlight_line := false
	// if b, ok := v.Settings().Get("highlight_line", highlight_line).(bool); ok {
	// 	highlight_line = b
	// }
	//	if b, ok := v.Settings().Get("inverse_caret_state", false).(bool); !b && ok {
	// 	if caret_style == termbox.AttrReverse {
	// 		caret_style = termbox.AttrUnderline
	// 	} else {
	// 		caret_style = termbox.AttrReverse
	// 	}
	// }
	// caret_style := termbox.AttrUnderline
	// if b, ok := v.Settings().Get("caret_style", "underline").(string); ok {
	// 	if b == "block" {
	// 		caret_style = termbox.AttrReverse
	// 	}
	// }

	data.Cull(viewport)
	recipe := make(Recipe)
	for _, v := range data {
		k := scheme.Spice(&v)
		rs := recipe[k]
		rs.AddAll(v.Regions.Regions())
		if rs.HasNonEmpty() {
			var last text.Region
			rgs := rs.Regions()
			for i, r := range rgs {
				if i > 0 && r.Begin() == last.End() {
					rs.Add(r.Cover(last))
				}
				last = r
			}
			recipe[k] = rs
		}
	}
	return recipe
}

// Transcribing the Recipe creates a linear step-by-step
// representation of it, which might or might not
// make it easier for Renderers to work with.
func (r Recipe) Transcribe() (ret TranscribedRecipe) {
	for flav, set := range r {
		for _, r := range set.Regions() {
			ret = append(ret, RenderUnit{Flavour: flav, Region: r})
		}
	}
	sort.Sort(&ret)
	return
}

// Just used to satisfy the sort.Interface interface, typically not used otherwise.
func (r *TranscribedRecipe) Len() int {
	return len(*r)
}

// Just used to satisfy the sort.Interface interface, typically not used otherwise.
func (r *TranscribedRecipe) Less(i, j int) bool {
	a, b := (*r)[i].Region, (*r)[j].Region
	if a.Begin() == b.Begin() {
		return a.End() > b.End()
	}
	return a.Begin() < b.Begin()
}

// Just used to satisfy the sort.Interface interface, typically not used otherwise.
func (r *TranscribedRecipe) Swap(i, j int) {
	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
}
