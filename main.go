// Posisition
// Independent
// Source
// Code
package main

import (
	"fmt"
	"gopkg.in/readline.v1"
	"strings"
)

// This function starts an interpertor
func main() {
	// Run command stuff here.
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word]*codeList),
		definedStackComments: make(map[word]string),
		predefinedWords:      make(map[word]GoWord),
		prefixWords:          make(map[word]*codeList),
	}
	m.loadPredefinedValues()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(
		`Postion
Independent
Source
Code
`)
	for {
		fmt.Print(">> ")
		line, err := rl.Readline()
		if strings.TrimSpace(line) == "exit" {
			fmt.Println("Exiting")
			return
		}
		if strings.TrimSpace(line) == "preload" {
			m.loadPredefinedValues()
		}
		if err != nil {
			panic(err)
		}
		words := getWordList(strings.TrimSpace(line))
		// fmt.Println(words)
		p := &codeList{
			idx:  0,
			code: words,
		}
		err = executeWordsOnMachine(m, p)
		if err != nil {
			fmt.Println("Error:")
			fmt.Println(err.Error())
		}
		fmt.Println("Data Stack:")
		for _, val := range m.values {
			fmt.Println(val.String())
		}
	}
}
