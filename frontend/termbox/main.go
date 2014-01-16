// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package main

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/quarnster/util/text"
	"lime/3rdparty/libs/gopy/lib"
	"lime/3rdparty/libs/termbox-go"
	"lime/backend"
	_ "lime/backend/commands"
	"lime/backend/sublime"
	"lime/backend/textmate"
	"lime/backend/util"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

var (
	lut = map[termbox.Key]backend.KeyPress{
		// Omission of these are intentional due to map collisions
		//		termbox.KeyCtrlTilde:      backend.KeyPress{Ctrl: true, Key: '~'},
		//		termbox.KeyCtrlBackslash:  backend.KeyPress{Ctrl: true, Key: '\\'},
		//		termbox.KeyCtrlSlash:      backend.KeyPress{Ctrl: true, Key: '/'},
		//		termbox.KeyCtrlUnderscore: backend.KeyPress{Ctrl: true, Key: '_'},
		//		termbox.KeyCtrlLsqBracket: backend.KeyPress{Ctrl: true, Key: '{'},
		//		termbox.KeyCtrlRsqBracket: backend.KeyPress{Ctrl: true, Key: '}'},
		// termbox.KeyCtrl3:
		// termbox.KeyCtrl8
		//		termbox.KeyCtrl2:      backend.KeyPress{Ctrl: true, Key: '2'},
		termbox.KeyCtrlSpace:  backend.KeyPress{Ctrl: true, Key: ' '},
		termbox.KeyCtrlA:      backend.KeyPress{Ctrl: true, Key: 'a'},
		termbox.KeyCtrlB:      backend.KeyPress{Ctrl: true, Key: 'b'},
		termbox.KeyCtrlC:      backend.KeyPress{Ctrl: true, Key: 'c'},
		termbox.KeyCtrlD:      backend.KeyPress{Ctrl: true, Key: 'd'},
		termbox.KeyCtrlE:      backend.KeyPress{Ctrl: true, Key: 'e'},
		termbox.KeyCtrlF:      backend.KeyPress{Ctrl: true, Key: 'f'},
		termbox.KeyCtrlG:      backend.KeyPress{Ctrl: true, Key: 'g'},
		termbox.KeyCtrlH:      backend.KeyPress{Ctrl: true, Key: 'h'},
		termbox.KeyCtrlJ:      backend.KeyPress{Ctrl: true, Key: 'j'},
		termbox.KeyCtrlK:      backend.KeyPress{Ctrl: true, Key: 'k'},
		termbox.KeyCtrlL:      backend.KeyPress{Ctrl: true, Key: 'l'},
		termbox.KeyCtrlN:      backend.KeyPress{Ctrl: true, Key: 'n'},
		termbox.KeyCtrlO:      backend.KeyPress{Ctrl: true, Key: 'o'},
		termbox.KeyCtrlP:      backend.KeyPress{Ctrl: true, Key: 'p'},
		termbox.KeyCtrlQ:      backend.KeyPress{Ctrl: true, Key: 'q'},
		termbox.KeyCtrlR:      backend.KeyPress{Ctrl: true, Key: 'r'},
		termbox.KeyCtrlS:      backend.KeyPress{Ctrl: true, Key: 's'},
		termbox.KeyCtrlT:      backend.KeyPress{Ctrl: true, Key: 't'},
		termbox.KeyCtrlU:      backend.KeyPress{Ctrl: true, Key: 'u'},
		termbox.KeyCtrlV:      backend.KeyPress{Ctrl: true, Key: 'v'},
		termbox.KeyCtrlW:      backend.KeyPress{Ctrl: true, Key: 'w'},
		termbox.KeyCtrlX:      backend.KeyPress{Ctrl: true, Key: 'x'},
		termbox.KeyCtrlY:      backend.KeyPress{Ctrl: true, Key: 'y'},
		termbox.KeyCtrlZ:      backend.KeyPress{Ctrl: true, Key: 'z'},
		termbox.KeyCtrl4:      backend.KeyPress{Ctrl: true, Key: '4'},
		termbox.KeyCtrl5:      backend.KeyPress{Ctrl: true, Key: '5'},
		termbox.KeyCtrl6:      backend.KeyPress{Ctrl: true, Key: '6'},
		termbox.KeyCtrl7:      backend.KeyPress{Ctrl: true, Key: '7'},
		termbox.KeyEnter:      backend.KeyPress{Key: backend.Enter},
		termbox.KeySpace:      backend.KeyPress{Key: ' '},
		termbox.KeyBackspace2: backend.KeyPress{Key: backend.Backspace},
		termbox.KeyArrowUp:    backend.KeyPress{Key: backend.Up},
		termbox.KeyArrowDown:  backend.KeyPress{Key: backend.Down},
		termbox.KeyArrowLeft:  backend.KeyPress{Key: backend.Left},
		termbox.KeyArrowRight: backend.KeyPress{Key: backend.Right},
		termbox.KeyDelete:     backend.KeyPress{Key: backend.Delete},
		termbox.KeyEsc:        backend.KeyPress{Key: backend.Escape},
		termbox.KeyPgup:       backend.KeyPress{Key: backend.PageUp},
		termbox.KeyPgdn:       backend.KeyPress{Key: backend.PageDown},
		termbox.KeyF1:         backend.KeyPress{Key: backend.F1},
		termbox.KeyF2:         backend.KeyPress{Key: backend.F2},
		termbox.KeyF3:         backend.KeyPress{Key: backend.F3},
		termbox.KeyF4:         backend.KeyPress{Key: backend.F4},
		termbox.KeyF5:         backend.KeyPress{Key: backend.F5},
		termbox.KeyF6:         backend.KeyPress{Key: backend.F6},
		termbox.KeyF7:         backend.KeyPress{Key: backend.F7},
		termbox.KeyF8:         backend.KeyPress{Key: backend.F8},
		termbox.KeyF9:         backend.KeyPress{Key: backend.F9},
		termbox.KeyF10:        backend.KeyPress{Key: backend.F10},
		termbox.KeyF11:        backend.KeyPress{Key: backend.F11},
		termbox.KeyF12:        backend.KeyPress{Key: backend.F12},
		termbox.KeyTab:        backend.KeyPress{Key: '\t'},
	}
	palLut    func(col textmate.Color) termbox.Attribute
	scheme    *textmate.Theme
	defaultBg = termbox.ColorBlack
	defaultFg = termbox.ColorWhite
	blink     bool
)

