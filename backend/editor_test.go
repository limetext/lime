// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"io/ioutil"
	"os"
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

func TestConfigLoading(t *testing.T) {
	editor := GetEditor()
	editor.loadSetting(NewPacket("testdata/Default.sublime-settings"))

	if editor.Settings().Has("tab_size") != true {
		t.Error("Expected editor settings to have tab_size")
	}

	tab_size := editor.Settings().Get("tab_size").(float64)
	if tab_size != 4 {
		t.Errorf("Expected tab_size to equal 4, got: %v", tab_size)
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
	// TODO: this should be uncomment after adding proper
	// settings hiererchy
	return
	var path string = "testdata/Default.sublime-settings"
	editor := GetEditor()
	editor.loadSetting(NewPacket(path))

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
