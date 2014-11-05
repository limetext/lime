package watch

import "testing"

type testAction int

func (t *testAction) double() {
	(*t) = (*t) * 2
}

func (t *testAction) triple() {
	(*t) = (*t) * 3
}

func TestApplyAction(t *testing.T) {
	var ta testAction = 1
	tests := []struct {
		acFn func()
		acEv int
		ev   int
		exp  int
	}{
		{ta.double, CREATE, CREATE, 2},
		{ta.double, CREATE, RENAME, 1},
		{ta.triple, CREATE | DELETE, DELETE, 3},
		{ta.triple, MODIFY, MODIFY | CREATE, 3},
	}
	for i, test := range tests {
		ac := action{test.acFn, test.acEv}
		ac.apply(test.ev)
		if int(ta) != test.exp {
			t.Errorf("Test %d: Expected ta variable be %d, but got %d", i, test.exp, ta)
		}
		ta = 1
	}
}

func TestApplyActions(t *testing.T) {
	var ta testAction = 1
	acs := actions{"key1": action{ta.double, CREATE}, "key2": {ta.triple, CREATE}}
	acs.apply(ALL, "key1")
	if int(ta) != 2 {
		t.Errorf("Expected ta variable be %d, but got %d", 2, ta)
	}
}

func TestApplyAll(t *testing.T) {
	var ta testAction = 1
	acs := actions{"key1": action{ta.double, CREATE}, "key2": {ta.triple, CREATE}}
	acs.applyAll(ALL)
	if int(ta) != 6 {
		t.Errorf("Expected ta variable be %d, but got %d", 6, ta)
	}
}

func TestNewEvent(t *testing.T) {
	tests := []struct {
		evs []int
		exp int
	}{
		{
			[]int{CREATE},
			CREATE,
		},
		{
			[]int{CREATE, RENAME, MODIFY, DELETE},
			ALL,
		},
	}
	for i, test := range tests {
		if ev := newEvent(test.evs); ev != test.exp {
			t.Errorf("Test %d: Expected  event %d, but got %d", i, test.exp, ev)
		}
	}
}
