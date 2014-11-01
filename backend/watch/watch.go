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
		Path() string
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
		path string
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
	path := watched.Path()
	fi, err := os.Stat(path)
	dir := err == nil && fi.IsDir()
	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if !dir && os.IsNotExist(err) {
		w.Watch(NewWatchedDir(filepath.Dir(path)))
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	// If exists in watchers we are already watching the path
	// no need to watch again just adding the action
	if exist(w.watchers, path) {
		w.watched[path] = append(w.watched[path], watched)
		return
	}
	// If the file is under one of watched dirs
	// no need to create watcher
	if !dir && exist(w.dirs, filepath.Dir(path)) {
		w.watched[path] = append(w.watched[path], watched)
		return
	}
	if err := w.wchr.Watch(path); err != nil {
		log4go.Error("Could not watch: %s", err)
		return
	}
	w.watchers = append(w.watchers, path)
	w.watched[path] = append(w.watched[path], watched)
	if dir {
		w.dirs = append(w.dirs, path)
		for _, p := range w.watchers {
			if filepath.Dir(p) != path {
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
	w.lock.Lock()
	defer w.lock.Unlock()
	path := watched.Path()
	watcheds, exst := w.watched[path]
	if !exst {
		return
	}
	l := len(w.watched[path])
	for i, wchd := range watcheds {
		if wchd == watched {
			w.watched[path][i], w.watched[path][l-1], w.watched[path] = w.watched[path][l-1], nil, w.watched[path][:l-1]
			l -= 1
			break
		}
	}
	if l == 0 {
		w.watchers = remove(w.watchers, path)
		delete(w.watched, path)
		if err := w.wchr.RemoveWatch(path); err != nil {
			log4go.Error("Couldn't unwatch file: %s", err)
		}
	}
	if !exist(w.dirs, path) {
		return
	}
	for p, _ := range w.watched {
		if filepath.Dir(p) == path && !exist(w.watchers, p) {
			if err := w.wchr.Watch(p); err != nil {
				log4go.Error("Could not watch: %s", err)
				continue
			}
			w.watchers = append(w.watchers, p)
		}
	}
	w.dirs = remove(w.dirs, path)
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
					w.watchers = remove(w.watchers, ev.Name)
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

func NewWatchedDir(path string) *WatchedDir {
	return &WatchedDir{path}
}

func (wd *WatchedDir) Path() string {
	return wd.path
}

func (wd *WatchedDir) Reload() {}

// Helper function checking an element exists in a slice
func exist(paths []string, path string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

// Helper function for removing an element from slice
func remove(slice []string, path string) []string {
	for i, el := range slice {
		if el == path {
			slice[i], slice = slice[len(slice)-1], slice[:len(slice)-1]
			break
		}
	}
	return slice
}
