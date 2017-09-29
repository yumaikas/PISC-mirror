// +build windows

package main

import (
	"pisc"
)

//#include <conio.h>
import "C"

var mode = "WINDOWS"

func init_term() error {
	return nil
	// stub, so that compilation doesn't fail
}

func de_init_term() error {
	return nil
}

func getch(m *pisc.Machine) error {
	char := C.getch()
	m.PusValue(pisc.Integer(char))
}
