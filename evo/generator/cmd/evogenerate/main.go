package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/jncornett/beans-engine/evo/generator"
	"github.com/jncornett/beans-engine/vm"
	"github.com/jncornett/beans-engine/vm/encoding/bytecode"
	"github.com/jncornett/beans-engine/vm/encoding/human"
)

const (
	minLength = 10
	maxLength = 25
)

func main() {
	var (
		length      = flag.Uint("length", 10, "set approximate length")
		variance    = flag.Float64("var", 0.1, "length variance")
		useBytecode = flag.Bool("bytecode", false, "output results in binary")
	)
	flag.Parse()
	rand.Seed(time.Now().Unix())
	n := randomLength(int(*length), *variance)
	var out []vm.Op
	for i := 0; i < n; i++ {
		out = append(out, generator.Default.Sample())
	}
	if *useBytecode {
		bytecode.NewEncoder(os.Stdout).Encode(out)
	} else {
		human.NewEncoder(os.Stdout).Encode(out)
	}
}

func randomLength(center int, variance float64) int {
	lean := float64(center) * variance
	min := float64(center) - lean
	if min < 0 {
		min = 0
	}
	return int(min) + rand.Intn(int(2*lean))
}
