package genome

import (
	"github.com/jncornett/beans-engine/evo/vm"
	"github.com/jncornett/beans-engine/pkg/discrete"
)

// RecombineVar ...
type RecombineVar struct {
	Length discrete.IntVar
	Switch discrete.BoolVar
}

// RecombineEntry ...
type RecombineEntry struct {
	Length int
	Switch bool
}

// Sample ...
func (rv RecombineVar) Sample() RecombineEntry {
	return RecombineEntry{
		Length: int(rv.Length.Sample()),
		Switch: rv.Switch.Sample(),
	}
}

// DefaultRecombine ...
var DefaultRecombine = RecombineVar{
	Length: discrete.Range(3, 4),
	Switch: discrete.Bernoulli(0.5),
}

// Recombine ...
func Recombine(rv RecombineVar, left, right []vm.Op) []vm.Op {
	var out []vm.Op
	for {
		switch {
		case len(left) == 0:
			out = append(out, right...)
			return out
		case len(right) == 0:
			out = append(out, left...)
			return out
		}
		re := rv.Sample()
		if re.Switch {
			n := re.Length
			if n > len(right) {
				n = len(right)
			}
			out = append(out, right[:n]...)
			right = right[n:]
		} else {
			n := re.Length
			if n > len(left) {
				n = len(left)
			}
			out = append(out, left[:n]...)
			left = left[n:]
		}
	}
}
