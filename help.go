package main

// These are the help words

func (m *machine) loadHelpWords() {
	// ( val name -- )
	m.predefinedWords["help"] = GoWord(func(m *machine) error {
		searchTerm := m.popValue().(String)
		if len(m.locals) <= 0 {
			return ErrNoLocalsExist
		}
		m.locals[len(m.locals)-1][varName.String()] = m.popValue()
		return nil
	})

}
