package main

func isBooleanWord(w word) bool {
	return w == "and" ||
		w == "or" ||
		w == "not"
}

func (m *machine) executeBooleanWord(w word) error {
	switch w {
	case "and":
		a := m.popValue().(Boolean)
		b := m.popValue().(Boolean)
		m.pushValue(Boolean(a && b))
	case "or":
		a := m.popValue().(Boolean)
		b := m.popValue().(Boolean)
		m.pushValue(Boolean(a || b))
	case "not":
		a := m.popValue().(Boolean)
		m.pushValue(Boolean(!a))
	}
	return nil
}
