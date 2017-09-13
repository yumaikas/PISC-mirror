// TODO: look for ways to split this code up better.
package pisc

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	EOF                       = fmt.Errorf("End of file")
	EncodingFault             = fmt.Errorf("Encoding error!")
	ConditionalTypeError      = fmt.Errorf("Expected a boolean value, but didn't find it.")
	WordDefParenExpectedError = fmt.Errorf("Word definitions require a stack effect comment!")
	QuotationTypeError        = fmt.Errorf("Expected Quotation value, but didn't find it.")
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

// Module The type for a function that extends PISC via adding words
type Module struct {
	Author    string
	Name      string
	License   string
	DocString string
	Load      func(m *Machine) error
}

// LoadModules Load a bunch of modules into this instance of the VM
func (m *Machine) LoadModules(modules ...Module) error {
	for _, mod := range modules {
		err := mod.Load(m)
		// Append the name
		m.LoadedModules = append(m.LoadedModules, mod.Name)
		if err != nil {
			return fmt.Errorf("Error loading %s, %s", mod.Name, err.Error())
		}
	}
	return nil
}

func (w word) String() string {
	return w.str
}

type codeSequence interface {
	nextWord() (*word, error)
	getCodePosition() CodePosition
	wrapError(error) error
	// Returns a codeSequence that starts a 0 for the same code
	cloneCode() codeSequence
}

type Machine struct {
	// TODO: add a stack pointer so that we can keep from re-allocating a lot.
	// stackPtr int
	Values []StackEntry

	// Keep a list of loaded modules names here, for doing some detection
	LoadedModules []string

	// This is reallocated when locals are used
	Locals []map[string]StackEntry
	// A map from words to slices of words.
	DefinedWords map[string]*CodeQuotation
	// A map from prefixes to prefix words
	PrefixWords map[string]*CodeQuotation
	// A map from words to predefined words (words built in go)
	PredefinedWords map[string]GoWord
	// TODO: try to work this out later.
	DefinedStackComments map[string]string
	// The top of the stack it the end of the []StackEntry slice.
	// Every so many entries, we may need to re-allocate the stack....
	HelpDocs map[string]string

	// Each time we are asked for a symbol, supply the value here, then increment
	SymbolIncr int64
	// Keep a default database around...
	NumDispatches int64

	DebugTrace string
}

func (m *Machine) LogAndResetDispatchCount(w io.Writer) {
	fmt.Fprintln(w, m.NumDispatches, "dispatches have occured")
	m.NumDispatches = 0
}

// Not the most efficient way, but should work for starting
func (m *Machine) trace(msg string) {
	m.DebugTrace += msg
}

// Append uses an allocation pattern via Go to amortize the number of allocations performed
func (m *Machine) PushValue(entry StackEntry) {
	m.Values = append(m.Values, entry)
}

func (m *Machine) genSymbol() {
	m.PushValue(Symbol(m.SymbolIncr))
	m.SymbolIncr++
}

func (m *Machine) PopValue() StackEntry {
	if len(m.Values) <= 0 {
		panic("Data underflow")
	}

	popped := m.Values[len(m.Values)-1]
	// Snip the length of value
	m.Values = m.Values[:len(m.Values)-1]
	return popped
}

func noop(m *Machine) error {
	return nil
}

func (m *Machine) ExecuteString(code string, pos CodePosition) error {
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
	prefixMatchRegex = regexp.MustCompile(`^[-_:~!@$%^&*<>+?.]+`)
)

var prefixChars = []rune{'-', '_', ':', '~', '!', '@', '$', '%', '^', '&', '*', '<', '>', '+', '?', '.'}

func isPrefixChar(r rune) bool {
	for _, c := range prefixChars {
		if r == c {
			return true
		}
	}
	return false
}

func (m *Machine) hasPrefixWord(w word) (*CodeQuotation, string, bool) {
	if !isPrefixChar(rune(w.str[0])) {
		return nil, "", false
	}
	var prefix string

	for i, r := range w.str {
		if !isPrefixChar(r) {
			prefix = string(w.str[0:i])
			seq, ok := m.PrefixWords[prefix]
			return seq, string(w.str[i:len(w.str)]), ok
		}
	}
	return nil, "", false
}

