// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"testing"
)

type CaseTest struct {
	in_region []Region
	in        string
	exp       string
}

func RunCaseTest(command string, testsuite *[]CaseTest, t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()

	for i, test := range *testsuite {
		v := w.NewFile()
		e := v.BeginEdit()
		v.Insert(e, 0, test.in)
		v.EndEdit(e)

		v.Sel().Clear()
		if test.in_region != nil {
			for _, r := range test.in_region {
				v.Sel().Add(r)
			}
		}
		ed.CommandHandler().RunTextCommand(v, command, nil)
		sr := v.Buffer().Substr(Region{0, v.Buffer().Size()})
		if sr != test.exp {
			t.Errorf("%s test %d failed: %v, %+v", command, i, sr, test)
		}
	}
}

func TestTitleCase(t *testing.T) {
	tests := []CaseTest{
		/*single selection*/
		{
			// Please note the bizarre  capitalization of the first L in he'Ll...  This is due to a bug in go's strings
			// library.  I'm going to try to get them to fix it...  If not, maybe we'll have
			// to write our own Title Casing function.
			[]Region{{24, 51}},
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give a man a match, and He'Ll Be Warm For A Minute, but set him on fire, and he'll be warm for the rest of his life.",
		},
		/*multiple selection*/
		{
			[]Region{{0, 17}, {52, 71}},
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give A Man A Match, and he'll be warm for a minute, But Set Him On Fire, and he'll be warm for the rest of his life.",
		},

		/*no selection*/
		{
			nil,
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
			"Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life.",
		},
		/*unicode*/
		{
			[]Region{{0, 12}},
			"ничего себе!",
			"Ничего Себе!",
		},
		/*asian characters*/

		{
			[]Region{{0, 9}},
			"千里之行﹐始于足下",
			"千里之行﹐始于足下",
		},
	}
	RunCaseTest("title_case", &tests, t)

}

func TestSwapCase(t *testing.T) {
	tests := []CaseTest{
		{
			[]Region{{0, 13}},
			"Hello, World!",
			"hELLO, wORLD!",
		},
		{
			[]Region{{0, 11}},
			"ПрИвЕт, МиР",
			"пРиВеТ, мИр",
		},
	}
	RunCaseTest("swap_case", &tests, t)
}

func TestUpperCase(t *testing.T) {
	tests := []CaseTest{
		/*single selection*/

		{
			[]Region{{0, 76}},
			"Try not to become a man of success, but rather try to become a man of value.",
			"TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE.",
		},
		/*multiple selection*/
		{
			[]Region{{0, 20}, {74, 76}},

			"Try not to become a man of success, but rather try to become a man of value.",
			"TRY NOT TO BECOME A man of success, but rather try to become a man of valuE.",
		},
		/*no selection*/
		{
			nil,

			"Try not to become a man of success, but rather try to become a man of value.",
			"Try not to become a man of success, but rather try to become a man of value.",
		},
		/*unicode*/

		{
			[]Region{{0, 74}},
			"чем больше законов и постановлений, тем больше разбойников и преступлений!",
			"ЧЕМ БОЛЬШЕ ЗАКОНОВ И ПОСТАНОВЛЕНИЙ, ТЕМ БОЛЬШЕ РАЗБОЙНИКОВ И ПРЕСТУПЛЕНИЙ!",
		},
		/*asian characters*/

		{
			[]Region{{0, 9}},
			"千里之行﹐始于足下",
			"千里之行﹐始于足下",
		},
	}
	RunCaseTest("upper_case", &tests, t)
}

func TestLowerCase(t *testing.T) {
	tests := []CaseTest{
		/*single selection*/
		{
			[]Region{{0, 76}},

			"TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE.",
			"try not to become a man of success, but rather try to become a man of value.",
		},
		/*multiple selection*/
		{
			[]Region{{0, 20}, {74, 76}},

			"TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE.",
			"try not to become a MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUe.",
		},
		/*no selection*/
		{
			nil,

			"Try not to become a man of success, but rather try to become a man of value.",
			"Try not to become a man of success, but rather try to become a man of value.",
		},
		/*unicode*/
		{
			[]Region{{0, 74}},

			"ЧЕМ БОЛЬШЕ ЗАКОНОВ И ПОСТАНОВЛЕНИЙ, ТЕМ БОЛЬШЕ РАЗБОЙНИКОВ И ПРЕСТУПЛЕНИЙ!",
			"чем больше законов и постановлений, тем больше разбойников и преступлений!",
		},
		/*asian characters*/

		{
			[]Region{{0, 9}},
			"千里之行﹐始于足下",
			"千里之行﹐始于足下",
		},
	}
	RunCaseTest("lower_case", &tests, t)
}
