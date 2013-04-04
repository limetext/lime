package sublime

import (
	"lime/3rdparty/libs/gopy/lib"
	"log"
)

func init() {
	methods := []py.Method{}
	py.Initialize()
	m, err := py.InitModule("sublime", methods)
	if err != nil {
		log.Fatal(err)
	}

	type class struct {
		name string
		c    *py.Class
	}
	classes := []class{
		{"Region", &_regionClass},
		{"RegionSet", &_region_setClass},
		{"View", &_viewClass},
		{"Window", &_windowClass},
		{"Edit", &_editClass},
		{"Settings", &_settingsClass},
	}
	for _, cl := range classes {
		c, err := cl.c.Create()
		if err != nil {
			log.Fatal(err)
		}
		if err := m.AddObject(cl.name, c); err != nil {
			log.Fatal(err)
		}
	}
}
