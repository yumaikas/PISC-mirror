package main

import (
	"fmt"
)

func (m *machine) loadVectorWords() error {

	m.addGoWord("vec-set-at",
		" ( vec val idx -- vec ) ",
		GoWord(func(m *machine) error {
			idx := int(m.popValue().(Integer))
			val := m.popValue()
			arr := m.values[len(m.values)-1].(Array)
			if idx > len(arr)-1 || idx < 0 {
				return fmt.Errorf("index out of bounds: %v", idx)
			}
			arr[idx] = val
			return nil
		}))

	m.addGoWord("vec-at",
		" ( vec idx -- elem ) ",
		GoWord(func(m *machine) error {
			idx := int(m.popValue().(Integer))
			arr := m.popValue().(Array)
			if idx > len(arr)-1 || idx < 0 {
				return fmt.Errorf("index out of bounds: %v", idx)
			}
			m.pushValue(arr[idx])
			return nil
		}))

	// ( -- vector )
	m.addGoWord("<vector>",
		" ( -- vector )",
		GoWord(func(m *machine) error {
			m.pushValue(Array(make([]stackEntry, 0)))
			return nil
		}))

	m.addGoWord("vec-each",
		" ( vec quot -- .. ) ",
		GoWord(func(m *machine) error {
			quot := m.popValue().(*quotation).toCode()
			vect := m.popValue().(Array)
			for _, elem := range vect {
				m.pushValue(elem)
				quot.idx = 0
				err := m.execute(quot)
				if err != nil {
					return err
				}
			}
			return nil
		}))

	m.addGoWord("vec-append",
		" ( vec elem -- newVect ) ",
		GoWord(func(m *machine) error {
			toAppend := m.popValue()
			arr := m.popValue().(Array)
			arr = append(arr, toAppend)
			m.pushValue(arr)
			return nil
		}))
	m.addGoWord("vec-prepend",
		" ( vec elem -- newVect ) ",
		GoWord(func(m *machine) error {
			toPrepend := m.popValue()
			arr := m.popValue().(Array)
			arr = append([]stackEntry{toPrepend}, arr...)
			m.pushValue(arr)
			return nil
		}))

	m.addGoWord("vec-popback",
		" ( vec -- vec elem ) ",
		GoWord(func(m *machine) error {
			arr := m.popValue().(Array)
			if len(arr) < 1 {
				return fmt.Errorf("no elements in array")
			}
			val := arr[len(arr)-1]
			arr = arr[:len(arr)-1]
			m.pushValue(arr)
			m.pushValue(val)
			return nil
		}))

	m.addGoWord("vec-popfront",
		" ( vec -- vec elem ) ",
		GoWord(func(m *machine) error {
			arr := m.popValue().(Array)
			if len(arr) < 1 {
				return fmt.Errorf("no elements in array")
			}
			val := arr[0]
			arr = arr[1:]
			m.pushValue(arr)
			m.pushValue(val)
			return nil
		}))
	return nil
}
