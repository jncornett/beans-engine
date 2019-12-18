package skua

import (
	"errors"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

// RunFunc ...
type RunFunc func(args []string) error

// SuggestionsFunc ...
type SuggestionsFunc func() []prompt.Suggest

// Command ...
type Command struct {
	Subcommands           map[string]Command
	Description           string
	Run                   RunFunc
	AdditionalSuggestions SuggestionsFunc
}

// Suggestions ...
func (cmd Command) Suggestions() []prompt.Suggest {
	var out []prompt.Suggest
	for name, sub := range cmd.Subcommands {
		out = append(out, prompt.Suggest{
			Text:        name,
			Description: sub.Description,
		})
	}
	if cmd.AdditionalSuggestions != nil {
		out = append(out, cmd.AdditionalSuggestions()...)
	}
	return out
}

// Exec ...
func (cmd Command) Exec(repl *Repl, args []string) error {
	if len(args) > 0 {
		sub, ok := cmd.Subcommands[args[0]]
		if ok {
			return sub.Exec(repl, args[1:])
		}
	}
	if cmd.Run != nil {
		return cmd.Run(args)
	}
	return nil
}

// Repl ...
type Repl struct {
	History  []string
	Commands map[string]Command
	OnError  func(error) string
}

// ErrQuit ...
var ErrQuit = errors.New("quit")

// Step ...
func (repl *Repl) Step() error {
	line := prompt.Input(
		// TODO make PS configurable
		fmt.Sprintf("[%2d] > ", len(repl.History)),
		repl.completer,
		prompt.OptionHistory(repl.History),
	)
	line = strings.TrimSpace(line)
	fields := strings.Fields(line)
	root := Command{Subcommands: repl.Commands}
	if err := root.Exec(repl, fields); err != nil {
		return err
	}
	repl.History = append(repl.History, line)
	return nil
}

// Loop ...
func (repl *Repl) Loop() {
	for {
		if err := repl.Step(); err != nil {
			if err == ErrQuit {
				return
			}
			repl.onError(err)
		}
	}
}

func (repl *Repl) completer(d prompt.Document) []prompt.Suggest {
	cmd := Command{Subcommands: repl.Commands}
	getPath := func(path []string) *Command {
		cur := cmd
		for _, part := range path {
			var ok bool
			cur, ok = cur.Subcommands[part]
			if !ok {
				return nil
			}
		}
		return &cur
	}
	line := d.CurrentLine()
	word := d.GetWordBeforeCursor()
	path := strings.Fields(line[:d.FindStartOfPreviousWordWithSpace()])
	if strings.HasSuffix(line, " ") {
		path = strings.Fields(line)
	}
	if len(path) == 0 {
		return prompt.FilterHasPrefix(cmd.Suggestions(), word, false)
	}
	sub := getPath(path)
	if sub == nil {
		return nil
	}
	return prompt.FilterHasPrefix(sub.Suggestions(), word, false)
}

func (repl *Repl) onError(err error) {
	if repl.OnError == nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	repl.OnError(err)
}
