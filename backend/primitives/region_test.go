package primitives

import (
	"testing"
)

func TestRegionIntersection(t *testing.T) {
	var tests = [][]Region{
		{{10, 20}, {25, 35}, {0, 0}},
		{{25, 35}, {10, 20}, {0, 0}},
		{{10, 25}, {20, 35}, {20, 25}},
		{{20, 35}, {10, 25}, {20, 25}},
		{{10, 25}, {15, 20}, {15, 20}},
		{{15, 20}, {10, 25}, {15, 20}},
	}
	for _, test := range tests {
		if res := test[0].Intersection(test[1]); res != test[2] {
			t.Errorf("Expected intersection %v, but got %v", test[2], res)
		}
	}
}
