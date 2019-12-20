package impl

import "github.com/jncornett/beans-engine/evo/vm"

// OpPush pushes a constant value onto the stack.
func OpPush(ctx vm.Context) {
	ctx.Stack().PushValue(ctx.Instr().Arg)
}

// OpPop pops a value from the current frame.
func OpPop(ctx vm.Context) {
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
	val, ok := ctx.Stack().GetValue(-1)
	if ok {
		ctx.Stack().PopValues(1)
	}
	op := ctx.Instr()
	if !val.Bool() {
		return
	}
	offset := int(op.Arg)
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
	op := ctx.Instr()
	val, ok := ctx.Registers().Load(int(op.Arg))
	if !ok {
		// load dynamic
		i, ok := ctx.Stack().GetValue(-1)
		if ok {
			ctx.Stack().PopValues(1)
			val, ok = ctx.Registers().Load(int(i))
		}
	}
	ctx.Stack().PushValue(val)
}

// OpStore ...
func OpStore(ctx vm.Context) {
	op := ctx.Instr()
	val, ok := ctx.Stack().GetValue(-1)
	if !ok {
		return
	}
	ok = ctx.Registers().Store(int(op.Arg), val)
	if !ok {
		// store dynamic
		i, ok := ctx.Stack().GetValue(-1)
		if ok {
			ctx.Stack().PopValues(1)
			ctx.Registers().Store(int(i), val)
		}
	}
}

// OpLabel ...
func OpLabel(ctx vm.Context) {
	ctx.Stack().Pop(1)
}

// Map ...
var Map = vm.Impl{
	vm.OpPush:    OpPush,
	vm.OpPop:     OpPop,
	vm.OpCall:    OpCall,
	vm.OpReturn:  OpReturn,
	vm.OpJumpIf:  OpJumpIf,
	vm.OpCompare: OpCompare,
	vm.OpNot:     OpNot,
	vm.OpInc:     OpInc,
	vm.OpDec:     OpDec,
	vm.OpLoad:    OpLoad,
	vm.OpStore:   OpStore,
	vm.OpLabel:   OpLabel,
}
