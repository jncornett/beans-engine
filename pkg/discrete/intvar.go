package discrete

import "math/rand"

import "sort"

// IntVar represents a discrete random variable.
type IntVar interface {
	Sample() int64
}

// IntVarFunc is a func that satisfies IntVar.
type IntVarFunc func() int64

// Sample ...
func (fn IntVarFunc) Sample() int64 {
	return fn()
}

// IntPoint represents a data point that maps an int value to a float.
type IntPoint struct {
	X int
	Y float64
}

// ConstIntVar is a random variable that represents a constant value.
type ConstIntVar int64

// Const ...
func Const(i int64) ConstIntVar {
	return ConstIntVar(i)
}

// Sample ...
func (iv ConstIntVar) Sample() int64 {
	return int64(iv)
}

// IntVarPoint represents a data point that maps a distribution to a float.
type IntVarPoint struct {
	X IntVar
	Y float64
}

// PiecewiseIntVar is a random variable that is defined by a piecwise CDF.
type PiecewiseIntVar []IntVarPoint

// FromPMF ...
func FromPMF(pmf []IntVarPoint) PiecewiseIntVar {
	var cdf []IntVarPoint
	total := 0.0
	for _, p := range pmf {
		if p.Y == 0 {
			continue
		}
		total += p.Y
		cdf = append(cdf, IntVarPoint{
			X: p.X,
			Y: total,
		})
	}
	// normalize cdf ranges
	for i, p := range cdf {
		cdf[i].Y = p.Y / total
	}
	return PiecewiseIntVar(cdf)
}

// Sample ...
func (iv PiecewiseIntVar) Sample() int64 {
	if len(iv) == 0 {
		return 0
	}
	f := rand.Float64()
	i := sort.Search(len(iv), func(i int) bool {
		return iv[i].Y >= f
	})
	if i < 0 {
		i = 0
	} else if i >= len(iv) {
		i = len(iv) - 1
	}
	return iv[i].X.Sample()
}

// RangeIntVar is a uniform random variable that is defined by a range.
type RangeIntVar struct{ Min, Max int64 }

// Range ...
func Range(min, max int64) RangeIntVar {
	return RangeIntVar{Min: min, Max: max}
}

// Sample ...
func (iv RangeIntVar) Sample() int64 {
	r := iv.Max - iv.Min
	if r < 1 {
		return 0
	}
	return iv.Min + rand.Int63n(r)
}

// SampleK ...
func SampleK(k int, iv IntVar) []int64 {
	var out []int64
	for i := 0; i < k; i++ {
		out = append(out, iv.Sample())
	}
	return out
}
