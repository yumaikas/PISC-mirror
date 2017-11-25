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

func getStackLen(m *Machine) error {
	l := len(m.Values)
	m.PushValue(Integer(l))
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

func logStack(m *Machine) error {
	arr := make([]string, len(m.Values))
	for i, v := range m.Values {
		arr[i] = v.String()
	}

	fmt.Println(strings.Join(arr, "\n"))
	return nil
}

func _showPrefixWords(m *Machine) error {
	for name := range m.PrefixWords {
		fmt.Println(name)
	}
	return nil
}

func _printDebugTrace(m *Machine) error {
	fmt.Println(m.DebugTrace)
	return nil
}

func _clearDebugTrace(m *Machine) error {
	m.DebugTrace = ""
	return nil
}

func _timeQuotation(m *Machine) error {
	start := time.Now()
	err := m.ExecuteQuotation()
	elapsed := time.Since(start)
	m.PushValue(String(fmt.Sprint("Code took ", elapsed)))
	return err
}

func _cpuPPROF(m *Machine) error {
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
}

func loadDebugCore(m *Machine) error {

	m.AddGoWordWithStack("log-stack",
		"( -- )",
		"Log the current state of the stack",
		logStack)

	m.AddGoWordWithStack("stack-len",
		"( -- stack-len )",
		"Get the length of the stack.",
		getStackLen)

	m.AddGoWordWithStack("count-go-words",
		"( -- num-go-words )",
		"Get the number of words in the machine that are defined in Go",
		getNumGoWords)

	m.AddGoWordWithStack("count-pisc-words",
		"( -- num-pisc-words )",
		"Get the number of words in this machine that are defined in PISC",
		getNumPiscWords)

	m.AddGoWordWithStack("count-prefix-words",
		"( -- num-prefix-words )",
		"Get the number of words that are defined as Prefixes",
		getNumPrefixWords)

	// ( -- )
	m.AddGoWordWithStack(
		"show-prefix-words",
		"( -- )",
		"Logs a list of prefix words to stdout",
		_showPrefixWords)

	// ( quot -- .. time )
	m.AddGoWordWithStack(
		"time",
		"( quot -- .. time )",
		"Times how long it takes to execute a quotation, leaving a string atop the stack describing the time",
		_timeQuotation)

	m.AddGoWordWithStack(
		"print-debug-trace",
		"( -- )",
		"Prints the extant debug trace",
		_printDebugTrace)

	m.AddGoWordWithStack(
		"clear-debug-trace",
		"( -- )",
		"Clears the debug trace",
		_clearDebugTrace)

	m.AddGoWordWithStack(
		"cpu-pprof",
		"( path quot -- )",
		"Runs the quotation in a PPROF session, saving the data to the file at path",
		_cpuPPROF)

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
