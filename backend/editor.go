// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/limetext/lime/backend/loaders"
	. "github.com/limetext/lime/backend/util"
	. "github.com/quarnster/util/text"
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
		watcher       *fsnotify.Watcher
		watchedFiles  map[string]Watched
	}

	// The Frontend interface defines the API
	// for functionality that is frontend specific.
	Frontend interface {
		// Probe the frontend for the currently
		// visible region of the given view.
		VisibleRegion(v *View) Region

		// Make the frontend show the specified region of the
		// given view.
		Show(v *View, r Region)

		// Sets the status message shown in the status bar
		StatusMessage(string)

		// Displays an error message to the usser
		ErrorMessage(string)

		// Displays a message dialog to the user
		MessageDialog(string)

		// Displays an ok / cancel dialog to the user.
		// "okname" if provided will be used as the text
		// instead of "Ok" for the ok button.
		// Returns true when ok was pressed, and false when
		// cancel was pressed.
		OkCancelDialog(msg string, okname string) bool
	}
	myLogWriter struct {
		log chan string
	}
	DummyFrontend struct {
		// Default return value for OkCancelDialog
		DefaultAction bool
	}
)

const (
	LIME_USER_PACKAGES_PATH = "../../3rdparty/bundles/"
	LIME_USER_PACKETS_PATH  = "../../3rdparty/bundles/User/"
	LIME_DEFAULTS_PATH      = "../../backend/packages/Default/"
)

func (h *DummyFrontend) StatusMessage(msg string) { log4go.Info(msg) }
func (h *DummyFrontend) ErrorMessage(msg string)  { log4go.Error(msg) }
func (h *DummyFrontend) MessageDialog(msg string) { log4go.Info(msg) }
func (h *DummyFrontend) OkCancelDialog(msg string, button string) bool {
	log4go.Info(msg)
	return h.DefaultAction
}
func (h *DummyFrontend) Show(v *View, r Region)       {}
func (h *DummyFrontend) VisibleRegion(v *View) Region { return Region{} }

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
			keyInput:     make(chan KeyPress, 32),
			watcher:      newWatcher(),
			watchedFiles: make(map[string]Watched),
		}
		ed.console.Settings().Set("is_widget", true)
		ed.Settings() // Just to initialize it
		log4go.Global.Close()
		log4go.Global.AddFilter("console", log4go.DEBUG, newMyLogWriter())
		go ed.inputthread()
		go ed.observeFiles()
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

func (e *Editor) loadKeybinding(pkg *packet) {
	var bindings KeyBindings
	if err := loaders.LoadJSON(pkg.Get().([]byte), &bindings); err != nil {
		log4go.Error(err)
	} else {
		log4go.Info("Loaded %s", pkg.Name())
		e.Watch(NewWatchedPackage(pkg))
	}
	e.keyBindings.merge(&bindings)
}

func (e *Editor) loadKeybindings() {
	for _, p := range packets.filter("keymap") {
		e.loadKeybinding(p)
	}
}

func (e *Editor) loadSetting(pkg *packet) {
	if err := loaders.LoadJSON(pkg.Get().([]byte), e.Settings()); err != nil {
		log4go.Error(err)
	} else {
		log4go.Info("Loaded %s", pkg.Name())
		e.Watch(NewWatchedPackage(pkg))
	}
}

func (e *Editor) loadSettings() {
	for _, p := range packets.filter("setting") {
		e.loadSetting(p)
	}
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

func (e *Editor) Watcher() *fsnotify.Watcher {
	return e.watcher
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

func (e *Editor) Watch(file Watched) {
	log4go.Finest("Watch(%v)", file)
	if err := e.watcher.Watch(file.Name()); err != nil {
		log4go.Error("Could not watch file: %v", err)
	} else {
		e.watchedFiles[file.Name()] = file
	}
}

func (e *Editor) UnWatch(name string) {
	if err := e.watcher.RemoveWatch(name); err != nil {
		log4go.Error("Couldn't unwatch file: %v", err)
	}
	log4go.Finest("UnWatch(%s)", name)
	delete(e.watchedFiles, name)
}

func (e *Editor) observeFiles() {
	for {
		select {
		case ev := <-e.watcher.Event:
			if ev.IsModify() {
				if f, exist := e.watchedFiles[ev.Name]; exist {
					f.Reload()
				}
			}
		case err := <-e.watcher.Error:
			log4go.Error("error:", err)
		}
	}
}

func newWatcher() (w *fsnotify.Watcher) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log4go.Error("Could not create watcher due to: %v", err)
	}
	return
}
