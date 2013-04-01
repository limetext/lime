package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"runtime"
)

type (
	EditorInfo struct {
		Arch     string
		Platform string
		Version  string
	}
	Editor interface {
		Windows() []*Window
		ActiveWindow() *Window
		Info() EditorInfo
		LogInput(bool)
		LogCommands(bool)
		CommandHandler() CommandHandler
		HandleInput(KeyPress)
	}
	editor struct {
		windows      []*Window
		activeWindow *Window
		loginput     bool
		cmdhandler   commandHandler
		keyBindings  KeyBindings
	}
)

var ed *editor

func GetEditor() Editor {
	if ed == nil {
		ed = &editor{
			cmdhandler: commandHandler{
				ApplicationCommands: make(appcmd),
				TextCommands:        make(textcmd),
				WindowCommands:      make(wndcmd),
			},
		}
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
		if err := LoadJSON(d, &bindings); err != nil {
			log4go.Error(err)
		} else {
			log4go.Info("Loaded %s", fn)
		}
		e.keyBindings.Merge(&bindings)
	}
}

func (e *editor) Windows() []*Window {
	return e.windows
}

func (e *editor) ActiveWindow() *Window {
	return e.activeWindow
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
