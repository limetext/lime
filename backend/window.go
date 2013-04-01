package backend

type Window struct {
	HasId
	HasSettings
	views []View
}

func (w *Window) NewView() *View {
	w.views = append(w.views, View{window: w})
	v := &w.views[len(w.views)-1]
	v.setBuffer(&Buffer{})
	v.selection.regions = []Region{{0, 0}}
	OnNew.Call(v)
	return v
}
