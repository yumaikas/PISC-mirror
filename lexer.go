package main

import (
	"fmt"
	"io"
	"unicode/utf8"
)

/*
The plan is to hav the source field indicate
where this code came from. Potential values are currently
stdin:
file:
db:
*/
type codePosition struct {
	lineNumber int
	offset     int
	source     string
}

type codeList struct {
	codePosition
	idx  int
	code string
	// This is for error handling purposes only, indicating where a given word defintion starts
	fileName string
}

func (c codeList) getcodePosition() codePosition {
	return c.codePosition
}

func (c codeList) wrapError(e error) error {
	return fmt.Errorf("%v\n in File: %v Line: %v column %v", e.Error(), c.source, c.lineNumber, c.offset)
}

func (c codeList) cloneCode() codeSequence {
	return &codeList{
		idx:          0,
		code:         c.code,
		codePosition: c.codePosition,
	}
}

func (c *codeList) nextWord() (*word, error) {
	currentWord := ""
	skipChar := false
	inString := false
	inLineComment := false
	currentLine := ""
	if c.idx >= len(c.code) {
		return &word{str: ""}, io.EOF
	}
	for _, v := range c.code[c.idx:] {
		width := utf8.RuneLen(v)
		c.idx += width
		c.offset += width
		if v == '\n' {
			// fmt.Println("Parsing:", currentLine)
			currentLine = ""
			c.lineNumber++
			c.offset = 0
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
				return nil, fmt.Errorf(fmt.Sprint(
					"Invalid escape sequence:", v,
					"current word: ", currentWord,
					"line:", c.lineNumber))
			}
		}

		switch v {
		case '\\':
			if inString {
				skipChar = true
			}
			continue
		case '#':
			inLineComment = true
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

type codeQuotation struct {
	idx   int
	words []*word
	codePosition
}

func (c *codeQuotation) nextWord() (*word, error) {
	if c.idx >= len(c.words) {
		return nil, io.EOF
	}
	c.idx++
	return c.words[c.idx-1], nil
}

func (c codeQuotation) wrapError(e error) error {
	return fmt.Errorf("%v\n in %v in quotation starting on Line: %v column %v", e.Error(), c.source, c.lineNumber, c.offset)
}

func (c codeQuotation) getcodePosition() codePosition {
	return c.codePosition
}

func (c codeQuotation) cloneCode() codeSequence {
	return &codeQuotation{
		idx:          0,
		words:        c.words,
		codePosition: c.codePosition,
	}
}
