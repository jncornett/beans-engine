package vm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackFrame_Push(t *testing.T) {
	tests := []struct {
		stack      []Value
		value      Value
		wantPushed bool
		wantStack  []Value
	}{
		{
			value:      42,
			wantPushed: true,
			wantStack:  []Value{42},
		},
		{
			stack:      []Value{31},
			value:      42,
			wantPushed: true,
			wantStack:  []Value{31, 42},
		},
		{
			stack: (func() []Value {
				values := make([]Value, MaxFrames)
				values[MaxFrames-1] = 31
				return values
			})(),
			value:      42,
			wantPushed: false,
			wantStack: (func() []Value {
				values := make([]Value, MaxFrames)
				values[MaxFrames-1] = 31
				return values
			})(),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v.push(%v)", tt.stack, tt.value), func(t *testing.T) {
			frame := StackFrame{Max: uint(len(tt.stack))}
			copy(frame.Data[:], tt.stack)
			pushed := frame.Push(tt.value)
			assert.Equal(t, tt.wantPushed, pushed)
			assert.Equal(t, tt.wantStack, frame.Values())
		})
	}
}

func TestStackFrame_Pop(t *testing.T) {
	tests := []struct {
		stack      []Value
		n          int
		wantPopped int
		wantStack  []Value
	}{
		{
			n: 0,
		},
		{
			n: 1,
		},
		{
			n: -1,
		},
		{
			stack:      []Value{42},
			n:          1,
			wantPopped: 1,
		},
		{
			stack:      []Value{31, 42},
			n:          1,
			wantPopped: 1,
			wantStack:  []Value{31},
		},
		{
			stack:      []Value{31, 42},
			n:          2,
			wantPopped: 2,
		},
		{
			stack:      []Value{31, 42},
			n:          3,
			wantPopped: 2,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v.pop(%v)", tt.stack, tt.n), func(t *testing.T) {
			frame := StackFrame{Max: uint(len(tt.stack))}
			copy(frame.Data[:], tt.stack)
			popped := frame.Pop(tt.n)
			assert.Equal(t, tt.wantPopped, popped)
			assert.Equal(t, tt.wantStack, frame.Values())
		})
	}
}

func TestStackFrame_Get(t *testing.T) {
	tests := []struct {
		stack   []Value
		idx     int
		wantVal Value
		wantOk  bool
	}{
		{
			idx: 0,
		},
		{
			idx: 1,
		},
		{
			idx: -1,
		},
		{
			stack:   []Value{42},
			idx:     0,
			wantVal: 42,
			wantOk:  true,
		},
		{
			stack: []Value{42},
			idx:   1,
		},
		{
			stack: []Value{42},
			idx:   2,
		},
		{
			stack:   []Value{42},
			idx:     -1,
			wantVal: 42,
			wantOk:  true,
		},
		{
			stack: []Value{42},
			idx:   -2,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v.get(%v)", tt.stack, tt.idx), func(t *testing.T) {
			frame := StackFrame{Max: uint(len(tt.stack))}
			copy(frame.Data[:], tt.stack)
			got, ok := frame.Get(tt.idx)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantVal, got)
		})
	}
}

func TestStack_PushValue(t *testing.T) {
	tests := []struct {
		name       string
		stack      []StackFrame
		v          Value
		wantPushed bool
		wantStack  []StackFrame
	}{
		{
			name:       "empty stack",
			v:          Value(42),
			wantPushed: true,
			wantStack:  []StackFrame{makeStackFrame(0, []Value{42})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := Stack{Max: uint(len(tt.stack))}
			copy(stack.Data[:], tt.stack)
			pushed := stack.PushValue(tt.v)
			assert.Equal(t, tt.wantPushed, pushed)
			assert.Equal(t, tt.wantStack, stack.Frames())
		})
	}
}

func makeStackFrame(ret int, vals []Value) StackFrame {
	frame := StackFrame{Return: ret, Max: uint(len(vals))}
	copy(frame.Data[:], vals)
	return frame
}
