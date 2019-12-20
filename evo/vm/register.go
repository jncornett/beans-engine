package vm

// Register ...
type Register []Value

// Load ...
func (reg *Register) Load(idx int) (val Value, ok bool) {
	if idx < 0 || idx >= len(*reg) {
		return 0, false
	}
	return (*reg)[idx], true
}

// Store ...
func (reg *Register) Store(idx int, val Value) (ok bool) {
	if idx < 0 || idx >= len(*reg) {
		return false
	}
	(*reg)[idx] = val
	return true
}
