package main

import "fmt"
import "reflect"

// GoWord a wrapper for functions that implement pieces of PISC
type GoWord func(*machine) error

// NilWord a wrapper for GoWords that should never fail
func NilWord(f func(*machine)) GoWord {
	return GoWord(func(m *machine) error {
		f(m)
		return nil
	})
}

func (m *machine) addGoWord(name, docstring string, impl GoWord) {
	m.helpDocs[name] = docstring
	m.predefinedWords[name] = impl
}

func t(m *machine) error {
	m.pushValue(Boolean(true))
	return nil
}

func f(m *machine) error {
	m.pushValue(Boolean(false))
	return nil
}

func dip(m *machine) error {
	quot := m.popValue().(*quotation).toCode()
	a := m.popValue()
	err := m.execute(quot)
	if err != nil {
		return err
	}
	m.pushValue(a)
	return nil
}

func pickDup(m *machine) error {
	distBack := int(m.popValue().(Integer))
	if distBack > len(m.values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.values))
	}
	m.pushValue(m.values[len(m.values)-distBack-1])
	return nil
}

func pickDrop(m *machine) error {
	distBack := int(m.popValue().(Integer))
	if distBack > len(m.values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.values))
	}
	valIdx := len(m.values) - distBack - 1
	val := m.values[valIdx]
	m.values = append(m.values[:valIdx], m.values[valIdx+1:]...)
	m.pushValue(val)
	return nil
}

func pickDel(m *machine) error {
	distBack := int(m.popValue().(Integer))
	if distBack > len(m.values)+1 {
		return fmt.Errorf("Cannot pick %v items back from stack of length %v", distBack, len(m.values))
	}
	valIdx := len(m.values) - distBack - 1
	m.values = append(m.values[:valIdx], m.values[valIdx+1:]...)
	return nil
}

func lenEntry(m *machine) error {
	length := m.popValue().(lenable).Length()
	m.pushValue(Integer(length))
	return nil
}

func errorFromEntry(m *machine) error {
	msg := m.popValue().String()
	return fmt.Errorf(msg)
}

func reflectEq(m *machine) error {
	a := m.popValue()
	b := m.popValue()
	m.pushValue(Boolean(reflect.DeepEqual(a, b)))
	return nil
}

var ModPISCCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "PISCCore",
	License:   "MIT",
	DocString: "Eventally, the small batch of core PISC words",
	Load:      loadPISCCore,
}

// These are the standard libraries that are currently trusted to not
var StandardModules = []PISCModule{
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

func (m *machine) loadForCLI() error {
	return m.LoadModules(append(StandardModules,
		ModIOCore, ModDebugCore, ModShellUtils)...)
}
func (m *machine) loadForDB() error {
	return m.LoadModules(append(StandardModules,
		ModBoltDB, ModIOCore, ModIOCore, ModShellUtils)...)
}

func (m *machine) loadForChatbot() error {
	return m.LoadModules(append(StandardModules,
		ModBoltDB, ModIRCKit)...)
}

func loadPISCCore(m *machine) error {
	if m.predefinedWords == nil {
		panic("Uninitialized stack machine!")
	}
	m.addGoWord("t", "( -- t )", GoWord(t))
	m.addGoWord("f", "( -- f )", GoWord(f))
	m.addGoWord("dip", "( a quot -- ... a )", GoWord(dip))
	m.predefinedWords["pick-dup"] = GoWord(pickDup)
	m.predefinedWords["pick-drop"] = GoWord(pickDrop)
	m.predefinedWords["pick-del"] = GoWord(pickDel)
	m.addGoWord("len", "( e -- lenOfE ) ", GoWord(lenEntry))
	m.addGoWord("eq", " ( a b -- same? ) ", GoWord(runEq))
	// Discourage use of reflection based eq via long name
	m.addGoWord("deep-slow-reflect-eq", "( a b -- same? )", GoWord(reflectEq))
	m.addGoWord("error", "( msg -- !! )", GoWord(errorFromEntry))
	return m.importPISCAsset("stdlib/std_lib.pisc")
}
