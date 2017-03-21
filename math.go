package main

import (
	"fmt"
	"math"
)

var MathCore = PISCModule{
	Author:    "Andreb Owen",
	Name:      "MathCore",
	License:   "MIT",
	DocString: `Basic math stuff.`,
	Load:      loadMathCore,
}

func loadMathCore(m *machine) error {
	return m.loadHigherMathWords()
}

type number interface {
	stackEntry
	Add(number) number
	Negate() number
	Multiply(number) number
	Divide(number) number
	LessThan(number) Boolean
}

func (m *machine) loadHigherMathWords() error {

	m.addGoWord("+", "( a b -- c )  addition", executeAdd)
	m.addGoWord("-", " ( a b -- c ) subtraction ", executeSubtract)
	m.addGoWord("*", " ( a b -- c ) multiplication ", executeMultiply)
	m.addGoWord("/", " ( a b -- c ) division ", executeDivide)
	m.addGoWord("mod", " ( a b -- c ) modulus ", executeModulus)
	m.addGoWord("<", " ( a b -- c ) numeric less than ", executeLessThan)
	m.addGoWord("zero?", " ( a -- isZero? ) returns if a is zero or not ", isZeroPred)

	// For now, PISC words are late-bound, so we can get away with this.
	err := m.importPISCAsset("stdlib/math.pisc")
	if err != nil {
		return err
	}

	// Why do we duplicate the work here?
	// Because we want both the >double word and the
	// math1arg words
	m.predefinedWords[">double"] = GoWord(func(m *machine) error {
		a := m.popValue()
		if i, ok := a.(Integer); ok {
			m.pushValue(Double(float64(i)))
			return nil
		}
		if d, ok := a.(Double); ok {
			m.pushValue(Double(d))
			return nil
		}
		return fmt.Errorf("Unexpected type %s cannot be coverted to double", a.Type())
	})

	var math1Arg = func(name string, mathyFunc func(float64) float64) {
		m.predefinedWords[name] = GoWord(func(m *machine) error {
			a := m.popValue()
			if i, ok := a.(Integer); ok {
				m.pushValue(Double(mathyFunc(float64(i))))
				return nil
			}
			if d, ok := a.(Double); ok {
				m.pushValue(Double(mathyFunc(float64(d))))
				return nil
			}
			return fmt.Errorf("Unexpected type %s cannot be coverted to double", a.Type())
		})
	}

	math1Arg("acos", math.Acos)
	math1Arg("acosh", math.Acosh)
	math1Arg("asin", math.Asin)
	math1Arg("asinh", math.Asinh)
	math1Arg("atan", math.Atan)
	math1Arg("atanh", math.Atanh)
	math1Arg("cbrt", math.Cbrt)
	math1Arg("ceil", math.Ceil)
	math1Arg("cos", math.Cos)
	math1Arg("cosh", math.Cosh)
	math1Arg("erf", math.Erf)
	math1Arg("erfc", math.Erfc)
	math1Arg("exp", math.Exp)
	math1Arg("exp2", math.Exp2)
	math1Arg("expm1", math.Expm1)
	math1Arg("floor", math.Floor)
	math1Arg("gamma", math.Gamma)
	math1Arg("j0", math.J0)
	math1Arg("j1", math.J1)
	math1Arg("log", math.Log)
	math1Arg("log10", math.Log10)
	math1Arg("log1p", math.Log1p)
	math1Arg("log2", math.Log2)
	math1Arg("logb", math.Logb)
	math1Arg("sin", math.Sin)
	math1Arg("sinh", math.Sinh)
	math1Arg("sqrt", math.Sqrt)
	math1Arg("tan", math.Tan)
	math1Arg("tanh", math.Tanh)
	math1Arg("trunc", math.Trunc)
	math1Arg("y0", math.Y0)
	math1Arg("y1", math.Y1)

	return nil
}

func isZeroPred(m *machine) error {
	a, ok := m.popValue().(number)
	if !ok {
		return fmt.Errorf("value %v was not a number", a)
	}
	if a == Integer(0) || a == Double(0.0) {
		m.pushValue(Boolean(true))
	} else {
		m.pushValue(Boolean(false))
	}
	return nil
}

func executeLessThan(m *machine) error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(b.LessThan(a))
	return nil
}

// Run add on doubles and ints
func executeAdd(m *machine) error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(a.Add(b))
	return nil
}

// Run subtract on doubles and ints
func executeSubtract(m *machine) error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(b.Add(a.Negate()))
	return nil
}

func executeMultiply(m *machine) error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(a.Multiply(b))
	return nil
}

func executeDivide(m *machine) error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(b.Divide(a))
	return nil
}

// Currently modulus is for ints only
func executeModulus(m *machine) error {
	a := m.popValue().(Integer)
	b := m.popValue().(Integer)
	m.pushValue(Integer(b % a))
	return nil
}

func (d Double) Negate() number {
	return number(-d)
}

func (i Integer) Negate() number {
	return number(-i)
}

// TODO: find a way to make this less copy-pasty

func (d Double) Add(n number) number {
	switch n.(type) {
	case Double:
		return number(d + n.(Double))
	case Integer:
		return number(Double(float64(d) + float64(int(n.(Integer)))))
	default:
		panic("Number type error!")
	}
}

func (i Integer) Add(n number) number {
	switch n.(type) {
	case Double:
		return number(Double(i) + n.(Double))
	case Integer:
		return number(Integer(int(i) + int(n.(Integer))))
	default:
		panic("Number type error!")
	}
}

func (i Integer) Multiply(n number) number {
	switch n.(type) {
	case Double:
		return number(Double(i) * n.(Double))
	case Integer:
		return number(Integer(int(i) * int(n.(Integer))))
	default:
		panic("Number type error!")
	}
}

func (d Double) Multiply(n number) number {
	switch n.(type) {
	case Double:
		return number(d * n.(Double))
	case Integer:
		return number(Double(float64(d) * float64(int(n.(Integer)))))
	default:
		panic("Number type error!")
	}
}

func (i Integer) Divide(n number) number {
	switch n.(type) {
	case Double:
		return number(Double(i) / n.(Double))
	case Integer:
		return number(Integer(int(i) / int(n.(Integer))))
	default:
		panic("Number type error!")
	}
}

func (d Double) Divide(n number) number {
	switch n.(type) {
	case Double:
		return number(d / n.(Double))
	case Integer:
		return number(Double(float64(d) / float64(int(n.(Integer)))))
	default:
		panic("Number type error!")
	}
}

func (i Integer) LessThan(n number) Boolean {
	switch n.(type) {
	case Double:
		return Boolean(Double(i) < n.(Double))
	case Integer:
		return Boolean(int(i) < int(n.(Integer)))
	default:
		panic("Number type error!")
	}
}

func (d Double) LessThan(n number) Boolean {
	switch n.(type) {
	case Double:
		return Boolean(d < n.(Double))
	case Integer:
		return Boolean(float64(d) < float64(int(n.(Integer))))
	default:
		panic("Number type error!")
	}
}
