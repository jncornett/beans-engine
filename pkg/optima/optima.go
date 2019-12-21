package optima

import (
	"errors"
	"sync"

	"github.com/jncornett/beans-engine/pkg/algorithm"
)

// Population ...
type Population interface {
	Len() int
	Cost(i int) float64
	Create(n int)
	Reap(n int)
	Swap(i, j int)
}

// PopulationFuncs ...
type PopulationFuncs struct {
	LenFunc    func() int
	CostFunc   func(i int) float64
	CreateFunc func(n int)
	ReapFunc   func(n int)
	SwapFunc   func(i, j int)
}

var _ Population = (*PopulationFuncs)(nil)

// Len ...
func (pop *PopulationFuncs) Len() int { return pop.LenFunc() }

// Cost ...
func (pop *PopulationFuncs) Cost(i int) float64 { return pop.CostFunc(i) }

// Create ...
func (pop *PopulationFuncs) Create(n int) { pop.CreateFunc(n) }

// Reap ...
func (pop *PopulationFuncs) Reap(n int) { pop.ReapFunc(n) }

// Swap ...
func (pop *PopulationFuncs) Swap(i, j int) { pop.SwapFunc(i, j) }

// Simulation ...
type Simulation struct {
	Size          int
	TargetCost    float64
	MaxIterations int
	ReapRatio     float64
	OnStep        func(minCost float64, steps int)
}

// Step ...
func (sim *Simulation) Step(pop Population) (minCost float64) {
	if pop.Len() > 0 {
		pop.Reap(int(float64(pop.Len()) * sim.ReapRatio))
	}
	pop.Create(sim.Size - pop.Len())
	costs := sim.computeCosts(pop.Len(), pop.Cost)
	if len(costs) == 0 {
		panic(errors.New("population size must not be zero"))
	}
	algorithm.Sort(
		func() int {
			return len(costs)
		},
		func(i, j int) {
			costs[i], costs[j] = costs[j], costs[i]
			pop.Swap(i, j)
		},
		func(i, j int) bool {
			return costs[i] < costs[j]
		},
	)
	return costs[0]
}

// Optimize ...
func (sim *Simulation) Optimize(pop Population) (minCost float64, steps int) {
	for steps = 0; steps < sim.MaxIterations; steps++ {
		minCost = sim.Step(pop)
		if sim.OnStep != nil {
			sim.OnStep(minCost, steps+1)
		}
		if minCost <= sim.TargetCost {
			break
		}
	}
	return minCost, steps
}

func (sim *Simulation) computeCosts(n int, costFn func(int) float64) []float64 {
	type costTuple struct {
		i    int
		cost float64
	}
	tuples := make(chan costTuple)
	go func() {
		defer close(tuples)
		var wg sync.WaitGroup
		defer wg.Wait()
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				tuples <- costTuple{i: i, cost: costFn(i)}
			}(i)
		}
	}()
	m := make(map[int]float64, n)
	for tuple := range tuples {
		m[tuple.i] = tuple.cost
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = m[i]
	}
	return out
}
