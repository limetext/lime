package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"github.com/quarnster/util/text"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/util"

	//	"time"
)

var (
	_ = backend.View{}
	_ = text.Region{}
)

var (
	_windowCommandGlueClass = py.Class{
		Name:    "sublime.WindowCommandGlue",
		Pointer: (*WindowCommandGlue)(nil),
	}
	_textCommandGlueClass = py.Class{
		Name:    "sublime.TextCommandGlue",
		Pointer: (*TextCommandGlue)(nil),
	}
	_applicationCommandGlueClass = py.Class{
		Name:    "sublime.ApplicationCommandGlue",
		Pointer: (*ApplicationCommandGlue)(nil),
	}
)

type (
	CommandGlue struct {
		py.BaseObject
		inner py.Object
	}
	WindowCommandGlue struct {
		py.BaseObject
		CommandGlue
	}
	TextCommandGlue struct {
		py.BaseObject
		CommandGlue
	}
	ApplicationCommandGlue struct {
		py.BaseObject
		CommandGlue
	}
)

func (c *CommandGlue) BypassUndo() bool {
	return false
}

func (c *CommandGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 1 {
		return fmt.Errorf("Expected only 1 argument not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return pyError(err)
	} else {
		c.inner = v
	}
	c.inner.Incref()
	return nil
}

func (c *CommandGlue) CreatePyArgs(args backend.Args) (ret *py.Dict, err error) {
	if r, err := toPython(args); err != nil {
		return nil, err
	} else {
		return r.(*py.Dict), nil
	}
}

func (c *CommandGlue) callBool(name string, args backend.Args) bool {
	gs := py.GilState_Ensure()
	defer gs.Release()

	var (
		pyargs, r py.Object
		err       error
	)
	if pyargs, err = c.CreatePyArgs(args); err != nil {
		log4go.Error(err)
		return false
	}
	defer pyargs.Decref()

	if r, err = c.CallMethodObjArgs(name, pyargs); err != nil {
		log4go.Error(err)
	} else {
		defer r.Decref()
		if r, ok := r.(*py.Bool); ok {
			return r.Bool()
		}
	}
	return true
}

func (c *CommandGlue) IsEnabled(args backend.Args) bool {
	return c.callBool("is_enabled", args)
}

func (c *CommandGlue) IsVisible(args backend.Args) bool {
	return c.callBool("is_visible", args)
}

func (c *CommandGlue) Description(args backend.Args) string {
	gs := py.GilState_Ensure()
	defer gs.Release()

	var (
		pyargs, r py.Object
		err       error
	)
	if pyargs, err = c.CreatePyArgs(args); err != nil {
		log4go.Error(err)
		return ""
	}
	defer pyargs.Decref()

	if r, err = c.CallMethodObjArgs("description", pyargs); err != nil {
		log4go.Error(err)
	} else {
		defer r.Decref()
		if r, ok := r.(*py.Unicode); ok {
			return r.String()
		}
	}
	return ""
}

func pyError(err error) error {
	// if m, err := py.Import("sys"); err == nil {
	// 	defer m.Decref()
	// 	if i, err := m.Dict().GetItemString("last_traceback"); err == nil {
	// 		defer i.Decref()
	// 		log4go.Debug("%v", i)
	// 	}
	// }
	return fmt.Errorf("%v", err)
}
func (c *TextCommandGlue) Run(v *backend.View, e *backend.Edit, args backend.Args) error {
	l := py.NewLock()
	defer l.Unlock()

	p0 := util.Prof.Enter("tc.run")
	defer p0.Exit()
	var (
		pyv, pye, pyargs, obj py.Object
		err                   error
	)
	if pyv, err = toPython(v); err != nil {
		return pyError(err)
	}
	defer pyv.Decref()

	if pye, err = toPython(e); err != nil {
		return pyError(err)
	}
	defer pye.Decref()

	if pyargs, err = c.CreatePyArgs(args); err != nil {
		return pyError(err)
	}
	defer pyargs.Decref()

	init := util.Prof.Enter("tc.init")
	if obj, err = c.inner.Base().CallFunctionObjArgs(pyv); err != nil {
		return pyError(err)
	}
	defer obj.Decref()
	init.Exit()
	// interrupt := true
	// defer func() { interrupt = false }()
	// go func() {
	// 	<-time.After(time.Second * 5)
	// 	if interrupt {
	// 		py.SetInterrupt()
	// 	}
	// }()
	exec := util.Prof.Enter("tc.exec")
	defer exec.Exit()
	if obj.Base().HasAttrString("run_") {
		// The plugin is probably trying to bypass the undostack...
		old := v.IsScratch()
		v.SetScratch(true)
		log4go.Finest("Discarded: %s", e)
		v.EndEdit(e)
		v.SetScratch(old)
		ret, err := obj.Base().CallMethodObjArgs("run_", pye, pyargs)
		if ret != nil {
			ret.Decref()
		}
		if err != nil {
			return pyError(err)
		}
	} else {
		ret, err := obj.Base().CallMethodObjArgs("run__", pye, pyargs)
		if ret != nil {
			ret.Decref()
		}
		if err != nil {
			return pyError(err)
		}
	}
	return nil
}

func (c *WindowCommandGlue) Run(w *backend.Window, args backend.Args) error {
	l := py.NewLock()
	defer l.Unlock()

	var (
		pyw, pyargs, obj py.Object
		err              error
	)
	log4go.Debug("WindowCommand: %v", args)
	if pyw, err = toPython(w); err != nil {
		return pyError(err)
	}
	defer pyw.Decref()

	if pyargs, err = c.CreatePyArgs(args); err != nil {
		return pyError(err)
	}
	defer pyargs.Decref()
	// interrupt := true
	// defer func() { interrupt = false }()
	// go func() {
	// 	<-time.After(time.Second * 5)
	// 	if interrupt {
	// 		py.SetInterrupt()
	// 	}
	// }()

	if obj, err = c.inner.Base().CallFunctionObjArgs(pyw); err != nil {
		return pyError(err)
	}
	defer obj.Decref()
	if ret, err := obj.Base().CallMethodObjArgs("run_", pyargs); err != nil {
		return pyError(err)
	} else {
		ret.Decref()
	}
	return nil
}

func (c *ApplicationCommandGlue) Run(args backend.Args) error {
	l := py.NewLock()
	defer l.Unlock()

	var (
		pyargs py.Object
		err    error
	)
	if pyargs, err = c.CreatePyArgs(args); err != nil {
		return pyError(err)
	}
	defer pyargs.Decref()
	// interrupt := true
	// defer func() { interrupt = false }()
	// go func() {
	// 	<-time.After(time.Second * 5)
	// 	if interrupt {
	// 		py.SetInterrupt()
	// 	}
	// }()

	if obj, err := c.inner.Base().CallFunctionObjArgs(); err != nil {
		return pyError(err)
	} else {
		defer obj.Decref()
		if ret, err := obj.Base().CallMethodObjArgs("run", pyargs); err != nil {
			return pyError(err)
		} else {
			ret.Decref()
		}
	}
	return nil
}

func (c *ApplicationCommandGlue) IsChecked(args backend.Args) bool {
	return c.callBool("is_checked", args)
}
