// +build !windows

package main

import (
	"fmt"
	"github.com/pkg/term"
	"pisc"
)

var mode = "POSIX"

var pty *term.Term
var buf = make([]byte, 1)

func init_term() error {
	var err error
	pty, err = term.Open("/dev/tty", term.RawMode)
	if err != nil {
		return err
	}
	return nil
}

func de_init_term() error {
	return pty.Restore()
}

func os_overload(m *pisc.Machine) error {
	return nil
}

func getch(m *pisc.Machine) error {
	num, err := pty.Read(buf)

	if err != nil {
		return err
	}
	if num < 1 {
		return fmt.Errorf("Didn't read at least 1 byte from the pty!")
	}

	m.PushValue(pisc.Integer(int(buf[0])))
	return nil
}
