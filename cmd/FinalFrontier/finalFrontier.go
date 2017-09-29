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
	m.AddGoWord("term-mode", "( -- mode )", getTermMode)
	m.AddGoWord("getkey", "( -- keyval )", getch)
	m.AddGoWord("game-script", "( filename -- ? ) ", loadGameScript)
	m.AddGoWord("quit-game", "( -- )", exit)
	return nil
}
