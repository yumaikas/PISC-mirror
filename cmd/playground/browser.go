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
	m.AddGoWordWithStack("create-paste-input", "( name:str -- )", "Create a paste input with the supplied name", addTextInput)
	m.AddGoWordWithStack("get-paste-text", "( name:str -- text )", "Get the paste text from the input with the supplied name", getTextOfInput)
	// TODO: Add in xhr/http bindings to make this able to drive APIs from the browser
	return nil
}
