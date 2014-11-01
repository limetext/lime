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

type dummy struct {
	name string
}

func (d *dummy) Name() string {
	return d.name
}

func (d *dummy) Reload() {}

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
		// 2 path refer same file or dir but different(e.g abs path and relative path)
	}
	for i, test := range tests {
		watcher := NewWatcher()
		for _, name := range test.paths {
			watcher.Watch(&dummy{name})
		}
		if len(watcher.watched) != len(test.expWatched) {
			t.Errorf("Test %d: Expected len of watched %d, but got %d", i, len(test.expWatched), len(watcher.watched))
		}
		for _, p := range test.expWatched {
			if _, exist := watcher.watched[p]; !exist {
				t.Errorf("Test %d: Expected %s exist in watched", i, p)
			}
		}
		if !equal(test.expWatchers, watcher.watchers) {
			t.Errorf("Test %d: Expected watchers %v, but got %v", i, test.expWatchers, watcher.watchers)
		}
	}
}

func TestUnWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	d := &dummy{name}
	watcher := NewWatcher()
	watcher.Watch(d)
	watcher.UnWatch(d)
	if len(watcher.watched) != 0 {
		t.Errorf("Expected watcheds be empty, but got %s", watcher.watched)
	}
}

func TestUnWatchDirectory(t *testing.T) {
	name := "testdata/dummy.txt"
	dir := "testdata"
	d := &dummy{name}
	d1 := NewWatchedDir(dir)
	watcher := NewWatcher()
	watcher.Watch(d)
	watcher.Watch(d1)
	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Fatalf("Expected watchers be equal to %s, but got %s", watcher.watchers, []string{"testdata"})
	}
	watcher.UnWatch(d1)
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", watcher.watchers, []string{name})
	}
}

func TestUnWatchOneOfSubscribers(t *testing.T) {
	name := "testdata/dummy.txt"
	d1 := &dummy{name}
	d2 := &dummy{name}
	watcher := NewWatcher()
	watcher.Watch(d1)
	watcher.Watch(d2)
	if len(watcher.watched[name]) != 2 {
		t.Fatalf("Expected watched[%s] length be %d, but got %d", name, 2, len(watcher.watched[name]))
	}
	watcher.UnWatch(d1)
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", watcher.watchers, []string{name})
	}
	if len(watcher.watched[name]) != 1 {
		t.Errorf("Expected watched[%s] length be %d, but got %d", name, 1, len(watcher.watched[name]))
	}
}

type dumView struct {
	text string
	name string
	lock sync.Mutex
}

func (d *dumView) Reload() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.text = "Reloaded"
}

func (d *dumView) Name() string {
	return d.name
}

func (d *dumView) Text() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.text
}

func TestObserve(t *testing.T) {
	name := "testdata/test.txt"
	watcher := NewWatcher()
	v := &dumView{name: name}
	watcher.Watch(v)
	go watcher.Observe()

	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestObserveDirectory(t *testing.T) {
	dir := "testdata"
	name := "testdata/test.txt"
	watcher := NewWatcher()
	v := &dumView{name: name}
	watcher.Watch(v)
	watcher.Watch(NewWatchedDir(dir))
	go watcher.Observe()

	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watcher.watchers)
	}
	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestCreateEvent(t *testing.T) {
	name := "testdata/new.txt"
	watcher := NewWatcher()
	v := &dumView{name: name}
	watcher.Watch(v)
	go watcher.Observe()

	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{name}, watcher.watchers)
	}
	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	os.Remove(name)
}

func TestDeleteEvent(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := NewWatcher()
	v := &dumView{name: name}
	watcher.Watch(v)
	go watcher.Observe()

	os.Remove(name)
	time.Sleep(time.Millisecond * 50)
	watcher.lock.Lock()
	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watcher.watchers)
	}
	watcher.lock.Unlock()
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestRenameEvent(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := NewWatcher()
	v := &dumView{name: name}
	watcher.Watch(v)
	go watcher.Observe()

	os.Rename(name, "testdata/rename.txt")
	time.Sleep(time.Millisecond * 50)
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	os.Rename("testdata/rename.txt", name)
}

func TestExist(t *testing.T) {
	test := struct {
		array []string
		elms  []string
		exps  []bool
	}{
		[]string{"a", "b", "c", "d"},
		[]string{"a", "t", "A"},
		[]bool{true, false, false},
	}
	for i, exp := range test.exps {
		if exist(test.array, test.elms[i]) != exp {
			t.Errorf("Expected in %v exist result of element %s be %v, but got %v", test.array, test.elms[i], exp, exist(test.array, test.elms[i]))
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		slice  []string
		remove string
		exp    []string
	}{
		{
			[]string{"a", "b", "c"},
			"a",
			[]string{"c", "b"},
		},
		{
			[]string{"a", "b", "c"},
			"k",
			[]string{"a", "b", "c"},
		},
	}
	for i, test := range tests {
		if exp := remove(test.slice, test.remove); !equal(exp, test.exp) {
			t.Errorf("Test %d: Expected %v be equal to %v", i, exp, test.exp)
		}
	}
}
