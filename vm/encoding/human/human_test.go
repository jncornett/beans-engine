package human

import (
	"testing"

	"github.com/jncornett/beans-engine/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
