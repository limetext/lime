package sublime

import (
	"code.google.com/p/log4go"
	"lime/3rdparty/libs/gopy/lib"
	"os"
	"strings"
)

func scanpath(path string, m *py.Module) {
	sys, err := py.Import("sys")
	if err != nil {
		log4go.Debug(err)
	} else {
		defer sys.Decref()
	}

	// This should probably be done by the Editor as it needs to scan through for themes, keybinding, settings etc
	if f, err := os.Open(path); err != nil {
		log4go.Warn(err)
	} else {
		defer f.Close()
		if dirs, err := f.Readdirnames(-1); err != nil {
			log4go.Warn(err)
		} else {
			for _, dir := range dirs {
				if dir != "Vintageous" && dir != "Default" && dir != "plugins" {
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
								//m.Incref()
								if s, err := py.NewUnicode(dir + "." + fn[:len(fn)-3]); err != nil {
									log4go.Error(err)
								} else {
									if r, err := m.Base().CallMethodObjArgs("reload_plugin", s); err != nil {
										log4go.Error(err)
									} else if r != nil {
										r.Decref()
									}
									// if i, err := sys.Base().CallMethodObjArgs("getrefcount", s); err != nil {
									// 	log4go.Error(err)
									// } else {
									// 	log4go.Debug("m refs: %d", i.(*py.Long).Int64())
									// 	i.Decref()
									// }
								}
							}
						}
					}
				}
			}
		}
	}
}

func Init() {
	py.Initialize()

	m, err := py.InitModule("sublime", sublime_methods)
	if err != nil {
		panic(err)
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
		{"OnQueryContextGlue", &_onQueryContextGlueClass},
		{"ViewEventGlue", &_viewEventGlueClass},
	}
	for _, cl := range classes {
		c, err := cl.c.Create()
		if err != nil {
			panic(err)
		}
		if err := m.AddObject(cl.name, c); err != nil {
			panic(err)
		}
	}
	py.AddToPath("../../backend/packages/")
	py.AddToPath("../../3rdparty/bundles/")
	py.AddToPath("../../backend/sublime/")
	if m, err := py.Import("sublime_plugin"); err != nil {
		panic(err)
	} else {
		// scanpath("../packages/", m)
		scanpath("../../3rdparty/bundles/", m)
	}
}
