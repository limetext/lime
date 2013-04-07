package sublime

import (
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/primitives"
)

var (
	_ = backend.View{}
	_ = primitives.Region{}
	_ = fmt.Errorf
)

func (o *Window) Py_active_view() (py.Object, error) {
	ret0 := backend.GetEditor().Frontend().ActiveView(o.data)
	var err error
	var pyret0 py.Object

	pyret0, err = _viewClass.Alloc(1)
	if err != nil {
	} else if v2, ok := pyret0.(*View); !ok {
		return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
	} else {
		v2.data = ret0
	}
	if err != nil {
		return nil, err
	}
	return pyret0, err
}
