// TODO: look for ways to split this code up better.
package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/asdine/storm"
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
	ExitingProgram            = fmt.Errorf("User called `quit`, terminating program")
)

type word struct {
	str  string
	impl GoWord
}

func (w word) String() string {
	return w.str
}

type codeSequence interface {
	nextWord() (*word, error)
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
	definedWords map[string]*codeQuotation
	// A map from prefixes to prefix words
	prefixWords map[string]*codeQuotation
	// A map from words to predefined words (words built in go)
	predefinedWords map[string]GoWord
	// TODO: try to work this out later.
	definedStackComments map[string]string
	// The top of the stack it the end of the []stackEntry slice.
	// Every so many entries, we may need to re-allocate the stack....
	helpDocs map[string]string

	// Each time we are asked for a symbol, supply the value here, then increment
	symbolIncr int64
	db         *storm.DB
}

//TODO: Optimize append pattern?
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

	// Allow pieces of the stack to be GCd?
	/*
		if cap(m.values) > 4*len(m.values) {
			fmt.Println("Copying stack for GC purposes")
			dest := make([]stackEntry, len(m.values))
			copy(m.values[:], dest[:])
			m.values = dest
		}
	*/
	popped := m.values[len(m.values)-1]
	// Snip the length of value
	m.values = m.values[:len(m.values)-1]
	return popped
}

func noop(m *machine) error {
	return nil
}

func (m *machine) executeString(code string, pos codePosition) error {
	p, err := stringToQuotation(code, pos)
	if err != nil {
		return err
	}
	return m.execute(p)
}

var (
	// spaceMatch       = regexp.MustCompile(`[\s\r\n]+`)
	// floatMatch       = regexp.MustCompile(`^-?\d+\.\d+$`)
	// intMatch         = regexp.MustCompile(`^-?\d+$`)
	prefixMatchRegex = regexp.MustCompile(`^[-\[\]:~!@$%^&*<>+?]+`)
)

var prefixChars = []rune{'-', '[', ']', ':', '~', '!', '@', '$', '%', '^', '&', '*', '<', '>', '+', '?'}

func isPrefixChar(r rune) bool {
	for _, c := range prefixChars {
		if r == c {
			return true
		}
	}
	return false
}

func (m *machine) hasPrefixWord(w word) (*codeQuotation, string, bool) {
	if !isPrefixChar(rune(w.str[0])) {
		return nil, "", false
	}
	var prefix string

	for i, r := range w.str {
		if !isPrefixChar(r) {
			prefix = string(w.str[0:i])
			seq, ok := m.prefixWords[prefix]
			return seq, string(w.str[i:len(w.str)]), ok
		}
	}
	return nil, "", false
}

type intPusher int

func (i intPusher) pushInt(m *machine) error {
	m.pushValue(Integer(i))
	return nil
}

// Both of the functions below are non-idomatic shenanegains, but chould present some pretty large gains...
func tryParseInt(w *word, intVal *int) bool {
	// This is a very easy early exit oppurtunity for this function.
	if !strings.ContainsRune("0123456789-+", rune(w.str[0])) {
		return false
	}
	var err error
	*intVal, err = strconv.Atoi(string(w.str))
	ip := intPusher(*intVal)
	// This is a method reference, will it work, and will it be more efficient than a closure or the like?
	w.impl = ip.pushInt
	return err == nil
}

type floatPusher float64

func (f floatPusher) pushFloat(m *machine) error {
	m.pushValue(Double(f))
	return nil
}
func tryParseFloat(w *word, floatVal *float64) bool {
	// This is a very easy early exit oppurtunity for this function.
	if !strings.ContainsRune("0123456789-+", rune(w.str[0])) {
		return false
	}
	var err error
	*floatVal, err = strconv.ParseFloat(string(w.str), 64)
	fp := floatPusher(*floatVal)
	w.impl = fp.pushFloat
	return err == nil
}

