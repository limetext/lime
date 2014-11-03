package watch

import (
	"github.com/howeyc/fsnotify"
	"github.com/limetext/lime/backend/log"
	"os"
	"path/filepath"
	"sync"
)

type (
	// Watched interface defines the methods that every
	// watched type should implement we use the Name() as
	// the path we should watch and Reload() as the action
	// we should take on create delete modify and rename event
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

	// Watching directory without any action on reload
	// this helps us to create less individual watchers
	// also we can apply actions on create events
	WatchedDir struct {
		name string
	}
)

func NewWatcher() (*Watcher, error) {
	wchr, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	watched := make(map[string][]Watched)
	watchers := make([]string, 0)
	dirs := make([]string, 0)

	return &Watcher{wchr: wchr, watched: watched, watchers: watchers, dirs: dirs}, nil
}

func (w *Watcher) Watch(watched Watched) error {
	log.Finest("Watch(%s)", watched.Name())
	name := watched.Name()
	fi, err := os.Stat(name)
	dir := err == nil && fi.IsDir()
	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if !dir && os.IsNotExist(err) {
		if err := w.Watch(NewWatchedDir(filepath.Dir(name))); err != nil {
			return err
		}
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	// If exists in watchers we are already watching the path
	// no need to watch again just adding the action
	if exist(w.watchers, name) {
		w.watched[name] = append(w.watched[name], watched)
		return nil
	}
	// If the file is under one of watched dirs
	// no need to create watcher
	if !dir && exist(w.dirs, filepath.Dir(name)) {
		w.watched[name] = append(w.watched[name], watched)
		return nil
	}
	if err := w.wchr.Watch(name); err != nil {
		return err
	}
	w.watchers = append(w.watchers, name)
	w.watched[name] = append(w.watched[name], watched)
	// If name refers to a directory and we created watcher on it we
	// will remove watchers created on files under this directory because
	// one watcher on the parent directory is enough for all of them
	if dir {
		w.dirs = append(w.dirs, name)
		for _, p := range w.watchers {
			if filepath.Dir(p) != name {
				continue
			}
			if err := w.wchr.RemoveWatch(p); err != nil {
				log.Error("Couldn't unwatch file: %s", err)
				continue
			}
			w.watchers = remove(w.watchers, p)
		}
	}
	return nil
}

func (w *Watcher) UnWatch(watched Watched) error {
	log.Finest("UnWatch(%s)", watched.Name())
	w.lock.Lock()
	defer w.lock.Unlock()
	name := watched.Name()
	watcheds := w.watched[name]
	l := len(w.watched[name])
	for i, wchd := range watcheds {
		if wchd == watched {
			w.watched[name][i], w.watched[name][l-1], w.watched[name] = w.watched[name][l-1], nil, w.watched[name][:l-1]
			l -= 1
			break
		}
	}
	if l != 0 {
		return nil
	}
	w.watchers = remove(w.watchers, name)
	delete(w.watched, name)
	if exist(w.watchers, name) {
		if err := w.wchr.RemoveWatch(name); err != nil {
			return err
		}
	}
	if !exist(w.dirs, name) {
		return nil
	}
	// If name refers to a watched directory we should put back
	// watchers on watching files under the directory
	for p, _ := range w.watched {
		if filepath.Dir(p) == name && !exist(w.watchers, p) {
			if err := w.wchr.Watch(p); err != nil {
				log.Error("Could not watch: %s", err)
				continue
			}
			w.watchers = append(w.watchers, p)
		}
	}
	w.dirs = remove(w.dirs, name)
	return nil
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
				watcheds := w.watched[ev.Name]
				exInDir := exist(w.dirs, ev.Name)
				if !exInDir {
					// We will apply directory actions to, if one of the files
					// inside the directory has changed
					watcheds = append(watcheds, w.watched[filepath.Dir(ev.Name)]...)
				}
				if len(watcheds) == 0 {
					return
				}
				for _, watched := range watcheds {
					watched.Reload()
				}
				if !exInDir {
					return
				}
				// If the ev.Name refers to a directory run all watched actions
				// for wathedq files under the directory
				for p, watcheds := range w.watched {
					if filepath.Dir(p) == ev.Name {
						for _, watched := range watcheds {
							watched.Reload()
						}
					}
				}
			}()
		case err := <-w.wchr.Error:
			log.Error("Watcher error: %s", err)
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
