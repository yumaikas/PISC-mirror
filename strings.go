package main

import (
	"bufio"
	//"regexp"
	"strconv"
	"strings"
)

var ModStringsCore = PISCModule{
	Author:    "Andrew Owen",
	Name:      "StringCore",
	License:   "MIT",
	DocString: "TODO: Fill this in",
	Load:      loadStringCore,
}

// TODO: Implement Lua Patterns
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

func _strToInt(m *machine) error {
	a := m.popValue().(String).String()
	if i, err := strconv.Atoi(a); err != nil {
		return err
	} else {
		m.pushValue(Integer(i))
		return nil
	}
}

// Potential TODO: Review for performance?
func _strJoin(m *machine) error {
	sep := m.popValue().(String).String()
	elems := m.popValue().(Array)
	str := ""
	for _, val := range elems[:len(elems)-1] {
		str += val.String() + sep
	}
	str += elems[len(elems)-1].String()
	m.pushValue(String(str))
	return nil
}

func _strSplit(m *machine) error {
	sep := m.popValue().(String).String()
	str := m.popValue().(String).String()
	strs := strings.Split(str, sep)
	toPush := make(Array, len(strs))
	for idx, val := range strs {
		toPush[idx] = String(val)
	}
	m.pushValue(toPush)
	return nil
}

func _strEmpty(m *machine) error {
	a := m.popValue().String()
	m.pushValue(Boolean(len(a) <= 0))
	return nil
}

func _strEq(m *machine) error {
	b := m.popValue().String()
	a := m.popValue().String()
	m.pushValue(Boolean(a == b))
	return nil
}

func _strToRuneReader(m *machine) error {
	a := m.popValue().(String)
	reader := strings.NewReader(string(a))
	bufReader := bufio.NewReader(reader)

	readerObj := makeReader(bufReader)
	m.pushValue(Dict(readerObj))
	return nil
}

func _eachChar(m *machine) error {
	quot := m.popValue().(*quotation).toCode()
	str := m.popValue().(String).String()
	var err error
	for _, r := range str {
		m.pushValue(String(string(r)))
		quot.idx = 0
		err = m.execute(quot)
		if err != nil {
			break
		}
	}
	return err
}

func _strReplace(m *machine) error {
	replace := m.popValue().String()
	pat := m.popValue().String()
	str := m.popValue().String()
	newstr := strings.Replace(str, pat, replace, -1)
	m.pushValue(String(newstr))
	return nil
}

func _strContains(m *machine) error {
	substr := m.popValue().String()
	str := m.popValue().String()
	m.pushValue(Boolean(strings.Contains(str, substr)))
	return nil
}

func _strEndsWith(m *machine) error {
	endStr := m.popValue().String()
	str := m.popValue().String()
	m.pushValue(Boolean(strings.HasSuffix(str, endStr)))
	return nil
}

func _strStartsWith(m *machine) error {
	prefix := m.popValue().String()
	str := m.popValue().String()
	m.pushValue(Boolean(strings.HasPrefix(str, prefix)))
	return nil
}

func _strSubstr(m *machine) error {
	end := m.popValue().(Integer)
	start := m.popValue().(Integer)
	str := m.popValue().String()

	m.pushValue(String(str[start:end]))
	return nil
}

func _strIdxOf(m *machine) error {
	substr := m.popValue().String()
	str := m.popValue().String()
	idx := strings.Index(str, substr)
	m.pushValue(Integer(idx))
	return nil
}

func _strRepeat(m *machine) error {
	numRepeats := int(m.popValue().(Integer))
	str := m.popValue().String()
	m.pushValue(String(strings.Repeat(str, numRepeats)))
	return nil
}

func loadStringCore(m *machine) error {

	m.addGoWord("str-concat", "( str-a str-b -- str-ab )", _concat)
	m.addGoWord(">string", "( anyVal -- str )", _toString)
	m.addGoWord("str>int", "( str -- int! )", _strToInt)
	m.addGoWord("str-join", "( vec sep -- str )", _strJoin)

	m.addGoWord("str-ends?", "( str endstr -- endswith? )", _strEndsWith)
	m.addGoWord("str-eq?", "( a b -- eq? )", _strEq)
	m.addGoWord("str-substr", "( str start end -- substr )", _strSubstr)
	m.addGoWord("str-idx-of", "( str sub -- idx )", _strIdxOf)
	m.addGoWord("str-split", "( str sep -- vec )", _strSplit)
	m.addGoWord("str-starts?", "( str prefix -- startswith? )", _strStartsWith)

	m.addGoWord("str-empty?", "( str -- empty? )", _strEmpty)
	m.addGoWord("str>rune-reader", "( str -- reader )", _strToRuneReader)
	m.addGoWord("each-char", "( str quot -- .. )", _eachChar)

	m.addGoWord("str-replace", "( str pat replace -- .. )", _strReplace)
	m.addGoWord("str-contains?", "( str cont -- contained? )", _strContains)

	m.addGoWord("str-repeat", "( str repeat-count -- 'str )", _strRepeat)

	return m.importPISCAsset("stdlib/strings.pisc")
}
