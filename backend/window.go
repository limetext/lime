package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"lime/backend/primitives"
)

type Window struct {
	primitives.HasId
	HasSettings
	views []*View
}

func (w *Window) NewView() *View {
	w.views = append(w.views, newView(w))
	v := w.views[len(w.views)-1]
	v.Settings().Parent = w
	v.setBuffer(&primitives.Buffer{})
	v.selection.Clear()
	v.selection.Add(primitives.Region{0, 0})
	OnNew.Call(v)
	return v
}

func (w *Window) Views() []*View {
	return w.views
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

func (w *Window) runCommand(c WindowCommand, name string, args Args) error {
	defer func() {
		if r := recover(); r != nil {
			log4go.Error("Paniced while running window command %s %v: %v", name, args, r)
		}
	}()
	return c.Run(w, args)
}
