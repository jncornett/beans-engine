package vm

// RuntimeHook ...
type RuntimeHook int

const (
	// RuntimeHookBeforeStep ...
	RuntimeHookBeforeStep RuntimeHook = iota
)

// RuntimeHandler ...
type RuntimeHandler func(*Runtime, *State, *RunResult) (ok bool)

// RuntimeHookConfig ...
type RuntimeHookConfig map[RuntimeHook][]RuntimeHandler

// Runtime ...
type Runtime struct {
	FrameSize int
	Impl      Impl
	Hooks     map[RuntimeHook][]RuntimeHandler
}

// RunResult ...
type RunResult struct {
	Interrupted bool
	Iterations  int
}

type runtimeContext struct {
	runtime *Runtime
	state   *State
	op      Op
}

func (ctx runtimeContext) Instr() Op {
	return ctx.Op()
}

func (ctx runtimeContext) Op() Op {
	return ctx.op
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
	Instr() Op
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
		if !r.hook(RuntimeHookBeforeStep, state, &result) {
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
func (r *Runtime) Exec(state *State, instr Op) {
	ctx := runtimeContext{
		runtime: r,
		state:   state,
		op:      instr,
	}
	if fn, ok := r.Impl[instr.Type]; ok {
		fn(&ctx)
	}
}

func (r *Runtime) hook(rh RuntimeHook, state *State, result *RunResult) (ok bool) {
	ok = true
	for _, h := range r.Hooks[rh] {
		ok = h(r, state, result)
		if !ok {
			break
		}
	}
	return ok
}

// AddHookFunc ...
func (r *Runtime) AddHookFunc(rh RuntimeHook, h RuntimeHandler) *Runtime {
	r.Hooks[rh] = append(r.Hooks[rh], h)
	return r
}

// AddHook ...
func (r *Runtime) AddHook(cfg RuntimeHookConfig) *Runtime {
	for rh, handlers := range cfg {
		r.Hooks[rh] = append(r.Hooks[rh], handlers...)
	}
	return r
}

// RemoveHooks ...
func (r *Runtime) RemoveHooks(rh RuntimeHook) *Runtime {
	delete(r.Hooks, rh)
	return r
}

// RuntimeWithMaxIterations ...
func RuntimeWithMaxIterations(max uint) RuntimeHookConfig {
	return RuntimeHookConfig{
		RuntimeHookBeforeStep: []RuntimeHandler{
			func(r *Runtime, state *State, result *RunResult) (ok bool) {
				return result.Iterations < int(max)
			},
		},
	}
}
