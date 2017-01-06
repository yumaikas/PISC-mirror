package main

func (m *machine) loadLoopWords() error {
	m.predefinedWords["times"] = GoWord(func(m *machine) error {

		toExec := m.popValue().(*quotation).toCode()
		nOfTimes := m.popValue().(Integer)
		for i := int(0); i < int(nOfTimes); i++ {
			toExec.idx = 0
			err := m.execute(toExec)
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
			m.execute(body)
		}
		return nil
	})

	return nil
}
