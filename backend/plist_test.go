package backend

import (
	"io/ioutil"
	"lime/backend/textmate"
	"testing"
)

func TestPlist(t *testing.T) {
	tests := []string{
		"../3rdparty/bundles/property-list.tmbundle/Syntaxes/Property List (XML).tmLanguage",
	}
	for _, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			var syn textmate.Language
			if err := LoadPlist(d, &syn); err != nil {
				t.Error(err)
			} else {
				t.Log(syn)
			}
		}
	}
}

func TestPlist2(t *testing.T) {
	tests := []string{
		"../3rdparty/bundles/monokai.tmbundle/Themes/Monokai.tmTheme",
	}
	for _, fn := range tests {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			var theme textmate.Theme
			if err := LoadPlist(d, &theme); err != nil {
				t.Error(err)
			} else {
				t.Log(theme)
			}
		}
	}
}
