package backend

import (
	"errors"
	"fmt"
	"github.com/quarnster/completion/util"
	"io/ioutil"
	"lime/backend/textmate"
	"testing"
)

type tmp map[string]*textmate.Language

func (t tmp) GetLanguage(id string) (*textmate.Language, error) {
	if v, ok := t[id]; !ok {
		return nil, errors.New("Can't handle id " + id)
	} else {
		return v, nil
	}
}

func TestTmLanguage(t *testing.T) {
	t2 := make(tmp)

	files := []string{
		"../3rdparty/bundles/property-list.tmbundle/Syntaxes/Property List (XML).tmLanguage",
		"../3rdparty/bundles/xml.tmbundle/Syntaxes/XML.plist",
		"../3rdparty/bundles/go.tmbundle/Syntaxes/Go.tmLanguage",
	}
	for _, fn := range files {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Errorf("Couldn't load file %s: %s", fn, err)
		} else {
			var l textmate.Language
			if err := LoadPlist(d, &l); err != nil {
				t.Fatal(err)
			} else {
				t2[l.ScopeName] = &l
			}
		}
	}
	textmate.Provider = t2

	type test struct {
		in  string
		out string
		syn string
	}
	tests := []test{
		{
			"../3rdparty/bundles/property-list.tmbundle/Syntaxes/Property List (XML).tmLanguage",
			"testdata/plist.tmlang",
			"text.xml.plist",
		},
		{
			"textmate_language_test.go",
			"testdata/go.tmlang",
			"source.go",
		},
	}
	for _, t3 := range tests {
		lp := textmate.LanguageParser{Language: t2[t3.syn]}

		var d0 string
		if d, err := ioutil.ReadFile(t3.in); err != nil {
			t.Errorf("Couldn't load file %s: %s", t3.in, err)
			continue
		} else {
			d0 = string(d)
		}
		lp.Parse(d0)

		str := fmt.Sprintf("%s", lp.RootNode())
		if d, err := ioutil.ReadFile(t3.out); err != nil {
			if err := ioutil.WriteFile(t3.out, []byte(str), 0644); err != nil {
				t.Error(err)
			}
		} else if diff := util.Diff(string(d), str); diff != "" {
			t.Error(diff)
		}
	}
}
