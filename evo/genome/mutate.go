package genome

import (
	"github.com/jncornett/beans-engine/evo/vm"
	"github.com/jncornett/beans-engine/pkg/discrete"
)

// Change ...
type Change int

const (
	// ChangeNone ...
	ChangeNone Change = iota
	// ChangeInsert ...
	ChangeInsert
	// ChangeDelete ...
	ChangeDelete
	// ChangeReplace ...
	ChangeReplace
	// ChangeMax ...
	ChangeMax
)

// ChangeVar ...
type ChangeVar struct {
	discrete.IntVar
}

// Sample ...
func (cv ChangeVar) Sample() Change {
	return Change(cv.IntVar.Sample())
}

// DefaultChange ...
var DefaultChange = ChangeVar{
	IntVar: discrete.FromPMF([]discrete.IntVarPoint{
		discrete.IntVarPoint{Y: 0.9, X: discrete.Const(int64(ChangeNone))},
		discrete.IntVarPoint{Y: 1.0, X: discrete.Range(int64(ChangeInsert), int64(ChangeMax))},
	}),
}

// Mutate ...
func Mutate(cv ChangeVar, ov OpVar, code []vm.Op) []vm.Op {
	var out []vm.Op
	for _, op := range code {
		switch cv.Sample() {
		case ChangeInsert:
			out = append(out, op, ov.Sample())
		case ChangeDelete:
			continue
		case ChangeReplace:
			out = append(out, ov.Sample())
		default:
			out = append(out, op)
		}
	}
	return out
}
