// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package main

import (
	"flag"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	_ "github.com/limetext/lime/backend/commands"
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/sublime"
	"github.com/limetext/lime/backend/textmate"
	"github.com/limetext/lime/backend/util"
	"github.com/limetext/termbox-go"
	. "github.com/limetext/text"
	"path"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

var (
	lut = map[termbox.Key]keys.KeyPress{
		// Omission of these are intentional due to map collisions
		//		termbox.KeyCtrlTilde:      keys.KeyPress{Ctrl: true, Key: '~'},
		//		termbox.KeyCtrlBackslash:  keys.KeyPress{Ctrl: true, Key: '\\'},
		//		termbox.KeyCtrlSlash:      keys.KeyPress{Ctrl: true, Key: '/'},
		//		termbox.KeyCtrlUnderscore: keys.KeyPress{Ctrl: true, Key: '_'},
		//		termbox.KeyCtrlLsqBracket: keys.KeyPress{Ctrl: true, Key: '{'},
		//		termbox.KeyCtrlRsqBracket: keys.KeyPress{Ctrl: true, Key: '}'},
		// termbox.KeyCtrl3:
		// termbox.KeyCtrl8
		//		termbox.KeyCtrl2:      keys.KeyPress{Ctrl: true, Key: '2'},
		termbox.KeyCtrlSpace:  {Ctrl: true, Key: ' '},
		termbox.KeyCtrlA:      {Ctrl: true, Key: 'a'},
		termbox.KeyCtrlB:      {Ctrl: true, Key: 'b'},
		termbox.KeyCtrlC:      {Ctrl: true, Key: 'c'},
		termbox.KeyCtrlD:      {Ctrl: true, Key: 'd'},
		termbox.KeyCtrlE:      {Ctrl: true, Key: 'e'},
		termbox.KeyCtrlF:      {Ctrl: true, Key: 'f'},
		termbox.KeyCtrlG:      {Ctrl: true, Key: 'g'},
		termbox.KeyCtrlH:      {Ctrl: true, Key: 'h'},
		termbox.KeyCtrlJ:      {Ctrl: true, Key: 'j'},
		termbox.KeyCtrlK:      {Ctrl: true, Key: 'k'},
		termbox.KeyCtrlL:      {Ctrl: true, Key: 'l'},
		termbox.KeyCtrlN:      {Ctrl: true, Key: 'n'},
		termbox.KeyCtrlO:      {Ctrl: true, Key: 'o'},
		termbox.KeyCtrlP:      {Ctrl: true, Key: 'p'},
		termbox.KeyCtrlQ:      {Ctrl: true, Key: 'q'},
		termbox.KeyCtrlR:      {Ctrl: true, Key: 'r'},
		termbox.KeyCtrlS:      {Ctrl: true, Key: 's'},
		termbox.KeyCtrlT:      {Ctrl: true, Key: 't'},
		termbox.KeyCtrlU:      {Ctrl: true, Key: 'u'},
		termbox.KeyCtrlV:      {Ctrl: true, Key: 'v'},
		termbox.KeyCtrlW:      {Ctrl: true, Key: 'w'},
		termbox.KeyCtrlX:      {Ctrl: true, Key: 'x'},
		termbox.KeyCtrlY:      {Ctrl: true, Key: 'y'},
		termbox.KeyCtrlZ:      {Ctrl: true, Key: 'z'},
		termbox.KeyCtrl4:      {Ctrl: true, Key: '4'},
		termbox.KeyCtrl5:      {Ctrl: true, Key: '5'},
		termbox.KeyCtrl6:      {Ctrl: true, Key: '6'},
		termbox.KeyCtrl7:      {Ctrl: true, Key: '7'},
		termbox.KeyEnter:      {Key: keys.Enter},
		termbox.KeySpace:      {Key: ' '},
		termbox.KeyBackspace2: {Key: keys.Backspace},
		termbox.KeyArrowUp:    {Key: keys.Up},
		termbox.KeyArrowDown:  {Key: keys.Down},
		termbox.KeyArrowLeft:  {Key: keys.Left},
		termbox.KeyArrowRight: {Key: keys.Right},
		termbox.KeyDelete:     {Key: keys.Delete},
		termbox.KeyEsc:        {Key: keys.Escape},
		termbox.KeyPgup:       {Key: keys.PageUp},
		termbox.KeyPgdn:       {Key: keys.PageDown},
		termbox.KeyF1:         {Key: keys.F1},
		termbox.KeyF2:         {Key: keys.F2},
		termbox.KeyF3:         {Key: keys.F3},
		termbox.KeyF4:         {Key: keys.F4},
		termbox.KeyF5:         {Key: keys.F5},
		termbox.KeyF6:         {Key: keys.F6},
		termbox.KeyF7:         {Key: keys.F7},
		termbox.KeyF8:         {Key: keys.F8},
		termbox.KeyF9:         {Key: keys.F9},
		termbox.KeyF10:        {Key: keys.F10},
		termbox.KeyF11:        {Key: keys.F11},
		termbox.KeyF12:        {Key: keys.F12},
		termbox.KeyTab:        {Key: '\t'},
	}
	palLut    func(col textmate.Color) termbox.Attribute
	scheme    *textmate.Theme
	defaultBg = termbox.ColorBlack
	defaultFg = termbox.ColorWhite
	blink     bool
)

// Command line flags
var (
	showConsole   = flag.Bool("console", false, "Display console")
	consoleHeight = flag.Int("consoleHeight", 20, "Height of console")
	rotateLog     = flag.Bool("rotateLog", false, "Rotate debug log")
)

const (
	render_chan_len = 2
)

type layout struct {
	x, y          int
	width, height int
	visible       Region
	lastUpdate    int
}

type tbfe struct {
	layout         map[*backend.View]layout
	status_message string
	dorender       chan bool
	shutdown       chan bool
	lock           sync.Mutex
	editor         *backend.Editor
	console        *backend.View
	currentView    *backend.View
	currentWindow  *backend.Window
}

type tbfeBufferDeltaObserver struct {
	t    *tbfe
	view *backend.View
}

// Creates and initializes the frontend.
func createFrontend() *tbfe {
	var t tbfe
	t.dorender = make(chan bool, render_chan_len)
	t.shutdown = make(chan bool, 2)
	t.layout = make(map[*backend.View]layout)

	t.editor = t.setupEditor()
	t.console = t.editor.Console()
	t.currentWindow = t.editor.NewWindow()

	// Assuming that all extra arguments are files
	if files := flag.Args(); len(files) > 0 {
		for _, file := range files {
			t.currentView = createNewView(file, t.currentWindow)
		}
	} else {
		t.currentView = t.currentWindow.NewFile()
	}

	t.console.Buffer().AddObserver(&t)
	t.setupCallbacks(t.currentView)

	path := path.Join("..", "..", "packages", "themes", "TextMate-Themes", "Monokai.tmTheme")
	if sc, err := textmate.LoadTheme(path); err != nil {
		log.Error(err)
	} else {
		scheme = sc
	}

	setColorMode()
	setSchemeSettings()

	w, h := termbox.Size()
	t.handleResize(h, w, true)

	// These might take a while
	t.editor.Init()
	go sublime.Init()

	return &t
}

func (t *tbfe) renderView(v *backend.View, lay layout) {
	p := util.Prof.Enter("render")
	defer p.Exit()

	sx, sy, w, h := lay.x, lay.y, lay.width, lay.height
	vr := lay.visible
	runes := v.Buffer().Substr(vr)
	x, y := sx, sy
	ex, ey := sx+w, sy+h

	style, _ := v.Settings().Get("caret_style", "underline").(string)
	inverse, _ := v.Settings().Get("inverse_caret_state", false).(bool)

	caretStyle := getCaretStyle(style, inverse)
	caretBlink, _ := v.Settings().Get("caret_blink", true).(bool)
	if caretBlink && blink {
		caretStyle = 0
	}
	tabSize := 4
	ts := v.Settings().Get("tab_size", tabSize)
	// TODO(.): crikey...
	if i, ok := ts.(int); ok {
		tabSize = i
	} else if f, ok := ts.(float64); ok {
		tabSize = int(f)
	}

	lineNumbers, _ := v.Settings().Get("line_numbers", true).(bool)

	recipie := v.Transform(scheme, vr).Transcribe()

	fg, bg := defaultFg, defaultBg
	sel := v.Sel()

	line, _ := v.Buffer().RowCol(vr.Begin())
	eofline, _ := v.Buffer().RowCol(v.Buffer().Size())
	lineNumberRenderSize := len(intToRunes(eofline))

	for i, r := range runes {
		o := vr.Begin() + i
		curr := 0
		fg, bg = defaultFg, defaultBg

		if lineNumbers {
			renderLineNumber(&line, &x, y, lineNumberRenderSize, fg, bg)
		}

		for curr < len(recipie) && (o >= recipie[curr].Region.Begin()) {
			if o < recipie[curr].Region.End() {
				fg = palLut(textmate.Color(recipie[curr].Flavour.Foreground))
				bg = palLut(textmate.Color(recipie[curr].Flavour.Background))
			}
			curr++
		}
		iscursor := sel.Contains(Region{o, o})
		if iscursor {
			fg = fg | caretStyle
			termbox.SetCell(x, y, ' ', fg, bg)
		}
		if r == '\t' {
			add := (x + 1 + (tabSize - 1)) &^ (tabSize - 1)
			for x < add {
				if x < ex {
					termbox.SetCell(x, y, ' ', fg, bg)
				}
				// A long cursor looks weird
				fg = fg & ^(termbox.AttrUnderline | termbox.AttrReverse)
				x++
			}
			continue
		} else if r == '\n' {
			x = sx
			if y++; y > ey {
				break
			} else if lineNumbers {
				renderLineNumber(&line, &x, y, lineNumberRenderSize, fg, bg)
			}
			continue
		}
		if x < ex {
			termbox.SetCell(x, y, r, fg, bg)
		}
		x++
	}
	fg, bg = defaultFg, defaultBg
	// Need this if the cursor is at the end of the buffer
	o := vr.Begin() + len(runes)
	iscursor := sel.Contains(Region{o, o})
	if iscursor {
		fg = fg | caretStyle
		termbox.SetCell(x, y, ' ', fg, bg)
	}

	// restore original caretStyle before blink modification
	caretStyle = caretStyle
}

func (t *tbfe) clip(v *backend.View, s, e int) Region {
	p := util.Prof.Enter("clip")
	defer p.Exit()
	t.lock.Lock()
	h := t.layout[v].height
	t.lock.Unlock()
	if e-s > h {
		e = s + h
	} else if e-s < h {
		s = e - h
	}
	if e2, _ := v.Buffer().RowCol(v.Buffer().TextPoint(e, 0)); e2 < e {
		e = e2
	}
	if s < 0 {
		s = 0
	}
	e = s + h
	r := Region{v.Buffer().TextPoint(s, 0), v.Buffer().TextPoint(e, 0)}
	return v.Buffer().LineR(r)
}

func (t *tbfe) Show(v *backend.View, r Region) {
	t.lock.Lock()
	l := t.layout[v]
	t.lock.Unlock()
	p := util.Prof.Enter("show")
	defer p.Exit()

	lv := l.visible

	s1, _ := v.Buffer().RowCol(lv.Begin())
	e1, _ := v.Buffer().RowCol(lv.End())
	s2, _ := v.Buffer().RowCol(r.Begin())
	e2, _ := v.Buffer().RowCol(r.End())

	r1 := Region{s1, e1}
	r2 := Region{s2, e2}

	r3 := r1.Cover(r2)
	diff := 0
	if d1, d2 := Abs(r1.Begin()-r3.Begin()), Abs(r1.End()-r3.End()); d1 > d2 {
		diff = r3.Begin() - r1.Begin()
	} else {
		diff = r3.End() - r1.End()
	}
	r3.A = r1.Begin() + diff
	r3.B = r1.End() + diff

	r3 = t.clip(v, r3.A, r3.B)
	l.visible = r3
	t.lock.Lock()
	t.layout[v] = l
	t.lock.Unlock()
	t.render()
}

func (t *tbfe) VisibleRegion(v *backend.View) Region {
	t.lock.Lock()
	r, ok := t.layout[v]
	t.lock.Unlock()
	if !ok || r.lastUpdate != v.Buffer().ChangeCount() {
		t.Show(v, r.visible)
		t.lock.Lock()
		r = t.layout[v]
		t.lock.Unlock()
	}
	return r.visible
}

func (t *tbfe) StatusMessage(msg string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.status_message = msg
}

func (t *tbfe) ErrorMessage(msg string) {
	log.Error(msg)
}

// TODO(q): Actually show a dialog
func (t *tbfe) MessageDialog(msg string) {
	log.Info(msg)
}

// TODO(q): Actually show a dialog
func (t *tbfe) OkCancelDialog(msg, ok string) bool {
	log.Info(msg, ok)
	return false
}

func (t *tbfe) scroll(b Buffer) {
	t.Show(backend.GetEditor().Console(), Region{b.Size(), b.Size()})
}

func (t *tbfe) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	t.scroll(changed_buffer)
}

