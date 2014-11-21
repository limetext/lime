// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/limetext/gopy/lib"
	"github.com/limetext/lime/backend"
	_ "github.com/limetext/lime/backend/commands"
	"github.com/limetext/lime/backend/keys"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/lime/backend/sublime"
	"github.com/limetext/lime/backend/textmate"
	"github.com/limetext/lime/backend/util"
	. "github.com/limetext/text"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	scheme *textmate.Theme
	blink  bool
	port   = flag.Int("port", 8080, "Configures which port to host lime on")
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
}

func (t *tbfe) Erased(changed_buffer Buffer, region_removed Region, data_removed []rune) {
	t.scroll(changed_buffer)
}

func (t *tbfe) Inserted(changed_buffer Buffer, region_inserted Region, data_inserted []rune) {
	t.scroll(changed_buffer)
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

	t.BroadcastData(map[string]interface{}{"type": "statusMessage", "msg": msg})
}

func (t *tbfe) ErrorMessage(msg string) {
	log.Error(msg)

	t.BroadcastData(map[string]interface{}{"type": "errorMessage", "msg": msg})
}

func (t *tbfe) MessageDialog(msg string) {
	log.Info(msg)

	t.BroadcastData(map[string]interface{}{"type": "messageDialog", "msg": msg})
}

// TODO: wait for client response, return true/false
func (t *tbfe) OkCancelDialog(msg, ok string) bool {
	log.Info(msg, ok)

	t.BroadcastData(map[string]interface{}{"type": "okCancelDialog", "msg": msg, "ok": ok})

	return false
}

func (t *tbfe) scroll(b Buffer) {
	t.Show(backend.GetEditor().Console(), Region{b.Size(), b.Size()})
}

var pc = 0

func (t *tbfe) render(w io.Writer) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic in renderthread: %v\n%s", r, string(debug.Stack()))
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
func (t *tbfe) key(w http.ResponseWriter, req *http.Request) {
	log.Debug("key: %s", req)
	kc := req.FormValue("keyCode")
	var kp keys.KeyPress
	v, _ := strconv.ParseInt(kc, 10, 32)

	if req.FormValue("altKey") == "true" {
		kp.Alt = true
	}
	if req.FormValue("ctrlKey") == "true" {
		kp.Ctrl = true
	}
	if req.FormValue("metaKey") == "true" {
		kp.Super = true
	}
	if req.FormValue("shiftKey") == "true" {
		kp.Shift = true
	}
	if !kp.Shift {
		v = int64(unicode.ToLower(rune(v)))
	}
	kp.Key = keys.Key(v)
	backend.GetEditor().HandleInput(kp)
}

func (t *tbfe) view(w http.ResponseWriter, req *http.Request) {
	log.Debug("view: %s", req)
	if t.dirty {
		t.dirty = false
		t.render(w)
	} else {
		w.WriteHeader(404)
	}
}

func (t *tbfe) theme(w http.ResponseWriter, req *http.Request) {
	log.Debug("theme: %s", req)

	reqpath, _ := url.QueryUnescape(req.RequestURI)

	// Make sure the URL starts with "/themes/"
	// Don't allow ".." in URLs
	if !strings.HasPrefix(reqpath, "/themes/") || strings.Index(reqpath, "..") != -1 {
		w.WriteHeader(404)
		return
	}

	filepath := path.Join("../../packages", reqpath)

	exists := false
	if s, err := os.Stat(filepath); err == nil {
		if !s.IsDir() {
			exists = true
		}
	}

	if exists {
		fi, err := os.Open(filepath)
		if err != nil {
			w.WriteHeader(500)
			log.Error(err)
			return
		}

		defer fi.Close()

		io.Copy(w, fi)
	} else {
		w.WriteHeader(404)
	}
}

func (t *tbfe) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := time.Now()
	w.Header().Set("Content-Type", "text/html")
	log.Debug("Serving client: %s", req)

	c := scheme.Spice(&render.ViewRegions{})

	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(404)
		panic(err)
	}

	r := strings.NewReplacer("{{foregroundColor}}", htmlcol(c.Foreground), "{{backgroundColor}}", htmlcol(c.Background))
	r.WriteString(w, string(html))

	log.Debug("Done serving client: %s", time.Since(s))
}

var clients []*websocket.Conn

