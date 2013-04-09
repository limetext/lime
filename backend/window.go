package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"lime/backend/primitives"
	"runtime/debug"
)

type Window struct {
	primitives.HasId
	HasSettings
	views []*View
}

func (w *Window) NewFile() *View {
	w.views = append(w.views, newView(w))
	v := w.views[len(w.views)-1]
	v.Settings().SetParent(w)
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

func (w *Window) runCommand(c WindowCommand, name string, args Args) error {
	defer func() {
		if r := recover(); r != nil {
			log4go.Error("Paniced while running window command %s %v: %v\n%s", name, args, r, string(debug.Stack()))
		}
	}()
	return c.Run(w, args)
}
