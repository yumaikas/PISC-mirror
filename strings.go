package main

import (
	"bufio"
	"strconv"
	"strings"
)

// TODO: Add more words to support strings here, we need a way to handle a lot more cases, like
// replacement, substringing, joining and so on.

func (m *machine) loadStringWords() error {
	m.predefinedWords["concat"] = NilWord(func(m *machine) {
		a := m.popValue().(String)
		b := m.popValue().(String)
		m.pushValue(String(b + a))
	})
	m.predefinedWords[">string"] = NilWord(func(m *machine) {
		a := m.popValue()
		if s, ok := a.(String); ok {
			m.pushValue(String(s))
			return
		}
		m.pushValue(String(a.String()))
	})

	m.predefinedWords["string>int"] = GoWord(func(m *machine) error {
		a := m.popValue().(String).String()
		if i, err := strconv.Atoi(a); err != nil {
			return err
		} else {
			m.pushValue(Integer(i))
			return nil
		}

	})

	// ( vec sep -- str )
	m.predefinedWords["str-join"] = NilWord(func(m *machine) {
		sep := m.popValue().(String).String()
		elems := m.popValue().(Array)
		str := ""
		for _, val := range elems[:len(elems)-1] {
			str += val.String() + sep
		}
		str += elems[len(elems)-1].String()
		m.pushValue(String(str))
	})

	// ( str sep -- vec )
	m.predefinedWords["str-split"] = NilWord(func(m *machine) {
		sep := m.popValue().(String).String()
		str := m.popValue().(String).String()
		strs := strings.Split(str, sep)
		toPush := make(Array, len(strs))
		for idx, val := range strs {
			toPush[idx] = String(val)
		}
		m.pushValue(toPush)
	})

	m.predefinedWords["empty?"] = NilWord(func(m *machine) {
		a := m.popValue().(String)
		if len(a) > 0 {
			m.pushValue(Boolean(len(a) > 0))
		}
	})

	m.predefinedWords["str-eq"] = NilWord(func(m *machine) {
		b := m.popValue().String()
		a := m.popValue().String()
		m.pushValue(Boolean(a == b))

	})
	// ( str -- obj )
	m.predefinedWords["str>rune-reader"] = NilWord(func(m *machine) {
		a := m.popValue().(String)
		reader := strings.NewReader(string(a))
		bufReader := bufio.NewReader(reader)

		readerObj := makeReader(bufReader)
		m.pushValue(Dict(readerObj))
	})
	// ( str quot -- .. )
	m.predefinedWords["each-char"] = NilWord(func(m *machine) {
		quot := m.popValue().(quotation)
		str := m.popValue().(String).String()
		code := &codeList{idx: 0, code: quot}
		for _, r := range str {
			m.pushValue(String(string(r)))
			code.idx = 0
			executeWordsOnMachine(m, code)
		}
	})

	return nil
}
