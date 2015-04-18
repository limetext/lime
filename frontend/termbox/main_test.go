// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"github.com/limetext/lime/backend"
	"github.com/limetext/termbox-go"
	. "github.com/limetext/text"
	"testing"
	"time"
)

func TestPadLineRunes(t *testing.T) {
	var testPadData = []struct {
		line     int
		total    int
		expected []rune
	}{
		{1, 3, []rune{' ', ' ', '1', ' '}},
		{10, 3, []rune{' ', '1', '0', ' '}},
		{100, 3, []rune{'1', '0', '0', ' '}},
	}

	for _, p := range testPadData {
		runes := intToRunes(p.line)
		padded := padLineRunes(runes, p.total)

		for i, r := range p.expected {
			if r != padded[i] {
				t.Error("Expected runes to be padded")
			}
		}
	}
}

func TestIntToRunes(t *testing.T) {
	actual := intToRunes(100)
	expected := []rune{'1', '0', '0'}
	for i, r := range expected {
		if actual[i] != r {
			t.Error("Expected int to be transformed to runes")
		}
	}
}

func TestGetCaretStyle(t *testing.T) {
	var testcases = []struct {
		style    string
		inverse  bool
		expected termbox.Attribute
	}{
		{"", false, termbox.AttrUnderline},
		{"block", false, termbox.AttrReverse},
		{"", true, termbox.AttrReverse},
	}

	for _, tc := range testcases {
		if style := getCaretStyle(tc.style, tc.inverse); style != tc.expected {
			t.Errorf("Expected %v, got %v", tc.expected, style)
		}
	}
}

func TestUpdateVisibleRegion(t *testing.T) {
	defer backend.DestroyEditor()
	var (
		fe tbfe
		e  = backend.GetEditor()
		w  = e.NewWindow()
		v  = w.NewFile()
	)

	fe.layout = make(map[*backend.View]layout)
	fe.layout[v] = layout{0, 0, 100, 100 - *consoleHeight - 1, Region{}, 0}
	fe.setupCallbacks(v)

	edit := v.BeginEdit()
	v.Insert(edit, 0, "foo")
	v.EndEdit(edit)

	if end := fe.layout[v].visible.End(); end != 3 {
		t.Fatalf("Expected 3, got %d", end)
	}
}

func TestCreateFrontend(t *testing.T) {
	defer backend.DestroyEditor()
	var frontend *tbfe
	frontendWasCreated := make(chan bool, 0)
	go func() {
		frontend = createFrontend()
		frontendWasCreated <- true
	}()

	select {
	case <-frontendWasCreated:
		break
	case <-time.After(2 * time.Second):
		t.Error("Frontend was not created within timeout")
	}

	frontend.lock.Lock()
	defer frontend.lock.Unlock()

	*showConsole = true
	if frontend.editor == nil {
		t.Error("Editor is nil")
	}

	if frontend.console == nil {
		t.Error("Current console is nil")
	}

	if frontend.currentWindow == nil {
		t.Error("Current window is nil")
	}

	if frontend.currentView == nil {
		t.Error("Current view is nil")
	}

	if _, ok := frontend.layout[frontend.currentView]; !ok {
		t.Error("Current view not in layout")
	}

	if _, ok := frontend.layout[frontend.console]; !ok {
		t.Error("Console view not in layout")
	}

	if len(frontend.layout) != 2 {
		t.Errorf("Layout has length %d, but should have length 2", len(frontend.layout))
	}
}

func TestLoopShutdown(t *testing.T) {
	defer backend.DestroyEditor()
	frontend := createFrontend()

	loopHasExited := make(chan bool, 0)
	go func() {
		frontend.loop()
		loopHasExited <- true
	}()

	frontend.shutdown <- true

	select {
	case <-loopHasExited:
		break
	case <-time.After(2 * time.Second):
		t.Error("Loop did not terminate within timeout")
	}
}

func TestHandleInput(t *testing.T) {
	frontend := createFrontend()

	event_a := termbox.Event{
		Type: termbox.EventKey,
		Ch:   'a',
	}
	event_b := termbox.Event{
		Type: termbox.EventKey,
		Ch:   'b',
	}
	expected := "ab"

	frontend.handleInput(event_a)
	frontend.handleInput(event_b)

	// Allow the event to be handled in the call chain
	time.Sleep(10 * time.Millisecond)

	buffer := frontend.currentView.Buffer()
	if substring := buffer.Substr(Region{A: 0, B: buffer.Size()}); substring != expected {
		t.Errorf("Expected %q to be in editor's buffer, but got %q.", expected, substring)
	}
}
