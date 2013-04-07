package sublime

import (
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"reflect"
)

func toPython(r interface{}) (py.Object, error) {
	switch t := r.(type) {
	case bool:
		if t {
			return py.True, nil
		} else {
			return py.False, nil
		}
	case int:
		return py.NewInt(t), nil
	case string:
		return py.NewString(t)
	case float64:
		return py.NewFloat(t)
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
	case *backend.Window:
		pyret0, err := _windowClass.Alloc(1)
		if err != nil {
			return nil, err
		} else if v2, ok := pyret0.(*Window); !ok {
			return nil, fmt.Errorf("Unable to convert return value to the right type?!: %s", pyret0.Type())
		} else {
			v2.data = t
			return v2, nil
		}
	case []interface{}:
		ret, err := py.NewTuple(int64(len(t)))
		if err != nil {
			return nil, err
		}
		for k, v := range t {
			if p, err := toPython(v); err != nil {
				return nil, err
			} else if err := ret.SetItem(int64(k), p); err != nil {
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
				return nil, err
			} else if err := ret.SetItemString(k, p); err != nil {
				return nil, err
			}
		}
		return ret, nil
	case nil:
		return py.None, nil
	default:
		switch t := reflect.ValueOf(r); t.Kind() {
		case reflect.Int:
			return toPython(int(t.Int()))
		case reflect.Slice:
			ret, err := py.NewList(int64(t.Len()))
			if err != nil {
				return nil, err
			}
			for i := 0; i < t.Len(); i++ {
				if p, err := toPython(t.Index(i).Interface()); err != nil {
					return nil, err
				} else if err := ret.SetItem(int64(i), p); err != nil {
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
	case *py.Int:
		return t.Int(), nil
	case *py.Bool:
		return t.Bool(), nil
	case *py.String:
		return t.String(), nil
	case *py.Float:
		return t.Float64(), nil
	case *Edit:
		return t.data, nil
	case *Region:
		return t.data, nil
	case *py.List:
		g := make([]interface{}, t.Size())
		for i, r := range t.Slice() {
			if v, err := fromPython(r); err != nil {
				return nil, err
			} else {
				g[i] = v
			}
		}
		return g, nil
	case *py.Tuple:
		g := make([]interface{}, t.Size())
		for i, r := range t.Slice() {
			if v, err := fromPython(r); err != nil {
				return nil, err
			} else {
				g[i] = v
			}
		}
		return g, nil
	case *py.Dict:
		if ms, err := t.MapString(); err != nil {
			return nil, err
		} else {
			m2 := make(backend.Args)
			for k, v := range ms {
				if v2, err := fromPython(v); err != nil {
					return nil, err
				} else {
					m2[k] = v2
				}
			}
			return m2, nil
		}
	default:
		return nil, fmt.Errorf("Cannot convert type %s from python to go", r.Type())
	}
}

func (o *Settings) Py_get(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
	var (
		arg1 string
		def  py.Object
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*py.String); !ok {
			return nil, fmt.Errorf("Expected type *py.String for backend.Settings.Get() arg1, not %s", v.Type())
		} else {
			arg1 = v2.String()
		}
	}
	if v, err := tu.GetItem(1); err == nil {
		def = v
	}
	r := o.data.Get(arg1, def)
	if r == def && r != nil {
		return def, nil
	} else if r == nil {
		return py.None, nil
	}
	return toPython(r)
}

func (o *Settings) Py_set(tu *py.Tuple, kw *py.Dict) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v2, ok := v.(*py.String); !ok {
			return nil, fmt.Errorf("Expected type *py.String for backend.Settings.Set() arg1, not %s", v.Type())
		} else {
			arg1 = v2.String()
		}
	}
	if v, err := tu.GetItem(1); err != nil {
		return nil, err
	} else {
		if v2, err := fromPython(v); err != nil {
			return nil, err
		} else {
			o.data.Set(arg1, v2)
		}
	}

	return py.None, nil
}