type intPusher int

func (i intPusher) pushInt(m *Machine) error {
	m.PushValue(Integer(i))
	return nil
}

// Both of the functions below are non-idomatic shenanegains, but chould present some pretty large gains...
func tryParseInt(w *word, intVal *int) bool {
	if len(w.str) <= 0 {
		return false
	}
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

func (f floatPusher) pushFloat(m *Machine) error {
	m.PushValue(Double(f))
	return nil
}

func tryParseFloat(w *word, floatVal *float64) bool {
	if len(w.str) <= 0 {
		return false
	}
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

func (m *Machine) execute(p *CodeQuotation) error {
	var old_idx = p.Idx
	p.Idx = 0
	var retErr = m.do_execute(p)
	p.Idx = old_idx
	return retErr
}

// This executes a given code sequence against a given machine
func (m *Machine) do_execute(p *CodeQuotation) (retErr error) {
	var err error
	var wordVal *word
	for err == nil {
		wordVal, err = p.nextWord()
		if err == io.EOF {
			return
		}
		if err != nil {
			return err
		}
		m.NumDispatches++
		var intVal int
		var floatVal float64
		switch {
		case wordVal.impl != nil:
			err = wordVal.impl(m)
		case wordVal.str == "quit":
			return ExitingProgram

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
			wordVal.impl = func(m *Machine) error {
				// fmt.Println("From impl")
				for i := 0; i <= commentWordLen; i++ {
					p.nextWord()
				}
				return nil
			}

		case strings.HasPrefix(wordVal.str, "#"):
			// Skip line comment, potentialy work with it later, but not now.
			continue

		// Both of these functions set the .impl values for the int word in question
		case tryParseInt(wordVal, &intVal):
			m.PushValue(Integer(intVal))
		case tryParseFloat(wordVal, &floatVal):
			m.PushValue(Double(floatVal))

		// This word needs to be defined before we can allow other things to be defined
		case wordVal.str == ":":
			err = m.readWordDefinition(p)
		case wordVal.str == ":PRE":
			err = m.readPrefixDefinition(p)
			// err = m.readPatternDefinition
		case wordVal.str == ":DOC":
			err = m.readWordDocumentation(p)
		case wordVal.str == "error":
			pos := p.getCodePosition()
			str := m.PopValue().String()
			output := fmt.Sprint(pos.Source, " ", pos.LineNumber+1, ":", pos.Offset, ": ", str)
			return fmt.Errorf("%s", output)

		case wordVal.str == "}":
			panic("Unbalanced }!")
		case wordVal.str == "{":
			// fmt.Println("Not cached")
			__quot := &CodeQuotation{
				Words:         make([]*word, 0),
				CodePositions: make([]CodePosition, 0),
			}
			// This is the word that will be patched in holding the Quotation words we're about to patch out.
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
				__quot.CodePositions = append(__quot.CodePositions, p.getCodePosition())
				__quot.Words = append(__quot.Words, wordVal)
			}
			_quot := &Quotation{
				inner:  __quot,
				locals: m.Locals[len(m.Locals)-1]}
			// Run the { } as a Quotation
			_quot.locals = m.Locals[len(m.Locals)-1]
			currIdx := len(m.Values)
			m.PushValue(_quot)
			//fmt.Println(_quot.inner.words)
			err := m.ExecuteQuotation()
			if err != nil {
				return p.wrapError(err)
			}
			vals := make([]StackEntry, len(m.Values)-currIdx)
			copy(vals, m.Values[currIdx:len(m.Values)])
			m.Values = m.Values[:currIdx]
			m.PushValue(&Vector{Elements: vals})

		case wordVal.str == "${":
			// Begin Quotation
			pos := p.getCodePosition()
			//quote := make([]*word, 0)
			quot := &CodeQuotation{
				Words:         make([]*word, 0),
				CodePositions: make([]CodePosition, 0),
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
				quot.CodePositions = append(quot.CodePositions, p.getCodePosition())
				quot.Words = append(quot.Words, wordVal)
			}
			// anchorWord.str = fmt.Sprint("herp", len(__quot.words))
			// Run the { } as a Quotation
			_quot := &Quotation{
				inner:  quot,
				locals: m.Locals[len(m.Locals)-1],
			}
			join, err := stringToQuotation(`"" str-join`, pos)
			if err != nil {
				panic(err)
			}
			_join := &Quotation{
				inner:  join,
				locals: m.Locals[len(m.Locals)-1],
			}
			_quot.locals = m.Locals[len(m.Locals)-1]
			currIdx := len(m.Values)
			// m.PushValue(_quot)
			err = m.execute(_quot.inner)
			if err != nil {
				return err
			}
			vals := make([]StackEntry, len(m.Values)-currIdx)
			copy(vals, m.Values[currIdx:len(m.Values)])
			m.Values = m.Values[:currIdx]
			m.PushValue(&Vector{Elements: vals})
			err = m.execute(_join.inner)

		case wordVal.str == "[":
			// Begin Quotation
			__quot := &CodeQuotation{
				Words:         make([]*word, 0),
				CodePositions: make([]CodePosition, 0),
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
				__quot.CodePositions = append(__quot.CodePositions, p.getCodePosition())
				__quot.Words = append(__quot.Words, wordVal)
			}
			_Quotation := &Quotation{
				inner:  __quot,
				locals: m.Locals[len(m.Locals)-1],
			}
			// _Quotation.locals = m.Locals[len(m.Locals)-1]
			m.PushValue(_Quotation)

		case wordVal.str == "]":
			// panic here.
			panic("Unbalanced ]")
		case strings.HasPrefix(wordVal.str, `"`):
			// This is the case for string literals
			// Slice out the " chracter
			strVal := strings.TrimSuffix(strings.TrimPrefix(wordVal.str, `"`), `"`)

			m.PushValue(String(strVal))
		case wordVal.str == ")":
			panic("Unexpected )")
		case wordVal.str == ";":
			panic("Unexpected ;")
		case wordIsWhitespace(*wordVal):
			wordVal.impl = func(m *Machine) error {
				return nil
			}
			// TODO: capture this space?
			continue
		case len(wordVal.str) == 0:
			wordVal.impl = func(m *Machine) error {
				return nil
			}
			continue
		default:
			if fn, ok := m.PredefinedWords[wordVal.str]; ok {
				if ok {
					err = fn(m)
					if err != nil {
						return err
					}
					wordVal.impl = fn
				}
			} else if val, ok := m.DefinedWords[wordVal.str]; ok {
				// Run the definition of this word on this machine.
				// val.idx = 0
				err = m.execute(val)
				if err != nil {
					return err
				}
				// This is a closure, so I'll need to be careful about it's performance
				wordVal.impl = func(m *Machine) error {
					// val.idx = 0
					return m.execute(val)
				}
			} else if prefixFunc, nonPrefix, ok := m.hasPrefixWord(*wordVal); ok {
				// Put the post-prefix string at the top of the stack, so it can
				// be used.
				m.PushValue(String(nonPrefix))
				err = m.execute(prefixFunc)
				if err != nil {
					return err
				}
				// Captures prefix/nonprefix
				wordVal.impl = func(m *Machine) error {
					m.PushValue(String(nonPrefix))
					return m.execute(prefixFunc)
				}
			} else if err = m.tryLocalWord(wordVal); err == LocalFuncRun || err == WordNotFound {
				if err == WordNotFound {
					return p.wrapError(fmt.Errorf("Undefined word: %v, %v", wordVal, err))
				}
				err = nil
			} else {
				return p.wrapError(err)
			}
		}
		if err != nil {
			// TODO: add some ways to debug here....
			fmt.Fprintln(os.Stderr, "Execution stopped during word:  ", wordVal, " error: ", err)
			return err
		}
	}
	return nil
}

var ErrNoLocals = fmt.Errorf("No locals to try !")
var LocalFuncRun = fmt.Errorf("Nothing was wrong")
var WordNotFound = fmt.Errorf("word was undefined")

func (c *Quotation) execute(m *Machine) error {

	var old_idx = c.inner.Idx
	c.inner.Idx = 0
	m.Locals = append(m.Locals, c.locals)

	var err = m.execute(c.inner)

	m.Locals = m.Locals[:len(m.Locals)-1]
	c.inner.Idx = old_idx
	return err
}

func (m *Machine) tryLocalWord(w *word) error {
	// TODO: In progress
	if len(m.Locals) > 0 {
		if localFunc, found := m.Locals[len(m.Locals)-1][w.str]; found {
			if fn, ok := localFunc.(*Quotation); ok {
				// Push the locals on
				m.Locals = append(m.Locals, fn.locals)
				err := fn.execute(m)
				// Take the local off
				m.Locals = m.Locals[:len(m.Locals)-1]
				if err != nil {
					return err
				}
				return LocalFuncRun
			} else if fn, ok := localFunc.(GoFunc); ok {
				// fmt.Println("HEX")
				err := fn(m)
				if err != nil {
					return err
				}
			} else {
				return WordNotFound
			}
			return LocalFuncRun
		} else {
			return fmt.Errorf("word not found %v", w.str)
		}
	}
	return ErrNoLocals
}

func (m *Machine) readWordBody(name word, c codeSequence) ([]*word, []CodePosition, error) {
	var err = error(nil)
	openPar, err2 := c.nextWord()
	if openPar.str != "(" {
		fmt.Fprintln(os.Stderr, "ERRR0!")
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
	m.DefinedStackComments[name.str] = strings.TrimSpace(stackComment)

	wordDef := make([]*word, 0)
	wordInfo := make([]CodePosition, 0)
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
		wordInfo = append(wordInfo, c.getCodePosition())
	}
	if hasLocal == true {
		wordDef = append([]*word{&word{str: "get-locals"}}, wordDef...)
		wordInfo = append([]CodePosition{CodePosition{Source: "Hardcoded"}}, wordInfo...)
		wordDef = append(wordDef, &word{str: "drop-locals"})
		wordInfo = append(wordInfo, CodePosition{Source: "Hardcoded"})
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
func (m *Machine) readWordDocumentation(c codeSequence) error {
	word, err := c.nextWord()
	if err != nil {
		return err
	}
	if _, found := m.PrefixWords[word.str]; !found {
	} else if _, found := m.PredefinedWords[word.str]; !found {
	} else if _, found := m.DefinedWords[word.str]; !found {
		return fmt.Errorf("No definition for word: %s", word)
	}
	// TODO: Make this it's own loop
	wordDef, _, err := m.readWordBody(*word, c)
	// Save the docs here
	m.HelpDocs[word.str] = stringFromWordDef(wordDef)
	return err
}

// Prefix (:PRE) definitions, which use a prefix
func (m *Machine) readPrefixDefinition(c codeSequence) error {
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
	m.PrefixWords[prefix.str] = &CodeQuotation{
		Idx:           0,
		Words:         wordDef,
		CodePositions: positions,
	}
	return nil
}

// Used for : defined words
func (m *Machine) readWordDefinition(c codeSequence) error {
	name, err := c.nextWord()
	wordDef, positions, err := m.readWordBody(*name, c)
	if err != nil {
		fmt.Println(c.getCodePosition())
		return err
	}
	m.DefinedWords[name.str] = &CodeQuotation{
		Idx:           0,
		Words:         wordDef,
		CodePositions: positions,
	}
	return nil
}

// Used for call word.
func (m *Machine) ExecuteQuotation() error {
	quoteVal := m.PopValue()
	if q, ok := quoteVal.(*Quotation); ok {
		// q.inner = q.inner.cloneCode()
		m.Locals = append(m.Locals, q.locals)
		err := m.execute(q.toCode())
		m.Locals = m.Locals[:len(m.Locals)-1]
		// Works if the err is nil or
		return err
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

func (m *Machine) runConditionalOperator() error {
	falseVal := m.PopValue()
	trueVal := m.PopValue()
	booleanVal := m.PopValue()
	if b, ok := booleanVal.(Boolean); ok {
		if b {
			m.PushValue(trueVal)
		} else {
			m.PushValue(falseVal)
		}
		return nil
	} else {
		// Return the stack to it's previous state, for debugging...?
		m.PushValue(booleanVal)
		m.PushValue(trueVal)
		m.PushValue(falseVal)
		return ConditionalTypeError
	}
}