const (
	console_height  = 20
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
	lock           sync.Mutex
}

func (t *tbfe) renderView(v *backend.View, lay layout) {
	p := util.Prof.Enter("render")
	defer p.Exit()

	sx, sy, w, h := lay.x, lay.y, lay.width, lay.height
	vr := lay.visible
	runes := v.Buffer().Substr(vr)
	x, y := sx, sy
	ex, ey := sx+w, sy+h

	tab_size, ok := v.Settings().Get("tab_size", 4).(int)
	if !ok {
		tab_size = 4
	}

	recipie := v.Transform(scheme, vr).Transcribe()

	curr := 0
	fg, bg := defaultFg, defaultBg
	sel := v.Sel()

	caret_blink := true
	if b, ok := v.Settings().Get("caret_blink", true).(bool); ok {
		caret_blink = b
	}

	highlight_line := false
	if b, ok := v.Settings().Get("highlight_line", highlight_line).(bool); ok {
		highlight_line = b
	}
	caret_style := termbox.AttrUnderline
	if b, ok := v.Settings().Get("caret_style", "underline").(string); ok {
		if b == "block" {
			caret_style = termbox.AttrReverse
		}
	}
	if b, ok := v.Settings().Get("inverse_caret_state", false).(bool); !b && ok {
		if caret_style == termbox.AttrReverse {
			caret_style = termbox.AttrUnderline
		} else {
			caret_style = termbox.AttrReverse
		}
	}

	if caret_blink && blink {
		caret_style = 0
	}

	shouldRenderLineNumbers, _ := v.Settings().Get("line_numbers", true).(bool)

	line, _ := v.Buffer().RowCol(vr.Begin())
	eofline, _ := v.Buffer().RowCol(v.Buffer().Size())
	lineNumberRenderSize := len(intToRunes(eofline))

	for i, r := range runes {
		o := vr.Begin() + i
		curr = 0
		fg, bg = defaultFg, defaultBg

		if shouldRenderLineNumbers {
			if x == 0 {
				lineRunes := padLineRunes(intToRunes(line), lineNumberRenderSize)

				for _, num := range lineRunes {
					termbox.SetCell(x, y, num, fg, bg)
					x++
				}

				line++
			}
		}

		for curr < len(recipie) && (o >= recipie[curr].Region.Begin()) {
			if o < recipie[curr].Region.End() {
				fg = palLut(textmate.Color(recipie[curr].Flavour.Foreground))
				bg = palLut(textmate.Color(recipie[curr].Flavour.Background))
			}
			curr++
		}
		if sel.Contains(Region{o, o}) {
			fg = fg | caret_style
		}
		if r == '\t' {
			add := (x + 1 + (tab_size - 1)) &^ (tab_size - 1)
			for x < add {
				if x < ex {
					termbox.SetCell(x, y, ' ', fg, bg)
				}
				fg = fg &^ termbox.AttrUnderline // Just looks weird with a long underline
				x++
			}
			continue
		} else if r == '\n' {
			x = sx
			y++
			if y > ey {
				break
			}
			continue
		}
		if x < ex {
			termbox.SetCell(x, y, r, fg, bg)
		}
		x++
	}

	if shouldRenderLineNumbers {
		if x == 0 {
			lineRunes := padLineRunes(intToRunes(line), lineNumberRenderSize)

			for _, num := range lineRunes {
				termbox.SetCell(x, y, num, fg, bg)
				x++
			}
		}
	}
}

