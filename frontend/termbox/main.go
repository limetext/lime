package main

import (
	"code.google.com/p/log4go"
	"fmt"
	"io/ioutil"
	"lime/3rdparty/libs/termbox-go"
	"lime/backend"
	"lime/backend/loaders"
	"lime/backend/primitives"
	"lime/backend/textmate"
	"strings"
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
)

func renderView(sx, sy, w, h int, v *backend.View) {
	sel := v.Sel()
	substr := v.Substr(primitives.Region{0, v.Size()})
	lines := strings.Split(substr, "\n")
	s, e := 0, len(lines)
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
		lfg, lbg  termbox.Attribute
	)

	tab_size, ok := v.Settings().Get("tab_size", 4).(int)
	if !ok {
		tab_size = 4
	}

	for i := range runes {
		sub2 += string(runes[i])
		if runes[i] == '\n' {
			x = sx
			y++
			if y > ey {
				break
			}
			continue
		} else if runes[i] == '\t' {
			add := (x + 1 + (tab_size - 1)) &^ (tab_size - 1)
			for x < add {
				termbox.SetCell(x, y, ' ', lfg, lbg)
				x++
			}
			continue
		}

		if x < ex {
			fg, bg := defaultFg, defaultBg
			o := off + len(sub2)
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
			// TODO(q): It should differ between a proper selection and just the cursor position
			r := primitives.Region{o, o}
			if sel.Contains(r) {
				// TODO(q): It should use the correct colors for selections and the cursor
				bg = 3
			}
			termbox.SetCell(x, y, runes[i], fg, bg)
		}
		x++
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		log4go.Exit(err)
	}

	ed := backend.GetEditor()
	ed.LogInput(true)
	ed.LogCommands(true)
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
		fmt.Println(c.Substr(primitives.Region{0, c.Size()}))
	}()

	w := ed.NewWindow()
	v := w.OpenFile("main.go", 0)

	if err := v.SetSyntaxFile("../../3rdparty/bundles/GoSublime/GoSublime.tmLanguage"); err != nil {
		log4go.Error("Unable to set syntax file: %s", err)
	}
	sel := v.Sel()
	sel.Clear()
	end := v.Buffer().Size() - 2
	sel.Add(primitives.Region{end - 24, end - 24})
	sel.Add(primitives.Region{end - 22, end - 22})
	sel.Add(primitives.Region{end - 20, end - 20})
	for {
		termbox.Clear(defaultFg, defaultBg)
		w, h := termbox.Size()
		renderView(0, 0, w, h-3, v)
		renderView(0, h-3, w, 3, c)

		termbox.Flush()
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			var kp backend.KeyPress

			if ev.Ch != 0 {
				kp.Key = backend.Key(ev.Ch)
			} else if v2, ok := lut[ev.Key]; ok {
				kp = v2
			}

			if ev.Key == termbox.KeyEsc {
				return
			}
			ed.HandleInput(kp)
		}
	}
}
