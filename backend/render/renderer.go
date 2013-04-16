package render

import (
	"image"
	"image/color"
	"lime/backend/primitives"
)

type (
	Flavour struct {
		Background color.RGBA
		Foreground color.RGBA
		Font       Font
	}

	RenderUnit struct {
		Flavour Flavour
		Region  primitives.Region
		Layout  image.Rectangle
	}

	Recipe map[Flavour]primitives.RegionSet

	ColourScheme interface {
		Spice(ViewRegions) Flavour
	}

	Renderer struct {
	}
)

func (r *Renderer) Transform(scheme ColourScheme, data ViewRegionMap, viewport primitives.Region) Recipe {
	data.Cull(viewport)
	recipe := make(Recipe)
	for _, v := range data {
		k := scheme.Spice(v)
		rs := recipe[k]
		rs.AddAll(v.Regions.Regions())
		recipe[k] = rs
	}
	return recipe
}

func (r Recipe) Transcribe() []RenderUnit {
	return nil
}
