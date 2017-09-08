package pisc

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

var ModDebugCore = Module{
	Author:    "Andrew Owen",
	Name:      "DebugCore",
	License:   "MIT",
	DocString: "The debug words used in PISC",
	Load:      loadDebugCore,
}

// Inspired by a computerphile video on how
// Postscript was a "big" language at 400
// "Operators" (words in PISC), so I can instruement
// the relative size of PISC in terms of surface area

func getNumGoWords(m *Machine) error {
	m.PushValue(Integer(len(m.PredefinedWords)))
	return nil
}

func getNumPiscWords(m *Machine) error {
	m.PushValue(Integer(len(m.DefinedWords)))
	return nil
}

func getNumPrefixWords(m *Machine) error {
	m.PushValue(Integer(len(m.PrefixWords)))
	return nil
}

func loadDebugCore(m *Machine) error {

	m.AddGoWord("count-go-words",
		"( -- num-go-words )",
		getNumGoWords)

	m.AddGoWord("count-pisc-words",
		"( -- num-pisc-words )",
		getNumPiscWords)

	m.AddGoWord("count-prefix-words",
		"( -- num-go-words )",
		getNumPrefixWords)

	// ( -- )
	m.PredefinedWords["show-prefix-words"] = NilWord(func(m *Machine) {
		for name := range m.PrefixWords {
			fmt.Println(name)
		}
	})
	// ( quot -- .. time )
	m.AddGoWord("time", "( quot -- .. time )", GoWord(func(m *Machine) error {
		words := &CodeQuotation{
			Idx:   0,
			Words: []*word{&word{str: "call"}},
		}
		start := time.Now()
		err := m.execute(words)
		elapsed := time.Since(start)
		m.PushValue(String(fmt.Sprint("Code took ", elapsed)))
		return err
	}))

	m.AddGoWord("print-debug-trace", "( -- )", func(m *Machine) error {
		fmt.Println(m.DebugTrace)
		return nil
	})

	m.AddGoWord("clear-debug-trace", "( -- )", func(m *Machine) error {
		m.DebugTrace = ""
		return nil
	})

	// ( filepath Quotation -- )
	m.PredefinedWords["cpu-pprof"] = GoWord(func(m *Machine) error {
		m.ExecuteString("swap", CodePosition{Source: "cpu-pprof GoWord"})
		path := m.PopValue().String()
		f, err := os.Create(path)
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return err
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
			return err
		}
		m.ExecuteQuotation()
		pprof.StopCPUProfile()
		return nil
	})

	// ( -- )
	m.PredefinedWords["dump-defined-words"] = GoWord(func(m *Machine) error {
		// var words = make(Array, 0)
		for name, seq := range m.PrefixWords {
			fmt.Println(":PRE", name, m.DefinedStackComments[name], DumpToString(seq), ";")
		}
		for name, seq := range m.DefinedWords {
			fmt.Println(":DOC", name, m.DefinedStackComments[name], m.HelpDocs[name], ";")
			fmt.Println(":", name, m.DefinedStackComments[name], DumpToString(seq), ";")
		}
		return nil
	})
	return m.ImportPISCAsset("stdlib/debug.pisc")
}

// TODO: See about this...
func DumpToString(c codeSequence) string {
	c = c.cloneCode()
	words := make([]string, 0)
	for {
		w, err := c.nextWord()
		if err == io.EOF {
			return strings.Join(words, " ")
		}
		if err != nil {
			panic("Unexpected error!!!")
		}
		words = append(words, w.str)
	}
}
