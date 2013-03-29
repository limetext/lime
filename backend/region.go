package backend

import "fmt"

type (
	Region struct {
		A, B int
	}
	RegionSet struct {
		Regions []Region
	}
)

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

func (r Region) Clip(other Region) Region {
	return Region{clamp(other.Begin(), other.End(), r.A), clamp(other.Begin(), other.End(), r.B)}
}

func (r *RegionSet) adjust(position, delta int) {
	for i := range r.Regions {
		if r.Regions[i].A > position {
			r.Regions[i].A += delta
		}
		if r.Regions[i].B > position {
			r.Regions[i].B += delta
		}
	}
}
