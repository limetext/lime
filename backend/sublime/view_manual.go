package sublime

import (
	"lime/3rdparty/libs/gopy/lib"
)

func (v *View) Py_has_non_empty_selection_region() (py.Object, error) {
	for _, r := range v.data.Sel().Regions() {
		if !r.Empty() {
			return py.True, nil
		}
	}
	return py.False, nil
}
