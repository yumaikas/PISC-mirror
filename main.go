// Posisition
// Independent
// Source
// Code
package main

import (
	"io"
	// "flag" TODO: Implement flags for file and burst modes
	"fmt"
	"os"
	"strings"

	"log"

	"runtime/pprof"

	"gopkg.in/readline.v1"
	cli "gopkg.in/urfave/cli.v1"
)

// This function starts an interpertor
func main() {
	app := cli.NewApp()
	app.Author = "Andrew Owen, @yumaikas"
	app.Name = "PISC, aka Posisition Independent Source Code"
	app.Usage = "A small stack based scripting langauge built for fun"
	app.Action = handleFlags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "interactive, i",
			Usage: "Run the interactive version of PISC",
		},
		cli.StringFlag{
			Name:  "command, c",
			Usage: "Expressions to run from the command line, before -i, if it exists",
		},
		cli.BoolFlag{
			Name:   "benchmark",
			Hidden: true,
			Usage:  "Run various benchmarks, using pprof, and print out pertinent information",
		},
	}
	app.Run(os.Args)
}

func initMachine() *machine {
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[word]codeSequence),
		definedStackComments: make(map[word]string),
		predefinedWords:      make(map[word]GoWord),
		prefixWords:          make(map[word]codeSequence),
		helpDocs:             make(map[word]string),
	}
	m.loadPredefinedValues()
	return m
}

func handleFlags(ctx *cli.Context) {
	m := initMachine()
	// Execute this before benchmarking since we aren't yet benchmarking file loads
	m.executeString(`"factorial.pisc" import`)
	if ctx.IsSet("benchmark") {
		f, err := os.Create("bench-cpu-recursion.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
		}
		m.executeString("100000 [ 12 factorial drop ] times")
		pprof.StopCPUProfile()
		f, err = os.Create("bench-cpu-iteration.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
		}
		m.executeString("100000 [ 12 factorial-loop drop ] times")
		defer pprof.StopCPUProfile()
		return
	}
	if ctx.IsSet("command") {
		line := ctx.String("command")
		p := &codeList{
			idx:  0,
			code: line,
			codePosition: codePosition{
				source: fmt.Sprint("args:"),
			},
		}
		m.execute(p)
	}
	if ctx.IsSet("interactive") {
		loadInteractive(m)
	}
}

func loadInteractive(m *machine) {

	// given_files := flag.Bool("f", false, "Sets the rest of the arguments to list of files")
	// Run command stuff here.

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
		err = m.execute(p)
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
