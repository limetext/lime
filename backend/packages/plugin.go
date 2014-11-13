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
	"path/filepath"
	"strings"
)

type (
	// Plugin is a Package type containing some files
	// with specific suffix that could be interpreted by
	// lime text api(currently python) and some
	// settings, snippets, commands and etc as packets
	Plugin struct {
		text.HasSettings
		path    string
		suffix  string
		files   []os.FileInfo
		packets Packets
	}
)

// Initializes a new plugin whith loading all of the
// settings, keymaps and etc. Suffix variable show's
// which file types we need for plugin for example if
// the plugin is written in python the suffix should
// be ".py". We will use this function at initialization
// to add user plugins and on new_plugin command
func NewPlugin(path string, suffix string) *Plugin {
	var p *Plugin = &Plugin{path: path, suffix: suffix}
	p.Reload()
	return p
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
	var (
		files []os.FileInfo
		pckts Packets
	)
	log.Info("Reloading plugin %s", p.Name())
	f, err := os.Open(p.path)
	if err != nil {
		log.Error("Couldn't open dir: %s", err)
		return
	}
	defer f.Close()
	fi, err := f.Readdir(-1)
	if err != nil {
		log.Error("Couldn't read dir: %s", err)
		return
	}
	for _, f := range fi {
		if p.suffix != "" && strings.HasSuffix(f.Name(), p.suffix) {
			files = append(files, f)
			continue
		}
		s := filepath.Ext(f.Name())
		for _, t := range types {
			if !strings.Contains(s, t) {
				continue
			}
			var pckt *Packet
			if t == "keymap" {
				pckt = NewPacket(pt.Join(p.path, f.Name()), new(keys.KeyBindings))
			} else {
				// We don't have any settings hierarchy for plugins at this moment
				pckt = NewPacket(pt.Join(p.path, f.Name()), p.Settings())
			}
			pckts = append(pckts, pckt)
		}
	}
	p.files = files
	p.packets = pckts
}

// When the plugin is initialized we won't
// load plugin packets until we are asked to
// so here we will load all plugin packets
func (p *Plugin) LoadPackets() {
	for _, pckt := range p.packets {
		pckt.Load()
	}
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
