package main

import (
	"os"
	"pisc"
)

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

func getTermMode(m *pisc.Machine) error {
	m.PushValue(pisc.String(mode))
	return nil
}

func exit(m *pisc.Machine) error {
	de_init_term()
	os.Exit(0)
	return nil
}

var ModFinalFrontier = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "ConsoleIO",
	License:   "MIT",
	DocString: "Provides basic OS-interaction words",
	Load:      loadFinalFrontier,
}

func loadFinalFrontier(m *pisc.Machine) error {
	init_term()
	m.AddGoWordWithStack("term-mode", "( -- mode:str )", "Get the current game mode", getTermMode)
	m.AddGoWordWithStack("getkey", "( -- keyval:str )", "Wait for a key to be pressed, return its value", getch)
	m.AddGoWordWithStack("game-script", "( filename:str -- ? ) ", "Execute a game script", loadGameScript)
	m.AddGoWordWithStack("quit-game", "( -- )", "Exit the game.", exit)
	os_overload(m)
	return nil
}
