package json

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/jncornett/beans-engine/evo/vm"
)

// OpField ...
type OpField struct {
	Type vm.OpCode `json:"o,omitempty"`
	Arg  vm.Value  `json:"a,omitempty"`
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
func Unmarshal(p []byte) ([]vm.Op, error) {
	var out []vm.Op
	if err := NewDecoder(bytes.NewReader(p)).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Decode ...
func (dec *Decoder) Decode(out *[]vm.Op) error {
	return dec.dec.Decode(out)
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
func Marshal(in []vm.Op) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encode ...
func (enc *Encoder) Encode(in []vm.Op) error {
	return enc.enc.Encode(in)
}
