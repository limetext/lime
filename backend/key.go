package backend

import (
	"bytes"
	"code.google.com/p/log4go"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
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
	OpEqual Op = iota
	OpNotEqual
	OpRegexMatch
	OpNotRegexMatch
	OpRegexContains
	OpNotRegexContains
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

type (
	Key      rune
	KeyPress struct {
		Key                     Key
		Shift, Super, Alt, Ctrl bool
	}
	Op         int
	KeyContext struct {
		rawKeyContext
	}
	rawKeyContext struct {
		Key      string
		Operator Op
		Operand  interface{}
		MatchAll bool `json:"match_all"`
	}

	KeyBinding struct {
		Keys     []KeyPress
		Command  string
		Args     Args
		Context  []KeyContext
		priority int
	}

	KeyBindings struct {
		Bindings []*KeyBinding
		keyOff   int
	}
)

func (k *Op) UnmarshalJSON(d []byte) error {
	var tmp string
	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}
	switch tmp {
	default:
		*k = OpEqual
	case "not_equal":
		*k = OpNotEqual
	case "regex_match":
		*k = OpRegexMatch
	case "not_regex_match":
		*k = OpNotRegexMatch
	case "regex_contains":
		*k = OpRegexContains
	case "not_regex_contains":
		*k = OpNotRegexContains
	}

	return nil
}

func (k *KeyContext) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, &k.rawKeyContext); err != nil {
		return err
	}
	if k.Operand == nil {
		k.Operand = true
	}
	return nil
}

func (k KeyPress) Index() (ret int) {
	ret = int(k.Key)
	if k.Shift {
		ret += shift
	}
	if k.Alt {
		ret += alt
	}
	if k.Ctrl {
		ret += ctrl
	}
	if k.Super {
		ret += super
	}
	return
}

func (k KeyPress) IsCharacter() bool {
	return unicode.IsPrint(rune(k.Key)) && !k.Super && !k.Ctrl
}

func (k *KeyPress) fix() {
	lower := Key(unicode.ToLower(rune(k.Key)))
	if lower != k.Key {
		k.Shift = true
		k.Key = lower
	}
}

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
				r := []Key(c)
				if len(r) != 1 {
					log4go.Warn("Unknown key value with %d bytes: %s", len(c), c)
					return nil
				}
				k.Key = Key(c[0])
				k.fix()
			}
		}
	}
	return nil
}

func (k *KeyBindings) Len() int {
	return len(k.Bindings)
}

func (k *KeyBindings) Less(i, j int) bool {
	return k.Bindings[i].Keys[k.keyOff].Index() < k.Bindings[j].Keys[k.keyOff].Index()
}

func (k *KeyBindings) Swap(i, j int) {
	k.Bindings[i], k.Bindings[j] = k.Bindings[j], k.Bindings[i]
}

func (k *KeyBindings) DropLessEqualKeys(count int) {
	for i := 0; i < len(k.Bindings); {
		if len(k.Bindings[i].Keys) <= count {
			k.Bindings[i] = k.Bindings[len(k.Bindings)-1]
			k.Bindings = k.Bindings[:len(k.Bindings)-1]
		} else {
			i++
		}
	}
	sort.Sort(k)
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

func (k *KeyBindings) merge(other *KeyBindings) {
	// TODO: what is the order really? Newer on top?
	k.Bindings = append(k.Bindings, other.Bindings...)
	k.DropLessEqualKeys(0)
}

func (k *KeyBindings) filter(ki int, ret *KeyBindings) {
	idx := sort.Search(k.Len(), func(i int) bool {
		return k.Bindings[i].Keys[k.keyOff].Index() >= ki
	})
	for i := idx; i < len(k.Bindings) && k.Bindings[i].Keys[k.keyOff].Index() == ki; i++ {
		ret.Bindings = append(ret.Bindings, k.Bindings[i])
	}
}

func (k *KeyBindings) Filter(kp KeyPress) (ret KeyBindings) {
	p := Prof.Enter("key.filter")
	defer p.Exit()

	kp.fix()
	k.DropLessEqualKeys(k.keyOff)
	ret.keyOff = k.keyOff + 1
	ki := kp.Index()

	k.filter(ki, &ret)

	if kp.IsCharacter() {
		k.filter(int(Any), &ret)
	}
	return
}

func (k *KeyBindings) Action(v *View) (kb *KeyBinding) {
	p := Prof.Enter("key.action")
	defer p.Exit()

	for i := range k.Bindings {
		if len(k.Bindings[i].Keys) > k.keyOff {
			continue
		}
		for _, c := range k.Bindings[i].Context {
			if OnQueryContext.Call(v, c.Key, c.Operator, c.Operand, c.MatchAll) != True {
				goto skip
			}
		}
		if kb == nil || kb.priority < k.Bindings[i].priority {
			kb = k.Bindings[i]
		}
	skip:
	}
	return
}

func (k Key) String() string {
	if v, ok := rkeylut[k]; ok {
		return v
	}
	return string(k)
}
func (k KeyPress) String() (ret string) {
	if k.Super {
		ret += "super+"
	}
	if k.Ctrl {
		ret += "ctrl+"
	}
	if k.Alt {
		ret += "alt+"
	}
	if k.Shift {
		ret += "shift+"
	}
	ret += fmt.Sprintf("%s", k.Key)
	return
}

func (k KeyBindings) String() string {
	var buf bytes.Buffer
	for _, b := range k.Bindings {
		buf.WriteString(fmt.Sprintf("%+v\n", b))
	}
	return buf.String()
}
