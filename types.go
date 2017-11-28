package pisc

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type StackEntry interface {
	String() string
	Type() string
}

type Lenable interface {
	Length() int
}

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}

func runEq(m *Machine) error {
	a := m.PopValue()
	b := m.PopValue()
	m.PushValue(Boolean(eq(a, b)))
	return nil
}

// Check if two arrays are referentially equal
func vectorRefEq(a, b *Vector) bool {
	if len(a.Elements) != len(b.Elements) {
		return false
	}
	if len(a.Elements) == 0 {
		return true
	}
	return &(a.Elements)[0] == &(b.Elements)[0]
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

func eq(a, b StackEntry) bool {

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
		return vectorRefEq(a.(*Vector), b.(*Vector))
	case "Dictionary":
		return mapRefEq(a.(Dict), b.(Dict))
	case "Quotation":
		return &a == &b
	case "Go Word":
		return &a == &b
	}
	// If we got here, something is borked
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
type Dict map[string]StackEntry

// Vectors are mutable pointers to slices.
type Vector struct {
	Elements []StackEntry
}

// String is the PISC wrapper around strings
type String string

// Symbol is used for unique symboles
type Symbol int64

// This is a separate sematic from arrays.
type Quotation struct {
	inner  *CodeQuotation
	locals Dict
}

func (q Quotation) toCode() *CodeQuotation {
	q.inner.Idx = 0
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

func (q *Quotation) String() string {
	return fmt.Sprint([]*word(q.inner.Words))
}

func (dict Dict) String() string {
	var key_order StackEntry
	// var strMethod StackEntry
	var found bool
	/* TODO: Figure this out later
	if strMethod, found = dict["tostring"]; found {
		if toCall, ok := strMethod.(*CodeQuotation) {
		}
	}
	*/
	buf := bytes.NewBufferString("map[")
	if key_order, found = dict["__ordering"]; found {
		if keys, ok := key_order.(*Vector); ok {
			for _, k := range keys.Elements {
				buf.WriteString(fmt.Sprint(k.String(), ":", dict[k.String()], " "))
			}
			buf.WriteString("]")
			return buf.String()
		}
		return fmt.Sprint(map[string]StackEntry(dict))
	} else {
		return fmt.Sprint(map[string]StackEntry(dict))
	}

}

func (v *Vector) String() string {
	elems := make([]string, len(v.Elements))
	for i, e := range v.Elements {
		elems[i] = e.String()
	}
	return "{ " + strings.Join(elems, " ") + " }"
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

func (q *Quotation) Type() string {
	return "Quotation"
}

func (v *Vector) Type() string {
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

func (v *Vector) Length() int {
	return len(v.Elements)
}

func (s String) Length() int {
	return len(s)
}
