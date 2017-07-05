package pisc

import "testing"

func TestRunPISCTests(t *testing.T) {
	m := initMachine()
	err := m.loadForCLI()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = m.ExecuteString(`"tests/all.pisc" import`, CodePosition{Source: "go test"})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
