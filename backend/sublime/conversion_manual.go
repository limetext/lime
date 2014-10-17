// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	"github.com/limetext/text"
	"reflect"
)

type (
	Tuple []interface{}
	List  []interface{}
)

func toPython(r interface{}) (py.Object, error) {
	if v := reflect.ValueOf(r); !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		py.None.Incref()
		return py.None, nil
	}
	switch t := r.(type) {
	case nil:
		py.None.Incref()
		return py.None, nil
	case bool:
		if t {
			py.True.Incref()
			return py.True, nil
		}
		py.False.Incref()
		return py.False, nil
	case int:
		return py.NewLong(int64(t)), nil
	case int64:
		return py.NewLong(t), nil
	case string:
		return py.NewUnicode(t)
	case float64:
		return py.NewFloat(t)
	case *text.Settings:
		pyret0, err := _settingsClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*Settings); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case *backend.Edit:
		pyret0, err := _editClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*Edit); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case *backend.View:
		pyret0, err := _viewClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*View); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case text.Region:
		pyret0, err := _regionClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*Region); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case *text.RegionSet:
		pyret0, err := _region_setClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*RegionSet); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case *backend.Window:
		if t == nil {
			panic(t)
		}
		pyret0, err := _windowClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*Window); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case List:
		ret, err := py.NewList(int64(len(t)))
		if err != nil {
			return nil, err
		}
		for k, v := range t {
			if p, err := toPython(v); err != nil {
				ret.Decref()
				return nil, err
			} else if err := ret.SetItem(int64(k), p); err != nil {
				ret.Decref()
				return nil, err
			}
		}
		return ret, nil
	case Tuple:
		ret, err := py.NewTuple(int64(len(t)))
		if err != nil {
			return nil, err
		}
		for k, v := range t {
			if p, err := toPython(v); err != nil {
				ret.Decref()
				return nil, err
			} else if err := ret.SetItem(int64(k), p); err != nil {
				ret.Decref()
				return nil, err
			}
		}
		return ret, nil
	case backend.Args:
		ret, err := py.NewDict()
		if err != nil {
			return nil, err
		}
		for k, v := range t {
			if p, err := toPython(v); err != nil {
				ret.Decref()
				return nil, err
			} else if err := ret.SetItemString(k, p); err != nil {
				p.Decref()
				ret.Decref()
				return nil, err
			} else {
				p.Decref()
			}
		}
		return ret, nil
	default:
		switch t := reflect.ValueOf(r); t.Kind() {
		case reflect.Int:
			return toPython(t.Int())
		case reflect.Slice:
			ret, err := py.NewList(int64(t.Len()))
			if err != nil {
				return nil, err
			}
			for i := 0; i < t.Len(); i++ {
				if p, err := toPython(t.Index(i).Interface()); err != nil {
					ret.Decref()
					return nil, err
				} else if err := ret.SetItem(int64(i), p); err != nil {
					ret.Decref()
					return nil, err
				}
			}
			return ret, nil
		}
		return nil, fmt.Errorf("Can't return type %v from go to python", reflect.TypeOf(t))
	}
}

func fromPython(r py.Object) (interface{}, error) {
	switch t := r.(type) {
	case *py.NoneObject:
		return nil, nil
	case *py.Long:
		return int(t.Int64()), nil
	case *py.Bool:
		return t.Bool(), nil
	case *py.Unicode:
		return t.String(), nil
	case *py.Float:
		return t.Float64(), nil
	case *Edit:
		return t.data, nil
	case *Region:
		return t.data, nil
	case *py.List:
		g := make(List, t.Size())
		for i, r := range t.Slice() {
			if v, err := fromPython(r); err != nil {
				return nil, err
			} else {
				g[i] = v
			}
		}
		return g, nil
	case *TextCommandGlue:
		t.Incref()
		return t, nil
	case *WindowCommandGlue:
		t.Incref()
		return t, nil
	case *ApplicationCommandGlue:
		t.Incref()
		return t, nil
	case *py.Tuple:
		g := make(Tuple, t.Size())
		for i, r := range t.Slice() {
			if v, err := fromPython(r); err != nil {
				return nil, err
			} else {
				g[i] = v
			}
		}
		return g, nil
	case *py.Dict:
		ms, err := t.MapString()
		if err != nil {
			return nil, err
		}
		m2 := make(backend.Args)
		for k, v := range ms {
			if v2, err := fromPython(v); err != nil {
				return nil, err
			} else {
				m2[k] = v2
			}
		}
		return m2, nil
	default:
		return nil, fmt.Errorf("Cannot convert type %s from python to go", r.Type())
	}
}
