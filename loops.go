package main

func (m *machine) loadLoopWords() error {
	m.predefinedWords["times"] = GoWord(func(m *machine) error {

		toExec := m.popValue().(quotation).toCode()
		nOfTimes := m.popValue().(Integer)
		for i := int(0); i < int(nOfTimes); i++ {
			toExec.idx = 0
			err := executeWordsOnMachine(m, toExec)
			if err != nil {
				return err
			}
		}
		return nil
	})
	// ( pred body -- .. )
	m.predefinedWords["while"] = GoWord(func(m *machine) error {
		body := m.popValue().(quotation).toCode()
		pred := m.popValue().(quotation).toCode()

		for {
			pred.idx = 0
			err := executeWordsOnMachine(m, pred)
			if err != nil {
				return err
			}

			if !bool(m.popValue().(Boolean)) {
				break
			}
			body.idx = 0
			executeWordsOnMachine(m, body)
		}
		return nil
	})

	/*

		m.predefinedWords["each"] = NilWord(func(m *machine) {

		})
					case "while":
						body := m.popValue().(quotation)
						pred := m.popValue().(quotation)

			case "while":
				body := m.popValue().(quotation)
				predicate := m.popValue().(quotation)

				for {

				}
	*/
	return nil
}
