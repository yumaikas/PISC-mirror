package pisc

import (
	"fmt"
)

var ModVectorCore = Module{
	Author:    "Andrew Owen",
	Name:      "VectorCore",
	License:   "MIT",
	DocString: "TODO: Fill this in",
	Load:      loadVectorCore,
}

func loadVectorCore(m *Machine) error {
	m.AddGoWord("vec-set-at", " ( vec val idx -- vec ) ", GoWord(vecSetAt))
	m.AddGoWord("vec-at", " ( vec idx -- elem ) ", GoWord(vecAt))
	m.AddGoWord("<vector>", " ( -- vector )", GoWord(makeVec))
	m.AddGoWord("vec-each", " ( vec quot -- .. ) ", GoWord(vecEach))
	m.AddGoWord("vec-append", " ( vec elem -- vec ) ", GoWord(vecAppend))
	m.AddGoWord("vec-push", ` ( vec elem -- )
		Mutates vector to contain elem`, GoWord(vecPush))
	m.AddGoWord("vec-pushback", " ( vec elem -- ) ", GoWord(vecPush))
	m.AddGoWord("vec-prepend", " ( vec elem -- vec ) ", GoWord(vecPrepend))
	m.AddGoWord("vec-popback", " ( vec -- vec elem ) ", GoWord(vecPopback))
	m.AddGoWord("vec-popfront", " ( vec -- vec elem ) ", GoWord(vecPopfront))
	// Return success if we can load out PISC file as well.
	return m.ImportPISCAsset("stdlib/vectors.pisc")
}

func vecSetAt(m *Machine) error {
	idx := int(m.PopValue().(Integer))
	val := m.PopValue()
	elems := m.Values[len(m.Values)-1].(*Vector).Elements
	if idx > len(elems)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	elems[idx] = val
	return nil
}

func vecAt(m *Machine) error {
	idx := int(m.PopValue().(Integer))
	elems := m.PopValue().(*Vector).Elements
	if idx > len(elems)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	m.PushValue(elems[idx])
	return nil
}

func makeVec(m *Machine) error {
	vec := &Vector{Elements: make([]StackEntry, 0)}
	m.PushValue(vec)
	return nil
}

func vecEach(m *Machine) error {
	quot := m.PopValue().(*Quotation)
	vect := m.PopValue().(*Vector)
	for _, elem := range vect.Elements {
		m.PushValue(elem)
		m.PushValue(quot)
		err := m.ExecuteQuotation()
		if err != nil {
			fmt.Println("ERROR HAPPENED")
			return err
		}
	}
	return nil
}

func vecPush(m *Machine) error {
	toAppend := m.PopValue()
	vec := m.PopValue().(*Vector)
	newElems := append(vec.Elements, toAppend)
	vec.Elements = newElems
	return nil
}

func vecAppend(m *Machine) error {
	toAppend := m.PopValue()
	vec := m.PopValue().(*Vector)
	newElems := append(vec.Elements, toAppend)
	vec.Elements = newElems
	m.PushValue(vec)
	return nil
}

func vecPrepend(m *Machine) error {
	toPrepend := m.PopValue()
	arr := m.PopValue().(*Vector)
	newElems := append([]StackEntry{toPrepend}, (arr.Elements)...)
	arr.Elements = newElems
	m.PushValue(arr)
	return nil
}

func vecPopback(m *Machine) error {
	vec := m.PopValue().(*Vector)
	if len(vec.Elements) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := (vec.Elements)[len(vec.Elements)-1]
	vec.Elements = vec.Elements[:len(vec.Elements)-1]
	m.PushValue(vec)
	m.PushValue(val)
	return nil
}

func vecPopfront(m *Machine) error {
	vec := m.PopValue().(*Vector)
	if len(vec.Elements) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := vec.Elements[0]
	vec.Elements = vec.Elements[1:]
	m.PushValue(vec)
	m.PushValue(val)
	return nil
}
