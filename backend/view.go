// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"code.google.com/p/log4go"
	"fmt"
	. "github.com/quarnster/util/text"
	"lime/backend/parser"
	"lime/backend/render"
	"lime/backend/textmate"
	. "lime/backend/util"
	"reflect"
	"runtime/debug"
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
	//
	// Think of it a bit like an SQL transaction.
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
	ret.Settings().AddOnChange("lime.view.syntax", func(name string) {
		ret.lock.Lock()
		defer ret.lock.Unlock()
		if name != "syntax" {
			return
		}
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

// Returns a string describing this edit object. Should typically not be manually called.
func (e *Edit) String() string {
	return fmt.Sprintf("%s: %v, %v, %v", e.command, e.args, e.bypassUndo, e.composite)
}

// Applies the action of this Edit object. Should typically not be manually called.
func (e *Edit) Apply() {
	e.composite.Apply()
}

// Reverses the application of this edit object. Should typically not be manually called.
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

// flush is called every time the underlying buffer is changed.
// It calls "Adjust" on all the regions associated with this view,
// triggers the "OnModified" event, and adds a reparse request
// to the parse go-routine.
func (v *View) flush(position, delta int) {
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		e := Prof.Enter("view.flush")
		defer e.Exit()
		// TODO(.): issue #211
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

// parsethread would be the go-routine used for dealing with reparsing the
// current buffer when it has been modified. Each opened view has its own
// go-routine parsethread which sits idle and waits for requests to be sent
// on this view's reparseChan.
//
// The Buffer's ChangeCount, as well as the parse request's "forced" attribute
// is used to determined if a parse actually needs to happen or not.
//
// If it is decided that a reparse should take place, a snapshot of the buffer is
// taken and a parse is performed. Upon completion of this parse operation,
// and if the snapshot of the buffer has not already become outdated,
// then the regions of the view associated with syntax highlighting is updated.
//
// Changes made the buffer during the time when there is no accurate
// parse of the buffer is a monkey patched version of the old syntax highlighting
// regions, which in most instances will be accurate.
//
// See package lime/backend/parser for more details.
func (v *View) parsethread() {
	pc := 0
	lastParse := -1
	doparse := func() (ret bool) {
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
		if len(source) == 0 {
			return
		}
		// TODO: Allow other parsers instead of this hardcoded textmate version
		pr, err := textmate.NewLanguageParser(source, sub)
		if err != nil {
			log4go.Error("Couldn't parse: %v", err)
			return
		}
		syn, err := parser.NewSyntaxHighlighter(pr)
		if err != nil {
			log4go.Error("Couldn't create syntaxhighlighter: %v", err)
			return
		}
		v.lock.Lock()
		defer v.lock.Unlock()
		// Only set if it isn't invalid already, otherwise the
		// current syntax highlighting will be more accurate
		// as it will have had incremental adjustments done to it
		if v.buffer.ChangeCount() != lastParse {
			return
		}
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
		return true
	}
	for pr := range v.reparseChan {
		if cc := v.buffer.ChangeCount(); lastParse != cc || pr.forced {
			lastParse = cc
			if doparse() {
				v.Settings().Set("lime.syntax.updated", lastParse)
			}
		}
	}
}

// Send a reparse request via the reparse channel.
// If "forced" is set to true, then a reparse will be made
// even if the buffer appears to not have changed.
//
// The actual parsing is done in a separate go-routine, for which the
// "lime.syntax.updated" setting will be set once it has finished.
func (v *View) reparse(forced bool) {
	if len(v.reparseChan) < cap(v.reparseChan) || forced {
		v.reparseChan <- parseReq{forced}
	}
}

// Returns the full concatenated nested scope name at point.
// See package lime/backend/parser for details.
func (v *View) ScopeName(point int) string {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.syntax != nil {
		return v.syntax.ScopeName(point)
	}
	return ""
}

// Returns the Region of the inner most scope that contains "point".
// See package lime/backend/parser for details.
func (v *View) ExtractScope(point int) Region {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.syntax != nil {
		return v.syntax.ScopeExtent(point)
	}
	return Region{}
}

// ScoreSelector takes a point and a selector string and returns a score
// as to how good that specific selector matches the scope name at
// that point.
func (v *View) ScoreSelector(point int, selector string) int {
	// TODO(.): The algorithm to determine the score has not been compared
	// to what ST3 actually does. Not sure if any plugin I personally
	// care about uses this functionality, and if it does if it requires
	// specific scores to be returned.
	//
	// The overall implementation should be fine as a higher score is
	// returned the more specific a selector is due to the innermost
	// scope name being concatenated last in the string returned by ScopeName
	if sn := v.ScopeName(point); len(sn) > 0 {
		return 1 + strings.Index(sn, selector)
	}
	return 0
}

// Sel returns a pointer to the RegionSet used by this View
// to mark possibly multiple cursor positions and selection
// regions.
//
// Some quick notes:
// The actual cursor position is always in Region.B.
// Region{0,0} is a cursor at the start of the text (before any characters in the text).
//
// Region{0,1} has the cursor at position 1 (after the first character),
// but also selects/highlights the first character. In this instance Region.A = 0, Region.B = 1,
// Region.Start() returns 0 and Region.End() returns 1.
//
// Region{1,0} has the cursor at position 0 (before the first character),
// but also selects/highlights the first character. Think holding shift and pressing left on your keyboard.
// In this instance Region.A = 1, Region.B = 0, Region.Start() returns 0 and Region.End() returns 1.
//
func (v *View) Sel() *RegionSet {
	// BUG(.): Sometimes Sel becomes empty. There should always be at a minimum 1 valid cursor.
	return &v.selection
}

// Returns the window this view belongs to.
func (v *View) Window() *Window {
	return v.window
}

// Returns the underlying buffer that this view is a view into.
func (v *View) Buffer() Buffer {
	return v.buffer
}

// Inserts text at the given position in the provided edit object.
// Tabs are depending on the view's settings translated to spaces.
// The return value is the length of the string that was inserted.
func (v *View) Insert(edit *Edit, point int, value string) int {
	if t, ok := v.Settings().Get("translate_tabs_to_spaces", false).(bool); ok && t && strings.Contains(value, "\t") {
		tab_size, ok := v.Settings().Get("tab_size", 4).(int)
		if !ok {
			tab_size = 4
		}

		lines := strings.Split(value, "\n")
		for i, li := range lines {
			for {
				idx := strings.Index(li, "\t")
				if idx == -1 {
					break
				}
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
			}
			lines[i] = li
		}
		value = strings.Join(lines, "\n")
	}
	edit.composite.AddExec(NewInsertAction(v.buffer, point, value))
	// TODO(.): I think this should rather be the number of runes inserted?
	// The spec states that len() of a string returns the number of bytes,
	// which isn't very useful as all other buffer values are IIRC in runes.
	// http://golang.org/ref/spec#Length_and_capacity
	return len(value)
}

// Adds an Erase action of the given region to the provided edit object.
func (v *View) Erase(edit *Edit, r Region) {
	edit.composite.AddExec(NewEraseAction(v.buffer, r))
}

// Adds a Replace action of the given region to the provided edit object.
func (v *View) Replace(edit *Edit, r Region, value string) {
	edit.composite.AddExec(NewReplaceAction(v.buffer, r, value))
}

// Creates a new Edit object. Think of it a bit like starting an SQL transaction.
// Another Edit object should not be created before ending the previous one.
//
// TODO(.): Is nesting edits ever valid? Perhaps a nil edit should be returned if the previous wasn't ended?
// What if it will never be ended? Leaving the buffer in a broken state where no more changes can be made to
// it is obviously not good and is the reason why ST3 removed the ability to manually create Edit objects
// to stop people from breaking the undo stack.
func (v *View) BeginEdit() *Edit {
	e := newEdit(v)
	v.editstack = append(v.editstack, e)
	return e
}

// Ends the given Edit object.
func (v *View) EndEdit(edit *Edit) {
	if edit.invalid {
		// This happens when nesting edits and the "child" edit ends after the parent edit.
		log4go.Fine("This edit has already been invalidated: %v, %v", edit, v.editstack)
		return
	}

	// Find the position of this edit object in this view's edit stack.
	// If plugins, commands, etc are well-behaved the ended edit should be
	// last in the stack, but "shit happens" and we cannot count on this being the case.
	i := len(v.editstack) - 1
	for i := len(v.editstack) - 1; i >= 0; i-- {
		if v.editstack[i] == edit {
			break
		}
	}
	if i == -1 {
		// TODO(.): Under what instances does this happen again?
		log4go.Error("This edit isn't even in the stack... where did it come from? %v, %v", edit, v.editstack)
		return
	}

	var selection_modified bool

	if l := len(v.editstack) - 1; i != l {
		// TODO(.): See TODO in BeginEdit
		log4go.Error("This edit wasn't last in the stack... %d !=  %d: %v, %v", i, l, edit, v.editstack)
	}

	// Invalidate all edits "below" and including this edit.
	for j := len(v.editstack) - 1; j >= i; j-- {
		current_edit := v.editstack[j]
		current_edit.invalid = true
		sel_same := reflect.DeepEqual(*v.Sel(), current_edit.savedSel)
		buf_same := v.buffer.ChangeCount() == current_edit.savedCount
		eq := (sel_same && buf_same && current_edit.composite.Len() == 0)
		if !eq && sel_same {
			// TODO(.): HUH? shouldn't that be !sel_same???
			selection_modified = true
		}
		if v.scratch || current_edit.bypassUndo || eq {
			continue
		}
		switch {
		case i == 0:
			// Good behaved, no nested edits!
			fallthrough
		case j != i:
			// BOO! Someone begun another edit without finishing the first one first.
			// In this instance, the parent edit ended before the child.
			// TODO(.): What would be the correct way to handle this?
			v.undoStack.Add(edit)
		default:
			// BOO! Also bad behaved. This edit object was begun after the parent begun,
			// but was at least finished before the parent finished.
			//
			// Adding it as a child of the parent edit so that undoing the parent
			// will undo this edit as well.
			v.editstack[i-1].composite.Add(current_edit)
		}
	}
	// Pop this edit and all the children from the edit stack.
	v.editstack = v.editstack[:i]
	if selection_modified {
		OnSelectionModified.Call(v)
	}
}

// Sets the scratch property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) SetScratch(s bool) {
	v.scratch = s
}

// Checks the scratch property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) IsScratch() bool {
	return v.scratch
}

// Sets the overwrite status property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) OverwriteStatus() bool {
	return v.overwrite
}

