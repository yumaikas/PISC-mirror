package main

import (
	"fmt"
	"pisc"
)

func main() {
	m := &pisc.Machine{
		Values:               make([]pisc.StackEntry, 0),
		DefinedWords:         make(map[string]*pisc.CodeQuotation),
		DefinedStackComments: make(map[string]string),
		PredefinedWords:      make(map[string]pisc.GoWord),
		PrefixWords:          make(map[string]*pisc.CodeQuotation),
		HelpDocs:             make(map[string]string),
	}

	m.LoadModules(append(
		pisc.StandardModules,
		pisc.ModDebugCore,
		pisc.ModIOCore,
		ModFinalFrontier)...)
	err := m.ExecuteString(`"scripts/main.pisc" game-script`, pisc.CodePosition{Source: "main.go"})
	if err != nil {
		fmt.Println(err.Error())
	}
}
