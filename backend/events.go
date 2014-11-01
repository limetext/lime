// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package backend

import (
	"github.com/limetext/lime/backend/log"
	"github.com/limetext/lime/backend/util"
	"strings"
)

type (
	// An event callback dealing with View events.
	ViewEventCallback func(v *View)
	// A ViewEvent is simply a bunch of ViewEventCallbacks.
	ViewEvent []ViewEventCallback

	// The return value returned from a QueryContextCallback.
	QueryContextReturn int

	// The context is queried when trying to figure out what action should be performed when
	// certain conditions are met.
	//
	// Context is just a string identifier, an optional comparison operator, an optional operand, and an optional
	// match_all boolean. The data of the context is optionally provided together with a key binding and the key's
	// action will only be considered if the context conditions are met.
	//
	// Exactly how these values are interpreted is up to the individual context handlers, which may be fully
	// customized by implementing the callback in a plugin.
	//
	// For instance pressing the key 'j' will have a different meaning when in a VI command mode emulation
	// and when in a VI insert mode emulation. A plugin would then define two key binding entries for 'j',
	// describe the key binding context to be able to discern which action is appropriate when 'j' is then pressed.
	QueryContextCallback func(v *View, key string, operator util.Op, operand interface{}, match_all bool) QueryContextReturn

	// A QueryContextEvent is simply a bunch of QueryContextCallbacks.
	QueryContextEvent []QueryContextCallback

	// A WindowEventCallback deals with Window events.
	WindowEventCallback func(w *Window)
	// A WindowEvent is simply a bunch of WindowEventCallbacks.
	WindowEvent []WindowEventCallback
)

const (
	True    QueryContextReturn = iota //< Returned when the context query matches.
	False                             //< Returned when the context query does not match.
	Unknown                           //< Returned when the QueryContextCallback does not know how to deal with the given context.
)

// Add the provided ViewEventCallback to this ViewEvent
// TODO(.): Support removing ViewEventCallbacks?
func (ve *ViewEvent) Add(cb ViewEventCallback) {
	*ve = append(*ve, cb)
}

// Trigger this ViewEvent by calling all the registered callbacks in order of registration.
func (ve *ViewEvent) Call(v *View) {
	log.Finest("%s(%v)", evNames[ve], v.Id())
	for _, ev := range *ve {
		ev(v)
	}
}

// Add the provided QueryContextCallback to the QueryContextEvent.
// TODO(.): Support removing QueryContextCallbacks?
func (qe *QueryContextEvent) Add(cb QueryContextCallback) {
	*qe = append(*qe, cb)
}

// Searches for a QueryContextCallback and returns the result of the first callback being able to deal with this
// context, or Unknown if no such callback was found.
func (qe QueryContextEvent) Call(v *View, key string, operator util.Op, operand interface{}, match_all bool) QueryContextReturn {
	log.Fine("Query context: %s, %v, %v, %v", key, operator, operand, match_all)
	for i := range qe {
		r := qe[i](v, key, operator, operand, match_all)
		if r != Unknown {
			return r
		}
	}
	log.Fine("Unknown context: %s", key)
	return Unknown
}

// Add the provided WindowEventCallback to this WindowEvent.
// TODO(.): Support removing WindowEventCallbacks?
func (we *WindowEvent) Add(cb WindowEventCallback) {
	*we = append(*we, cb)
}

// Trigger this WindowEvent by calling all the registered callbacks in order of registration.
func (we *WindowEvent) Call(w *Window) {
	log.Finest("%s(%v)", wevNames[we], w.Id())
	for _, ev := range *we {
		ev(w)
	}
}

var (
	OnNew               ViewEvent //< Called when a new view is created
	OnLoad              ViewEvent //< Called when loading a view's buffer has finished
	OnActivated         ViewEvent //< Called when a view gains input focus.
	OnDeactivated       ViewEvent //< Called when a view loses input focus.
	OnPreClose          ViewEvent //< Called when a view is about to be closed.
	OnClose             ViewEvent //< Called when a view has been closed.
	OnPreSave           ViewEvent //< Called just before a view's buffer is saved.
	OnPostSave          ViewEvent //< Called after a view's buffer has been saved.
	OnModified          ViewEvent //< Called when the contents of a view's underlying buffer has changed.
	OnSelectionModified ViewEvent //< Called when a view's Selection/cursor has changed.

	OnNewWindow    WindowEvent       //< Called when a new window has been created.
	OnQueryContext QueryContextEvent //< Called when context is being queried.
)

var (
	evNames = map[*ViewEvent]string{
		&OnNew:               "OnNew",
		&OnLoad:              "OnLoad",
		&OnActivated:         "OnActivated",
		&OnDeactivated:       "OnDeactivated",
		&OnPreClose:          "OnPreClose",
		&OnClose:             "OnClose",
		&OnPreSave:           "OnPreSave",
		&OnPostSave:          "OnPostSave",
		&OnModified:          "OnModified",
		&OnSelectionModified: "OnSelectionModified",
	}
	wevNames = map[*WindowEvent]string{
		&OnNewWindow: "OnNewWindow",
	}
)

func init() {
	// Register functionality dealing with a couple of built in contexts
	OnQueryContext.Add(func(v *View, key string, operator util.Op, operand interface{}, match_all bool) QueryContextReturn {
		if strings.HasPrefix(key, "setting.") && operator == util.OpEqual {
			c, ok := v.Settings().Get(key[8:]).(bool)
			if c && ok {
				return True
			}
			return False
		} else if key == "num_selections" {
			opf, _ := operand.(float64)
			op := int(opf)

			switch operator {
			case util.OpEqual:
				if op == v.Sel().Len() {
					return True
				}
				return False
			case util.OpNotEqual:
				if op != v.Sel().Len() {
					return True
				}
				return False
			}
		}
		return Unknown
	})

	OnLoad.Add(func(view *View) {
		editor := GetEditor()
		editor.Watch(view)
	})
}
