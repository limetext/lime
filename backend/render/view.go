package render

import (
	"github.com/quarnster/util/text"
)

const (
	DRAW_EMPTY ViewRegionFlags = (1 << iota)
	HIDE_ON_MINIMAP
	DRAW_EMPTY_AS_OVERWRITE
	DRAW_NO_FILL
	DRAW_NO_OUTLINE
	DRAW_SOLID_UNDERLINE
	DRAW_STIPPLED_UNDERLINE
	DRAW_SQUIGGLY_UNDERLINE
	PERSISTENT
	HIDDEN
	FOREGROUND
	SELECTION
	HIGHLIGHT
	DRAW_TEXT
	DEFAULT ViewRegionFlags = 0
)

type (
	ViewRegionMap   map[string]ViewRegions
	ViewRegionFlags int
	ViewRegions     struct {
		Regions text.RegionSet
		Scope   string
		Icon    string
		Flags   ViewRegionFlags
	}
)

func (vrm *ViewRegionMap) Cull(viewport text.Region) {
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

func (vr *ViewRegions) Cull(viewport text.Region) {
	nr := []text.Region{}
	for _, r := range vr.Regions.Regions() {
		if viewport.Intersects(r) {
			nr = append(nr, viewport.Intersection(r))
		}
	}
	vr.Regions.Clear()
	vr.Regions.AddAll(nr)
}