func wordIsWhitespace(w word) bool {
	for _, r := range w.str {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// This executes a given code sequence against a given machine
func (m *machine) execute(p *codeQuotation) (retErr error) {
	var err error
	var wordVal *word
	var old_idx = p.idx
	p.idx = 0
	// This isn't efficient, but it's a way to keep from setting stuff to 0 everywhere else.
	defer func() {
		p.idx = old_idx
	}()
	defer func() {
		pErr := recover()
		if pErr != nil {
			retErr = fmt.Errorf("%s", pErr)
		}
		if retErr != nil {
			fmt.Println("Error while executing", wordVal, ":", p.wrapError(retErr))
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
		var intVal int
		var floatVal float64
		switch {
		case wordVal.impl != nil:
			err = wordVal.impl(m)
		case wordVal.str == "quit":
			return ExitingProgram
		// Comments are going to be exclusively of the /*  */ variety for now.
		case wordVal.str == "/*":
			commentWordLen := 0
			for err == nil {
				wordVal, err = p.nextWord()
				commentWordLen++
				if err != nil {
					break
				}
				if wordVal.str == "*/" {
					break
				}
			}
			wordVal.impl = func(m *machine) error {
				// fmt.Println("From impl")
				for i := 0; i <= commentWordLen; i++ {
					p.nextWord()
				}
				return nil
			}
		case strings.HasPrefix(wordVal.str, "#"):
			// Skip line comment, potentialy work with it later, but not now.
			continue
		case tryParseInt(wordVal, &intVal):
			m.pushValue(Integer(intVal))
		case tryParseFloat(wordVal, &floatVal):
			m.pushValue(Double(floatVal))
		// This word needs to be defined before we can allow other things to be defined
		case wordVal.str == ":":
			err = m.readWordDefinition(p)
		case wordVal.str == ":PRE":
			err = m.readPrefixDefinition(p)
			// err = m.readPatternDefinition
		case wordVal.str == ":DOC":
			err = m.readWordDocumentation(p)
		case wordVal.str == "typeof":
			m.pushValue(String(m.popValue().Type()))
		case wordVal.str == "stack-empty?":
			if len(m.values) == 0 {
				m.pushValue(Boolean(true))
			} else {
				m.pushValue(Boolean(false))
			}
		case wordVal.str == "}":
			panic("Unbalanced }!")
		case wordVal.str == "{":
			// fmt.Println("Not cached")
			__quot := &codeQuotation{
				words:         make([]*word, 0),
				codePositions: make([]codePosition, 0),
			}
			// This is the word that will be patched in holding the quotation words we're about to patch out.
			// quote := make([]*word, 0)
			depth := 0
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal.str == "${" {
					depth++
				}
				if wordVal.str == "{" {
					depth++
				}
				if wordVal.str == "}" {
					depth--
				}
				// have to accomodate for the decrement in f
				if wordVal.str == "}" && depth == -1 {
					break
				}
				__quot.codePositions = append(__quot.codePositions, p.getcodePosition())
				__quot.words = append(__quot.words, wordVal)
			}
			_quot := &quotation{
				inner:  __quot,
				locals: m.locals[len(m.locals)-1]}
			// Run the { } as a quotation
			_quot.locals = m.locals[len(m.locals)-1]
			currIdx := len(m.values)
			m.pushValue(_quot)
			//fmt.Println(_quot.inner.words)
			err := m.executeQuotation()
			if err != nil {
				return p.wrapError(err)
			}
			vals := make([]stackEntry, len(m.values)-currIdx)
			copy(vals, m.values[currIdx:len(m.values)])
			m.values = m.values[:currIdx]
			m.pushValue(Array(vals))

		case wordVal.str == "${":
			// Begin quotation
			pos := p.getcodePosition()
			//quote := make([]*word, 0)
			quot := &codeQuotation{
				words:         make([]*word, 0),
				codePositions: make([]codePosition, 0),
			}
			depth := 0
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal.str == "{" {
					depth++
				}
				if wordVal.str == "${" {
					depth++
				}
				if wordVal.str == "}" {
					depth--
				}
				// have to accomodate for the decrement in f
				if wordVal.str == "}" && depth == -1 {
					break
				}
				quot.codePositions = append(quot.codePositions, p.getcodePosition())
				quot.words = append(quot.words, wordVal)
			}
			// anchorWord.str = fmt.Sprint("herp", len(__quot.words))
			// Run the { } as a quotation
			_quot := &quotation{
				inner:  quot,
				locals: m.locals[len(m.locals)-1],
			}
			join, err := stringToQuotation(`"" str-join`, pos)
			if err != nil {
				panic(err)
			}
			_join := &quotation{
				inner:  join,
				locals: m.locals[len(m.locals)-1],
			}
			_quot.locals = m.locals[len(m.locals)-1]
			currIdx := len(m.values)
			// m.pushValue(_quot)
			err = m.execute(_quot.inner)
			//err := m.executeQuotation()
			if err != nil {
				return err
			}
			vals := make([]stackEntry, len(m.values)-currIdx)
			copy(vals, m.values[currIdx:len(m.values)])
			m.values = m.values[:currIdx]
			m.pushValue(Array(vals))
			err = m.execute(_join.inner)

		case wordVal.str == "[":
			// Begin quotation
			__quot := &codeQuotation{
				words:         make([]*word, 0),
				codePositions: make([]codePosition, 0),
			}
			depth := 0
			for err == nil {
				wordVal, err = p.nextWord()
				if wordVal.str == "[" {
					depth++
				}
				if wordVal.str == "]" {
					depth--
				}
				// have to accomodate for the decrement in f
				if wordVal.str == "]" && depth == -1 {
					break
				}
				__quot.codePositions = append(__quot.codePositions, p.getcodePosition())
				__quot.words = append(__quot.words, wordVal)
			}
			_quotation := &quotation{
				inner:  __quot,
				locals: m.locals[len(m.locals)-1],
			}
			_quotation.locals = m.locals[len(m.locals)-1]
			m.pushValue(_quotation)

		case isMathWord(*wordVal):
			err = m.executeMathWord(wordVal)
		case wordVal.str == "]":
			// panic here.
			panic("Unbalanced ]")
		case strings.HasPrefix(wordVal.str, `"`):
			// This is the case for string literals
			// Slice out the " chracter
			strVal := strings.TrimSuffix(strings.TrimPrefix(wordVal.str, `"`), `"`)

			m.pushValue(String(strVal))
		case wordVal.str == "?":
			// If there is an error, this will stop the loop.
			wordVal.impl = func(m *machine) error { return m.runConditionalOperator() }
			err = m.runConditionalOperator()
		case wordVal.str == "call":
			wordVal.impl = func(m *machine) error { return m.executeQuotation() }
			err = m.executeQuotation()
		case wordVal.str == ")":
			panic("Unexpected )")
		case wordVal.str == ";":
			panic("Unexpected ;")
		case wordIsWhitespace(*wordVal):
			wordVal.impl = func(m *machine) error {
				return nil
			}
			// TODO: capture this space?
			continue
		case len(wordVal.str) == 0:
			wordVal.impl = func(m *machine) error {
				return nil
			}
			continue
		default:
			if fn, ok := m.predefinedWords[wordVal.str]; ok {
				if ok {
					err = fn(m)
					if err != nil {
						return err
					}
					wordVal.impl = fn
				}
			} else if val, ok := m.definedWords[wordVal.str]; ok {
				// Run the definition of this word on this machine.
				// val.idx = 0
				err = m.execute(val)
				if err != nil {
					return err
				}
				// This is a closure, so I'll need to be careful about it's performance
				wordVal.impl = func(m *machine) error {
					// val.idx = 0
					return m.execute(val)
				}
			} else if prefixFunc, nonPrefix, ok := m.hasPrefixWord(*wordVal); ok {
				// Put the post-prefix string at the top of the stack, so it can
				// be used.
				m.pushValue(String(nonPrefix))
				err = m.execute(prefixFunc)
				// Captures prefix/nonprefix
				wordVal.impl = func(m *machine) error {
					m.pushValue(String(nonPrefix))
					return m.execute(prefixFunc)
				}
			} else if err = m.tryLocalWord(wordVal); err == LocalFuncRun {
				err = nil
			} else {
				return p.wrapError(fmt.Errorf("Undefined word: %v, %v", wordVal, err))
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

func (c *quotation) execute(m *machine) error {
	return m.execute(c.inner)
}

func (m *machine) tryLocalWord(w *word) error {
	// TODO: In progress
	if len(m.locals) > 0 {
		if localFunc, found := m.locals[len(m.locals)-1][w.str]; found {
			if fn, ok := localFunc.(*quotation); ok {
				w.impl = fn.execute
				err := w.impl(m)
				if err != nil {
					return err
				}
				return LocalFuncRun
			} else if fn, ok := localFunc.(GoFunc); ok {
				w.impl = GoWord(fn)
				err := fn(m)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("value is not a word")
			}
			return LocalFuncRun
		}
		return LocalFuncRun
	}
	return ErrNoLocals
}

func (m *machine) readWordBody(name word, c codeSequence) ([]*word, []codePosition, error) {
	var err = error(nil)
	openPar, err2 := c.nextWord()
	if openPar.str != "(" {
		fmt.Println("ERRR0!")
		return nil, nil, WordDefParenExpectedError
	}
	// TODO: come back and clean this up.
	if err != nil || err2 != nil {
		return nil, nil, fmt.Errorf("Errors %s | %s", err.Error(), err2.Error())
	}

	stackComment := "( "
	var wordVal *word
	// read the stack annotation for the word
	{
		pushes := 0
		pops := 0
		// Counting pushes
		for err == nil {
			wordVal, err = c.nextWord()
			if err != nil {
				return nil, nil, err
			}

			if wordVal.str == "--" {
				stackComment += "-- "
				break
			}

			if wordVal.str == ")" {
				return nil, nil, ParenBeforeStackDashError
			}
			pushes++
			stackComment += wordVal.str + " "
		}
		// Counting pops
		for err == nil {
			wordVal, err = c.nextWord()
			if err != nil {
				return nil, nil, err
			}

			if wordVal.str == "--" {
				return nil, nil, UnexpectedStackDashError
			}
			if wordVal.str == ")" {
				stackComment += ")"
				break
			}
			pops++
			stackComment += wordVal.str + " "
		}
	}
	// fmt.Println("stackComment is", stackComment)
	m.definedStackComments[name.str] = strings.TrimSpace(stackComment)

	wordDef := make([]*word, 0)
	wordInfo := make([]codePosition, 0)
	hasLocal := false
	for err == nil {
		wordVal, err := c.nextWord()
		if err != nil {
			return nil, nil, err
		}
		if isLocalWordPrefix(*wordVal) {
			hasLocal = true
		}
		if wordVal.str == ";" {
			break
		}
		wordDef = append(wordDef, wordVal)
		wordInfo = append(wordInfo, c.getcodePosition())
	}
	if hasLocal == true {
		wordDef = append([]*word{&word{str: "get-locals"}}, wordDef...)
		wordInfo = append([]codePosition{codePosition{source: "Hardcoded"}}, wordInfo...)
		wordDef = append(wordDef, &word{str: "drop-locals"})
		wordInfo = append(wordInfo, codePosition{source: "Hardcoded"})
	}
	// fmt.Println(wordDef, wordInfo)
	return wordDef, wordInfo, nil
}

func getNonPrefixOf(w word) string {
	return prefixMatchRegex.ReplaceAllString(w.str, "")
}

// Prefix words can only start with symbols like :!@#$%^&*
func getPrefixOf(w word) string {
	return prefixMatchRegex.FindString(strings.TrimSpace(w.str))
}

//
// NB. We're going to allocate a lot for now.
func stringFromWordDef(definition []*word) string {
	// Copied from std lib's strings.Join
	if len(definition) == 0 {
		return ""
	}
	if len(definition) == 1 {
		return definition[0].str
	}
	n := len(" ") * (len(definition) - 1)
	for i := 0; i < len(definition); i++ {
		n += len(definition[i].str)
	}
	b := make([]byte, n)
	bp := copy(b, definition[0].str)

	for _, s := range definition[1:] {
		bp += copy(b[bp:], " ")
		bp += copy(b[bp:], s.str)
	}
	return string(b)
}

// Document comments, which end in a ;
func (m *machine) readWordDocumentation(c codeSequence) error {
	word, err := c.nextWord()
	if err != nil {
		return err
	}
	if _, found := m.prefixWords[word.str]; !found {
	} else if _, found := m.predefinedWords[word.str]; !found {
	} else if _, found := m.definedWords[word.str]; !found {
		return fmt.Errorf("No definition for word: %s", word)
	}
	// TODO: Make this it's own loop
	wordDef, _, err := m.readWordBody(*word, c)
	// Save the docs here
	m.helpDocs[word.str] = stringFromWordDef(wordDef)
	return err
}

// Prefix (:PRE) definitions, which use a prefix
func (m *machine) readPrefixDefinition(c codeSequence) error {
	prefix, err := c.nextWord()
	if err != nil {
		return err
	}
	if !prefixMatchRegex.MatchString(prefix.str) {
		return InvalidPrefixCharError
	}
	wordDef, positions, err := m.readWordBody(*prefix, c)
	if err != nil {
		return err
	}
	m.prefixWords[prefix.str] = &codeQuotation{
		idx:           0,
		words:         wordDef,
		codePositions: positions,
	}
	return nil
}

// Used for : defined words
func (m *machine) readWordDefinition(c codeSequence) error {
	name, err := c.nextWord()
	wordDef, positions, err := m.readWordBody(*name, c)
	if err != nil {
		fmt.Println(c.getcodePosition())
		return err
	}
	m.definedWords[name.str] = &codeQuotation{
		idx:           0,
		words:         wordDef,
		codePositions: positions,
	}
	return nil
}

// Used for call word.
func (m *machine) executeQuotation() error {
	quoteVal := m.popValue()
	if q, ok := quoteVal.(*quotation); ok {
		// q.inner = q.inner.cloneCode()
		m.locals = append(m.locals, q.locals)
		m.execute(q.toCode())
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
