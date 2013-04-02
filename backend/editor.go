package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	"io/ioutil"
	"lime/backend/loaders"
	. "lime/backend/primitives"
	"runtime"
)

type (
	EditorInfo struct {
		Arch     string
		Platform string
		Version  string
	}
	Editor interface {
		NewWindow() *Window
		Windows() []*Window
		ActiveWindow() *Window
		Info() EditorInfo
		LogInput(bool)
		LogCommands(bool)
		CommandHandler() CommandHandler
		HandleInput(KeyPress)
		Console() *View
	}
	editor struct {
		windows      []*Window
		activeWindow *Window
		loginput     bool
		cmdhandler   commandHandler
		keyBindings  KeyBindings
		console      *View
	}
)

type myLogWriter struct {
}

func (m *myLogWriter) LogWrite(rec *log4go.LogRecord) {
	c := GetEditor().Console()
	e := c.BeginEdit()
	c.Insert(e, c.Size(), fmt.Sprintf("%s: %s\n", rec.Level, rec.Message))
	c.EndEdit(e)
}

func (m *myLogWriter) Close() {
}

var ed *editor

func GetEditor() Editor {
	if ed == nil {
		ed = &editor{
			cmdhandler: commandHandler{
				ApplicationCommands: make(appcmd),
				TextCommands:        make(textcmd),
				WindowCommands:      make(wndcmd),
			},
			console: &View{
				buffer:  &Buffer{},
				scratch: true,
			},
		}
		log4go.Global.Close()
		log4go.Global.AddFilter("console", log4go.DEBUG, &myLogWriter{})
		ed.loadKeybindings()
	}
	return ed
}

func (e *editor) loadKeybindings() {
	fn := "../../backend/json/testdata/Default (OSX).sublime-keymap"
	if d, err := ioutil.ReadFile(fn); err != nil {
		log4go.Error("Couldn't load file %s: %s", fn, err)
	} else {
		var bindings KeyBindings
		if err := loaders.LoadJSON(d, &bindings); err != nil {
			log4go.Error(err)
		} else {
			log4go.Info("Loaded %s", fn)
		}
		e.keyBindings.Merge(&bindings)
	}
}

func (e *editor) Console() *View {
	return e.console
}

func (e *editor) Windows() []*Window {
	return e.windows
}

func (e *editor) ActiveWindow() *Window {
	return e.activeWindow
}

func (e *editor) NewWindow() *Window {
	e.windows = append(e.windows, &Window{})
	w := e.windows[len(e.windows)-1]
	e.activeWindow = w
	return w
}

func (e *editor) Info() EditorInfo {
	return EditorInfo{
		runtime.GOARCH,
		runtime.GOOS,
		"0",
	}
}

func (e *editor) CommandHandler() CommandHandler {
	return &e.cmdhandler
}

func (e *editor) HandleInput(kp KeyPress) {
	if e.loginput {
		log4go.Debug("Key: %v", kp)
	}
	// TODO: multi-stage keys
	possible_actions := e.keyBindings.Filter(kp)
	if possible_actions.Len() == 1 {
		action := possible_actions.Bindings[0]
		// TODO: context
		e.CommandHandler().RunTextCommand(nil, action.Command, action.Args)
	}
}

func (e *editor) LogInput(l bool) {
	e.loginput = l
}
func (e *editor) LogCommands(bool) {
	e.cmdhandler.log = true
}
