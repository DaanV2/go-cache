package hash

import "math"

type Range struct {
	min uint64
	max uint64
}

func NewRange() Range {
	return Range{
		min: math.MaxUint64,
		max: 0,
	}
}

func (r *Range) Update(value uint64) {
	r.min = min(r.min, value)
	r.max = max(r.max, value)
}

func (r *Range) Has(value uint64) bool {
	return r.min <= value && value <= r.max
}
