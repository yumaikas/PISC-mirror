package main

type GoWord func(*machine) error

func NilWord(f func(*machine)) GoWord {
	return GoWord(func(m *machine) error {
		f(m)
		return nil
	})
}

func (m *machine) loadPredefinedValues() {
	if m.predefinedWords == nil {
		panic("Uninitialized stack machine!")
	}
	m.predefinedWords["t"] = NilWord(func(m *machine) {
		m.pushValue(Boolean(true))
	})
	m.predefinedWords["f"] = NilWord(func(m *machine) {
		m.pushValue(Boolean(false))
	})
	m.predefinedWords["dip"] = NilWord(func(m *machine) {
		quot := m.popValue().(quotation)
		a := m.popValue()
		executeWordsOnMachine(m, &codeList{idx: 0, code: quot})
		m.pushValue(a)
	})
	m.predefinedWords["pick-dup"] = NilWord(func(m *machine) {
		distBack := int(m.popValue().(Integer))
		m.pushValue(m.values[len(m.values)-distBack-1])
	})
	m.predefinedWords["pick-drop"] = NilWord(func(m *machine) {
		distBack := int(m.popValue().(Integer))
		valIdx := len(m.values) - distBack - 1
		val := m.values[valIdx]
		m.values = append(m.values[:valIdx], m.values[valIdx+1:]...)
		m.pushValue(val)
	})
	m.predefinedWords["pick-del"] = NilWord(func(m *machine) {
		distBack := int(m.popValue().(Integer))
		valIdx := len(m.values) - distBack - 1
		m.values = append(m.values[:valIdx], m.values[valIdx+1:]...)
	})
	m.predefinedWords["len"] = NilWord(func(m *machine) {
		length := m.popValue().(lenable).Length()
		m.pushValue(Integer(length))
	})

	words := getWordList(`
		:PRE # ( name -- .. ) ".pisc" concat import ;
		"std_lib.pisc" import `)
	code := &codeList{
		idx:  0,
		code: words,
	}

	m.loadLocalWords()
	m.loadStringWords()
	m.loadBooleanWords()
	m.loadLoopWords()
	m.loadDictWords()
	m.loadVectorWords()
	m.loadSymbolWords()
	m.loadHigherMathWords()
	m.loadHelpWords()

	m.loadIOWords()
	err := executeWordsOnMachine(m, code)
	if err != nil {
		panic(err)
	}

}