func (t *tbfe) clip(v *backend.View, s, e int) Region {
	p := util.Prof.Enter("clip")
	defer p.Exit()
	h := t.layout[v].height
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
	if l.visible.Covers(r) {
		return
	}
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
	log4go.Error(msg)
}

// TODO(q): Actually show a dialog
func (t *tbfe) MessageDialog(msg string) {
	log4go.Info(msg)
}

// TODO(q): Actually show a dialog
func (t *tbfe) OkCancelDialog(msg, ok string) {
	log4go.Info(msg, ok)
}

func (t *tbfe) scroll(b Buffer, pos, delta int) {
	t.Show(backend.GetEditor().Console(), Region{b.Size(), b.Size()})
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
				log4go.Error("Panic in renderthread: %v\n%s", r, string(debug.Stack()))
				if pc > 1 {
					panic(r)
				}
				pc++
			}
		}()
		termbox.Clear(defaultFg, defaultBg)
		w, h := termbox.Size()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				termbox.SetCell(x, y, ' ', defaultFg, defaultBg)
			}
		}
		t.lock.Lock()
		vs := make([]*backend.View, 0, len(t.layout))
		l := make([]layout, 0, len(t.layout))
		for k, v := range t.layout {
			vs = append(vs, k)
			l = append(l, v)
		}
		t.lock.Unlock()
		for i, v := range vs {
			t.renderView(v, l[i])
		}
		t.lock.Lock()
		runes := []rune(t.status_message)
		t.lock.Unlock()
		for i := 0; i < w && i < len(runes); i++ {
			termbox.SetCell(i, h-1, runes[i], defaultFg, defaultBg)
		}
		termbox.Flush()
	}
	for a := range t.dorender {
		_ = a
		log4go.Finest("Rendering")
		dorender()
	}
}

