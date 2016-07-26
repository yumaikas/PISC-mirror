package main

func (m *machine) loadLoopWords() error {
	m.predefinedWords["times"] = GoWord(func(m *machine) error {

		body := m.popValue().(quotation)
		nOfTimes := m.popValue().(Integer)
		toExec := &codeList{idx: 0, code: body}
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
	m.predefinedWords["while"] = NilWord(func(m *machine) {
		body := m.popValue().(quotation)
		pred := m.popValue().(quotation)
		predExec := &codeList{idx: 0, code: pred}
		bodyExec := &codeList{idx: 0, code: body}

		for {
			predExec.idx = 0
			executeWordsOnMachine(m, predExec)

			if !bool(m.popValue().(Boolean)) {
				break
			}
			bodyExec.idx = 0
			executeWordsOnMachine(m, bodyExec)
		}
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
