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

	e = v.BeginEdit()
	ed.CommandHandler().RunTextCommand(v, "title_case", nil)
	v.EndEdit(e)

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

	e = v.BeginEdit()
	ed.CommandHandler().RunTextCommand(v, "title_case", nil)
	v.EndEdit(e)

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

	e = v.BeginEdit()
	ed.CommandHandler().RunTextCommand(v, "title_case", nil)
	v.EndEdit(e)

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

	e = v.BeginEdit()
	ed.CommandHandler().RunTextCommand(v, "title_case", nil)
	v.EndEdit(e)

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

	e = v.BeginEdit()
	ed.CommandHandler().RunTextCommand(v, "swap_case", nil)
	v.EndEdit(e)

	result := v.Buffer().Substr(Region{0, v.Buffer().Size()})

	if result != helloSwapped {

		t.Errorf(`TitleCaseCommand Failed.
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

	e = v.BeginEdit()
	ed.CommandHandler().RunTextCommand(v, "swap_case", nil)
	v.EndEdit(e)

	result = v.Buffer().Substr(Region{0, v.Buffer().Size()})

	if result != privetSwapped {

		t.Errorf(`TitleCaseCommand Failed.
			Expected: %s
			Got: %s`, privetSwapped, result)
	}

}
