// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	// "code.google.com/p/log4go" Later
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
)

type (
	Package interface {
		// Returns the name of the package for sth
		// like settings is the filename and for
		// plugins is the dir name
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
		settings []*Setting
		keymaps  []*KeyMap
		path     string
		files    []os.FileInfo
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

// We store all scaned packages here with appropriate
// key like plugins, settings, keymaps, etc
// plugins specific settings or keymaps won't be in here
// we should access them from the plugin itself
var Packages = make(map[string][]Package)

// Initializes a new plugin whith loading all of the
// settings, keymaps and python files inside the path
func NewPlugin(path string) *Plugin {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	fi, err := f.Readdir(-1)
	if err != nil {
		return nil
	}
	p := &Plugin{path: path}
	files := make([]os.FileInfo, 0)
	for _, f := range fi {
		if strings.HasSuffix(f.Name(), ".py") {
			files = append(files, f)
		} else if strings.HasSuffix(f.Name(), ".sublime-settings") {
			p.settings = append(p.settings, NewSetting(path+string(os.PathSeparator)+f.Name()))
		} else if strings.HasSuffix(f.Name(), ".sublime-keymap") {
			p.keymaps = append(p.keymaps, NewKeyMap(path+string(os.PathSeparator)+f.Name()))
		}
	}
	p.files = files
	return p
}

func (p *Plugin) Get() interface{} {
	return p.files
}

func (p *Plugin) Name() string {
	return path.Base(p.path)
}

func (p *Plugin) Settings() []*Setting {
	return p.settings
}

func (p *Plugin) KeyMaps() []*KeyMap {
	return p.keymaps
}

func (p *Plugin) Reload() {
	f, err := os.Open(p.path)
	if err != nil {
		return
	}
	defer f.Close()
	fi, err := f.Readdir(-1)
	if err != nil {
		return
	}
	p.files = make([]os.FileInfo, 0)
	p.settings = make([]*Setting, 0)
	p.keymaps = make([]*KeyMap, 0)
	for _, f := range fi {
		if strings.HasSuffix(f.Name(), ".py") {
			p.files = append(p.files, f)
		} else if strings.HasSuffix(f.Name(), ".sublime-settings") {
			p.settings = append(p.settings, NewSetting(p.path+string(os.PathSeparator)+f.Name()))
		} else if strings.HasSuffix(f.Name(), ".sublime-keymap") {
			p.keymaps = append(p.keymaps, NewKeyMap(p.path+string(os.PathSeparator)+f.Name()))
		}
	}
}

func NewSetting(path string) *Setting {
	return &Setting{path, nil}
}

func (p *Setting) Get() interface{} {
	if p.data == nil {
		d, err := ioutil.ReadFile(p.path)
		if err != nil {
			return nil
		}
		p.data = d
	}
	return p.data
}

func (p *Setting) Name() string {
	return path.Base(p.path)
}

func (p *Setting) Reload() {
	d, err := ioutil.ReadFile(p.path)
	if err != nil {
		return
	}
	p.data = d
}

func NewKeyMap(path string) *KeyMap {
	return &KeyMap{path, nil}
}

func (p *KeyMap) Get() interface{} {
	if p.data == nil {
		d, err := ioutil.ReadFile(p.path)
		if err != nil {
			return nil
		}
		p.data = d
	}
	return p.data
}

func (p *KeyMap) Name() string {
	return path.Base(p.path)
}

func (p *KeyMap) Reload() {
	d, err := ioutil.ReadFile(p.path)
	if err != nil {
		return
	}
	p.data = d
}

func add(key string, p Package) {
	if !reflect.ValueOf(p).IsNil() {
		Packages[key] = append(Packages[key], p)
	}
}

func Scanpath(path string) []*Plugin {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	dirs, err := f.Readdirnames(-1)
	if err != nil {
		return nil
	}
	plugins := make([]*Plugin, 0)
	for _, dir := range dirs {
		dir2 := path + dir
		f2, err := os.Open(dir2)
		if err != nil {
			continue
		}
		defer f2.Close()
		fi, err := f2.Readdir(-1)
		if err != nil {
			continue
		}
		for _, f := range fi {
			if strings.HasSuffix(f.Name(), ".py") {
				plugins = append(plugins, NewPlugin(dir2))
				break
			}
		}
	}
	return plugins
}

// Loading the default packages
func init() {
	add("settings", NewSetting(DEFAULT_SUBLIME_SETTINGS))
	add("keymaps", NewKeyMap(DEFAULT_SUBLIME_KEYBINDINGS))
	plugins := Scanpath(SUBLIME_USER_PACKAGES_PATH)
	for _, p := range plugins {
		add("plugins", p)
	}
}
