// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"fmt"
	. "github.com/limetext/text"
)

type (
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

// Returns a string describing this Edit object. Should typically not be manually called.
func (e *Edit) String() string {
	return fmt.Sprintf("%s: %v, %v, %v", e.command, e.args, e.bypassUndo, e.composite)
}

// Applies the action of this Edit object. Should typically not be manually called.
func (e *Edit) Apply() {
	e.composite.Apply()
}

// Reverses the application of this Edit object. Should typically not be manually called.
func (e *Edit) Undo() {
	e.composite.Undo()
	e.v.Sel().Clear()
	for _, r := range e.savedSel.Regions() {
		e.v.Sel().Add(r)
	}
}
