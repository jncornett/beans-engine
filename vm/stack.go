package vm

// StackFrameData ...
type StackFrameData [FrameSize]Value

// StackFrame ...
type StackFrame struct {
	Return int
	Data   StackFrameData
	Max    uint
}

// Push ...
func (frame *StackFrame) Push(val Value) (pushed bool) {
	if frame.Max >= uint(len(frame.Data)) {
		// stack is full
		return false
	}
	frame.Data[frame.Max] = val
	frame.Max++
	return true
}

// Get ...
func (frame *StackFrame) Get(idx int) (val Value, ok bool) {
	i, ok := offsetIndex(int(frame.Max), idx)
	if !ok {
		return 0, false
	}
	return frame.Data[i], true
}

// Pop ...
func (frame *StackFrame) Pop(n int) (popped int) {
	if n < 0 {
		return 0
	}
	pop := uint(n)
	if pop > frame.Max {
		pop = frame.Max
	}
	frame.Max -= pop
	return int(pop)
}

// Values ...
func (frame *StackFrame) Values() []Value {
	max := frame.Max
	if max > uint(len(frame.Data)) {
		max = uint(len(frame.Data))
	}
	if max == 0 {
		return nil // makes testing easier
	}
	return frame.Data[:max]
}

// Stack ...
type Stack struct {
	Data [MaxFrames]StackFrame
	Max  uint
}

// Push ...
func (stack *Stack) Push(iptr int) (pushed bool) {
	if stack.Max >= uint(len(stack.Data)) {
		return false
	}
	stack.Data[stack.Max] = StackFrame{Return: iptr}
	stack.Max++
	return true
}

// PushValue ...
func (stack *Stack) PushValue(val Value) (pushed bool) {
	frame, ok := stack.Get(-1)
	if !ok {
		stack.Push(0)
		frame, ok = stack.Get(-1)
		if !ok {
			return false
		}
	}
	return frame.Push(val)
}

// PopValues ...
func (stack *Stack) PopValues(n int) (popped int) {
	frame, ok := stack.Get(-1)
	if !ok {
		return 0
	}
	return frame.Pop(n)
}

// GetValue ...
func (stack *Stack) GetValue(idx int) (val Value, ok bool) {
	frame, ok := stack.Get(-1)
	if !ok {
		return 0, false
	}
	return frame.Get(idx)
}

// Get ...
func (stack *Stack) Get(idx int) (frame *StackFrame, ok bool) {
	i, ok := offsetIndex(int(stack.Max), idx)
	if !ok {
		return nil, false
	}
	return &stack.Data[i], true
}

// Pop ...
func (stack *Stack) Pop(n int) (popped int) {
	if n < 0 {
		return 0
	}
	pop := uint(n)
	if pop > stack.Max {
		pop = stack.Max
	}
	stack.Max -= pop
	return int(pop)
}

// Frames ...
func (stack *Stack) Frames() []StackFrame {
	max := stack.Max
	if max > uint(len(stack.Data)) {
		max = uint(len(stack.Data))
	}
	if max == 0 {
		return nil // makes testing easier
	}
	return stack.Data[:max]
}
