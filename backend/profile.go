package backend

import (
	"fmt"
	"sync"
	"time"
)

type (
	Profileentry struct {
		Calls   int
		Tottime time.Duration
	}
	Profiler struct {
		mutex sync.Mutex
		data  map[string]Profileentry
	}
	ProfToken struct {
		Name  string
		start time.Time
	}
)

var Prof = Profiler{data: make(map[string]Profileentry)}

func (p *Profiler) Enter(name string) ProfToken {
	return ProfToken{name, time.Now()}
}

func (pt *ProfToken) Exit() {
	Prof.mutex.Lock()
	defer Prof.mutex.Unlock()
	e := Prof.data[pt.Name]
	e.Calls++
	e.Tottime += time.Since(pt.start)
	Prof.data[pt.Name] = e
}

func (p Profiler) String() string {
	ret := ""
	for k, v := range p.data {
		ret += fmt.Sprintf("%s %d, %s, %s\n", k, v.Calls, v.Tottime, v.Tottime/time.Duration(v.Calls))
	}
	return ret
}
