package primitives

import "fmt"

type (
	Region struct {
		A, B int
	}
	RegionSet struct {
		regions []Region
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

func (r *RegionSet) Adjust(position, delta int) {
	for i := range r.regions {
		if r.regions[i].A >= position {
			r.regions[i].A += delta
		}
		if r.regions[i].B >= position {
			r.regions[i].B += delta
		}
	}
}

func (r *RegionSet) Substract(r2 Region) {
	r.Adjust(r2.Begin(), r2.Size())
}

func (r *RegionSet) Add(r2 Region) {
	// TODO: should be kept in sorted order + merged with any intersecting regions
	r.regions = append(r.regions, r2)
}

func (r *RegionSet) Clear() {
	r.regions = r.regions[0:0]
}

func (r *RegionSet) Get(i int) Region {
	return r.regions[i]
}

func (r *RegionSet) Len() int {
	return len(r.regions)
}

func (r *RegionSet) Contains(r2 Region) bool {
	for i := range r.regions {
		if r.regions[i] == r2 || (r.regions[i].Contains(r2.Begin()) && r.regions[i].Contains(r2.End())) {
			return true
		}
	}
	return false
}
