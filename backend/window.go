package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"lime/backend/primitives"
)

type Window struct {
	HasId
	HasSettings
	views []View
}

func (w *Window) NewView() *View {
	w.views = append(w.views, View{window: w})
	v := &w.views[len(w.views)-1]
	v.settings.Parent = w
	v.setBuffer(&primitives.Buffer{})
	v.selection.Clear()
	v.selection.Add(primitives.Region{0, 0})
	OnNew.Call(v)
	return v
}

func (w *Window) OpenFile(filename string, flags int) *View {
	v := w.NewView()
	v.SetScratch(true)
	e := v.BeginEdit()
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

// TODO(q): ActiveView should return the actual active view
func (w *Window) ActiveView() *View {
	if len(w.views) > 0 {
		return &w.views[0]
	}
	return nil
}
