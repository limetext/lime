package backend

import (
	//	"bytes"
	"code.google.com/p/log4go"
	"fmt"
	//	"github.com/quarnster/parser"
	//	"io/ioutil"
	//	"lime/backend/loaders"
	. "github.com/quarnster/util/text"
	"lime/backend/parser"
	"lime/backend/render"
	"lime/backend/textmate"
	. "lime/backend/util"
	"reflect"
	"runtime/debug"
	//	"sort"
	"strings"
	"sync"
)

type (
	// A View provides a view into a specific underlying buffer
	// with its own set of selections, settings, viewport, etc.
	// Multiple Views can share the same underlying data buffer.
	View struct {
		HasSettings
		HasId
		name        string
		window      *Window
		buffer      Buffer
		selection   RegionSet
		undoStack   UndoStack
		scratch     bool
		overwrite   bool
		cursyntax   string
		syntax      parser.SyntaxHighlighter
		regions     render.ViewRegionMap
		editstack   []*Edit
		lock        sync.Mutex
		modified    bool
		reparseChan chan parseReq
	}
	parseReq struct {
		forced bool
	}

	// The Edit object is an internal type passed as an argument
	// to a TextCommand. All text operations need to be associated
	// with a valid Edit object.
	Edit struct {
		invalid    bool
		composite  CompositeAction
		savedSel   RegionSet
		savedCount int
		command    string
		args       Args
		v          *View
		bypassUndo bool
	}
)

func newView(w *Window) *View {
	ret := &View{window: w, regions: make(render.ViewRegionMap)}
	ret.Settings().AddOnChange("lime.view.syntax", func() {
		ret.lock.Lock()
		defer ret.lock.Unlock()
		syn, _ := ret.Settings().Get("syntax", "").(string)
		if syn != ret.cursyntax {
			ret.cursyntax = syn
			defer ret.reparse(true)
		}
	})

	ret.reparseChan = make(chan parseReq, 32)
	go ret.parsethread()
	ret.Settings().Set("is_widget", false)
	return ret
}

func newEdit(v *View) *Edit {
	ret := &Edit{
		v:          v,
		savedCount: v.buffer.ChangeCount(),
	}
	for _, r := range v.Sel().Regions() {
		ret.savedSel.Add(r)
	}
	return ret
}

func (e *Edit) String() string {
	return fmt.Sprintf("%s: %v, %v, %v", e.command, e.args, e.bypassUndo, e.composite)
}

func (e *Edit) Apply() {
	e.composite.Apply()
}

func (e *Edit) Undo() {
	e.composite.Undo()
	e.v.Sel().Clear()
	for _, r := range e.savedSel.Regions() {
		e.v.Sel().Add(r)
	}
}

func (v *View) setBuffer(b Buffer) error {
	if v.buffer != nil {
		return fmt.Errorf("There is already a buffer set")
	}
	v.buffer = b
	// TODO(q): Dynamically load the correct syntax file
	b.AddCallback(func(_ Buffer, position, delta int) {
		v.flush(position, delta)
	})
	return nil
}

func (v *View) flush(position, delta int) {
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		e := Prof.Enter("view.flush")
		defer e.Exit()
		v.selection.Adjust(position, delta)
		if v.syntax != nil {
			v.syntax.Adjust(position, delta)
		}
		for k, v2 := range v.regions {
			v2.Regions.Adjust(position, delta)
			v.regions[k] = v2
		}
	}()
	OnModified.Call(v)
	v.reparse(false)
}

func (v *View) parsethread() {
	pc := 0
	lastParse := -1
	doparse := func() {
		p := Prof.Enter("syntax.parse")
		defer p.Exit()
		defer func() {
			if r := recover(); r != nil {
				log4go.Error("Panic in parse thread: %v\n%s", r, string(debug.Stack()))
				if pc > 0 {
					panic(r)
				}
				pc++
			}
		}()
		b := v.Buffer()
		b.Lock()
		sub := b.Substr(Region{0, b.Size()})
		b.Unlock()
		source, _ := v.Settings().Get("syntax", "").(string)
		if len(source) != 0 {
			// TODO
			if pr, err := textmate.NewLanguageParser(source, sub); err != nil {
				log4go.Error("Couldn't parse: %v", err)
			} else if syn, err := parser.NewSyntaxHighlighter(pr); err != nil {
				log4go.Error("Couldn't create syntaxhighlighter: %v", err)
			} else {
				v.lock.Lock()
				defer v.lock.Unlock()
				// Only set if it isn't invalid already, otherwise the
				// current syntax highlighting will be more accurate
				// as it will have had incremental adjustments done to it
				if v.buffer.ChangeCount() == lastParse {
					v.syntax = syn
					for k := range v.regions {
						if strings.HasPrefix(k, "lime.syntax") {
							delete(v.regions, k)
						}
					}
					for k, v2 := range syn.Flatten() {
						if v2.Regions.HasNonEmpty() {
							v.regions[k] = v2
						}
					}
				}
			}
		}
	}
	for pr := range v.reparseChan {
		if cc := v.buffer.ChangeCount(); lastParse != cc || pr.forced {
			lastParse = cc
			doparse()
			v.Settings().Set("lime.syntax.updated", lastParse)
		}
	}
}

func (v *View) reparse(forced bool) {
	if len(v.reparseChan) < cap(v.reparseChan) || forced {
		v.reparseChan <- parseReq{forced}
	}
}

func (v *View) ScopeName(point int) string {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.syntax != nil {
		return v.syntax.ScopeName(point)
	}
	return ""
}

