package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"lime/backend/primitives"
	"runtime/debug"
	"sync"
)

type Window struct {
	primitives.HasId
	primitives.HasSettings
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
	v.setBuffer(primitives.NewBuffer())
	v.selection.Clear()
	v.selection.Add(primitives.Region{0, 0})
	OnNew.Call(v)
	w.SetActiveView(v)
	return v
}

func (w *Window) Views() []*View {
	w.lock.Lock()
	defer w.lock.Unlock()
	ret := make([]*View, 0, len(w.views))
	copy(ret, w.views)
	return ret
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

func (w *Window) SetActiveView(v *View) {
	// w.lock.Lock()
	// defer w.lock.Unlock()

	if w.active_view != nil {
		OnDeactivated.Call(w.active_view)
	}
	w.active_view = v
	if w.active_view != nil {
		OnActivated.Call(w.active_view)
	}
}

func (w *Window) ActiveView() *View {
	return w.active_view
}

func (w *Window) runCommand(c WindowCommand, name string, args Args) error {
	defer func() {
		if r := recover(); r != nil {
			log4go.Error("Paniced while running window command %s %v: %v\n%s", name, args, r, string(debug.Stack()))
		}
	}()
	return c.Run(w, args)
}
