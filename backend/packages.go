// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"os"
	"path/filepath"
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

	Pckts []*Packet

	Plugin struct {
		path    string
		suffix  string
		files   []os.FileInfo
		packets Pckts
	}

	// Packets are small packages containing 1 file
	// individual settings, keymaps, snippets and etc
	// are Packet
	Packet struct {
		path string
		data []byte
	}
)

const (
	DEFAULT_SUBLIME_SETTINGS    = "../../backend/packages/Default/Default.sublime-settings"
	DEFAULT_SUBLIME_KEYBINDINGS = "../../backend/packages/Default/Default.sublime-keymap"
	SUBLIME_USER_PACKAGES_PATH  = "../../3rdparty/bundles/"
)

// Valid packet types
// TODO: command, snippet and etc should be here
var types = []string{"setting", "keymap"}

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
	pckts := make([]*Packet, 0)
	for _, f := range fi {
		if suffix != "" && strings.HasSuffix(f.Name(), suffix) {
			files = append(files, f)
		} else {
			s := filepath.Ext(f.Name())
			for _, t := range types {
				if strings.Contains(s, t) {
					pckts = append(pckts, NewPacket(path+string(os.PathSeparator)+f.Name()))
				}
			}
		}
	}
	return &Plugin{path, suffix, files, pckts}
}

func (p *Plugin) Get() interface{} {
	return p.files
}

func (p *Plugin) Name() string {
	return p.path
}

func (p *Plugin) Packets() []*Packet {
	return p.packets
}

func (p *Plugin) Reload() {
	log4go.Info("Reloading plugin %s", p.Name())
	p1 := NewPlugin(p.path, p.suffix)
	p.files = p1.files
	p.packets = p1.packets
	for _, pckt := range p.packets {
		pckt.Reload()
	}
}

func NewPacket(path string) *Packet {
	return &Packet{path, nil}
}

func loadData(path string) []byte {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log4go.Error("Couldn't read file: %s", err)
		return nil
	}
	return d
}

func (p *Packet) Get() interface{} {
	if p.data == nil {
		p.data = loadData(p.path)
	}
	return p.data
}

func (p *Packet) Name() string {
	return p.path
}

func (p *Packet) Reload() {
	log4go.Info("Reloading %s", p.Name())
	p.data = loadData(p.path)
	e := GetEditor()
	e.loadSetting(p)
}

func (p *Pckts) Type(key string) []*Packet {
	pckts := make([]*Packet, 0)
	for _, pckt := range Packets {
		if strings.Contains(filepath.Ext(pckt.Name()), key) {
			pckts = append(pckts, pckt)
		}
	}
	return pckts
}

func add(p *Packet) {
	if !reflect.ValueOf(p).IsNil() {
		Packets = append(Packets, p)
	}
}

// Scaning path for finding plugins that contain files
// whith specific suffix
func ScanPath(path string, suffix string) []*Plugin {
	ps := make([]*Plugin, 0)
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
		if dir != "Vintageous" && dir != "Default" && dir != "plugins" && dir != "Closetag" {
			// TODO obviously
			continue
		}
		dir2 := path + dir
		f2, err := os.Open(dir2)
		if err != nil {
			log4go.Warn(err)
			continue
		}
		defer f2.Close()
		fi, err := f2.Readdir(-1)
		if err != nil {
			log4go.Warn(err)
		}
		for _, f := range fi {
			fn := f.Name()
			if !strings.HasSuffix(fn, ".py") {
				ps = append(ps, NewPlugin(dir2, suffix))
				break
			}
		}
	}
	return ps
}

// We'll store loaded packets on startup here
// plugins specific packets won't be in here
// they should be accessed from the plugin itself
var Packets Pckts

// Loading the default packets
func init() {
	Packets = make([]*Packet, 0)

	add(NewPacket(DEFAULT_SUBLIME_SETTINGS))
	add(NewPacket(DEFAULT_SUBLIME_KEYBINDINGS))
}
