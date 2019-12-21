package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jncornett/beans-engine/evo/genome"
	"github.com/jncornett/beans-engine/evo/vm"
	"github.com/jncornett/beans-engine/evo/vm/encoding/evo"
	"github.com/jncornett/beans-engine/evo/vm/impl"
	"github.com/jncornett/beans-engine/pkg/discrete"
	"github.com/jncornett/beans-engine/pkg/optima"
)

const (
	defaultPopulationSize    = 100
	defaultMaxIterations     = 100000
	defaultRuntimeIterations = 100
	defaultReapRatio         = 0.5
	defaultInputSize         = 8
	defaultCodeSize          = 100
)

// Args ...
type Args struct {
	Size    int     `help:"population size"`
	Target  float64 `help:"target cost"`
	Max     int     `help:"max iterations"`
	Timeout int     `help:"vm timeout in steps"`
	Input   int     `help:"number of vm registers"`
}

func main() {
	args := Args{
		Size:    defaultPopulationSize,
		Max:     defaultMaxIterations,
		Input:   defaultInputSize,
		Timeout: defaultRuntimeIterations,
		Target:  1,
	}
	arg.MustParse(&args)
	rand.Seed(time.Now().Unix())
	if err := run(&args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func run(args *Args) error {
	// FIXME do some additional parameter valiation...
	log.Printf("Args: %+v\n", args)
	stepLogDebounce := debounceDuration(2 * time.Second)
	var codes [][]vm.Op
	sim := optima.Simulation{
		Size:          args.Size,
		TargetCost:    args.Target,
		MaxIterations: args.Max,
		ReapRatio:     defaultReapRatio,
		OnStep: func(minCost float64, steps int) {
			stepLogDebounce(func() {
				log.Printf("Step: cost=%v, steps=%v, popSize=%v\n", minCost, steps, len(codes))
			})
		},
	}
	runtime := vm.Runtime{
		Impl:  impl.Map,
		Hooks: vm.RuntimeWithMaxIterations(uint(args.Timeout)),
	}
	costFunc := CostFunc123(func(i int) []int8 {
		registers := make(vm.Register, args.Input)
		state := &vm.State{Registers: registers}
		state.Script.Code = codes[i]
		_ = runtime.Run(state)
		out := make([]int8, 3)
		for i, v := range registers[:3] {
			if i >= len(out) {
				break
			}
			out[i] = int8(v)
		}
		return out
	})
	pop := &optima.PopulationFuncs{
		LenFunc: func() int { return len(codes) },
		CostFunc: func(i int) float64 {
			cost := costFunc(i)
			cost += 0.1 * float64(len(codes[i]))
			return cost
		},
		CreateFunc: func(n int) {
			if len(codes) == 0 {
				// start from beginning
				for i := 0; i < args.Size; i++ {
					codes = append(codes, genome.SampleN(genome.Default, defaultCodeSize))
				}
				return
			}
			bv := discrete.Bernoulli(0.8)
			pv := discrete.Range(0, int64(len(codes)))
			for i := 0; i < n; i++ {
				var code []vm.Op
				if bv.Sample() {
					code = genome.Mutate(genome.DefaultChange, genome.Default, codes[int(pv.Sample())])
				} else {
					left := codes[int(pv.Sample())]
					right := codes[int(pv.Sample())]
					code = genome.Recombine(genome.DefaultRecombine, left, right)
				}
				codes = append(codes, code)
			}
		},
		ReapFunc: func(n int) { codes = codes[:len(codes)-n] },
		SwapFunc: func(i, j int) { codes[i], codes[j] = codes[j], codes[i] },
	}
	cost, steps := sim.Optimize(pop)
	log.Printf("Done: cost=%v, steps=%v\n", cost, steps)
	evo.NewEncoder(os.Stdout).Encode(codes[0])
	return nil
}

func randomRegisters(n int) vm.Register {
	reg := make(vm.Register, n)
	for i := 0; i < n; i++ {
		reg[i] = vm.Value(i)
	}
	rand.Shuffle(n, func(i, j int) {
		reg[i], reg[j] = reg[j], reg[i]
	})
	return reg
}

// SortCost ...
type SortCost struct {
	LenFunc            func() int
	SwapFunc           func(i, j int)
	LessFunc           func(i, j int) bool
	Swaps, Comparisons int
}

func (a *SortCost) Len() int { return a.LenFunc() }
func (a *SortCost) Swap(i, j int) {
	a.Swaps++
	a.SwapFunc(i, j)
}
func (a *SortCost) Less(i, j int) bool {
	a.Comparisons++
	return a.LessFunc(i, j)
}

type syncTime struct {
	t   time.Time
	mux sync.RWMutex
}

func (st *syncTime) Get() time.Time {
	st.mux.RLock()
	defer st.mux.RUnlock()
	return st.t
}

func (st *syncTime) Set(t time.Time) {
	st.mux.Lock()
	defer st.mux.Unlock()
	st.t = t
}

func debounceDuration(d time.Duration) func(func()) {
	var triggered syncTime
	return func(fn func()) {
		now := time.Now()
		if now.Sub(triggered.Get()) >= d {
			triggered.Set(now)
			fn()
		}
	}
}

// CostFunc123 ...
func CostFunc123(computeFn func(int) []int8) func(int) float64 {
	return func(i int) float64 {
		out := computeFn(i)
		for len(out) < 3 {
			out = append(out, math.MaxInt8-1)
		}
		return math.Abs(float64(out[0])-1) + math.Abs(float64(out[1])-2) + math.Abs(float64(out[2]-3))
	}
}

// CostFuncSortedList ...
func CostFuncSortedList(inputSize int, sortFn func(i int, input []int8) []int8) func(int) float64 {
	return func(i int) float64 {
		// Create input
		input := make([]int8, inputSize)
		for i := 0; i < len(input); i++ {
			input[i] = int8(i)
		}
		rand.Shuffle(len(input), func(i, j int) { input[i], input[j] = input[j], input[i] })
		output := sortFn(i, input)
		// // Create registers, using 2nd half as input
		// registers := make(vm.Register, 2*inputSize)
		// copy(registers[inputSize:], input)
		// state := &vm.State{Registers: registers}
		// state.Script.Code = codes[i]
		// _ = runtime.Run(state)
		// output := make([]vm.Value, args.Input)
		// copy(output, state.Registers[args.Input:])
		// compute sort distance
		tmp := make([]int8, len(output))
		copy(tmp, output)
		sortDistance := SortCost{
			LenFunc:  func() int { return len(tmp) },
			SwapFunc: func(i, j int) { tmp[i], tmp[j] = tmp[j], tmp[i] },
			LessFunc: func(i, j int) bool { return tmp[i] < tmp[j] },
		}
		sort.Sort(&sortDistance)
		// report results
		// log.Println("---")
		// log.Println("Input: ", input)
		// log.Println("Output:", output)
		// log.Println("Swaps: ", sortDistance.swaps)
		// log.Println("Comps: ", sortDistance.comparisons)
		return float64(sortDistance.Swaps + sortDistance.Comparisons)
	}
}
