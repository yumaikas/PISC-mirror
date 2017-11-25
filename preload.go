package pisc

import "fmt"
import "reflect"
import "strings"

// GoWord a wrapper for functions that implement pieces of PISC
type GoWord func(*Machine) error

// NilWord a wrapper for GoWords that should never fail
func NilWord(f func(*Machine)) GoWord {
	return GoWord(func(m *Machine) error {
		f(m)
		return nil
	})
}

func docify(doc string) string {
	return strings.Replace(doc, "<code>", "`", -1)
}

func (m *Machine) AddGoWordWithStack(name, stackEffect, docstring string, impl GoWord) {
	m.DefinedStackComments[name] = stackEffect
	m.HelpDocs[name] = docify(docstring)
	m.PredefinedWords[name] = impl
}

func (m *Machine) AppendToHelpTopic(topic, content string) {
	if help, found := m.HelpDocs[topic]; found {
		m.HelpDocs[topic] = help + docify(content)
		return
	}
	m.HelpDocs[topic] = docify(content)
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
	quot := m.PopValue().(*Quotation).toCode()
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

func isModuleLoaded(m *Machine) error {
	modName := m.PopValue().String()
	for _, mod := range m.LoadedModules {
		if modName == mod {
			m.PushValue(Boolean(true))
			return nil
		}
	}
	m.PushValue(Boolean(false))
	return nil
}

var ModPISCCore = Module{
	Author:    "Andrew Owen",
	Name:      "PISCCore",
	License:   "MIT",
	DocString: "Eventally, the small batch of core PISC words",
	Load:      loadPISCCore,
}

// These are the standard libraries that are currently trusted to not cause problems in general
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
	ModPISCCore,
}

func loadPISCCore(m *Machine) error {
	if m.PredefinedWords == nil {
		panic("Uninitialized stack machine!")
	}
	m.AddGoWordWithStack("t", "( -- t )", "The True constant", t)
	m.AddGoWordWithStack("f", "( -- f )", "The False constant", f)
	m.AddGoWordWithStack(
		"dip", "( a quot -- ... a )",
		"Execute a quoation without the top of the stack, and then restore the element at the top",
		dip)

	m.AddGoWordWithStack(
		"stack-empty?",
		"( -- empty? )",
		"Returns true if the stack is empty, false if isn't",
		isStackEmpty)

	m.AddGoWordWithStack(
		"typeof",
		"( a -- typeofa )",
		"Get the type of the value at the top of the stack", typeof)

	m.AddGoWordWithStack("?",
		"( a b ? -- a/b )",
		"The contintional operator: Takes a, b and a boolean. \n"+
			"Returns a if the boolean is true, b if it is false",
		condOperator)

	m.AddGoWordWithStack(
		"call",
		"( quot -- ... )",
		"Call the quotation or callable that is on the stack",
		call)

	m.AddGoWordWithStack("len",
		"( e -- lenOfE )",
		"Check the length of a given collection or string using a Go-side Length interface ",
		lenEntry)
	m.AddGoWordWithStack("eq",
		" ( a b -- same? ) ",
		"Run shallow equality",
		runEq)
	// Discourage use of reflection based eq via long name
	m.AddGoWordWithStack(
		"deep-slow-reflect-eq",
		"( a b -- same? )",
		"Run a deep, refection based comparison. Slower than reflect-eq, but easier to use for vectors",
		reflectEq)
	m.AddGoWordWithStack(
		"error",
		"( message -- !! )",
		"Create an error from msg",
		errorFromEntry)
	m.AddGoWordWithStack(
		"module-loaded?",
		"( module-name -- loaded? )",
		"Check to see if a module with the given name is loaded",
		isModuleLoaded)

	// These words are important for combinators,
	// but aren't documented on the PISC side
	m.PredefinedWords["pick-dup"] = GoWord(pickDup)
	m.PredefinedWords["pick-drop"] = GoWord(pickDrop)
	m.PredefinedWords["pick-del"] = GoWord(pickDel)

	return m.ImportPISCAsset("stdlib/std_lib.pisc")
}
