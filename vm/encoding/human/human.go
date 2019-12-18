package human

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jncornett/beans-engine/vm"
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
func Unmarshal(p []byte) ([]vm.Instruction, error) {
	var out []vm.Instruction
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

func parseInstruction(opName string, fields []string) (vm.Instruction, error) {
	var (
		arg    vm.Value
		option vm.Value
	)
	op, ok := OpCodes[opName]
	if !ok {
		return vm.Instruction{}, fmt.Errorf("unknown opcode: %q", opName)
	}
	if len(fields) > 0 {
		var err error
		arg, err = parseValue(fields[0])
		if err != nil {
			return vm.Instruction{}, fmt.Errorf("could not parse opcode arg: %w", err)
		}
		if len(fields) > 1 {
			option, err = parseValue(fields[1])
			if err != nil {
				return vm.Instruction{}, fmt.Errorf("could not parse opcode option: %w", err)
			}
		}
	}
	return vm.Instruction{
		Op:     op,
		Arg:    arg,
		Option: option,
	}, nil
}

// DecodeLine ...
func DecodeLine(line string) (instr vm.Instruction, ok bool, err error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, ";") {
		return vm.Instruction{}, false, nil
	}
	fields := strings.Fields(line)
	instr, err = DecodeArgs(fields...)
	if err != nil {
		return vm.Instruction{}, false, err
	}
	return instr, true, nil
}

// DecodeArgs ...
func DecodeArgs(fields ...string) (vm.Instruction, error) {
	if len(fields) == 0 {
		return vm.Instruction{}, errors.New("no fields provided")
	}
	instr, err := parseInstruction(fields[0], fields[1:])
	if err != nil {
		return vm.Instruction{}, err
	}
	return instr, nil
}

// Decode ...
func (dec *Decoder) Decode(out *[]vm.Instruction) error {
	i := 0
	for dec.scan.Scan() {
		i++
		instr, ok, err := DecodeLine(dec.scan.Text())
		if err != nil {
			return fmt.Errorf("at %d: %w", i, err)
		}
		if !ok {
			continue
		}
		*out = append(*out, vm.Instruction(instr))
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
func Marshal(in []vm.Instruction) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeLine ...
func EncodeLine(instr vm.Instruction) string {
	return fmt.Sprintf("%s\t%d\t%d", strings.ToLower(instr.Op.String()), instr.Arg, instr.Option)
}

// Encode ...
func (enc *Encoder) Encode(in []vm.Instruction) error {
	for i, instr := range in {
		if _, err := fmt.Fprintln(enc.w, EncodeLine(instr)); err != nil {
			return fmt.Errorf("at %d: %w", i+1, err)
		}
	}
	return nil
}
