package main

type stackEntry interface{}

type number interface {
	Add(number) number
}

type Boolean bool
type Integer int
type Double float64
type Dict map[string]stackEntry
type Array []stackEntry
type String string

// This is a separate sematic from arrays.
type quotation []word

// TODO: come back and implement this for Double and Integer
func (d Double) Add(n number) number {
	switch n.(type) {
	case Double:
		return number(d + n.(Double))
	case Integer:
		return number(Double(float64(d) + float64(int(n.(Integer)))))
	default:
		panic("Number addition type error!")
	}
}

func (i Integer) Add(n number) number {
	switch n.(type) {
	case Double:
		return number(Double(i) + n.(Double))
	case Integer:
		return number(Integer(float64(i) + float64(int(n.(Integer)))))
	default:
		panic("Number addition type error!")
	}
}
