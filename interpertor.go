// TODO: look for ways to split this code up better.
package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	EOF                       = fmt.Errorf("End of file")
	EncodingFault             = fmt.Errorf("Encoding error!")
	ConditionalTypeError      = fmt.Errorf("Expected a boolean value, but didn't find it.")
	WordDefParenExpectedError = fmt.Errorf("Word definitions require a stack effect comment!")
	QuotationTypeError        = fmt.Errorf("Expected quotation value, but didn't find it.")
	InvalidAddTypeError       = fmt.Errorf("Expected two integer values, but didn't find them.")
	UnexpectedStackDashError  = fmt.Errorf("Found unexpected -- in stack annotation")
	ParenBeforeStackDashError = fmt.Errorf("Found ) before -- in stack annotation")
	InvalidPrefixCharError    = fmt.Errorf("Found invalid character in prefix definition")
)

type word string

type codeSequence interface {
	nextWord() (word, error)
	getcodePosition() codePosition
	wrapError(error) error
	// Returns a codeSequence that starts a 0 for the same code
	cloneCode() codeSequence
}

type machine struct {
	// TODO: add a stack pointer so that we can keep from re-allocating a lot.
	// stackPtr int
	values []stackEntry
	// This is reallocated when locals are used
	locals []map[string]stackEntry
	// A map from words to slices of words.
	definedWords map[word]codeSequence
	// A map from prefixes to prefix words
	prefixWords map[word]codeSequence
	// A map from words to predefined words (words built in go)
	predefinedWords map[word]GoWord
	// TODO: try to work this out later.
	definedStackComments map[word]string
	// The top of the stack it the end of the []stackEntry slice.
	// Every so many entries, we may need to re-allocate the stack....
	helpDocs map[word]string

	// Each time we are asked for a symbol, supply the value here, then increment
	symbolIncr int64
}

func (m *machine) pushValue(entry stackEntry) {
	m.values = append(m.values, entry)
}

func (m *machine) genSymbol() {
	m.pushValue(Symbol(m.symbolIncr))
	m.symbolIncr++
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

func runCode(code string) *machine {
	// words := getWordList(strings.TrimSpace(code))
	p := &codeList{
		idx:  0,
		code: code,
	}
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word]codeSequence),
		definedStackComments: make(map[word]string),
	}
	executeWordsOnMachine(m, p)
	return m
}

var (
	spaceMatch       = regexp.MustCompile(`[\s\r\n]+`)
	floatMatch       = regexp.MustCompile(`^-?\d+\.\d+$`)
	intMatch         = regexp.MustCompile(`^-?\d+$`)
	prefixMatchRegex = regexp.MustCompile(`^[-\[\]:!@$%^&*<>+]+`)
)

