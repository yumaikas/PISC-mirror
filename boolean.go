package pisc

var ModBoolCore = Module{
	Author:    "Andrew Owen",
	Name:      "BoolCore",
	License:   "MIT",
	DocString: "The 3 basic boolean operation, and/or and not",
	Load:      loadBoolCore,
}

func _and(m *Machine) error {
	a := m.PopValue().(Boolean)
	b := m.PopValue().(Boolean)
	m.PushValue(Boolean(a && b))
	return nil
}

func _or(m *Machine) error {
	a := m.PopValue().(Boolean)
	b := m.PopValue().(Boolean)
	m.PushValue(Boolean(a || b))
	return nil
}

func _not(m *Machine) error {
	a := m.PopValue().(Boolean)
	m.PushValue(Boolean(!a))
	return nil
}

func loadBoolCore(m *Machine) error {
	m.AppendToHelpTopic("Bools", "PISC boolean expressions are generally not short-circuited without using quotations")
	m.AddGoWordWithStack("and", "( a b -- a&b )", "Boolean And", _and)
	m.AddGoWordWithStack("or", "( a b -- a||b )", "Boolean OR", _or)
	m.AddGoWordWithStack("not", "( a  -- not-a )", "Boolean NOT", _not)

	return m.ImportPISCAsset("stdlib/bools.pisc")
}
