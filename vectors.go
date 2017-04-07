package main

import (
	"fmt"
)

var ModVectorCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "VectorCore",
	License:   "MIT",
	DocString: "TODO: Fill this in",
	Load:      loadVectorCore,
}

func loadVectorCore(m *machine) error {
	m.addGoWord("vec-set-at", " ( vec val idx -- vec ) ", GoWord(vecSetAt))
	m.addGoWord("vec-at", " ( vec idx -- elem ) ", GoWord(vecAt))
	m.addGoWord("<vector>", " ( -- vector )", GoWord(makeVec))
	m.addGoWord("vec-each", " ( vec quot -- .. ) ", GoWord(vecEach))
	m.addGoWord("vec-append", " ( vec elem -- newVect ) ", GoWord(vecAppend))
	m.addGoWord("vec-prepend", " ( vec elem -- newVect ) ", GoWord(vecPrepend))
	m.addGoWord("vec-popback", " ( vec -- vec elem ) ", GoWord(vecPopback))
	m.addGoWord("vec-popfront", " ( vec -- vec elem ) ", GoWord(vecPopfront))
	// Return success if we can load out PISC file as well.
	return m.importPISCAsset("stdlib/vectors.pisc")
}

func vecSetAt(m *machine) error {
	idx := int(m.popValue().(Integer))
	val := m.popValue()
	arr := m.values[len(m.values)-1].(Array)
	if idx > len(arr)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	arr[idx] = val
	return nil
}

func vecAt(m *machine) error {
	idx := int(m.popValue().(Integer))
	arr := m.popValue().(Array)
	if idx > len(arr)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	m.pushValue(arr[idx])
	return nil
}

func makeVec(m *machine) error {
	m.pushValue(Array(make([]stackEntry, 0)))
	return nil
}

func vecEach(m *machine) error {
	quot := m.popValue().(*quotation)
	vect := m.popValue().(Array)
	for _, elem := range vect {
		m.pushValue(elem)
		m.pushValue(quot)
		err := m.executeQuotation()
		if err != nil {
			fmt.Println("ERROR HAPPENED")
			return err
		}
	}
	return nil
}

func vecAppend(m *machine) error {
	toAppend := m.popValue()
	arr := m.popValue().(Array)
	arr = append(arr, toAppend)
	m.pushValue(arr)
	return nil
}

func vecPrepend(m *machine) error {
	toPrepend := m.popValue()
	arr := m.popValue().(Array)
	arr = append([]stackEntry{toPrepend}, arr...)
	m.pushValue(arr)
	return nil
}

func vecPopback(m *machine) error {
	arr := m.popValue().(Array)
	if len(arr) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	m.pushValue(arr)
	m.pushValue(val)
	return nil
}

func vecPopfront(m *machine) error {
	arr := m.popValue().(Array)
	if len(arr) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := arr[0]
	arr = arr[1:]
	m.pushValue(arr)
	m.pushValue(val)
	return nil
}
