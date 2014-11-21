// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/packages"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/lime/backend/util"
	"github.com/limetext/lime/backend/watch"
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
		log.Info("Python sez: %s", i)
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
				log.Debug("Error in callback: %v", err)
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
		{"OP_EQUAL", int(util.OpEqual)},
		{"OP_NOT_EQUAL", int(util.OpNotEqual)},
		{"OP_REGEX_MATCH", int(util.OpRegexMatch)},
		{"OP_NOT_REGEX_MATCH", int(util.OpNotRegexMatch)},
		{"OP_REGEX_CONTAINS", int(util.OpRegexContains)},
		{"OP_NOT_REGEX_CONTAINS", int(util.OpNotRegexContains)},
		{"INHIBIT_WORD_COMPLETIONS", 0},
		{"INHIBIT_EXPLICIT_COMPLETIONS", 0},
		{"LITERAL", 0},
		{"IGNORECASE", 0},
		{"CLASS_WORD_START", int(backend.CLASS_WORD_START)},
		{"CLASS_WORD_END", int(backend.CLASS_WORD_END)},
		{"CLASS_PUNCTUATION_START", int(backend.CLASS_PUNCTUATION_START)},
		{"CLASS_PUNCTUATION_END", int(backend.CLASS_PUNCTUATION_END)},
		{"CLASS_SUB_WORD_START", int(backend.CLASS_SUB_WORD_START)},
		{"CLASS_SUB_WORD_END", int(backend.CLASS_SUB_WORD_END)},
		{"CLASS_LINE_START", int(backend.CLASS_LINE_START)},
		{"CLASS_LINE_END", int(backend.CLASS_LINE_END)},
		{"CLASS_EMPTY_LINE", int(backend.CLASS_EMPTY_LINE)},
		{"CLASS_MIDDLE_WORD", int(backend.CLASS_MIDDLE_WORD)},
		{"CLASS_WORD_START_WITH_PUNCTUATION", int(backend.CLASS_WORD_START_WITH_PUNCTUATION)},
		{"CLASS_WORD_END_WITH_PUNCTUATION", int(backend.CLASS_WORD_END_WITH_PUNCTUATION)},
		{"CLASS_OPENING_PARENTHESIS", int(backend.CLASS_OPENING_PARENTHESIS)},
		{"CLASS_CLOSING_PARENTHESIS", int(backend.CLASS_CLOSING_PARENTHESIS)},
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
	py.AddToPath(backend.LIME_PACKAGES_PATH)
	py.AddToPath(backend.LIME_USER_PACKAGES_PATH)
	py.AddToPath(path.Join("..", "..", "backend", "sublime"))
}

// Wrapper for packages.Plugin and py.Module
// merges Plugin.Reload and loadPlugin for watcher
type plugin struct {
	pl *packages.Plugin
	m  *py.Module
}

func newPlugin(p *packages.Plugin, m *py.Module) *plugin {
	return &plugin{p, m}
}

func (p *plugin) FileChanged(name string) {
	p.pl.Reload()
	p.loadPlugin()
}

func (p *plugin) Name() string {
	return p.pl.Name()
}

func (p *plugin) loadPlugin() {
	fi := p.pl.Get().([]os.FileInfo)
	for _, f := range fi {
		fn := f.Name()
		s, err := py.NewUnicode(path.Base(p.pl.Name()) + "." + fn[:len(fn)-3])
		if err != nil {
			log.Error(err)
			return
		}
		if r, err := p.m.Base().CallMethodObjArgs("reload_plugin", s); err != nil {
			log.Error(err)
		} else if r != nil {
			r.Decref()
		}
	}
	p.pl.LoadPackets()
}

var watcher *watch.Watcher

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
		log.Debug(err)
	} else {
		defer sys.Decref()
	}

	if watcher, err = watch.NewWatcher(); err != nil {
		log.Error("Couldn't create watcher: %s", err)
	}

	plugins := packages.ScanPlugins(backend.LIME_USER_PACKAGES_PATH, ".py")
	for _, p := range plugins {
		// TODO: add all plugins after supporting all commands
		if p.Name() == path.Join("..", "..", "packages", "Vintageous") {
			pl := newPlugin(p, m)
			pl.loadPlugin()
			if err := watcher.Watch(pl.Name(), pl); err != nil {
				log.Error("Couldn't watch %s: %s", pl.Name(), err)
			}
		}
	}

	go watcher.Observe()
}
