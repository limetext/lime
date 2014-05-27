package backend

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPlugin(t *testing.T) {
	tests := []struct {
		path   string
		suffix string
		files  []string
		pkcts  []string
	}{
		{
			"testdata/Vintageous",
			".py",
			[]string{"action_cmds.py", "state.py", "transformers.py"},
			[]string{
				"testdata/Vintageous/Vintageous.sublime-settings",
				"testdata/Vintageous/Default.sublime-keymap",
			},
		},
	}
	for i, test := range tests {
		p := NewPlugin(test.path, test.suffix)
		if p.Name() != test.path {
			t.Errorf("Test %d: Expected plugin name %s but, got %s", i, test.path, p.Name())
		}
		for _, f := range test.files {
			found := false
			for _, fi := range p.Get().([]os.FileInfo) {
				if f == fi.Name() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Test %d: Expected to find %s in %s plugin", i, f, p.Name())
			}
		}
		for _, f := range test.pkcts {
			found := false
			for _, pckt := range p.Packets() {
				if f == pckt.Name() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Test %d: Expected to find %s in %s plugin packets", i, f, p.Name())
			}
		}
	}
}

func TestPluginReload(t *testing.T) {
	p := NewPlugin("testdata/Closetag", ".vim")
	if err := ioutil.WriteFile("testdata/Closetag/test.vim", []byte("testing"), 0644); err != nil {
		t.Fatalf("Couldn't write file: %s", err)
	}
	if err := ioutil.WriteFile("testdata/Closetag/test.settings", []byte("testing packets"), 0644); err != nil {
		t.Fatalf("Couldn't write file: %s", err)
	}
	p.Reload()
	fi := p.Get().([]os.FileInfo)
	found := false
	for _, f := range fi {
		if f.Name() == "test.vim" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find test.vim file in %s", p.Name())
	}
	found = false
	for _, p := range p.Packets() {
		if p.Name() == "testdata/Closetag/test.settings" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected to find testdata/Closetag/test.settings file in %s", p.Name())
	}
}

func TestScanPath(t *testing.T) {
	tests := []struct {
		path   string
		suffix string
		expect []string
	}{
		{
			"testdata/",
			".py",
			[]string{
				"testdata/Vintageous",
			},
		},
		{
			"testdata/",
			".vim",
			[]string{
				"testdata/Closetag",
			},
		},
	}
	for i, test := range tests {
		plugins := ScanPath(test.path, test.suffix)
		for _, f := range test.expect {
			found := false
			for _, p := range plugins {
				if f == p.Name() {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Test %d: Expected ScanPath find %s plugin", i, f)
			}
		}
	}
}

func TestPacket(t *testing.T) {
	tests := []struct {
		path string
		data string
	}{
		{
			"testdata/Vintageous/Vintageous.sublime-settings",
			"Testing packages",
		},
	}
	for i, test := range tests {
		s := NewPacket(test.path)
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

func TestPckts(t *testing.T) {
	test := struct {
		pckts  []string
		expect map[string][]string
	}{
		[]string{
			"testdata/Default.sublime-settings",
			"testdata/Vintageous/Default.sublime-keymap",
			"testdata/Vintageous/Vintageous.sublime-settings",
		},
		map[string][]string{
			"setting": []string{
				"testdata/Default.sublime-settings",
				"testdata/Vintageous/Vintageous.sublime-settings",
			},
			"keymap": []string{
				"testdata/Vintageous/Default.sublime-keymap",
			},
		},
	}
	for _, p := range test.pckts {
		add(NewPacket(p))
	}
	for key, ns := range test.expect {
		ps := Packets.Type(key)
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
