// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/packages"
	"github.com/limetext/lime/backend/parser"
	"github.com/limetext/lime/backend/render"
	"github.com/limetext/lime/backend/textmate"
	. "github.com/limetext/lime/backend/util"
	"github.com/limetext/rubex"
	. "github.com/limetext/text"
	"io/ioutil"
	"os"
	"path"
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
		reparseChan chan parseReq
	}
	parseReq struct {
		forced bool
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
			ret.reparse(true)
			ret.loadSettings()
		}
	})

	ret.reparseChan = make(chan parseReq, 32)
	go ret.parsethread()
	ret.Settings().Set("is_widget", false)
	return ret
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

// Flush is called every time the underlying buffer is changed.
// It calls Adjust() on all the regions associated with this view,
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

// parsethread() would be the go-routine used for dealing with reparsing the
// current buffer when it has been modified. Each opened view has its own
// go-routine parsethread() which sits idle and waits for requests to be sent
// on this view's reparseChan.
//
// The Buffer's ChangeCount, as well as the parse request's "forced" attribute
// is used to determined if a parse actually needs to happen or not.
//
// If it is decided that a reparse should take place, a snapshot of the Buffer is
// taken and a parse is performed. Upon completion of this parse operation,
// and if the snapshot of the buffer has not already become outdated,
// then the regions of the view associated with syntax highlighting is updated.
//
// Changes made to the Buffer during the time when there is no accurate
// parse of the buffer is a monkey-patched version of the old syntax highlighting
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
				log.Error("Panic in parse thread: %v\n%s", r, string(debug.Stack()))
				if pc > 0 {
					panic(r)
				}
				pc++
			}
		}()

		b := v.Buffer()
		sub := b.Substr(Region{0, b.Size()})

		source, _ := v.Settings().Get("syntax", "").(string)
		if len(source) == 0 {
			return
		}

		// TODO: Allow other parsers instead of this hardcoded textmate version
		pr, err := textmate.NewLanguageParser(source, sub)
		if err != nil {
			log.Error("Couldn't parse: %v", err)
			return
		}

		syn, err := parser.NewSyntaxHighlighter(pr)
		if err != nil {
			log.Error("Couldn't create syntaxhighlighter: %v", err)
			return
		}

		// Only set if it isn't invalid already, otherwise the
		// current syntax highlighting will be more accurate
		// as it will have had incremental adjustments done to it
		if v.buffer.ChangeCount() != lastParse {
			return
		}

		v.lock.Lock()
		defer v.lock.Unlock()

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

	v.lock.Lock()
	ch := v.reparseChan
	v.lock.Unlock()
	defer v.cleanup()
	if ch == nil {
		return
	}

	for pr := range ch {
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
// even if the Buffer appears to not have changed.
//
// The actual parsing is done in a separate go-routine, for which the
// "lime.syntax.updated" setting will be set once it has finished.
//
// Note that it's presumed that the function calling this function
// has locked the view!
func (v *View) reparse(forced bool) {
	if v.isClosed() {
		// No point in issuing a re-parse if the view has been closed
		return
	}
	if len(v.reparseChan) < cap(v.reparseChan) || forced {
		v.reparseChan <- parseReq{forced}
	}
}

// Will load view settings respect to current syntax
// e.g if current syntax is Python settings order will be:
// Packages/Python/Python.sublime-settings
// Packages/Python/Python (Windows).sublime-settings
// Packages/User/Python.sublime-settings
// <Buffer Specific Settings>
func (v *View) loadSettings() {
	syntax := v.Settings().Get("syntax", "").(string)

	if syntax == "" {
		v.Settings().SetParent(v.window)
		return
	}

	defSettings, usrSettings, platSettings := &HasSettings{}, &HasSettings{}, &HasSettings{}

	defSettings.Settings().SetParent(v.window)
	platSettings.Settings().SetParent(defSettings)
	usrSettings.Settings().SetParent(platSettings)
	v.Settings().SetParent(usrSettings)

	ed := GetEditor()
	if r, err := rubex.Compile(`([A-Za-z]+?)\.(?:[^.]+)$`); err != nil {
		log.Error(err)
		return
	} else if s := r.FindStringSubmatch(syntax); s != nil {
		p := path.Join(LIME_PACKAGES_PATH, s[1], s[1]+".sublime-settings")
		ed.loadSetting(packages.NewPacket(p, defSettings.Settings()))

		p = path.Join(LIME_PACKAGES_PATH, s[1], s[1]+" ("+ed.plat()+").sublime-settings")
		ed.loadSetting(packages.NewPacket(p, platSettings.Settings()))

		p = path.Join(LIME_USER_PACKETS_PATH, s[1]+".sublime-settings")
		ed.loadSetting(packages.NewPacket(p, usrSettings.Settings()))
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

// Returns the Region of the innermost scope that contains "point".
// See package lime/backend/parser for details.
func (v *View) ExtractScope(point int) Region {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.syntax != nil {
		return v.syntax.ScopeExtent(point)
	}
	return Region{}
}

// ScoreSelector() takes a point and a selector string and returns a score
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

// Sel() returns a pointer to the RegionSet used by this View
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

// Returns the window this View belongs to.
func (v *View) Window() *Window {
	return v.window
}

// Returns the underlying Buffer that this View is a view into.
func (v *View) Buffer() Buffer {
	return v.buffer
}

// Inserts text at the given position in the provided edit object.
// Tabs are (sometimes, depending on the View's settings) translated to spaces.
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

// Adds an Erase action of the given Region to the provided Edit object.
func (v *View) Erase(edit *Edit, r Region) {
	edit.composite.AddExec(NewEraseAction(v.buffer, r))
}

// Adds a Replace action of the given Region to the provided Edit object.
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
		// This happens when nesting Edits and the child Edit ends after the parent edit.
		log.Fine("This edit has already been invalidated: %v, %v", edit, v.editstack)
		return
	}

	// Find the position of this Edit object in this View's Edit stack.
	// If plugins, commands, etc are well-behaved the ended edit should be
	// last in the stack, but shit happens and we cannot count on this being the case.
	i := len(v.editstack) - 1
	for i := len(v.editstack) - 1; i >= 0; i-- {
		if v.editstack[i] == edit {
			break
		}
	}
	if i == -1 {
		// TODO(.): Under what instances does this happen again?
		log.Error("This edit isn't even in the stack... where did it come from? %v, %v", edit, v.editstack)
		return
	}

	var selection_modified bool

	if l := len(v.editstack) - 1; i != l {
		// TODO(.): See TODO in BeginEdit
		log.Error("This edit wasn't last in the stack... %d !=  %d: %v, %v", i, l, edit, v.editstack)
	}

	// Invalidate all Edits "below" and including this Edit.
	for j := len(v.editstack) - 1; j >= i; j-- {
		current_edit := v.editstack[j]
		current_edit.invalid = true
		sel_same := reflect.DeepEqual(*v.Sel(), current_edit.savedSel)
		buf_same := v.buffer.ChangeCount() == current_edit.savedCount
		eq := (sel_same && buf_same && current_edit.composite.Len() == 0)
		if !eq && !sel_same {
			selection_modified = true
		}
		if v.IsScratch() || current_edit.bypassUndo || eq {
			continue
		}
		switch {
		case i == 0:
			// Well-behaved, no nested edits!
			fallthrough
		case j != i:
			// BOO! Someone began another Edit without finishing the first one.
			// In this instance, the parent Edit ended before the child.
			// TODO(.): What would be the correct way to handle this?
			v.undoStack.Add(edit)
		default:
			// BOO! Also poorly-behaved. This Edit object began after the parent began,
			// but was finished before the parent finished.
			//
			// Add it as a child of the parent Edit so that undoing the parent
			// will undo this edit as well.
			v.editstack[i-1].composite.Add(current_edit)
		}
	}
	// Pop this Edit and all the children off the Edit stack.
	v.editstack = v.editstack[:i]
	if selection_modified {
		OnSelectionModified.Call(v)
	}
}

// Sets the scratch property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) SetScratch(s bool) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.scratch = s
}

