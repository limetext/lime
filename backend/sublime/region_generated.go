// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file was generated as part of a build step and shouldn't be manually modified

package sublime

import (
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	"github.com/quarnster/util/text"
)

var (
	_ = backend.View{}
	_ = text.Region{}
	_ = fmt.Errorf
)

var _regionClass = py.Class{
	Name:    "sublime.Region",
	Pointer: (*Region)(nil),
}

type Region struct {
	py.BaseObject
	data text.Region
}

func (o *Region) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() > 2 {
		return fmt.Errorf("Expected at most 2 arguments")
	}

	if args.Size() > 0 {
		if v, err := args.GetItem(0); err != nil {
			return err
		} else {
			if v3, err2 := fromPython(v); err2 != nil {
				return err2
			} else {
				if v2, ok := v3.(int); !ok {
					return fmt.Errorf("Expected type int for Region.A, not %s", v.Type())
				} else {
					o.data.A = v2
				}
			}
		}
	}

	if args.Size() > 1 {
		if v, err := args.GetItem(1); err != nil {
			return err
		} else {
			if v3, err2 := fromPython(v); err2 != nil {
				return err2
			} else {
				if v2, ok := v3.(int); !ok {
					return fmt.Errorf("Expected type int for Region.B, not %s", v.Type())
				} else {
					o.data.B = v2
				}
			}
		}
	}

	return nil
}
func (o *Region) Py_begin() (py.Object, error) {
	ret0 := o.data.Begin()
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_clip(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 text.Region
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(text.Region); !ok {
				return nil, fmt.Errorf("Expected type text.Region for text.Region.Clip() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Clip(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_contains(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 int
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(int); !ok {
				return nil, fmt.Errorf("Expected type int for text.Region.Contains() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Contains(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_cover(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 text.Region
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(text.Region); !ok {
				return nil, fmt.Errorf("Expected type text.Region for text.Region.Cover() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Cover(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_covers(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 text.Region
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(text.Region); !ok {
				return nil, fmt.Errorf("Expected type text.Region for text.Region.Covers() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Covers(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_empty() (py.Object, error) {
	ret0 := o.data.Empty()
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_end() (py.Object, error) {
	ret0 := o.data.End()
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_intersection(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 text.Region
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(text.Region); !ok {
				return nil, fmt.Errorf("Expected type text.Region for text.Region.Intersection() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Intersection(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_intersects(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 text.Region
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(text.Region); !ok {
				return nil, fmt.Errorf("Expected type text.Region for text.Region.Intersects() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	ret0 := o.data.Intersects(arg1)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) Py_size() (py.Object, error) {
	ret0 := o.data.Size()
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func (o *Region) PyStr() string {
	return o.data.String()
}

func (o *Region) PyGet_a() (py.Object, error) {
	ret := o.data.A
	var pyret py.Object
	var err error

	pyret, err = toPython(ret)
	return pyret, err

}

func (o *Region) PySet_a(v py.Object) error {
	if v3, err2 := fromPython(v); err2 != nil {
		return err2
	} else {
		if v2, ok := v3.(int); !ok {
			return fmt.Errorf("Expected type int for Region.A, not %s", v.Type())
		} else {
			o.data.A = v2
		}
	}
	return nil
}

func (o *Region) PyGet_b() (py.Object, error) {
	ret := o.data.B
	var pyret py.Object
	var err error

	pyret, err = toPython(ret)
	return pyret, err

}

func (o *Region) PySet_b(v py.Object) error {
	if v3, err2 := fromPython(v); err2 != nil {
		return err2
	} else {
		if v2, ok := v3.(int); !ok {
			return fmt.Errorf("Expected type int for Region.B, not %s", v.Type())
		} else {
			o.data.B = v2
		}
	}
	return nil
}
