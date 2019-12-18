package vm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		register Register
		idx      int
		wantOk   bool
		wantVal  Value
	}{
		{
			register: Register{1, 2, 3},
			idx:      0,
			wantOk:   true,
			wantVal:  1,
		},
		{
			register: Register{1, 2, 3},
			idx:      -1,
		},
		{
			register: Register{1, 2, 3},
			idx:      2,
			wantOk:   true,
			wantVal:  3,
		},
		{
			register: Register{1, 2, 3},
			idx:      3,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v.Load(%v)", tt.register, tt.idx), func(t *testing.T) {
			got, ok := tt.register.Load(tt.idx)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantVal, got)
		})
	}
}

func TestStore(t *testing.T) {
	tests := []struct {
		register     Register
		idx          int
		val          Value
		wantOk       bool
		wantRegister Register
	}{
		{
			register:     Register{},
			idx:          0,
			val:          7,
			wantOk:       false,
			wantRegister: Register{},
		},
		{
			register:     Register{},
			idx:          -1,
			val:          7,
			wantOk:       false,
			wantRegister: Register{},
		},
		{
			register:     Register{},
			idx:          1,
			val:          7,
			wantOk:       false,
			wantRegister: Register{},
		},
		{
			register:     Register{1, 2, 3},
			idx:          0,
			val:          7,
			wantOk:       true,
			wantRegister: Register{7, 2, 3},
		},
		{
			register:     Register{1, 2, 3},
			idx:          -1,
			val:          7,
			wantRegister: Register{1, 2, 3},
		},
		{
			register:     Register{1, 2, 3},
			idx:          1,
			val:          7,
			wantOk:       true,
			wantRegister: Register{1, 7, 3},
		},
		{
			register:     Register{1, 2, 3},
			idx:          3,
			val:          7,
			wantRegister: Register{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v.Load(%v)", tt.register, tt.idx), func(t *testing.T) {
			ok := tt.register.Store(tt.idx, tt.val)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantRegister, tt.register)
		})
	}
}
