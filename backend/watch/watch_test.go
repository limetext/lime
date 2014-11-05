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

func watch(t *testing.T, watcher *Watcher, name, key string, act func()) {
	if err := watcher.Watch(name, key, act); err != nil {
		t.Fatalf("Couldn' Watch %s : %s", name, err)
	}
}

func unwatch(t *testing.T, watcher *Watcher, name, key string) {
	if err := watcher.UnWatch(name, "test"); err != nil {
		t.Fatalf("Couldn' UnWatch %s : %s", name, err)
	}
}

func dummy() {}

func TestNewWatcher(t *testing.T) {
	watcher := newWatcher(t)
	if len(watcher.dirs) != 0 {
		t.Errorf("Expected len(dirs) of new watcher %d, but got %d", 0, len(watcher.dirs))
	}
	if len(watcher.watchers) != 0 {
		t.Errorf("Expected len(watchers) of new watcher %d, but got %d", 0, len(watcher.watchers))
	}
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
	for i, test := range tests {
		watcher := newWatcher(t)
		for j, name := range test.paths {
			watch(t, watcher, name, string(j), dummy)
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
	}
}

func Testwatch(t *testing.T) {
	watcher := newWatcher(t)
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

}

func TestFlushDir(t *testing.T) {

}

func TestUnWatch(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	watch(t, watcher, name, "test", dummy)
	unwatch(t, watcher, name, "test")
	if len(watcher.watched) != 0 {
		t.Errorf("Expected watcheds be empty, but got %s", watcher.watched)
	}
}

func TestUnWatchAll(t *testing.T) {

}

func TestUnWatchDirectory(t *testing.T) {
	name := "testdata/dummy.txt"
	dir := "testdata"
	watcher := newWatcher(t)
	watch(t, watcher, name, "test", dummy)
	watch(t, watcher, dir, "test", nil)
	if !equal(watcher.watchers, []string{"testdata"}) {
		t.Fatalf("Expected watchers be equal to %s, but got %s", []string{"testdata"}, watcher.watchers)
	}
	unwatch(t, watcher, dir, "test")
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", []string{name}, watcher.watchers)
	}
}

func TestUnWatchOneOfSubscribers(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	watch(t, watcher, name, "test", dummy)
	watch(t, watcher, name, "test2", dummy)
	if len(watcher.watched[name]) != 2 {
		t.Fatalf("Expected watched[%s] length be %d, but got %d", name, 2, len(watcher.watched[name]))
	}
	unwatch(t, watcher, name, "test")
	if !equal(watcher.watchers, []string{name}) {
		t.Errorf("Expected watchers be equal to %s, but got %s", []string{name}, watcher.watchers)
	}
	if len(watcher.watched[name]) != 1 {
		t.Errorf("Expected watched[%s] length be %d, but got %d", name, 1, len(watcher.watched[name]))
	}
}

func TestunWatch(t *testing.T) {

}

func TestRemoveWatch(t *testing.T) {

}

func TestRemoveDir(t *testing.T) {

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
	watcher := newWatcher(t)
	v := &dumView{name: name}
	watch(t, watcher, name, "test", v.Reload)
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
	watcher := newWatcher(t)
	v := &dumView{name: dir}
	watch(t, watcher, dir, "test", v.Reload)
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
	watcher := newWatcher(t)
	v := &dumView{name: name}
	watch(t, watcher, name, "test", v.Reload)
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
	os.Remove(name)
}

func TestDeleteEvent(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	v := &dumView{name: name}
	watch(t, watcher, name, "test", v.Reload)
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
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	ioutil.WriteFile(name, []byte(""), 0644)
}

func TestRenameEvent(t *testing.T) {
	name := "testdata/dummy.txt"
	watcher := newWatcher(t)
	v := &dumView{name: name}
	watch(t, watcher, name, "test", v.Reload)
	go watcher.Observe()

	os.Rename(name, "testdata/rename.txt")
	time.Sleep(time.Millisecond * 50)
	if v.Text() != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text())
	}
	os.Rename("testdata/rename.txt", name)
}
