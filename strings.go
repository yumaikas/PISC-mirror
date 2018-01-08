package pisc

import (
	"bufio"
	"unicode"
	"unicode/utf8"
	//"regexp"
	//"fmt"
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

var emptyStr = String("")

// Potential TODO: Review for performance?
func _strJoin(m *Machine) error {
	sep := m.PopValue().(String).String()
	elems := m.PopValue().(*Vector).Elements
	str := ""
	if len(elems) == 0 {
		m.PushValue(emptyStr)
		return nil
	}
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
	m.PushValue(&Vector{Elements: toPush})
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

// ( str -- reversed-str )
func _strReverse(m *Machine) error {
	str := m.PopValue().String()

	// Code from: http://rosettacode.org/wiki/Reverse_a_string#Go
	r := make([]rune, len(str))
	start := len(str)
	for _, c := range str {
		if c != utf8.RuneError {
			start--
			r[start] = c
		}
	}
	m.PushValue(String(string(r[start:])))
	return nil
}

// ( str -- reversed-bytes )
func _strReverseBytes(m *Machine) error {
	str := m.PopValue().String()

	// Code inspired by: http://rosettacode.org/wiki/Reverse_a_string#Go
	r := make([]byte, len(str))
	for i := 0; i < len(str); i++ {
		r[i] = str[len(str)-1-i]
	}

	m.PushValue(String(string(r)))
	return nil
}

// ( str -- reversed-Graphemes )
func _strReverseGraphemes(m *Machine) error {
	str := m.PopValue().String()

	// Code inspired by: http://rosettacode.org/wiki/Reverse_a_string#Go
	if str == "" {
		m.PushValue(String(""))
		return nil
	}

	source := []rune(str)
	result := make([]rune, len(source))
	start := len(result)

	for i := 0; i < len(source); {
		// Continue past invalid UTF-8
		if source[i] == utf8.RuneError {
			i++
			continue
		}

		j := i + 1
		for j < len(source) && (unicode.Is(unicode.Mn, source[j]) ||
			unicode.Is(unicode.Me, source[j]) ||
			unicode.Is(unicode.Mc, source[j])) {
			j++
		}

		for k := j - 1; k >= i; k-- {
			start--
			result[start] = source[k]
		}
		i = j
	}
	m.PushValue(String(string(result[start:])))
	return nil
}

func loadStringCore(m *Machine) error {

	m.AddGoWord("str-concat", "( str-a str-b -- str-ab )", _concat)
	m.AddGoWordWithStack(
		"str-concat",
		"( str-a str-b -- str-ab )",
		"Concatenate str-a and str-b, allocating a new string",
		_concat)
	m.AddGoWordWithStack(">string",
		"( anyVal -- str )",
		"Cast the top value of the stack to a string.",
		_toString)
	m.AddGoWordWithStack("str>int",
		"( str -- int! )",
		"Attempt to parse a string from the top of the stack into an int.",
		_strToInt)
	m.AddGoWordWithStack(
		"str-join",
		"( vec sep -- str )",
		"Casts each element of vec to a string, and then joins them with sep in the middle",
		_strJoin)

	m.AddGoWordWithStack(
		"str-ends?",
		"( str endstr -- endswith? )",
		"Checks to see if str ends with endstr",
		_strEndsWith)
	m.AddGoWordWithStack(
		"str-eq?",
		"( a b -- eq? )",
		"Compares two strings for equality.",
		_strEq)
	m.AddGoWordWithStack(
		"str-substr",
		"( str start end -- substr )",
		"Takes a substring of str over the range [start:end]",
		_strSubstr)
	m.AddGoWordWithStack("str-idx-of",
		"( str sub -- idx )",
		"Find the index of the first instance of sub in str",
		_strIdxOf)
	m.AddGoWordWithStack("str-split",
		"( str sep -- vec )",
		"Split str into a vector of strings, based on sep",
		_strSplit)
	m.AddGoWordWithStack("str-starts?",
		"( str prefix -- startswith? )",
		"Check to see if str starts with prefix.",
		_strStartsWith)

	m.AddGoWordWithStack("str-empty?",
		"( str -- empty? )",
		"Returns true if the string has a length of 0 or less",
		_strEmpty)
	m.AddGoWordWithStack("str>rune-reader",
		"( str -- reader )",
		"Creates a rune-reader from str",
		_strToRuneReader)
	m.AddGoWordWithStack("each-char",
		"( str quot -- .. )",
		"Executes quot for each character in str",
		_eachChar)

	m.AddGoWordWithStack("str-replace",
		"( str pat replace -- .. )",
		"Returns a new string based on str with `pat` swapped for `replace`",
		_strReplace)
	m.AddGoWordWithStack("str-contains?",
		"( str cont -- contained? )",
		"Returns true if `cont` is in str, false otherwise",
		_strContains)
	m.AddGoWordWithStack("str-count",
		"( str sep -- count )",
		"Counts the number of instances of sep in str",
		_strCount)

	m.AddGoWordWithStack("str-repeat",
		"( str repeat-count -- 'str )",
		"Returns a new string that is str repeated `repeat-count` times",
		_strRepeat)
	m.AddGoWordWithStack("str-trim",
		"( str -- 'str )",
		"Trims whitespace off the front and back of `str`",
		_strTrim)

	m.AddGoWordWithStack("str-upper",
		" ( str -- upper-str ) ",
		"Returns an upper-cased copy of str",
		_strUpper)
	m.AddGoWordWithStack("str-lower",
		" ( str -- lower-str ) ",
		"Returns an lower-cased copy of str",
		_strLower)

	m.AddGoWordWithStack("str-reverse",
		" ( str -- reversed-runes ) ",
		"Splits a string into runes, and then reverses the runes into a new string",
		_strReverse)
	m.AddGoWordWithStack("str-reverse-bytes",
		" ( str -- reversed-bytes ) ",
		"Splits a string into bytes, and then reverses the bytes into a new string",
		_strReverseBytes)
	m.AddGoWordWithStack("str-reverse-graphemes",
		" ( str -- reversed-graphemes ) ",
		"Reverses str by grapheme clusters",
		_strReverseGraphemes)

	return m.ImportPISCAsset("stdlib/strings.pisc")
}
