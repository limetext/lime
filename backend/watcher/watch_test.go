package watcher

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func dum() {
	return
}

func TestWatch(t *testing.T) {
	tests := []struct {
		paths       map[string]func()
		expWatched  []string
		expWatchers []string
	}{
		{
			map[string]func(){
				"testdata/dummy.txt": dum,
				"testdata/test.txt":  dum,
			},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
		},
		{
			map[string]func(){
				"testdata":           dum,
				"testdata/dummy.txt": dum,
				"testdata/test.txt":  dum,
			},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
		{
			map[string]func(){
				"testdata/dummy.txt": dum,
				"testdata/test.txt":  dum,
				"testdata":           dum,
			},
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
		},
		// 2 path refer same file or dir but different(e.g abs path and relative path)
	}
	for i, test := range tests {
		for path, action := range test.paths {
			Watch(path, action)
		}
		if len(watched) != len(test.expWatched) {
			t.Errorf("Test %d: Expected len of watched %d, but got %d", i, len(test.expWatched), len(watched))
		}
		for _, p := range test.expWatched {
			if _, exist := watched[p]; !exist {
				t.Errorf("Test %d: Expected %s exist in watched", i, p)
			}
		}
		if !reflect.DeepEqual(test.expWatchers, watchers) {
			t.Errorf("Test %d: Expected watchers %v, but got %v", i, test.expWatchers, watchers)
		}
		for _, path := range test.expWatchers {
			UnWatch(path)
		}
		reinit()
	}
}

func TestUnWatch(t *testing.T) {
	tests := []struct {
		watchs      []string
		unWatchs    []string
		expWatched  []string
		expWatchers []string
	}{
		{
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/dummy.txt"},
			[]string{"testdata/test.txt"},
			[]string{"testdata/test.txt"},
		},
		{
			[]string{"testdata", "testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata"},
			[]string{"testdata/dummy.txt", "testdata/test.txt"},
			[]string{"testdata/test.txt", "testdata/dummy.txt"},
		},
	}
	for i, test := range tests {
		for _, path := range test.watchs {
			Watch(path, dum)
		}
		for _, path := range test.unWatchs {
			UnWatch(path)
		}
		if len(watched) != len(test.expWatched) {
			t.Errorf("Test %d: Expected len of watched %d, but got %d", i, len(test.expWatched), len(watched))
		}
		for _, p := range test.expWatched {
			if _, exist := watched[p]; !exist {
				t.Errorf("Test %d: Expected %s exist in watched", i, p)
			}
		}
		if !reflect.DeepEqual(test.expWatchers, watchers) {
			t.Errorf("Test %d: Expected watchers %v, but got %v", i, test.expWatchers, watchers)
		}
		for _, path := range test.expWatchers {
			UnWatch(path)
		}
		reinit()
	}
}

type dumView struct {
	Text string
	Name string
}

func (d *dumView) Reload() {
	d.Text = "Reloaded"
}

func (d *dumView) Rename() {
	d.Name = "Renamed"
}

func TestObserve(t *testing.T) {
	path := "testdata/test.txt"
	v := new(dumView)
	Watch(path, v.Reload)
	go Observe()

	if err := ioutil.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if v.Text != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text)
	}
	ioutil.WriteFile(path, []byte(""), 0644)
	UnWatch(path)
	reinit()
}

func TestObserveDirectory(t *testing.T) {
	dir := "testdata"
	path := "testdata/test.txt"
	v := new(dumView)
	Watch(path, v.Reload)
	Watch(dir, nil)
	go Observe()

	if !reflect.DeepEqual(watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watchers)
	}
	if err := ioutil.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if v.Text != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text)
	}
	ioutil.WriteFile(path, []byte(""), 0644)
	UnWatch(dir)
	reinit()
}

func TestObserveCreateEvent(t *testing.T) {
	path := "testdata/new.txt"
	v := new(dumView)
	Watch(path, v.Reload)
	go Observe()

	if !reflect.DeepEqual(watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{path}, watchers)
	}
	if err := ioutil.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error: %s", err)
	}
	time.Sleep(time.Millisecond * 50)
	if v.Text != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text)
	}
	os.Remove(path)
	UnWatch(path)
	reinit()
}

func TestObserveDeleteEvent(t *testing.T) {
	path := "testdata/dummy.txt"
	v := new(dumView)
	Watch(path, v.Reload)
	go Observe()

	os.Remove(path)
	time.Sleep(time.Millisecond * 50)
	if v.Text != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text)
	}
	if !reflect.DeepEqual(watchers, []string{"testdata"}) {
		t.Errorf("Expected watchers be equal to %v, but got %v", []string{"testdata"}, watchers)
	}
	UnWatch("testdata")
	ioutil.WriteFile(path, []byte(""), 0644)
	reinit()
}

func TestObserveRenameEvent(t *testing.T) {
	path := "testdata/dummy.txt"
	v := new(dumView)
	Watch(path, v.Reload)
	go Observe()

	os.Rename(path, "testdata/rename.txt")
	time.Sleep(time.Millisecond * 50)
	if v.Text != "Reloaded" {
		t.Errorf("Expected dumView Text %s, but got %s", "Reloaded", v.Text)
	}
	os.Rename("testdata/rename.txt", path)
	UnWatch(path)
	reinit()
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
		if exp := remove(test.slice, test.remove); !reflect.DeepEqual(exp, test.exp) {
			t.Errorf("Test %d: Expected %v be equal to %v", i, exp, test.exp)
		}
	}
}
