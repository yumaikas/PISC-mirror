package pisc

import "fmt"
import "reflect"

// GoWord a wrapper for functions that implement pieces of PISC
type GoWord func(*Machine) error

// NilWord a wrapper for GoWords that should never fail
func NilWord(f func(*Machine)) GoWord {
	return GoWord(func(m *Machine) error {
		f(m)
		return nil
	})
}

func (m *Machine) AddGoWord(name, docstring string, impl GoWord) {
	m.HelpDocs[name] = docstring
	m.PredefinedWords[name] = impl
}

func t(m *Machine) error {
	m.PushValue(Boolean(true))
	return nil
}

func f(m *Machine) error {
	m.PushValue(Boolean(false))
	return nil
}

func condOperator(m *Machine) error {
	return m.runConditionalOperator()
}

func call(m *Machine) error {
	return m.ExecuteQuotation()
}

func isStackEmpty(m *Machine) error {
	m.PushValue(Boolean(len(m.Values) == 0))
	return nil
}

func typeof(m *Machine) error {
	m.PushValue(String(m.PopValue().Type()))
	return nil
}

func dip(m *Machine) error {
	quot := m.PopValue().(*quotation).toCode()
	a := m.PopValue()
	err := m.execute(quot)
	if err != nil {
		return err
	}
	m.PushValue(a)
	return nil
}

func pickDup(m *Machine) error {
	distBack := int(m.PopValue().(Integer))
	if distBack > len(m.Values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.Values))
	}
	m.PushValue(m.Values[len(m.Values)-distBack-1])
	return nil
}

func pickDrop(m *Machine) error {
	distBack := int(m.PopValue().(Integer))
	if distBack > len(m.Values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.Values))
	}
	valIdx := len(m.Values) - distBack - 1
	val := m.Values[valIdx]
	m.Values = append(m.Values[:valIdx], m.Values[valIdx+1:]...)
	m.PushValue(val)
	return nil
}

func pickDel(m *Machine) error {
	distBack := int(m.PopValue().(Integer))
	if distBack > len(m.Values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.Values))
	}
	valIdx := len(m.Values) - distBack - 1
	m.Values = append(m.Values[:valIdx], m.Values[valIdx+1:]...)
	return nil
}

func lenEntry(m *Machine) error {
	length := m.PopValue().(Lenable).Length()
	m.PushValue(Integer(length))
	return nil
}

func errorFromEntry(m *Machine) error {
	msg := m.PopValue().String()
	return fmt.Errorf(msg)
}

func reflectEq(m *Machine) error {
	a := m.PopValue()
	b := m.PopValue()
	m.PushValue(Boolean(reflect.DeepEqual(a, b)))
	return nil
}

var ModPISCCore = Module{
	Author:    "Andrew Owen",
	Name:      "PISCCore",
	License:   "MIT",
	DocString: "Eventally, the small batch of core PISC words",
	Load:      loadPISCCore,
}

// These are the standard libraries that are currently trusted to not
var StandardModules = []Module{
	ModLoopCore,
	ModLocalsCore,
	ModDictionaryCore,
	ModHelpCore,
	ModStringsCore,
	ModMathCore,
	ModBoolCore,
	ModVectorCore,
	ModSymbolCore,
	ModRandomCore,
	ModMetaQuoation,
	ModPISCCore,
}

func loadPISCCore(m *Machine) error {
	if m.PredefinedWords == nil {
		panic("Uninitialized stack machine!")
	}
	m.AddGoWord("t", "( -- t )", GoWord(t))
	m.AddGoWord("f", "( -- f )", GoWord(f))
	m.AddGoWord("dip", "( a quot -- ... a )", GoWord(dip))
	m.AddGoWord("stack-empty?", "( -- empty? )", GoWord(isStackEmpty))
	m.AddGoWord("typeof", "( a -- typeofa )", GoWord(typeof))
	m.AddGoWord("?", "( a b ? -- a/b )", GoWord(condOperator))
	m.AddGoWord("call", "( quot -- ... )", GoWord(call))
	m.PredefinedWords["pick-dup"] = GoWord(pickDup)
	m.PredefinedWords["pick-drop"] = GoWord(pickDrop)
	m.PredefinedWords["pick-del"] = GoWord(pickDel)
	m.AddGoWord("len", "( e -- lenOfE ) ", GoWord(lenEntry))
	m.AddGoWord("eq", " ( a b -- same? ) ", GoWord(runEq))
	// Discourage use of reflection based eq via long name
	m.AddGoWord("deep-slow-reflect-eq", "( a b -- same? )", GoWord(reflectEq))
	m.AddGoWord("error", "( msg -- !! )", GoWord(errorFromEntry))
	return m.importPISCAsset("stdlib/std_lib.pisc")
}