// This executes a given code sequence against a given machine
func executeWordsOnMachine(m *machine, p codeSequence) (retErr error) {
	var err error
	var wordVal word
	defer func() {
		pErr := recover()
		if pErr != nil {
			retErr = fmt.Errorf("%s", pErr)
		}
		if retErr != nil {
			fmt.Println("Error while executing", wordVal, ":", retErr)
		}
	}()
	for err == nil {
		// fmt.Println(intMatch.MatchString(string(wordVal)))
		wordVal, err = p.nextWord()
		if err == io.EOF {
			return
		}
		if err != nil {
			return err
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
		case strings.HasPrefix(string(wordVal), "#"):
			// Skip line comment, potentialy work with it later, but not now.
			continue
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
		// This word needs to be defined before we can allow other things to be defined
		case wordVal == ":":
			err = m.readWordDefinition(p)
		case wordVal == ":PRE":
			err = m.readPrefixDefinition(p)
			// err = m.readPatternDefinition
		case wordVal == ":DOC":
			err = m.readWordDocumentation(p)
		case wordVal == "typeof":
			m.pushValue(String(m.popValue().Type()))
		case wordVal == "stack-empty?":
			if len(m.values) == 0 {
				m.pushValue(Boolean(true))
			} else {
				m.pushValue(Boolean(false))
			}
		case wordVal == "{":
			// Begin quotation
			pos := p.getcodePosition()
			quote := make([]word, 0)
			depth := 0
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal == "{" {
					depth++
				}
				if wordVal == "}" {
					depth--
				}
				// have to accomodate for the decrement in f
				if wordVal == "}" && depth == -1 {
					break
				}
				quote = append(quote, wordVal)
			}
			currIdx := len(m.values)
			// Run the { } as a quotation
			m.pushValue(quotation{
				code:         quote,
				codePosition: pos,
				locals:       m.locals[len(m.locals)-1]})
			err := m.executeQuotation()
			if err != nil {
				return err
			}
			vals := make([]stackEntry, len(m.values)-currIdx)
			copy(vals, m.values[currIdx:len(m.values)])
			m.values = m.values[:currIdx]
			m.pushValue(Array(vals))

		case wordVal == "[":
			// Begin quotation
			pos := p.getcodePosition()
			quote := make([]word, 0)
			depth := 0
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal == "[" {
					depth++
				}
				if wordVal == "]" {
					depth--
				}
				// have to accomodate for the decrement in f
				if wordVal == "]" && depth == -1 {
					break
				}
				quote = append(quote, wordVal)
			}
			m.pushValue(quotation{
				code:         quote,
				codePosition: pos,
				locals:       m.locals[len(m.locals)-1]})

		case isMathWord(wordVal):
			err = m.executeMathWord(wordVal)
		case wordVal == "]":
			// panic here.
			panic("Unbalanced ]")
		case strings.HasPrefix(string(wordVal), `"`):
			// This is the case for string literals
			// Slice out the " chracter
			strVal := strings.TrimSuffix(strings.TrimPrefix(string(wordVal), `"`), `"`)

			m.pushValue(String(strVal))
		case wordVal == "?":
			// If there is an error, this will stop the loop.
			err = m.runConditionalOperator()
		case wordVal == "call":
			err = m.executeQuotation()
		case wordVal == ")":
			panic("Unexpected )")
		case wordVal == ";":
			panic("Unexpected ;")
		case spaceMatch.MatchString(string(wordVal)):
			// TODO: capture this space?
			continue
		case len(wordVal) == 0:
			continue
		default:
			if fn, ok := m.predefinedWords[wordVal]; ok {
				if ok {
					err = fn(m)
					if err != nil {
						return err
					}
				}
			} else if val, ok := m.definedWords[wordVal]; ok {
				// Run the definition of this word on this machine.
				err = executeWordsOnMachine(m, val.cloneCode())
				if err != nil {
					return err
				}
			} else if prefixFunc, ok := m.prefixWords[getPrefixOf(wordVal)]; ok {
				// Put the post-prefix string at the top of the stack, so it can
				// be used.
				m.pushValue(String(getNonPrefixOf(wordVal)))
				err = executeWordsOnMachine(m, prefixFunc.cloneCode())
			} else if err = m.tryLocalWord(string(wordVal)); err == LocalFuncRun {
				err = nil
				// return p.wrapError(fmt.Errorf("Undefined word: %v", wordVal))
			} else {
				return p.wrapError(fmt.Errorf("Undefined word: %v", wordVal))
			}
			// Evaluate a defined word, or complain if a word is not defined.

			// Plan of attack: Expand word definition, and push terms into current spot on stack.
			// Hrm....
		}
		if err != nil {
			// TODO: add some ways to debug here....
			fmt.Println("Execution stopped during word:  ", wordVal, " error: ", err)
			return err
		}
	}
	return nil
}

var ErrNoLocals = fmt.Errorf("No locals to try !")
var LocalFuncRun = fmt.Errorf("Nothing was wrong")

func (m *machine) tryLocalWord(wordName string) error {
	// TODO: In progress
	if len(m.locals) > 0 {
		if localFunc, found := m.locals[len(m.locals)-1][string(wordName)]; found {
			if fn, ok := localFunc.(quotation); ok {
				code := &codeQuotation{
					idx:          0,
					words:        fn.code,
					codePosition: fn.codePosition,
				}
				err := executeWordsOnMachine(m, code)
				if err != nil {
					return err
				}
			} else if fn, ok := localFunc.(GoFunc); ok {
				err := fn(m)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Value is not a word!")
			}
			return LocalFuncRun
		}
	}
	return ErrNoLocals
}

