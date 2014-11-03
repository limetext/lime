// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/packages"
	. "github.com/limetext/lime/backend/util"
	"github.com/limetext/lime/backend/watch"
	. "github.com/limetext/text"
	"path"
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
		windows      []*Window
		activeWindow *Window
		logInput     bool
		cmdHandler   commandHandler
		keyBindings  keys.KeyBindings
		console      *View
		frontend     Frontend
		keyInput     chan (keys.KeyPress)
		*watch.Watcher
		clipboardSetter func(string) error
		clipboardGetter func() (string, error)
		clipboard       string
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

	DummyFrontend struct {
		m sync.Mutex
		// Default return value for OkCancelDialog
		defaultAction bool
	}
)

var (
	LIME_USER_PACKAGES_PATH = path.Join("..", "..", "3rdparty", "bundles")
	LIME_USER_PACKETS_PATH  = path.Join("..", "..", "3rdparty", "bundles", "User")
	LIME_PACKAGES_PATH      = path.Join("..", "..", "packages")
	LIME_DEFAULTS_PATH      = path.Join("..", "..", "packages", "Default")

	// All user individual settings, snippets, etc.
	// will be in here for later loading
	packets packages.Packets
)

func (h *DummyFrontend) SetDefaultAction(action bool) {
	h.m.Lock()
	defer h.m.Unlock()
	h.defaultAction = action
}
func (h *DummyFrontend) StatusMessage(msg string) { log.Info(msg) }
func (h *DummyFrontend) ErrorMessage(msg string)  { log.Error(msg) }
func (h *DummyFrontend) MessageDialog(msg string) { log.Info(msg) }
func (h *DummyFrontend) OkCancelDialog(msg string, button string) bool {
	log.Info(msg)
	h.m.Lock()
	defer h.m.Unlock()
	return h.defaultAction
}
func (h *DummyFrontend) Show(v *View, r Region)       {}
func (h *DummyFrontend) VisibleRegion(v *View) Region { return Region{} }

var (
	ed  *Editor
	edl sync.Mutex
)

