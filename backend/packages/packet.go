// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"encoding/json"
	"github.com/limetext/lime/backend/loaders"
	"github.com/limetext/lime/backend/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type (
	// Packets are small Packages containing 1 file.
	// Individual settings, keymaps, snippets, etc. are Packets
	Packet struct {
		path string
		// the Packet content will be Unmarshal to this variable
		// so on reload we know where we should unmarshal it again
		marshalTo json.Unmarshaler
	}

	// Useful for managing packets for plugins
	// and loading user packets for editor
	Packets []*Packet
)

// Initializes new packet with specific path
func NewPacket(path string, marshal json.Unmarshaler) *Packet {
	return &Packet{path, marshal}
}

func (p *Packet) Name() string {
	return p.path
}

// Returns packet file data if any error occurred
// on reading file we will return nil
func (p *Packet) Get() interface{} {
	e := []byte(`{}`)
	if p.group() == "keymap" {
		e = []byte(`[]`)
	}

	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		log.Finest("%s doesn't exist yet", p.path)
		return e
	}

	d, err := ioutil.ReadFile(p.path)
	if err != nil {
		log.Error("Couldn't read file: %s", err)
		return e
	}
	return d
}

// Forces editor to load the packet again
func (p *Packet) FileChanged() {
	p.Load()
}

func (p *Packet) Load() error {
	return loaders.LoadJSON(p.Get().([]byte), p)
}

func (p *Packet) MarshalTo() json.Unmarshaler {
	return p.marshalTo
}

func (p *Packet) UnmarshalJSON(data []byte) error {
	return p.marshalTo.UnmarshalJSON(data)
}

// Returns packet type(settings, commands, keymaps, ...)
func (p *Packet) group() string {
	for _, key := range types {
		if strings.Contains(filepath.Ext(p.Name()), key) {
			return key
		}
	}
	return ""
}

// Returns Packets with specific type
func (p Packets) Filter(key string) Packets {
	var pckts Packets
	for _, pckt := range p {
		if strings.Contains(filepath.Ext(pckt.Name()), key) {
			pckts = append(pckts, pckt)
		}
	}
	return pckts
}
