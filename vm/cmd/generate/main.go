package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jncornett/beans-engine/vm"
	"github.com/jncornett/beans-engine/vm/encoding/bytecode"
	"github.com/jncornett/beans-engine/vm/encoding/ndjson"
)

func main() {
	var (
		length = flag.Uint("length", 10, "set approximate length")
		json   = flag.Bool("ndjson", false, "output result in ndjson")
	)
	flag.Parse()
	rand.Seed(time.Now().Unix())
	if err := run(*length, *json); err != nil {
		log.Fatal(err)
	}
}

func run(length uint, json bool) error {
	n := randomCenter(int(length), 0.25)
	var out []vm.Instruction
	for i := 0; i < n; i++ {
		out = append(out, vm.Instruction{
			Op:     randomOp(),
			Arg:    randomValue(1.0),
			Option: randomValue(1.0),
		})
	}
	var marshal func([]vm.Instruction) ([]byte, error)
	if json {
		marshal = ndjson.Marshal
	} else {
		marshal = bytecode.Marshal
	}
	b, err := marshal(out)
	if err != nil {
		return fmt.Errorf("failed to marshal bytecode: %w", err)
	}
	fmt.Print(string(b))
	return nil
}

func randomCenter(val int, factor float64) int {
	x := int(float64(val) * factor)
	lo := val - x
	hi := val + x
	r := hi - lo
	return lo + rand.Intn(r)
}

func randomOp() vm.OpCode {
	i := rand.Intn(len(vm.OpCodes))
	return vm.OpCodes[i]
}

func randomValue(stddev float64) vm.Value {
	return vm.Value(rand.NormFloat64() * stddev)
}
