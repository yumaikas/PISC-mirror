// +build windows

package main

import (
	"pisc"
    "github.com/mattn/go-colorable"
    "io"
    "fmt"
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
    m.AddGoWord("priv_puts", "( str -- )", printWithAnsiColor)
    return nil
}


func getch(m *pisc.Machine) error {
	char := C.getch()
	m.PushValue(pisc.Integer(char))
    return nil
}
