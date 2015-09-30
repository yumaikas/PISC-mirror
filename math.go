package main

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