func (t *tbfe) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	t.scroll(changed_buffer)
}

func (t tbfe) setupCallbacks(view *backend.View) {
	// Ensure that the visible region currently presented is
	// inclusive of the insert/erase delta.
	view.Buffer().AddObserver(&tbfeBufferDeltaObserver{t: &t, view: view})

	backend.OnNew.Add(func(v *backend.View) {
		v.Settings().AddOnChange("lime.frontend.termbox.render", func(name string) { t.render() })
	})

	backend.OnModified.Add(func(v *backend.View) {
		t.render()
	})

	backend.OnSelectionModified.Add(func(v *backend.View) {
		t.render()
	})
}

func (t *tbfe) setupEditor() *backend.Editor {
	ed := backend.GetEditor()
	ed.SetFrontend(t)
	ed.LogInput(false)
	ed.LogCommands(false)

	return ed
}

func (t *tbfe) render() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if len(t.dorender) < cap(t.dorender) {
		t.dorender <- true
	}
}

func (t *tbfe) renderthread() {
	pc := 0
	dorender := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic in renderthread: %v\n%s", r, string(debug.Stack()))
				if pc > 1 {
					panic(r)
				}
				pc++
			}
		}()
		termbox.Clear(defaultFg, defaultBg)

		t.lock.Lock()
		vs := make([]*backend.View, 0, len(t.layout))
		l := make([]layout, 0, len(t.layout))
		for k, v := range t.layout {
			vs = append(vs, k)
			l = append(l, v)
		}
		runes := []rune(t.status_message)
		t.lock.Unlock()

		w, h := termbox.Size()
		for i := 0; i < w && i < len(runes); i++ {
			termbox.SetCell(i, h-2, runes[i], defaultFg, defaultBg)
		}

		for i, v := range vs {
			t.renderView(v, l[i])
		}

		termbox.Flush()
	}

	for _ = range t.dorender {
		log.Finest("Rendering")
		dorender()
	}
}

