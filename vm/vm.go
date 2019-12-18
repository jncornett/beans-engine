package vm

const (
	// FrameSize ...
	FrameSize = 8
	// MaxFrames ...
	MaxFrames = 8
)

// Value ...
type Value int8

// Bool ...
func (v Value) Bool() bool {
	return v != 0
}

// Not ...
func (v Value) Not() Value {
	if v == 0 {
		return 1
	}
	return 0
}

// Instruction ...
type Instruction struct {
	Op     OpCode
	Arg    Value
	Option Value
}

// State ...
type State struct {
	Script    Script
	Stack     Stack
	Registers Register
}

// FrameSnapshot ...
type FrameSnapshot struct {
	Return int
	Values []Value
}

// Snapshot ...
type Snapshot struct {
	Iptr      int
	Stack     []FrameSnapshot
	Registers []Value
}

// Snapshot ...
func (state *State) Snapshot() Snapshot {
	registers := make([]Value, len(state.Registers))
	copy(registers, state.Registers)
	frames := state.Stack.Frames()
	stack := make([]FrameSnapshot, 0, len(frames))
	for _, frame := range frames {
		stack = append(stack, FrameSnapshot{
			Return: frame.Return,
			Values: frame.Values(),
		})
	}
	return Snapshot{
		Iptr:      state.Script.Iptr,
		Stack:     stack,
		Registers: registers,
	}
}

// OpImpl ...
type OpImpl func(ctx Context)

// Impl ...
type Impl map[OpCode]OpImpl

func clampIndex(max, i int) int {
	if i < 0 {
		return 0
	}
	if i > max {
		return max
	}
	return i
}

func clampOffset(max, idx int) int {
	i, _ := offsetIndex(max, idx)
	return clampIndex(max, i)
}

func offsetIndex(max int, idx int) (i int, ok bool) {
	if idx < 0 {
		idx = max + idx
	}
	if idx < 0 || idx >= max {
		return 0, false
	}
	return idx, true
}
