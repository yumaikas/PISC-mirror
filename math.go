package pisc

import (
	"fmt"
	"math"
)

var ModMathCore = Module{
	Author:    "Andreb Owen",
	Name:      "MathCore",
	License:   "MIT",
	DocString: `Basic math stuff.`,
	Load:      loadMathCore,
}

func loadMathCore(m *Machine) error {
	return m.loadHigherMathWords()
}

type number interface {
	StackEntry
	Add(number) number
	Negate() number
	Multiply(number) number
	Divide(number) number
	LessThan(number) Boolean
}

func (m *Machine) loadHigherMathWords() error {

	m.AddGoWord("+", "( a b -- c )  addition", executeAdd)
	m.AddGoWord("-", " ( a b -- c ) subtraction ", executeSubtract)
	m.AddGoWord("*", " ( a b -- c ) multiplication ", executeMultiply)
	m.AddGoWord("/", " ( a b -- c ) division ", executeDivide)
	m.AddGoWord("mod", " ( a b -- c ) modulus ", executeModulus)
	m.AddGoWord("<", " ( a b -- c ) numeric less than ", executeLessThan)
	m.AddGoWord("zero?", " ( a -- isZero? ) returns if a is zero or not ", isZeroPred)

	// For now, PISC words are late-bound, so we can get away with this.
	err := m.ImportPISCAsset("stdlib/math.pisc")
	if err != nil {
		return err
	}

	// Why do we duplicate the work here?
	// Because we want both the >double word and the
	// math1arg words
	m.PredefinedWords[">double"] = GoWord(func(m *Machine) error {
		a := m.PopValue()
		if i, ok := a.(Integer); ok {
			m.PushValue(Double(float64(i)))
			return nil
		}
		if d, ok := a.(Double); ok {
			m.PushValue(Double(d))
			return nil
		}
		return fmt.Errorf("Unexpected type %s cannot be coverted to double", a.Type())
	})

	var math1Arg = func(name string, mathyFunc func(float64) float64) {
		m.PredefinedWords[name] = GoWord(func(m *Machine) error {
			a := m.PopValue()
			if i, ok := a.(Integer); ok {
				m.PushValue(Double(mathyFunc(float64(i))))
				return nil
			}
			if d, ok := a.(Double); ok {
				m.PushValue(Double(mathyFunc(float64(d))))
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

func isZeroPred(m *Machine) error {
	a, ok := m.PopValue().(number)
	if !ok {
		return fmt.Errorf("value %v was not a number", a)
	}
	if a == Integer(0) || a == Double(0.0) {
		m.PushValue(Boolean(true))
	} else {
		m.PushValue(Boolean(false))
	}
	return nil
}

func executeLessThan(m *Machine) error {
	a := m.PopValue().(number)
	b := m.PopValue().(number)
	m.PushValue(b.LessThan(a))
	return nil
}

// Run add on doubles and ints
func executeAdd(m *Machine) error {
	a := m.PopValue().(number)
	b := m.PopValue().(number)
	m.PushValue(a.Add(b))
	return nil
}

// Run subtract on doubles and ints
func executeSubtract(m *Machine) error {
	a := m.PopValue().(number)
	b := m.PopValue().(number)
	m.PushValue(b.Add(a.Negate()))
	return nil
}

func executeMultiply(m *Machine) error {
	a := m.PopValue().(number)
	b := m.PopValue().(number)
	m.PushValue(a.Multiply(b))
	return nil
}

func executeDivide(m *Machine) error {
	a := m.PopValue().(number)
	b := m.PopValue().(number)
	m.PushValue(b.Divide(a))
	return nil
}

// Currently modulus is for ints only
func executeModulus(m *Machine) error {
	a := m.PopValue().(Integer)
	b := m.PopValue().(Integer)
	m.PushValue(Integer(b % a))
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
