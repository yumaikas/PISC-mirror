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
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Execute a file as a bit of pisc",
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
		definedWords:         make(map[string]*codeQuotation),
		definedStackComments: make(map[string]string),
		predefinedWords:      make(map[string]GoWord),
		prefixWords:          make(map[string]*codeQuotation),
		helpDocs:             make(map[string]string),
	}
	m.loadPredefinedValues()
	return m
}

func handleFlags(ctx *cli.Context) {
	m := initMachine()
	m.logAndResetDispatchCount(os.Stderr)
	// Execute this before benchmarking since we aren't yet benchmarking file loads
	if ctx.IsSet("benchmark") {
		err := m.executeString(`"factorial.pisc" import`, codePosition{source: "pre-benchmark import"})
		f, err := os.Create("bench-cpu-recursion.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return
		}
		pos := codePosition{source: "Benchmark recursive"}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
		}
		err = m.executeString("100000 [ 12 factorial drop ] times", pos)
		if err != nil {
			log.Fatal("Recursive benchmark failed:", err)
		}
		pprof.StopCPUProfile()
		f, err = os.Create("bench-cpu-iteration.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return
		}
		pos = codePosition{source: "Benchmark loop"}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
			return
		}
		err = m.executeString("100000 [ 12 factorial-loop drop ] times", pos)
		if err != nil {
			log.Fatal("Recursive benchmark failed:", err)
			pprof.StopCPUProfile()
			return
		}
		pprof.StopCPUProfile()
		return
	}
	if ctx.IsSet("file") {
		m.pushValue(String(ctx.String("file")))
		err := m.executeString("import", codePosition{
			source: "argument line",
		})
		if err != nil {
			log.Fatal("Error running file")
		}
		m.logAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("command") {
		line := ctx.String("command")
		p, err := stringToQuotation(line, codePosition{source: "args"})
		if err != nil {
			log.Fatal("Error in command: ", err)
		}
		err = m.execute(p)
		if err != nil {
			log.Fatal("Error in command: ", err)
		}
		m.logAndResetDispatchCount(os.Stderr)
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

		err = m.executeString(line, codePosition{source: fmt.Sprint("stdin:", numEntries)})
		if err == ExitingProgram {
			fmt.Fprintln(os.Stderr, "Exiting program")
			return
		}
		if err != nil {
			fmt.Println("Error:")
			fmt.Println(err.Error())
			return
		}
		m.logAndResetDispatchCount(os.Stderr)
		fmt.Fprintln(os.Stderr, "Data Stack:")
		for _, val := range m.values {
			fmt.Println(val.String(), fmt.Sprint("<", val.Type(), ">"))
		}
	}

}
