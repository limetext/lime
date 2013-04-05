package backend

import (
	"code.google.com/p/log4go"
	"io/ioutil"
	"lime/backend/loaders"
	. "lime/backend/primitives"
	"runtime"
)

type (
	Editor struct {
		HasSettings
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
	c.Insert(e, c.Buffer().Size(), log4go.FormatLogRecord(log4go.FORMAT_DEFAULT, rec))
	c.EndEdit(e)
}

func (m *myLogWriter) Close() {
}

var ed *Editor

func GetEditor() *Editor {
	if ed == nil {
		ed = &Editor{
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
		ed.Settings() // Just to initialize it
		log4go.Global.Close()
		log4go.Global.AddFilter("console", log4go.DEBUG, &myLogWriter{})
		ed.loadKeybindings()
		ed.loadSettings()
	}
	return ed
}

func (e *Editor) loadKeybinding(fn string) {
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
func (e *Editor) loadKeybindings() {
	// TODO(q): should search for keybindings
	e.loadKeybinding("../../backend/packages/Default/Default.sublime-keymap")
}

func (e *Editor) loadSetting(fn string) {
	if d, err := ioutil.ReadFile(fn); err != nil {
		log4go.Error("Couldn't load file %s: %s", fn, err)
	} else {
		var settings settingsMap
		if err := loaders.LoadJSON(d, &settings); err != nil {
			log4go.Error(err)
		} else {
			log4go.Info("Loaded %s", fn)
		}
		e.settings.merge(settings)
	}
}

func (e *Editor) loadSettings() {
	// TODO(q): should search for settings
	e.loadSetting("../../backend/packages/Default/Default.sublime-settings")
}

func (e *Editor) Console() *View {
	return e.console
}

func (e *Editor) Windows() []*Window {
	return e.windows
}

func (e *Editor) ActiveWindow() *Window {
	return e.activeWindow
}

func (e *Editor) NewWindow() *Window {
	e.windows = append(e.windows, &Window{})
	w := e.windows[len(e.windows)-1]
	w.Settings().Parent = e
	e.activeWindow = w
	return w
}

func (e *Editor) Arch() string {
	return runtime.GOARCH
}

func (e *Editor) Platform() string {
	return runtime.GOOS
}

func (e *Editor) Version() string {
	return "0"
}

func (e *Editor) CommandHandler() CommandHandler {
	return &e.cmdhandler
}

func (e *Editor) HandleInput(kp KeyPress) {
	if e.loginput {
		log4go.Debug("Key: %v", kp)
	}
	// TODO: multi-stage keys
	possible_actions := e.keyBindings.Filter(kp)
	if possible_actions.Len() == 1 {
		action := possible_actions.Bindings[0]
		// TODO: context
		// TODO: what's the command precedence?
		if err := e.CommandHandler().RunTextCommand(e.ActiveWindow().ActiveView(), action.Command, action.Args); err != nil {
			log4go.Debug("Couldn't run textcommand: %s", err)
			if err := e.CommandHandler().RunWindowCommand(e.ActiveWindow(), action.Command, action.Args); err != nil {
				log4go.Debug("Couldn't run windowcommand: %s", err)
				if err := e.CommandHandler().RunApplicationCommand(action.Command, action.Args); err != nil {
					log4go.Debug("Couldn't run applicationcommand: %s", err)
				}
			}
		}
	} else if possible_actions.Len() == 0 && possible_actions.keyOff == 1 && (!kp.Ctrl && !kp.Alt && !kp.Super) {
		// presume insert
		if err := e.CommandHandler().RunTextCommand(e.ActiveWindow().ActiveView(), "insert", Args{"characters": string(kp.Key)}); err != nil {
			log4go.Debug("Couldn't run textcommand: %s", err)
		}
	}
}

func (e *Editor) LogInput(l bool) {
	e.loginput = l
}

func (e *Editor) LogCommands(bool) {
	e.cmdhandler.log = true
}

func (e *Editor) StatusMessage(msg string) {
	log4go.Info(msg)
}

func (e *Editor) ErrorMessage(msg string) {
	log4go.Error(msg)
}

// TODO(q): Actually show a dialog
func (e *Editor) MessageDialog(msg string) {
	log4go.Info(msg)
}

func (e *Editor) RunCommand(name string, args Args) {
	e.CommandHandler().RunApplicationCommand(name, args)
}
