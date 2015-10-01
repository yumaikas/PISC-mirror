// TODO: look for ways to split this code up better.
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	EOF                       = fmt.Errorf("End of file")
	EncodingFault             = fmt.Errorf("Encoding error!")
	ConditionalTypeError      = fmt.Errorf("Expected a boolean value, but didn't find it.")
	WordDefParenExpectedError = fmt.Errorf("Word definitions require a stack effect commnet!")
	QuotationTypeError        = fmt.Errorf("Expected quotation value, but didn't find it.")
	InvalidAddTypeError       = fmt.Errorf("Expected two integer values, but didn't find them.")
)

type word string

type codeList struct {
	idx    int
	code   []word
	spaces []string
	debug  bool
}

type codeSequence interface {
	nextWord() (word, error)
	currentSpace() string
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

func getWordList(code string) ([]word, []string) {
	words := make([]word, 0)
	for _, v := range spaceMatch.Split(code, -1) {
		words = append(words, word(v))
	}
	return words, spaceMatch.FindAllString(code, -1)
}

func runCode(code string) *machine {
	words, spaces := getWordList(strings.TrimSpace(code))
	p := &codeList{
		idx:    0,
		code:   words,
		spaces: spaces,
	}
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word][]word),
		definedStackComments: make(map[word]string),
	}
	executeWordsOnMachine(m, p)
	return m
}

var (
	spaceMatch = regexp.MustCompile(`\s+`)
	floatMatch = regexp.MustCompile(`-?\d+\.\d+`)
	intMatch   = regexp.MustCompile(`-?\d+`)
)

// TODO: run a tokenizer on the code that handles string literals more appropriately.
// This executes a given code sequence against a given machine
func executeWordsOnMachine(m *machine, p codeSequence) {
	var err error
	var wordVal word
	for err == nil {
		// fmt.Println(intMatch.MatchString(string(wordVal)))
		wordVal, err = p.nextWord()
		if err != nil {
			return
		}
		switch {
		// Comments are going to be exclusively of the /*  */ variety for now.
		case wordVal == "/*":
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal == "*/" {
					break
				}
			}
		case floatMatch.MatchString(string(wordVal)):
			floatVal, floatErr := strconv.ParseFloat(string(wordVal), 64)
			if floatErr != nil {
				panic(floatErr)
			}
			m.pushValue(Double(floatVal))
		case intMatch.MatchString(string(wordVal)):
			intVal, intErr := strconv.Atoi(string(wordVal))
			if intErr != nil {
				panic(intErr)
			}
			m.pushValue(Integer(intVal))
		case wordVal == "dup":
			stackVal := m.popValue()
			m.pushValue(stackVal)
			m.pushValue(stackVal)
		case wordVal == "drop":
			m.popValue()
		case wordVal == "2drop":
			m.popValue()
			m.popValue()
		case wordVal == "3drop":
			m.popValue()
			m.popValue()
			m.popValue()
			// Math words are: +, -, *, /, div, and mod
		case isMathWord(wordVal):
			m.executeMathWord(wordVal)
		case isBooleanWord(wordVal):
			m.executeBooleanWord(wordVal)
		case isStringWord(wordVal):
			m.executeStringWord(wordVal)
		case wordVal == "f":
			m.pushValue(Boolean(false))
		case wordVal == "t":
			// Push true onto stack
			m.pushValue(Boolean(true))
		case wordVal == "[":
			// Begin quotation
			quote := make([]word, 0)
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal == "]" {
					break
				}
				quote = append(quote, wordVal)
			}
			m.pushValue(quotation(quote))
		case wordVal == "]":
			// panic here.
			panic("Unbalanced ]")
		case strings.HasPrefix(string(wordVal), `"`):
			// This is the case for string literals
			strVal := ""
			// Slice out the " chracter
			strVal += strings.TrimPrefix(string(wordVal), `"`) + p.currentSpace()

			if strings.HasSuffix(string(wordVal), `"`) && len(wordVal) > 1 {
				// Slice out the " at the end.
				strVal = strings.TrimPrefix(string(wordVal), `"`)
				strVal = strings.TrimSuffix(strVal, `"`)
				m.pushValue(String(strVal))
				continue
			}
			for err == nil {
				wordVal, err = p.nextWord()
				if strings.HasSuffix(string(wordVal), `"`) {
					strVal += strings.TrimSuffix(string(wordVal), `"`)
					break
				}
				strVal += string(wordVal) + p.currentSpace()
			}
			m.pushValue(String(strVal))
		case wordVal == "?":
			// If there is an error, this will stop the loop.
			err = m.runConditionalOperator()
		case wordVal == "call":
			err := m.executeQuotation()
			if err != nil {
				panic(err)
			}
		case wordVal == ":":
			m.readWordDefinition(p)
		case wordVal == ")":
			panic("Unexpected )")
		case wordVal == ";":
			panic("Unexpected ;")
		case spaceMatch.MatchString(string(wordVal)):
			// TODO: capture this space?
			continue
		default:
			if val, ok := m.definedWords[wordVal]; ok {
				// Run the definition of this word on this machine.
				executeWordsOnMachine(m, &codeList{idx: 0, code: val})
			} else {
				panic("Undefined word: " + string(wordVal))
			}
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
		fmt.Println("ERRR0!")
		return WordDefParenExpectedError
	}
	// TODO: come back and clean this up.
	if err != nil || err2 != nil {
		return fmt.Errorf("Errors %s | %s", err.Error(), err2.Error())
	}

	stackComment := ""
	var wordVal word
	// read the stack comment for the word
	for err == nil {
		wordVal, err = c.nextWord()
		if err != nil {
			return err
		}
		if wordVal == ")" {
			break
		}
		stackComment += string(wordVal) + " "
	}
	fmt.Println("stackComment is", stackComment)
	m.definedStackComments[name] = strings.TrimSpace(stackComment)

	wordDef := make([]word, 0)
	for err == nil {
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

// Used for call word.
func (m *machine) executeQuotation() error {
	quoteVal := m.popValue()
	if q, ok := quoteVal.(quotation); ok {
		executeWordsOnMachine(m, &codeList{idx: 0, code: q})
		return nil
	} else {
		return QuotationTypeError
	}
}

type _type int

const (
	type_int _type = iota
	type_double
	type_else
)

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
		// Return the stack to it's previous state, for debugging...?
		m.pushValue(booleanVal)
		m.pushValue(trueVal)
		m.pushValue(falseVal)
		return ConditionalTypeError
	}
}
func (c *codeList) currentSpace() string {
	if len(c.spaces) == 0 {
		return " "
	}
	if c.idx >= len(c.spaces) {
		return ""
	}
	return c.spaces[c.idx-1]
}

func (c *codeList) nextWord() (word, error) {
	if c.idx < len(c.code) {
		retval := c.code[c.idx]
		c.idx++
		// fmt.Println(retval)
		return retval, nil
	}
	// fmt.Println("EOF")
	return word(""), EOF
}
