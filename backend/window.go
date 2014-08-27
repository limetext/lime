// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"github.com/quarnster/util/text"
	"io/ioutil"
	"runtime/debug"
	"sync"
)

type Window struct {
	text.HasId
	text.HasSettings
	views       []*View
	active_view *View
	lock        sync.Mutex
}

func (w *Window) NewFile() *View {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.views = append(w.views, newView(w))
	v := w.views[len(w.views)-1]

	v.Settings().SetParent(w)
	v.setBuffer(text.NewBuffer())
	v.selection.Clear()
	v.selection.Add(text.Region{0, 0})
	OnNew.Call(v)
	w.SetActiveView(v)

	return v
}

func (w *Window) Views() []*View {
	w.lock.Lock()
	defer w.lock.Unlock()
	ret := make([]*View, len(w.views))
	copy(ret, w.views)
	return ret
}

func (w *Window) remove(v *View) {
	w.lock.Lock()
	defer w.lock.Unlock()
	for i, vv := range w.views {
		if v == vv {
			end := len(w.views) - 1
			if i != end {
				copy(w.views[i:], w.views[i+1:])
			}
			w.views = w.views[:end]
			return
		}
	}
	log4go.Error("Wanted to remove view %+v, but it doesn't appear to be a child of this window", v)
}

func (w *Window) OpenFile(filename string, flags int) *View {
	v := w.NewFile()

	v.SetScratch(true)
	e := v.BeginEdit()
	v.Buffer().SetFileName(filename)
	if d, err := ioutil.ReadFile(filename); err != nil {
		log4go.Error("Couldn't load file %s: %s", filename, err)
	} else {
		v.Insert(e, 0, string(d))
	}
	v.EndEdit(e)
	v.SetScratch(false)
	OnLoad.Call(v)

	return v
}

func (w *Window) SetActiveView(v *View) *View{
	// w.lock.Lock()
	// defer w.lock.Unlock()

	if w.active_view != nil {
		OnDeactivated.Call(w.active_view)
	}
	w.active_view = v
	if w.active_view != nil {
		OnActivated.Call(w.active_view)
	}

	return v
}

func (w *Window) ActiveView() *View {
	return w.active_view
}

func (w *Window) ActiveViewIndex() int {
	for pos, view := range w.views {
		if view == w.active_view {
			return pos
		}
	}
	return -1
}

func (w *Window) NextView() *View {
	index := w.ActiveViewIndex()
	if index == -1 {
		return nil
	}

	index++
	if index == len(w.views) {
		return w.SetActiveView(w.views[0])
	} else {
		return w.SetActiveView(w.views[index])
	}
}

func (w *Window) PrevView() *View {
	index := w.ActiveViewIndex()
	if index == -1 {
		return nil
	}

	index--
	if index < 0 {
		return w.SetActiveView(w.views[len(w.views)-1])
	} else {
		return w.SetActiveView(w.views[index])
	}
}

func (w *Window) Close() {
	w.CloseAllViews()
	ed := GetEditor()
	ed.remove(w)
}

func (w *Window) CloseAllViews() {
	for _, v := range w.views {
		v.Close()
	}
}

func (w *Window) runCommand(c WindowCommand, name string) error {
	defer func() {
		if r := recover(); r != nil {
			log4go.Error("Paniced while running window command %s %v: %v\n%s", name, c, r, string(debug.Stack()))
		}
	}()
	return c.Run(w)
}
