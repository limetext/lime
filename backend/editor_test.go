// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"io/ioutil"
	"testing"
	"time"
)

type DummyWatchedFile struct {
	name string
}

func (d *DummyWatchedFile) Name() string {
	return d.name
}

func (d *DummyWatchedFile) Reload() {
	// noop
}

func TestConfigLoading(t *testing.T) {
	editor := GetEditor()
	editor.loadSetting("testdata/Default.sublime-settings")

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
	observedFile := &DummyWatchedFile{"editor_test.go"}
	editor.Watch(observedFile)

	if editor.watchedFiles["editor_test.go"] != observedFile {
		t.Fatal("Expected editor to watch the specified file")
	}
}

func TestWatchingSettings(t *testing.T) {
	var path string = "testdata/Default.sublime-settings"
	editor := GetEditor()
	editor.loadSetting(path)

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