// Checks the scratch property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) IsScratch() bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.scratch
}

// Sets the overwrite status property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) OverwriteStatus() bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.overwrite
}

// Checks the overwrite status property of the view.
// TODO(.): Couldn't this just be a value in the View's Settings?
func (v *View) SetOverwriteStatus(s bool) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.overwrite = s
}

// Returns whether the underlying Buffer has any unsaved modifications.
// Note that Scratch buffers are never considered dirty.
func (v *View) IsDirty() bool {
	if v.IsScratch() {
		return false
	}
	lastSave, _ := v.buffer.Settings().Get("lime.last_save_change_count", -1).(int)
	return v.buffer.ChangeCount() != lastSave
}

func (v *View) FileChanged(filename string) {
	log.Finest("Reloading %s", filename)

	if saving, ok := v.Settings().Get("lime.saving", false).(bool); ok && saving {
		// This reload was triggered by ourselves saving to this file, so don't reload it
		return
	}
	if !GetEditor().Frontend().OkCancelDialog("File was changed by another program, reload?", "reload") {
		return
	}

	if d, err := ioutil.ReadFile(filename); err != nil {
		log.Error("Could not read file: %s\n. Error was: %v", filename, err)
	} else {
		edit := v.BeginEdit()
		end := v.Buffer().Size()
		v.Replace(edit, Region{0, end}, string(d))
		v.EndEdit(edit)
	}
}

