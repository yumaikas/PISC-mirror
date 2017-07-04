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
	m.AddGoWord("vec-append", " ( vec elem -- newVect ) ", GoWord(vecAppend))
	m.AddGoWord("vec-prepend", " ( vec elem -- newVect ) ", GoWord(vecPrepend))
	m.AddGoWord("vec-popback", " ( vec -- vec elem ) ", GoWord(vecPopback))
	m.AddGoWord("vec-popfront", " ( vec -- vec elem ) ", GoWord(vecPopfront))
	// Return success if we can load out PISC file as well.
	return m.importPISCAsset("stdlib/vectors.pisc")
}

func vecSetAt(m *Machine) error {
	idx := int(m.PopValue().(Integer))
	val := m.PopValue()
	arr := m.Values[len(m.Values)-1].(Array)
	if idx > len(arr)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	arr[idx] = val
	return nil
}

func vecAt(m *Machine) error {
	idx := int(m.PopValue().(Integer))
	arr := m.PopValue().(Array)
	if idx > len(arr)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	m.PushValue(arr[idx])
	return nil
}

func makeVec(m *Machine) error {
	m.PushValue(Array(make([]StackEntry, 0)))
	return nil
}

func vecEach(m *Machine) error {
	quot := m.PopValue().(*Quotation)
	vect := m.PopValue().(Array)
	for _, elem := range vect {
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

func vecAppend(m *Machine) error {
	toAppend := m.PopValue()
	arr := m.PopValue().(Array)
	arr = append(arr, toAppend)
	m.PushValue(arr)
	return nil
}

func vecPrepend(m *Machine) error {
	toPrepend := m.PopValue()
	arr := m.PopValue().(Array)
	arr = append([]StackEntry{toPrepend}, arr...)
	m.PushValue(arr)
	return nil
}

func vecPopback(m *Machine) error {
	arr := m.PopValue().(Array)
	if len(arr) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	m.PushValue(arr)
	m.PushValue(val)
	return nil
}

func vecPopfront(m *Machine) error {
	arr := m.PopValue().(Array)
	if len(arr) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := arr[0]
	arr = arr[1:]
	m.PushValue(arr)
	m.PushValue(val)
	return nil
}
