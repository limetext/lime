package primitives

import (
	"testing"
)

// Verified against ST3
func TestRegionIntersects(t *testing.T) {
	type Test struct {
		a, b Region
		c    bool
	}
	var tests = []Test{
		{Region{10, 20}, Region{25, 35}, false},
		{Region{25, 35}, Region{10, 20}, false},
		{Region{10, 25}, Region{20, 35}, true},
		{Region{20, 35}, Region{10, 25}, true},
		{Region{10, 25}, Region{15, 20}, true},
		{Region{15, 20}, Region{10, 25}, true},
		{Region{5, 10}, Region{10, 23}, false},
		{Region{5, 10}, Region{5, 10}, true},
		{Region{0, 0}, Region{0, 0}, true},
		{Region{1, 1}, Region{1, 1}, true},
		{Region{23, 24}, Region{10, 23}, false},
	}
	for _, test := range tests {
		if res := test.a.Intersects(test.b); res != test.c {
			t.Errorf("Expected %v, but got %v, %v", test.c, res, test)
		}
	}
}

// Verified against ST3
func TestRegionIntersection(t *testing.T) {
	var tests = [][]Region{
		{{10, 20}, {25, 35}, {0, 0}},
		{{25, 35}, {10, 20}, {0, 0}},
		{{10, 25}, {20, 35}, {20, 25}},
		{{20, 35}, {10, 25}, {20, 25}},
		{{10, 25}, {15, 20}, {15, 20}},
		{{15, 20}, {10, 25}, {15, 20}},
		{{5, 10}, {10, 23}, {0, 0}},
		{{5, 10}, {5, 10}, {5, 10}},
		{{1, 1}, {1, 1}, {0, 0}},
	}
	for _, test := range tests {
		if res := test[0].Intersection(test[1]); res != test[2] {
			t.Errorf("Expected intersection %v, but got %v, %v", test[2], res, test)
		}
	}
}

// Verified against ST3
func TestContains(t *testing.T) {
	type Test struct {
		r   Region
		pos int
		c   bool
	}
	tests := []Test{
		{Region{0, 0}, 0, true},
		{Region{10, 10}, 10, true},
		{Region{10, 11}, 10, true},
		{Region{10, 11}, 11, true},
		{Region{10, 11}, 12, false},
		{Region{10, 11}, 9, false},
	}
	for _, test := range tests {
		if res := test.r.Contains(test.pos); res != test.c {
			t.Errorf("Expected %v, but got %v, %v, %v", test.c, res, test.r, test.pos)
		}
	}
}
