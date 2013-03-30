package textmate

import (
	"fmt"
)

type (
	Regex string

	Language struct {
		FileTypes      []string
		FirstLineMatch string
		Patterns       []Pattern
		Repository     map[string]Pattern
		ScopeName      string
	}

	Named struct {
		Name string
	}

	Captures map[string]Named

	Pattern struct {
		Named
		Include       string
		Match         Regex
		Captures      Captures
		Begin         Regex
		BeginCaptures Captures
		End           Regex
		EndCaptures   Captures
		Patterns      []Pattern
	}
)

func (p Pattern) String() string {
	return fmt.Sprintf("%-60s %s", p.Name, p.Match)
}

func (s Language) String() (ret string) {
	ret = s.ScopeName + "\n"
	for i := range s.Patterns {
		ret += fmt.Sprintf("\t%s\n", s.Patterns[i])
	}
	return
}
