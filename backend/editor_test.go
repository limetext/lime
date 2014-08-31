// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

type DummyWatched struct {
	name string
}

func (d *DummyWatched) Name() string {
	return d.name
}

func (d *DummyWatched) Reload() {
	// noop
}

func TestGetEditor(t *testing.T) {
	e := GetEditor()
	if e == nil {
		t.Error("Expected an editor, but got nil")
	}
}

func TestLoadKeyBinding(t *testing.T) {
	var kb KeyBindings

	editor := GetEditor()
	editor.loadKeyBinding(NewPacket("testdata/Default.sublime-keymap", new(KeyBindings)))

	editor.keyBindings.filter(69, &kb)
	if kb.Len() == 69 {
		t.Errorf("Expected editor to have key %d bound, but it didn't", 69)
	}
}

func TestLoadKeyBindings(t *testing.T) {
	editor := GetEditor()
	editor.loadKeyBindings()

	editor.keyBindings.Len()
	if editor.keyBindings.Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}
}

func TestLoadSetting(t *testing.T) {
	editor := GetEditor()
	editor.loadSetting(NewPacket("testdata/Default.sublime-settings", editor.Settings()))

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}

	tab_size := editor.Settings().Get("tab_size").(float64)
	if tab_size != 4 {
		t.Errorf("Expected tab_size to equal 4, got: %v", tab_size)
	}
}

func TestLoadSettings(t *testing.T) {
	LIME_USER_PACKAGES_PATH = path.Join("..", "3rdparty", "bundles")
	LIME_USER_PACKETS_PATH = path.Join("..", "3rdparty", "bundles", "User")
	LIME_DEFAULTS_PATH = path.Join("packages", "Default")

	editor := GetEditor()
	editor.loadSettings()

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}

	plat := editor.Settings().Parent()
	switch editor.Platform() {
	case "windows":
		if plat.Settings().Get("font_face", "") != "Consolas" {
			t.Errorf("Expected windows font_face be Consolas, but is %s", plat.Settings().Get("font_face", ""))
		}
	case "darwin":
		if plat.Settings().Get("font_face", "") != "Menlo Regular" {
			t.Errorf("Expected OSX font_face be Menlo Regular, but is %s", plat.Settings().Get("font_face", ""))
		}
	default:
		if plat.Settings().Get("font_face", "") != "Monospace" {
			t.Errorf("Expected Linux font_face be Monospace, but is %s", plat.Settings().Get("font_face", ""))
		}
	}
}

func TestInit(t *testing.T) {
	editor := GetEditor()
	editor.Init()

	editor.keyBindings.Len()
	if editor.keyBindings.Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}
}

func TestWatch(t *testing.T) {
	editor := GetEditor()
	observedFile := &DummyWatched{"editor_test.go"}
	editor.Watch(observedFile)

	if editor.watchedFiles["editor_test.go"] != observedFile {
		t.Fatal("Expected editor to watch the specified file")
	}
}

func TestWatchOnSaveAs(t *testing.T) {
	var testfile string = "testdata/Default.sublime-settings"
	tests := []struct {
		as string
	}{
		{
			"User.sublime-settings",
		},
		{
			"testdata/User.sublime-settings",
		},
	}
	ed := GetEditor()
	w := ed.NewWindow()
	for i, test := range tests {
		v := w.OpenFile(testfile, 0)
		if err := v.SaveAs(test.as); err != nil {
			t.Fatalf("Test %d: Can't save to `%s`: %s", i, test.as, err)
		}
		if _, exist := ed.watchedFiles[test.as]; !exist {
			t.Errorf("Test %d: Should watch %s file", i, test.as)
		}
		if err := os.Remove(test.as); err != nil {
			t.Errorf("Test %d: Couldn't remove test file %s", i, test.as)
		}
	}
}

func TestWatchingSettings(t *testing.T) {
	// TODO: This won't pass until the settings hierarchy is set up properly.
	return

	var path string = "testdata/Default.sublime-settings"
	editor := GetEditor()
	editor.loadSetting(NewPacket(path, editor.Settings()))

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal("Error in reading the default settings")
	}

	data := []byte("{\n\t\"tab_size\": 8\n}")
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatal("Error in writing to setting")
	}
	time.Sleep(time.Millisecond * 10)
	if tab_size := editor.Settings().Get("tab_size").(float64); tab_size != 8 {
		t.Errorf("Expected tab_size equal to 8, but got %v", tab_size)
	}

	err = ioutil.WriteFile(path, buf, 0644)
	if err != nil {
		t.Fatal("Error in writing the default back to setting")
	}
}

func TestNewWindow(t *testing.T) {
	ed := GetEditor()
	l := len(ed.Windows())
	_ = ed.NewWindow()

	if len(ed.Windows()) != l+1 {
		t.Errorf("Expected 1 window, but got %d", len(ed.Windows()))
	}
}

func TestRemoveWindow(t *testing.T) {
	ed := GetEditor()
	l := len(ed.Windows())

	w := ed.NewWindow()
	ed.remove(w)

	if len(ed.Windows()) != l {
		t.Errorf("Expected the window to be removed, but %d still remain", len(ed.Windows()))
	}
}

func TestSetActiveWindow(t *testing.T) {
	ed := GetEditor()

	w1 := ed.NewWindow()
	w2 := ed.NewWindow()

	if ed.ActiveWindow() != w2 {
		t.Error("Expected the newest window to be active, but it wasn't")
	}

	ed.SetActiveWindow(w1)

	if ed.ActiveWindow() != w1 {
		t.Error("Expected the first window to be active, but it wasn't")
	}
}

func TestSetFrontend(t *testing.T) {
	f := DummyFrontend{}

	ed := GetEditor()
	ed.SetFrontend(&f)

	if ed.Frontend() != &f {
		t.Errorf("Expected a DummyFrontend to be set, but got %T", ed.Frontend())
	}
}

func TestClipboard(t *testing.T) {
	ed := GetEditor()
	s := "test"

	ed.SetClipboard(s)

	if ed.GetClipboard() != s {
		t.Errorf("Expected %s to be on the clipboard, but got %s", s, ed.GetClipboard())
	}
}

func TestHandleInput(t *testing.T) {
	ed := GetEditor()
	kp := KeyPress{Key: 'i'}

	ed.HandleInput(kp)

	if ki := <-ed.keyInput; ki != kp {
		t.Errorf("Expected %s to be on the input buffer, but got %s", kp, ki)
	}
}
