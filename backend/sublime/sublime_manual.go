package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
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

func sublime_Console(tu *py.Tuple, kwargs *py.Dict) (py.Object, error) {
	if tu.Size() != 1 {
		return nil, fmt.Errorf("Unexpected argument count: %d", tu.Size())
	}
	if i, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		log4go.Info("Python sez: %s", i)
	}
	return toPython(nil)
}

func init() {
	sublime_methods = append(sublime_methods, py.Method{Name: "console", Func: sublime_Console})
	backend.GetEditor()
	py.InitializeEx(false)
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
}

// TODO
func Init() {
	if m, err := py.Import("sublime_plugin"); err != nil {
		panic(err)
	} else {
		// scanpath("../packages/", m)
		scanpath("../../3rdparty/bundles/", m)
	}
}
