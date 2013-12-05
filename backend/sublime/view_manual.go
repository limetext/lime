// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"github.com/quarnster/util/text"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
)

var _ = log4go.Error

type Refcounted py.Object

func (v *View) Py_has_non_empty_selection_region() (py.Object, error) {
	for _, r := range v.data.Sel().Regions() {
		if !r.Empty() {
			return toPython(true)
		}
	}
	return toPython(false)
}

func (view *View) Py_show(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
	var (
		arg1 text.Region
	)
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*Region); !ok {
		if v2, ok := v.(*py.Long); !ok {
			return nil, fmt.Errorf("Expected type *Region or *Int for primitives.Buffer.Substr() arg1, not %s", v.Type())
		} else {
			arg1.A = int(v2.Int64())
			arg1.B = arg1.A + 1
		}
	} else {
		arg1 = v2.data
	}
	backend.GetEditor().Frontend().Show(view.data, arg1)
	return toPython(nil)
}

func (o *View) Py_substr(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 text.Region
	)
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*Region); !ok {
		if v2, ok := v.(*py.Long); !ok {
			return nil, fmt.Errorf("Expected type *Region or *Int for primitives.Buffer.Substr() arg1, not %s", v.Type())
		} else {
			arg1.A = int(v2.Int64())
			arg1.B = arg1.A + 1
		}
	} else {
		arg1 = v2.data
	}
	ret0 := o.data.Buffer().Substr(arg1)
	var pyret0 py.Object

	pyret0, err = py.NewUnicode(ret0)
	if err != nil {
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	return pyret0, err
}

func (o *View) Py_add_regions(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
	var (
		arg1 string
		arg2 []text.Region
	)
	if tu.Size() < 2 {
		return nil, fmt.Errorf("Not the expected argument size: %d", tu.Size())
	}
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*py.Unicode); !ok {
		return nil, fmt.Errorf("Expected type *py.Unicode for backend.View.AddRegions() arg1, not %s", v.Type())
	} else {
		arg1 = v2.String()
	}
	v, err = tu.GetItem(1)
	if err != nil {
		return nil, err
	}
	v2, ok := v.(*py.List)
	if !ok {
		return nil, fmt.Errorf("Expected type *py.List for backend.View.AddRegions() arg2, not %s", v.Type())
	}
	d := v2.Slice()
	arg2 = make([]text.Region, len(d))
	for i, o := range d {
		if v, err := fromPython(o); err != nil {
			return nil, err
		} else if v2, ok := v.(text.Region); !ok {
			return nil, fmt.Errorf("Expected non-region item in list passed to backend.View.AddRegions(): %s", o.Type())
		} else {
			arg2[i] = v2
		}
	}
	o.data.AddRegions(arg1, arg2, "", "", 0) // TODO
	return toPython(nil)
}

func (o *View) Py_command_history(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 int
		arg2 bool
	)
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*py.Long); !ok {
		return nil, fmt.Errorf("Expected type *py.Long for backend.View.CommandHistory() arg1, not %s", v.Type())
	} else {
		arg1 = int(v2.Int64())
	}
	if v, err := tu.GetItem(1); err == nil {
		if v2, ok := v.(*py.Bool); !ok {
			return nil, fmt.Errorf("Expected type *py.Bool for backend.View.CommandHistory() arg2, not %s", v.Type())
		} else {
			arg2 = v2.Bool()
		}
	}
	ret0, ret1, ret2 := o.data.CommandHistory(arg1, arg2)
	var pyret0 py.Object

	pyret0, err = py.NewUnicode(ret0)
	if err != nil {
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	var pyret1 py.Object

	// Sublime compatibility
	if len(ret1) == 0 {
		pyret1, err = toPython(nil)
	} else {
		pyret1, err = toPython(ret1)
	}
	if err != nil {
		pyret0.Decref()
		// TODO: do the py objs need to be freed?
		return nil, err
	}

	var pyret2 py.Object

	pyret2 = py.NewLong(int64(ret2))
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
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*py.Unicode); !ok {
		return nil, fmt.Errorf("Expected type *py.Unicode for backend.View.RunCommand() arg1, not %s", v.Type())
	} else {
		arg1 = v2.String()
	}
	arg2 = make(backend.Args)
	if v, err := tu.GetItem(1); err == nil {
		v2, ok := v.(*py.Dict)
		if !ok {
			return nil, fmt.Errorf("Expected type *py.Dict for backend.View.RunCommand() arg2, not %s", v.Type())
		}
		if v, err := fromPython(v2); err != nil {
			return nil, err
		} else {
			arg2 = v.(backend.Args)
		}
	}
	backend.GetEditor().CommandHandler().RunTextCommand(o.data, arg1, arg2)
	return toPython(nil)
}

func (o *View) Py_visible_region() (py.Object, error) {
	ret0 := backend.GetEditor().Frontend().VisibleRegion(o.data)
	var err error
	var pyret0 py.Object

	pyret0, err = _regionClass.Alloc(1)
	if err != nil {
		return nil, err
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

func (o *View) Py_line(tu *py.Tuple) (py.Object, error) {
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if _, ok := v.(*Region); ok {
		return o.liner(tu)
	} else {
		return o.line(tu)
	}
}

func (o *View) Py_full_line(tu *py.Tuple) (py.Object, error) {
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if _, ok := v.(*Region); ok {
		return o.fullliner(tu)
	}
	return o.fullline(tu)
}

func (o *View) Py_word(tu *py.Tuple) (py.Object, error) {
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if _, ok := v.(*Region); ok {
		return o.wordr(tu)
	}
	return o.word(tu)
}

func (o *View) Py_set_syntax_file(tu *py.Tuple) (py.Object, error) {
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v, ok := v.(*py.Unicode); ok {
		o.data.Settings().Set("syntax", v.String())
	}
	return toPython(nil)
}