// Saves the file
func (v *View) Save() error {
	return v.SaveAs(v.buffer.FileName())
}

// Saves the file to the specified filename
func (v *View) SaveAs(name string) (err error) {
	log.Fine("SaveAs(%s)", name)
	v.Settings().Set("lime.saving", true)
	defer v.Settings().Erase("lime.saving")
	var atomic bool
	OnPreSave.Call(v)
	if atomic, _ = v.Settings().Get("atomic_save", true).(bool); v.buffer.FileName() == "" || !atomic {
		if err := v.nonAtomicSave(name); err != nil {
			return err
		}
	} else {
		n, err := ioutil.TempDir(path.Dir(v.buffer.FileName()), "lime")
		if err != nil {
			return err
		}
		tmpf := path.Join(n, "tmp")
		if err := v.nonAtomicSave(tmpf); err != nil {
			return err
		}
		if err := os.Rename(tmpf, name); err != nil {
			// When we want to save as a file in another directory
			// we can't go with os.Rename so we need to force
			// not atomic saving sometimes as 4th test in TestSaveAsOpenFile
			if err := v.nonAtomicSave(name); err != nil {
				return err
			}
		}
		if err := os.RemoveAll(n); err != nil {
			return err
		}
	}

	ed := GetEditor()
	if fn := v.buffer.FileName(); fn != name {
		v.Buffer().SetFileName(name)
		if fn != "" {
			ed.UnWatch(fn, v)
		}
		ed.Watch(name, v)
	}

	v.buffer.Settings().Set("lime.last_save_change_count", v.buffer.ChangeCount())
	OnPostSave.Call(v)
	return nil
}

func (v *View) nonAtomicSave(name string) error {
	data := []byte(v.buffer.Substr(Region{0, v.buffer.Size()}))
	if err := ioutil.WriteFile(name, data, 0644); err != nil {
		return err
	}
	return nil
}

// Returns the CommandHistory entry at the given relative index.
//
// When "modifying_only" is set to true, only commands that actually changed
// the Buffer in some way (as opposed to just moving the cursor around) are counted as
// an index. That would be a "hard" command as it is referred to in UndoStack.Undo.
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
			log.Error("Paniced while running text command %s %v: %v\n%s", name, cmd, r, string(debug.Stack()))
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

