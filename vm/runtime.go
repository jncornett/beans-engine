package vm

// Runtime ...
type Runtime struct {
	MaxIterations int
	FrameSize     int
	Impl          Impl
}

// RunResult ...
type RunResult struct {
	Interrupted bool
	Iterations  int
}

type runtimeContext struct {
	runtime     *Runtime
	state       *State
	instruction Instruction
}

func (ctx runtimeContext) Instr() Instruction {
	return ctx.instruction
}

func (ctx runtimeContext) Stack() *Stack {
	return &ctx.state.Stack
}

func (ctx runtimeContext) Script() *Script {
	return &ctx.state.Script
}

func (ctx runtimeContext) PopFrame() (*StackFrame, bool) {
	frame, ok := ctx.state.Stack.Get(-1)
	if !ok {
		return nil, false
	}
	if ctx.state.Stack.Pop(1) < 1 {
		return nil, false
	}
	return frame, true
}

func (ctx runtimeContext) PopValue() (Value, bool) {
	frame, ok := ctx.state.Stack.Get(-1)
	if !ok {
		return 0, false
	}
	val, ok := frame.Get(-1)
	if !ok {
		return 0, false
	}
	if frame.Pop(1) < 1 {
		return 0, false
	}
	return val, true
}

func (ctx runtimeContext) Registers() *Register {
	return &ctx.state.Registers
}

// Context ...
type Context interface {
	Instr() Instruction
	Stack() *Stack
	PopFrame() (*StackFrame, bool)
	PopValue() (Value, bool)
	Script() *Script
	Registers() *Register
}

// Run ...
func (r *Runtime) Run(state *State) (result RunResult) {
	for {
		// Halt check
		if r.MaxIterations > 0 && (result.Iterations >= r.MaxIterations) {
			result.Interrupted = true
			break
		}
		result.Iterations++
		if !r.Step(state) {
			break
		}
	}
	return result
}

// Step executes a single instruction in the state.
// Step returns true if the program is not halted.
func (r *Runtime) Step(state *State) (ok bool) {
	next, ok := state.Script.Next()
	if !ok {
		return false
	}
	r.Exec(state, next)
	return true
}

// Exec executes an arbitrary instruction against state.
func (r *Runtime) Exec(state *State, instr Instruction) {
	ctx := runtimeContext{
		runtime:     r,
		state:       state,
		instruction: instr,
	}
	if fn, ok := r.Impl[instr.Op]; ok {
		fn(&ctx)
	}
}
