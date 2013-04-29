package main

import (
	"code.google.com/p/log4go"
	"github.com/salviati/go-qt5/qt5"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	_ "lime/backend/commands"
	// "lime/backend/loaders"
	"lime/backend/primitives"
	// "lime/backend/sublime"
	"image/color"
	"lime/backend/textmate"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	_      = primitives.Region{}
	_      = color.RGBA{}
	wnds   = make(map[*backend.Window]QLimeWindow)
	scheme *textmate.Theme
)

type (
	QLimeWindow struct {
		tw *qt5.TabWidget
		w  *backend.Window
	}
	QLimeView struct {
		*qt5.Widget
		v *backend.View
	}
)

func newQLimeView(v *backend.View) *QLimeView {
	log4go.Debug("new_QLimeView entered")
	defer log4go.Debug("new_QLimeView exited")
	var ret QLimeView
	ret.Widget = qt5.NewWidget()
	ret.v = v
	ret.Widget.OnPaintEvent(func(ev *qt5.PaintEvent) {
		p := qt5.NewPainter()
		defer p.Close()
		p.Begin(ret)
		b := v.Buffer()
		ps := p.Font().PointSize()

		pen := qt5.NewPen()
		p.SetPen(pen)
		brush := qt5.NewBrush()
		brush.SetStyle(qt5.SolidPattern)
		def := scheme.Settings[0]
		p.SetBrush(brush)
		f := p.Font()
		f.SetFixedPitch(true)
		p.SetFont(f)

		brush.SetColor(color.RGBA(def.Settings["background"]))
		p.DrawRect(ev.Rect())
		is_widget, ok := v.Settings().Get("is_widget", false).(bool)
		is_widget = ok && is_widget
		pen.SetColor(color.RGBA(def.Settings["background"]))
		p.SetPen(pen)

		for y := 0; y < 20; y++ {
			pos := b.TextPoint(y, 0)
			line := b.Line(pos)

			if is_widget {
				p.DrawText(qt5.Point{0, (y + 1) * (ps + 2)}, b.Substr(line))
			} else {
				for line.Contains(pos) {
					scope := primitives.Region{pos, pos}
					sn := v.ScopeName(pos)
					for line.Contains(pos) {
						pos++
						if v.ScopeName(pos) != sn {
							scope.B = pos
							break
						}
					}
					is := line.Intersection(scope)
					c := color.RGBA(def.Settings["foreground"])
					s := scheme.ClosestMatchingSetting(sn)
					if v, ok := s.Settings["foreground"]; ok {
						c = color.RGBA(v)
					}
					pen.SetColor(c)
					p.SetPen(pen)
					_, col := b.RowCol(line.A)
					p.DrawText(qt5.Point{col * ps / 2, (y + 1) * (ps + 2)}, b.Substr(is))
					line.A = is.End()
				}
			}
		}
	})
	ret.Widget.OnResizeEvent(func(ev *qt5.ResizeEvent) {
		if w, ok := v.Settings().Get("is_widget", false).(bool); ok && !w {
			ret.Widget.SetMinimumSize(qt5.Size{600, 100})
		}
	})
	v.Settings().Set("lime.qt.widget", &ret)
	return &ret
}

func new_window(w *backend.Window) {
	log4go.Debug("new_window entered")
	defer log4go.Debug("new_window exited")
	qw := qt5.NewWidget()
	qw.Show()
	qw.SetSizev(600, 400)
	tw := qt5.NewTabWidget()
	lbox := qt5.NewVBoxLayout()
	lbox.AddWidget(tw)
	c := newQLimeView(backend.GetEditor().Console())
	sa := qt5.NewScrollArea()
	sa.SetWidget(c)
	lbox.AddWidget(sa)
	qw.SetLayout(lbox)
	wnds[w] = QLimeWindow{tw, w}
}

func new_view(v *backend.View) {
	log4go.Debug("new_view entered")
	defer log4go.Debug("new_view exited")
	qw := wnds[v.Window()]
	w := newQLimeView(v)
	v.Settings().Set("syntax", "../../3rdparty/bundles/GoSublime/GoSublime.tmLanguage")

	w.SetSizev(600, 400)
	// w := qt5.NewWidget()
	sa := qt5.NewScrollArea()
	sa.SetWidget(w)
	qw.tw.AddTab(sa, v.Buffer().Name(), nil)
}

func view_modified(v *backend.View) {

	// 	if w, ok := v.Settings().Get("lime.qt.widget", nil).(*qt5.LineEdit); ok {
	// 		w.SetText(v.Buffer().String())
	// 	}
}

func main() {
	py.InitializeEx(false)
	defer py.Finalize()
	e := backend.GetEditor()
	log4go.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

	if sc, err := textmate.LoadTheme("../../3rdparty/bundles/TextMate-Themes/GlitterBomb.tmTheme"); err != nil {
		log4go.Error(err)
	} else {
		scheme = sc
		log4go.Debug("scheme: %v", scheme)
	}

	backend.OnNewWindow.Add(new_window)
	backend.OnNew.Add(new_view)
	backend.OnModified.Add(view_modified)
	go e.Init()
	qt5.Main(func() {
		w := e.NewWindow()
		w.OpenFile("main.go", 0)
		w.OpenFile("../../backend/editor.go", 0)
		qt5.Run()
	})
}
