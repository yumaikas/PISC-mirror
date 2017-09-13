package pisc

import (
	"bufio"
	//"regexp"
	"strconv"
	"strings"
)

var ModStringsCore = Module{
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

func _concat(m *Machine) error {
	a := m.PopValue().(String)
	b := m.PopValue().(String)
	m.PushValue(String(b + a))
	return nil
}

func _toString(m *Machine) error {
	a := m.PopValue()
	if s, ok := a.(String); ok {
		m.PushValue(String(s))
		return nil
	}
	m.PushValue(String(a.String()))
	return nil
}

func _strToInt(m *Machine) error {
	a := m.PopValue().(String).String()
	if i, err := strconv.Atoi(a); err != nil {
		return err
	} else {
		m.PushValue(Integer(i))
		return nil
	}
}

// Potential TODO: Review for performance?
func _strJoin(m *Machine) error {
	sep := m.PopValue().(String).String()
	elems := *m.PopValue().(Vector).Elements
	str := ""
	for _, val := range elems[:len(elems)-1] {
		str += val.String() + sep
	}
	str += elems[len(elems)-1].String()
	m.PushValue(String(str))
	return nil
}

func _strSplit(m *Machine) error {
	sep := m.PopValue().(String).String()
	str := m.PopValue().(String).String()
	strs := strings.Split(str, sep)
	toPush := make([]StackEntry, len(strs))
	for idx, val := range strs {
		toPush[idx] = String(val)
	}
	m.PushValue(Vector{Elements: &toPush})
	return nil
}

func _strEmpty(m *Machine) error {
	a := m.PopValue().String()
	m.PushValue(Boolean(len(a) <= 0))
	return nil
}

func _strEq(m *Machine) error {
	b := m.PopValue().String()
	a := m.PopValue().String()
	m.PushValue(Boolean(a == b))
	return nil
}

func _strToRuneReader(m *Machine) error {
	a := m.PopValue().(String)
	reader := strings.NewReader(string(a))
	bufReader := bufio.NewReader(reader)

	readerObj := MakeReader(bufReader)
	m.PushValue(Dict(readerObj))
	return nil
}

func _eachChar(m *Machine) error {
	quot := m.PopValue().(*Quotation).toCode()
	str := m.PopValue().(String).String()
	var err error
	for _, r := range str {
		m.PushValue(String(string(r)))
		quot.Idx = 0
		err = m.execute(quot)
		if err != nil {
			break
		}
	}
	return err
}

func _strReplace(m *Machine) error {
	replace := m.PopValue().String()
	pat := m.PopValue().String()
	str := m.PopValue().String()
	newstr := strings.Replace(str, pat, replace, -1)
	m.PushValue(String(newstr))
	return nil
}

func _strContains(m *Machine) error {
	substr := m.PopValue().String()
	str := m.PopValue().String()
	m.PushValue(Boolean(strings.Contains(str, substr)))
	return nil
}

func _strEndsWith(m *Machine) error {
	endStr := m.PopValue().String()
	str := m.PopValue().String()
	m.PushValue(Boolean(strings.HasSuffix(str, endStr)))
	return nil
}

func _strStartsWith(m *Machine) error {
	prefix := m.PopValue().String()
	str := m.PopValue().String()
	m.PushValue(Boolean(strings.HasPrefix(str, prefix)))
	return nil
}

func _strSubstr(m *Machine) error {
	end := m.PopValue().(Integer)
	start := m.PopValue().(Integer)
	str := m.PopValue().String()

	m.PushValue(String(str[start:end]))
	return nil
}

func _strIdxOf(m *Machine) error {
	substr := m.PopValue().String()
	str := m.PopValue().String()
	idx := strings.Index(str, substr)
	m.PushValue(Integer(idx))
	return nil
}

func _strRepeat(m *Machine) error {
	numRepeats := int(m.PopValue().(Integer))
	str := m.PopValue().String()
	m.PushValue(String(strings.Repeat(str, numRepeats)))
	return nil
}

func _strTrim(m *Machine) error {
	str := m.PopValue().String()
	str = strings.Trim(str, "\t\n\r ")
	m.PushValue(String(str))
	return nil
}

// ( str sep -- count )
func _strCount(m *Machine) error {
	sep := m.PopValue().String()
	str := m.PopValue().String()
	m.PushValue(Integer(strings.Count(str, sep)))
	return nil
}

// ( str -- upperstr )
func _strUpper(m *Machine) error {
	str := m.PopValue().String()
	m.PushValue(String(strings.ToUpper(str)))
	return nil
}

// ( str -- lowerstr )
func _strLower(m *Machine) error {
	str := m.PopValue().String()
	m.PushValue(String(strings.ToLower(str)))
	return nil
}
func loadStringCore(m *Machine) error {

	m.AddGoWord("str-concat", "( str-a str-b -- str-ab )", _concat)
	m.AddGoWord(">string", "( anyVal -- str )", _toString)
	m.AddGoWord("str>int", "( str -- int! )", _strToInt)
	m.AddGoWord("str-join", "( vec sep -- str )", _strJoin)

	m.AddGoWord("str-ends?", "( str endstr -- endswith? )", _strEndsWith)
	m.AddGoWord("str-eq?", "( a b -- eq? )", _strEq)
	m.AddGoWord("str-substr", "( str start end -- substr )", _strSubstr)
	m.AddGoWord("str-idx-of", "( str sub -- idx )", _strIdxOf)
	m.AddGoWord("str-split", "( str sep -- vec )", _strSplit)
	m.AddGoWord("str-starts?", "( str prefix -- startswith? )", _strStartsWith)

	m.AddGoWord("str-empty?", "( str -- empty? )", _strEmpty)
	m.AddGoWord("str>rune-reader", "( str -- reader )", _strToRuneReader)
	m.AddGoWord("each-char", "( str quot -- .. )", _eachChar)

	m.AddGoWord("str-replace", "( str pat replace -- .. )", _strReplace)
	m.AddGoWord("str-contains?", "( str cont -- contained? )", _strContains)
	m.AddGoWord("str-count", "( str sep -- count )", _strCount)

	m.AddGoWord("str-repeat", "( str repeat-count -- 'str )", _strRepeat)
	m.AddGoWord("str-trim", "( str -- 'str )", _strTrim)

	m.AddGoWord("str-upper", " ( str -- upper-str ) ", _strUpper)
	m.AddGoWord("str-lower", " ( str -- lower-str ) ", _strLower)

	return m.ImportPISCAsset("stdlib/strings.pisc")
}
