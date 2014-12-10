// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/text"
	"os"
	pt "path"
	"strings"
)

type (
	// Plugin is a Package type containing some files
	// with specific suffix that could be interpreted by
	// lime text api(currently python) and some
	// settings, snippets, commands and etc as packets
	Plugin struct {
		text.HasSettings
		keys.HasKeyBindings
		path             string
		suffix           string
		files            []os.FileInfo
		defaultSettings  *text.HasSettings
		platformSettings *text.HasSettings
		defaultBindings  *keys.HasKeyBindings
	}
)

// Initializes a new plugin whith loading all of the
// settings, keymaps and etc. Suffix variable show's
// which file types we need for plugin for example if
// the plugin is written in python the suffix should
// be ".py". We will use this function at initialization
// to add user plugins and on new_plugin command
func NewPlugin(path string, suffix string) (p *Plugin) {
	p = &Plugin{path: path, suffix: suffix}
	p.defaultSettings = new(text.HasSettings)
	p.platformSettings = new(text.HasSettings)
	p.defaultBindings = new(keys.HasKeyBindings)

	p.Settings().SetParent(p.platformSettings)
	p.platformSettings.Settings().SetParent(p.defaultSettings)

	p.KeyBindings().SetParent(p.defaultBindings)
	return
}

func (p *Plugin) Name() string {
	return p.path
}

// Returns slice of files with plugin suffix
// loaded at initialization
func (p *Plugin) Get() interface{} {
	return p.files
}

// On plugin reload we will scan for plugin files
// and packets in plugin path
func (p *Plugin) Reload() {
	var files []os.FileInfo
	log.Info("Reloading plugin %s", p.Name())
	f, err := os.Open(p.path)
	if err != nil {
		log.Errorf("Couldn't open dir: %s", err)
		return
	}
	defer f.Close()
	fi, err := f.Readdir(-1)
	if err != nil {
		log.Errorf("Couldn't read dir: %s", err)
		return
	}
	for _, f := range fi {
		if p.suffix != "" && strings.HasSuffix(f.Name(), p.suffix) {
			files = append(files, f)
		}
	}
	p.files = files
}

// Scaning path for finding plugins that contain files
// whith specific suffix
func ScanPlugins(path string, suffix string) []*Plugin {
	var plugins []*Plugin
	f, err := os.Open(path)
	if err != nil {
		log.Warn(err)
		return nil
	}
	defer f.Close()
	dirs, err := f.Readdirnames(-1)
	if err != nil {
		log.Warn(err)
		return nil
	}
	for _, dir := range dirs {
		dir2 := pt.Join(path, dir)
		f2, err := os.Open(dir2)
		if err != nil {
			log.Warn(err)
			continue
		}
		defer f2.Close()
		fi, err := f2.Readdir(-1)
		if err != nil {
			continue
		}
		for _, f := range fi {
			fn := f.Name()
			if strings.HasSuffix(fn, suffix) {
				plugins = append(plugins, NewPlugin(dir2, suffix))
				break
			}
		}
	}
	return plugins
}
