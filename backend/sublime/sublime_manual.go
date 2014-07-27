// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	"github.com/limetext/lime/backend/render"
	"os"
	"path"
	"time"
)

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

func sublime_set_timeout(tu *py.Tuple, kwargs *py.Dict) (py.Object, error) {
	var (
		pyarg py.Object
	)
	if tu.Size() != 2 {
		return nil, fmt.Errorf("Unexpected argument count: %d", tu.Size())
	}
	if i, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		pyarg = i
	}
	if i, err := tu.GetItem(1); err != nil {
		return nil, err
	} else if v, err := fromPython(i); err != nil {
		return nil, err
	} else if v2, ok := v.(int); !ok {
		return nil, fmt.Errorf("Expected int not %s", i.Type())
	} else {
		pyarg.Incref()
		go func() {
			time.Sleep(time.Millisecond * time.Duration(v2))
			l := py.NewLock()
			defer l.Unlock()
			defer pyarg.Decref()
			if ret, err := pyarg.Base().CallFunctionObjArgs(); err != nil {
				log4go.Debug("Error in callback: %v", err)
			} else {
				ret.Decref()
			}
		}()
	}
	return toPython(nil)
}

func init() {
	sublime_methods = append(sublime_methods, py.Method{Name: "console", Func: sublime_Console}, py.Method{Name: "set_timeout", Func: sublime_set_timeout})
	backend.GetEditor()
	l := py.InitAndLock()
	defer l.Unlock()
	//	py.InitializeEx(false)
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
	type constant struct {
		name     string
		constant int
	}
	constants := []constant{
		{"OP_EQUAL", int(backend.OpEqual)},
		{"OP_NOT_EQUAL", int(backend.OpNotEqual)},
		{"OP_REGEX_MATCH", int(backend.OpRegexMatch)},
		{"OP_NOT_REGEX_MATCH", int(backend.OpNotRegexMatch)},
		{"OP_REGEX_CONTAINS", int(backend.OpRegexContains)},
		{"OP_NOT_REGEX_CONTAINS", int(backend.OpNotRegexContains)},
		{"INHIBIT_WORD_COMPLETIONS", 0},
		{"INHIBIT_EXPLICIT_COMPLETIONS", 0},
		{"LITERAL", 0},
		{"IGNORECASE", 0},
		{"CLASS_WORD_START", 1},
		{"CLASS_WORD_END", 2},
		{"CLASS_PUNCTUATION_START", 4},
		{"CLASS_PUNCTUATION_END", 8},
		{"CLASS_SUB_WORD_START", 16},
		{"CLASS_SUB_WORD_END", 32},
		{"CLASS_LINE_START", 64},
		{"CLASS_LINE_END", 128},
		{"CLASS_EMPTY_LINE", 256},
		{"CLASS_MIDDLE_WORD", 512},
		{"CLASS_WORD_START_WITH_PUNCTUATION", 1024},
		{"CLASS_WORD_END_WITH_PUNCTUATION", 2048},
		{"CLASS_OPENING_PARENTHESIS", 4096},
		{"CLASS_CLOSING_PARENTHESIS", 8192},
		{"DRAW_EMPTY", int(render.DRAW_EMPTY)},
		{"HIDE_ON_MINIMAP", int(render.HIDE_ON_MINIMAP)},
		{"DRAW_EMPTY_AS_OVERWRITE", int(render.DRAW_EMPTY_AS_OVERWRITE)},
		{"DRAW_NO_FILL", int(render.DRAW_NO_FILL)},
		{"DRAW_NO_OUTLINE", int(render.DRAW_NO_OUTLINE)},
		{"DRAW_SOLID_UNDERLINE", int(render.DRAW_SOLID_UNDERLINE)},
		{"DRAW_STIPPLED_UNDERLINE", int(render.DRAW_STIPPLED_UNDERLINE)},
		{"DRAW_SQUIGGLY_UNDERLINE", int(render.DRAW_SQUIGGLY_UNDERLINE)},
		{"PERSISTENT", int(render.PERSISTENT)},
		{"HIDDEN", int(render.HIDDEN)},
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
	for _, c := range constants {
		if err := m.AddIntConstant(c.name, c.constant); err != nil {
			panic(err)
		}
	}
	py.AddToPath("../../backend/packages/")
	py.AddToPath("../../3rdparty/bundles/")
	py.AddToPath("../../backend/sublime/")
}

func loadPlugin(p *backend.Plugin, m *py.Module) {
	fi := p.Get().([]os.FileInfo)
	for _, f := range fi {
		fn := f.Name()
		s, err := py.NewUnicode(path.Base(p.Name()) + "." + fn[:len(fn)-3])
		if err != nil {
			log4go.Error(err)
			return
		}
		if r, err := m.Base().CallMethodObjArgs("reload_plugin", s); err != nil {
			log4go.Error(err)
		} else if r != nil {
			r.Decref()
		}
	}
	p.LoadPackets()
	watch(backend.NewWatchedPackage(p))
}

var (
	watcher        *fsnotify.Watcher
	watchedPlugins map[string]*backend.WatchedPackage
)

func watch(plugin *backend.WatchedPackage) {
	log4go.Finest("Watch(%v)", plugin)
	if err := watcher.Watch(plugin.Name()); err != nil {
		log4go.Error("Could not watch plugin: %v", err)
	} else {
		watchedPlugins[plugin.Name()] = plugin
	}
}

func unWatch(name string) {
	if err := watcher.RemoveWatch(name); err != nil {
		log4go.Error("Couldn't unwatch file: %v", err)
	}
	log4go.Finest("UnWatch(%s)", name)
	delete(watchedPlugins, name)
}

func observePlugins(m *py.Module) {
	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsModify() {
				if p, exist := watchedPlugins[path.Dir(ev.Name)]; exist {
					p.Reload()
					loadPlugin(p.Package().(*backend.Plugin), m)
				}
			}
		case err := <-watcher.Error:
			log4go.Error("error:", err)
		}
	}
}

// TODO
func Init() {
	l := py.NewLock()
	defer l.Unlock()
	m, err := py.Import("sublime_plugin")
	if err != nil {
		panic(err)
	}
	sys, err := py.Import("sys")
	if err != nil {
		log4go.Debug(err)
	} else {
		defer sys.Decref()
	}

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log4go.Error("Could not create watcher due to: %v", err)
	}
	watchedPlugins = make(map[string]*backend.WatchedPackage)
	go observePlugins(m)

	plugins := backend.ScanPlugins(backend.LIME_USER_PACKAGES_PATH, ".py")
	for _, p := range plugins {
		// TODO: add all plugins after supporting all commands
		if p.Name() == "../../3rdparty/bundles/Vintageous" {
			loadPlugin(p, m)
		}
	}
}
