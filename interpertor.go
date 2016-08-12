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
	WordDefParenExpectedError = fmt.Errorf("Word definitions require a stack effect comment!")
	QuotationTypeError        = fmt.Errorf("Expected quotation value, but didn't find it.")
	InvalidAddTypeError       = fmt.Errorf("Expected two integer values, but didn't find them.")
	UnexpectedStackDashError  = fmt.Errorf("Found unexpected -- in stack annotation")
	ParenBeforeStackDashError = fmt.Errorf("Found ) before -- in stack annotation")
	InvalidPrefixCharError    = fmt.Errorf("Found invalid character in prefix definition")
)

type word string

type codeList struct {
	idx   int
	code  []word
	debug bool
}

type codeSequence interface {
	nextWord() (word, error)
}

type machine struct {
	// TODO: add a stack pointer so that we can keep from re-allocating a lot.
	// stackPtr int
	values []stackEntry
	// This is reallocated when locals are used
	locals []map[string]stackEntry
	// A map from words to slices of words.
	definedWords map[word]*codeList
	// A map from prefixes to prefix words
	prefixWords map[word]*codeList
	// A map from words to predefined words (words built in go)
	predefinedWords map[word]GoWord
	// TODO: try to work this out later.
	definedStackComments map[word]string
	// The top of the stack it the end of the []stackEntry slice.
	// Every so many entries, we may need to re-allocate the stack....

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

func getWordList(code string) []word {
	words := make([]word, 0)
	currentWord := ""
	skipChar := false
	inString := false
	lineno := 1
	currentLine := ""
	for _, v := range code {
		if v == '\n' {
			// fmt.Println("Parsing:", currentLine)
			currentLine = ""
			lineno++
		}
		currentLine += string(v)
		if skipChar {
			switch v {
			case 'n':
				currentWord += "\n"
				skipChar = false
				continue
			case 't':
				currentWord += "\t"
				skipChar = false
				continue
			case 'r':
				currentWord += "\n"
				skipChar = false
				continue
			case '\\':
				currentWord += `\`
				skipChar = false
				continue
			case '"':
				currentWord += `"`
				skipChar = false
				continue
			default:
				panic(fmt.Sprint(
					"Invalid escape sequence:", v,
					"current word: ", currentWord,
					"line:", lineno))
			}
		}
		switch v {
		case '\\':
			skipChar = true
			continue
		case '"':
			if inString {
				currentWord += "\""
				inString = false
				continue
			} else {
				inString = true
				currentWord += string(v)
			}
		case ' ':
			fallthrough
		case '\t':
			fallthrough
		case '\n':
			fallthrough
		case '\r':
			if inString {
				currentWord += string(v)
			} else {
				words = append(words, word(currentWord))
				currentWord = ""
			}
		default:
			currentWord += string(v)
		}
	}
	words = append(words, word(currentWord))
	if inString {
		panic("Unterminated string!")
	}
	return words
}

func runCode(code string) *machine {
	words := getWordList(strings.TrimSpace(code))
	p := &codeList{
		idx:  0,
		code: words,
	}
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word]*codeList),
		definedStackComments: make(map[word]string),
	}
	executeWordsOnMachine(m, p)
	return m
}

var (
	spaceMatch       = regexp.MustCompile(`[\s\r\n]+`)
	floatMatch       = regexp.MustCompile(`^-?\d+\.\d+$`)
	intMatch         = regexp.MustCompile(`^-?\d+$`)
	prefixMatchRegex = regexp.MustCompile(`^[-\[\]:!@#$%^&*<>]+`)
)

