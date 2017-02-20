package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type stackEntry interface {
	String() string
	Type() string
}

type lenable interface {
	Length() int
}

func runEq(m *machine) error {
	a := m.popValue()
	b := m.popValue()
	m.pushValue(Boolean(eq(a, b)))
	return nil
}

// Check if two arrays are
func arrayRefEq(a, b Array) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	return &a[0] == &b[0]
}

func mapRefEq(a, b Dict) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	aVal := reflect.ValueOf(a).Pointer()
	bVal := reflect.ValueOf(b).Pointer()
	return aVal == bVal
}

func eq(a, b stackEntry) bool {

	if a.Type() != b.Type() {
		return false
	}
	// Otherwise, the types are the same and we can compare them
	switch a.Type() {
	case "Boolean":
		return a.(Boolean) == b.(Boolean)
	case "String":
		return a.String() == b.String()
	case "Integer":
		return a.(Integer) == b.(Integer)
	case "Double":
		return a.(Double) == b.(Double)
	case "Symbol":
		return a.(Symbol) == b.(Symbol)
	case "Vector":
		return arrayRefEq(a.(Array), b.(Array))
	case "Dictionary":
		return mapRefEq(a.(Dict), b.(Dict))
	case "Quotation":
		return &a == &b
	case "Go Word":
		return &a == &b
	}
	// If we got here, something
	panic("eq failed!!!")

}

// Boolean is the PISC wrapper around bools
type Boolean bool

// Integer is the PISC wrapper around ints
type Integer int

// Double is the PISC wrapper aroudn float64
type Double float64

// Dict is the wrapper around dictionaries.
// TODO: Find a way to support dicts with arbitrary keys, not just strings
type Dict map[string]stackEntry

// Array is the wrapper around slices of stackEntry
type Array []stackEntry

// String is the PISC wrapper around strings
type String string

// Symbol is used for unique symboles
type Symbol int64

// This is a separate sematic from arrays.
type quotation struct {
	inner  *codeQuotation
	locals Dict
}

func (q quotation) toCode() *codeQuotation {
	q.inner.idx = 0
	return q.inner
}

type GoFunc GoWord

func (g GoFunc) String() string {
	return "<Native Code>"
}

func (s String) String() string {
	return string(s)
}

func (s Symbol) String() string {
	return "#" + fmt.Sprint(int(s))
}

func (b Boolean) String() string {
	if bool(b) {
		return "t"
	}
	return "f"
}

func (i Integer) String() string {
	return strconv.Itoa(int(i))
}

func (d Double) String() string {
	return fmt.Sprint(float64(d))
}

func (q quotation) String() string {
	return fmt.Sprint([]*word(q.inner.words))
}

func (dict Dict) String() string {
	var key_order stackEntry
	var found bool
	buf := bytes.NewBufferString("map[")
	if key_order, found = dict["__ordering"]; found {
		if keys, ok := key_order.(Array); ok {
			for _, k := range keys {
				buf.WriteString(fmt.Sprint(k.String(), ":", dict[k.String()], " "))
			}
			buf.WriteString("]")
			return buf.String()
		}
		return fmt.Sprint(map[string]stackEntry(dict))
	} else {
		return fmt.Sprint(map[string]stackEntry(dict))
	}

}

func (a Array) String() string {
	return fmt.Sprint([]stackEntry(a))
}

func (g GoFunc) Type() string {
	return "Go Word"
}

func (i Integer) Type() string {
	return "Integer"
}

func (d Double) Type() string {
	return "Double"
}

func (q quotation) Type() string {
	return "Quotation"
}

func (a Array) Type() string {
	return "Vector"
}

func (dict Dict) Type() string {
	return "Dictionary"
}

func (b Boolean) Type() string {
	return "Boolean"
}

func (s String) Type() string {
	return "String"
}

func (s Symbol) Type() string {
	return "Symbol"
}

func (dict Dict) Length() int {
	return len(dict)
}

func (a Array) Length() int {
	return len(a)
}

func (s String) Length() int {
	return len(s)
}
