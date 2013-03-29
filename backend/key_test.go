package backend

import (
	"io/ioutil"
	"testing"
)

func TestKeyFilter(t *testing.T) {
	fn := "json/testdata/Default (OSX).sublime-keymap"
	if d, err := ioutil.ReadFile(fn); err != nil {
		t.Logf("Couldn't load file %s: %s", fn, err)
	} else {
		var bindings KeyBindings
		if err := LoadJSON(d, &bindings); err != nil {
			t.Error(err)
		} else {
			t.Log(bindings)
		}

		if b2 := bindings.Filter(KeyPress{Key: 'j', Ctrl: true}); b2.Len() != 3 {
			t.Errorf("Not of the expected length: %d, %s", 3, b2)
		} else if b3 := b2.Filter(KeyPress{Key: 's'}); b3.Len() != 1 {
			t.Errorf("Not of the expected length: %d, %s", 1, b3)
		}
	}
}
