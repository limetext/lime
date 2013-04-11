package backend

import (
	"fmt"
	"time"
)

type (
	Profileentry struct {
		Calls   int
		Tottime time.Duration
	}
	Profiler  map[string]Profileentry
	ProfToken struct {
		Name  string
		start time.Time
	}
)

var Prof = make(Profiler)

func (p *Profiler) Enter(name string) ProfToken {
	return ProfToken{name, time.Now()}
}

func (pt *ProfToken) Exit() {
	e := Prof[pt.Name]
	e.Calls++
	e.Tottime += time.Since(pt.start)
	Prof[pt.Name] = e
}

func (p Profiler) String() string {
	ret := ""
	for k, v := range p {
		ret += fmt.Sprintf("%s %d, %s, %s\n", k, v.Calls, v.Tottime, v.Tottime/time.Duration(v.Calls))
	}
	return ret
}
