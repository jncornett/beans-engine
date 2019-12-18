package defaultimpl

import (
	"github.com/jncornett/beans-engine/vm"
)

// OpPushConst pushes a constant value onto the stack.
func OpPushConst(ctx vm.Context) {
	ctx.Stack().PushValue(ctx.Instr().Arg)
}

// OpPopValue pops a value from the current frame.
func OpPopValue(ctx vm.Context) {
	ctx.Stack().PopValues(1)
}

// OpCall pushes a new frame onto the stack.
func OpCall(ctx vm.Context) {
	iptr, ok := ctx.Script().FindNextLabel(ctx.Instr().Arg)
	if !ok {
		return
	}
	ctx.Script().Jump(iptr + 1)
}

// OpReturn pops a frame off of the stack and resets the instruction pointer.
func OpReturn(ctx vm.Context) {
	frame, ok := ctx.PopFrame()
	if !ok {
		return // no frames to process
	}
	ctx.Script().Jump(frame.Return)
}

// OpJumpIf ...
func OpJumpIf(ctx vm.Context) {
	val, _ := ctx.Stack().GetValue(-1)
	instr := ctx.Instr()
	if val.Bool() != instr.Arg.Bool() {
		return
	}
	offset := int(instr.Option)
	if offset == 0 {
		offset = 1 // default to skipping the next instruction
	}
	ctx.Script().JumpOffset(offset)
}

// OpCompare ...
func OpCompare(ctx vm.Context) {
	rhs, _ := ctx.PopValue()
	lhs, _ := ctx.PopValue()
	ctx.Stack().PushValue(lhs - rhs)
}

// OpNot ...
func OpNot(ctx vm.Context) {
	val, _ := ctx.PopValue()
	ctx.Stack().PushValue(val.Not())
}

// OpInc ...
func OpInc(ctx vm.Context) {
	val, _ := ctx.PopValue()
	step := ctx.Instr().Arg
	if step == 0 {
		step = 1
	}
	ctx.Stack().PushValue(val + step)
}

// OpDec ...
func OpDec(ctx vm.Context) {
	val, _ := ctx.PopValue()
	step := ctx.Instr().Arg
	if step == 0 {
		step = 1
	}
	ctx.Stack().PushValue(val - step)
}

// OpLoad ...
func OpLoad(ctx vm.Context) {
	instr := ctx.Instr()
	var reg int
	if instr.Option.Bool() {
		// load dynamic
		val, _ := ctx.PopValue()
		reg = int(val)
	} else {
		reg = int(instr.Arg)
	}
	if val, ok := ctx.Registers().Load(reg); ok {
		ctx.Stack().PushValue(val)
	}
}

// OpStore ...
func OpStore(ctx vm.Context) {
	instr := ctx.Instr()
	var reg int
	if instr.Option.Bool() {
		// load dynamic
		val, _ := ctx.PopValue()
		reg = int(val)
	} else {
		reg = int(instr.Arg)
	}
	if val, ok := ctx.PopValue(); ok {
		ctx.Registers().Store(reg, val)
	}
}

// OpLabel ...
func OpLabel(ctx vm.Context) {
	ctx.Stack().Pop(1)
}

// Map ...
var Map = vm.Impl{
	vm.OpPushConst: OpPushConst,
	vm.OpPop:       OpPopValue,
	vm.OpCall:      OpCall,
	vm.OpReturn:    OpReturn,
	vm.OpJumpIf:    OpJumpIf,
	vm.OpCompare:   OpCompare,
	vm.OpNot:       OpNot,
	vm.OpInc:       OpInc,
	vm.OpDec:       OpDec,
	vm.OpLoad:      OpLoad,
	vm.OpStore:     OpStore,
	vm.OpLabel:     OpLabel,
}