// TODO: run a tokenizer on the code that handles string literals more appropriately.
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
		// This word needs to be defined before we can allow other things to be defined
		case wordVal == ":":
			err = m.readWordDefinition(p)
		case wordVal == ":PRE":
			err = m.readPrefixDefinition(p)
			// err = m.readPatternDefinition
		case wordVal == "typeof":
			m.pushValue(String(m.popValue().Type()))
		case wordVal == "stack-empty?":
			if len(m.values) == 0 {
				m.pushValue(Boolean(true))
			} else {
				m.pushValue(Boolean(false))
			}
		case wordVal == "extern-call":
			wordName := m.popValue().(String).String()
			fn := m.predefinedWords[word(wordName)]
			err = fn(m)
		case wordVal == "[":
			// Begin quotation
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
			m.pushValue(quotation(quote))
		// Local words that use a parser
		/* case strings.HasPrefix(string(wordVal), ":"):
			m.locals[len(m.locals)-1][string(wordVal[1:])] = m.popValue()
		case strings.HasPrefix(string(wordVal), "$"):
			val, found := m.locals[len(m.locals)-1][string(wordVal[1:])]
			if !found {
				return ErrLocalNotFound
			}
			m.pushValue(val)
		*/

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
			err := m.executeQuotation()
			if err != nil {
				panic(err)
			}
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
			if val, ok := m.definedWords[wordVal]; ok {
				// Run the definition of this word on this machine.
				val.idx = 0
				err = executeWordsOnMachine(m, val)
				if err != nil {
					return err
				}
			} else if prefixFunc, ok := m.prefixWords[getPrefixOf(wordVal)]; ok {
				prefixFunc.idx = 0
				// Put the post-prefix string at the top of the stack, so it can
				// be used.
				m.pushValue(String(getNonPrefixOf(wordVal)))
				err = executeWordsOnMachine(m, prefixFunc)
			} else if localFunc, ok := m.locals[len(m.locals)-1][string(wordVal)]; ok {
				if fn, ok := localFunc.(quotation); ok {
					code := &codeList{idx: 0, code: fn}
					err = executeWordsOnMachine(m, code)
					if err != nil {
						return err
					}
				} else {
					cl := p.(*codeList)
					fmt.Println(cl.code[cl.idx-1])
					// fmt.Println(p.(*codeList).code)

					panic(fmt.Sprint("Undefined word: ", wordVal))
				}
			} else {
				cl := p.(*codeList)
				fmt.Println(cl.code[cl.idx-1])
				// fmt.Println(p.(*codeList).code)

				panic(fmt.Sprint("Undefined word: ", wordVal))
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

// IN PROGRESS: This function.
func (m *machine) readVectorLiteral(c codeSequence) error {
	/*
		array := make([]stackEntry, 0)
		depth := 0
		p := &codeList{
			code: make([]word, 1),
			idx:  0,
		}

		for err == nil {
			wordVal, err = c.nextWord()
			p.code[0] = wordVal
			// this could be dicey... it's not threadsafe by default..
			currStackLen := len(m.values) - 1
			executeWordsOnMachine(m, p)
			if len(m.values)-currStackLen != 1 {
				return fmt.Errorf("Word in array returned more than 1 value, this isn't allowed!")
			}
			array = append(array, m.popValue())
			// TODO: Handle nested array literals via recursive calls into readVectorLiteral
			// !IDEA: Have a isLiteralWord and evaluateLiteralToMachine that wrap all the literals and push them to the machine
			// -IDEA: Some kind of array that allows for more than 1 value to come from a call
			if wordVal == "{" {
				depth++
			}
			if wordVal == "}" {
				depth--
			}
			// have to accomodate for the decrement in f in the previous function.
			if wordVal == "}" && depth == -1 {
				break
			}

		}
		m.pushValue(quotation(quote))
	*/
	return nil
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
	return word(prefixMatchRegex.FindString(string(w)))
}

// Prefix (:PRE) definitions, which use a prefix
func (m *machine) readPrefixDefinition(c codeSequence) error {
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
	m.prefixWords[prefix] = &codeList{
		idx:  0,
		code: wordDef,
	}
	return nil
}

// Used for : defined words
func (m *machine) readWordDefinition(c codeSequence) error {
	name, err := c.nextWord()
	wordDef, err := m.readWordBody(c)
	if err != nil {
		return err
	}
	m.definedWords[name] = &codeList{
		idx:  0,
		code: wordDef,
	}
	return nil
}

// func (m *machine) readQuotation(c *codeSequence) error {
// }

// Used for call word.
func (m *machine) executeQuotation() error {
	quoteVal := m.popValue()
	if q, ok := quoteVal.(quotation); ok {
		executeWordsOnMachine(m, &codeList{idx: 0, code: q})
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
