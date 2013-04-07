package sublime

import (
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/primitives"
)

func (v *View) Py_has_non_empty_selection_region() (py.Object, error) {
	for _, r := range v.data.Sel().Regions() {
		if !r.Empty() {
			return py.True, nil
		}
	}
	py.False.Incref()
	return py.False, nil
}

func (v *View) Py_show(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
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
	backend.GetEditor().Frontend().Show(v.data, arg1)
	py.None.Incref()
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
	py.None.Incref()
	return py.None, nil
}

func (o *View) Py_command_history(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 int
		arg2 bool
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*py.Int); !ok {
			return nil, fmt.Errorf("Expected type *py.Int for backend.View.CommandHistory() arg1, not %s", v.Type())
		} else {
			arg1 = v2.Int()
		}
	}
	if v, err := tu.GetItem(1); err == nil {
		if v2, ok := v.(*py.Bool); !ok {
			return nil, fmt.Errorf("Expected type *py.Bool for backend.View.CommandHistory() arg2, not %s", v.Type())
		} else {
			arg2 = v2.Bool()
		}
	}
	ret0, ret1, ret2 := o.data.CommandHistory(arg1, arg2)
	var err error
	var pyret0 py.Object

	pyret0, err = py.NewString(ret0)
	if err != nil {
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	var pyret1 py.Object

	pyret1, err = toPython(ret1)
	if err != nil {
		pyret0.Decref()
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	var pyret2 py.Object

	pyret2 = py.NewInt(int(ret2))
	if err != nil {
		pyret0.Decref()
		pyret1.Decref()
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	return py.PackTuple(pyret0, pyret1, pyret2)
}

func (o *View) Py_run_command(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
		arg2 backend.Args
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*py.String); !ok {
			return nil, fmt.Errorf("Expected type *py.String for backend.View.RunCommand() arg1, not %s", v.Type())
		} else {
			arg1 = v2.String()
		}
	}
	arg2 = make(backend.Args)
	if v, err := tu.GetItem(1); err == nil {
		if v2, ok := v.(*py.Dict); !ok {
			return nil, fmt.Errorf("Expected type *py.Dict for backend.View.RunCommand() arg2, not %s", v.Type())
		} else {
			if v, err := fromPython(v2); err != nil {
				return nil, err
			} else {
				arg2 = v.(backend.Args)
			}
		}
	}
	backend.GetEditor().CommandHandler().RunTextCommand(o.data, arg1, arg2)
	py.None.Incref()
	return py.None, nil
}

func (o *View) Py_visible_region() (py.Object, error) {
	ret0 := backend.GetEditor().Frontend().VisibleRegion(o.data)
	var err error
	var pyret0 py.Object

	pyret0, err = _regionClass.Alloc(1)
	if err != nil {
	} else if v2, ok := pyret0.(*Region); !ok {
		return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
	} else {
		v2.data = ret0
	}
	if err != nil {
		return nil, err
	}
	return pyret0, err
}