func (v *View) ExtractScope(point int) Region {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.syntax != nil {
		return v.syntax.ScopeExtent(point)
	}
	return Region{}
}

func (v *View) ScoreSelector(point int, selector string) int {
	if sn := v.syntax.ScopeName(point); len(sn) > 0 {
		return 1 + strings.Index(sn, selector)
	}
	return 0
}

func (v *View) Sel() *RegionSet {
	return &v.selection
}

func (v *View) Window() *Window {
	return v.window
}

func (v *View) Buffer() Buffer {
	return v.buffer
}

func (v *View) Insert(edit *Edit, point int, value string) int {
	if t, ok := v.Settings().Get("translate_tabs_to_spaces", false).(bool); ok && t && strings.Contains(value, "\t") {
		tab_size, ok := v.Settings().Get("tab_size", 4).(int)
		if !ok {
			tab_size = 4
		}

		lines := strings.Split(value, "\n")
		for i, li := range lines {
			for {
				if idx := strings.Index(li, "\t"); idx != -1 {
					ai := idx
					if i == 0 {
						_, col := v.buffer.RowCol(point)
						ai = col + 1
					}
					add := 1 + ((ai + (tab_size - 1)) &^ (tab_size - 1))
					spaces := ""
					for j := ai; j < add; j++ {
						spaces += " "
					}
					li = li[:idx] + spaces + li[idx+1:]
					continue
				}
				break
			}
			lines[i] = li
		}
		value = strings.Join(lines, "\n")
	}
	edit.composite.AddExec(NewInsertAction(v.buffer, point, value))
	return len(value)
}

func (v *View) Erase(edit *Edit, r Region) {
	edit.composite.AddExec(NewEraseAction(v.buffer, r))
}

func (v *View) Replace(edit *Edit, r Region, value string) {
	edit.composite.AddExec(NewReplaceAction(v.buffer, r, value))
}

func (v *View) BeginEdit() *Edit {
	e := newEdit(v)
	v.editstack = append(v.editstack, e)
	return e
}

func (v *View) EndEdit(e *Edit) {
	if e.invalid {
		log4go.Fine("This edit has already been invalidated: %v, %v", e, v.editstack)
		return
	}
	i := len(v.editstack) - 1
	for i := len(v.editstack) - 1; i >= 0; i-- {
		if v.editstack[i] == e {
			break
		}
	}
	if i == -1 {
		log4go.Error("This edit isn't even in the stack... where did it come from? %v, %v", e, v.editstack)
		return
	}

	var selmod bool

	if l := len(v.editstack) - 1; i != l {
		log4go.Error("This edit wasn't last in the stack... %d !=  %d: %v, %v", i, l, e, v.editstack)
	}
	for j := len(v.editstack) - 1; j >= i; j-- {
		ce := v.editstack[j]
		ce.invalid = true
		is := reflect.DeepEqual(*v.Sel(), ce.savedSel)
		ib := v.buffer.ChangeCount() == ce.savedCount
		eq := (is && ib && ce.composite.Len() == 0)
		if !eq && is {
			selmod = true
		}
		if !v.scratch && !ce.bypassUndo && !eq {
			if i == 0 || j != i {
				// Presume someone forgot to add it in the j != i case
				v.undoStack.Add(e)
			} else {
				// This edit belongs to another edit
				v.editstack[i-1].composite.Add(ce)
			}
		}
	}
	v.editstack = v.editstack[:i]
	if selmod {
		OnSelectionModified.Call(v)
	}
}

func (v *View) SetScratch(s bool) {
	v.scratch = s
}

func (v *View) IsScratch() bool {
	return v.scratch
}

func (v *View) OverwriteStatus() bool {
	return v.overwrite
}

func (v *View) SetOverwriteStatus(s bool) {
	v.overwrite = s
}

func (v *View) CommandHistory(idx int, modifying_only bool) (name string, args Args, count int) {
	// TODO: merge history when possible
	if i := v.undoStack.index(idx, modifying_only); i != -1 {
		e := v.undoStack.actions[i]
		return e.command, e.args, 1
	}
	return "", nil, 0
}

func (v *View) runCommand(cmd TextCommand, name string) error {
	e := v.BeginEdit()
	e.command = name
	//	e.args = args
	e.bypassUndo = cmd.BypassUndo()

	defer func() {
		v.EndEdit(e)
		if r := recover(); r != nil {
			log4go.Error("Paniced while running text command %s %v: %v\n%s", name, cmd, r, string(debug.Stack()))
		}
	}()
	p := Prof.Enter("view.cmd." + name)
	defer p.Exit()
	return cmd.Run(v, e)
}

func (v *View) AddRegions(key string, regions []Region, scope, icon string, flags render.ViewRegionFlags) {
	vr := render.ViewRegions{Scope: scope, Icon: icon, Flags: flags}
	vr.Regions.AddAll(regions)

	v.lock.Lock()
	defer v.lock.Unlock()
	v.regions[key] = vr
}

func (v *View) GetRegions(key string) (ret []Region) {
	v.lock.Lock()
	defer v.lock.Unlock()
	vr := v.regions[key]
	rs := vr.Regions.Regions()
	ret = make([]Region, len(rs))
	copy(ret, rs)
	return
}

func (v *View) EraseRegions(key string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	delete(v.regions, key)
}

func (v *View) UndoStack() *UndoStack {
	return &v.undoStack
}

func (v *View) Transform(scheme render.ColourScheme, viewport Region) render.Recipe {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.syntax == nil {
		return nil
	}
	rr := make(render.ViewRegionMap)
	for k, v := range v.regions {
		rr[k] = v.Clone()
	}
	return render.Transform(scheme, rr, viewport)
}
