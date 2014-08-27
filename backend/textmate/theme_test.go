// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package textmate

import (
	"fmt"
	"github.com/limetext/lime/backend/loaders"
	"github.com/limetext/lime/backend/util"
	"io/ioutil"
	"testing"
)

func TestLoadTheme(t *testing.T) {
	type Test struct {
		in  string
		out string
	}
	tests := []Test{
		{
			"testdata/Monokai.tmTheme",
			"testdata/Monokai.tmTheme.res",
		},
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

func TestLoadThemeFromPlist(t *testing.T) {
	f := "testdata/Monokai.tmTheme"
	th, err := LoadTheme(f)
	if err != nil {
		t.Errorf("Tried to load %s, but got an error: %v", f, err)
	}

	n := "Monokai"
	if th.Name != n {
		t.Errorf("Tried to load %s, but got %s", f, th)
	}
}

func TestLoadThemeFromNonPlist(t *testing.T) {
	f := "testdata/Monokai.tmTheme.res"
	_, err := LoadTheme(f)
	if err == nil {
		t.Errorf("Tried to load %s, expecting an error, but didn't get one", f)
	}
}

func TestLoadThemeFromMissingFile(t *testing.T) {
	f := "testdata/MissingFile"
	_, err := LoadTheme(f)
	if err == nil {
		t.Errorf("Tried to load %s, expecting an error, but didn't get one", f)
	}
}
