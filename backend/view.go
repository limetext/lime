package backend

import (
	"bytes"
	"code.google.com/p/log4go"
	"fmt"
	"github.com/quarnster/parser"
	"io/ioutil"
	"lime/backend/loaders"
	. "lime/backend/primitives"
	"lime/backend/render"
	"lime/backend/textmate"
	. "lime/backend/util"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
)

type (
	View struct {
		HasSettings
		HasId
		name          string
		window        *Window
		buffer        Buffer
		selection     RegionSet
		undoStack     UndoStack
		scratch       bool
		overwrite     bool
		syntax        textmate.LanguageParser
		rootNode      *parser.Node
		lastScopeNode *parser.Node
		lastScopeBuf  bytes.Buffer
		lastScopeName string
		regions       render.ViewRegionMap
		editstack     []*Edit
		lock          sync.Mutex
		modified      bool
		reparseChan   chan parseReq
	}
	parseReq struct {
		syntax textmate.LanguageParser
		forced bool
	}
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
	v.syntax.Language = &textmate.Language{}
	b.AddCallback(func(_ Buffer, a, b int) {
		v.flush(a, b)
	})
	return nil
}

func (v *View) flush(a, b int) {
	func() {
		v.lock.Lock()
		defer v.lock.Unlock()

		e := Prof.Enter("view.flush")
		defer e.Exit()
		v.selection.Adjust(a, b)
		for k, v2 := range v.regions {
			v2.Regions.Adjust(a, b)
			v.regions[k] = v2
		}
	}()
	v.reparse(false)
}

func (v *View) parsethread() {
	pc := 0
	doparse := func(pr parseReq) {
		defer func() {
			if r := recover(); r != nil {
				log4go.Error("Panic in parse thread: %v\n%s", r, string(debug.Stack()))
				if pc > 0 {
					panic(r)
				}
				pc++
			}
		}()
		pr.syntax.Parse(v.buffer.Substr(Region{0, v.Buffer().Size()}))
		v.lock.Lock()
		v.rootNode = pr.syntax.RootNode()
		v.lastScopeNode = nil
		v.lastScopeBuf.Reset()
		v.lock.Unlock()
	}
	lastParse := -1
	for pr := range v.reparseChan {
		if cc := v.buffer.ChangeCount(); lastParse != cc || pr.forced {
			lastParse = cc
			doparse(pr)
		}
	}
}
func (v *View) reparse(forced bool) {
	v.lock.Lock()
	defer v.lock.Unlock()
	if len(v.reparseChan) < cap(v.reparseChan) || forced {
		v.reparseChan <- parseReq{v.syntax, forced}
	}
}

func (v *View) SetSyntaxFile(f string) error {
	var lang textmate.Language
	if d, err := ioutil.ReadFile(f); err != nil {
		return err
	} else if err := loaders.LoadPlist(d, &lang); err != nil {
		return err
	} else {
		v.lock.Lock()
		v.syntax.Language = &lang
		v.lock.Unlock()
		v.reparse(true)
	}
	return nil
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

	var bufmod, selmod bool

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
		if !eq && ib {
			bufmod = true
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
	if len(v.editstack) == 0 {
		if bufmod {
			OnModified.Call(v)
		}
		if selmod {
			OnSelectionModified.Call(v)
		}
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

func (v *View) findScope(search parser.Range, node *parser.Node) *parser.Node {
	idx := sort.Search(len(node.Children), func(i int) bool {
		return node.Children[i].Range.Start >= search.Start || node.Children[i].Range.Contains(search)
	})
	for idx < len(node.Children) {
		c := node.Children[idx]
		if c.Range.Start > search.End {
			break
		}
		if c.Range.Contains(search) {
			if node.Name != " " && node != v.lastScopeNode {
				v.lastScopeBuf.WriteByte(' ')
				v.lastScopeBuf.WriteString(node.Name)
			}
			return v.findScope(search, node.Children[idx])
		}
		idx++
	}
	if node != v.lastScopeNode && node.Range.Contains(search) && node.Name != "" {
		v.lastScopeBuf.WriteByte(' ')
		v.lastScopeBuf.WriteString(node.Name)
		return node
	}
	return nil
}

func (v *View) updateScope(point int) {
	if v.rootNode == nil {
		return
	}

	search := parser.Range{point, point + 1}
	if v.lastScopeNode != nil && v.lastScopeNode.Range.Contains(search) {
		if len(v.lastScopeNode.Children) != 0 {
			if no := v.findScope(search, v.lastScopeNode); no != v.lastScopeNode {
				v.lastScopeNode = no
				v.lastScopeName = v.lastScopeBuf.String()
			}
		}
	} else {
		v.lastScopeNode = nil
		v.lastScopeBuf.Reset()
		v.lastScopeNode = v.findScope(search, v.rootNode)
		v.lastScopeName = v.lastScopeBuf.String()
	}
}

func (v *View) ExtractScope(point int) Region {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.updateScope(point)
	if v.lastScopeNode != nil {
		r := v.lastScopeNode.Range
		return Region{r.Start, r.End}
	}
	return Region{}
}

func (v *View) ScopeName(point int) string {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.updateScope(point)
	return v.lastScopeName
}

func (v *View) CommandHistory(idx int, modifying_only bool) (name string, args Args, count int) {
	// TODO: merge history when possible
	if i := v.undoStack.index(idx, modifying_only); i != -1 {
		e := v.undoStack.actions[i]
		return e.command, e.args, 1
	}
	return "", nil, 0
}

func (v *View) runCommand(cmd TextCommand, name string, args Args) error {
	e := v.BeginEdit()
	e.command = name
	e.args = args
	e.bypassUndo = cmd.BypassUndo()

	defer func() {
		v.EndEdit(e)
		if r := recover(); r != nil {
			log4go.Error("Paniced while running text command %s %v: %v\n%s", name, args, r, string(debug.Stack()))
		}
	}()
	p := Prof.Enter("view.cmd." + name)
	defer p.Exit()
	return cmd.Run(v, e, args)
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
