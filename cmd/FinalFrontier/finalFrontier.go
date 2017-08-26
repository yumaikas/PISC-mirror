package main

import (
	"os"
	"pisc"
)

//#include<conio.h>
import "C"

// Windows only right now!
func getch(m *pisc.Machine) error {
	char := C.getch()
	m.PushValue(pisc.Integer(char))
	return nil
}

func loadGameScript(m *pisc.Machine) error {
	assetKey := m.PopValue().String()
	data, err := Asset(string(assetKey))
	if err != nil {
		return err
	}
	err = m.ExecuteString(string(data), pisc.CodePosition{Source: "file:" + string(assetKey)})
	if err != nil {
		return err
	}
	return nil
}

func exit(m *pisc.Machine) error {
	os.Exit(0)
	return nil
}

var ModFinalFrontier = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "ConsoleIO",
	License:   "MIT",
	DocString: "Provides getkey on Windows",
	Load:      loadFinalFrontier,
}

func loadFinalFrontier(m *pisc.Machine) error {
	m.AddGoWord("getkey", "( -- keyval )", getch)
	m.AddGoWord("game-script", "( filename -- ? ) ", loadGameScript)
	m.AddGoWord("quit-game", "( -- )", exit)
	return nil
}
