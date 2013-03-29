package backend

type (
	ViewEventCallback func(v *View)
	ViewEvent         []ViewEventCallback
)

func (ve ViewEvent) Add(cb ViewEventCallback) {
	ve = append(ve, cb)
}

func (ve ViewEvent) Call(v *View) {
	for i := range ve {
		ve[i](v)
	}
}

var (
	OnNew               ViewEvent
	OnLoad              ViewEvent
	OnActivated         ViewEvent
	OnDeactivated       ViewEvent
	OnClose             ViewEvent
	OnPreSave           ViewEvent
	OnPostSave          ViewEvent
	OnModified          ViewEvent
	OnSelectionModified ViewEvent
)
