package commands

import (
	//"github.com/quarnster/util/text"
	. "lime/backend"
	"strings"
	"unicode"
)

type (
	// The TitleCaseCommand transforms all selections
	// to be in Title Case.  For instance, the text:
	// "this is some sample text"
	// turns in to:
	// "This Is Some Sample Text"
	TitleCaseCommand struct {
		DefaultCommand
	}

	// The SwapCaseCommand transforms all selections
	// so that each character in the selection
	// is the opposite case.  For example, the text:
	// "Hello, World!"
	// turns in to:
	// "hELLO, wORLD!"
	SwapCaseCommand struct {
		DefaultCommand
	}
)

func (c *TitleCaseCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		if r.Size() == 0 {
			// return error or just ignore?
		} else {
			//v.Replace(e, r, c.Characters)
			t := v.Buffer().Substr(r)
			v.Replace(e, r, strings.Title(t))
		}
	}
	return nil
}

func (c *SwapCaseCommand) Run(v *View, e *Edit) error {
	sel := v.Sel()
	for i := 0; i < sel.Len(); i++ {
		r := sel.Get(i)
		if r.Size() == 0 {
			// return error or just ignore?
		} else {
			//v.Replace(e, r, c.Characters)
			text := v.Buffer().Substr(r)
			swapped := make([]rune, 0)
			for _, c := range text {
				if unicode.IsUpper(c) {
					swapped = append(swapped, unicode.ToLower(c))
				} else {
					swapped = append(swapped, unicode.ToUpper(c))
				}
			}
			v.Replace(e, r, string(swapped))
		}
	}
	return nil
}

func init() {
	register([]cmd{
		{"title_case", &TitleCaseCommand{}},
		{"swap_case", &SwapCaseCommand{}},
	})
}
