package watch

import (
	"code.google.com/p/log4go"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
	"sync"
)

type (
	Watched interface {
		Name() string
		Reload()
	}

	Watcher struct {
		wchr     *fsnotify.Watcher
		watched  map[string][]Watched // All watched paths
		watchers []string             // helper variable for paths we created watcher on
		dirs     []string             // helper variable for dirs we are watching
		lock     sync.Mutex
	}

	WatchedDir struct {
		name string
	}
)

func NewWatcher() *Watcher {
	wchr, err := fsnotify.NewWatcher()
	if err != nil {
		log4go.Error("Could not create watcher due to: %s", err)
		return nil
	}
	watched := make(map[string][]Watched)
	watchers := make([]string, 0)
	dirs := make([]string, 0)

	return &Watcher{wchr: wchr, watched: watched, watchers: watchers, dirs: dirs}
}

func (w *Watcher) Watch(watched Watched) {
	log4go.Finest("Watch(%s)", watched.Name())
	name := watched.Name()
	fi, err := os.Stat(name)
	dir := err == nil && fi.IsDir()
	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if !dir && os.IsNotExist(err) {
		w.Watch(NewWatchedDir(filepath.Dir(name)))
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	// If exists in watchers we are already watching the path
	// no need to watch again just adding the action
	if exist(w.watchers, name) {
		w.watched[name] = append(w.watched[name], watched)
		return
	}
	// If the file is under one of watched dirs
	// no need to create watcher
	if !dir && exist(w.dirs, filepath.Dir(name)) {
		w.watched[name] = append(w.watched[name], watched)
		return
	}
	if err := w.wchr.Watch(name); err != nil {
		log4go.Error("Could not watch: %s", err)
		return
	}
	w.watchers = append(w.watchers, name)
	w.watched[name] = append(w.watched[name], watched)
	if dir {
		w.dirs = append(w.dirs, name)
		for _, p := range w.watchers {
			if filepath.Dir(p) != name {
				continue
			}
			if err := w.wchr.RemoveWatch(p); err != nil {
				log4go.Error("Couldn't unwatch file: %s", err)
				continue
			}
			w.watchers = remove(w.watchers, p)
		}
	}
}

func (w *Watcher) UnWatch(watched Watched) {
	log4go.Finest("UnWatch(%s)", watched.Name())
	w.lock.Lock()
	defer w.lock.Unlock()
	name := watched.Name()
	watcheds, exst := w.watched[name]
	if !exst {
		return
	}
	l := len(w.watched[name])
	for i, wchd := range watcheds {
		if wchd == watched {
			w.watched[name][i], w.watched[name][l-1], w.watched[name] = w.watched[name][l-1], nil, w.watched[name][:l-1]
			l -= 1
			break
		}
	}
	if l == 0 {
		w.watchers = remove(w.watchers, name)
		delete(w.watched, name)
		if err := w.wchr.RemoveWatch(name); err != nil {
			log4go.Error("Couldn't unwatch file: %s", err)
		}
	}
	if !exist(w.dirs, name) {
		return
	}
	for p, _ := range w.watched {
		if filepath.Dir(p) == name && !exist(w.watchers, p) {
			if err := w.wchr.Watch(p); err != nil {
				log4go.Error("Could not watch: %s", err)
				continue
			}
			w.watchers = append(w.watchers, p)
		}
	}
	w.dirs = remove(w.dirs, name)
}

func (w *Watcher) Observe() {
	for {
		select {
		case ev := <-w.wchr.Event:
			func() {
				// The watcher will be removed if the file is deleted
				// so we need to watch the parent directory for when the
				// file is created again
				if ev.IsDelete() {
					w.lock.Lock()
					w.watchers = remove(w.watchers, ev.Name)
					w.lock.Unlock()
					w.Watch(NewWatchedDir(filepath.Dir(ev.Name)))
				}
				w.lock.Lock()
				defer w.lock.Unlock()
				watcheds, exst := w.watched[ev.Name]
				if !exst {
					return
				}
				for _, watched := range watcheds {
					watched.Reload()
				}
				if !exist(w.dirs, ev.Name) {
					return
				}
				for p, watcheds := range w.watched {
					if filepath.Dir(p) == ev.Name && !exist(w.watchers, p) {
						for _, watched := range watcheds {
							watched.Reload()
						}
					}
				}
			}()
		case err := <-w.wchr.Error:
			log4go.Error("Watcher error: %s", err)
		}
	}
}

func NewWatchedDir(name string) *WatchedDir {
	return &WatchedDir{name}
}

func (wd *WatchedDir) Name() string {
	return wd.name
}

func (wd *WatchedDir) Reload() {}

// Helper function checking an element exists in a slice
func exist(paths []string, name string) bool {
	for _, p := range paths {
		if p == name {
			return true
		}
	}
	return false
}

// Helper function for removing an element from slice
func remove(slice []string, name string) []string {
	for i, el := range slice {
		if el == name {
			slice[i], slice = slice[len(slice)-1], slice[:len(slice)-1]
			break
		}
	}
	return slice
}
