package main

import (
	// "bufio"
	"fmt"
	"io/ioutil"
	// "os"
)

func isIOWord(w word) bool {
	return w == "import" || w == "priv_puts" || w == "in_line"
}
func (m *machine) executeIOWord(w word) error {
	switch w {
	case "import":
		fileName := m.popValue().(String)
		data, err := ioutil.ReadFile(string(fileName))
		if err != nil {
			return err
		}
		words, spaces := getWordList(string(data))
		p := &codeList{
			idx:    0,
			code:   words,
			spaces: spaces,
		}
		err = executeWordsOnMachine(m, p)
		if err != nil {
			return err
		}
	case "in_line":

	case "priv_puts":
		data := m.popValue().(String)
		fmt.Println(string(data))
	}
	return nil
}

// TODO: add words for printing and such here.
/*
func isStringWord(w word) bool {
	return w == "" ||
		w == "or" ||
		w == "not"
}

func (m *machine) executeStringWord(w word) error {
	switch w {
	case "and":
		a := m.popValue().(Boolean)
		b := m.popValue().(Boolean)
		m.pushValue(Boolean(a && b))
	case "or":
		a := m.popValue().(Boolean)
		b := m.popValue().(Boolean)
		m.pushValue(Boolean(a || b))
	case "not":
		a := m.popValue().(Boolean)
		m.pushValue(Boolean(!a))
	}
	return nil
}
*/
