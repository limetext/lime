package primitives

import (
	"testing"
)

func TestActions(t *testing.T) {
	type Test struct {
		action   Action
		expected string
	}
	const init = "hello world"
	var buffer = NewBuffer()
	buffer.Insert(0, init)
	tests := []Test{
		{NewInsertAction(buffer, 0, "hello"), "hellohello world"},
		{NewInsertAction(buffer, 1, "hello"), "hhelloello world"},
		{NewInsertAction(buffer, 11, "hello"), "hello worldhello"},
		{NewInsertAction(buffer, 12, "hello"), "hello worldhello"},
		{NewEraseAction(buffer, Region{0, 0}), "hello world"},
		{NewEraseAction(buffer, Region{0, 1}), "ello world"},
		{NewEraseAction(buffer, Region{1, 0}), "ello world"},
		{NewEraseAction(buffer, Region{0, 3}), "lo world"},
		{NewEraseAction(buffer, Region{3, 0}), "lo world"},
		{NewEraseAction(buffer, Region{3, 6}), "helworld"},
		{NewEraseAction(buffer, Region{12, 12}), "hello world"},
		{NewReplaceAction(buffer, Region{0, 0}, "test"), "testhello world"},
		{NewReplaceAction(buffer, Region{0, 1}, "H"), "Hello world"},
		{NewReplaceAction(buffer, Region{1, 0}, "H"), "Hello world"},
		{NewReplaceAction(buffer, Region{0, 3}, "'"), "'lo world"},
		{NewReplaceAction(buffer, Region{3, 0}, "'"), "'lo world"},
		{NewReplaceAction(buffer, Region{3, 6}, "vetica "), "helvetica world"},
		{NewReplaceAction(buffer, Region{0, 5}, "goodbye"), "goodbye world"},
		{NewReplaceAction(buffer, Region{12, 12}, "."), "hello world."},
	}
	for i, test := range tests {
		test.action.Apply()
		if d := buffer.String(); d != test.expected {
			t.Errorf("Apply %d, Expected %s, but got %s", i, test.expected, d)
		}
		test.action.Undo()
		if d := buffer.String(); d != init {
			t.Errorf("Undo %d, Expected %s, but got %s", i, init, d)
		}
	}
	bigdata := make([]rune, 1024*1024)
	fill(bigdata)
	buffer.Insert(0, string(bigdata))
	a := NewInsertAction(buffer, buffer.Size(), "test")
	a.Apply()
	if end := buffer.Substr(Region{buffer.Size() - 4, buffer.Size()}); end != "test" {
		t.Errorf("%s != test", end)
	}
}

func TestActionsUtf(t *testing.T) {
	type Test struct {
		action   Action
		expected string
	}
	const init = "€þıœəßðĸʒ×ŋµåäö𝄞"
	var buffer = NewBuffer()
	buffer.Insert(0, init)

	tests := []Test{
		{NewInsertAction(buffer, 0, "𝄞€ŋ"), "𝄞€ŋ€þıœəßðĸʒ×ŋµåäö𝄞"},
		{NewInsertAction(buffer, 1, "𝄞€ŋ"), "€𝄞€ŋþıœəßðĸʒ×ŋµåäö𝄞"},
		{NewInsertAction(buffer, 10, "𝄞€ŋ"), "€þıœəßðĸʒ×𝄞€ŋŋµåäö𝄞"},
		{NewEraseAction(buffer, Region{0, 0}), "€þıœəßðĸʒ×ŋµåäö𝄞"},
		{NewEraseAction(buffer, Region{0, 1}), "þıœəßðĸʒ×ŋµåäö𝄞"},
		{NewEraseAction(buffer, Region{0, 3}), "œəßðĸʒ×ŋµåäö𝄞"},
		{NewEraseAction(buffer, Region{6, 9}), "€þıœəß×ŋµåäö𝄞"},
		{NewReplaceAction(buffer, Region{3, 6}, "𝄞öäå"), "€þı𝄞öäåðĸʒ×ŋµåäö𝄞"},
	}
	for i, test := range tests {
		test.action.Apply()
		if d := buffer.String(); d != test.expected {
			t.Errorf("Apply %d\n\tExpected %v\n\tBut got  %v", i, []rune(test.expected), []rune(d))
		}
		test.action.Undo()
		if d := buffer.String(); d != init {
			t.Errorf("Undo %d, Expected %s, but got %s", i, init, d)
		}
	}
}
