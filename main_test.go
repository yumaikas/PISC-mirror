package main

import "testing"

func TestParser(t *testing.T) {
	m := runCode(`1 2 dup`)

	a := m.popValue().(Integer)
	b := m.popValue().(Integer)
	c := m.popValue().(Integer)
	if a != Integer(2) && b != Integer(2) && c != Integer(3) {
		t.Fail()
		t.Log("Stack values were: ", a, b, c)
	}
}

func TestExecution(t *testing.T) {
	m := runCode(`t [ 1 ] [ 2 ] ? call
 f [ 1 ] [ 2 ] ? call`)
	a := m.popValue().(Integer)
	b := m.popValue().(Integer)
	if a != Integer(2) && b != Integer(1) {
		t.Fail()
		t.Log("Stack values were", a, b, m.values)
	}
}

func TestCall(t *testing.T) {
	m := runCode(`[ 1 ] call`)
	a := m.popValue().(Integer)
	if a != Integer(1) {
		t.Fail()
		t.Log("Stack values were", a, m.values)
	}

}

func TestWordDefinition(t *testing.T) {
	m := runCode(`: if ( ? a b -- x ) ? call ;
		f [ 2 ] [ 4 ] if`)
	a := m.popValue().(Integer)
	ifWord := m.definedWords[word("if")]
	ifComment := m.definedStackComments[word("if")]
	if a != Integer(4) ||
		ifWord[0] != word("?") ||
		ifWord[1] != "call" ||
		ifComment != "? a b -- x" {
		t.Fail()
		t.Log("Not everything was as expected")
		t.Log(m.definedWords)
		t.Log(m.definedStackComments)
		t.Log(m.values)
	}
}

func TestIntAddition(t *testing.T) {
	m := runCode(`1 2 +`)
	a := m.popValue().(Integer)
	if a != Integer(3) {
		t.Fail()
		t.Log("Stack values:", a, m.values)
	}
}

// This is just a test to make sure that performance isn't hugely sucky.
func TestManyIntAddition(t *testing.T) {
	m := runCode(`1 2 3 4 5 6 300 2 3 2 3 2 3 2 3 2 3 + + + + + + + + + + + + + + + +`)
	a := m.popValue().(Integer)
	if a != Integer(346) {
		t.Fail()
		t.Log("Stack values:", a, m.values)
	}
}
