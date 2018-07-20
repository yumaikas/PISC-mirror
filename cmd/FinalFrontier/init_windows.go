// +build windows

package main

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"io"
	"pisc"
)

//#include <conio.h>
import "C"

var mode = "WINDOWS"
var writer io.Writer

func init_term() error {
	writer = colorable.NewColorableStdout()
	return nil
	// stub, so that compilation doesn't fail
}

func de_init_term() error {
	return nil
}

func printWithAnsiColor(m *pisc.Machine) error {
	str := m.PopValue().(pisc.String)
	_, err := fmt.Fprint(writer, str)
	return err
}

func os_overload(m *pisc.Machine) error {
	// Overwritting priv_puts to handle ANSI color codes on windows
	m.AddGoWordWithStack("priv_puts", "( output:str -- )", "Emit the string with ANSI colors", printWithAnsiColor)
	return nil
}

func getch(m *pisc.Machine) error {
	char := C.getch()
	m.PushValue(pisc.Integer(char))
	return nil
}
