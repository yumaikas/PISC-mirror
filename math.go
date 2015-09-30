package main

type number interface {
	Add(number) number
	Negate() number
}

// Proxy the casting work over
func (m *machine) executeAdd() error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(a.Add(b))
	return nil
}

func (m *machine) executeSubtract() error {
	a := m.popValue().(number)
	b := m.popValue().(number)
	m.pushValue(a.Add(b.Negate()))
	return nil
}

func (d Double) Negate() number {
	return number(-d)
}

func (i Integer) Negate() number {
	return number(-i)
}

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
