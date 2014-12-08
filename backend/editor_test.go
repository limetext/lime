// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/packages"
	"path"
	"testing"
)

func init() {
	LIME_PACKAGES_PATH = path.Join("..", "packages")
	LIME_USER_PACKAGES_PATH = path.Join("..", "packages", "User")
	LIME_DEFAULTS_PATH = path.Join("..", "packages", "Default")
}

func TestGetEditor(t *testing.T) {
	editor := GetEditor()
	if editor == nil {
		t.Error("Expected an editor, but got nil")
	}
}

func TestLoadKeyBinding(t *testing.T) {
	editor := GetEditor()
	pkg := packages.NewPacket("testdata/Default.sublime-keymap", editor.KeyBindings())
	editor.load(pkg)

	kb := editor.KeyBindings().Filter(keys.KeyPress{Key: 'i'})
	if expectedLen := 3; kb.Len() != expectedLen {
		t.Errorf("Expected to have %d keys in the filter, but it had %d", expectedLen, kb.Len())
	}
}

func TestLoadKeyBindings(t *testing.T) {
	editor := GetEditor()
	editor.loadKeyBindings()

	if editor.defaultBindings.KeyBindings().Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}
}

func TestLoadSetting(t *testing.T) {
	editor := GetEditor()
	editor.load(packages.NewPacket("testdata/Default.sublime-settings", editor.Settings()))

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
	}

	tab_size := editor.Settings().Get("tab_size").(float64)
	if tab_size != 4 {
		t.Errorf("Expected tab_size to equal 4, got: %v", tab_size)
	}
}

func TestLoadSettings(t *testing.T) {
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

	if editor.defaultBindings.KeyBindings().Len() <= 0 {
		t.Errorf("Expected editor to have some keys bound, but it didn't")
	}

	if editor.Settings().Parent().Settings().Parent().Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size, but it didn't")
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

	// Put back whatever was already there.
	clip := editor.GetClipboard()
	defer editor.SetClipboard(clip)

	s := "test0"

	editor.SetClipboard(s)

	if editor.GetClipboard() != s {
		t.Errorf("Expected %q to be on the clipboard, but got %q", s, editor.GetClipboard())
	}

	s = "test1"

	editor.SetClipboard(s)

	if editor.GetClipboard() != s {
		t.Errorf("Expected %q to be on the clipboard, but got %q", s, editor.GetClipboard())
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
