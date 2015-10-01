package main

// TODO: Add more words to support strings here, we need a way to handle a lot more cases, like
// replacement, substringing, joining and so on.
func isStringWord(w word) bool {
	return w == "concat" || w == "empty?" || w == ">string"
}

func (m *machine) executeStringWord(w word) error {
	switch w {
	case "concat":
		a := m.popValue().(String)
		b := m.popValue().(String)
		m.pushValue(String(b + a))
	case ">string":
		a := m.popValue()
		if s, ok := a.(String); ok {
			m.pushValue(String(s))
			return nil
		}
		m.pushValue(String(a.String()))
	case "empty?":
		a := m.popValue().(String)
		if len(a) > 0 {
			m.pushValue(Boolean(false))
		}
	}
	return nil
}
