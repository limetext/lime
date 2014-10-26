// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/packages"
	. "github.com/limetext/text"
	"io/ioutil"
)

type (
	Watched interface {
		Name() string
		Reload()
	}

	WatchedUserFile struct {
		view *View
	}

	WatchedPackage struct {
		pkg packages.Package
	}
)

func NewWatchedUserFile(view *View) *WatchedUserFile {
	return &WatchedUserFile{view}
}

func (o WatchedUserFile) String() string {
	return fmt.Sprintf("%s (%d)", o.Name(), o.view.Id())
}

func (o *WatchedUserFile) Name() string {
	return o.view.Buffer().FileName()
}

func (o *WatchedUserFile) Reload() {
	log.LogFinest("\"%v\".Reload()", o)

	view := o.view
	filename := o.Name()

	if saving, ok := view.Settings().Get("lime.saving", false).(bool); ok && saving {
		// This reload was triggered by ourselves saving to this file, so don't reload it
		return
	}
	if !GetEditor().Frontend().OkCancelDialog("File was changed by another program, reload?", "reload") {
		return
	}

	if d, err := ioutil.ReadFile(filename); err != nil {
		log.LogError("Could not read file: %s\n. Error was: %v", filename, err)
	} else {
		edit := view.BeginEdit()
		end := view.Buffer().Size()
		view.Replace(edit, Region{0, end}, string(d))
		view.EndEdit(edit)
	}
}

func NewWatchedPackage(pkg packages.Package) *WatchedPackage {
	return &WatchedPackage{pkg}
}

func (o *WatchedPackage) Name() string {
	return o.pkg.Name()
}

func (o *WatchedPackage) Reload() {
	o.pkg.Reload()
}

func (o *WatchedPackage) Package() packages.Package {
	return o.pkg
}
