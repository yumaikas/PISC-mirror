package main

import (
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

// This is a separate sematic from arrays.
type quotation []word

func (s String) String() string {
	return "<<" + string(s) + ">>"
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
	return fmt.Sprint([]word(q))
}

func (dict Dict) String() string {
	return fmt.Sprint(map[string]stackEntry(dict))
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

func (dict Dict) Length() int {
	return len(dict)
}

func (a Array) Length() int {
	return len(a)
}