func GetEditor() *Editor {
	edl.Lock()
	defer edl.Unlock()
	if ed == nil {
		ed = &Editor{
			cmdHandler: commandHandler{
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
			keyInput: make(chan keys.KeyPress, 32),
		}
		var err error
		if ed.Watcher, err = watch.NewWatcher(); err != nil {
			log.Error("Couldn't create watcher: %s", err)
		}
		ed.console.Settings().Set("is_widget", true)
		ed.Settings() // Just to initialize it
		log.AddFilter("console", log.DEBUG, log.NewLogWriter(ed.handleLog))
		go ed.inputthread()
		go ed.Observe()
	}
	return ed
}

func (e *Editor) Frontend() Frontend {
	return e.frontend
}

func (e *Editor) SetFrontend(f Frontend) {
	e.frontend = f
}

func setClipboard(n string) error {
	return clipboard.WriteAll(n)
}

func getClipboard() (string, error) {
	return clipboard.ReadAll()
}

func (e *Editor) Init() {
	ed.SetClipboardFuncs(setClipboard, getClipboard)
	ed.loadDefaultPackets()
	ed.loadKeyBindings()
	ed.loadSettings()
}

func (e *Editor) SetClipboardFuncs(setter func(string) error, getter func() (string, error)) {
	e.clipboardSetter = setter
	e.clipboardGetter = getter
}

func (e *Editor) loadDefaultPackets() {
	paths := []string{
		LIME_DEFAULTS_PATH,
		LIME_USER_PACKETS_PATH,
	}
	for _, path := range paths {
		for _, p := range packages.ScanPackets(path) {
			packets = append(packets, p)
		}
	}
}

func (e *Editor) loadKeyBinding(pkg *packages.Packet) {
	if err := pkg.Load(); err != nil {
		log.Error(err)
	} else {
		log.Info("Loaded %s", pkg.Name())
		e.Watch(pkg)
	}
	e.keyBindings.Merge(pkg.MarshalTo().(*keys.KeyBindings))
}

func (e *Editor) loadKeyBindings() {
	for _, p := range packets.Filter("keymap") {
		e.loadKeyBinding(p)
	}
}

func (e *Editor) loadSetting(pkg *packages.Packet) {
	if err := pkg.Load(); err != nil {
		log.Error(err)
	} else {
		log.Info("Loaded %s", pkg.Name())
		e.Watch(pkg)
	}
}

func (e *Editor) loadSettings() {
	defSettings, platSettings := &HasSettings{}, &HasSettings{}
	platSettings.Settings().SetParent(defSettings)
	ed.Settings().SetParent(platSettings)

	p := path.Join(LIME_DEFAULTS_PATH, "Preferences.sublime-settings")
	defPckt := packages.NewPacket(p, defSettings.Settings())
	e.loadSetting(defPckt)

	p = path.Join(LIME_DEFAULTS_PATH, "Preferences ("+e.plat()+").sublime-settings")
	platPckt := packages.NewPacket(p, platSettings.Settings())
	e.loadSetting(platPckt)

	p = path.Join(LIME_USER_PACKETS_PATH, "Preferences.sublime-settings")
	userPckt := packages.NewPacket(p, e.Settings())
	e.loadSetting(userPckt)
}

func (e *Editor) PackagesPath() string {
	return LIME_USER_PACKAGES_PATH
}

func (e *Editor) Console() *View {
	return e.console
}

func (e *Editor) Windows() []*Window {
	edl.Lock()
	defer edl.Unlock()
	ret := make([]*Window, len(e.windows))
	copy(ret, e.windows)
	return ret
}

func (e *Editor) SetActiveWindow(w *Window) {
	e.activeWindow = w
}

func (e *Editor) ActiveWindow() *Window {
	return e.activeWindow
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

func (e *Editor) remove(w *Window) {
	edl.Lock()
	defer edl.Unlock()
	for i, ww := range e.windows {
		if w == ww {
			end := len(e.windows) - 1
			if i != end {
				copy(e.windows[i:], e.windows[i+1:])
			}
			e.windows = e.windows[:end]
			return
		}
	}
	log.Error("Wanted to remove window %+v, but it doesn't appear to be a child of this editor", w)
}

func (e *Editor) Arch() string {
	return runtime.GOARCH
}

func (e *Editor) Platform() string {
	return runtime.GOOS
}

func (e *Editor) plat() string {
	switch e.Platform() {
	case "windows":
		return "Windows"
	case "darwin":
		return "OSX"
	}
	return "Linux"
}

func (e *Editor) Version() string {
	return "0"
}

func (e *Editor) CommandHandler() CommandHandler {
	return &e.cmdHandler
}

func (e *Editor) HandleInput(kp keys.KeyPress) {
	e.keyInput <- kp
}

func (e *Editor) inputthread() {
	pc := 0
	var lastBindings keys.KeyBindings
	doinput := func(kp keys.KeyPress) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic in inputthread: %v\n%s", r, string(debug.Stack()))
				if pc > 0 {
					panic(r)
				}
				pc++
			}
		}()
		p := Prof.Enter("hi")
		defer p.Exit()

		lvl := log.FINE
		if e.logInput {
			lvl++
		}
		log.Logf(lvl, "Key: %v", kp)
		if lastBindings.KeyOff() == 0 {
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

		qc := func(key string, operator Op, operand interface{}, match_all bool) bool {
			return OnQueryContext.Call(v, key, operator, operand, match_all) == True
		}

		if action := possible_actions.Action(qc); action != nil {
			p2 := Prof.Enter("hi.perform")
			e.RunCommand(action.Command, action.Args)
			p2.Exit()
		} else if possible_actions.KeyOff() > 1 {
			lastBindings = e.keyBindings
			goto try_again
		} else if kp.IsCharacter() {
			p2 := Prof.Enter("hi.character")
			log.Finest("kp: %v, pos: %v", kp, possible_actions)
			if err := e.CommandHandler().RunTextCommand(v, "insert", Args{"characters": string(rune(kp.Key))}); err != nil {
				log.Debug("Couldn't run textcommand: %s", err)
			}
			p2.Exit()
		}
	}
	for kp := range e.keyInput {
		doinput(kp)
	}
}

func (e *Editor) LogInput(l bool) {
	e.logInput = l
}

func (e *Editor) LogCommands(l bool) {
	e.cmdHandler.log = l
}

func (e *Editor) RunCommand(name string, args Args) {
	// TODO?
	var (
		wnd *Window
		v   *View
	)
	if wnd = e.ActiveWindow(); wnd != nil {
		v = wnd.ActiveView()
	}

	// TODO: what's the command precedence?
	if c := e.cmdHandler.TextCommands[name]; c != nil {
		if err := e.CommandHandler().RunTextCommand(v, name, args); err != nil {
			log.Debug("Couldn't run textcommand: %s", err)
		}
	} else if c := e.cmdHandler.WindowCommands[name]; c != nil {
		if err := e.CommandHandler().RunWindowCommand(wnd, name, args); err != nil {
			log.Debug("Couldn't run windowcommand: %s", err)
		}
	} else if c := e.cmdHandler.ApplicationCommands[name]; c != nil {
		if err := e.CommandHandler().RunApplicationCommand(name, args); err != nil {
			log.Debug("Couldn't run applicationcommand: %s", err)
		}
	} else {
		log.Debug("Couldn't find command to run")
	}
}

func (e *Editor) SetClipboard(n string) {
	if err := e.clipboardSetter(n); err != nil {
		log.Error("Could not set clipboard: %v", err)
	}

	// Keep a local copy in case the system clipboard isn't working
	e.clipboard = n
}

func (e *Editor) GetClipboard() string {
	if n, err := e.clipboardGetter(); err == nil {
		return n
	} else {
		log.Error("Could not get clipboard: %v", err)
	}

	return e.clipboard
}

func (ed *Editor) handleLog(s string) {
	c := ed.Console()
	f := fmt.Sprintf("%08d %d %s", c.Buffer().Size(), len(s), s)
	e := c.BeginEdit()
	c.Insert(e, c.Buffer().Size(), f)
	c.EndEdit(e)
}
