// Posisition
// Independent
// Source
// Code
package main

import (
	"io"
	// "flag" TODO: Implement flags for file and burst modes
	"fmt"
	"gopkg.in/readline.v1"
	"os"
	"strings"
)

// This function starts an interpertor
func main() {
	// given_files := flag.Bool("f", false, "Sets the rest of the arguments to list of files")
	// Run command stuff here.
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word]codeSequence),
		definedStackComments: make(map[word]string),
		predefinedWords:      make(map[word]GoWord),
		prefixWords:          make(map[word]codeSequence),
		helpDocs:             make(map[word]string),
	}
	m.loadPredefinedValues()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Stdout:          os.Stderr,
		Stderr:          os.Stderr,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Fprintln(
		os.Stderr,
		`Postion
Independent
Source
Code`)
	numEntries := 0
	for {
		// fmt.Print(">> ")
		line, err := rl.Readline()
		if strings.TrimSpace(line) == "exit" {
			fmt.Fprintln(os.Stderr, "Exiting")
			return
		}
		if strings.TrimSpace(line) == "preload" {
			m.loadPredefinedValues()
		}
		if err == io.EOF {
			fmt.Fprintln(os.Stderr, "Exiting program")
			return
		}
		if err != nil {
			panic(err)
		}
		numEntries++
		// fmt.Println(words)
		p := &codeList{
			idx:  0,
			code: line,
			codePosition: codePosition{
				source: fmt.Sprint("stdin:", numEntries),
			},
		}
		err = executeWordsOnMachine(m, p)
		if err == ExitingProgram {
			fmt.Fprintln(os.Stderr, "Exiting program")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:")
			fmt.Fprintln(os.Stderr, err.Error())
		}
		fmt.Fprintln(os.Stderr, "Data Stack:")
		for _, val := range m.values {
			fmt.Fprintln(os.Stderr, val.String(), fmt.Sprint("<", val.Type(), ">"))
		}
	}
}
