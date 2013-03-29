package backend

import (
	"code.google.com/p/log4go"
	"strings"
)

const (
	Left rune = 0x2190 + iota
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
	F1 rune = 0x2701 + iota
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
)

var keylut = map[string]rune{
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
}

type (
	KeyPress struct {
		Key                     rune
		Shift, Super, Alt, Ctrl bool
	}

	KeyContext struct {
		Key, Operator string
		Operand       interface{}
		MatchAll      bool `json:"match_all"`
	}

	KeyBinding struct {
		Keys    []KeyPress
		Command string
		Args    map[string]interface{}
		Context []KeyContext
	}
)

func (k *KeyPress) UnmarshalJSON(d []byte) error {
	combo := strings.Split(string(d[1:len(d)-1]), "+")
	for _, c := range combo {
		lower := strings.ToLower(c)
		switch lower {
		case "super":
			k.Super = true
		case "ctrl":
			k.Ctrl = true
		case "alt":
			k.Alt = true
		case "shift":
			k.Shift = true
		default:
			if v, ok := keylut[lower]; ok {
				k.Key = v
			} else {
				r := []rune(lower)
				if len(r) != 1 {
					log4go.Warn("Unknown key value with %d bytes: %s", len(c), c)
					return nil
				}
				if lower != c {
					k.Shift = true
				}
				k.Key = r[0]
			}
		}
	}
	return nil
}
