package ndjson

import (
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
