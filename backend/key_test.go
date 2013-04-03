package backend

import (
	"io/ioutil"
	"lime/backend/loaders"
	"testing"
)

func TestLoadKeyBindings(t *testing.T) {
	tests := []string{
		"loaders/json/testdata/Default (OSX).sublime-keymap",
		"/Users/quarnster/Library/Application Support/Sublime Text 3/Packages/Default/Default (Linux).sublime-keymap",
		"/Users/quarnster/Library/Application Support/Sublime Text 3/Packages/Default/Default (OSX).sublime-keymap",
		"/Users/quarnster/Library/Application Support/Sublime Text 3/Packages/Default/Default (Windows).sublime-keymap",
	}
	for _, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			var bindings KeyBindings
			if err := loaders.LoadJSON(d, &bindings); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestKeyFilter(t *testing.T) {
	fn := "loaders/json/testdata/Default (OSX).sublime-keymap"
	if d, err := ioutil.ReadFile(fn); err != nil {
		t.Errorf("Couldn't load file %s: %s", fn, err)
	} else {
		var bindings KeyBindings
		if err := loaders.LoadJSON(d, &bindings); err != nil {
			t.Error(err)
		}

		if b2 := bindings.Filter(KeyPress{Key: 'j', Ctrl: true}); b2.Len() != 3 {
			t.Errorf("Not of the expected length: %d, %s", 3, b2)
		} else if b3 := b2.Filter(KeyPress{Key: 's'}); b3.Len() != 1 {
			t.Errorf("Not of the expected length: %d, %s", 1, b3)
		}
	}
}
