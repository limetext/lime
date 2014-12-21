// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/limetext/lime/backend/log"
	"gopkg.in/fsnotify.v1"
)

type (
	FileChangedCallback interface {
		FileChanged(string)
	}
	FileCreatedCallback interface {
		FileCreated(string)
	}
	FileRemovedCallback interface {
		FileRemoved(string)
	}
	FileRenamedCallback interface {
		FileRenamed(string)
	}

	// Wrapper around fsnotify watcher to suit lime needs
	// Enables:
	// 		- Watching directories, we will have less individual watchers
	// 		- Have multiple subscribers on one file or directory resolves #285
	// 		- Watching a path which doesn't exist yet
	// 		- Watching and applying action on certain events
	Watcher struct {
		wchr     *fsnotify.Watcher
		watched  map[string][]interface{}
		watchers []string // helper variable for paths we created watcher on
		dirs     []string // helper variable for dirs we are watching
		lock     sync.Mutex
	}
)

func NewWatcher() (*Watcher, error) {
	wchr, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{wchr: wchr}
	w.watched = make(map[string][]interface{})
	w.watchers = make([]string, 0)
	w.dirs = make([]string, 0)

	return w, nil
}

func (w *Watcher) Watch(name string, cb interface{}) error {
	log.Finest("Watch(%s)", name)
	fi, err := os.Stat(name)
	isDir := err == nil && fi.IsDir()
	// If the file doesn't exist currently we will add watcher for file
	// directory and look for create event inside the directory
	if os.IsNotExist(err) {
		log.Fine("%s doesn't exist, Watching parent directory", name)
		if err := w.Watch(filepath.Dir(name), nil); err != nil {
			return err
		}
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	if err := w.add(name, cb); err != nil {
		if !isDir {
			return err
		}
		if exist(w.dirs, name) {
			log.Debug("%s is watched already", name)
			return nil
		}
	}
	// If exists in watchers we are already watching the path
	// Or
	// If the file is under one of watched dirs
	//
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

func (w *Watcher) add(name string, cb interface{}) error {
	numok := 0
	if _, ok := cb.(FileChangedCallback); ok {
		numok++
	}
	if _, ok := cb.(FileCreatedCallback); ok {
		numok++
	}
	if _, ok := cb.(FileRemovedCallback); ok {
		numok++
	}
	if _, ok := cb.(FileRenamedCallback); ok {
		numok++
	}
	if numok == 0 {
		return errors.New("The callback argument does satisfy any File*Callback interfaces")
	}
	w.watched[name] = append(w.watched[name], cb)
	return nil
}

func (w *Watcher) watch(name string) error {
	if err := w.wchr.Add(name); err != nil {
		return err
	}
	w.watchers = append(w.watchers, name)
	return nil
}

// Remove watchers created on files under this directory because
// one watcher on the parent directory is enough for all of them
func (w *Watcher) flushDir(name string) {
	if exist(w.dirs, name) {
		return
	}
	w.dirs = append(w.dirs, name)
	for _, p := range w.watchers {
		if filepath.Dir(p) == name && !exist(w.dirs, p) {
			if err := w.removeWatch(p); err != nil {
				log.Errorf("Couldn't unwatch file %s: %s", p, err)
			}
		}
	}
}

func (w *Watcher) UnWatch(name string, cb interface{}) error {
	log.Finest("UnWatch(%s)", name)
	w.lock.Lock()
	defer w.lock.Unlock()
	if cb == nil {
		return w.unWatch(name)
	}
	for i, c := range w.watched[name] {
		if c == cb {
			w.watched[name][i] = w.watched[name][len(w.watched[name])-1]
			w.watched[name][len(w.watched[name])-1] = nil
			w.watched[name] = w.watched[name][:len(w.watched[name])-1]
			break
		}
	}
	if len(w.watched[name]) == 0 {
		w.unWatch(name)
	}
	return nil
}

func (w *Watcher) unWatch(name string) error {
	delete(w.watched, name)
	if err := w.removeWatch(name); err != nil {
		return err
	}
	return nil
}

func (w *Watcher) removeWatch(name string) error {
	if err := w.wchr.Remove(name); err != nil {
		return err
	}
	w.watchers = remove(w.watchers, name)
	if exist(w.dirs, name) {
		w.removeDir(name)
	}
	return nil
}

// Put back watchers on watching files under the directory
func (w *Watcher) removeDir(name string) {
	for p, _ := range w.watched {
		if filepath.Dir(p) == name {
			if err := w.watch(p); err != nil {
				log.Errorf("Could not watch: %s", err)
				continue
			}
		}
	}
	w.dirs = remove(w.dirs, name)
}

func (w *Watcher) Observe() {
	for {
		select {
		case ev := <-w.wchr.Events:
			func() {
				w.lock.Lock()
				defer w.lock.Unlock()
				w.apply(ev)
				name := ev.Name
				// If the name refers to a directory run all watched
				// callbacks for wathed files under the directory
				if exist(w.dirs, name) {
					for p, _ := range w.watched {
						if filepath.Dir(p) == name {
							ev.Name = p
							w.apply(ev)
						}
					}
				}
				dir := filepath.Dir(name)
				// The watcher will be removed if the file is deleted
				// so we need to watch the parent directory for when the
				// file is created again
				if ev.Op&fsnotify.Remove != 0 {
					w.watchers = remove(w.watchers, name)
					w.lock.Unlock()
					w.Watch(dir, nil)
					w.lock.Lock()
				}
				// We will apply parent directory FileChanged callbacks to,
				// if one of the files inside the directory has changed
				if cbs, exist := w.watched[dir]; ev.Op&fsnotify.Write != 0 && exist {
					for _, cb := range cbs {
						if c, ok := cb.(FileChangedCallback); ok {
							c.FileChanged(dir)
						}
					}
				}

			}()
		case err := <-w.wchr.Errors:
			log.Errorf("Watcher error: %s", err)
		}
	}
}

func (w *Watcher) apply(ev fsnotify.Event) {
	for _, cb := range w.watched[ev.Name] {
		if ev.Op&fsnotify.Create != 0 {
			if c, ok := cb.(FileCreatedCallback); ok {
				c.FileCreated(ev.Name)
			}
		}
		if ev.Op&fsnotify.Write != 0 {
			if c, ok := cb.(FileChangedCallback); ok {
				c.FileChanged(ev.Name)
			}
		}
		if ev.Op&fsnotify.Remove != 0 {
			if c, ok := cb.(FileRemovedCallback); ok {
				c.FileRemoved(ev.Name)
			}
		}
		if ev.Op&fsnotify.Rename != 0 {
			if c, ok := cb.(FileRenamedCallback); ok {
				c.FileRenamed(ev.Name)
			}
		}
	}
}
