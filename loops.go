package pisc

var ModLoopCore = Module{
	Author:    "Andrew Owen",
	Name:      "LoopCore",
	License:   "MIT",
	DocString: "times is the standard counted loop, while takes two conditions",
	Load:      loadLoopCore,
}

func loadLoopCore(m *Machine) error {
	m.PredefinedWords["times"] = GoWord(func(m *Machine) error {
		toExec := m.PopValue().(*Quotation)
		nOfTimes := m.PopValue().(Integer)
		for i := int(0); i < int(nOfTimes); i++ {
			m.PushValue(toExec)
			err := m.ExecuteQuotation()
			if err != nil {
				return err
			}
		}
		return nil
	})

	// ( pred body -- .. )
	m.PredefinedWords["while"] = GoWord(func(m *Machine) error {
		body := m.PopValue().(*Quotation).toCode()
		pred := m.PopValue().(*Quotation).toCode()

		for {
			pred.Idx = 0
			err := m.execute(pred)
			if err != nil {
				return err
			}

			if !bool(m.PopValue().(Boolean)) {
				break
			}
			body.Idx = 0
			err = m.execute(body)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return m.importPISCAsset("stdlib/loops.pisc")
}
