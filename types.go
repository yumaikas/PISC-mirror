package main

import (
	"fmt"
	"strconv"
)

type stackEntry interface {
	String() string
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
	return fmt.Sprint(q)
}

func (dict Dict) String() string {
	return fmt.Sprint(dict)
}

func (a Array) String() string {
	return fmt.Sprint(a)
}
