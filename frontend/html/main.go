package main

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/quarnster/util/text"
	"io"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	_ "lime/backend/commands"
	"lime/backend/render"
	"lime/backend/sublime"
	"lime/backend/textmate"
	"lime/backend/util"
	"net/http"
	"runtime/debug"
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
}

func htmlcol(c render.Colour) string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

func (t *tbfe) renderView(wr io.Writer, v *backend.View, lay layout) {
	p := util.Prof.Enter("render")
	defer p.Exit()

	vr := lay.visible
	runes := v.Buffer().Substr(vr)

	recipie := v.Transform(scheme, vr).Transcribe()
	highlight_line := false
	if b, ok := v.Settings().Get("highlight_line", highlight_line).(bool); ok {
		highlight_line = b
	}

	lastEnd := 0
	for _, reg := range recipie {
		if lastEnd != reg.Region.Begin() {
			io.WriteString(wr, runes[lastEnd:reg.Region.Begin()])
		}
		fmt.Fprintf(wr, "<span style=\"color:#%s; background-color:#%s\">%s</span>", htmlcol(reg.Flavour.Foreground), htmlcol(reg.Flavour.Background), runes[reg.Region.Begin():reg.Region.End()])
		lastEnd = reg.Region.End()
	}
	if lastEnd != vr.End() {
		io.WriteString(wr, v.Buffer().Substr(Region{lastEnd, vr.End()}))
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
	if !l.visible.Covers(r) {
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

var pc = 0

func (t *tbfe) render(w io.Writer) {
	defer func() {
		if r := recover(); r != nil {
			log4go.Error("Panic in renderthread: %v\n%s", r, string(debug.Stack()))
			if pc > 1 {
				panic(r)
			}
			pc++
		}
	}()
	vs := make([]*backend.View, 0, len(t.layout))
	l := make([]layout, 0, len(t.layout))
	for k, v := range t.layout {
		vs = append(vs, k)
		l = append(l, v)
	}
	for i, v := range vs {
		t.renderView(w, v, l[i])
	}
	//	runes := []rune(t.status_message)
}
func (t *tbfe) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	log4go.Debug("Serving client: %s", req)

	c := scheme.Spice(&render.ViewRegions{})

	fmt.Fprintf(w, "<html><body style=\"white-space:pre; color:#%s; background-color:#%s\">", htmlcol(c.Foreground), htmlcol(c.Background))
	t.render(w)
	io.WriteString(w, "</body></html>")
	log4go.Debug("Done serving client")
}

func (t *tbfe) loop() {
	// backend.OnNew.Add(func(v *backend.View) {
	// 	v.Settings().AddOnChange("lime.frontend.html.render", func() { t.render() })
	// })
	// backend.OnModified.Add(func(v *backend.View) {
	// 	t.render()
	// })
	// backend.OnSelectionModified.Add(func(v *backend.View) {
	// 	t.render()
	// })

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
	//	end := v.Buffer().Size() - 2
	sel.Add(Region{0, 0})
	// sel.Add(Region{end - 22, end - 22})
	// sel.Add(Region{end - 16, end - 20})
	// sel.Add(Region{end - 13, end - 10})

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

	go func() {
		ed.Init()
		sublime.Init()
	}()
	log4go.Debug("serving")
	if err := http.ListenAndServe("localhost:8080", t); err != nil {
		log4go.Error("Error serving: %s", err)
	}
	log4go.Debug("Done")
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
