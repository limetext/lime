// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"fmt"
	"github.com/limetext/lime/backend/log"
	"strings"
	"unicode"
)

type (
	// A Key press with the given Key
	// and modifiers.
	KeyPress struct {
		Key                     Key
		Shift, Super, Alt, Ctrl bool
	}
)

// Returns an index used for sorting key presses.
// TODO(.): This is in no way a unique index with quite a lot of collisions and potentially resulting
// in bad lookups.
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

// Returns whether this KeyPress is a print character or not.
func (k KeyPress) IsCharacter() bool {
	return unicode.IsPrint(rune(k.Key)) && !k.Super && !k.Ctrl
}

// Modifies the KeyPress so that it's Key is a unicode lower case
// rune and if it was in uppercase before this modification, the
// "Shift" modifier is also enabled.
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
					log.Global.LogWarning("Unknown key value with %d bytes: %s", len(c), c)
					return nil
				}
				k.Key = Key(c[0])
				k.fix()
			}
		}
	}
	return nil
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
