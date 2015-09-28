package main

import "testing"

func TestParser(t *testing.T) {
	return
	m := runCode(`1 2 dup`)

	t.Fail()
	t.Log(m)

}

func TestExecution(t *testing.T) {
	m := runCode(`t [ 1 ] [ 2 ] ? call
 f [ 1 ] [ 2 ] ? call`)
	t.Fail()
	t.Log(m.values)
}

func TestCall(t *testing.T) {
	m := runCode(`[ 1 ] call`)
	t.Fail()
	t.Log(m.values)
}
