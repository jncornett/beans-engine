package algorithm

import "sort"

// Sorter ...
type Sorter struct {
	SortLen  func() int
	SortSwap func(i, j int)
	SortLess func(i, j int) bool
}

func (s Sorter) Len() int           { return s.SortLen() }
func (s Sorter) Swap(i, j int)      { s.SortSwap(i, j) }
func (s Sorter) Less(i, j int) bool { return s.SortLess(i, j) }

// Sort ...
func Sort(
	sortLen func() int,
	sortSwap func(i, j int),
	sortLess func(i, j int) bool,
) {
	sort.Sort(&Sorter{
		SortLen:  sortLen,
		SortSwap: sortSwap,
		SortLess: sortLess,
	})
}
