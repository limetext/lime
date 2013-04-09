package textmate

import (
	"fmt"
	"github.com/quarnster/completion/util"
	"io/ioutil"
	"testing"
)

func TestTmLanguage(t *testing.T) {
	files := []string{
		"testdata/Property List (XML).tmLanguage",
		"testdata/XML.plist",
		"testdata/Go.tmLanguage",
	}
	for _, fn := range files {
		if err := Provider.Load(fn); err != nil {
			t.Fatal(err)
		}
	}

	type test struct {
		in  string
		out string
		syn string
	}
	tests := []test{
		{
			"testdata/Property List (XML).tmLanguage",
			"testdata/plist.tmlang",
			"text.xml.plist",
		},
		{
			"testdata/main.go",
			"testdata/main.go.res",
			"source.go",
		},
		{
			"testdata/utf.go",
			"testdata/utf.go.res",
			"source.go",
		},
	}
	for _, t3 := range tests {
		l, err := Provider.GetLanguage(t3.syn)
		if err != nil {
			t.Error(err)
			continue
		}
		lp := LanguageParser{Language: l}

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
