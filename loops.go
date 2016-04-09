package main

// TODO: add these words to the main loop.
func isLoopWord(w word) bool {
	return w == "times"
}

func (m *machine) executeLoopWord(w word) error {
	switch w {
	case "times":
		body := m.popValue().(quotation)
		nOfTimes := m.popValue().(Integer)
		toExec := &codeList{
			idx:    0,
			code:   body,
			spaces: make([]string, 0),
		}
		for i := int(0); i < int(nOfTimes); i++ {
			toExec.idx = 0
			executeWordsOnMachine(m, toExec)
		}
		/*
			case "while":
				body := m.popValue().(quotation)
				pred := m.popValue().(quotation)
		*/
	}
	return nil
}
