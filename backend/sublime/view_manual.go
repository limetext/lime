package sublime

import (
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend/primitives"
)

func (v *View) Py_has_non_empty_selection_region() (py.Object, error) {
	for _, r := range v.data.Sel().Regions() {
		if !r.Empty() {
			return py.True, nil
		}
	}
	return py.False, nil
}

func (v *View) Py_show(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
	//TODO
	return py.None, nil
}

func (o *View) Py_substr(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 primitives.Region
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*Region); !ok {
			if v2, ok := v.(*py.Int); !ok {
				return nil, fmt.Errorf("Expected type *Region or *Int for primitives.Buffer.Substr() arg1, not %s", v.Type())
			} else {
				arg1.A = v2.Int()
				arg1.B = arg1.A + 1
			}
		} else {
			arg1 = v2.data
		}
	}
	ret0 := o.data.Buffer().Substr(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = py.NewString(ret0)
	if err != nil {
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	return pyret0, err
}

func (o *View) Py_add_regions(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
	var (
		arg1 string
		arg2 []primitives.Region
	)
	if tu.Size() < 2 {
		return nil, fmt.Errorf("Not the expected argument size: %d", tu.Size())
	}
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*py.String); !ok {
			return nil, fmt.Errorf("Expected type *py.String for backend.View.AddRegions() arg1, not %s", v.Type())
		} else {
			arg1 = v2.String()
		}
	}

	if v, err := tu.GetItem(1); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*py.List); !ok {
			return nil, fmt.Errorf("Expected type *py.List for backend.View.AddRegions() arg2, not %s", v.Type())
		} else {
			d := v2.Slice()
			arg2 = make([]primitives.Region, len(d))
			for i, o := range d {
				if v, err := fromPython(o); err != nil {
					return nil, err
				} else if v2, ok := v.(primitives.Region); !ok {
					return nil, fmt.Errorf("Expected non-region item in list passed to backend.View.AddRegions(): %s", o.Type())
				} else {
					arg2[i] = v2
				}
			}
		}
	}
	o.data.AddRegions(arg1, arg2)
	return py.None, nil
}
