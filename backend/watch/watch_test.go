// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

func equal(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for _, a := range a1 {
		if !exist(a2, a) {
			return false
		}
	}
	return true
}

func newWatcher(t *testing.T) *Watcher {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("Couldn't create watcher: %s", err)
	}
	return watcher
}

func watch(t *testing.T, watcher *Watcher, name string, cb interface{}) {
	if err := watcher.Watch(name, cb); err != nil {
		t.Fatalf("Couldn' Watch %s : %s", name, err)
	}
}

func unwatch(t *testing.T, watcher *Watcher, name string, cb interface{}) {
	if err := watcher.UnWatch(name, cb); err != nil {
		t.Fatalf("Couldn' UnWatch %s : %s", name, err)
	}
}

func TestNewWatcher(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	if len(watcher.dirs) != 0 {
		t.Errorf("Expected len(dirs) of new watcher %d, but got %d", 0, len(watcher.dirs))
	}
	if len(watcher.watchers) != 0 {
		t.Errorf("Expected len(watchers) of new watcher %d, but got %d", 0, len(watcher.watchers))
	}
}

type dummy struct {
	text string
	name string
	lock sync.Mutex
}

func (d *dummy) FileChanged(name string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.text = "Changed"
}

func (d *dummy) FileCreated(name string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.text = "Created"
}

func (d *dummy) FileRemoved(name string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.text = "Removed"
}

func (d *dummy) FileRenamed(name string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.text = "Renamed"
}

func (d *dummy) Text() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.text
}

func TestWatch(t *testing.T) {
	tests := []struct {
		paths       []string
		expWatched  []string
		expWatchers []string
	}{
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
		},
		{
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt", "testdata"},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
	}
	d := new(dummy)
	for i, test := range tests {
		watcher := newWatcher(t)
		for _, name := range test.paths {
			watch(t, watcher, name, d)
		}
		if len(watcher.watched) != len(test.expWatched) {
			t.Errorf("Test %d: Expected watched %v keys equal to %v", i, watcher.watched, test.expWatched)
		}
		for _, p := range test.expWatched {
			if _, exist := watcher.watched[p]; !exist {
				t.Errorf("Test %d: Expected %s exist in watched", i, p)
			}
		}
		if !equal(test.expWatchers, watcher.watchers) {
			t.Errorf("Test %d: Expected watchers %v, but got %v", i, test.expWatchers, watcher.watchers)
		}
		watcher.wchr.Close()
	}
}

func Testwatch(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	if err := watcher.watch("testdata/dummy.txt"); err != nil {
		t.Fatalf("Couldn't watch %s", "testdata/dummy.txt")
	}
	if err := watcher.watch("testdata/test.txt"); err != nil {
		t.Fatalf("Couldn't watch %s", "testdata/test.txt")
	}
	if !equal(watcher.watchers, []string{"testdata/dummy.txt", "testdata/test.txt"}) {
		t.Errorf("Expected watchers %v, but got %v", []string{"testdata/dummy.txt", "testdata/test.txt"}, watcher.watchers)
	}
}

func TestAdd(t *testing.T) {
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := new(dummy)
	watcher.add("test", d)
	if cb := watcher.watched["test"][0]; cb != d {
		t.Errorf("Expected watcher['test'][0] callback equal to %d, but got %d", d, cb)
	}
}

func TestFlushDir(t *testing.T) {
	name := "testdata/dummy.txt"
	dir := "testdata"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := new(dummy)
	watch(t, watcher, name, d)
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers equal to %v, but got %v", []string{name}, watcher.watchers)
	}
	watcher.flushDir(dir)
	if !equal(watcher.dirs, []string{dir}) {
		t.Errorf("Expected dirs equal to %v, but got %v", []string{dir}, watcher.dirs)
	}
	if !equal(watcher.watchers, []string{}) {
		t.Errorf("Expected watchers equal to %v, but got %v", []string{}, watcher.watchers)
	}
}

func TestUnWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := new(dummy)
	watch(t, watcher, name, d)
	unwatch(t, watcher, name, d)
	if len(watcher.watched) != 0 {
		t.Errorf("Expected watcheds be empty, but got %v", watcher.watched)
	}
}

func TestUnWatchAll(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if l := len(watcher.watched[name]); l != 2 {
		t.Errorf("Expected len of watched['%s'] be %d, but got %d", name, 2, l)
	}
	unwatch(t, watcher, name, nil)
	if _, exist := watcher.watched[name]; exist {
		t.Errorf("Expected all %s watched be removed", name)
	}
	if !equal(watcher.watchers, []string{}) {
		t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	}
}

