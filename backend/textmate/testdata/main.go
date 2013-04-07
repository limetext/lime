package main

import (
	"code.google.com/p/log4go"
	"fmt"
	"io/ioutil"
	"lime/3rdparty/libs/termbox-go"
	"lime/backend"
	"lime/backend/loaders"
	"lime/backend/primitives"
	"lime/backend/sublime"
	"lime/backend/textmate"
	"strings"
	"time"
)

var (
	lut = map[termbox.Key]backend.KeyPress{
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
		termbox.KeyCtrl2:      backend.KeyPress{Ctrl: true, Key: '2'},
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
		termbox.KeyTab:        backend.KeyPress{Key: '\t'},
	}
	schemelut = make(map[string][2]termbox.Attribute)
	defaultBg = termbox.ColorBlack
	defaultFg = termbox.ColorWhite
	blink     bool
)

const console_height = 20

type tbfe struct {
	visibleregion  map[*backend.View]primitives.Region
	status_message string
	active_window  *backend.Window
}

func (t *tbfe) ActiveWindow() *backend.Window {
	return t.active_window
}

func (t *tbfe) ActiveView(w *backend.Window) *backend.View {
	if w == nil {
		return nil
	}
	if v := w.Views(); len(v) > 0 {
		return v[0]
	}
	return nil
}

