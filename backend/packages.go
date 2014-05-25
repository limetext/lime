// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type (
	Package interface {
		// Returns the path of the package
		Name() string

		// Returns the useful data that we need
		// from this package for example for a
		// plugin will be the python files or for
		// a keymap will be the file data
		Get() interface{}

		// Reloads the package data
		Reload()
	}

	Plugin struct {
		path     string
		suffix   string
		files    []os.FileInfo
		settings []*Setting
		keymaps  []*KeyMap
	}

	Setting struct {
		path string
		data []byte
	}

	KeyMap struct {
		path string
		data []byte
	}
)

const (
	DEFAULT_SUBLIME_SETTINGS    = "../../backend/packages/Default/Default.sublime-settings"
	DEFAULT_SUBLIME_KEYBINDINGS = "../../backend/packages/Default/Default.sublime-keymap"
	SUBLIME_USER_PACKAGES_PATH  = "../../3rdparty/bundles/"
)

// We store default packages here with appropriate
// key like plugins, settings, keymaps, etc plugins
// specific settings or keymaps won't be in here
// they should be accessed from the plugin itself
var Packages = make(map[string][]Package)

// Initializes a new plugin whith loading all of the
// settings, keymaps and python files inside the path
func NewPlugin(path string, suffix string) *Plugin {
	f, err := os.Open(path)
	if err != nil {
		log4go.Error("Couldn't open dir: %s", err)
		return nil
	}
	defer f.Close()
	fi, err := f.Readdir(-1)
	if err != nil {
		log4go.Error("Couldn't read dir: %s", err)
		return nil
	}
	files := make([]os.FileInfo, 0)
	sets := make([]*Setting, 0)
	keys := make([]*KeyMap, 0)
	for _, f := range fi {
		if strings.HasSuffix(f.Name(), suffix) {
			files = append(files, f)
		} else if strings.HasSuffix(f.Name(), ".sublime-settings") {
			sets = append(sets, NewSetting(path+string(os.PathSeparator)+f.Name()))
		} else if strings.HasSuffix(f.Name(), ".sublime-keymap") {
			keys = append(keys, NewKeyMap(path+string(os.PathSeparator)+f.Name()))
		}
	}
	return &Plugin{path, suffix, files, sets, keys}
}

func (p *Plugin) Get() interface{} {
	return p.files
}

func (p *Plugin) Name() string {
	return p.path
}

func (p *Plugin) Settings() []*Setting {
	return p.settings
}

func (p *Plugin) KeyMaps() []*KeyMap {
	return p.keymaps
}

func (p *Plugin) Reload() {
	p = NewPlugin(p.path, p.suffix)
	for _, s := range p.settings {
		s.Reload()
	}
	for _, k := range p.keymaps {
		k.Reload()
	}
}

func NewSetting(path string) *Setting {
	return &Setting{path, nil}
}

func loadData(path string) []byte {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log4go.Error("Couldn't read file: %s", err)
		return nil
	}
	return d
}

func (p *Setting) Get() interface{} {
	if p.data == nil {
		p.data = loadData(p.path)
	}
	return p.data
}

func (p *Setting) Name() string {
	return p.path
}

func (p *Setting) Reload() {
	p.data = loadData(p.path)
	e := GetEditor()
	e.loadSetting(p)
}

func NewKeyMap(path string) *KeyMap {
	return &KeyMap{path, nil}
}

func (p *KeyMap) Get() interface{} {
	if p.data == nil {
		p.data = loadData(p.path)
	}
	return p.data
}

func (p *KeyMap) Name() string {
	return p.path
}

func (p *KeyMap) Reload() {
	p.data = loadData(p.path)
	e := GetEditor()
	e.loadKeybinding(p)
}

func add(key string, p Package) {
	if !reflect.ValueOf(p).IsNil() {
		Packages[key] = append(Packages[key], p)
	}
}

// Loading the default packages
func init() {
	add("settings", NewSetting(DEFAULT_SUBLIME_SETTINGS))
	add("keymaps", NewKeyMap(DEFAULT_SUBLIME_KEYBINDINGS))
	// TODO: for now we just load vintageous plugin but we should
	// scan this path and add all plugins
	add("plugins", NewPlugin(SUBLIME_USER_PACKAGES_PATH+string(os.PathSeparator)+"Vintageous", ".py"))
}
