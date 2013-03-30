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
	}
	d0 := ""
	for i, fn := range files {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			if i == 0 {
				d0 = string(d)
			}
			var l textmate.Language
			if err := LoadPlist(d, &l); err != nil {
				t.Fatal(err)
			} else {
				t2[l.ScopeName] = &l
			}
		}
	}
	textmate.Provider = t2
	lp := textmate.LanguageParser{Language: t2["text.xml.plist"]}

	lp.Parse(d0)

	const expfile = "testdata/plist.tmlang"
	str := fmt.Sprintf("%s", lp.RootNode())
	if d, err := ioutil.ReadFile(expfile); err != nil {
		if err := ioutil.WriteFile(expfile, []byte(str), 0644); err != nil {
			t.Error(err)
		}
	} else if diff := util.Diff(str, string(d)); diff != "" {
		t.Error(diff)
	}
}
