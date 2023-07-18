package utils

import (
	"fmt"

	"github.com/Netflix/go-expect"
	pseudotty "github.com/creack/pty"
	"github.com/hinshun/vt10x"
)

var (
	expectNewConsole = expect.NewConsole
)

// NewVT10XConsole returns a new expect.Console that multiplexes the
// Stdin/Stdout to a VT10X terminal, allowing Console to interact with an
// application sending ANSI escape sequences.
func NewVT10XConsole(opts ...expect.ConsoleOpt) (*expect.Console, error) {
	pty, tty, err := pseudotty.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open pseudotty: %v", err)
	}

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expectNewConsole(append(opts, expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))...)
	if err != nil {
		return nil,  fmt.Errorf("failed to create console: %v", err)
	}

	return c, nil
}
