package primitives

import (
	"testing"
)

func TestActionsUtf(t *testing.T) {
	type Test struct {
		action   Action
		expected string
	}
	const init = "€þıœəßðĸʒ×ŋµåäö𝄞"

	buffer := Buffer{data: []rune(init)}
	tests := []Test{
		{NewInsertAction(&buffer, 0, "𝄞€ŋ"), "𝄞€ŋ€þıœəßðĸʒ×ŋµåäö𝄞"},
		{NewInsertAction(&buffer, 1, "𝄞€ŋ"), "€𝄞€ŋþıœəßðĸʒ×ŋµåäö𝄞"},
		{NewInsertAction(&buffer, 10, "𝄞€ŋ"), "€þıœəßðĸʒ×𝄞€ŋŋµåäö𝄞"},
		{NewEraseAction(&buffer, Region{0, 0}), "€þıœəßðĸʒ×ŋµåäö𝄞"},
		{NewEraseAction(&buffer, Region{0, 1}), "þıœəßðĸʒ×ŋµåäö𝄞"},
		{NewEraseAction(&buffer, Region{0, 3}), "œəßðĸʒ×ŋµåäö𝄞"},
		{NewEraseAction(&buffer, Region{6, 9}), "€þıœəß×ŋµåäö𝄞"},
		{NewReplaceAction(&buffer, Region{3, 6}, "𝄞öäå"), "€þı𝄞öäåðĸʒ×ŋµåäö𝄞"},
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