func (t *tbfe) WebsocketServer(ws *websocket.Conn) {
	clients = append(clients, ws)

	// Send status message
	if t.status_message != "" {
		websocket.JSON.Send(ws, map[string]string{"type": "statusMessage", "msg": t.status_message})
	}

	// Send cursor position
	websocket.JSON.Send(ws, t.GetSelectionMessage(backend.GetEditor().ActiveWindow().ActiveView()))

	// Send editor content
	var buf bytes.Buffer
	t.render(bufio.NewWriter(&buf))
	websocket.Message.Send(ws, buf.Bytes())
	buf.Reset()

	var data map[string]interface{}
	var kp keys.KeyPress
	for {
		err := websocket.JSON.Receive(ws, &data)
		if err != nil {
			log.Error(err)
			return
		}
		//log.LogDebug("Received: %s", data)

		msgType := data["type"].(string)

		if msgType == "key" {
			kp.Alt = data["altKey"].(bool)
			kp.Ctrl = data["ctrlKey"].(bool)
			kp.Super = data["metaKey"].(bool)
			kp.Shift = data["shiftKey"].(bool)

			if keyName, ok := data["key"].(string); ok {
				if utf8.RuneCountInString(keyName) == 1 { // One char
					r, _ := utf8.DecodeRuneInString(keyName)
					kp.Key = keys.Key(int64(r))
				} else {
					// TODO: automatic lookup instead of this manual lookup
					// See https://github.com/limetext/lime/pull/421/files#r19269236
					keymap := map[string]keys.Key{
						"ArrowLeft":   keys.Left,
						"ArrowUp":     keys.Up,
						"ArrowRight":  keys.Right,
						"ArrowDown":   keys.Down,
						"Left":        keys.Left,
						"Up":          keys.Up,
						"Right":       keys.Right,
						"Down":        keys.Down,
						"Enter":       keys.Enter,
						"Escape":      keys.Escape,
						"Backspace":   keys.Backspace,
						"Delete":      keys.Delete,
						"Del":         keys.Delete, // Deprecated: some old browsers still use "Del" instead of "Delete"
						"KeypadEnter": keys.KeypadEnter,
						"F1":          keys.F1,
						"F2":          keys.F2,
						"F3":          keys.F3,
						"F4":          keys.F4,
						"F5":          keys.F5,
						"F6":          keys.F6,
						"F7":          keys.F7,
						"F8":          keys.F8,
						"F9":          keys.F9,
						"F10":         keys.F10,
						"F11":         keys.F11,
						"F12":         keys.F12,
						"Insert":      keys.Insert,
						"PageUp":      keys.PageUp,
						"PageDown":    keys.PageDown,
						"Home":        keys.Home,
						"End":         keys.End,
						"Break":       keys.Break,
					}

					if key, ok := keymap[keyName]; ok {
						kp.Key = key
					} else {
						log.Debug("Unknown key: %s", keyName)
						continue
					}
				}
			} else {
				v := int64(data["keyCode"].(float64))
				if !kp.Shift {
					v = int64(unicode.ToLower(rune(v)))
				}
				kp.Key = keys.Key(v)
			}

			backend.GetEditor().HandleInput(kp)
		} else if msgType == "command" {
			command := data["name"].(string)
			//args := data["args"].([]string) //TODO: add arguments support

			ed := backend.GetEditor()
			go ed.RunCommand(command, make(backend.Args))
		} else {
			log.Info("Unhandled message type: %s", msgType)
		}
	}
}

func (t *tbfe) BroadcastData(data map[string]interface{}) {
	for _, ws := range clients {
		websocket.JSON.Send(ws, data)
	}
}

func (t *tbfe) SetDirty() {
	t.dirty = true

	var buf bytes.Buffer
	t.render(bufio.NewWriter(&buf))
	for _, ws := range clients {
		websocket.Message.Send(ws, buf.Bytes())
	}
}

func (t *tbfe) GetSelectionMessage(v *backend.View) map[string]interface{} {
	return map[string]interface{}{
		"type": "selection",
		"sel":  v.Sel().Regions(),
	}
}

func (t *tbfe) loop() {
	backend.OnNew.Add(func(v *backend.View) {
		v.Settings().AddOnChange("lime.frontend.html.render", func(name string) {
			if name != "lime.syntax.updated" {
				return
			}
			t.SetDirty()
		})
	})
	// TODO: maybe not useful?
	/*backend.OnModified.Add(func(v *backend.View) {
		t.SetDirty()
	})*/
	backend.OnSelectionModified.Add(func(v *backend.View) {
		t.BroadcastData(t.GetSelectionMessage(v))
	})

	ed := backend.GetEditor()
	ed.SetFrontend(t)
	ed.LogInput(false)
	ed.LogCommands(false)
	c := ed.Console()
	if sc, err := textmate.LoadTheme("../../packages/themes/TextMate-Themes/Monokai.tmTheme"); err != nil {
		log.Error(err)
	} else {
		scheme = sc
	}

	defer func() {
		fmt.Println(util.Prof)
	}()

	w := ed.NewWindow()
	v := w.OpenFile("main.go", 0)
	//v.Settings().Set("trace", true)
	v.Settings().Set("syntax", "../../packages/go.tmbundle/Syntaxes/Go.tmLanguage")
	c.Buffer().AddObserver(t)

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
	log.Debug("Serving on port %d", *port)
	http.HandleFunc("/", t.ServeHTTP)
	http.HandleFunc("/view", t.view)
	http.HandleFunc("/key", t.key)
	http.HandleFunc("/themes/", t.theme)
	http.Handle("/ws", websocket.Handler(t.WebsocketServer))
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil); err != nil {
		log.Error("Error serving: %s", err)
	}
	log.Debug("Done")
}

func main() {
	flag.Parse()
	log.AddFilter("file", log.FINEST, log.NewConsoleLogWriter())
	defer func() {
		py.NewLock()
		py.Finalize()
	}()

	var t tbfe
	t.dorender = make(chan bool, render_chan_len)
	t.layout = make(map[*backend.View]layout)
	t.loop()
}