func (t *tbfe) handleResize(height, width int, init bool) {
	// This should handle multiple views in a less hardcoded fashion.
	// After all, it is possible to *not* have a view in a window.
	t.lock.Lock()
	if init {
		t.layout[t.currentView] = layout{0, 0, 0, 0, Region{}, 0}
		if *showConsole {
			t.layout[t.console] = layout{0, 0, 0, 0, Region{}, 0}
		}
	}

	if *showConsole {
		view_layout := t.layout[t.currentView]
		view_layout.height = height - *consoleHeight - 4
		view_layout.width = width

		console_layout := t.layout[t.console]
		console_layout.y = height - *consoleHeight - 2
		console_layout.width = width
		console_layout.height = *consoleHeight - 1

		t.layout[t.console] = console_layout
		t.layout[t.currentView] = view_layout
	} else {
		view_layout := t.layout[t.currentView]
		view_layout.height = height - 3
		view_layout.width = width
		t.layout[t.currentView] = view_layout
	}
	t.lock.Unlock()

	// Ensure that the new visible region is recalculated
	t.Show(t.currentView, t.VisibleRegion(t.currentView))
}

func (t *tbfe) handleInput(ev termbox.Event) {
	if ev.Key == termbox.KeyCtrlQ {
		t.shutdown <- true
	}

	var kp keys.KeyPress
	if ev.Ch != 0 {
		kp.Key = keys.Key(ev.Ch)
	} else if v2, ok := lut[ev.Key]; ok {
		kp = v2
	} else {
		return
	}

	t.editor.HandleInput(kp)
}

