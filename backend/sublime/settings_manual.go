package sublime

import (
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
)

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
	default:
		return nil, fmt.Errorf("Can't return type %s from Settings.Get() to python", t)
	}
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
		switch t := v.(type) {
		case *py.Int:
			o.data.Set(arg1, t.Int())
		case *py.Bool:
			o.data.Set(arg1, t.Bool())
		case *py.String:
			o.data.Set(arg1, t.String())
		case *py.Float:
			o.data.Set(arg1, t.Float64())
		default:
			return nil, fmt.Errorf("Can't set setting \"%s\" with a type of %s", arg1, v.Type())
		}
	}

	return py.None, nil
}
