package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"

	"github.com/c-bata/go-prompt"
	"github.com/jncornett/beans-engine/pkg/skua"
	"github.com/jncornett/beans-engine/vm"
	"github.com/jncornett/beans-engine/vm/impl/defaultimpl"
	"github.com/jncornett/beans-engine/vm/encoding/bytecode"
	"github.com/jncornett/beans-engine/vm/encoding/human"
	"github.com/naoina/toml"
)

// NumRegisters ...
const NumRegisters = 8

func main() {
	var (
		forceRepl = flag.Bool("repl", false, "interactive mode")
	)
	flag.Parse()
	if err := run(*forceRepl, flag.Args()); err != nil {
		log.Fatal(err)
	}
}

func run(forceRepl bool, args []string) error {
	var (
		state = vm.State{
			Registers: make(vm.Register, NumRegisters),
		}
		runtime = vm.Runtime{
			Impl: defaultimpl.Map,
		}
	)
	var fileLoaded bool
	if args := flag.Args(); len(args) > 0 {
		fileLoaded = true
		code, err := load(args[0])
		if err != nil {
			return err
		}
		state.Script = vm.Script{Code: code}
	}
	if !forceRepl && fileLoaded {
		runtime.Run(&state)
		snap := state.Snapshot()
		b, err := toml.Marshal(&snap)
		if err != nil {
			return err
		}
		_, err = fmt.Print(string(b))
		return err
	}
	repl := newRepl(&state, &runtime)
	repl.Loop()
	return nil
}

func newRepl(state *vm.State, runtime *vm.Runtime) *skua.Repl {
	const contextLines = 2
	printList := func(start, iptr int, code []vm.Instruction) error {
		for i, instr := range code {
			s := " "
			if i+start == iptr {
				s = ">"
			}
			line := human.EncodeLine(instr)
			fmt.Printf(
				"%s %s %s\n",
				aurora.Red(s),
				aurora.Gray(12, fmt.Sprintf("%2d:", i+1)),
				line,
			)
		}
		return nil
	}
	opCodeSuggestions := make([]prompt.Suggest, 0, len(vm.OpCodes))
	for name, op := range human.OpCodes {
		opCodeSuggestions = append(opCodeSuggestions, prompt.Suggest{
			Text:        strings.ToLower(name),
			Description: op.String(),
		})
	}
	opCodeCompleter := func(d prompt.Document) []prompt.Suggest {
		line := d.CurrentLine()
		fields := strings.Fields(line)
		if len(fields) == 0 || (len(fields) == 1 && strings.HasSuffix(line, " ")) || len(fields) > 1 {
			return nil
		}
		return prompt.FilterHasPrefix(opCodeSuggestions, d.GetWordBeforeCursor(), false)
	}
	return &skua.Repl{
		Commands: map[string]skua.Command{
			"script": skua.Command{
				Description: "run code",
				Subcommands: map[string]skua.Command{
					"list": skua.Command{
						Description: "print script",
						Subcommands: map[string]skua.Command{
							"all": skua.Command{
								Description: "print entire script",
								Run:         func([]string) error { return printList(0, state.Script.Iptr, state.Script.Code) },
							},
						},
						Run: func([]string) error {
							start := state.Script.Iptr - contextLines
							if start < 0 {
								start = 0
							}
							end := state.Script.Iptr + contextLines + 1
							if end > len(state.Script.Code) {
								end = len(state.Script.Code)
							}
							return printList(start, state.Script.Iptr, state.Script.Code[start:end])
						},
					},
					"new": skua.Command{
						Description: "write new script",
						Run: func([]string) error {
							var code []vm.Instruction
						Loop:
							for {
								line := strings.TrimSpace(prompt.Input(fmt.Sprintf(" %2d: ", len(code)+1), opCodeCompleter))
								if fields := strings.Fields(line); len(fields) > 0 {
									switch fields[0] {
									case "eof", "q", "quit":
										break Loop
									}
								}
								instr, ok, err := human.DecodeLine(line)
								if err != nil {
									fmt.Printf("error: %v\n", err)
									continue
								}
								if !ok {
									continue
								}
								code = append(code, instr)
							}
							state.Script = vm.Script{Code: code}
							return nil
						},
					},
					"load": skua.Command{
						Description: "load script from file",
						Run: func(args []string) error {
							var filename string
							if len(args) > 0 {
								filename = args[0]
							}
							code, err := load(filename)
							if err != nil {
								return err
							}
							state.Script = vm.Script{Code: code}
							return nil
						},
					},
					"save": skua.Command{
						Description: "save script to file",
						Subcommands: map[string]skua.Command{
							"bytecode": skua.Command{
								Description: "output to bytecode",
								Run: func(args []string) error {
									filename := firstString(args)
									b, err := bytecode.Marshal(state.Script.Code)
									if err != nil {
										return err
									}
									return saveFile(filename, b)
								},
							},
						},
						Run: func(args []string) error {
							filename := firstString(args)
							b, err := human.Marshal(state.Script.Code)
							if err != nil {
								return err
							}
							return saveFile(filename, b)
						},
					},
					"run": skua.Command{
						Description: "run script",
						Run: func([]string) error {
							runtime.Run(state)
							return nil
						},
					},
					"step": skua.Command{
						Description: "step script",
						Run: func([]string) error {
							runtime.Step(state)
							return nil
						},
					},
					"eval": skua.Command{
						Description: "eval instruction",
						Run: func(args []string) error {
							instr, err := human.DecodeArgs(args...)
							if err != nil {
								return err
							}
							runtime.Exec(state, instr)
							return nil
						},
						AdditionalSuggestions: func() []prompt.Suggest { return opCodeSuggestions },
					},
					"reset": skua.Command{
						Description: "reset instruction pointer",
						Run: func([]string) error {
							state.Script.Reset()
							return nil
						},
					},
				},
			},
			"dump": skua.Command{
				Description: "inspect state",
				Subcommands: map[string]skua.Command{
					"registers": skua.Command{
						Description: "inspect registers",
						Run: func([]string) error {
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Register #", "Value"})
							for i, val := range state.Registers {
								table.Append([]string{strconv.Itoa(i), strconv.Itoa(int(val))})
							}
							table.Render()
							return nil
						},
					},
					"stack": skua.Command{
						Description: "inspect stack",
						Run: func([]string) error {
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Frame #", "Values"})
							for i, frame := range state.Stack.Frames() {
								table.Append([]string{strconv.Itoa(i), fmt.Sprintf("%v", frame.Values())})
							}
							table.Render()
							return nil
						},
					},
				},
			},
			"quit": skua.Command{
				Description: "quit",
				Run:         func([]string) error { return skua.ErrQuit },
			},
		},
	}
}

func firstString(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func saveFile(filename string, p []byte) error {
	var w io.Writer
	if filename == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	_, err := w.Write(p)
	return err
}

func load(filename string) ([]vm.Instruction, error) {
	var (
		b   []byte
		err error
	)
	if filename == "" {
		b, err = ioutil.ReadAll(os.Stdin)
	} else {
		b, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return nil, err
	}
	code, err := bytecode.Unmarshal(b)
	if err != nil {
		code, err = human.Unmarshal(b)
		if err != nil {
			return nil, err
		}
	}
	return code, err
}
