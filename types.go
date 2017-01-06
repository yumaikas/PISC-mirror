package main

import (
	"bytes"
	"fmt"
	"strconv"
)

type stackEntry interface {
	String() string
	Type() string
}

type lenable interface {
	Length() int
}

type Boolean bool
type Integer int
type Double float64
type Dict map[string]stackEntry
type Array []stackEntry
type String string
type Symbol int64

// This is a separate sematic from arrays.
type quotation struct {
	code   []*word
	locals Dict
	codePosition
}

func (q quotation) toCode() *codeQuotation {
	return &codeQuotation{
		idx:          0,
		words:        q.code,
		codePosition: q.codePosition,
	}
}

type GoFunc GoWord

func (g GoFunc) String() string {
	return "<Native Code>"
}

func (g GoFunc) Type() string {
	return "Go Word"
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
	return fmt.Sprint([]*word(q.code))
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
	return "ERROR: Symbol observed on stack!"
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
