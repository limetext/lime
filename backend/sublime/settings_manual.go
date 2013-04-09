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
		if v2, ok := v.(*py.Unicode); !ok {
			return nil, fmt.Errorf("Expected type *py.Unicode for backend.Settings.Get() arg1, not %s", v.Type())
		} else {
			arg1 = v2.String()
		}
	}
	if v, err := tu.GetItem(1); err == nil {
		def = v
	}
	r := o.data.Get(arg1, def)
	if r == def && r != nil {
		def.Incref()
		return def, nil
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
		if v2, ok := v.(*py.Unicode); !ok {
			return nil, fmt.Errorf("Expected type *py.Unicode for backend.Settings.Set() arg1, not %s", v.Type())
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
	return toPython(nil)
}
