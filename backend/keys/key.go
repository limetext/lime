// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"unicode"
)

type (
	Key rune
)

const (
	Left Key = 0x2190 + iota
	Up
	Right
	Down

	Enter       = '\n'
	Escape      = 0x001B
	Backspace   = 0x0008
	Delete      = 0x007F
	KeypadEnter = '\n'
)

const (
	// map to dingbats...
	F1 Key = 0x2701 + iota
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	Insert
	PageUp
	PageDown
	Home
	End
	Break
	Any Key = unicode.MaxRune
)

const (
	shift = (1 << (29 - iota))
	ctrl
	alt
	super
)

var keylut = map[string]Key{
	"up":            Up,
	"left":          Left,
	"right":         Right,
	"down":          Down,
	"enter":         Enter,
	"tab":           '\t',
	"escape":        Escape,
	"space":         ' ',
	"f1":            F1,
	"f2":            F2,
	"f3":            F3,
	"f4":            F4,
	"f5":            F5,
	"f6":            F6,
	"f7":            F7,
	"f8":            F8,
	"f9":            F9,
	"f10":           F10,
	"f11":           F11,
	"f12":           F12,
	"backspace":     Backspace,
	"delete":        Delete,
	"keypad_enter":  KeypadEnter,
	"insert":        Insert,
	"pageup":        PageUp,
	"pagedown":      PageDown,
	"home":          Home,
	"end":           End,
	"break":         Break,
	"forward_slash": '/',
	"backquote":     '`',
	"\\\"":          '"',
	"plus":          '+',
	"minus":         '-',
	"equals":        '=',
	"<character>":   Any,
}

var rkeylut = map[Key]string{
	Up:        "up",
	Left:      "left",
	Right:     "right",
	Down:      "down",
	Enter:     "enter",
	'\t':      "tab",
	Escape:    "escape",
	' ':       "space",
	F1:        "f1",
	F2:        "f2",
	F3:        "f3",
	F4:        "f4",
	F5:        "f5",
	F6:        "f6",
	F7:        "f7",
	F8:        "f8",
	F9:        "f9",
	F10:       "f10",
	F11:       "f11",
	F12:       "f12",
	Backspace: "backspace",
	Delete:    "delete",
	Insert:    "insert",
	PageUp:    "pageup",
	PageDown:  "pagedown",
	Home:      "home",
	End:       "end",
	Break:     "break",
	'/':       "forward_slash",
	'`':       "backquote",
	'"':       "\\\"",
	'+':       "plus",
	'-':       "minus",
	'=':       "equals",
	Any:       "<character>",
}

func (k Key) String() string {
	if v, ok := rkeylut[k]; ok {
		return v
	}
	return string(k)
}
