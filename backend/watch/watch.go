// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/howeyc/fsnotify"
	"github.com/limetext/lime/backend/log"
)

// Wrapper around fsnotify watcher to suit lime needs
// Enables:
// 		- Watching directories, we will have less individual watchers
// 		- Have multiple subscribers on one file or directory resolves #285
// 		- Watching a path which doesn't exist yet
// 		- Watching and applying action on certain events
type Watcher struct {
	wchr     *fsnotify.Watcher
	watched  map[string]actions // All watched paths
	watchers []string           // helper variable for paths we created watcher on
	dirs     []string           // helper variable for dirs we are watching
	lock     sync.Mutex
}

func NewWatcher() (*Watcher, error) {
	wchr, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	watched := make(map[string]actions)
	watchers := make([]string, 0)
	dirs := make([]string, 0)

	return &Watcher{wchr: wchr, watched: watched, watchers: watchers, dirs: dirs}, nil
}

func (w *Watcher) Watch(name, key string, act func(), events ...int) error {
	log.Finest("Watch(%s)", name)
	if key == "" {
		return errors.New(fmt.Sprintf("Couldn't Watch(%s) with empty key", name))
	}
	fi, err := os.Stat(name)
	isDir := err == nil && fi.IsDir()
	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if os.IsNotExist(err) {
		log.Debug("File doesn't exist, Watching parent dir")
		if err := w.Watch(filepath.Dir(name), "Dir", nil); err != nil {
			return err
		}
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	if exist(w.dirs, name) && act == nil {
		return nil
	}
	w.add(name, key, act, newEvent(events))
	// If exists in watchers we are already watching the path
	// no need to watch again just adding the action
	// Or
	// If the file is under one of watched dirs
	// no need to create watcher
	if exist(w.watchers, name) || (!isDir && exist(w.dirs, filepath.Dir(name))) {
		return nil
	}
	if err := w.watch(name); err != nil {
		return err
	}
	if isDir {
		w.flushDir(name)
	}
	return nil
}

func (w *Watcher) watch(name string) error {
	if err := w.wchr.Watch(name); err != nil {
		return err
	}
	w.watchers = append(w.watchers, name)
	return nil
}

func (w *Watcher) add(name, key string, act func(), ev int) {
	_, exist := w.watched[name]
	if !exist {
		w.watched[name] = make(actions)
	}
	w.watched[name][key] = action{act, ev}
}

// Remove watchers created on files under this directory because
// one watcher on the parent directory is enough for all of them
func (w *Watcher) flushDir(name string) {
	if exist(w.dirs, name) {
		return
	}
	w.dirs = append(w.dirs, name)
	for _, p := range w.watchers {
		if filepath.Dir(p) != name {
			continue
		}
		if err := w.wchr.RemoveWatch(p); err != nil {
			log.Error("Couldn't unwatch file %s: %s", p, err)
			continue
		}
		w.watchers = remove(w.watchers, p)
	}
}

func (w *Watcher) UnWatch(name, key string) error {
	log.Finest("UnWatch(%s)", name)
	w.lock.Lock()
	defer w.lock.Unlock()
	if key == "" {
		return w.unWatch(name)
	}
	delete(w.watched[name], key)
	if len(w.watched[name]) == 0 {
		return w.unWatch(name)
	}
	return nil
}

func (w *Watcher) unWatch(name string) error {
	if err := w.removeWatch(name); err != nil {
		return err
	}
	delete(w.watched, name)
	if exist(w.dirs, name) {
		w.removeDir(name)
	}
	return nil
}

func (w *Watcher) removeWatch(name string) error {
	if err := w.wchr.RemoveWatch(name); err != nil {
		return err
	}
	w.watchers = remove(w.watchers, name)
	return nil
}

// Put back watchers on watching files under the directory
func (w *Watcher) removeDir(name string) {
	for p, _ := range w.watched {
		if filepath.Dir(p) == name {
			if err := w.wchr.Watch(p); err != nil {
				log.Error("Could not watch: %s", err)
				continue
			}
			w.watchers = append(w.watchers, p)
		}
	}
	w.dirs = remove(w.dirs, name)
}

// Moves action with provided key(all acitons if key is empty)
// To another path(dest)
func (w *Watcher) Move(name, dest, key string) error {
	log.Finest("Moving watched actions with key %s from %s to %s", key, name, dest)
	w.lock.Lock()
	acs, ex := w.watched[name]
	w.lock.Unlock()
	if !ex {
		return errors.New(fmt.Sprintf("Moving path(%s) doesn't exists in watcheds", name))
	}
	if key != "" {
		ac, ex := acs[key]
		if !ex {
			return errors.New(fmt.Sprintf("key(%s) doesn't exists in actions", key))
		}
		if err := w.move(name, dest, key, ac); err != nil {
			return err
		}
		return w.UnWatch(name, key)
	}
	for k, ac := range acs {
		if err := w.move(name, dest, k, ac); err != nil {
			return err
		}
	}
	return w.UnWatch(name, key)
}

func (w *Watcher) move(name, dest, key string, ac action) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	_, ex := w.watched[dest]
	if ex {
		w.watched[dest][key] = ac
		return nil
	}
	w.lock.Unlock()
	if err := w.Watch(dest, key, ac.fn, ac.ev); err != nil {
		return err
	}
	w.lock.Lock()
	return nil
}

func (w *Watcher) Observe() {
	for {
		select {
		case ev := <-w.wchr.Event:
			func() {
				w.lock.Lock()
				defer w.lock.Unlock()
				event := evnt(*ev)
				if acs, ex := w.watched[ev.Name]; ex {
					acs.applyAll(event)
				}
				if !exist(w.dirs, ev.Name) {
					// The watcher will be removed if the file is deleted
					// so we need to watch the parent directory for when the
					// file is created again
					if ev.IsDelete() {
						w.watchers = remove(w.watchers, ev.Name)
						w.lock.Unlock()
						w.Watch(filepath.Dir(ev.Name), "Dir", nil)
						w.lock.Lock()
					}
					// We will apply parent directory actions to, if one of the files
					// inside the directory has changed
					if acs, ex := w.watched[filepath.Dir(ev.Name)]; ex {
						acs.applyAll(event)
					}
					return
				}
				// If the ev.Name refers to a directory run all watched actions
				// for wathed files under the directory
				for p, acs := range w.watched {
					if filepath.Dir(p) == ev.Name {
						acs.applyAll(event)
					}
				}
			}()
		case err := <-w.wchr.Error:
			log.Error("Watcher error: %s", err)
		}
	}
}

func evnt(ev fsnotify.FileEvent) int {
	event := 0
	if ev.IsCreate() {
		event |= CREATE
	}
	if ev.IsDelete() {
		event |= DELETE
	}
	if ev.IsModify() {
		event |= MODIFY
	}
	if ev.IsRename() {
		event |= RENAME
	}
	return event
}
