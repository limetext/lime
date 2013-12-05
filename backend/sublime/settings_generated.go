// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file was generated as part of a build step and shouldn't be manually modified
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

var _settingsClass = py.Class{
	Name:    "sublime.Settings",
	Pointer: (*Settings)(nil),
}

type Settings struct {
	py.BaseObject
	data *text.Settings
}

func (o *Settings) PyInit(args *py.Tuple, kwds *py.Dict) error {
	return fmt.Errorf("Can't initialize type Settings")
}
func (o *Settings) Py_add_on_change(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
		arg2 text.OnChangeCallback
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for text.Settings.AddOnChange() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	if v, err := tu.GetItem(1); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(text.OnChangeCallback); !ok {
				return nil, fmt.Errorf("Expected type text.OnChangeCallback for text.Settings.AddOnChange() arg2, not %s", v.Type())
			} else {
				arg2 = v2
			}
		}
	}
	o.data.AddOnChange(arg1, arg2)
	return toPython(nil)
}

func (o *Settings) Py_clear_on_change(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for text.Settings.ClearOnChange() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	o.data.ClearOnChange(arg1)
	return toPython(nil)
}

func (o *Settings) Py_erase(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for text.Settings.Erase() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	o.data.Erase(arg1)
	return toPython(nil)
}

func (o *Settings) Py_has(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for text.Settings.Has() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Has(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Settings) Py_id() (py.Object, error) {
	ret0 := o.data.Id()
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}
