package main

import (
	"fmt"
)

func (m *machine) loadVectorWords() error {

	m.predefinedWords["vec-at"] = GoWord(func(m *machine) error {
		// ( vec idx -- elem )
		idx := int(m.popValue().(Integer))
		arr := m.popValue().(Array)
		if idx > len(arr)-1 || idx < 0 {
			return fmt.Errorf("Index out of bounds!")
		}
		m.pushValue(arr[idx])
		return nil
	})

	// ( -- vector )
	m.predefinedWords["<vector>"] = NilWord(func(m *machine) {
		m.pushValue(Array(make([]stackEntry, 0)))
	})
	// ( vec quot -- .. )
	m.predefinedWords["vec-each"] = GoWord(func(m *machine) error {
		quot := m.popValue().(quotation).toCode()
		vect := m.popValue().(Array)
		for _, elem := range vect {
			m.pushValue(elem)
			quot.idx = 0
			err := executeWordsOnMachine(m, quot)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// ( vec elem -- newVect )
	m.predefinedWords["vec-append"] = NilWord(func(m *machine) {
		toAppend := m.popValue()
		arr := m.popValue().(Array)
		arr = append(arr, toAppend)
		m.pushValue(arr)
	})
	// ( vec elem -- newVect )
	m.predefinedWords["vec-prepend"] = NilWord(func(m *machine) {
		toPrepend := m.popValue()
		arr := m.popValue().(Array)
		arr = append([]stackEntry{toPrepend}, arr...)
		m.pushValue(arr)
	})

	// ( vec -- vec elem )
	m.predefinedWords["vec-popback"] = GoWord(func(m *machine) error {
		arr := m.popValue().(Array)
		if len(arr) < 1 {
			return fmt.Errorf("No elements in array!")
		}
		val := arr[len(arr)-1]
		arr = arr[:len(arr)-1]
		m.pushValue(arr)
		m.pushValue(val)
		return nil
	})
	// ( vec -- vec elem)
	m.predefinedWords["vec-popfront"] = GoWord(func(m *machine) error {
		arr := m.popValue().(Array)
		if len(arr) < 1 {
			return fmt.Errorf("No elements in array!")
		}
		val := arr[0]
		arr = arr[1:]
		m.pushValue(arr)
		m.pushValue(val)
		return nil
	})
	return nil
}
