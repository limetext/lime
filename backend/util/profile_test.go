// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"testing"
	"time"
)

func TestSortByName(t *testing.T) {
	Prof.data["A"] = ProfileEntry{Calls: 1, Tottime: time.Duration(5)}
	Prof.data["B"] = ProfileEntry{Calls: 1, Tottime: time.Duration(4)}
	Prof.data["C"] = ProfileEntry{Calls: 1, Tottime: time.Duration(3)}
	Prof.data["D"] = ProfileEntry{Calls: 1, Tottime: time.Duration(2)}
	results := Prof.SortByName()
	if results[0].Name != "A" {
		t.Error("TestSortByName expected A, but got %s", results[0].Name)
	}
	if results[1].Name != "B" {
		t.Error("TestSortByName expected B, but got %s", results[1].Name)
	}
	if results[2].Name != "C" {
		t.Error("TestSortByName expected C, but got %s", results[2].Name)
	}
	if results[3].Name != "D" {
		t.Error("TestSortByName expected D, but got %s", results[3].Name)
	}
}

func TestSortByTime(t *testing.T) {
	Prof.data["A"] = ProfileEntry{Calls: 1, Tottime: time.Duration(1)}
	Prof.data["C"] = ProfileEntry{Calls: 1, Tottime: time.Duration(2)}
	Prof.data["D"] = ProfileEntry{Calls: 1, Tottime: time.Duration(3)}
	Prof.data["B"] = ProfileEntry{Calls: 1, Tottime: time.Duration(4)}
	results := Prof.SortByTotalTime()
	if results[0].Name != "A" {
		t.Error("TestSortByTime expected A, but got %s", results[0].Name)
	}
	if results[1].Name != "C" {
		t.Error("TestSortByTime expected C, but got %s", results[1].Name)
	}
	if results[2].Name != "D" {
		t.Error("TestSortByTime expected D, but got %s", results[2].Name)
	}
	if results[3].Name != "B" {
		t.Error("TestSortByTime expected B, but got %s", results[3].Name)
	}
}

func TestSortByAvgTime(t *testing.T) {
	Prof.data["A"] = ProfileEntry{Calls: 1, Tottime: time.Duration(1)}
	Prof.data["C"] = ProfileEntry{Calls: 1, Tottime: time.Duration(2)}
	Prof.data["D"] = ProfileEntry{Calls: 1, Tottime: time.Duration(3)}
	Prof.data["B"] = ProfileEntry{Calls: 1, Tottime: time.Duration(4)}
	results := Prof.SortByAvgTime()
	if results[0].Name != "A" {
		t.Error("TestSortByTime expected A, but got %s", results[0].Name)
	}
	if results[1].Name != "C" {
		t.Error("TestSortByTime expected C, but got %s", results[1].Name)
	}
	if results[2].Name != "D" {
		t.Error("TestSortByTime expected D, but got %s", results[2].Name)
	}
	if results[3].Name != "B" {
		t.Error("TestSortByTime expected B, but got %s", results[3].Name)
	}
}

func TestEnter(t *testing.T) {
	profiler := &Profiler{}
	pt := profiler.Enter("some name")

	if &pt == nil {
		t.Error("Returned a nil ProfileToken")
	}
}

func TestExit(t *testing.T) {
	ptoken := &ProfToken{}
	ptoken.Exit()
}

func TestLess(t *testing.T) {
	ps := &prsorter{nil, func(i, j int) bool {
		return i < j
	}}
	result := ps.Less(4, 5)
	if !result {
		t.Error("Less returned incorrect result for 4 < 5")
	}
}

func TestLen(t *testing.T) {
	ps := &prsorter{nil, func(i, j int) bool {
		return i < j
	}}
	if ps.Len() != 0 {
		t.Error("Incorrect length of Data (nil data should be 0 len)")
	}
}

func TestSwap(t *testing.T) {
	pr0 := &ProfileResult{"pr1", *&ProfileEntry{}}
	pr1 := &ProfileResult{"pr2", *&ProfileEntry{}}
	prList := []ProfileResult{*pr0, *pr1}

	ps := &prsorter{prList, func(i, j int) bool { return i < j }}
	ps.Swap(0, 1)

	if ps.data[0] != *pr1 || ps.data[1] != *pr0 {
		t.Error("prsorter swapped incorrectly!")
	}
}

func TestString(t *testing.T) {
	Prof.String()
}
