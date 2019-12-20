package discrete

import "math/rand"

// BoolVar ...
type BoolVar interface {
	Sample() bool
}

// BoolVarFunc ...
type BoolVarFunc func() bool

// Sample ...
func (fn BoolVarFunc) Sample() bool {
	return fn()
}

// BernoulliBoolVar ...
type BernoulliBoolVar float64

// Bernoulli ...
func Bernoulli(p float64) BernoulliBoolVar {
	if p < 0 || p >= 1.0 {
		panic("p must be in the range [0.0, 1.0)")
	}
	return BernoulliBoolVar(p)
}

// Sample ...
func (bv BernoulliBoolVar) Sample() bool {
	return rand.Float64() > float64(bv)
}
