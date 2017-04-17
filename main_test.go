package main

import "testing"

func TestRunPISCTests(t *testing.T) {
	m := initMachine()
	err := m.loadForCLI()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = m.executeString(`"tests/all.pisc" import`, codePosition{source: "go test"})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
