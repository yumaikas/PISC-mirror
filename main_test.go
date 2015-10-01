package main

import "testing"

func TestParser(t *testing.T) {
	m := runCode(`1 2 dup`)

	a := m.popValue().(Integer)
	b := m.popValue().(Integer)
	c := m.popValue().(Integer)
	if a != Integer(2) || b != Integer(2) || c != Integer(1) {
		t.Fail()
		t.Log("Stack values were: ", a, b, c)
	}
}

func TestExecution(t *testing.T) {
	m := runCode(`t [ 1 ] [ 2 ] ? call
 f [ 1 ] [ 2 ] ? call`)
	a := m.popValue().(Integer)
	b := m.popValue().(Integer)
	if a != Integer(2) || b != Integer(1) {
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

func TestAdditionTypePromotion(t *testing.T) {
	m := runCode(`1 2.0 +`)
	a := m.popValue().(Double)
	if a != Double(3) {
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

func TestMath(t *testing.T) {
	m := runCode(`1 2 -`)
	a := m.popValue().(Integer)
	if a != Integer(-1) {
		t.Fail()
		t.Log(m.values)
	}
}

func TestBooleanOr(t *testing.T) {
	m := runCode(`t f or`)
	a := m.popValue().(Boolean)
	if a != Boolean(true) {
		t.Fail()
		t.Log(m.values)
	}
}

func TestBooleanAnd(t *testing.T) {
	m := runCode(`t f and`)
	a := m.popValue().(Boolean)
	if a != Boolean(false) {
		t.Fail()
		t.Log(m.values)
	}
}

func TestEven(t *testing.T) {
	m := runCode(`: even? ( x -- ? ) 2 mod zero? ; 2 even? 3 even? 4 even?`)
	a := m.popValue().(Boolean)
	b := m.popValue().(Boolean)
	c := m.popValue().(Boolean)
	if a != Boolean(true) || b != Boolean(false) || c != Boolean(true) {
		t.Fail()
		t.Log(m.values)
	}
}

func TestRecursion(t *testing.T) {
	m := runCode(`: countDown ( n x -- x ) 1 - dup dup zero? [ ] [ countDown ] ? call ; 3 countDown`)
	a := m.popValue().(Integer)
	if a != Integer(0) {
		t.Fail()
		t.Log(m.values)
	}
}

func TestString(t *testing.T) {
	m := runCode(`"aa" " as " "bo" concat concat`)
	if m.popValue() != String("aa as bo") {
		t.Fail()
		t.Log(m.values)
	}
}

func TestStringConv(t *testing.T) {
	m := runCode(`1.34 >string`)
	if m.popValue() != String("1.34") {
		t.Fail()
		t.Log(m.values)
	}
}
