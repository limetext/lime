package backend

import (
	"io/ioutil"
	"testing"
)

func TestPlist(t *testing.T) {
	tests := []string{
		"../3rdparty/bundles/c.tmbundle/Syntaxes/C.plist",
	}
	for _, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			var syn Syntax
			if err := LoadPlist(d, &syn); err != nil {
				t.Error(err)
			} else {
				t.Log(syn)
			}
		}
	}
}
