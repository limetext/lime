// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package packages

import (
	"io/ioutil"
	"testing"
)

func TestPacket(t *testing.T) {
	tests := []struct {
		path string
		data string
	}{
		{
			"testdata/Vintage/Vintageous.sublime-settings",
			"Testing packages",
		},
	}
	for i, test := range tests {
		s := NewPacket(test.path, nil)
		d, err := ioutil.ReadFile(test.path)
		if err != nil {
			t.Fatalf("Test %d: Can't read file: %s", i, err)
		}
		if d1 := s.Get().([]byte); string(d) != string(d1) {
			t.Errorf("Test %d: Expected %v but, got %s", i, d, string(d1))
		}
		if err := ioutil.WriteFile(test.path, []byte(test.data), 0644); err != nil {
			t.Fatalf("Test %d: Can't write file: %s", i, err)
		}
		s.Reload()
		if d1 := s.Get().([]byte); test.data != string(d1) {
			t.Errorf("Test %d: Expected %s but, got %s", i, test.data, string(d1))
		}
		if err := ioutil.WriteFile(test.path, d, 0644); err != nil {
			t.Fatalf("Test %d: Can't write back file: %s", i, err)
		}
	}
}

func TestPackets(t *testing.T) {
	test := struct {
		pckts  []string
		expect map[string][]string
	}{
		[]string{
			"testdata/Default.sublime-settings",
			"testdata/Vintage/Default.sublime-keymap",
			"testdata/Vintage/Vintageous.sublime-settings",
		},
		map[string][]string{
			"setting": {
				"testdata/Default.sublime-settings",
				"testdata/Vintage/Vintageous.sublime-settings",
			},
			"keymap": {
				"testdata/Vintage/Default.sublime-keymap",
			},
		},
	}

	packets := Packets{}
	for _, p := range test.pckts {
		packets = append(packets, NewPacket(p, nil))
	}
	for key, ns := range test.expect {
		ps := packets.Filter(key)
		for _, p := range ns {
			found := false
			for _, p1 := range ps {
				if p == p1.Name() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %s in %s plugin packets", p, ps)
			}
		}
	}
}

func TestScanPackets(t *testing.T) {
	tests := []struct {
		path   string
		expect []string
	}{
		{
			"../../packages",
			[]string{
				"../../packages/Default/Default.sublime-keymap",
			},
		},
		{
			"testdata",
			[]string{
				// "testdata/Vintage/Default.sublime-commands", not supporting commands yet
				"testdata/Vintage/Default.sublime-keymap",
				"testdata/Default.sublime-keymap",
			},
		},
	}
	for i, test := range tests {
		packets := ScanPackets(test.path)
		for _, f := range test.expect {
			found := false
			for _, p := range packets {
				if f == p.Name() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Test %d: Expected ScanPackets find %s packet", i, f)
			}
		}
	}
}
