package bytecode

import (
	"bytes"
	"github.com/jncornett/beans-engine/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	code := []vm.Op{
		vm.Op{Type: vm.OpPush, Arg: 42},
		vm.Op{Type: vm.OpPush, Arg: 43},
		vm.Op{Type: vm.OpCompare},
		vm.Op{Type: vm.OpJumpIf},
		vm.Op{Type: vm.OpJumpIf},
		vm.Op{Type: vm.OpNoop},
	}
	b, err := Marshal(code)
	require.NoError(t, err)
	got, err := Unmarshal(b)
	require.NoError(t, err)
	assert.Equal(t, code, got)
}

func TestEncodeDecode(t *testing.T) {
	code1 := []vm.Op{
		vm.Op{Type: vm.OpPush, Arg: 42},
		vm.Op{Type: vm.OpPush, Arg: 43},
		vm.Op{Type: vm.OpCompare},
		vm.Op{Type: vm.OpJumpIf},
		vm.Op{Type: vm.OpJumpIf},
		vm.Op{Type: vm.OpNoop},
	}
	code2 := []vm.Op{
		vm.Op{Type: vm.OpPush, Arg: 44},
		vm.Op{Type: vm.OpPush, Arg: 92},
		vm.Op{Type: vm.OpCompare},
		vm.Op{Type: vm.OpJumpIf},
		vm.Op{Type: vm.OpNoop},
	}
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	err := enc.Encode(code1)
	require.NoError(t, err)
	err = enc.Encode(code2)
	require.NoError(t, err)
	dec := NewDecoder(&buf)
	var gotCode1 []vm.Op
	err = dec.Decode(&gotCode1)
	require.NoError(t, err)
	require.Equal(t, code1, gotCode1)
	var gotCode2 []vm.Op
	err = dec.Decode(&gotCode2)
	require.NoError(t, err)
	require.Equal(t, code2, gotCode2)
}
