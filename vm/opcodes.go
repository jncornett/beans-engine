package vm

import "strconv"

// OpCode ...
type OpCode int

const (
	// OpNoop ...
	OpNoop OpCode = iota
	// OpPush ...
	OpPush
	// OpPop ...
	OpPop
	// OpCall ...
	OpCall
	// OpReturn ...
	OpReturn
	// OpJumpIf ...
	OpJumpIf
	// OpCompare ...
	OpCompare
	// OpNot ...
	OpNot
	// OpInc ...
	OpInc
	// OpDec ...
	OpDec
	// OpLoad ...
	OpLoad
	// OpStore ...
	OpStore
	// OpLabel ...
	OpLabel
	// OpSyscall ...
	OpSyscall
	// OpMax ...
	OpMax
)

func (op OpCode) String() string {
	switch op {
	case OpNoop:
		return "Noop"
	case OpPush:
		return "Push"
	case OpPop:
		return "Pop"
	case OpCall:
		return "Call"
	case OpReturn:
		return "Return"
	case OpJumpIf:
		return "JumpIf"
	case OpCompare:
		return "Compare"
	case OpNot:
		return "Not"
	case OpInc:
		return "Inc"
	case OpDec:
		return "Dec"
	case OpLoad:
		return "Load"
	case OpStore:
		return "Store"
	case OpLabel:
		return "Label"
	case OpSyscall:
		return "Syscall"
	case OpMax:
		return "Max"
	default:
		return "OpCode(" + strconv.Itoa(int(op)) + ")"
	}
}

// OpCodes ...
var OpCodes = []OpCode{
	OpNoop,
	OpPush,
	OpPop,
	OpCall,
	OpReturn,
	OpJumpIf,
	OpCompare,
	OpNot,
	OpInc,
	OpDec,
	OpLoad,
	OpStore,
	OpLabel,
	OpSyscall,
}
