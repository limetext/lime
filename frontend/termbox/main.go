package main

import (
	"code.google.com/p/log4go"
	"fmt"
	"lime/3rdparty/libs/termbox-go"
	"lime/backend"
	"lime/backend/primitives"
	"strings"
)

var lut = map[termbox.Key]backend.KeyPress{
	termbox.KeyCtrlA: backend.KeyPress{Ctrl: true, Key: 'a'},
	termbox.KeyCtrlB: backend.KeyPress{Ctrl: true, Key: 'b'},
	termbox.KeyCtrlC: backend.KeyPress{Ctrl: true, Key: 'c'},
	termbox.KeyCtrlD: backend.KeyPress{Ctrl: true, Key: 'd'},
	termbox.KeyCtrlE: backend.KeyPress{Ctrl: true, Key: 'e'},
	termbox.KeyCtrlF: backend.KeyPress{Ctrl: true, Key: 'f'},
	termbox.KeyCtrlG: backend.KeyPress{Ctrl: true, Key: 'g'},
	termbox.KeyCtrlH: backend.KeyPress{Ctrl: true, Key: 'h'},
	termbox.KeyCtrlI: backend.KeyPress{Ctrl: true, Key: 'i'},
	termbox.KeyCtrlJ: backend.KeyPress{Ctrl: true, Key: 'j'},
	termbox.KeyCtrlK: backend.KeyPress{Ctrl: true, Key: 'k'},
	termbox.KeyCtrlL: backend.KeyPress{Ctrl: true, Key: 'l'},
	termbox.KeyCtrlM: backend.KeyPress{Ctrl: true, Key: 'm'},
	termbox.KeyCtrlN: backend.KeyPress{Ctrl: true, Key: 'n'},
	termbox.KeyCtrlO: backend.KeyPress{Ctrl: true, Key: 'o'},
	termbox.KeyCtrlP: backend.KeyPress{Ctrl: true, Key: 'p'},
	termbox.KeyCtrlQ: backend.KeyPress{Ctrl: true, Key: 'q'},
	termbox.KeyCtrlR: backend.KeyPress{Ctrl: true, Key: 'r'},
	termbox.KeyCtrlS: backend.KeyPress{Ctrl: true, Key: 's'},
	termbox.KeyCtrlT: backend.KeyPress{Ctrl: true, Key: 't'},
	termbox.KeyCtrlU: backend.KeyPress{Ctrl: true, Key: 'u'},
	termbox.KeyCtrlV: backend.KeyPress{Ctrl: true, Key: 'v'},
	termbox.KeyCtrlW: backend.KeyPress{Ctrl: true, Key: 'w'},
	termbox.KeyCtrlX: backend.KeyPress{Ctrl: true, Key: 'x'},
	termbox.KeyCtrlY: backend.KeyPress{Ctrl: true, Key: 'y'},
	termbox.KeyCtrlZ: backend.KeyPress{Ctrl: true, Key: 'z'},
	termbox.KeyCtrl2: backend.KeyPress{Ctrl: true, Key: '2'},
	termbox.KeyCtrl3: backend.KeyPress{Ctrl: true, Key: '3'},
	termbox.KeyCtrl4: backend.KeyPress{Ctrl: true, Key: '4'},
	termbox.KeyCtrl5: backend.KeyPress{Ctrl: true, Key: '5'},
	termbox.KeyCtrl6: backend.KeyPress{Ctrl: true, Key: '6'},
	termbox.KeyCtrl7: backend.KeyPress{Ctrl: true, Key: '7'},
	termbox.KeyCtrl8: backend.KeyPress{Ctrl: true, Key: '8'},
}

func renderView(sx, sy, w, h int, v *backend.View) {
	substr := v.Substr(primitives.Region{0, v.Size()})
	lines := strings.Split(substr, "\n")
	s, e := 0, len(lines)
	if e > 1 {
		e = e - 1
		if e > h {
			s = e - h
		}
	}
	lines = lines[s:e]
	runes := []rune(strings.Join(lines, "\n"))
	x, y := sx, sy
	ex, ey := sx+w, sy+h
	for i := range runes {
		if runes[i] == '\n' {
			x = sx
			y++
			if y > ey {
				break
			}
			continue
		}
		if x < ex {
			termbox.SetCell(x, y, runes[i], 1, 0)
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
	if err := termbox.SetColorMode(termbox.ColorMode256); err != nil {
		log4go.Error("Unable to use 256 color mode: %s", err)
	} else {
		log4go.Debug("Using 256 color mode")
		pal := make([]termbox.RGB, 256)
		pal[1] = termbox.RGB{255, 0, 255}
		termbox.SetColorPalette(pal)
	}
	defer func() {
		termbox.Close()
		fmt.Println(c.Substr(primitives.Region{0, c.Size()}))
	}()
	w := ed.NewWindow()
	v := w.OpenFile("main.go", 0)

	for {
		termbox.Clear(0, 0)
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
			} else if v, ok := lut[ev.Key]; ok {
				kp = v
			}
			if ev.Key == termbox.KeyEsc {
				return
			}
			ed.HandleInput(kp)
		}
	}
}
