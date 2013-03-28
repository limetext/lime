package backend

import "fmt"

type Region struct {
	A, B int
}

func (r Region) String() string {
	return fmt.Sprintf("(%d, %d)", r.A, r.B)
}

func (r Region) Begin() int {
	return min(r.A, r.B)
}

func (r Region) End() int {
	return max(r.A, r.B)
}

func (r Region) Contains(point int) bool {
	return point >= r.Begin() && point < r.End()
}

func (r Region) Empty() bool {
	return r.A == r.B
}

func (r Region) Size() int {
	return r.End() - r.Begin()
}

func (r Region) Cover(other Region) Region {
	return Region{min(r.Begin(), other.Begin()), max(r.End(), other.End())}
}
