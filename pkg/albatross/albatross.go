// Package albatross is a REPL builder.
package albatross

import (
	"errors"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora"
)

// DefaultPs ...
func DefaultPs(i int) string {
	return fmt.Sprintf("[%3d] > ", i+1)
}

// App ...
type App struct {
	Root    *Cmd
	Ps      func(int) string
	History []string
}

// New ...
func New(name string) *App {
	return &App{
		Root: Command(name),
		Ps:   DefaultPs,
	}
}

var (
	// ErrHelp ...
	ErrHelp = errors.New("help")
	// ErrQuit ...
	ErrQuit = errors.New("quit")
)

// Eval ...
func (app *App) Eval(tokens []string) error {
	cmd := app.Root
	var args []string
	for i, arg := range tokens {
		sub := cmd.Find(arg)
		if sub == nil {
			args = tokens[i:]
			break
		}
	}
	err := ErrHelp
	if cmd.Run != nil {
		err = cmd.Run(args)
	}
	if err == ErrHelp {
		return app.Help(tokens, args)
	}
	return err
}

// Help ...
func (app *App) Help(tokens, args []string) error {
	panic("not implemented")
}

// Completer ...
func (app *App) Completer(d prompt.Document) []prompt.Suggest {
	panic("not implemented")
}

// Step ...
func (app *App) Step() error {
	for {
		line := strings.TrimSpace(prompt.Input(app.Ps(len(app.History)), app.Completer, prompt.OptionHistory(app.History)))
		if line == "" {
			continue
		}
		app.History = append(app.History, line)
		return app.Eval(strings.Fields(line))
	}
}

// Loop ...
func (app *App) Loop() {
	for {
		if err := app.Step(); err != nil {
			if err == ErrQuit {
				return
			}
			fmt.Println(aurora.Red("error:"), err.Error())
		}
	}
}

// Command ...
func (app *App) Command(sub *Cmd) *App {
	app.Root.Command(sub)
	return app
}

// Do ..
func (app *App) Do(fn func(args []string) error) *App {
	app.Root.Run = fn
	return app
}

// Cmd ...
type Cmd struct {
	Name        string
	Aliases     []string
	Subcommands []*Cmd
	Args        []*Arg
	Help        string
	Run         func(args []string) error
}

// Command ...
func Command(name string) *Cmd {
	return &Cmd{Name: name}
}

// Find ...
func (cmd *Cmd) Find(name string) *Cmd {
	for _, sub := range cmd.Subcommands {
		if name == sub.Name {
			return sub
		}
		for _, alias := range sub.Aliases {
			if name == alias {
				return sub
			}
		}
	}
	return nil
}

// Command ...
func (cmd *Cmd) Command(sub *Cmd) *Cmd {
	cmd.Subcommands = append(cmd.Subcommands, sub)
	return cmd
}

// Arg ...
func (cmd *Cmd) Arg(name, help string, completer func(token string) []string) *Cmd {
	cmd.Args = append(cmd.Args, &Arg{
		Name: name,
	})
	return cmd
}

// Alias ...
func (cmd *Cmd) Alias(name string) *Cmd {
	cmd.Aliases = append(cmd.Aliases, name)
	return cmd
}

// Desc ...
func (cmd *Cmd) Desc(help string) *Cmd {
	cmd.Help = help
	return cmd
}

// Do ...
func (cmd *Cmd) Do(fn func(args []string) error) *Cmd {
	cmd.Run = fn
	return cmd
}

// Arg ...
type Arg struct {
	Name      string
	Help      string
	Completer func(token string) []string
}
