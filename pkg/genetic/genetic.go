package genetic

import "github.com/jncornett/beans-engine/pkg/algorithm"

// Population ...
type Population interface {
	Size() int
	Cost(index int) float64
	Repopulate()
	Swap(i, j int)
	Drop(n int)
}

// Params ...
type Params struct {
	TargetCost    float64
	MaxIterations int
	Reap          func(sortedCosts []float64) (drop int)
}

// Optimize ...
func Optimize(pop Population, params Params) (cost float64) {
	type costEntry struct {
		i    int
		cost float64
	}
	i := 1
	for {
		costs := make([]float64, 0, pop.Size())
		for i := 0; i < pop.Size(); i++ {
			costs = append(costs, pop.Cost(i))
		}
		algorithm.Sort(
			func() int { return pop.Size() },
			func(i, j int) {
				pop.Swap(i, j)
				costs[i], costs[j] = costs[j], costs[i]
			},
			func(i, j int) bool { return costs[i] < costs[j] },
		)
		if i >= params.MaxIterations || costs[0] <= params.TargetCost {
			return costs[0]
		}
		i++
		pop.Drop(params.Reap(costs))
		pop.Repopulate()
	}
}
