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

	t.Fail()
	t.Log(m.definedWords)
	t.Log(m.definedStackComments)
	t.Log(m.values)
}
