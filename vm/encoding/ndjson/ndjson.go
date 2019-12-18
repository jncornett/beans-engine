package ndjson

import (
	"bytes"
	"encoding/json"
	"github.com/jncornett/beans-engine/vm"
	"io"
)

// InstructionField ...
type InstructionField struct {
	Op     vm.OpCode `json:"o,omitempty"`
	Arg    vm.Value  `json:"a,omitempty"`
	Option vm.Value  `json:"p,omitempty"`
}

// Decoder ...
type Decoder struct {
	dec *json.Decoder
}

// NewDecoder ...
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		dec: json.NewDecoder(r),
	}
}

// Unmarshal ...
func Unmarshal(p []byte) ([]vm.Instruction, error) {
	var out []vm.Instruction
	if err := NewDecoder(bytes.NewReader(p)).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Decode ...
func (dec *Decoder) Decode(out *[]vm.Instruction) error {
	for {
		var instr InstructionField
		if err := dec.dec.Decode(&instr); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		*out = append(*out, vm.Instruction(instr))
	}
}

// Encoder ...
type Encoder struct {
	enc *json.Encoder
}

// NewEncoder ...
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		enc: json.NewEncoder(w),
	}
}

// Marshal ...
func Marshal(in []vm.Instruction) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encode ...
func (enc *Encoder) Encode(in []vm.Instruction) error {
	for _, instr := range in {
		if err := enc.enc.Encode((*InstructionField)(&instr)); err != nil {
			return err
		}
	}
	return nil
}