// Checks the overwrite status property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) SetOverwriteStatus(s bool) {
	v.overwrite = s
}

// Returns the CommandHistory entry at the given index.
// TODO(.): Yeah, I don't know what that means either ;)
func (v *View) CommandHistory(idx int, modifying_only bool) (name string, args Args, count int) {
	// TODO(.): merge history when possible
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

// AddRegions lets users mark text regions in a view with a scope name, gutter icon and ViewRegionflags
// which are then optionally used to alter the display of those regions.
//
// Typical uses would be to draw squiggly lines under misspelled words, show an icon in the gutter to
// indicate a breakpoint, keeping track of snippet or auto-completion fields, highlight code compilation
// warnings, etc.
//
// The regions will be automatically adjusted as appropriate when the underlying buffer is changed.
func (v *View) AddRegions(key string, regions []Region, scope, icon string, flags render.ViewRegionFlags) {
	vr := render.ViewRegions{Scope: scope, Icon: icon, Flags: flags}
	vr.Regions.AddAll(regions)

	v.lock.Lock()
	defer v.lock.Unlock()
	v.regions[key] = vr
}

// Returns the regions associated by the given key.
func (v *View) GetRegions(key string) (ret []Region) {
	v.lock.Lock()
	defer v.lock.Unlock()
	vr := v.regions[key]
	rs := vr.Regions.Regions()
	ret = make([]Region, len(rs))
	copy(ret, rs)
	return
}

// Removes the regions associated with the given key from the view.
func (v *View) EraseRegions(key string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	delete(v.regions, key)
}

// Returns the UndoStack of this view. Tread lightly.
func (v *View) UndoStack() *UndoStack {
	return &v.undoStack
}

// Transform takes a ColourScheme and a viewport and returns a Recipe suitable
// for rendering the contents of this view that is visible in that viewport.
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
	rs := render.ViewRegions{Flags: render.SELECTION}
	rs.Regions.AddAll(v.selection.Regions())
	rr["lime.selection"] = rs
	return render.Transform(scheme, rr, viewport)
}
