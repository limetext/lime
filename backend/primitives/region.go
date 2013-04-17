package primitives

import (
	"fmt"
)

type (
	Region struct {
		A, B int
	}
)

func (r Region) String() string {
	return fmt.Sprintf("(%d, %d)", r.A, r.B)
}

func (r Region) Begin() int {
	return Min(r.A, r.B)
}

func (r Region) End() int {
	return Max(r.A, r.B)
}

func (r Region) Contains(point int) bool {
	return point >= r.Begin() && point <= r.End()
}

func (r Region) Covers(r2 Region) bool {
	return r.Contains(r2.Begin()) && r2.End() <= r.End()
}

func (r Region) Empty() bool {
	return r.A == r.B
}

func (r Region) Size() int {
	return r.End() - r.Begin()
}

func (r Region) Cover(other Region) Region {
	return Region{Min(r.Begin(), other.Begin()), Max(r.End(), other.End())}
}

func (r Region) Clip(other Region) Region {
	return Region{Clamp(other.Begin(), other.End(), r.A), Clamp(other.Begin(), other.End(), r.B)}
}

func (r Region) Intersects(other Region) bool {
	return r == other || r.Intersection(other).Size() > 0
}

func (r Region) Intersection(other Region) (ret Region) {
	if r.Contains(other.Begin()) || other.Contains(r.Begin()) {
		r2 := Region{Max(r.Begin(), other.Begin()), Min(r.End(), other.End())}
		if r2.Size() != 0 {
			ret = r2
		}
	}

	return ret
}
