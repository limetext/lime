package watcher

import (
	"code.google.com/p/log4go"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
	"sync"
)

var (
	watched  map[string][]func() // All watched paths including their actions
	watchers []string            // helper variable for paths we created watcher on
	dirs     []string            // helper variable for dirs we are watching
	wchr     *fsnotify.Watcher
	lock     sync.Mutex
)

func existIn(paths []string, path string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

func remove(slice []string, path string) []string {
	for i, el := range slice {
		if el == path {
			slice[i], slice = slice[len(slice)-1], slice[:len(slice)-1]
			break
		}
	}
	return slice
}

func Watch(path string, action func()) {
	log4go.Finest("Watch(%s)", path)
	fi, err := os.Stat(path)
	dir := err == nil && fi.IsDir()
	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if !dir && os.IsNotExist(err) {
		Watch(filepath.Dir(path), nil)
	}
	// If the path points to a file and there is no action
	// Don't watch
	if !dir && action == nil {
		log4go.Error("No action for watching the file")
		return
	}
	lock.Lock()
	defer lock.Unlock()
	// If exists in watchers we are already watching the path
	// no need to watch again just adding the action
	if existIn(watchers, path) {
		if action != nil {
			watched[path] = append(watched[path], action)
		}
		return
	}
	// If the file is under one of watched dirs
	// no need to create watcher
	if !dir && existIn(dirs, filepath.Dir(path)) {
		watched[path] = append(watched[path], action)
		return
	}
	if err := wchr.Watch(path); err != nil {
		log4go.Error("Could not watch: %s", err)
		return
	}
	watchers = append(watchers, path)
	watched[path] = append(watched[path], action)
	if dir {
		dirs = append(dirs, path)
		for _, p := range watchers {
			if filepath.Dir(p) == path {
				if err := wchr.RemoveWatch(p); err != nil {
					log4go.Error("Couldn't unwatch file: %s", err)
					return
				}
				watchers = remove(watchers, p)
			}
		}
	}
}

func UnWatch(path string) {
	lock.Lock()
	defer lock.Unlock()
	log4go.Finest("UnWatch(%s)", path)
	if existIn(watchers, path) {
		if existIn(dirs, path) {
			for p, _ := range watched {
				if filepath.Dir(p) == path && !existIn(watchers, p) {
					if err := wchr.Watch(p); err != nil {
						log4go.Error("Could not watch: %s", err)
						return
					}
					watchers = append(watchers, p)
				}
			}
		}
		if err := wchr.RemoveWatch(path); err != nil {
			log4go.Error("Couldn't unwatch file: %s", err)
			return
		}
		watchers = remove(watchers, path)
	}
	dirs = remove(dirs, path)
	delete(watched, path)
}

func Observe() {
	for {
		select {
		case ev := <-wchr.Event:
			// The watcher will be removed if the file is deleted
			// so we need to watch the parent directory for when the
			// file is created again
			if ev.IsDelete() {
				remove(watchers, ev.Name)
				Watch(filepath.Dir(ev.Name), nil)
			}
			func() {
				lock.Lock()
				defer lock.Unlock()
				if actions, exist := watched[ev.Name]; exist {
					for _, action := range actions {
						if action != nil {
							action()
						}
					}
				}
				if existIn(dirs, ev.Name) {
					for p, actions := range watched {
						if filepath.Dir(p) == ev.Name && !existIn(watchers, p) {
							for _, action := range actions {
								action()
							}
						}
					}
				}
			}()
		case err := <-wchr.Error:
			log4go.Error("Watcher error: %s", err)
		}
	}
}

func init() {
	var err error
	if wchr, err = fsnotify.NewWatcher(); err != nil {
		log4go.Error("Could not create watcher due to: %s", err)
		return
	}
	watched = make(map[string][]func())
	watchers = nil
	dirs = nil
}
