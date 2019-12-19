package generator

import (
	"github.com/jncornett/beans-engine/pkg/discrete"
	"github.com/jncornett/beans-engine/vm"
)

// OpCodeVar defines a random variable over OpCodes.
type OpCodeVar struct {
	discrete.IntVar
}

// Sample ...
func (ov OpCodeVar) Sample() vm.OpCode {
	i := ov.IntVar.Sample()
	if i < 0 || i >= int64(vm.OpMax) {
		return vm.OpNoop
	}
	return vm.OpCode(i)
}

// ValueVar defines a random variable over Values.
type ValueVar struct {
	discrete.IntVar
}

// Sample ...
func (vv ValueVar) Sample() vm.Value {
	i := vv.IntVar.Sample()
	if i < vm.MinValue || i > vm.MaxValue {
		return 0
	}
	return vm.Value(i)
}

// OpConfig ...
type OpConfig struct {
	Weight float64
	Arg    ValueVar
}

// OpVar defines a random variable over instructions.
type OpVar struct {
	Type OpCodeVar
	Arg  map[vm.OpCode]ValueVar
}

// NewOpVar ...
func NewOpVar(config map[vm.OpCode]OpConfig) OpVar {
	var opCodePMF []discrete.IntVarPoint
	argVarMap := make(map[vm.OpCode]ValueVar)
	for opCode, c := range config {
		opCodePMF = append(opCodePMF, discrete.IntVarPoint{
			X: discrete.Const(int64(opCode)),
			Y: c.Weight,
		})
		argVarMap[opCode] = c.Arg
	}
	return OpVar{
		Type: OpCodeVar{discrete.FromPMF(opCodePMF)},
		Arg:  argVarMap,
	}
}

// Sample ...
func (ov OpVar) Sample() vm.Op {
	op := ov.Type.Sample()
	var val vm.Value
	if vv, ok := ov.Arg[op]; ok {
		val = vv.Sample()
	}
	return vm.Op{Type: op, Arg: val}
}

// Default ...
var Default = NewOpVar(map[vm.OpCode]OpConfig{
	vm.OpNoop:    OpConfig{Weight: 3, Arg: ValueVar{discrete.Const(0)}},
	vm.OpPush:    OpConfig{Weight: 2, Arg: ValueVar{discrete.Range(0, 9)}},
	vm.OpPop:     OpConfig{Weight: 2, Arg: ValueVar{discrete.Const(0)}},
	vm.OpCall:    OpConfig{Weight: 1, Arg: ValueVar{discrete.Range(0, 9)}},
	vm.OpReturn:  OpConfig{Weight: 1, Arg: ValueVar{discrete.Const(0)}},
	vm.OpJumpIf:  OpConfig{Weight: 2, Arg: ValueVar{discrete.Range(-8, 9)}},
	vm.OpCompare: OpConfig{Weight: 1, Arg: ValueVar{discrete.Const(0)}},
	vm.OpNot:     OpConfig{Weight: 1, Arg: ValueVar{discrete.Const(0)}},
	vm.OpInc:     OpConfig{Weight: 2, Arg: ValueVar{discrete.Range(0, 3)}},
	vm.OpDec:     OpConfig{Weight: 2, Arg: ValueVar{discrete.Range(0, 3)}},
	vm.OpLoad:    OpConfig{Weight: 2, Arg: ValueVar{discrete.Range(0, 9)}},
	vm.OpStore:   OpConfig{Weight: 2, Arg: ValueVar{discrete.Range(0, 9)}},
	vm.OpLabel:   OpConfig{Weight: 1, Arg: ValueVar{discrete.Range(0, 9)}},
})
