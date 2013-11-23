// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package sublime

import (
	"fmt"
	"github.com/quarnster/util/text"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
)

var (
	_ = backend.View{}
	_ = text.Region{}
	_ = fmt.Errorf
)

func (o *Window) Py_open_file(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
		arg2 int
	)
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*py.Unicode); !ok {
		return nil, fmt.Errorf("Expected type *py.Unicode for backend.Window.OpenFile() arg1, not %s", v.Type())
	} else {
		arg1 = v2.String()
	}
	if tu.Size() > 1 {
		v, err := tu.GetItem(1)
		if err != nil {
			return nil, err
		}
		if v2, ok := v.(*py.Long); !ok {
			return nil, fmt.Errorf("Expected type *py.Long for backend.Window.OpenFile() arg2, not %s", v.Type())
		} else {
			arg2 = int(v2.Int64())
		}
	}
	ret0 := o.data.OpenFile(arg1, arg2)
	var pyret0 py.Object

	pyret0, err = _viewClass.Alloc(1)
	if err != nil {
		return nil, err
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

func (o *Window) Py_run_command(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
		arg2 backend.Args
	)
	v, err := tu.GetItem(0)
	if err != nil {
		return nil, err
	}
	if v2, ok := v.(*py.Unicode); !ok {
		return nil, fmt.Errorf("Expected type *py.Unicode for backend.Window.RunCommand() arg1, not %s", v.Type())
	} else {
		arg1 = v2.String()
	}
	arg2 = make(backend.Args)
	if v, err := tu.GetItem(1); err == nil {
		v2, ok := v.(*py.Dict)
		if !ok {
			return nil, fmt.Errorf("Expected type *py.Dict for backend.Window.RunCommand() arg2, not %s", v.Type())
		}
		if v, err := fromPython(v2); err != nil {
			return nil, err
		} else {
			arg2 = v.(backend.Args)
		}
	}
	backend.GetEditor().CommandHandler().RunWindowCommand(o.data, arg1, arg2)
	return toPython(nil)
}
