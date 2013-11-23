// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/quarnster/util/text"
	"io/ioutil"
	"lime/backend/loaders"
	. "lime/backend/util"
	"runtime"
	"runtime/debug"
	"sync"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	GetEditor()
}

type (
	Editor struct {
		HasSettings
		windows       []*Window
		active_window *Window
		loginput      bool
		cmdhandler    commandHandler
		keyBindings   KeyBindings
		console       *View
		frontend      Frontend
		clipboard     string
		keyInput      chan (KeyPress)
	}
	Frontend interface {
		VisibleRegion(v *View) Region
		Show(v *View, r Region)
		StatusMessage(string)
		ErrorMessage(string)
		MessageDialog(string)
		OkCancelDialog(msg string, okname string)
	}
	myLogWriter struct {
		log chan string
	}
	DummyFrontend struct{}
)

func (h *DummyFrontend) StatusMessage(msg string)      {}
func (h *DummyFrontend) ErrorMessage(msg string)       {}
func (h *DummyFrontend) MessageDialog(msg string)      {}
func (h *DummyFrontend) OkCancelDialog(string, string) {}
func (h *DummyFrontend) Show(v *View, r Region)        {}
func (h *DummyFrontend) VisibleRegion(v *View) Region  { return Region{} }

func newMyLogWriter() *myLogWriter {
	ret := &myLogWriter{make(chan string, 100)}
	go ret.handle()
	return ret
}

func (m *myLogWriter) handle() {
	for fl := range m.log {
		c := GetEditor().Console()
		f := fmt.Sprintf("%08d %d %s", c.Buffer().Size(), len(fl), fl)
		e := c.BeginEdit()
		c.Insert(e, c.Buffer().Size(), f)
		c.EndEdit(e)
	}
}

func (m *myLogWriter) LogWrite(rec *log4go.LogRecord) {
	p := Prof.Enter("log")
	defer p.Exit()
	fl := log4go.FormatLogRecord(log4go.FORMAT_DEFAULT, rec)
	m.log <- fl
}

func (m *myLogWriter) Close() {
	fmt.Println("Closing...")
	close(m.log)
}

var (
	ed  *Editor
	edl sync.Mutex
)

func GetEditor() *Editor {
	edl.Lock()
	defer edl.Unlock()
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
				buffer:  NewBuffer(),
				scratch: true,
			},
			keyInput: make(chan KeyPress, 32),
		}
		ed.console.Settings().Set("is_widget", true)
		ed.Settings() // Just to initialize it
		log4go.Global.Close()
		log4go.Global.AddFilter("console", log4go.DEBUG, newMyLogWriter())
		go ed.inputthread()
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

func (e *Editor) Init() {
	ed.loadKeybindings()
	ed.loadSettings()
}

func (e *Editor) loadKeybinding(fn string) {
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		log4go.Error("Couldn't load file %s: %s", fn, err)
	}
	var bindings KeyBindings
	if err := loaders.LoadJSON(d, &bindings); err != nil {
		log4go.Error(err)
	} else {
		log4go.Info("Loaded %s", fn)
	}
	e.keyBindings.merge(&bindings)
}

func (e *Editor) loadKeybindings() {
	// TODO(q): should search for keybindings
	e.loadKeybinding("../../backend/packages/Default/Default.sublime-keymap")
	e.loadKeybinding("../../3rdparty/bundles/Vintageous/Default.sublime-keymap")
}

func (e *Editor) loadSetting(fn string) {
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		log4go.Error("Couldn't load file %s: %s", fn, err)
	}
	if err := loaders.LoadJSON(d, e.Settings()); err != nil {
		log4go.Error(err)
	} else {
		log4go.Info("Loaded %s", fn)
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
	edl.Lock()
	defer edl.Unlock()
	ret := make([]*Window, 0, len(e.windows))
	copy(ret, e.windows)
	return ret
}

func (e *Editor) SetActiveWindow(w *Window) {
	e.active_window = w
}

func (e *Editor) ActiveWindow() *Window {
	return e.active_window
}

func (e *Editor) NewWindow() *Window {
	edl.Lock()
	e.windows = append(e.windows, &Window{})
	w := e.windows[len(e.windows)-1]
	edl.Unlock()
	w.Settings().SetParent(e)
	ed.SetActiveWindow(w)
	OnNewWindow.Call(w)
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
	e.keyInput <- kp
}

func (e *Editor) inputthread() {
	pc := 0
	var lastBindings KeyBindings
	doinput := func(kp KeyPress) {
		defer func() {
			if r := recover(); r != nil {
				log4go.Error("Panic in inputthread: %v\n%s", r, string(debug.Stack()))
				if pc > 0 {
					panic(r)
				}
				pc++
			}
		}()
		p := Prof.Enter("hi")
		defer p.Exit()

		lvl := log4go.FINE
		if e.loginput {
			lvl++
		}
		log4go.Logf(lvl, "Key: %v", kp)
		if lastBindings.keyOff == 0 {
			lastBindings = e.keyBindings
		}
	try_again:
		possible_actions := lastBindings.Filter(kp)
		lastBindings = possible_actions

		// TODO?
		var (
			wnd *Window
			v   *View
		)
		if wnd = e.ActiveWindow(); wnd != nil {
			v = wnd.ActiveView()
		}

		if action := possible_actions.Action(v); action != nil {
			p2 := Prof.Enter("hi.perform")
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
			p2.Exit()
		} else if possible_actions.keyOff > 1 {
			lastBindings = e.keyBindings
			goto try_again
		} else if kp.IsCharacter() {
			p2 := Prof.Enter("hi.character")
			log4go.Finest("kp: %v, pos: %v", kp, possible_actions)
			if err := e.CommandHandler().RunTextCommand(v, "insert", Args{"characters": string(rune(kp.Key))}); err != nil {
				log4go.Debug("Couldn't run textcommand: %s", err)
			}
			p2.Exit()
		}
	}
	for kp := range e.keyInput {
		doinput(kp)
	}
}

func (e *Editor) LogInput(l bool) {
	e.loginput = l
}

func (e *Editor) LogCommands(l bool) {
	e.cmdhandler.log = l
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
