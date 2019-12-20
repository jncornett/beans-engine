package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScript_Peek(t *testing.T) {
	tests := []struct {
		name      string
		script    Script
		wantInstr Op
		wantOk    bool
	}{
		{
			name: "empty",
		},
		{
			name: "start",
			script: Script{
				Code: []Op{
					Op{Type: OpPush, Arg: 1},
				},
			},
			wantInstr: Op{Type: OpPush, Arg: 1},
			wantOk:    true,
		},
		{
			name: "end",
			script: Script{
				Code: []Op{
					Op{Type: OpPush, Arg: 1},
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
		wantInstr Op
		wantOk    bool
		wantIptr  int
	}{
		{
			name: "empty",
		},
		{
			name: "start",
			script: Script{
				Code: []Op{
					Op{Type: OpPush, Arg: 1},
				},
			},
			wantIptr:  1,
			wantInstr: Op{Type: OpPush, Arg: 1},
			wantOk:    true,
		},
		{
			name: "end",
			script: Script{
				Code: []Op{
					Op{Type: OpPush, Arg: 1},
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
