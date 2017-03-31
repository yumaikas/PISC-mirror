package main

import (
	"bufio"
	//"regexp"
	"strconv"
	"strings"
)

type stringPattern struct {
}

// TODO: Start replacing anonymous functions with named functions so that profiling is a little clearer

// TODO: Add more words to support strings here, we need a way to handle a lot more cases, like
// replacement, substringing, joining and so on.

func _concat(m *machine) error {
	a := m.popValue().(String)
	b := m.popValue().(String)
	m.pushValue(String(b + a))
	return nil
}

func _toString(m *machine) error {
	a := m.popValue()
	if s, ok := a.(String); ok {
		m.pushValue(String(s))
		return nil
	}
	m.pushValue(String(a.String()))
	return nil
}

func (m *machine) loadStringWords() error {
	m.predefinedWords["concat"] = _concat
	m.predefinedWords[">string"] = _toString

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

	m.predefinedWords["str-empty?"] = NilWord(func(m *machine) {
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
		quot := m.popValue().(*quotation).toCode()
		str := m.popValue().(String).String()
		for _, r := range str {
			m.pushValue(String(string(r)))
			quot.idx = 0
			m.execute(quot)
		}
	})

	// ( str pat replace -- replaced )
	m.predefinedWords["str-replace"] = NilWord(func(m *machine) {
		replace := m.popValue().String()
		pat := m.popValue().String()
		str := m.popValue().String()
		newstr := strings.Replace(str, pat, replace, -1)
		m.pushValue(String(newstr))
	})

	// ( str cont -- ? )
	m.predefinedWords["str-contains"] = NilWord(func(m *machine) {
		substr := m.popValue().String()
		str := m.popValue().String()
		if strings.Contains(str, substr) {
			m.pushValue(Boolean(true))
			return
		}
		m.pushValue(Boolean(false))
	})

	m.addGoWord("str-substr", "( str start end -- substr )", func(m *machine) error {
		end := m.popValue().(Integer)
		start := m.popValue().(Integer)
		str := m.popValue().String()

		m.pushValue(String(str[start:end]))
		return nil
	})
	m.addGoWord("str-idx-of", "( str sub -- idx )", func(m *machine) error {
		substr := m.popValue().String()
		str := m.popValue().String()
		idx := strings.Index(str, substr)
		m.pushValue(Integer(idx))
		return nil
	})

	return nil
}
