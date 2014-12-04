// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/limetext/lime/backend/util"
	"sort"
)

type (
	// A single KeyBinding for which after pressing the given
	// sequence of Keys, and the Context matches,
	// the Command will be invoked with the provided Args.
	KeyBinding struct {
		Keys     []KeyPress
		Command  string
		Args     map[string]interface{}
		Context  []KeyContext
		priority int
	}

	KeyBindings struct {
		Bindings []*KeyBinding
		seqIndex int // The index we are in a multiple key sequence keybinding
		parent   *KeyBindings
	}
)

// Returns the number of KeyBindings.
func (k *KeyBindings) Len() int {
	return len(k.Bindings)
}

// Compares one KeyBinding to another for sorting purposes.
func (k *KeyBindings) Less(i, j int) bool {
	return k.Bindings[i].Keys[k.seqIndex].Index() < k.Bindings[j].Keys[k.seqIndex].Index()
}

// Swaps the two KeyBindings at the given positions.
func (k *KeyBindings) Swap(i, j int) {
	k.Bindings[i], k.Bindings[j] = k.Bindings[j], k.Bindings[i]
}

// Drops all KeyBindings that are a sequence of key presses less or equal
// to the given number.
func (k *KeyBindings) DropLessEqualKeys(count int) {
	for {
		for i := 0; i < len(k.Bindings); {
			if len(k.Bindings[i].Keys) <= count {
				k.Bindings[i] = k.Bindings[len(k.Bindings)-1]
				k.Bindings = k.Bindings[:len(k.Bindings)-1]
			} else {
				i++
			}
		}
		sort.Sort(k)
		if k.parent == nil {
			break
		}
		k = k.parent
	}
}

func (k *KeyBindings) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, &k.Bindings); err != nil {
		return err
	}
	for i := range k.Bindings {
		k.Bindings[i].priority = i
	}
	k.DropLessEqualKeys(0)
	return nil
}

func (k *KeyBindings) SetParent(p *KeyBindings) {
	k.parent = p
	// All parents and childs seqIndex must be equal
	p.seqIndex = k.seqIndex
}

func (k *KeyBindings) Parent() *KeyBindings {
	return k.parent
}

func (k *KeyBindings) filter(ki int, ret *KeyBindings) {
	for {
		idx := sort.Search(k.Len(), func(i int) bool {
			return k.Bindings[i].Keys[k.seqIndex].Index() >= ki
		})
		for i := idx; i < len(k.Bindings) && k.Bindings[i].Keys[k.seqIndex].Index() == ki; i++ {
			ret.Bindings = append(ret.Bindings, k.Bindings[i])
		}
		if k.parent == nil {
			break
		}
		k = k.parent
		if ret.parent == nil {
			ret.SetParent(new(KeyBindings))
		}
		ret = ret.parent
	}
}

// Filters the KeyBindings, returning a new KeyBindings object containing
// a subset of matches for the given key press.
func (k *KeyBindings) Filter(kp KeyPress) (ret KeyBindings) {
	p := Prof.Enter("key.filter")
	defer p.Exit()

	kp.fix()
	k.DropLessEqualKeys(k.seqIndex)
	ret.seqIndex = k.seqIndex + 1
	ki := kp.Index()

	k.filter(ki, &ret)

	if kp.IsCharacter() {
		k.filter(int(Any), &ret)
	}
	return
}

// Tries to resolve all the current KeyBindings in k to a single
// action. If any action is appropriate as determined by context,
// the return value will be the specific KeyBinding that is possible
// to execute now, otherwise it is nil.
func (k *KeyBindings) Action(qc func(key string, operator Op, operand interface{}, match_all bool) bool) (kb *KeyBinding) {
	p := Prof.Enter("key.action")
	defer p.Exit()

	for {
		for i := range k.Bindings {
			if len(k.Bindings[i].Keys) > k.seqIndex {
				// This key binding is of a key sequence longer than what is currently
				// probed for. For example, the binding is for the sequence ['a','b','c'], but
				// the user has only pressed ['a','b'] so far.
				continue
			}
			for _, c := range k.Bindings[i].Context {
				if !qc(c.Key, c.Operator, c.Operand, c.MatchAll) {
					goto skip
				}
			}
			if kb == nil || kb.priority < k.Bindings[i].priority {
				kb = k.Bindings[i]
			}
		skip:
		}
		if kb != nil || k.parent == nil {
			break
		}
		k = k.parent
	}
	return
}

func (k *KeyBindings) SeqIndex() int {
	return k.seqIndex
}

func (k KeyBindings) String() string {
	var buf bytes.Buffer
	for _, b := range k.Bindings {
		buf.WriteString(fmt.Sprintf("%+v\n", b))
	}
	return buf.String()
}
