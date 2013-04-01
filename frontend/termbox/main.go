package main

import (
	"code.google.com/p/log4go"
	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		log4go.Exit(err)
	}
	defer termbox.Close()
	for {
		ev := termbox.PollEvent()
		log4go.Info("ev: %v, %d", ev, ev.Key == termbox.KeyEsc)
		if ev.Key == termbox.KeyEsc {
			break
		}
	}
}
