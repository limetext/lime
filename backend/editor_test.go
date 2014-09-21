// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/packages"
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
	editor := GetEditor()
	if editor == nil {
		t.Error("Expected an editor, but got nil")
	}
}

func TestLoadKeyBinding(t *testing.T) {
	editor := GetEditor()
	editor.loadKeyBinding(packages.NewPacket("testdata/Default.sublime-keymap", new(keys.KeyBindings)))

	kb := editor.keyBindings.Filter(keys.KeyPress{Key: 'i'})
	if kb.Len() == 69 {
		t.Errorf("Expected to have %d keys in the filter, but it had %d", 69, kb.Len())
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
	editor.loadSetting(packages.NewPacket("testdata/Default.sublime-settings", editor.Settings()))

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
	LIME_DEFAULTS_PATH = path.Join("..", "packages", "Default")

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

	editor := GetEditor()
	w := editor.NewWindow()
	defer w.Close()

	for i, test := range tests {
		v := w.OpenFile(testfile, 0)

		if err := v.SaveAs(test.as); err != nil {
			t.Fatalf("Test %d: Can't save to `%s`: %s", i, test.as, err)
		}

		if v.IsDirty() {
			t.Errorf("Test %d: Expected the view to be clean, but it wasn't", i)
		}

		if _, exist := editor.watchedFiles[test.as]; !exist {
			t.Errorf("Test %d: Should watch %s file", i, test.as)
		}

		v.Close()

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
	editor.loadSetting(packages.NewPacket(path, editor.Settings()))

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
	editor := GetEditor()
	l := len(editor.Windows())

	w := editor.NewWindow()
	defer w.Close()

	if len(editor.Windows()) != l+1 {
		t.Errorf("Expected 1 window, but got %d", len(editor.Windows()))
	}
}

func TestRemoveWindow(t *testing.T) {
	editor := GetEditor()
	l := len(editor.Windows())

	w0 := editor.NewWindow()
	defer w0.Close()

	editor.remove(w0)

	if len(editor.Windows()) != l {
		t.Errorf("Expected the window to be removed, but %d still remain", len(editor.Windows()))
	}

	w1 := editor.NewWindow()
	defer w1.Close()

	w2 := editor.NewWindow()
	defer w2.Close()

	editor.remove(w1)

	if len(editor.Windows()) != l+1 {
		t.Errorf("Expected the window to be removed, but %d still remain", len(editor.Windows()))
	}
}

func TestSetActiveWindow(t *testing.T) {
	editor := GetEditor()

	w1 := editor.NewWindow()
	defer w1.Close()

	w2 := editor.NewWindow()
	defer w2.Close()

	if editor.ActiveWindow() != w2 {
		t.Error("Expected the newest window to be active, but it wasn't")
	}

	editor.SetActiveWindow(w1)

	if editor.ActiveWindow() != w1 {
		t.Error("Expected the first window to be active, but it wasn't")
	}
}

func TestSetFrontend(t *testing.T) {
	f := DummyFrontend{}

	editor := GetEditor()
	editor.SetFrontend(&f)

	if editor.Frontend() != &f {
		t.Errorf("Expected a DummyFrontend to be set, but got %T", editor.Frontend())
	}
}

func TestClipboard(t *testing.T) {
	editor := GetEditor()
	s := "test"

	editor.SetClipboard(s)

	if editor.GetClipboard() != s {
		t.Errorf("Expected %s to be on the clipboard, but got %s", s, editor.GetClipboard())
	}
}

func TestHandleInput(t *testing.T) {
	editor := GetEditor()
	kp := keys.KeyPress{Key: 'i'}

	editor.HandleInput(kp)

	if ki := <-editor.keyInput; ki != kp {
		t.Errorf("Expected %s to be on the input buffer, but got %s", kp, ki)
	}
}
