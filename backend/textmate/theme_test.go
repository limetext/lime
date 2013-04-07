package textmate

import (
	"fmt"
	"github.com/quarnster/completion/util"
	"io/ioutil"
	"lime/backend/loaders"
	"testing"
)

func TestLoadTheme(t *testing.T) {
	type Test struct {
		in  string
		out string
	}
	tests := []Test{
		{"testdata/Monokai.tmTheme", "testdata/Monokai.tmTheme.res"},
	}
	for _, test := range tests {
		if d, err := ioutil.ReadFile(test.in); err != nil {
			t.Logf("Couldn't load file %s: %s", test.in, err)
		} else {
			var theme Theme
			if err := loaders.LoadPlist(d, &theme); err != nil {
				t.Error(err)
			} else {
				str := fmt.Sprintf("%s", theme)
				if d, err := ioutil.ReadFile(test.out); err != nil {
					if err := ioutil.WriteFile(test.out, []byte(str), 0644); err != nil {
						t.Error(err)
					}
				} else if diff := util.Diff(string(d), str); diff != "" {
					t.Error(diff)
				}

			}
		}
	}
}
