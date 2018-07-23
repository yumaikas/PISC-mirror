package pisc

import "errors"

var ModLoopCore = Module{
	Author:    "Andrew Owen",
	Name:      "LoopCore",
	License:   "MIT",
	DocString: "Implements basic loops",
	Load:      loadLoopCore,
}

var ErrBreakLoop = errors.New("Breaking out of a loop")
var ErrContinueLoop = errors.New("Early exit on a Loop")

func _doTimes(m *Machine) error {
	toExec := m.PopValue().(*Quotation)
	nOfTimes := m.PopValue().(Integer)
	for i := int(0); i < int(nOfTimes); i++ {
		err := m.CallQuote(toExec)
		if IsLoopError(err) && LoopShouldEnd(err) {
			return nil
		}
		if IsLoopError(err) && !LoopShouldEnd(err) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func LoopShouldEnd(e error) bool {
	return e == ErrBreakLoop
}

func IsLoopError(e error) bool {
	return e == ErrBreakLoop || e == ErrContinueLoop
}

func _doWhile(m *Machine) error {
	body := m.PopValue().(*Quotation)
	pred := m.PopValue().(*Quotation)
	for {
		err := m.CallQuote(pred)
		if err != nil {
			return err
		}
		if !bool(m.PopValue().(Boolean)) {
			break
		}
		err = m.CallQuote(body)
		if IsLoopError(err) && LoopShouldEnd(err) {
			return nil
		}
		if IsLoopError(err) && !LoopShouldEnd(err) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func _break(m *Machine) error {
	return ErrBreakLoop
}

func _continue(m *Machine) error {
	return ErrContinueLoop
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

	m.AddGoWordWithStack(
		"break",
		"( -- ! )",
		"Break out of a loop body",
		_break)
	m.AddGoWordWithStack(
		"continue",
		"( -- ! )",
		"End this iteration of a loop body",
		_continue)

	return m.ImportPISCAsset("stdlib/loops.pisc")
}