func (t *tbfe) loop() {
	timechan := make(chan bool, 0)

	// Only set up the timers if we should actually blink the cursor
	// This should somehow be changable on an OnSettingsChanged callback
	if p := t.editor.Settings().Get("caret_blink", true).(bool); p {
		duration := time.Second / 2
		if p, ok := t.editor.Settings().Get("caret_blink_phase", 1.0).(float64); ok {
			duration = time.Duration(float64(time.Second)*p) / 2
		}
		timer := time.NewTimer(duration)

		defer func() {
			timer.Stop()
			close(timechan)
		}()

		go func() {
			for _ = range timer.C {
				timechan <- true
				timer.Reset(duration)
			}
		}()
	}

	// Due to termbox still running, we can't close evchan
	evchan := make(chan termbox.Event, 32)
	go func() {
		for {
			evchan <- termbox.PollEvent()
		}
	}()

	for {
		p := util.Prof.Enter("mainloop")
		select {
		case ev := <-evchan:
			mp := util.Prof.Enter("evchan")
			switch ev.Type {
			case termbox.EventError:
				log.Debug("error occured")
				return
			case termbox.EventResize:
				t.handleResize(ev.Height, ev.Width, false)
			case termbox.EventKey:
				t.handleInput(ev)
				blink = false
			}
			mp.Exit()

		case <-timechan:
			blink = !blink
			t.render()

		case <-t.shutdown:
			return
		}
		p.Exit()
	}
}

