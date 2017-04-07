package main

var ModLoopCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "LoopCore",
	License:   "MIT",
	DocString: "times is the standard counted loop, while takes two conditions",
	Load:      loadLoopCore,
}

func loadLoopCore(m *machine) error {
	m.predefinedWords["times"] = GoWord(func(m *machine) error {
		toExec := m.popValue().(*quotation)
		nOfTimes := m.popValue().(Integer)
		for i := int(0); i < int(nOfTimes); i++ {
			m.pushValue(toExec)
			err := m.executeQuotation()
			if err != nil {
				return err
			}
		}
		return nil
	})
	// ( pred body -- .. )
	m.predefinedWords["while"] = GoWord(func(m *machine) error {
		body := m.popValue().(*quotation).toCode()
		pred := m.popValue().(*quotation).toCode()

		for {
			pred.idx = 0
			err := m.execute(pred)
			if err != nil {
				return err
			}

			if !bool(m.popValue().(Boolean)) {
				break
			}
			body.idx = 0
			err = m.execute(body)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return m.importPISCAsset("stdlib/loops.pisc")
}
