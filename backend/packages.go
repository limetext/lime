// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"encoding/json"
	"io/ioutil"
	"os"
	pt "path"
	"path/filepath"
	"strings"
)

type (
	Package interface {
		// Returns the path of the package
		Name() string

		// Depending on the implemented package
		// returns useful data for python plugin is
		// python files for setting is file content
		Get() interface{}

		// Reloads package data
		Reload()
	}

	// Plugin is a Package type containing some files
	// with specific suffix that could be interpreted by
	// lime text api(currently python) and some
	// settings, snippets, commands and etc as packets
	Plugin struct {
		path    string
		suffix  string
		files   []os.FileInfo
		packets pckts
	}

	// Packets are small packages containing 1 file.
	// Individual settings, keymaps, snippets and etc
	// are Packet
	packet struct {
		path string
		// the packet content will be Unmarshal to this variable
		// so on reload we know where we should unmarshal it again
		marshalTo json.Unmarshaler
	}

	// Useful for managing packets for plugins
	// and loading user packets for editor
	pckts []*packet
)

// This is useful when we are loading new plugin or
// scanning for user settings, snippets and etc we
// will add files which their suffix contains one of
// these keywords
var types = []string{"settings", "keymap"}

// Initializes a new plugin whith loading all of the
// settings, keymaps and etc. Suffix variable show's
// which file types we need for plugin for example if
// the plugin is written in python the suffix should
// be ".py". We will use this function at initialization
// to add user plugins and on new_plugin command
func NewPlugin(path string, suffix string) *Plugin {
	var p *Plugin = &Plugin{path, suffix, nil, nil}
	p.Reload()
	return p
}

// Returns slice of files with plugin suffix
// loaded at initialization
func (p *Plugin) Get() interface{} {
	return p.files
}

func (p *Plugin) Name() string {
	return p.path
}

// When the plugin is initialized we won't
// load plguin packets until we are asked to
// so here we will load all plugin packets
func (p *Plugin) LoadPackets() {
	for _, pckt := range p.packets {
		pckt.Reload()
	}
}

// On plugin reload we will scan for plugin files
// and packets in plugin path
func (p *Plugin) Reload() {
	var (
		files []os.FileInfo
		pckts []*packet
	)
	log4go.Info("Reloading plugin %s", p.Name())
	f, err := os.Open(p.path)
	if err != nil {
		log4go.Error("Couldn't open dir: %s", err)
		return
	}
	defer f.Close()
	fi, err := f.Readdir(-1)
	if err != nil {
		log4go.Error("Couldn't read dir: %s", err)
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
			var pckt *packet
			if t == "keymap" {
				pckt = NewPacket(pt.Join(p.path, f.Name()), new(KeyBindings))
			} else {
				// We don't have any settings hierarchy for plugins at this moment
				pckt = NewPacket(pt.Join(p.path, f.Name()), nil)
			}
			pckts = append(pckts, pckt)
		}
	}
	p.files = files
	p.packets = pckts
}

// Initializes new packet with specific path
func NewPacket(path string, marshal json.Unmarshaler) *packet {
	return &packet{path, marshal}
}

// Returns packet file data if any error occurred
// on reading file we will return nil
func (p *packet) Get() interface{} {
	d, err := ioutil.ReadFile(p.path)
	if err != nil {
		log4go.Error("Couldn't read file: %s", err)
		return []byte{}
	}
	return d
}

func (p *packet) Name() string {
	return p.path
}

// Forces editor to load the packet again
func (p *packet) Reload() {
	ed := GetEditor()
	if p.group() == "settings" {
		ed.loadSetting(p)
	} else if p.group() == "keymap" {
		ed.loadKeyBinding(p)
	}
}

// Returns packet type(settings, commands, keymaps, ...)
func (p *packet) group() string {
	for _, key := range types {
		if strings.Contains(filepath.Ext(p.Name()), key) {
			return key
		}
	}
	return ""
}

func (p *packet) UnmarshalJSON(data []byte) error {
	return p.marshalTo.UnmarshalJSON(data)
}

// Returns packets with specific type
func (p pckts) filter(key string) []*packet {
	var pckts []*packet
	for _, pckt := range p {
		if strings.Contains(filepath.Ext(pckt.Name()), key) {
			pckts = append(pckts, pckt)
		}
	}
	return pckts
}

// Scaning path for finding plugins that contain files
// whith specific suffix
func ScanPlugins(path string, suffix string) []*Plugin {
	var plugins []*Plugin
	f, err := os.Open(path)
	if err != nil {
		log4go.Warn(err)
		return nil
	}
	defer f.Close()
	dirs, err := f.Readdirnames(-1)
	if err != nil {
		log4go.Warn(err)
		return nil
	}
	for _, dir := range dirs {
		dir2 := pt.Join(path, dir)
		f2, err := os.Open(dir2)
		if err != nil {
			log4go.Warn(err)
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

// Initialize scan for loading user and limetext defaults
// except settings because for settings we have a hierarchy
// i.e commands, snippets etc
func scanPackets(path string) []*packet {
	var packets []*packet
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log4go.Error("Error on walking: %s", err)
			return err
		}
		s := filepath.Ext(info.Name())
		for _, t := range types {
			if t != "settings" && strings.Contains(s, t) {
				packets = append(packets, NewPacket(path, new(KeyBindings)))
			}
		}
		return nil
	}
	if err := filepath.Walk(path, walkFn); err != nil {
		log4go.Error("Can't walk: %s", err)
	}
	return packets
}

// All user individual settings, snippets etc
// will be in here for later loading by editor
var packets pckts

// Loading the default packets
func init() {
	pcts := scanPackets(LIME_DEFAULTS_PATH)
	pcts = append(pcts, scanPackets(LIME_USER_PACKETS_PATH)...)
	for _, p := range pcts {
		packets = append(packets, p)
	}
}