func TestUnWatchDirectory(t *testing.T) {
	name := "testdata/dummy.txt"
	dir := "testdata"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := new(dummy)
	watch(t, watcher, name, d)
	watch(t, watcher, dir, nil)
	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Fatalf("Expected watchers be equal to %s, but got %s", []string{"testdata"}, watcher.watchers)
	}
	unwatch(t, watcher, dir, nil)
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", []string{name}, watcher.watchers)
	}
}

func TestUnWatchOneOfSubscribers(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if len(watcher.watched[name]) != 2 {
		t.Fatalf("Expected watched[%s] length be %d, but got %d", name, 2, len(watcher.watched[name]))
	}
	unwatch(t, watcher, name, d1)
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", []string{name}, watcher.watchers)
	}
	if len(watcher.watched[name]) != 1 {
		t.Errorf("Expected watched[%s] length be %d, but got %d", name, 1, len(watcher.watched[name]))
	}
}

func TestunWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d1 := new(dummy)
	d2 := new(dummy)
	watch(t, watcher, name, d1)
	watch(t, watcher, name, d2)
	if err := watcher.unWatch(name); err != nil {
		t.Fatalf("Couldn't unWatch %s: %s", name, err)
	}
	if _, exist := watcher.watched[name]; exist {
		t.Errorf("Expected all %s watched be removed", name)
	}
	if !equal(watcher.watchers, []string{}) {
		t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	}
}

func TestRemoveWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := new(dummy)
	watch(t, watcher, name, d)
	watcher.removeWatch(name)
	if !equal(watcher.watchers, []string{}) {
		t.Errorf("Expected watchers be empty but got %v", watcher.watchers)
	}
}

func TestRemoveDir(t *testing.T) {
	name := "testdata/dummy.txt"
	dir := "testdata"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := new(dummy)
	watch(t, watcher, dir, d)
	watch(t, watcher, name, d)
	if !equal(watcher.watchers, []string{dir}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", []string{dir}, watcher.watchers)
	}
	if !equal(watcher.dirs, []string{dir}) {
		t.Errorf("Expected dirs be equal to %s, but got %s", []string{dir}, watcher.dirs)
	}
	watcher.removeDir(dir)
	if !equal(watcher.dirs, []string{}) {
		t.Errorf("Expected dirs be empty but got %v", watcher.dirs)
	}
	if !equal(watcher.watchers, []string{dir, name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", []string{name}, watcher.watchers)
	}
}

func TestObserve(t *testing.T) {
	name := "testdata/test.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := &dummy{name: name}
	watch(t, watcher, name, d)
	go watcher.Observe()

	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if d.Text() != "Changed" {
		t.Errorf("Expected dumView Text %s, but got %s", "Changed", d.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestObserveDirectory(t *testing.T) {
	dir := "testdata"
	name := "testdata/test.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := &dummy{name: dir}
	watch(t, watcher, dir, d)
	go watcher.Observe()

	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watcher.watchers)
	}
	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if d.Text() != "Changed" {
		t.Errorf("Expected dumView Text %s, but got %s", "Changed", d.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestCreateEvent(t *testing.T) {
	name := "testdata/new.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := &dummy{name: name}
	watch(t, watcher, name, d)
	go watcher.Observe()

	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watcher.watchers)
	}
	if err := ioutil.WriteFile(name, []byte(""), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if d.Text() != "Created" {
		t.Errorf("Expected dumView Text %s, but got %s", "Created", d.Text())
	}
	os.Remove(name)
}

func TestDeleteEvent(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := &dummy{name: name}
	watch(t, watcher, name, d)
	go watcher.Observe()

	if err := os.Remove(name); err != nil {
		t.Fatalf("Couldn't remove file %s: %s", name, err)
	}
	time.Sleep(time.Millisecond * 50)
	watcher.lock.Lock()
	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watcher.watchers)
	}
	watcher.lock.Unlock()
	if d.Text() != "Removed" {
		t.Errorf("Expected dumView Text %s, but got %s", "Removed", d.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestRenameEvent(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	defer watcher.wchr.Close()
	d := &dummy{name: name}
	watch(t, watcher, name, d)
	go watcher.Observe()

	os.Rename(name, "testdata/rename.txt")
	time.Sleep(time.Millisecond * 50)
	if d.Text() != "Renamed" {
		t.Errorf("Expected dumView Text %s, but got %s", "Renamed", d.Text())
	}
	os.Rename("testdata/rename.txt", name)
}
