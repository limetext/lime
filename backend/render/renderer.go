package render

import (
	"github.com/quarnster/util/text"
	"image"
	"image/color"
)

type (
	Colour  color.RGBA
	Flavour struct {
		Background Colour
		Foreground Colour
		Font       Font
	}

	RenderUnit struct {
		Flavour Flavour
		Region  text.Region
		Layout  image.Rectangle
	}

	Recipe map[Flavour]text.RegionSet

	ColourScheme interface {
		Spice(*ViewRegions) Flavour
	}

	Renderer struct {
	}
)

func (r *Renderer) Transform(scheme ColourScheme, data ViewRegionMap, viewport text.Region) Recipe {
	data.Cull(viewport)
	recipe := make(Recipe)
	for _, v := range data {
		k := scheme.Spice(&v)
		rs := recipe[k]
		rs.AddAll(v.Regions.Regions())
		recipe[k] = rs
	}
	return recipe
}

func (r Recipe) Transcribe() []RenderUnit {
	return nil
}