// Returns the Regions associated with the given key.
func (v *View) GetRegions(key string) (ret []Region) {
	v.lock.Lock()
	defer v.lock.Unlock()
	vr := v.regions[key]
	rs := vr.Regions.Regions()
	ret = make([]Region, len(rs))
	copy(ret, rs)
	return
}

// Removes the Regions associated with the given key from the view.
func (v *View) EraseRegions(key string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	delete(v.regions, key)
}

// Returns the UndoStack of this view. Tread lightly.
func (v *View) UndoStack() *UndoStack {
	return &v.undoStack
}

// Transform() takes a ColourScheme and a viewport and returns a Recipe suitable
// for rendering the contents of this View that is visible in that viewport.
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

func (v *View) cleanup() {
	v.lock.Lock()
	defer v.lock.Unlock()

	// TODO(.): There can be multiple views into a single Buffer,
	// need to do some reference counting to see when it should be
	// closed
	v.buffer.Close()
	v.buffer = nil
}

func (v *View) isClosed() bool {
	return v.reparseChan == nil
}

// Initiate the "close" operation of this view.
// Returns "true" if the view was closed. Otherwise returns "false".
func (v *View) Close() bool {
	OnPreClose.Call(v)
	if v.IsDirty() {
		close_anyway := GetEditor().Frontend().OkCancelDialog("File has been modified since last save, close anyway?", "Close")
		if !close_anyway {
			return false
		}
	}
	if n := v.buffer.FileName(); n != "" {
		GetEditor().UnWatch(n, v)
	}

	// Call the event first while there's still access possible to the underlying
	// buffer
	OnClose.Call(v)

	v.window.remove(v)

	// Closing the reparseChan, and setting to nil will eventually clean up other resources
	// when the parseThread exits
	v.lock.Lock()
	defer v.lock.Unlock()
	close(v.reparseChan)
	v.reparseChan = nil

	return true
}

const (
	CLASS_WORD_START = 1 << iota
	CLASS_WORD_END
	CLASS_PUNCTUATION_START
	CLASS_PUNCTUATION_END
	CLASS_SUB_WORD_START
	CLASS_SUB_WORD_END
	CLASS_LINE_START
	CLASS_LINE_END
	CLASS_EMPTY_LINE
	CLASS_MIDDLE_WORD
	CLASS_WORD_START_WITH_PUNCTUATION
	CLASS_WORD_END_WITH_PUNCTUATION
	CLASS_OPENING_PARENTHESIS
	CLASS_CLOSING_PARENTHESIS
)

