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

	buffer := Buffer{data: init}
	tests := []Test{
		{NewInsertAction(&buffer, 0, "hello"), "hellohello world"},
		{NewInsertAction(&buffer, 1, "hello"), "hhelloello world"},
		{NewInsertAction(&buffer, 11, "hello"), "hello worldhello"},
		{NewInsertAction(&buffer, 12, "hello"), "hello worldhello"},
		{NewEraseAction(&buffer, Region{0, 0}), "hello world"},
		{NewEraseAction(&buffer, Region{0, 1}), "ello world"},
		{NewEraseAction(&buffer, Region{1, 0}), "ello world"},
		{NewEraseAction(&buffer, Region{0, 3}), "lo world"},
		{NewEraseAction(&buffer, Region{3, 0}), "lo world"},
		{NewEraseAction(&buffer, Region{3, 6}), "helworld"},
		{NewEraseAction(&buffer, Region{12, 12}), "hello world"},
		{NewReplaceAction(&buffer, Region{0, 0}, "test"), "testhello world"},
		{NewReplaceAction(&buffer, Region{0, 1}, "H"), "Hello world"},
		{NewReplaceAction(&buffer, Region{1, 0}, "H"), "Hello world"},
		{NewReplaceAction(&buffer, Region{0, 3}, "'"), "'lo world"},
		{NewReplaceAction(&buffer, Region{3, 0}, "'"), "'lo world"},
		{NewReplaceAction(&buffer, Region{3, 6}, "vetica "), "helvetica world"},
		{NewReplaceAction(&buffer, Region{0, 5}, "goodbye"), "goodbye world"},
		{NewReplaceAction(&buffer, Region{12, 12}, "."), "hello world."},
	}
	for i, test := range tests {
		test.action.Apply()
		if buffer.data != test.expected {
			t.Fatalf("Apply %d, Expected %s, but got %s", i, test.expected, buffer.data)
		}
		test.action.Undo()
		if buffer.data != init {
			t.Fatalf("Undo %d, Expected %s, but got %s", i, init, buffer.data)
		}
	}
}