func (bdo *tbfeBufferDeltaObserver) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	ensureVisibleRegionContainsInsertOrEraseDelta(bdo.t, bdo.view, region_removed.A-region_removed.B)
}

func (bdo *tbfeBufferDeltaObserver) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	ensureVisibleRegionContainsInsertOrEraseDelta(bdo.t, bdo.view, region_inserted.B-region_inserted.A)
}

func ensureVisibleRegionContainsInsertOrEraseDelta(t *tbfe, view *backend.View, delta int) {
	t.lock.Lock()
	visible := t.layout[view].visible
	t.lock.Unlock()
	t.Show(view, Region{visible.Begin(), visible.End() + delta})
}

func intToRunes(n int) (runes []rune) {
	lineStr := strconv.FormatInt(int64(n), 10)

	return []rune(lineStr)
}

func padLineRunes(line []rune, totalLineSize int) (padded []rune) {
	currentLineSize := len(line)
	if currentLineSize < totalLineSize {
		padding := (totalLineSize - currentLineSize)

		for i := 0; i < padding; i++ {
			padded = append(padded, ' ')
		}
	}

	padded = append(padded, line...)
	padded = append(padded, ' ')

	return
}

func renderLineNumber(line, x *int, y, lineNumberRenderSize int, fg, bg termbox.Attribute) {
	if *x == 0 {
		lineRunes := padLineRunes(intToRunes(*line), lineNumberRenderSize)

		for _, num := range lineRunes {
			termbox.SetCell(*x, y, num, fg, bg)
			*x++
		}

		*line++
	}

}

