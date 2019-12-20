// Package cli contains utilities for building CLIs for evo.
package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jncornett/beans-engine/evo/vm"
	"github.com/jncornett/beans-engine/evo/vm/encoding"
	"github.com/jncornett/beans-engine/evo/vm/encoding/evo"
	"github.com/jncornett/beans-engine/evo/vm/encoding/evox"
	"github.com/jncornett/beans-engine/evo/vm/encoding/json"
)

// Encoding ...
type Encoding string

const (
	// EncodingNone ...
	EncodingNone Encoding = ""
	// EncodingJSON ...
	EncodingJSON Encoding = "json"
	// EncodingEvo ...
	EncodingEvo Encoding = "evo"
	// EncodingEvoX ...
	EncodingEvoX Encoding = "evox"
)

// StdioFilename ...
const StdioFilename = "-"

// AvailableEncodings ...
var AvailableEncodings = []Encoding{
	EncodingJSON,
	EncodingEvo,
	EncodingEvoX,
}

// Extensions ...
var Extensions = map[string]Encoding{
	".json": EncodingJSON,
	".evo":  EncodingEvo,
	".evox": EncodingEvoX,
}

// NewDecoder ...
func NewDecoder(e Encoding, r io.Reader) (encoding.Decoder, error) {
	switch e {
	case EncodingJSON:
		return json.NewDecoder(r), nil
	case EncodingEvo:
		return evo.NewDecoder(r), nil
	case EncodingEvoX:
		return evox.NewDecoder(r), nil
	default:
		return nil, fmt.Errorf("unknown encoding: %q", string(e))
	}
}

// NewUnmarshaler ...
func NewUnmarshaler(e Encoding) (func([]byte) ([]vm.Op, error), error) {
	switch e {
	case EncodingJSON:
		return json.Unmarshal, nil
	case EncodingEvo:
		return evo.Unmarshal, nil
	case EncodingEvoX:
		return evox.Unmarshal, nil
	default:
		return nil, fmt.Errorf("unknown encoding: %q", string(e))
	}
}

// Unmarshal ...
func Unmarshal(e Encoding, p []byte) ([]vm.Op, error) {
	fn, err := NewUnmarshaler(e)
	if err != nil {
		return nil, err
	}
	return fn(p)
}

// NewEncoder ...
func NewEncoder(e Encoding, w io.Writer) (encoding.Encoder, error) {
	switch e {
	case EncodingJSON:
		return json.NewEncoder(w), nil
	case EncodingEvo:
		return evo.NewEncoder(w), nil
	case EncodingEvoX:
		return evox.NewEncoder(w), nil
	default:
		return nil, fmt.Errorf("unknown encoding: %q", string(e))
	}
}

// NewMarshaler ...
func NewMarshaler(e Encoding) (func([]vm.Op) ([]byte, error), error) {
	switch e {
	case EncodingJSON:
		return json.Marshal, nil
	case EncodingEvo:
		return evo.Marshal, nil
	case EncodingEvoX:
		return evox.Marshal, nil
	default:
		return nil, fmt.Errorf("unknown encoding: %q", string(e))
	}
}

// Marshal ...
func Marshal(e Encoding, code []vm.Op) ([]byte, error) {
	fn, err := NewMarshaler(e)
	if err != nil {
		return nil, err
	}
	return fn(code)
}

// Load ...
func Load(filename string, e Encoding) ([]vm.Op, error) {
	if e == EncodingNone {
		e, _ = Extensions[filepath.Ext(filename)]
	}
	if e == EncodingNone {
		return nil, fmt.Errorf("could not determine encoding for %q", filename)
	}
	var r io.Reader
	if filename == StdioFilename {
		r = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	dec, err := NewDecoder(e, r)
	if err != nil {
		return nil, err
	}
	var out []vm.Op
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Save ...
func Save(filename string, e Encoding, code []vm.Op) error {
	log.Println("Save", filename, e, code)
	if e == EncodingNone {
		e, _ = Extensions[filepath.Ext(filename)]
	}
	if e == EncodingNone {
		return fmt.Errorf("could not determine encoding for %q", filename)
	}
	var w io.Writer
	if filename == StdioFilename {
		w = os.Stdout
	} else {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	enc, err := NewEncoder(e, w)
	if err != nil {
		return err
	}
	return enc.Encode(code)
}

// GuessEncoding ...
func GuessEncoding(filename string) Encoding {
	return Extensions[filepath.Ext(filename)]
}
