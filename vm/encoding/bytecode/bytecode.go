package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jncornett/beans-engine/vm"
	"io"
)

// MagicField ...
type MagicField [4]byte

// Magic ...
var Magic = MagicField{4, 3, 2, 1}

// VersionField ...
type VersionField [4]byte

func (f VersionField) String() string {
	return fmt.Sprintf("%v.%v.%v.%v",
		int(f[0]),
		int(f[1]),
		int(f[2]),
		int(f[3]),
	)
}

// Version ...
var Version = VersionField{0, 0, 1, 0}

// LengthField ...
type LengthField uint64

// OpField ...
type OpField struct {
	Type int8
	Arg  int8
}

// DefaultByteOrder ...
var DefaultByteOrder = binary.LittleEndian

// Decoder ...
type Decoder struct {
	r         io.Reader
	ByteOrder binary.ByteOrder
}

// NewDecoder ...
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r, ByteOrder: DefaultByteOrder}
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
	var magic MagicField
	if err := dec.read(&magic); err != nil {
		return fmt.Errorf("invalid magic field: %w", err)
	}
	if !bytes.Equal(Magic[:], magic[:]) {
		return fmt.Errorf("wrong magic detected: want %v, got %v", Magic, magic)
	}
	var version VersionField
	if err := dec.read(&version); err != nil {
		return fmt.Errorf("invalid version field: %w", err)
	}
	if !bytes.Equal(Version[:], version[:]) {
		return fmt.Errorf("version mismatch: want %v, got %v", Version, version)
	}
	var length LengthField
	if err := dec.read(&length); err != nil {
		return fmt.Errorf("invalid length field: %w", err)
	}
	for i := 0; LengthField(i) < length; i++ {
		var op OpField
		if err := dec.read(&op); err != nil {
			return fmt.Errorf("at op %d: %w", i, err)
		}
		*out = append(*out, vm.Op{
			Type: vm.OpCode(op.Type),
			Arg:  vm.Value(op.Arg),
		})
	}
	return nil
}

func (dec *Decoder) read(out interface{}) error {
	return binary.Read(dec.r, dec.ByteOrder, out)
}

// Encoder ...
type Encoder struct {
	w         io.Writer
	ByteOrder binary.ByteOrder
}

// NewEncoder ...
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, ByteOrder: DefaultByteOrder}
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
	if err := enc.write(&Magic); err != nil {
		return fmt.Errorf("failed to write magic field: %w", err)
	}
	if err := enc.write(&Version); err != nil {
		return fmt.Errorf("failed to write version field: %w", err)
	}
	length := LengthField(len(in))
	if err := enc.write(&length); err != nil {
		return fmt.Errorf("failed to write length field: %w", err)
	}
	for i, op := range in {
		field := OpField{
			Type: int8(op.Type),
			Arg:  int8(op.Arg),
		}
		if err := enc.write(&field); err != nil {
			return fmt.Errorf("at op %d: %w", i, err)
		}
	}
	return nil
}

func (enc *Encoder) write(in interface{}) error {
	return binary.Write(enc.w, enc.ByteOrder, in)
}
