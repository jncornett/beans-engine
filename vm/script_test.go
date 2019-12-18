package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScript_Peek(t *testing.T) {
	tests := []struct {
		name      string
		script    Script
		wantInstr Instruction
		wantOk    bool
	}{
		{
			name: "empty",
		},
		{
			name: "start",
			script: Script{
				Code: []Instruction{
					Instruction{Op: OpPushConst, Arg: 1, Option: 2},
				},
			},
			wantInstr: Instruction{Op: OpPushConst, Arg: 1, Option: 2},
			wantOk:    true,
		},
		{
			name: "end",
			script: Script{
				Code: []Instruction{
					Instruction{Op: OpPushConst, Arg: 1, Option: 2},
				},
				Iptr: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instr, ok := tt.script.Peek()
			assert.Equal(t, tt.wantInstr, instr)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestScript_Next(t *testing.T) {
	tests := []struct {
		name      string
		script    Script
		wantInstr Instruction
		wantOk    bool
		wantIptr  int
	}{
		{
			name: "empty",
		},
		{
			name: "start",
			script: Script{
				Code: []Instruction{
					Instruction{Op: OpPushConst, Arg: 1, Option: 2},
				},
			},
			wantIptr:  1,
			wantInstr: Instruction{Op: OpPushConst, Arg: 1, Option: 2},
			wantOk:    true,
		},
		{
			name: "end",
			script: Script{
				Code: []Instruction{
					Instruction{Op: OpPushConst, Arg: 1, Option: 2},
				},
				Iptr: 1,
			},
			wantIptr: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instr, ok := tt.script.Next()
			assert.Equal(t, tt.wantInstr, instr)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantIptr, tt.script.Iptr)
		})
	}

}
