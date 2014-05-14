// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	// "code.google.com/p/log4go" Later
	"fmt"
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
		// Returns the path that the package exists
		Path() string
	}

	Plugin struct {
		setting *Setting
		keymap  *KeyMap
		path    string
		files   []os.FileInfo
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

var Packages = make(map[string][]Package)

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
			p.setting = NewSetting(path + string(os.PathSeparator) + f.Name())
		} else if strings.HasSuffix(f.Name(), ".sublime-keymap") {
			p.keymap = NewKeyMap(path + string(os.PathSeparator) + f.Name())
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

func (p *Plugin) Path() string {
	return p.path
}

func (p *Plugin) Setting() *Setting {
	return p.setting
}

func (p *Plugin) KeyMap() *KeyMap {
	return p.keymap
}

func NewSetting(path string) *Setting {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return &Setting{path, d}
}

func (p *Setting) Get() interface{} {
	return p.data
}

func (p *Setting) Name() string {
	return path.Base(p.path)
}

func (p *Setting) Path() string {
	return p.path
}

func NewKeyMap(path string) *KeyMap {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &KeyMap{path, d}
}

func (p *KeyMap) Get() interface{} {
	return p.data
}

func (p *KeyMap) Name() string {
	return path.Base(p.path)
}

func (p *KeyMap) Path() string {
	return p.path
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

func init() {
	add("settings", NewSetting(DEFAULT_SUBLIME_SETTINGS))
	add("keymaps", NewKeyMap(DEFAULT_SUBLIME_KEYBINDINGS))

	f, err := os.Open(SUBLIME_USER_PACKAGES_PATH)
	if err != nil {
		return
	}
	defer f.Close()
	dirs, err := f.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, dir := range dirs {
		dir2 := SUBLIME_USER_PACKAGES_PATH + dir
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
				add("plugins", NewPlugin(dir2))
				break
			}
		}
	}
}
