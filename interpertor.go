package main

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var (
	EOF                  = fmt.Errorf("End of file")
	EncodingFault        = fmt.Errorf("Encoding error!")
	ConditionalTypeError = fmt.Errorf("Expected a boolean value, but didn't find it.")
)

type word string

type script struct {
	idx  int
	code string
}

type codeSequence interface {
	nextWord() (word, error)
}

type machine struct {
	// TODO: add a stack pointer so that we can keep from re-allocating a lot.
	// stackPtr int
	values []stackEntry
	// A map from words to slices of words.
	definedWords map[word][]word
	// The top of the stack it the end of the []stackEntry slice.
	// Every so many entries, we may need to re-allocate the stack....
}

func (m *machine) pushValue(entry stackEntry) {
	m.values = append(m.values, entry)
}

func (m *machine) popValue() stackEntry {
	if len(m.values) <= 0 {
		panic("Data underflow")
	}

	popped := m.values[len(m.values)-1]
	// Snip the length of value
	m.values = m.values[:len(m.values)-1]
	return popped
}

func runCode(code string) {
	intMatch := regexp.MustCompile(`\d+`)
	// floatMatch := regexp.MustCompile(`\d+\.\d+`)

	p := script{
		idx:  0,
		code: code,
	}

	m := machine{values: make([]stackEntry, 0)}

	var err error
	var wordVal word
	for err == nil {
		wordVal, err = p.nextWord()
		switch {
		case intMatch.MatchString(string(wordVal)):
			intVal, intErr := strconv.Atoi(string(wordVal))
			if intErr != nil {
				panic(intErr)
			}
			m.pushValue(intVal)
		case wordVal == "dup":
			stackVal := m.popValue()
			m.pushValue(stackVal)
			m.pushValue(stackVal)
		case wordVal == "f":
			m.pushValue(Boolean(false))
		case wordVal == "t":
			m.pushValue(Boolean(true))
			// Push true onoto stack
		case wordVal == "[":
			// Begin quotation
			quotation := make([]word, 0)
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal == "]" {
					break
				}
				quotation = append(quotation, wordVal)
			}
			m.pushValue(quotation)
		case wordVal == "]":
			// panic here.
			panic("Unbalanced ]")
		case wordVal == "?":
			// If there is an error, this will stop the loop.
			err = m.runConditionalOperator()
		case wordVal == "call":

		default:
			// Evaluate a defined word, or complain if a word is not defined.

			// Plan of attack: Expand word definition, and push terms into current spot on stack.
			// Hrm....
		}
		if err != nil {
			fmt.Println("Execution stopped at index", p.idx, "Error: ", err)
			return
		}
	}
}

func (m *machine) runConditionalOperator() error {
	falseVal := m.popValue()
	trueVal := m.popValue()
	booleanVal := m.popValue()

	if b, ok := booleanval.(Boolean); ok {
		if b {
			m.pushValue(trueVal)
		} else {
			m.pushValue(falseVal)
		}
		return nil
	} else {
		m.pushValue(booleanVal)
		m.pushValue(trueVal)
		m.pushValue(falseVal)
		return ConditionalTypeError
	}
}

func (s *script) nextWord() (word, error) {
	var readingInto string
	for idx, r := range s.code[s.idx:] {
		if r == utf8.RuneError {
			return word(""), EncodingFault
		}
		if unicode.IsSpace(r) {
			return word(readingInto), nil
		}
		readingInto += string(r)
	}
	return word(readingInto), EOF
}
