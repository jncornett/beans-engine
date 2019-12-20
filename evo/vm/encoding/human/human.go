package human

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jncornett/beans-engine/evo/vm"
)

// OpCodes ...
var OpCodes = (func() map[string]vm.OpCode {
	out := make(map[string]vm.OpCode)
	for _, op := range vm.OpCodes {
		out[strings.ToLower(op.String())] = op
	}
	return out
})()

// Decoder ...
type Decoder struct {
	scan *bufio.Scanner
}

// NewDecoder ...
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scan: bufio.NewScanner(r),
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

var bases = map[string]int{
	"x": 16,
	"o": 8,
	"b": 2,
}

func parseValue(s string) (vm.Value, error) {
	base := 10
	for selector, b := range bases {
		test := strings.ToLower(s)
		long := "0" + selector
		if strings.HasPrefix(test, long) {
			s = s[len(long):]
			base = b
			break
		}
		if strings.HasPrefix(test, selector) {
			s = s[len(selector):]
			base = b
			break
		}
	}
	v, err := strconv.ParseInt(s, base, 8)
	if err != nil {
		return 0, err
	}
	return vm.Value(v), nil
}

func parseOp(opName string, fields []string) (vm.Op, error) {
	var arg vm.Value
	op, ok := OpCodes[opName]
	if !ok {
		return vm.Op{}, fmt.Errorf("unknown opcode: %q", opName)
	}
	if len(fields) > 0 {
		var err error
		arg, err = parseValue(fields[0])
		if err != nil {
			return vm.Op{}, fmt.Errorf("could not parse opcode arg: %w", err)
		}
	}
	return vm.Op{
		Type: op,
		Arg:  arg,
	}, nil
}

// DecodeLine ...
func DecodeLine(line string) (op vm.Op, ok bool, err error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, ";") {
		return vm.Op{}, false, nil
	}
	fields := strings.Fields(line)
	op, err = DecodeArgs(fields...)
	if err != nil {
		return vm.Op{}, false, err
	}
	return op, true, nil
}

// DecodeArgs ...
func DecodeArgs(fields ...string) (vm.Op, error) {
	if len(fields) == 0 {
		return vm.Op{}, errors.New("no fields provided")
	}
	op, err := parseOp(fields[0], fields[1:])
	if err != nil {
		return vm.Op{}, err
	}
	return op, nil
}

// Decode ...
func (dec *Decoder) Decode(out *[]vm.Op) error {
	i := 0
	for dec.scan.Scan() {
		i++
		op, ok, err := DecodeLine(dec.scan.Text())
		if err != nil {
			return fmt.Errorf("at %d: %w", i, err)
		}
		if !ok {
			continue
		}
		*out = append(*out, vm.Op(op))
	}
	if err := dec.scan.Err(); err != nil {
		return err
	}
	return nil
}

// Encoder ...
type Encoder struct {
	w io.Writer
}

// NewEncoder ...
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Marshal ...
func Marshal(in []vm.Op) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeLine ...
func EncodeLine(op vm.Op) string {
	return fmt.Sprintf("%s\t%d", strings.ToLower(op.Type.String()), op.Arg)
}

// Encode ...
func (enc *Encoder) Encode(in []vm.Op) error {
	for i, op := range in {
		if _, err := fmt.Fprintln(enc.w, EncodeLine(op)); err != nil {
			return fmt.Errorf("at %d: %w", i+1, err)
		}
	}
	return nil
}
