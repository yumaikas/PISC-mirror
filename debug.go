package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

func (m *machine) loadDebugWords() error {
	// ( -- )
	m.predefinedWords["show-prefix-words"] = NilWord(func(m *machine) {
		for name := range m.prefixWords {
			fmt.Println(name)
		}
	})
	// ( quot -- .. time )
	m.predefinedWords["time"] = GoWord(func(m *machine) error {
		words := &codeQuotation{
			idx:   0,
			words: []word{"call"},
		}
		start := time.Now()
		m.execute(words)
		elapsed := time.Since(start)
		m.pushValue(String(fmt.Sprint("Code took ", elapsed)))
		return nil
	})
	// ( filepath quotation -- )
	m.predefinedWords["cpu-pprof"] = GoWord(func(m *machine) error {
		m.executeString("swap")
		path := m.popValue().String()
		f, err := os.Create(path)
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return err
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
			return err
		}
		m.executeQuotation()
		pprof.StopCPUProfile()
		return nil
	})

	// ( -- )
	m.predefinedWords["dump-defined-words"] = GoWord(func(m *machine) error {
		// var words = make(Array, 0)
		for name, seq := range m.prefixWords {
			fmt.Println(":PRE", name, m.definedStackComments[name], DumpToString(seq), ";")
		}
		for name, seq := range m.definedWords {
			fmt.Println(":DOC", name, m.definedStackComments[name], m.helpDocs[name], ";")
			fmt.Println(":", name, m.definedStackComments[name], DumpToString(seq), ";")
		}
		return nil
	})
	return nil
}

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
		words = append(words, string(w))
	}
	return strings.Join(words, " ")
}