func (t *tbfe) renderView(sx, sy, w, h int, v *backend.View) {
	sel := v.Sel()
	substr := v.Buffer().Data()
	vr := t.VisibleRegion(v)
	lines := strings.Split(substr, "\n")
	s, _ := v.Buffer().RowCol(vr.Begin())
	e, _ := v.Buffer().RowCol(vr.End())
	if e > 1 {
		e = e - 1
		if e > h {
			s = e - h
		}
	}
	off := len(strings.Join(lines[:s], "\n"))
	lines = lines[s:e]
	runes := []rune(strings.Join(lines, "\n"))
	x, y := sx, sy
	ex, ey := sx+w, sy+h

	sub2 := ""
	var (
		lastScope string
		lfg, lbg  = defaultFg, defaultBg
	)

	tab_size, ok := v.Settings().Get("tab_size", 4).(int)
	if !ok {
		tab_size = 4
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
	caret_blink := true
	if b, ok := v.Settings().Get("caret_blink", true).(bool); ok {
		caret_blink = b
	}

	for i := range runes {
		sub2 += string(runes[i])

		if x < ex {
			o := off + len(sub2)
			r := primitives.Region{o, o}
			fg, bg := lfg, lbg
			scope := v.ScopeName(o)
			if scope != lastScope {
				lastScope = scope
				na := scope
				for len(na) > 0 {
					sn := na
					i := strings.LastIndex(sn, " ")
					if i != -1 {
						sn = sn[i+1:]
					}
					if c, ok := schemelut[sn]; ok {
						fg, bg = c[0], c[1]
						break
					}
					if i2 := strings.LastIndex(na, "."); i2 == -1 {
						break
					} else if i > i2 {
						na = na[:i]
					} else {
						na = strings.TrimSpace(na[:i2])
					}
				}
				lfg, lbg = fg, bg
			} else {
				fg, bg = lfg, lbg
			}
			for _, r2 := range sel.Regions() {
				if r2.B == r.B {
					if !caret_blink || blink {
						if r2.Contains(o) {
							fg |= termbox.AttrReverse
						} else {
							fg |= caret_style
						}
					}
					break
				} else if r2.Contains(o) {
					fg |= termbox.AttrReverse
					break
				}
			}
			if runes[i] == '\t' {
				add := (x + 1 + (tab_size - 1)) &^ (tab_size - 1)

				for x < add {
					termbox.SetCell(x, y, ' ', fg, bg)
					x++
				}
				continue
			} else if runes[i] == '\n' {
				termbox.SetCell(x, y, ' ', fg, bg)
				x = sx
				y++
				if y > ey {
					break
				}
				continue
			}
			termbox.SetCell(x, y, runes[i], fg, bg)
		}
		x++
	}
}

func (t *tbfe) clip(v *backend.View, r primitives.Region) primitives.Region {
	s, _ := v.Buffer().RowCol(r.Begin())
	e, _ := v.Buffer().RowCol(r.End())

	_, h := termbox.Size()
	h -= console_height
	if e-s > h {
		e = s + h
	}
	if e2, _ := v.Buffer().RowCol(v.Buffer().TextPoint(e, 1)); e2 < e {
		e = e2
	}
	if e-s < h {
		s = e - h
	}
	if s < 1 {
		s = 1
	}
	r.A = v.Buffer().Line(v.Buffer().TextPoint(s, 1)).A
	r.B = v.Buffer().Line(v.Buffer().TextPoint(e, 1)).B
	return r
}

func (t *tbfe) Show(v *backend.View, r primitives.Region) {
	t.visibleregion[v] = t.clip(v, primitives.Region{r.Begin(), v.Buffer().Size()})
}

func (t *tbfe) VisibleRegion(v *backend.View) primitives.Region {
	if r, ok := t.visibleregion[v]; ok {
		return r
	} else {
		t.Show(v, primitives.Region{0, 0})
		return t.visibleregion[v]
	}
}

func (t *tbfe) StatusMessage(msg string) {
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

func (t *tbfe) scroll(b *primitives.Buffer, pos, delta int) {
	t.Show(backend.GetEditor().Console(), primitives.Region{b.Size(), b.Size()})
}

func (t *tbfe) loop() {
	ed := backend.GetEditor()
	ed.SetFrontend(t)
	//ed.LogInput(true)
	//ed.LogCommands(true)
	c := ed.Console()
	var (
		scheme textmate.Theme
	)

	if d, err := ioutil.ReadFile("../../3rdparty/bundles/TextMate-Themes/GlitterBomb.tmTheme"); err != nil {
		log4go.Error("Unable to load colorscheme definition: %s", err)
	} else if err := loaders.LoadPlist(d, &scheme); err != nil {
		log4go.Error("Unable to load colorscheme definition: %s", err)
	}

	var (
		palLut  func(col textmate.Color) termbox.Attribute
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
			pal = append(pal, tc)
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
		schemelut[s.Scope] = [2]termbox.Attribute{fi, bi}
	}
	if mode256 {
		termbox.SetColorPalette(pal)
	}
	defer func() {
		termbox.Close()
		fmt.Println(c.Buffer().Data())
	}()

	w := ed.NewWindow()
	t.active_window = w
	v := w.OpenFile("main.go", 0)
	v.Settings().Set("trace", true)
	c.Buffer().AddCallback(t.scroll)

	if err := v.SetSyntaxFile("../../3rdparty/bundles/GoSublime/GoSublime.tmLanguage"); err != nil {
		log4go.Error("Unable to set syntax file: %s", err)
	}
	sel := v.Sel()
	sel.Clear()
	end := v.Buffer().Size() - 2
	sel.Add(primitives.Region{end - 24, end - 24})
	// sel.Add(primitives.Region{end - 22, end - 22})
	// sel.Add(primitives.Region{end - 16, end - 20})
	// sel.Add(primitives.Region{end - 13, end - 10})

	evchan := make(chan termbox.Event)

	go func() {
		for {
			evchan <- termbox.PollEvent()
		}
	}()

	sublime.Init()
	for {
		blink = !blink
		termbox.Clear(defaultFg, defaultBg)
		w, h := termbox.Size()

		t.renderView(0, 0, w, h-console_height, v)
		t.renderView(0, h-(console_height), w, (console_height - 1), c)
		runes := []rune(t.status_message)
		for i := 0; i < w && i < len(runes); i++ {
			termbox.SetCell(i, h-1, runes[i], defaultFg, defaultBg)
		}
		termbox.Flush()

		blink_phase := time.Second
		if p, ok := ed.Settings().Get("caret_blink_phase", 1.0).(float64); ok {
			blink_phase = time.Duration(float64(time.Second) * p)
		}

		select {
		case ev := <-evchan:
			switch ev.Type {
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
		case <-time.After(blink_phase / 2):
			// Divided by two since we're only doing a simple toggle blink
			// TODO(q): Shouldn't redraw if blink is disabled...
		}
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		log4go.Exit(err)
	}

	var t tbfe
	t.visibleregion = make(map[*backend.View]primitives.Region)
	t.loop()
}
