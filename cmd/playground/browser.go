package main

import (
	"pisc"

	"github.com/gopherjs/gopherjs/js"
)

var ModPlayground = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "ModPlayground",
	License:   "MIT",
	DocString: "Words for interacting with the browser playground",
	Load:      loadModPlayground,
}

func addTextInput(m *pisc.Machine) error {
	name := m.PopValue().String()
	js.Global.Get("app").Call("addPasteInput", name)
	return nil
}

func getTextOfInput(m *pisc.Machine) error {
	name := m.PopValue().String()
	val := js.Global.Get("app").Call("getTextOfPaste", name).String()
	m.PushValue(pisc.String(val))
	return nil
}

func loadModPlayground(m *pisc.Machine) error {
	m.AddGoWord("create-paste-input", "( name -- )", addTextInput)
	m.AddGoWord("get-paste-text", "( name -- text )", getTextOfInput)
	return nil
}
