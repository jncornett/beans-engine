package bytecode

import (
	"bytes"
	"github.com/jncornett/beans-engine/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	code := []vm.Instruction{
		vm.Instruction{Op: vm.OpPushConst, Arg: 42},
		vm.Instruction{Op: vm.OpPushConst, Arg: 43},
		vm.Instruction{Op: vm.OpCompare},
		vm.Instruction{Op: vm.OpJumpIf, Option: 2},
		vm.Instruction{Op: vm.OpJumpIf},
		vm.Instruction{Op: vm.OpNoop},
	}
	b, err := Marshal(code)
	require.NoError(t, err)
	got, err := Unmarshal(b)
	require.NoError(t, err)
	assert.Equal(t, code, got)
}

func TestEncodeDecode(t *testing.T) {
	code1 := []vm.Instruction{
		vm.Instruction{Op: vm.OpPushConst, Arg: 42},
		vm.Instruction{Op: vm.OpPushConst, Arg: 43},
		vm.Instruction{Op: vm.OpCompare},
		vm.Instruction{Op: vm.OpJumpIf, Option: 2},
		vm.Instruction{Op: vm.OpJumpIf},
		vm.Instruction{Op: vm.OpNoop},
	}
	code2 := []vm.Instruction{
		vm.Instruction{Op: vm.OpPushConst, Arg: 44},
		vm.Instruction{Op: vm.OpPushConst, Arg: 92},
		vm.Instruction{Op: vm.OpCompare},
		vm.Instruction{Op: vm.OpJumpIf},
		vm.Instruction{Op: vm.OpNoop},
	}
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	err := enc.Encode(code1)
	require.NoError(t, err)
	err = enc.Encode(code2)
	require.NoError(t, err)
	dec := NewDecoder(&buf)
	var gotCode1 []vm.Instruction
	err = dec.Decode(&gotCode1)
	require.NoError(t, err)
	require.Equal(t, code1, gotCode1)
	var gotCode2 []vm.Instruction
	err = dec.Decode(&gotCode2)
	require.NoError(t, err)
	require.Equal(t, code2, gotCode2)
}