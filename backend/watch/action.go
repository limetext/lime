// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package watch

// File events
const (
	CREATE = 1 << iota
	MODIFY
	DELETE
	RENAME

	ALL = CREATE | MODIFY | DELETE | RENAME
)

type (
	action struct {
		fn func() // The function that will be applyed
		ev int    // The events that we should apply the function
	}

	// Array of actions with a key as an identifier
	actions map[string]action
)

func (ac action) apply(ev int) {
	if ac.ev&ev != 0 {
		if ac.fn != nil {
			ac.fn()
		}
	}
}

func (acs actions) apply(ev int, keys ...string) {
	for _, key := range keys {
		if ac, exist := acs[key]; exist {
			ac.apply(ev)
		}
	}
}

func (acs actions) applyAll(ev int) {
	for _, ac := range acs {
		ac.apply(ev)
	}
}

func newEvent(evs []int) int {
	if len(evs) == 0 {
		return ALL
	}
	event := 0
	for _, ev := range evs {
		event |= ev
	}
	return event
}
