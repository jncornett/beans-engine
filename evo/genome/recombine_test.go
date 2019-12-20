package genome

import (
	"testing"

	"github.com/jncornett/beans-engine/evo/vm/encoding/evo"
	"github.com/stretchr/testify/assert"
)

const leftCode = `
label 7
push 1
push 2
push 3
pop
jumpif 2
call 7
`

const rightCode = `
push 7
store 4
push 42
store 3
load 4
dec
store 4
load 4
jumpif -3
`

func TestRecombine(t *testing.T) {
	left, _ := evo.Unmarshal([]byte(leftCode))
	right, _ := evo.Unmarshal([]byte(rightCode))
	got := Recombine(DefaultRecombine, left, right)
	assert.Greater(t, len(got), 0)
	assert.NotEqual(t, left, got)
	assert.NotEqual(t, right, got)
}
