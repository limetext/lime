// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type (
	ProfileEntry struct {
		Calls   int
		Tottime time.Duration
	}
	Profiler struct {
		mutex sync.Mutex
		data  map[string]ProfileEntry
	}
	ProfToken struct {
		Name  string
		start time.Time
	}
	ProfileResult struct {
		Name string
		ProfileEntry
	}
	prsorter struct {
		data []ProfileResult
		less func(i, j int) bool
	}
)

var Prof = Profiler{data: make(map[string]ProfileEntry)}

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

func (ps *prsorter) Less(i, j int) bool {
	return ps.less(i, j)
}
func (ps *prsorter) Len() int {
	return len(ps.data)
}

func (ps *prsorter) Swap(i, j int) {
	ps.data[i], ps.data[j] = ps.data[j], ps.data[i]
}

func (p *Profiler) SortByName() []ProfileResult {
	ret := p.Results()
	sort.Sort(&prsorter{ret, func(i, j int) bool { return ret[i].Name < ret[j].Name }})
	return ret
}

func (p *Profiler) SortByTotalTime() []ProfileResult {
	ret := p.Results()
	sort.Sort(&prsorter{ret, func(i, j int) bool { return ret[i].Tottime < ret[j].Tottime }})
	return ret
}

func (p *Profiler) SortByAvgTime() []ProfileResult {
	ret := p.Results()
	sort.Sort(&prsorter{ret, func(i, j int) bool { return ret[i].Avg() < ret[j].Avg() }})
	return ret
}

func (p *Profiler) Results() (ret []ProfileResult) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k, v := range p.data {
		ret = append(ret, ProfileResult{k, v})
	}
	return ret
}

func (v *ProfileResult) Avg() time.Duration {
	return v.Tottime / time.Duration(v.Calls)
}

func (v ProfileResult) String() string {
	return fmt.Sprintf("%64s %6d, %20s, %20s", v.Name, v.Calls, v.Tottime, v.Avg())
}

func (p *Profiler) String() (ret string) {
	pr := p.SortByAvgTime()
	ret = fmt.Sprintf("%64s %6s, %20s, %20s\n", "Name", "Calls", "Total Time", "Average")
	for _, v := range pr {
		ret += fmt.Sprintf("%s\n", v)
	}
	return ret
}