func (t *tbfe) loop() {
	backend.OnNew.Add(func(v *backend.View) {
		v.Settings().AddOnChange("lime.frontend.termbox.render", func(name string) { t.render() })
	})
	backend.OnModified.Add(func(v *backend.View) {
		t.render()
	})
	backend.OnSelectionModified.Add(func(v *backend.View) {
		t.render()
	})

	ed := backend.GetEditor()
	ed.SetFrontend(t)
	ed.LogInput(false)
	ed.LogCommands(false)
	c := ed.Console()
	if sc, err := textmate.LoadTheme("../../3rdparty/bundles/TextMate-Themes/GlitterBomb.tmTheme"); err != nil {
		log4go.Error(err)
	} else {
		scheme = sc
	}

	var (
		pal     = make([]termbox.RGB, 0, 256)
		mode256 bool
	)

	if err := termbox.SetColorMode(termbox.ColorMode256); err != nil {
		log4go.Error("Unable to use 256 color mode: %s", err)
	} else {
		log4go.Debug("Using 256 color mode")
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
			log4go.Debug("Adding colour: %d %+v %+v", l, col, tc)
			pal = append(pal, tc)
			termbox.SetColorPalette(pal)
			return termbox.Attribute(l)
		}
	}
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
		for _, setting := range []string{"caret", "highlight", "invisibles", "selection"} {
			if col, ok := s.Settings[setting]; ok {
				i := palLut(col)
				if setting == "selection" {
					fmt.Println(col, i)
				}
			}
		}
	}
	evchan := make(chan termbox.Event, 32)
	defer func() {
		close(evchan)
		termbox.Close()
		fmt.Println(util.Prof)
	}()

	w := ed.NewWindow()
	v := w.OpenFile("main.go", 0)
	v.Settings().Set("trace", true)
	v.Settings().Set("syntax", "../../3rdparty/bundles/go.tmbundle/Syntaxes/Go.tmLanguage")
	c.Buffer().AddCallback(t.scroll)

	sel := v.Sel()
	sel.Clear()
	sel.Add(Region{0, 0})

	go func() {
		for {
			evchan <- termbox.PollEvent()
		}
	}()

	{
		w, h := termbox.Size()
		t.lock.Lock()
		t.layout[v] = layout{0, 0, w, h - console_height - 1, Region{}, 0}
		t.layout[c] = layout{0, h - console_height + 1, w, console_height - 5, Region{}, 0}
		t.lock.Unlock()
		t.Show(v, Region{1, 1})
	}
	t.Show(v, Region{100, 100})
	t.Show(v, Region{1, 1})

	go func() {
		ed.Init()
		sublime.Init()
	}()

	for {
		p := util.Prof.Enter("mainloop")

		blink_phase := time.Second
		if p, ok := ed.Settings().Get("caret_blink_phase", 1.0).(float64); ok {
			blink_phase = time.Duration(float64(time.Second) * p)
		}

		// Divided by two since we're only doing a simple toggle blink
		timer := time.NewTimer(blink_phase / 2)
		select {
		case ev := <-evchan:
			mp := util.Prof.Enter("evchan")
			limit := 3
		loop:
			switch ev.Type {
			case termbox.EventError:
				log4go.Debug("error occured")
				return
			case termbox.EventKey:
				var kp backend.KeyPress

				if ev.Ch != 0 {
					kp.Key = backend.Key(ev.Ch)
				} else if v2, ok := lut[ev.Key]; ok {
					kp = v2
				} else {
					break
				}

				if ev.Key == termbox.KeyCtrlQ {
					return
				}
				ed.HandleInput(kp)

				blink = false
			}
			if len(evchan) > 0 {
				limit--
				ev = <-evchan
				goto loop
			}
			mp.Exit()
		case <-timer.C:
			// TODO(q): Shouldn't redraw if blink is disabled...

			blink = !blink
			t.render()
		}
		timer.Stop()
		p.Exit()
	}
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

func main() {
	log4go.AddFilter("file", log4go.FINEST, log4go.NewFileLogWriter("debug.log", true))
	defer func() {
		py.NewLock()
		py.Finalize()
	}()

	if err := termbox.Init(); err != nil {
		log4go.Exit(err)
	}

	var t tbfe
	t.dorender = make(chan bool, render_chan_len)
	t.layout = make(map[*backend.View]layout)
	go t.renderthread()
	t.loop()
}
