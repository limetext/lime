package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/primitives"
)

var (
	_ = backend.View{}
	_ = primitives.Region{}
)

var (
	_onQueryContextGlueClass = py.Class{
		Name:    "sublime.OnQueryContextGlue",
		Pointer: (*OnQueryContextGlue)(nil),
	}
)

type (
	OnQueryContextGlue struct {
		py.BaseObject
		inner py.Object
	}
)

func (c *OnQueryContextGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 1 {
		return fmt.Errorf("Expected only 1 argument not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return err
	} else {
		c.inner = v
	}
	// TODO: look into ref counting convention
	c.inner.Incref()
	c.Incref()

	backend.OnQueryContext.Add(c.onQueryContext)
	return nil
}

func (c *OnQueryContextGlue) onQueryContext(v *backend.View, key, operator string, operand interface{}, match_all bool) backend.QueryContextReturn {
	if pv, err := toPython(v); err != nil {
		log4go.Error(err)
	} else if pk, err := toPython(key); err != nil {
		log4go.Error(err)
	} else if po, err := toPython(operator); err != nil {
		log4go.Error(err)
	} else if poa, err := toPython(operand); err != nil {
		log4go.Error(err)
	} else if pm, err := toPython(match_all); err != nil {
		log4go.Error(err)
	} else if ret, err := c.inner.Base().CallFunctionObjArgs(pv, pk, po, poa, pm); err != nil {
		log4go.Error(err)
	} else if r2, ok := ret.(*py.Bool); ok {
		if r2.Bool() {
			return backend.True
		} else {
			return backend.False
		}
	} else {
		log4go.Debug("other: %v", ret)
	}
	return backend.Unknown
}
