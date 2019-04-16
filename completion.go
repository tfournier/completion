package completion

import (
	"fmt"
	"io"
	"os"
	"path"
)

// ICommand return interface for choice shell
type ICommand interface {
	Zsh() ICompletion
}

// ICompletion return interface for choice output
type ICompletion interface {
	ToString() string
	ToWriter(w io.Writer) error
	ToFile(filename string) error
}

type completion struct {
	string
}

func (cmd *command) Zsh() ICompletion {
	return completion{zsh(cmd)}
}

func (c completion) ToString() string {
	return c.string
}

func (c completion) ToWriter(w io.Writer) error {
	_, err := fmt.Fprint(w, c.string)
	return err
}

func (c completion) ToFile(filename string) error {

	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return c.ToWriter(outFile)
}

type command struct {
	Name        string
	Description string
	Alias       []string
	Flags       []flag
	Arguments   []string
	SubCommands []*command
	Parent      *command
}

func (cmd command) HasAlias() bool {
	return len(cmd.Alias) > 0
}

func (cmd command) HasFlags() bool {
	return len(cmd.Flags) > 0
}

func (cmd command) HasArguments() bool {
	return len(cmd.Arguments) > 0
}

func (cmd command) HasSubCommands() bool {
	return len(cmd.SubCommands) > 0
}

func (cmd command) HasParent() bool {
	return cmd.Parent != nil
}

type flag struct {
	Name        string
	Shorthand   string
	Description string
}

func (flag flag) HasShorthand() bool {
	return len(flag.Shorthand) > 0
}