func getCaretStyle(style string, inverse bool) termbox.Attribute {
	caret_style := termbox.AttrUnderline

	if style == "block" {
		caret_style = termbox.AttrReverse
	}

	if inverse {
		if caret_style == termbox.AttrReverse {
			caret_style = termbox.AttrUnderline
		} else {
			caret_style = termbox.AttrReverse
		}
	}

	return caret_style
}

func setColorMode() {
	var (
		mode256 bool
		pal     = make([]termbox.RGB, 0, 256)
	)

	if err := termbox.SetColorMode(termbox.ColorMode256); err != nil {
		log.Error("Unable to use 256 color mode: %s", err)
	} else {
		log.Debug("Using 256 color mode")
		mode256 = true
	}

	if !mode256 {
		pal = pal[:10] // Not correct, but whatever
		pal[termbox.ColorBlack] = termbox.RGB{0, 0, 0}
		pal[termbox.ColorWhite] = termbox.RGB{255, 255, 255}
		pal[termbox.ColorRed] = termbox.RGB{255, 0, 0}
		pal[termbox.ColorGreen] = termbox.RGB{0, 255, 0}
		pal[termbox.ColorBlue] = termbox.RGB{0, 0, 255}
		pal[termbox.ColorMagenta] = termbox.RGB{255, 0, 255}
		pal[termbox.ColorYellow] = termbox.RGB{255, 255, 0}
		pal[termbox.ColorCyan] = termbox.RGB{0, 255, 255}

		diff := func(i, j byte) int {
			v := int(i) - int(j)
			if v < 0 {
				return -v
			}
			return v
		}
		palLut = func(col textmate.Color) termbox.Attribute {
			mindist := 10000000
			mini := 0
			for i, c := range pal {
				if dist := diff(c.R, col.R) + diff(c.G, col.G) + diff(c.B, col.B); dist < mindist {
					mindist = dist
					mini = i
				}
			}
			return termbox.Attribute(mini)
		}
	} else {
		palLut = func(col textmate.Color) termbox.Attribute {
			tc := termbox.RGB{col.R, col.G, col.B}
			for i, c := range pal {
				if c == tc {
					return termbox.Attribute(i)
				}
			}
			l := len(pal)
			log.Debug("Adding colour: %d %+v %+v", l, col, tc)
			pal = append(pal, tc)
			termbox.SetColorPalette(pal)
			return termbox.Attribute(l)
		}
	}
}

func setSchemeSettings() {
	for i, s := range scheme.Settings {
		var (
			fi = defaultFg
			bi = defaultBg
		)
		if fg, ok := s.Settings["foreground"]; ok {
			fi = palLut(fg)
			if i == 0 {
				defaultFg = fi
			}
		}
		if bg, ok := s.Settings["background"]; ok {
			bi = palLut(bg)
			if i == 0 {
				defaultBg = bi
			}
		}
	}
}

func createNewView(filename string, window *backend.Window) *backend.View {
	syntax := "../../packages/go.tmbundle/Syntaxes/Go.tmLanguage"
	v := window.OpenFile(filename, 0)

	v.Settings().Set("trace", true)
	v.Settings().Set("syntax", syntax)

	return v
}

func main() {
	flag.Parse()

	log.AddFilter("file", log.FINEST, log.NewFileLogWriter("debug.log", *rotateLog))
	defer func() {
		py.NewLock()
		py.Finalize()
	}()

	if err := termbox.Init(); err != nil {
		log.Close(err)
		return
	}

	defer func() {
		termbox.Close()
		log.Debug(util.Prof)
		if err := recover(); err != nil {
			log.Critical(err)
			panic(err)
		}
	}()

	t := createFrontend()
	go t.renderthread()
	t.loop()
}
