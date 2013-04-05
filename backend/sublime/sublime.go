package sublime

import (
	"code.google.com/p/log4go"
	"lime/3rdparty/libs/gopy/lib"
	"log"
	"os"
	"strings"
)

func scanpath(path string, m *py.Module) {
	// This should probably be done by the Editor as it needs to scan through for themes, keybinding, settings etc
	if f, err := os.Open(path); err != nil {
		log4go.Warn(err)
	} else {
		defer f.Close()
		if dirs, err := f.Readdirnames(-1); err != nil {
			log4go.Warn(err)
		} else {
			for _, dir := range dirs {
				if dir != "Vintageous" && dir != "Default" {
					// TODO obviously
					continue
				}
				dir2 := path + dir
				if f2, err := os.Open(dir2); err != nil {
					log4go.Warn(err)
				} else {
					defer f2.Close()
					if fi, err := f2.Readdir(-1); err != nil {
						log4go.Warn(err)
					} else {
						for _, f := range fi {
							if fn := f.Name(); strings.HasSuffix(fn, ".py") {
								m.Base().CallMethod("reload_plugin", "s", dir+"."+fn[:len(fn)-3])
							}
						}
					}
				}
			}
		}
	}
}

func init() {
	py.Initialize()

	m, err := py.InitModule("sublime", sublime_methods)
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
		{"WindowCommandGlue", &_windowCommandGlueClass},
		{"TextCommandGlue", &_textCommandGlueClass},
		{"ApplicationCommandGlue", &_applicationCommandGlueClass},
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
	py.AddToPath("../packages/")
	py.AddToPath("../../3rdparty/bundles/")
	py.AddToPath(".")
	if m, err := py.Import("sublime_plugin"); err != nil {
		log.Fatal(err)
	} else {
		scanpath("../packages/", m)
		scanpath("../../3rdparty/bundles/", m)
	}
}
