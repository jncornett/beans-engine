package vm

// Script ...
type Script struct {
	Code []Op
	Iptr int
}

// Peek ...
func (script Script) Peek() (instr Op, ok bool) {
	if script.Iptr >= len(script.Code) {
		return Op{}, false
	}
	return script.Code[script.Iptr], true
}

// Next ...
func (script *Script) Next() (instr Op, ok bool) {
	instr, ok = script.Peek()
	if ok {
		script.Iptr++
	}
	return instr, ok
}

// FindNextLabel ...
func (script *Script) FindNextLabel(val Value) (iptr int, ok bool) {
	// FIXME linear runtime
	for i := script.Iptr; i < len(script.Code); i++ {
		instr := script.Code[i]
		if instr.Type != OpLabel {
			continue
		}
		if instr.Arg != val {
			continue
		}
		return i, true
	}
	if script.Iptr < len(script.Code) {
		for i := 0; i < script.Iptr; i++ {
			instr := script.Code[i]
			if instr.Type != OpLabel {
				continue
			}
			if instr.Arg != val {
				continue
			}
			return i, true
		}
	}
	return 0, false
}

// JumpOffset ...
func (script *Script) JumpOffset(offset int) (iptr int) {
	return script.Jump(script.Iptr + offset)
}

// Jump ...
func (script *Script) Jump(to int) (iptr int) {
	script.Iptr = clampIndex(len(script.Code), to)
	return script.Iptr
}

// Done ...
func (script Script) Done() bool {
	return script.Iptr >= len(script.Code)
}

// Reset ...
func (script *Script) Reset() {
	script.Iptr = 0
}
