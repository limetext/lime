package sublime

import "lime/backend"

type (
	Region backend.Region
	View   struct {
		bv      *backend.View
		scratch bool
	}

	Settings backend.Settings
	Window   backend.Window
)

func (v *View) Size() int {
	return v.bv.Buffer().Size()
}

func (v *View) IsScratch() bool {
	return v.scratch
}

func (v *View) SetScratch(s bool) {
	v.scratch = s
}

func (v *View) Settings() *Settings {
	return (*Settings)(v.bv.Settings())
}

func (v *View) Window() *Window {
	return (*Window)(v.bv.Window())
}

func (v *View) Substr(r Region) string {
	return v.bv.Buffer().Substr(backend.Region(r))
}

func LoadSettings(name string) *Settings {
	return nil
}
