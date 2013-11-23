// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
package commands

import (
	. "github.com/quarnster/util/text"
	. "lime/backend"
	"testing"
)

func TestTitleCase(t *testing.T) {

	quote := "Give a man a match, and he'll be warm for a minute, but set him on fire, and he'll be warm for the rest of his life."
	expected := "Give a man a match, and He'Ll Be Warm For A Minute, but set him on fire, and he'll be warm for the rest of his life."
	// Please note the bizarre  capitalization of the first L in he'Ll...  This is due to a bug in go's strings
	// library.  I'm going to try to get them to fix it...  If not, maybe we'll have
	// to write our own Title Casing function.
	expectedMulti := "Give A Man A Match, and he'll be warm for a minute, But Set Him On Fire, and he'll be warm for the rest of his life."

	// SINGLE SELECTION
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, quote)
	v.EndEdit(e)

	v.Sel().Clear()

	v.Sel().Add(Region{24, 51})

	ed.CommandHandler().RunTextCommand(v, "title_case", nil)

	result := v.Buffer().Substr(Region{0, v.Buffer().Size()})

	if result != expected {
		t.Errorf(`TitleCaseCommand Failed.
			Expected: %s
			Got: %s`, expected, result)
	}

	// MULTIPLE SELECTIONS
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, quote)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, 17})
	v.Sel().Add(Region{52, 71})

	ed.CommandHandler().RunTextCommand(v, "title_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})
	if result != expectedMulti {
		t.Errorf(`TitleCaseCommand Failed.
			Expected: %s
			Got: %s`, expectedMulti, result)
	}

	// NO SELECTION
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, quote)
	v.EndEdit(e)
	v.Sel().Clear()

	ed.CommandHandler().RunTextCommand(v, "title_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})
	if result != quote {
		t.Errorf(`TitleCaseCommand Failed.
			Expected: %s
			Got: %s`, quote, result)
	}

	// UNICODE CHARACTER TEST
	russian := "ничего себе!"
	titleRussian := "Ничего Себе!"
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, russian)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "title_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})
	if result != titleRussian {
		t.Errorf(`TitleCaseCommand Failed.
			Expected: %s
			Got: %s`, titleRussian, result)
	}
}

func TestSwapCase(t *testing.T) {

	hello := "Hello, World!"
	helloSwapped := "hELLO, wORLD!"
	privet := "ПрИвЕт, МиР"
	privetSwapped := "пРиВеТ, мИр"

	// ASCII Test
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, hello)
	v.EndEdit(e)

	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "swap_case", nil)

	result := v.Buffer().Substr(Region{0, v.Buffer().Size()})
	if result != helloSwapped {
		t.Errorf(`SwapCaseCommand Failed.
			Expected: %s
			Got: -%s-`, helloSwapped, result)
	}

	// Unicode Test
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, privet)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "swap_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})
	if result != privetSwapped {
		t.Errorf(`SwapCaseCommand Failed.
			Expected: %s
			Got: %s`, privetSwapped, result)
	}
}

func TestUpperCase(t *testing.T){
	en := "Try not to become a man of success, but rather try to become a man of value."
	ru := "чем больше законов и постановлений, тем больше разбойников и преступлений!"
	zh := "千里之行﹐始于足下"
	enup := "TRY NOT TO BECOME A MAN OF SUCCESS, BUT RATHER TRY TO BECOME A MAN OF VALUE."
	ruup := "ЧЕМ БОЛЬШЕ ЗАКОНОВ И ПОСТАНОВЛЕНИЙ, ТЕМ БОЛЬШЕ РАЗБОЙНИКОВ И ПРЕСТУПЛЕНИЙ!"

	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	
	// ASCII Test
	e := v.BeginEdit()
	v.Insert(e, 0, en)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "upper_case", nil)

	result := v.Buffer().Substr(Region{0, v.Buffer().Size()})

	if result != enup {
		t.Errorf(`UpperCaseCommand Failed.
			Expected: %s
			Got: %s`, enup, result)
	}

	// Unicode Test
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, ru)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "upper_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})	

	if result != ruup {
		t.Errorf(`UpperCaseCommand Failed.
			Expected: %s
			Got: %s`, ruup, result)
	}

	// Unicode (Chinese) Test
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, zh)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "upper_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})	

	if result != zh {
		t.Errorf(`UpperCaseCommand Failed.
			Expected: %s
			Got: %s`, zh, result)
	}

}

func TestLowerCase(t *testing.T){
	en := "We make a Living by WHAt we get, but WE make a lIfe by whAt we gIVe."
	ru := "Все счастливые сЕмьи пОхОжи друг на друга, КАЖДАЯ несчастливая СемьЯ несчастлива ПО-СВоЕмУ."
	ja := "行動のともなわないビジョンは、ただの白日夢。ビジョンのない行動は、ただの悪夢。"
	enlo := "we make a living by what we get, but we make a life by what we give."
	rulo := "все счастливые семьи похожи друг на друга, каждая несчастливая семья несчастлива по-своему."

	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	
	// ASCII Test
	e := v.BeginEdit()
	v.Insert(e, 0, en)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "lower_case", nil)

	result := v.Buffer().Substr(Region{0, v.Buffer().Size()})

	if result != enlo {
		t.Errorf(`LowerCaseCommand Failed.
			Expected: %s
			Got: %s`, enlo, result)
	}

	// Unicode Test
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, ru)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "lower_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})	

	if result != rulo {
		t.Errorf(`LowerCaseCommand Failed.
			Expected: %s
			Got: %s`, rulo, result)
	}

	// Unicode (Japanese) Test
	v = w.NewFile()
	e = v.BeginEdit()
	v.Insert(e, 0, ja)
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{0, v.Buffer().Size()})

	ed.CommandHandler().RunTextCommand(v, "lower_case", nil)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})	

	if result != ja {
		t.Errorf(`LowerCaseCommand Failed.
			Expected: %s
			Got: %s`, ja, result)
	}

}
