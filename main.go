// Posisition
// Independent
// Source
// Code
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// This function starts an interpertor
func main() {
	// Run command stuff here.
	in := bufio.NewReader(os.Stdin)
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word][]word),
		definedStackComments: make(map[word]string),
	}

	fmt.Print(
		`Postion
Independent
Source
Code
`)
	for {
		fmt.Print(">> ")
		line, err := in.ReadString('\n')
		if line == "exit" {
			fmt.Println("Exiting")
			return
		}
		if err != nil {
			panic(err)
		}
		words, spaces := getWordList(strings.TrimSpace(line))
		p := &codeList{
			idx:    0,
			code:   words,
			spaces: spaces,
		}
		executeWordsOnMachine(m, p)
		fmt.Println("Data Stack:")
		for _, val := range m.values {
			fmt.Println(val.String())
		}
	}
}
