package pisc

import (
	"fmt"
	"io"
	"unicode/utf8"
)

/*
The plan is to hav the Source field indicate
where this code came from. Potential values are currently
stdin:
file:
db:
*/
type CodePosition struct {
	LineNumber int
	Offset     int
	Source     string
}

type CodeList struct {
	CodePosition
	Idx  int
	Code string
	// This is for error handling purposes only, indicating where a given word defintion starts
	FileName string
}

func (c CodeList) getCodePosition() CodePosition {
	return c.CodePosition
}

func (c CodeList) wrapError(e error) error {
	return fmt.Errorf("%v\n in File: %v Line: %v column %v", e.Error(), c.Source, c.LineNumber, c.Offset)
}

func (c CodeList) cloneCode() codeSequence {
	return &CodeList{
		Idx:          0,
		Code:         c.Code,
		CodePosition: c.CodePosition,
	}
}

func stringToQuotation(code string, pos CodePosition) (*CodeQuotation, error) {
	basis := &CodeList{
		Idx:          0,
		Code:         code,
		CodePosition: pos,
	}
	quot := &CodeQuotation{
		Idx:           0,
		Words:         make([]*word, 0),
		CodePositions: make([]CodePosition, 0),
	}

	var err error
	var _word *word
	for err == nil {
		_word, err = basis.nextWord()
		if err == io.EOF {
			return quot, nil
		}
		if err != nil {
			return nil, err
		}
		quot.Words = append(quot.Words, _word)
		quot.CodePositions = append(quot.CodePositions, basis.getCodePosition())
	}
	return quot, nil
}

func (c *CodeList) nextWord() (*word, error) {
	currentWord := ""
	skipChar := false
	inString := false
	inLineComment := false
	currentLine := ""
	if c.Idx >= len(c.Code) {
		return &word{str: ""}, io.EOF
	}
	for _, v := range c.Code[c.Idx:] {
		width := utf8.RuneLen(v)
		c.Idx += width
		c.Offset += width
		if v == '\n' {
			// fmt.Println("Parsing:", currentLine)
			currentLine = ""
			c.LineNumber++
			c.Offset = 0
		}
		currentLine += string(v)
		if inLineComment {
			switch v {
			case '\n':
				fallthrough
			case '\r':
				return &word{str: currentWord}, nil
			default:
				currentWord += string(v)
				continue
			}
		}

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
				return nil, fmt.Errorf("Invalid escape sequence: %v current word: %v line: %v",
					v, currentWord, c.LineNumber)
			}
		}

		switch v {
		case '\\':
			if inString {
				skipChar = true
			}
			continue
		case '#':
			if !inString {
				inLineComment = true
			}
			currentWord += string(v)
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
			} else if len(currentWord) > 0 {
				return &word{str: currentWord}, nil
			} else {
				// Skip leading whitespace
				continue
			}
		default:
			currentWord += string(v)
		}
	}
	if inString {
		return nil, fmt.Errorf("Unterminated string!")
	}
	return &word{str: currentWord}, nil
}

type CodeQuotation struct {
	Idx           int
	Words         []*word
	CodePositions []CodePosition
}

func (c *CodeQuotation) nextWord() (*word, error) {
	if c.Idx >= len(c.Words) {
		return nil, io.EOF
	}
	c.Idx++
	return c.Words[c.Idx-1], nil
}

func (c *CodeQuotation) wrapError(e error) error {
	fmt.Println(c.Words)
	// return e
	return fmt.Errorf("%v\n in %v in quotation starting on Line: %v column %v",
		e.Error(),
		c.CodePositions[c.Idx-1].Source,
		c.CodePositions[c.Idx-1].LineNumber,
		c.CodePositions[c.Idx-1].Offset)
}

func (c *CodeQuotation) getCodePosition() CodePosition {
	return c.CodePositions[c.Idx-1]
}

func (c *CodeQuotation) cloneCode() codeSequence {
	return &CodeQuotation{
		Idx:           0,
		Words:         c.Words,
		CodePositions: c.CodePositions,
	}
}
