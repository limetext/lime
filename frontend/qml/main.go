// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"code.google.com/p/log4go"
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	_ "github.com/limetext/lime/backend/commands"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/lime/backend/sublime"
	"github.com/limetext/lime/backend/textmate"
	"github.com/limetext/lime/backend/util"
	"github.com/niemeyer/qml"
	. "github.com/quarnster/util/text"
	"image/color"
	"io"
	"sync"
)

var (
	scheme *textmate.Theme
	blink  bool
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
	dirty          bool
	Len            int
}

func htmlcol(c render.Colour) string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
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

func (t *tbfe) FormatLine(v *backend.View, line int) string {
	buf := bytes.NewBuffer(nil)
	vr := v.Buffer().Line(v.Buffer().TextPoint(line, 0))
	log4go.Debug("FormatLine: %d, %s", line, vr)
	if vr.Size() == 0 {
		return ""
	}
	recipie := v.Transform(scheme, vr).Transcribe()
	highlight_line := false
	if b, ok := v.Settings().Get("highlight_line", highlight_line).(bool); ok {
		highlight_line = b
	}
	lastEnd := vr.Begin()
	for _, reg := range recipie {
		if lastEnd != reg.Region.Begin() {
			fmt.Fprintf(buf, "<span>%s</span>", v.Buffer().Substr(Region{lastEnd, reg.Region.Begin()}))
		}
		fmt.Fprintf(buf, "<span style=\"white-space:pre; color:#%s; background:#%s\">%s</span>", htmlcol(reg.Flavour.Foreground), htmlcol(reg.Flavour.Background), v.Buffer().Substr(reg.Region))
		lastEnd = reg.Region.End()
	}
	if lastEnd != vr.End() {
		io.WriteString(buf, v.Buffer().Substr(Region{lastEnd, vr.End()}))
	}
	return buf.String()
}

func (t *tbfe) DefaultBg() color.RGBA {
	c := scheme.Spice(&render.ViewRegions{})
	c.Background.A = 0xff
	return color.RGBA(c.Background)
}

func (t *tbfe) DefaultFg() color.RGBA {
	c := scheme.Spice(&render.ViewRegions{})
	c.Foreground.A = 0xff
	return color.RGBA(c.Foreground)
}

func (t *tbfe) loop() {
	qml.Init(nil)
	engine := qml.NewEngine()

	engine.Context().SetVar("lines", t)
	engine.Context().SetVar("frontend", t)
	engine.Context().SetVar("editor", backend.GetEditor())

	backend.OnNew.Add(func(v *backend.View) {
		v.Settings().AddOnChange("lime.frontend.html.render", func(name string) { t.dirty = true })
	})
	backend.OnModified.Add(func(v *backend.View) {
		t.dirty = true
	})
	backend.OnSelectionModified.Add(func(v *backend.View) {
		t.dirty = true
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

	defer func() {
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

	{
		w, h := 800, 600
		t.lock.Lock()
		t.layout[v] = layout{0, 0, w, h - console_height - 1, Region{}, 0}
		t.layout[c] = layout{0, h - console_height + 1, w, console_height - 5, Region{}, 0}
		t.lock.Unlock()
		t.Show(v, Region{1, 1})
	}
	t.Show(v, Region{100, 100})
	t.Show(v, Region{1, 1})
	t.Len, _ = v.Buffer().RowCol(v.Buffer().Size())

	ed.Init()
	sublime.Init()

	component, err := engine.LoadFile("main.qml")
	if err != nil {
		log4go.Exit(err)
	}
	window := component.CreateWindow(nil)
	window.Show()
	qml.Changed(t, &t.Len)

	log4go.Debug("Done")
	window.Wait()
}

func main() {
	log4go.AddFilter("file", log4go.FINEST, log4go.NewConsoleLogWriter())
	defer func() {
		py.NewLock()
		py.Finalize()
	}()

	var t tbfe
	t.dorender = make(chan bool, render_chan_len)
	t.layout = make(map[*backend.View]layout)
	t.loop()
}
