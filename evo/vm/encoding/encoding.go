package encoding

import "github.com/jncornett/beans-engine/evo/vm"

// Encoder ...
type Encoder interface {
	Encode([]vm.Op) error
}

// Decoder ...
type Decoder interface {
	Decode(*[]vm.Op) error
}
