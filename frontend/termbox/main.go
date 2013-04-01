package main

import (
	"code.google.com/p/log4go"
	"github.com/nsf/termbox-go"
	"lime/backend"
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

func main() {
	if err := termbox.Init(); err != nil {
		log4go.Exit(err)
	}
	defer termbox.Close()
	e := backend.GetEditor()
	e.LogInput(true)
	e.LogCommands(true)
	for {
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
			e.HandleInput(kp)
		}
	}
}