func (m *machine) readWordBody(c codeSequence) ([]word, error) {
	var err = error(nil)
	openPar, err2 := c.nextWord()
	if openPar != "(" {
		fmt.Println("ERRR0!")
		return nil, WordDefParenExpectedError
	}
	// TODO: come back and clean this up.
	if err != nil || err2 != nil {
		return nil, fmt.Errorf("Errors %s | %s", err.Error(), err2.Error())
	}

	stackComment := ""
	var wordVal word
	// read the stack annotation for the word
	{
		pushes := 0
		pops := 0
		// Counting pushes
		for err == nil {
			wordVal, err = c.nextWord()
			if err != nil {
				return nil, err
			}

			if wordVal == "--" {
				break
			}

			if wordVal == ")" {
				return nil, ParenBeforeStackDashError
			}
			pushes++
			stackComment += string(wordVal) + " "
		}
		// Counting pops
		for err == nil {
			wordVal, err = c.nextWord()
			if err != nil {
				return nil, err
			}

			if wordVal == "--" {
				return nil, UnexpectedStackDashError
			}
			if wordVal == ")" {
				break
			}
			pops++
			stackComment += string(wordVal) + " "
		}
	}
	// fmt.Println("stackComment is", stackComment)
	// m.definedStackComments[name] = strings.TrimSpace(stackComment)

	wordDef := make([]word, 0)
	hasLocal := false
	for err == nil {
		wordVal, err := c.nextWord()
		if err != nil {
			return nil, err
		}
		if isLocalWordPrefix(wordVal) {
			hasLocal = true
		}
		if wordVal == ";" {
			break
		}
		wordDef = append(wordDef, wordVal)
	}
	if hasLocal == true {
		wordDef = append([]word{"get-locals"}, wordDef...)
		wordDef = append(wordDef, "drop-locals")
	}
	return wordDef, nil
}

func getNonPrefixOf(w word) word {
	return word(prefixMatchRegex.ReplaceAllString(string(w), ""))
}

// Prefix words can only start with symbols like :!@#$%^&*
func getPrefixOf(w word) word {
	return word(prefixMatchRegex.FindString(strings.TrimSpace(string(w))))
}

//
// NB. We're going to allocate a lot for now.
func stringFromWordDef(definition []word) string {
	// Copied from std lib's strings.Join
	if len(definition) == 0 {
		return ""
	}
	if len(definition) == 1 {
		return string(definition[0])
	}
	n := len(" ") * (len(definition) - 1)
	for i := 0; i < len(definition); i++ {
		n += len(definition[i])
	}
	b := make([]byte, n)
	bp := copy(b, definition[0])

	for _, s := range definition[1:] {
		bp += copy(b[bp:], " ")
		bp += copy(b[bp:], s)
	}
	return string(b)
}

// Document comments, which end in a ;
func (m *machine) readWordDocumentation(c codeSequence) error {
	word, err := c.nextWord()
	if err != nil {
		return err
	}
	if _, found := m.prefixWords[word]; !found {
	} else if _, found := m.predefinedWords[word]; !found {
	} else if _, found := m.definedWords[word]; !found {
		return fmt.Errorf("No definition for word: %s", word)
	}
	// TODO: Make this it's own loop
	wordDef, err := m.readWordBody(c)
	// Save the docs here
	m.helpDocs[word] = stringFromWordDef(wordDef)
	return err
}

// Prefix (:PRE) definitions, which use a prefix
func (m *machine) readPrefixDefinition(c codeSequence) error {
	pos := c.getcodePosition()
	prefix, err := c.nextWord()
	if err != nil {
		return err
	}
	if !prefixMatchRegex.MatchString(string(prefix)) {
		return InvalidPrefixCharError
	}
	wordDef, err := m.readWordBody(c)
	if err != nil {
		return err
	}
	m.prefixWords[prefix] = &codeQuotation{
		idx:          0,
		words:        wordDef,
		codePosition: pos,
	}
	return nil
}

// Used for : defined words
func (m *machine) readWordDefinition(c codeSequence) error {
	pos := c.getcodePosition()
	name, err := c.nextWord()
	wordDef, err := m.readWordBody(c)
	if err != nil {
		return err
	}
	m.definedWords[name] = &codeQuotation{
		idx:          0,
		words:        wordDef,
		codePosition: pos,
	}
	return nil
}

// func (m *machine) readQuotation(c *codeSequence) error {
// }

// Used for call word.
func (m *machine) executeQuotation() error {
	quoteVal := m.popValue()
	if q, ok := quoteVal.(quotation); ok {
		m.locals = append(m.locals, q.locals)
		executeWordsOnMachine(m, &codeQuotation{
			idx:          0,
			words:        q.code,
			codePosition: q.codePosition,
		})
		m.locals = m.locals[:len(m.locals)-1]
		return nil
	} else if q, ok := quoteVal.(GoFunc); ok {
		return q(m)

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

/*
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
*/
