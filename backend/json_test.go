package backend

import (
	"io/ioutil"
	"testing"
)

func TestJson(t *testing.T) {
	tests := []string{
		"json/testdata/Default (OSX).sublime-keymap",
		"/Users/quarnster/Library/Application Support/Sublime Text 3/Packages/Default/Default (Linux).sublime-keymap",
		"/Users/quarnster/Library/Application Support/Sublime Text 3/Packages/Default/Default (OSX).sublime-keymap",
		"/Users/quarnster/Library/Application Support/Sublime Text 3/Packages/Default/Default (Windows).sublime-keymap",
	}
	for _, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			var bindings KeyBindings
			if err := LoadJSON(d, &bindings); err != nil {
				t.Error(err)
			} else {
				t.Log(bindings)
			}
		}
	}
}