// Classifies point, returning a bitwise OR of zero or more of defined flags
func (v *View) Classify(point int) (res int) {
	var a, b string = "", ""
	ws := v.Settings().Get("word_separators", "[!\"#$%&'()*+,\\-./:;<=>?@\\[\\\\\\]^_`{|}~]").(string)
	if point > 0 {
		a = v.buffer.Substr(Region{point - 1, point})
	}
	if point < v.buffer.Size() {
		b = v.buffer.Substr(Region{point, point + 1})
	}

	// Special cases
	if v.buffer.Size() == 0 || point < 0 || point > v.buffer.Size() {
		res = 3520
		return
	}
	if re, err := rubex.Compile("[A-Z]"); err != nil {
		log.Error(err)
	} else {
		if re.MatchString(b) && !re.MatchString(a) {
			res |= CLASS_SUB_WORD_START
			res |= CLASS_SUB_WORD_END
		}
	}
	if a == "," {
		res |= CLASS_OPENING_PARENTHESIS
	}
	if b == "," {
		res |= CLASS_CLOSING_PARENTHESIS
	}
	if a == "," && b == "," {
		res = 0
		return
	}
	// Punc start & end
	if re, err := rubex.Compile(ws); err != nil {
		log.Error(err)
	} else {
		if (re.MatchString(b) || b == "") && !re.MatchString(a) {
			res |= CLASS_PUNCTUATION_START
		}
		if (re.MatchString(a) || a == "") && !re.MatchString(b) {
			res |= CLASS_PUNCTUATION_END
		}
		// Word start & end
		if re1, err := rubex.Compile("\\w"); err != nil {
			log.Error(err)
		} else if re2, err := rubex.Compile("\\s"); err != nil {
			log.Error(err)
		} else {
			if re1.MatchString(b) && (re.MatchString(a) || re2.MatchString(a) || a == "") {
				res |= CLASS_WORD_START
			}
			if re1.MatchString(a) && (re.MatchString(b) || re2.MatchString(b) || b == "") {
				res |= CLASS_WORD_END
			}
		}
	}
	// SubWord start & end

	// Line start & end
	if a == "\n" || a == "" {
		res |= CLASS_LINE_START
	}
	if b == "\n" || b == "" {
		res |= CLASS_LINE_END
	}
	// Empty line
	if (a == "\n" && b == "\n") || (a == "" && b == "") {
		res |= CLASS_EMPTY_LINE
	}
	// Middle word
	if re, err := rubex.Compile("\\w"); err != nil {
		log.Error(err)
	} else {
		if re.MatchString(a) && re.MatchString(b) {
			res |= CLASS_MIDDLE_WORD
		}
	}
	// Word start & end with punc
	if re, err := rubex.Compile("\\s"); err != nil {
		log.Error(err)
	} else {
		if (res&CLASS_PUNCTUATION_START == CLASS_PUNCTUATION_START) && (re.MatchString(a) || a == "") {
			res |= CLASS_WORD_START_WITH_PUNCTUATION
		}
		if (res&CLASS_PUNCTUATION_END == CLASS_PUNCTUATION_END) && (re.MatchString(b) || b == "") {
			res |= CLASS_WORD_END_WITH_PUNCTUATION
		}
	}
	// Openning & closing parentheses
	if re, err := rubex.Compile("[(\\[{]"); err != nil {
		log.Error(err)
	} else {
		if re.MatchString(a) || re.MatchString(b) {
			res |= CLASS_OPENING_PARENTHESIS
		}
		if re.MatchString(a) && a == b {
			res = 0
			return
		}
	}
	if re, err := rubex.Compile("[)\\]}]"); err != nil {
		log.Error(err)
	} else {
		if re.MatchString(a) || re.MatchString(b) {
			res |= CLASS_CLOSING_PARENTHESIS
		}
		if re.MatchString(a) && a == b {
			res = 0
			return
		}
	}
	return
}

// Finds the next location after point that matches the given classes
// Searches backward if forward is false
func (v *View) FindByClass(point int, forward bool, classes int) int {
	i := -1
	if forward {
		i = 1
	}
	size := v.buffer.Size()
	// Sublime doesn't consider initial point even if it matches.
	for p := point + i; ; p += i {
		if p <= 0 {
			return 0
		}
		if p >= size {
			return size
		}
		if v.Classify(p)&classes != 0 {
			return p
		}
	}
}

// Expands the selection until the point on each side matches the given classes
func (v *View) ExpandByClass(r Region, classes int) Region {
	// Sublime doesn't consider the points the region starts on.
	// If not already on edge of buffer, expand by 1 in both directions.
	a := r.A
	if a > 0 {
		a -= 1
	} else if a < 0 {
		a = 0
	}

	b := r.B
	size := v.buffer.Size()
	if b < size {
		b += 1
	} else if b > size {
		b = size
	}

	for ; a > 0 && (v.Classify(a)&classes == 0); a -= 1 {
	}
	for ; b < size && (v.Classify(b)&classes == 0); b += 1 {
	}
	return Region{a, b}
}
