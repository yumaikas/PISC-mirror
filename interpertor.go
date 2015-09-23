package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	EOF                       = fmt.Errorf("End of file")
	EncodingFault             = fmt.Errorf("Encoding error!")
	ConditionalTypeError      = fmt.Errorf("Expected a boolean value, but didn't find it.")
	WordDefParenExpectedError = fmt.Errorf("Word definitions require a stack effect commnet!")
)

type word string

type script struct {
	idx  int
	code string
}

type codeList struct {
	idx  int
	code []word
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
	// TODO: try to work this out later.
	definedStackComments map[word]string
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
	p := &script{
		idx:  0,
		code: code,
	}
	m := machine{
		values: make([]stackEntry, 0),
	}
	//
	executeWordsOnMachine(m, p)
}

var (
	intMatch = regexp.MustCompile(`\d+`)
	// floatMatch = regexp.MustCompile(`\d+\.\d+`)
)

// This executes a given code sequence against a given machine
func executeWordsOnMachine(m machine, p codeSequence) {
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
			// m.executeQuotation()
		case wordVal == ":":
			m.readWordDefinition(p)
		default:
			// Evaluate a defined word, or complain if a word is not defined.

			// Plan of attack: Expand word definition, and push terms into current spot on stack.
			// Hrm....
		}
		if err != nil {
			// TODO: add some ways to debug here....
			fmt.Println("Execution stopped Error: ", err)
			return
		}
	}
}

func (m *machine) readWordDefinition(c codeSequence) error {
	name, err := c.nextWord()
	openPar, err2 := c.nextWord()
	if openPar != "(" {
		return WordDefParenExpectedError
	}
	// TODO: come back and clean this up.
	if err != nil || err2 != nil {
		return fmt.Errorf("Errors %s | %s", err.Error(), err2.Error())
	}

	stackComment := ""
	// read the stack comment for the word
	for err != nil {
		wordVal, err := c.nextWord()
		if err != nil {
			return err
		}
		if wordVal == ")" {
			break
		}
		stackComment += string(word) + " "
	}
	m.definedStackComments[name] = strings.Trim(stackComment)

	wordDef := make([]word, 0)
	for err != nil {
		wordVal, err := c.nextWord()
		if err != nil {
			return err
		}
		if wordVal == ";" {
			break
		}
		wordDef = append(wordDef, wordVal)
	}
	m.definedWords[name] = wordDef
	return nil
}

func (m *machine) runConditionalOperator() error {
	falseVal := m.popValue()
	trueVal := m.popValue()
	booleanVal := m.popValue()

	if b, ok := booleanVal.(Boolean); ok {
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

func (c *codeList) nextWord() (word, error) {
	if c.idx < len(c.code) {
		retval := c.code[c.idx]
		c.idx++
		return retval, nil
	}
	return word(""), EOF
}

func (s *script) nextWord() (word, error) {
	var readingInto string
	for idx, r := range s.code[s.idx:] {
		s.idx = idx
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
