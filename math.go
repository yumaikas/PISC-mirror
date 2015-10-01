package main

type number interface {
	stackEntry
	Add(number) number
	Negate() number
	Multiply(number) number
	Divide(number) number
	LessThan(number) Boolean
}

func isMathWord(w word) bool {
	return w == "+" ||
		w == "-" ||
		w == "*" ||
		w == "/" ||
		// w == "div" ||
		w == "mod" ||
		w == "<" ||
		w == "zero?"
}

func (m *machine) executeMathWord(w word) error {
	switch w {
	case "+":
		return m.executeAdd()
	case "-":
		return m.executeSubtract()
	case "*":
		return m.executeMultiply()
	case "/":
		return m.executeDivide()
	case "zero?":
		a := m.popValue().(number)
		if a == Integer(0) || a == Double(0.0) {
			m.pushValue(Boolean(true))
		} else {
			m.pushValue(Boolean(false))
		}
	case "mod":
		m.executeModulus()
	}
	return nil
}

// Run add on doubles and ints
func (m *machine) executeAdd() error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(a.Add(b))
	return nil
}

// Run subtract on doubles and ints
func (m *machine) executeSubtract() error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(b.Add(a.Negate()))
	return nil
}

func (m *machine) executeMultiply() error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(a.Multiply(b))
	return nil
}

func (m *machine) executeDivide() error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(b.Divide(a))
	return nil
}

// Currently modulus is for ints only
func (m *machine) executeModulus() error {
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
