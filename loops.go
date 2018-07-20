package pisc

var ModLoopCore = Module{
	Author:    "Andrew Owen",
	Name:      "LoopCore",
	License:   "MIT",
	DocString: "times is the standard counted loop, while takes two conditions",
	Load:      loadLoopCore,
}

func _doTimes(m *Machine) error {
	toExec := m.PopValue().(*Quotation).toCode()
	nOfTimes := m.PopValue().(Integer)
	for i := int(0); i < int(nOfTimes); i++ {
		err := m.execute(toExec)
		if err != nil {
			return err
		}
	}
	return nil
}

func _doWhile(m *Machine) error {
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
}

func loadLoopCore(m *Machine) error {
	m.AddGoWordWithStack(
		"times",
		"( n quot -- ... )",
		"Call quot n times",
		_doTimes)
	// ( pred body -- .. )
	m.AddGoWordWithStack(
		"while",
		"( pred quot -- ... )",
		"Call quot while pred leave true at the top of the stack",
		_doWhile)

	return m.ImportPISCAsset("stdlib/loops.pisc")
}
