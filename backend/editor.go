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
		windows                   []*Window
		loginput                  bool
		cmdhandler                commandHandler
		keyBindings, lastBindings KeyBindings
		console                   *View
		frontend                  Frontend
		clipboard                 string
	}
	Frontend interface {
		ActiveWindow() *Window
		ActiveView(*Window) *View
		VisibleRegion(v *View) Region
		Show(v *View, r Region)
		StatusMessage(string)
		ErrorMessage(string)
		MessageDialog(string)
		OkCancelDialog(msg string, okname string)
	}
	myLogWriter struct {
	}
	DummyFrontend struct{}
)

func (h *DummyFrontend) ActiveWindow() *Window {
	if w := GetEditor().Windows(); len(w) > 0 {
		return w[0]
	}
	return nil
}

func (h *DummyFrontend) ActiveView(w *Window) *View {
	if v := w.Views(); len(v) > 0 {
		return v[0]
	}
	return nil
}
func (h *DummyFrontend) StatusMessage(msg string)      {}
func (h *DummyFrontend) ErrorMessage(msg string)       {}
func (h *DummyFrontend) MessageDialog(msg string)      {}
func (h *DummyFrontend) OkCancelDialog(string, string) {}
func (h *DummyFrontend) Show(v *View, r Region)        {}
func (h *DummyFrontend) VisibleRegion(v *View) Region  { return Region{} }

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
				verbose:             true,
			},
			frontend: &DummyFrontend{},
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
		//		initBasicCommands()
	}
	return ed
}

func (e *Editor) Frontend() Frontend {
	return e.frontend
}

func (e *Editor) SetFrontend(f Frontend) {
	e.frontend = f
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
		e.keyBindings.merge(&bindings)
	}
}
func (e *Editor) loadKeybindings() {
	// TODO(q): should search for keybindings
	e.loadKeybinding("../../backend/packages/Default/Default.sublime-keymap")
	e.loadKeybinding("../../3rdparty/bundles/Vintageous/Default.sublime-keymap")
}

func (e *Editor) loadSetting(fn string) {
	if d, err := ioutil.ReadFile(fn); err != nil {
		log4go.Error("Couldn't load file %s: %s", fn, err)
	} else {
		if err := loaders.LoadJSON(d, e.Settings()); err != nil {
			log4go.Error(err)
		} else {
			log4go.Info("Loaded %s", fn)
		}
	}
}

func (e *Editor) loadSettings() {
	// TODO(q): should search for settings
	e.loadSetting("../../backend/packages/Default/Default.sublime-settings")
}

func (e *Editor) PackagesPath() string {
	return "../../3rdparty/bundles/"
}

func (e *Editor) Console() *View {
	return e.console
}

func (e *Editor) Windows() []*Window {
	return e.windows
}

func (e *Editor) NewWindow() *Window {
	e.windows = append(e.windows, &Window{})
	w := e.windows[len(e.windows)-1]
	w.Settings().SetParent(e)
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
		log4go.Info("Key: %v", kp)
	}
	if e.lastBindings.keyOff == 0 {
		e.lastBindings = e.keyBindings
	}
try_again:
	possible_actions := e.lastBindings.Filter(kp)
	e.lastBindings = possible_actions

	// TODO?
	wnd := e.Frontend().ActiveWindow()
	v := e.Frontend().ActiveView(wnd)

	if action := possible_actions.Action(v); action != nil {
		// TODO: what's the command precedence?
		if c := e.cmdhandler.TextCommands[action.Command]; c != nil {
			if err := e.CommandHandler().RunTextCommand(v, action.Command, action.Args); err != nil {
				log4go.Debug("Couldn't run textcommand: %s", err)
			}
		} else if c := e.cmdhandler.WindowCommands[action.Command]; c != nil {
			if err := e.CommandHandler().RunWindowCommand(wnd, action.Command, action.Args); err != nil {
				log4go.Debug("Couldn't run windowcommand: %s", err)
			}
		} else if err := e.CommandHandler().RunApplicationCommand(action.Command, action.Args); err != nil {
			log4go.Debug("Couldn't run applicationcommand: %s", err)
		}
	} else if possible_actions.keyOff > 1 {
		e.lastBindings = e.keyBindings
		goto try_again
	} else if kp.IsCharacter() {
		log4go.Finest("kp: %v, pos: %v", kp, possible_actions)
		if err := e.CommandHandler().RunTextCommand(v, "insert", Args{"characters": string(rune(kp.Key))}); err != nil {
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

func (e *Editor) RunCommand(name string, args Args) {
	e.CommandHandler().RunApplicationCommand(name, args)
}

func (e *Editor) SetClipboard(n string) {
	e.clipboard = n
}

func (e *Editor) GetClipboard() string {
	return e.clipboard
}
